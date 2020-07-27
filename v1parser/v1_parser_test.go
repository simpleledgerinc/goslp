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

func TestTotalOutputAmountSend(t *testing.T) {
	scriptPubKey, _ := hex.DecodeString("6a04534c500001010453454e4420c4b0d62156b3fa5c8f3436079b5394f7edc1bef5dc1cd2f9d0c4d46f82cca47908000000000000000108000000000000000408000000000000005a")
	slpMsg, _ := ParseSLP(scriptPubKey)
	amt, _ := slpMsg.TotalSlpMsgOutputValue()
	if amt.Cmp(big.NewInt(95)) != 0 {
		t.Fatal("incorrect total amount from SEND script")
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
