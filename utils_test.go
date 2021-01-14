package goslp_test

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/gcash/bchd/wire"
	"github.com/simpleledgerinc/goslp"
)

func TestGetSlpTokenIDGenesis(t *testing.T) {
	txnHex := "0100000001fb6080a6ca752a808e4f86a6225164b3348b10572610d85001d00d1a9b151629030000006441e9e0035a15773bf8f86b65415c4827a9cec018abe15bc162b9974f3001a0a5ff751d35500837b9e8e77d322a8289fa92ff718be4701e5fa38f3e2726c729fb7b412102afef5c197947afa712fd6094935531935d24834cb2f0fefd811691e7230eb82bfeffffff040000000000000000416a04534c500001810747454e45534953034244441b426974636f696e20446f6e6174696f6e73204469726563746f72794c004c000100010208000000000000000122020000000000001976a914294e1c12d3f976f2dd5bd10467c4c605d6996b8e88ac22020000000000001976a914294e1c12d3f976f2dd5bd10467c4c605d6996b8e88ac9d350200000000001976a914e7abe33c8b9d58366b3114a8979509fc80420ad288ac52050a00"

	tx := wire.NewMsgTx(1)
	serializedTx, err := hex.DecodeString(txnHex)
	if err != nil {
		t.Fatal(err.Error())
	}

	err = tx.BchDecode(bytes.NewReader(serializedTx), wire.ProtocolVersion, wire.LatestEncoding)
	if err != nil {
		t.Fatal(err.Error())
	}

	_, err = goslp.GetSlpTokenID(tx)
	if err != nil {
		t.Fatal(err)
	}
}
