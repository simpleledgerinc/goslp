package metadatamaker

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
)

// MintBatonVout used so that vout value can be set as nil
type MintBatonVout struct {
	vout int
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

	return encodeSlpScript([][]byte{
		[]byte{uint8(versionType)},
		[]byte("GENESIS"),
		ticker,
		name,
		documentURL,
		documentHash,
		[]byte{uint8(decimals)},
		mintBatonVoutBytes,
		makeU64BigEndianBytes(quantity),
	})
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

	return encodeSlpScript([][]byte{
		[]byte{uint8(versionType)},
		[]byte("MINT"),
		tokenIDHex,
		mintBatonVoutBytes,
		makeU64BigEndianBytes(quantity),
	})
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

	amounts := make([][]byte, len(slpAmounts))
	for i, v := range slpAmounts {
		amt := makeU64BigEndianBytes(v)
		amounts[i] = amt
	}

	chunks := make([][]byte, 3+len(amounts))
	chunks[0] = []byte{uint8(versionType)}
	chunks[1] = []byte("SEND")
	chunks[2] = tokenIDHex
	for i, amt := range amounts {
		chunks[i+3] = amt
	}

	return encodeSlpScript(chunks)
}

func pushSlpData(buf []byte) ([]byte, error) {
	bufLen := len(buf)

	if bufLen == 0 {
		return []byte{0x4C, 0x00}, nil
	} else if bufLen < 0x4E {
		return bytes.Join([][]byte{{uint8(bufLen)}, buf}, []byte{}), nil
	} else if bufLen < 0xFF {
		return bytes.Join([][]byte{{0x4C, uint8(bufLen)}, buf}, []byte{}), nil
	} else if bufLen < 0xFFFF {
		tmp := make([]byte, 2)
		binary.LittleEndian.PutUint16(tmp, uint16(bufLen))
		return bytes.Join([][]byte{{0x4D}, tmp, buf}, []byte{}), nil
	} else if bufLen < 0xFFFFFFFF {
		tmp := make([]byte, 4)
		binary.LittleEndian.PutUint32(tmp, uint32(bufLen))
		return bytes.Join([][]byte{{0x4E}, tmp, buf}, []byte{}), nil
	} else {
		return nil, fmt.Errorf("pushSlpData cannot support more than 0xFFFFFFFF elements")
	}
}

func makeU64BigEndianBytes(v uint64) []byte {
	tmp := make([]byte, 8)
	binary.BigEndian.PutUint64(tmp, v)
	return tmp
}

func encodeSlpScript(chunks [][]byte) ([]byte, error) {
	encoded := make([][]byte, len(chunks)+2)
	encoded[0] = []byte{0x6A} // OP_RETURN
	encoded[1] = []byte("\x04SLP\x00")
	for i, chunk := range chunks {
		pushChunk, err := pushSlpData(chunk)
		if err != nil {
			return nil, err
		}
		encoded[i+2] = pushChunk
	}
	return bytes.Join(encoded, []byte{}), nil
}
