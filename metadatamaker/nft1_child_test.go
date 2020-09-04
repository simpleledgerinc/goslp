package metadatamaker

import (
	"testing"

	"github.com/simpleledgerinc/goslp/v1parser"
)

func TestNft1ChildCreateOpReturnGenesis(t *testing.T) {
	ticker := []byte("TEST")
	name := []byte("some name")
	documentURL := []byte("")
	documentHash := []byte("")
	decimals := 0
	quantity := uint64(1)

	slpMsg, err := NFT1ChildGenesis(
		ticker,
		name,
		documentURL,
		documentHash,
		decimals,
		quantity,
	)
	if err != nil {
		t.Fatal(err.Error())
	}
	if slpMsg == nil {
		t.Fatal("parameters were not marshalled")
	}
	_, err = v1parser.ParseSLP(slpMsg)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestNft1ChildCreateOpReturnSend(t *testing.T) {
	slpMsg, err := NFT1ChildSend(
		make([]byte, 32),
		[]uint64{0, 1},
	)
	if err != nil {
		t.Error(err.Error())
	}
	_, err = v1parser.ParseSLP(slpMsg)
	if err != nil {
		t.Error(err.Error())
	}
}
