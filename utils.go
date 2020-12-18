package goslp

import (
	"errors"
	"fmt"

	"github.com/gcash/bchd/wire"

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
	case *v1parser.SlpGenesis:
		hash := tx.TxHash()
		var tokenID []byte
		// reverse the bytes here since tokenID is coming from txn hash
		for i := len(hash[:]) - 1; i >= 0; i-- {
			tokenID = append(tokenID, hash[i])
		}
		return tokenID, nil
	case *v1parser.SlpMint:
		return msg.TokenID(), nil
	case *v1parser.SlpSend:
		return msg.TokenID(), nil
	default:
		return nil, fmt.Errorf("unknown error has occurred")
	}
}

func contains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
