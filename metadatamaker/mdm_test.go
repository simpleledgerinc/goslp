package metadatamaker

import (
	"testing"

	"github.com/simpleledgerinc/goslp/v1parser"
)

func TestCreateOpReturnGenesis(t *testing.T) {
	versionType := 1
	ticker := []byte("TEST")
	name := []byte("some name")
	documentURL := []byte("")
	documentHash := []byte("")
	decimals := 0
	mintBatonVout := &MintBatonVout{vout: 2}
	quantity := uint64(1)

	slpMsg, err := CreateOpReturnGenesis(
		versionType,
		ticker,
		name,
		documentURL,
		documentHash,
		decimals,
		mintBatonVout,
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

func TestCreateOpReturnGenesisBadMintVout(t *testing.T) {
	vouts := []*MintBatonVout{{vout: 0}, {vout: 1}}
	for _, vout := range vouts {
		versionType := 1
		ticker := []byte("TEST")
		name := []byte("some name")
		documentURL := []byte("")
		documentHash := []byte("")
		decimals := 0
		quantity := uint64(1)

		_, err := CreateOpReturnGenesis(
			versionType,
			ticker,
			name,
			documentURL,
			documentHash,
			decimals,
			vout,
			quantity,
		)
		if err.Error() != "mintBatonVout out of range (0x02 < > 0xFF)" {
			t.Fatal(err.Error())
		}
	}
}

func TestCreateOpReturnMintBadMintVout(t *testing.T) {
	vouts := []*MintBatonVout{{vout: 0}, {vout: 1}}
	for _, vout := range vouts {
		versionType := 1
		tokenID := make([]byte, 32)
		quantity := uint64(1)

		_, err := CreateOpReturnMint(
			versionType,
			tokenID,
			vout,
			quantity,
		)
		if err.Error() != "mintBatonVout out of range (0x02 < > 0xFF)" {
			t.Error(err.Error())
		}
	}
}

func TestCreateOpReturnSend(t *testing.T) {
	slpMsg, err := CreateOpReturnSend(
		1,
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
