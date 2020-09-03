package metadatamaker

import (
	"bytes"
	"encoding/binary"
	"errors"
)

// MintBatonVout mostly used as optional container
type MintBatonVout struct {
	vout int
}

// https://golang.org/ref/spec#Slice_types
// max len is int-1 which is size of the default integer on target build
func pushdata(buf []byte) []byte {
	bufLen := len(buf)

	if bufLen == 0 {
		return []byte{0x4C, 0x00}
	} else if bufLen < 0x4E {
		return bytes.Join([][]byte{{uint8(bufLen)}, buf}, []byte{})
	} else if bufLen < 0xFF {
		return bytes.Join([][]byte{{0x4C, uint8(bufLen)}, buf}, []byte{})
	} else if bufLen < 0xFFFF {
		tmp := make([]byte, 2)
		binary.LittleEndian.PutUint16(tmp, uint16(bufLen))
		return bytes.Join([][]byte{{0x4D}, tmp, buf}, []byte{})
	} else if bufLen < 0xFFFFFFFF {
		tmp := make([]byte, 4)
		binary.LittleEndian.PutUint32(tmp, uint32(bufLen))
		return bytes.Join([][]byte{{0x4E}, tmp, buf}, []byte{})
	} else {
		panic("pushdata cannot support more than 0xFFFFFFFF elements")
	}
}

// TODO we can use better name for this
func makeU64BigEndianBytes(v uint64) []byte {
	tmp := make([]byte, 8)
	binary.BigEndian.PutUint64(tmp, v)
	return tmp
}

// CreateOpReturnGenesis creates serialized Genesis op_return message
func CreateOpReturnGenesis(
	versionType int,
	ticker []byte,
	name []byte,
	documentURL []byte,
	documentHash []byte,
	decimals int,
	mintBatonVout *MintBatonVout,
	quantity uint64,
) ([]byte, error) {
	if versionType != 0x01 && versionType != 0x41 && versionType != 0x81 {
		return nil, errors.New("unknown versionType")
	}

	if len(documentHash) != 0 && len(documentHash) != 32 {
		return nil, errors.New("documentHash must be either 0 or 32 hex bytes")
	}
	if decimals < 0 || decimals > 9 {
		return nil, errors.New("decimals out of range")
	}

	if versionType == 0x41 {
		if quantity != 1 {
			return nil, errors.New("quantity must be 1 for NFT1 child genesis")
		}

		if decimals != 0 {
			return nil, errors.New("decimals must be 0 for NFT1 child genesis")
		}

		if mintBatonVout != nil {
			return nil, errors.New("mintBatonVout must be null for NFT1 child genesis")
		}
	}

	mintBatonVoutBytes := []byte{}
	if mintBatonVout != nil {
		if mintBatonVout.vout < 2 || mintBatonVout.vout > 0xFF {
			return nil, errors.New("mintBatonVout out of range (0x02 < > 0xFF)")
		}
		mintBatonVoutBytes = []byte{uint8(mintBatonVout.vout)}
	}

	buf := bytes.Join([][]byte{
		[]byte{0x6A}, // OP_RETURN
		pushdata([]byte("SLP\x00")),
		pushdata([]byte{uint8(versionType)}),
		pushdata([]byte("GENESIS")),
		pushdata(ticker),
		pushdata(name),
		pushdata(documentURL),
		pushdata(documentHash),
		pushdata([]byte{uint8(decimals)}),
		pushdata(mintBatonVoutBytes),
		pushdata(makeU64BigEndianBytes(quantity)),
	}, []byte{})

	return buf, nil
}

// CreateOpReturnMint creates serialized Mint op_return message
func CreateOpReturnMint(
	versionType int,
	tokenIDHex []byte,
	mintBatonVout *MintBatonVout,
	quantity uint64) ([]byte, error) {
	if versionType != 0x01 && versionType != 0x41 && versionType != 0x81 {
		return nil, errors.New("unknown versionType")
	}

	if len(tokenIDHex) != 32 {
		return nil, errors.New("tokenIdHex must be 32 bytes")
	}

	mintBatonVoutBytes := []byte{}
	if mintBatonVout != nil {
		if mintBatonVout.vout < 2 || mintBatonVout.vout > 0xFF {
			return nil, errors.New("mintBatonVout out of range (0x02 < > 0xFF)")
		}
		mintBatonVoutBytes = []byte{uint8(mintBatonVout.vout)}
	}

	buf := bytes.Join([][]byte{
		{0x6A}, // OP_RETURN
		pushdata([]byte("SLP\x00")),
		pushdata([]byte{uint8(versionType)}),
		pushdata([]byte("MINT")),
		pushdata(tokenIDHex),
		pushdata(mintBatonVoutBytes),
		pushdata(makeU64BigEndianBytes(quantity)),
	}, []byte{})

	return buf, nil
}

// CreateOpReturnSend create serialized Send op_return message
func CreateOpReturnSend(
	versionType int,
	tokenIDHex []byte,
	slpAmounts []uint64) ([]byte, error) {
	if versionType != 0x01 && versionType != 0x41 && versionType != 0x81 {
		return nil, errors.New("unknown versionType")
	}

	if len(tokenIDHex) != 32 {
		return nil, errors.New("tokenIdHex must be 32 bytes")
	}

	if len(slpAmounts) < 1 {
		return nil, errors.New("send requires at least one amount")
	}

	if len(slpAmounts) > 19 {
		return nil, errors.New("too many slp amounts")
	}

	amountPushdatas := make([][]byte, len(slpAmounts))
	for i, v := range slpAmounts {
		amountPushdatas[i] = pushdata(makeU64BigEndianBytes(v))
	}

	buf := bytes.Join([][]byte{
		{0x6A}, // OP_RETURN
		pushdata([]byte("SLP\x00")),
		pushdata([]byte{uint8(versionType)}),
		pushdata([]byte("SEND")),
		pushdata(tokenIDHex),
		bytes.Join(amountPushdatas, []byte{}),
	}, []byte{})

	return buf, nil
}
