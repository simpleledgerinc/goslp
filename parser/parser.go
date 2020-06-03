package parser

import (
	"encoding/binary"
	"encoding/hex"
)

// SlpGenesis is an unmarshalled Genesis OP_RETURN
type SlpGenesis struct {
	Ticker, Name, DocumentURI, DocumentHash []byte
	Decimals, MintBatonVout                 int
	Qty                                     uint64
}

// TickerAsUtf8 converts ticker field bytes to string using utf8 decoding
func (g *SlpGenesis) TickerAsUtf8() string {
	return string(g.Ticker)
}

// NameAsUtf8 converts name field bytes to string using utf8 decoding
func (g *SlpGenesis) NameAsUtf8() string {
	return string(g.Name)
}

// DocumentURIAsUtf8 converts documentURI field bytes to string using utf8 decoding
func (g *SlpGenesis) DocumentURIAsUtf8() string {
	return string(g.DocumentURI)
}

// DocumentHashAsHex converts documentHash field bytes to string using hexidecimal encoding
func (g *SlpGenesis) DocumentHashAsHex() string {
	return hex.EncodeToString(g.DocumentHash)
}

// SlpMint is an unmarshalled Mint OP_RETURN
type SlpMint struct {
	TokenID       []byte
	MintBatonVout int
	Qty           uint64
}

// TokenIDAsHex converts TokenId field bytes to string using hexidecimal encoding
func (m *SlpMint) TokenIDAsHex() string {
	return hex.EncodeToString(m.TokenID)
}

// SlpSend is an unmarshalled Send OP_RETURN
type SlpSend struct {
	TokenID []byte
	Amounts []uint64
}

// TokenIDAsHex converts TokenId field bytes to string using hexidecimal encoding
func (s *SlpSend) TokenIDAsHex() string {
	return hex.EncodeToString(s.TokenID)
}

// SlpOpReturn represents a generic interface for
// any type of unmarshalled SLP OP_RETURN message
type SlpOpReturn interface {
	// TODO: once tests are added may need to add ToMap to simplify
	//		 interaction with the SLP unit tests
	//ToMap(raw bool) map[string]string
}

// ParseResult ...
type ParseResult struct {
	TokenType       int
	TransactionType string
	Data            SlpOpReturn
}

