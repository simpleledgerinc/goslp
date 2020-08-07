package v1parser

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"reflect"
	"testing"
)

func TestSlpMessageUnitTests(t *testing.T) {
	resp, err := http.Get("https://raw.githubusercontent.com/simpleledger/slp-unit-test-data/master/script_tests.json")
	if err != nil {
		t.Fatal("cannot download unit tests")
	}
	data, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	type TestCase struct {
		Msg    string
		Script string
		Code   *float64
	}
	var tests []TestCase
	err = json.Unmarshal(data, &tests)
	if err != nil {
		t.Fatal(err.Error())
	}
	for _, test := range tests {
		slpbuf, _ := hex.DecodeString(test.Script)
		_, err := ParseSLP(slpbuf)
		if err != nil {
			if test.Code != nil {
				fmt.Println(test.Msg)
				continue
			}
			t.Fatal("goslp parser did not throw an error")
		}
		fmt.Println(test.Msg)
	}
}

func TestGetOutputAmountSend(t *testing.T) {
	scriptPubKey, _ := hex.DecodeString("6a04534c500001010453454e4420c4b0d62156b3fa5c8f3436079b5394f7edc1bef5dc1cd2f9d0c4d46f82cca47908000000000000000108000000000000000408000000000000005a")
	slpMsg, _ := ParseSLP(scriptPubKey)
	amt, _ := slpMsg.GetVoutAmount(3)
	if amt.Cmp(big.NewInt(90)) != 0 {
		t.Fatal("incorrect amount parsed for index")
	}
}

func TestGetOutputAmountSendNotNegative(t *testing.T) {
	scriptPubKey, _ := hex.DecodeString("6a04534c500001010453454e4420c4b0d62156b3fa5c8f3436079b5394f7edc1bef5dc1cd2f9d0c4d46f82cca47908ffffffffffffffff")
	slpMsg, _ := ParseSLP(scriptPubKey)
	amt, _ := slpMsg.GetVoutAmount(1)
	if amt.Cmp(big.NewInt(0)) < 1 {
		t.Fatal("amount is less than zero")
	}
	if amt.Cmp(new(big.Int).SetUint64(18446744073709551615)) != 0 {
		t.Fatal("amount is incorrect")
	}
}

func TestGetOutputAmountMint(t *testing.T) {
	scriptPubKey, _ := hex.DecodeString("6a04534c50000101044d494e5420d6876f0fce603be43f15d34348bb1de1a8d688e1152596543da033a060cff798010208000000000bebc200")
	slpMsg, _ := ParseSLP(scriptPubKey)
	amt, _ := slpMsg.GetVoutAmount(1)
	if amt.Cmp(big.NewInt(200000000)) != 0 {
		t.Fatal("incorrect amount parsed for index")
	}
}

func TestGetOutputAmountMintVout2(t *testing.T) {
	scriptPubKey, _ := hex.DecodeString("6a04534c50000101044d494e5420d6876f0fce603be43f15d34348bb1de1a8d688e1152596543da033a060cff798010208000000000bebc200")
	slpMsg, _ := ParseSLP(scriptPubKey)
	amt, _ := slpMsg.GetVoutAmount(2)
	if amt.Cmp(big.NewInt(0)) != 0 {
		t.Fatal("incorrect amount parsed for index")
	}
}

func TestGetOutputAmountMintVoutNotNegative(t *testing.T) {
	scriptPubKey, _ := hex.DecodeString("6a04534c50000101044d494e5420d6876f0fce603be43f15d34348bb1de1a8d688e1152596543da033a060cff798010208ffffffffffffffff")
	slpMsg, _ := ParseSLP(scriptPubKey)
	amt, _ := slpMsg.GetVoutAmount(1)
	if amt.Cmp(big.NewInt(0)) < 1 {
		t.Fatal("amount is less than zero")
	}
	if amt.Cmp(new(big.Int).SetUint64(18446744073709551615)) != 0 {
		t.Fatal("amount is incorrect")
	}
}

func TestGetOutputAmountGenesis(t *testing.T) {
	scriptPubKey, _ := hex.DecodeString("6a04534c500001010747454e45534953074f6e65436f696e074f6e65436f696e4c5468747470733a2f2f7468656e6578747765622e636f6d2f68617264666f726b2f323031392f31322f32332f6f6e65636f696e2d63727970746f63757272656e63792d7363616d2d6e6565642d746f2d6b6e6f772f4c00010401020800000061c9f36800")
	slpMsg, _ := ParseSLP(scriptPubKey)
	amt, _ := slpMsg.GetVoutAmount(1)
	if amt.Cmp(big.NewInt(420000000000)) != 0 {
		t.Fatal("incorrect amount parsed for index")
	}
}

func TestGetOutputAmountGenesisVout2(t *testing.T) {
	scriptPubKey, _ := hex.DecodeString("6a04534c500001010747454e45534953074f6e65436f696e074f6e65436f696e4c5468747470733a2f2f7468656e6578747765622e636f6d2f68617264666f726b2f323031392f31322f32332f6f6e65636f696e2d63727970746f63757272656e63792d7363616d2d6e6565642d746f2d6b6e6f772f4c00010401020800000061c9f36800")
	slpMsg, _ := ParseSLP(scriptPubKey)
	amt, _ := slpMsg.GetVoutAmount(2)
	if amt.Cmp(big.NewInt(0)) != 0 {
		t.Fatal("incorrect amount parsed for index")
	}
}

