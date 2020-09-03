package v1parser

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type jsData struct {
	symbol                string
	name                  string
	documentURI           string
	documentSha256        string
	decimals              int
	genesisOrMintQuantity string
	sendOutputs           []string
	transactionType       string
	mintBatonVout         *int
}

type jsResp struct {
	success bool
	data    jsData
}

// Fuzz implements interface for github.com/dvyuko/go-fuzz
func Fuzz(data []byte) int {

	// print input
	fmt.Println(hex.EncodeToString(data))

	// get JS parser result
	resp, err := http.Get("http://localhost:8077/" + hex.EncodeToString(data))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var j jsResp
	err = json.NewDecoder(resp.Body).Decode(&j)
	if err != nil {
		panic(err)
	}

	// get goslp parser result
	slpMsg, err := ParseSLP(data)

	// check parsing errors throws in both impls.
	if err != nil {
		if !j.success {
			return -1
		}
		panic(err.Error())
	}

	// check matching values on each parsed property
	switch _type := slpMsg.Data.(type) {
	case *SlpGenesis:
		if j.data.transactionType != "GENESIS" {
			panic("transaction type does not match")
		}
		if string(_type.DocumentURI) != j.data.documentURI {
			panic("document uri string does not match")
		}
		if hex.EncodeToString(_type.DocumentHash) != j.data.documentSha256 {
			panic("document hash hex does not match")
		}
		if string(_type.Name) != j.data.name {
			panic("name string does not match")
		}
		if _type.Decimals != j.data.decimals {
			panic("decimals does not match")
		}
		if j.data.mintBatonVout != nil {
			if _type.MintBatonVout != *j.data.mintBatonVout {
				panic("mint baton does not match")
			}
		}
		jsQty, err := strconv.ParseUint(j.data.genesisOrMintQuantity, 0, 64)
		if err != nil {
			panic("bad uint64 value in slp-validate")
		}
		if _type.Qty != jsQty {
			panic("genesis qty does not match")
		}
	case *SlpSend:
		if j.data.transactionType != "SEND" {
			panic("transaction type does not match")
		}
		if len(j.data.sendOutputs) == 0 || len(_type.Amounts) == 0 {
			panic("cannot have 0 outputs in send")
		}
		for i, amt := range _type.Amounts {
			jsQty, err := strconv.ParseUint(j.data.sendOutputs[i], 0, 64)
			if err != nil {
				panic("bad uint64 send value in slp-validate")
			}
			if amt != jsQty {
				panic("send qty does not match")
			}
		}
	case *SlpMint:
		if j.data.transactionType != "MINT" {
			panic("transaction type does not match")
		}
		jsQty, err := strconv.ParseUint(j.data.genesisOrMintQuantity, 0, 64)
		if err != nil {
			panic("bad uint64 value in slp-validate")
		}
		if _type.Qty != jsQty {
			panic("genesis qty does not match")
		}
		if j.data.mintBatonVout != nil {
			if _type.MintBatonVout != *j.data.mintBatonVout {
				panic("mint baton does not match")
			}
		}
	default:
		return 0
	}

	// parsing results are lexically correct and was parsed successfully
	return 1
}