// parseSLP unmarshalls an SLP message from a transaction scriptPubKey
func parseSLP(scriptPubKey []byte) ParseResult {
	it := 0
	itObj := scriptPubKey

	const OP_0 int = 0x00
	const OP_RETURN int = 0x6a
	const OP_PUSHDATA1 int = 0x4c
	const OP_PUSHDATA2 int = 0x4d
	const OP_PUSHDATA4 int = 0x4e

	PARSE_CHECK := func(v bool, str string) {
		if v {
			panic(str)
		}
	}

	extractU8 := func() int {
		r := uint8(itObj[it : it+1][0])
		it++
		return int(r)
	}

	extractU16 := func(littleEndian bool) int {
		var r uint16
		if littleEndian {
			r = binary.LittleEndian.Uint16(itObj[it : it+2])
		} else {
			r = binary.BigEndian.Uint16(itObj[it : it+2])
		}
		it += 2
		return int(r)
	}

	extractU32 := func(littleEndian bool) int {
		var r uint32
		if littleEndian {
			r = binary.LittleEndian.Uint32(itObj[it : it+4])
		} else {
			r = binary.BigEndian.Uint32(itObj[it : it+4])
		}
		it += 4
		return int(r)
	}

	extractU64 := func(littleEndian bool) int {
		var r uint64
		if littleEndian {
			r = binary.LittleEndian.Uint64(itObj[it : it+8])
		} else {
			r = binary.BigEndian.Uint64(itObj[it : it+8])
		}
		return int(r)
	}

	PARSE_CHECK(len(itObj) == 0, "scriptpubkey cannot be empty")
	PARSE_CHECK(int(itObj[it]) != OP_RETURN, "scriptpubkey not op_return")
	PARSE_CHECK(len(itObj) < 10, "scriptpubkey too small")
	it++

	extractPushdata := func() int {
		if it == len(itObj) {
			return -1
		}
		cnt := extractU8()
		if cnt > OP_0 && cnt < OP_PUSHDATA1 {
			if it+cnt > len(itObj) {
				it--
				return -1
			}
			return cnt
		} else if cnt == OP_PUSHDATA1 {
			if it+1 >= len(itObj) {
				it--
				return -1
			}
			return extractU8()
		} else if cnt == OP_PUSHDATA2 {
			if it+2 >= len(itObj) {
				it--
				return -1
			}
			return extractU16(true)
		} else if cnt == OP_PUSHDATA4 {
			if it+4 >= len(itObj) {
				it--
				return -1
			}
			return extractU32(true)
		}
		// other opcodes not allowed
		it--
		return -1
	}

	bufferToBN := func() int {
		if len(itObj) == 1 {
			return extractU8()
		}
		if len(itObj) == 2 {
			return extractU16(false)
		}
		if len(itObj) == 4 {
			return extractU32(false)
		}
		if len(itObj) == 8 {
			return extractU64(false)
		}
		panic("extraction of number from buffer failed")
	}

	checkValidTokenID := func(tokenID []byte) bool {
		return len(tokenID) == 32
	}

	chunks := make([][]byte, 0)
	for _len := extractPushdata(); _len >= 0; _len = extractPushdata() {
		buf := make([]byte, _len)
		copy(buf, itObj[it:it+_len])
		PARSE_CHECK(it+_len > len(itObj), "pushdata data extraction failed")
		it += _len
		chunks = append(chunks, buf)
		if len(chunks) == 1 {
			lokadID := chunks[0]
			PARSE_CHECK(len(lokadID) != 4, "lokad id wrong size")
			PARSE_CHECK(
				string(lokadID[0]) != "S" ||
					string(lokadID[1]) != "L" ||
					string(lokadID[2]) != "P" ||
					lokadID[3] != 0x00, "SLP not in first chunk",
			)
		}
	}

	PARSE_CHECK(it != len(itObj), "trailing data")
	PARSE_CHECK(len(chunks) == 0, "chunks empty")

	cit := 0
	CHECK_NEXT := func() {
		cit++
		PARSE_CHECK(cit == len(chunks), "parsing ended early")
		it = 0
		itObj = chunks[cit]
	}
	CHECK_NEXT()

	tokenTypeBuf := itObj
	PARSE_CHECK(len(tokenTypeBuf) != 1 && len(tokenTypeBuf) != 2,
		"token_type string length must be 1 or 2")
	tokenType := bufferToBN()

	PARSE_CHECK(tokenType != 0x01 &&
		tokenType != 0x41 &&
		tokenType != 0x81,
		"token_type not token-type1, nft1-group, or nft1-child")
	CHECK_NEXT()

	transactionType := string(itObj)
	if transactionType == "GENESIS" {
		PARSE_CHECK(len(chunks) != 10, "wrong number of chunks")
		CHECK_NEXT()

		ticker := itObj
		CHECK_NEXT()

		name := itObj
		CHECK_NEXT()

		documentURI := itObj
		CHECK_NEXT()

		documentHash := itObj
		PARSE_CHECK(len(documentHash) != 0 && len(documentHash) != 32, "documentHash must be size 0 or 32")
		CHECK_NEXT()

		decimalsBuf := itObj
		PARSE_CHECK(len(decimalsBuf) != 1, "decimals string length must be 1")
		CHECK_NEXT()

		decimals := bufferToBN()
		PARSE_CHECK(decimals > 9, "decimals biger than 9")
		CHECK_NEXT()

		mintBatonVoutBuf := itObj
		mintBatonVout := 0
		PARSE_CHECK(len(mintBatonVoutBuf) >= 2, "mintBatonVout string must be 0 or 1")
		if len(mintBatonVoutBuf) > 0 {
			mintBatonVout = bufferToBN()
			PARSE_CHECK(mintBatonVout < 2, "mintBatonVout must be at least 2")
		}
		CHECK_NEXT()

		qtyBuf := itObj
		PARSE_CHECK(len(qtyBuf) != 8, "initialQty Must be provided as an 8-byte buffer")
		qty := bufferToBN()

		if tokenType == 0x41 {
			PARSE_CHECK(decimals != 0, "NFT1 child token must have divisibility set to 0 decimal places")
			PARSE_CHECK(mintBatonVout != 0, "NFT1 child token must not have a minting baton")
			PARSE_CHECK(qty != 1, "NFT1 child token must have quantity of 1")
		}

		return ParseResult{
			TokenType:       tokenType,
			TransactionType: transactionType,
			Data: SlpGenesis{
				Ticker:        ticker,
				Name:          name,
				DocumentURI:   documentURI,
				DocumentHash:  documentHash,
				Decimals:      decimals,
				MintBatonVout: mintBatonVout,
				Qty:           uint64(qty),
			},
		}
	} else if transactionType == "MINT" {
		PARSE_CHECK(tokenType == 0x41, "NFT1 Child cannot have MINT transaction type.")

		PARSE_CHECK(len(chunks) != 6, "wrong number of chunks")
		CHECK_NEXT()

		tokenID := itObj
		PARSE_CHECK(!checkValidTokenID(tokenID), "tokenID invalid size")
		CHECK_NEXT()

		mintBatonVoutBuf := itObj
		mintBatonVout := 0
		PARSE_CHECK(len(mintBatonVoutBuf) >= 2, "mint_baton_vout string length must be 0 or 1")
		if len(mintBatonVoutBuf) > 0 {
			mintBatonVout = bufferToBN()
			PARSE_CHECK(mintBatonVout < 2, "mint_baton_vout must be at least 2")
		}
		CHECK_NEXT()

		addiitionalQtyBuf := itObj
		PARSE_CHECK(len(addiitionalQtyBuf) != 8, "additional_qty must be provided as an 8-byte buffer")
		qty := bufferToBN()

		return ParseResult{
			TokenType:       tokenType,
			TransactionType: transactionType,
			Data: SlpMint{
				TokenID:       tokenID,
				MintBatonVout: mintBatonVout,
				Qty:           uint64(qty),
			},
		}
	} else if transactionType == "SEND" {
		PARSE_CHECK(len(chunks) < 4, "wrong number of chunks")
		CHECK_NEXT()

		tokenID := itObj
		PARSE_CHECK(!checkValidTokenID(tokenID), "tokenId invalid size")
		CHECK_NEXT()

		amounts := make([]uint64, 0)
		for cit != len(chunks) {
			amountBuf := itObj
			PARSE_CHECK(len(amountBuf) != 8, "amount string size not 8 bytes")

			value := uint64(bufferToBN())
			amounts = append(amounts, value)

			cit++
			if cit < len(chunks) {
				itObj = chunks[cit]
			}
			it = 0
		}

		PARSE_CHECK(len(amounts) == 0, "token_amounts size is 0")
		PARSE_CHECK(len(amounts) > 19, "token_amounts size is greater than 19")

		return ParseResult{
			TokenType:       tokenType,
			TransactionType: transactionType,
			Data: SlpSend{
				TokenID: tokenID,
				Amounts: amounts,
			},
		}
	}

	panic("impossible parsing result")
}
