package goslp

import (
	"errors"
	"fmt"

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
	tokenType := uint8(slpMsg.TokenType())
	return &tokenType, nil
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

	switch msg := slpMsg.(type) {
	case v1parser.SlpGenesis:
		hash := tx.TxHash()
		var tokenID []byte
		// reverse the bytes here since tokenID is coming from txn hash
		for i := len(hash[:]) - 1; i >= 0; i-- {
			tokenID = append(tokenID, hash[i])
		}
		return tokenID, nil
	case v1parser.SlpMint:
		return msg.TokenID(), nil
	case v1parser.SlpSend:
		return msg.TokenID(), nil
	default:
		return nil, fmt.Errorf("unknown error has occurred")
	}
}

// ConvertSlpToCashAddress converts an slp formatted address to cash formatted address
func ConvertSlpToCashAddress(addr Address, params *chaincfg.Params) (bchutil.Address, error) {
	var (
		bchAddr bchutil.Address
		err     error
	)
	switch a := addr.(type) {
	case *AddressPubKeyHash:
		hash := a.Hash160()
		bchAddr, err = bchutil.NewAddressPubKeyHash(hash[:], params)
		if err != nil {
			return nil, err
		}
	case *AddressScriptHash:
		hash := a.Hash160()
		bchAddr, err = bchutil.NewAddressScriptHash(hash[:], params)
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("address being converted must be type goslp.AddressPubKeyHash or goslp.AddressScriptHash")
	}
	return bchAddr, nil
}

// ConvertCashToSlpAddress converts a cash formatted address to slp formatted address
func ConvertCashToSlpAddress(addr Address, params *chaincfg.Params) (bchutil.Address, error) {
	var (
		bchAddr bchutil.Address
		err     error
	)
	switch a := addr.(type) {
	case *bchutil.AddressPubKeyHash:
		hash := a.Hash160()
		bchAddr, err = NewAddressPubKeyHash(hash[:], params)
		if err != nil {
			return nil, err
		}
	case *bchutil.AddressScriptHash:
		hash := a.Hash160()
		bchAddr, err = NewAddressScriptHash(hash[:], params)
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("address being converted must be type bchutil.AddressPubKeyHash or bchutil.AddressScriptHash")
	}
	return bchAddr, nil
}

func contains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
