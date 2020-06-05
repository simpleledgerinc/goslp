package metadatamaker

import (
	"bytes"
	"encoding/binary"
	"errors"
)

// mostly used as optional container
type MintBatonVout struct {
	vout int
}

// https://golang.org/ref/spec#Slice_types
// max len is int-1 which is size of the default integer on target build
// TODO find out if we need error handling here?
// we might want to do system sanity check at compile time
// or possibly use something other...
// this is kind of edge case either way
func pushdata(buf []byte) []byte {
	bufLen := len(buf)

	if bufLen == 0 {
		return []byte{0x4C, 0x00}
	} else if bufLen < 0x4E {
		return bytes.Join([][]byte{[]byte{uint8(bufLen)}, buf}, []byte{})
	} else if bufLen < 0xFF {
		return bytes.Join([][]byte{[]byte{0x4C, uint8(bufLen)}, buf}, []byte{})
	} else if bufLen < 0xFFFF {
		tmp := make([]byte, 2)
		binary.LittleEndian.PutUint16(tmp, uint16(bufLen))
		return bytes.Join([][]byte{[]byte{0x4D}, tmp, buf}, []byte{})
	} else if bufLen < 0xFFFFFFFF {
		tmp := make([]byte, 4)
		binary.LittleEndian.PutUint32(tmp, uint32(bufLen))
		return bytes.Join([][]byte{[]byte{0x4E}, tmp, buf}, []byte{})
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

func CreateOpReturnGenesis(
	versionType int,
	ticker []byte,
	name []byte,
	documentUrl []byte,
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
		pushdata(documentUrl),
		pushdata(documentHash),
		pushdata([]byte{uint8(decimals)}),
		pushdata(mintBatonVoutBytes),
		pushdata(makeU64BigEndianBytes(quantity)),
	}, []byte{})

	return buf, nil
}

func CreateOpReturnMint(versionType int, tokenIdHex []byte, mintBatonVout *MintBatonVout, quantity uint64) ([]byte, error) {
	if versionType != 0x01 && versionType != 0x41 && versionType != 0x81 {
		return nil, errors.New("unknown versionType")
	}

	if len(tokenIdHex) != 32 {
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
		[]byte{0x6A}, // OP_RETURN
		pushdata([]byte("SLP\x00")),
		pushdata([]byte{uint8(versionType)}),
		pushdata([]byte("MINT")),
		pushdata(tokenIdHex),
		pushdata(mintBatonVoutBytes),
		pushdata(makeU64BigEndianBytes(quantity)),
	}, []byte{})

	return buf, nil
}

func CreateOpReturnSend(versionType int, tokenIdHex []byte, slpAmounts []uint64) ([]byte, error) {
	if versionType != 0x01 && versionType != 0x41 && versionType != 0x81 {
		return nil, errors.New("unknown versionType")
	}

	if len(tokenIdHex) != 32 {
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
		[]byte{0x6A}, // OP_RETURN
		pushdata([]byte("SLP\x00")),
		pushdata([]byte{uint8(versionType)}),
		pushdata([]byte("SEND")),
		pushdata(tokenIdHex),
		bytes.Join(amountPushdatas, []byte{}),
	}, []byte{})

	return buf, nil
}
