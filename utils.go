package goslp

import (
	"errors"

	"github.com/gcash/bchd/chaincfg"
	"github.com/gcash/bchd/wire"
	"github.com/gcash/bchutil"
	"github.com/simpleledgerinc/goslp/v1parser"
)

// GetSlpVersionType returns the SLP version number regardless of version/type
func GetSlpVersionType(slpPkScript []byte) (*uint8, error) {

	// TODO: replace the following with a method to only parse the version number
	slpMsg, err := v1parser.ParseSLP(slpPkScript)
	if err != nil {
		return nil, errors.New("unable to parse slp version")
	}
	_type := uint8(slpMsg.TokenType)
	return &_type, nil
}

// GetSlpTokenID returns the Token ID regardless of SLP version/type
func GetSlpTokenID(tx *wire.MsgTx) ([]byte, error) {

	slpVersion, err := GetSlpVersionType(tx.TxOut[0].PkScript)
	if err != nil {
		return nil, err
	}

	if !contains([]int{0x01, 0x41, 0x81}, int(*slpVersion)) {
		return nil, errors.New("cannot parse token id for an unknown slp version type")
	}

	slpMsg, err := v1parser.ParseSLP(tx.TxOut[0].PkScript)
	if err != nil {
		return nil, err
	}

	if slpMsg.TransactionType == "SEND" {
		return slpMsg.Data.(v1parser.SlpSend).TokenID, nil
	} else if slpMsg.TransactionType == "MINT" {
		return slpMsg.Data.(v1parser.SlpMint).TokenID, nil
	} else if slpMsg.TransactionType == "GENESIS" {
		hash := tx.TxHash()
		var tokenID []byte
		// reverse the bytes here since tokenID is coming from txn hash
		for i := len(hash[:]) - 1; i >= 0; i-- {
			tokenID = append(tokenID, hash[i])
		}
		return tokenID, nil
	} else {
		panic("unknown error has occured")
	}
}

// ConvertSlpToCashAddress converts an slp formatted address to cash formatted address
func ConvertSlpToCashAddress(addr Address, params *chaincfg.Params) bchutil.Address {
	var bchAddr bchutil.Address
	switch a := addr.(type) {
	case *AddressPubKeyHash:
		hash := a.Hash160()
		bchAddr, _ = bchutil.NewAddressPubKeyHash(hash[:], params)
	case *AddressScriptHash:
		hash := a.Hash160()
		bchAddr, _ = bchutil.NewAddressScriptHash(hash[:], params)
	}
	return bchAddr
}

func contains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
