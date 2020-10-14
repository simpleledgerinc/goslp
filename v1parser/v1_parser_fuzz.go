package v1parser

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type jsParserResult struct {
	Symbol                string
	Name                  string
	DocumentUri           string
	DocumentSha256        string
	Decimals              float64
	GenesisOrMintQuantity string
	SendOutputs           []string
	TransactionType       string
	BatonVout             *float64
	VersionType           float64
	TokenIdHex            string
	ContainsBaton         bool
}

type jsResp struct {
	Success  bool
	Data     jsParserResult
	ErrorMsg *string
}

// Fuzz implements interface for running github.com/dvyuko/go-fuzz
// For more information see ../fuzz/README.md
func Fuzz(data []byte) int {

	// ignore excessively large packages
	if len(data) > 100000 {
		return -1
	}

	// print input
	fmt.Println(hex.EncodeToString(data))

	// get JS parser result
	resp, err := http.Get("http://localhost:8077/" + hex.EncodeToString(data))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// print the js parsed result
	r, _ := ioutil.ReadAll(resp.Body)

	var j *jsResp
	err = json.Unmarshal(r, &j)
	if err != nil {
		panic(err)
	}

	// get goslp.v1parser parser result
	slpMsg, err := ParseSLP(data)

	// check parsing errors throws in both parser implementations.
	if err != nil {
		if !j.Success {
			return -1
		}
		panic(err.Error())
	} else if j.Success == false {
		panic("js parser or json unmarshaling failed, but goslp did not return an error")
	}

	// check matching values on each parsed property
	switch msg := slpMsg.(type) {
	case *SlpGenesis:
		if msg.Decimals != int(j.Data.Decimals) {
			panic("decimals does not match")
		}
		if j.Data.BatonVout != nil {
			if msg.MintBatonVout != int(*j.Data.BatonVout) {
				panic("mint baton does not match")
			}
		}
		jsQty, err := strconv.ParseUint(j.Data.GenesisOrMintQuantity, 0, 64)
		if err != nil {
			panic("bad uint64 value in slp-validate")
		}
		if msg.Qty != jsQty {
			panic("genesis qty does not match")
		}
	case *SlpSend:
		if j.Data.TransactionType != transactionTypeSend {
			panic("transaction type does not match")
		}
		if len(j.Data.SendOutputs) == 0 || len(msg.Amounts) == 0 {
			panic("cannot have 0 outputs in send")
		}
		if strings.ToLower(j.Data.TokenIdHex) != hex.EncodeToString(msg.TokenID()) {
			panic("token ID mismatch")
		}
		for i, amt := range msg.Amounts {
			jsQty, err := strconv.ParseUint(j.Data.SendOutputs[i+1], 10, 64)
			if err != nil {
				panic("bad uint64 send value in slp-validate")
			}
			if amt != jsQty {
				fmt.Println(amt)
				fmt.Println(jsQty)
				panic("send qty does not match")
			}
		}
	case *SlpMint:
		if j.Data.TransactionType != transactionTypeMint {
			panic("transaction type does not match")
		}
		if strings.ToLower(j.Data.TokenIdHex) != hex.EncodeToString(msg.TokenID()) {
			panic("token ID mismatch")
		}
		jsQty, err := strconv.ParseUint(j.Data.GenesisOrMintQuantity, 0, 64)
		if err != nil {
			panic("bad uint64 value in slp-validate")
		}
		if msg.Qty != jsQty {
			panic("genesis qty does not match")
		}
		if j.Data.BatonVout != nil {
			if msg.MintBatonVout != int(*j.Data.BatonVout) {
				panic("mint baton does not match")
			}
		}
	default:
		panic("parser did not throw an error but did not match any parse result")
	}

	// parsing results are lexically correct and was parsed successfully
	return 1
}
