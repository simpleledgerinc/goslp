package v1parser

import (
	"encoding/hex"
	"math/big"
	"reflect"
	"testing"
)

func TestGetOutputAmountSend(t *testing.T) {
	scriptPubKey, _ := hex.DecodeString("6a04534c500001010453454e4420c4b0d62156b3fa5c8f3436079b5394f7edc1bef5dc1cd2f9d0c4d46f82cca47908000000000000000108000000000000000408000000000000005a")
	slpMsg, _ := ParseSLP(scriptPubKey)
	amt, _ := slpMsg.GetVoutAmount(3)
	if amt.Cmp(big.NewInt(90)) != 0 {
		t.Fatal("incorrect amount parsed for index")
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