func TestGetOutputAmountGenesisVoutNotNegative(t *testing.T) {
	scriptPubKey, _ := hex.DecodeString("6a04534c500001010747454e45534953074f6e65436f696e074f6e65436f696e4c5468747470733a2f2f7468656e6578747765622e636f6d2f68617264666f726b2f323031392f31322f32332f6f6e65636f696e2d63727970746f63757272656e63792d7363616d2d6e6565642d746f2d6b6e6f772f4c000104010208ffffffffffffffff")
	slpMsg, _ := ParseSLP(scriptPubKey)
	amt, _ := slpMsg.GetVoutAmount(1)
	fmt.Println(amt.Text(10))
	if amt.Cmp(big.NewInt(0)) < 1 {
		t.Fatal("amount is less than zero")
	}
	if amt.Cmp(new(big.Int).SetUint64(18446744073709551615)) != 0 {
		t.Fatal("amount is incorrect")
	}
}

func TestTotalOutputAmountSend(t *testing.T) {
	scriptPubKey, _ := hex.DecodeString("6a04534c500001010453454e4420c4b0d62156b3fa5c8f3436079b5394f7edc1bef5dc1cd2f9d0c4d46f82cca47908000000000000000108000000000000000408000000000000005a")
	slpMsg, _ := ParseSLP(scriptPubKey)
	amt, _ := slpMsg.TotalSlpMsgOutputValue()
	if amt.Cmp(big.NewInt(95)) != 0 {
		t.Fatal("incorrect total amount from SEND script")
	}
}

func TestGetTotalOutputAmountSendNotNegative(t *testing.T) {
	scriptPubKey, _ := hex.DecodeString("6a04534c500001010453454e442044e3d05a07091091a63a4074287a784fcd96c26095682e05c22c4bd4e5bf8681080000000000000001080000000000000001080000000000000001080000000000000001080000000000000001080000000000000001080000000000000001080000000000000001088ac7230489e80000")
	slpMsg, _ := ParseSLP(scriptPubKey)
	amt, _ := slpMsg.TotalSlpMsgOutputValue()
	fmt.Println(amt.Text(10))
	if amt.Cmp(big.NewInt(0)) < 1 {
		t.Fatal("amount is less than zero")
	}
}

func TestTotalOutputAmountMint(t *testing.T) {
	scriptPubKey, _ := hex.DecodeString("6a04534c50000101044d494e5420d6876f0fce603be43f15d34348bb1de1a8d688e1152596543da033a060cff798010208000000000bebc200")
	slpMsg, _ := ParseSLP(scriptPubKey)
	amt, _ := slpMsg.TotalSlpMsgOutputValue()
	if amt.Cmp(big.NewInt(200000000)) != 0 {
		t.Fatal("incorrect total amount from script")
	}
}

func TestTotalOutputAmountGenesis(t *testing.T) {
	scriptPubKey, _ := hex.DecodeString("6a04534c500001010747454e45534953074f6e65436f696e074f6e65436f696e4c5468747470733a2f2f7468656e6578747765622e636f6d2f68617264666f726b2f323031392f31322f32332f6f6e65636f696e2d63727970746f63757272656e63792d7363616d2d6e6565642d746f2d6b6e6f772f4c00010401020800000061c9f36800")
	slpMsg, _ := ParseSLP(scriptPubKey)
	amt, _ := slpMsg.TotalSlpMsgOutputValue()
	if amt.Cmp(big.NewInt(420000000000)) != 0 {
		t.Fatal("incorrect total amount from script")
	}
}

func TestGenesisParseSlp(t *testing.T) {
	scriptPubKey, _ := hex.DecodeString("6a04534c500001010747454e455349534c004c004c004c0001004c0008ffffffffffffffff")
	slpMsg, err := ParseSLP(scriptPubKey)
	if err != nil {
		t.Fatal(err.Error())
	}
	if slpMsg.Data.(SlpGenesis).Qty != 18446744073709551615 {
		t.Fatal("incorrect mint qty")
	}
	if slpMsg.TransactionType != "GENESIS" {
		t.Fatal("not genesis")
	}
}

func TestSendParseSlp(t *testing.T) {
	scriptPubKey, _ := hex.DecodeString("6a04534c500001010453454e4420d6876f0fce603be43f15d34348bb1de1a8d688e1152596543da033a060cff798080000000165a0bc00")
	slpMsg, err := ParseSLP(scriptPubKey)
	if err != nil {
		t.Fatal(err.Error())
	}
	if slpMsg.TransactionType != "SEND" {
		t.Fatal("not send")
	}
	if slpMsg.Data.(SlpSend).Amounts[0] != 6000000000 {
		t.Fatal("incorrect send qty")
	}
	tokenID, err := hex.DecodeString("d6876f0fce603be43f15d34348bb1de1a8d688e1152596543da033a060cff798")
	if !reflect.DeepEqual(slpMsg.Data.(SlpSend).TokenID, tokenID) {
		t.Fatal("incorrect tokenID")
	}
}

func TestMintParseSlp(t *testing.T) {
	scriptPubKey, _ := hex.DecodeString("6a04534c50000101044d494e5420d6876f0fce603be43f15d34348bb1de1a8d688e1152596543da033a060cff7980102080000000017d78400")
	slpMsg, err := ParseSLP(scriptPubKey)
	if err != nil {
		t.Fatal(err.Error())
	}
	if slpMsg.TransactionType != "MINT" {
		t.Fatal("not mint")
	}
	if slpMsg.Data.(SlpMint).Qty != 400000000 {
		t.Fatal("incorrect mint qty")
	}
	tokenID, err := hex.DecodeString("d6876f0fce603be43f15d34348bb1de1a8d688e1152596543da033a060cff798")
	if !reflect.DeepEqual(slpMsg.Data.(SlpMint).TokenID, tokenID) {
		t.Fatal("incorrect tokenID")
	}
}
