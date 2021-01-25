package v1parser

import (
	"encoding/binary"
	"errors"
	"math/big"
)

// TokenType is an uint16 representing the slp version type
type TokenType uint16

const (
	// TokenTypeFungible01 version type used for ParseResult.TokenType
	TokenTypeFungible01 TokenType = 0x01
	// TokenTypeNft1Child41 version type used for ParseResult.TokenType
	TokenTypeNft1Child41 TokenType = 0x41
	// TokenTypeNft1Group81 version type used for ParseResult.TokenType
	TokenTypeNft1Group81 TokenType = 0x81
	// TransactionTypeGenesis transaction type used for ParseResult.TransactionType
	transactionTypeGenesis string = "GENESIS"
	// TransactionTypeMint transaction type used for ParseResult.TransactionType
	transactionTypeMint string = "MINT"
	// TransactionTypeSend transaction type used for ParseResult.TransactionType
	transactionTypeSend string = "SEND"
)

// ParseResult returns the parsed result.
type ParseResult interface {
	TokenType() TokenType
	TokenID() []byte
	GetVoutValue(vout int) (*big.Int, bool)
	TotalSlpMsgOutputValue() (*big.Int, error)
}

// SlpGenesis is an unmarshalled Genesis ParseResult
type SlpGenesis struct {
	tokenType                               TokenType
	Ticker, Name, DocumentURI, DocumentHash []byte
	Decimals, MintBatonVout                 int
	Qty                                     uint64
}

// TokenType returns the TokenType per the ParserResult interface
func (r SlpGenesis) TokenType() TokenType {
	return r.tokenType
}

// TokenID returns the TokenID per the ParserResult interface
func (r SlpGenesis) TokenID() []byte {
	return nil
}

// GetVoutValue returns the output amount or boolean flag indicating
// the index is a mint baton for a given transaction output index.
// Out of range vout returns nil.
func (r SlpGenesis) GetVoutValue(vout int) (*big.Int, bool) {

	if r.MintBatonVout == vout {
		return nil, true
	}

	if vout == 0 {
		return nil, false
	}

	if vout == 1 {
		return new(big.Int).SetUint64(r.Qty), false
	}
	return nil, false
}

// TotalSlpMsgOutputValue computes the output amount transferred in a transaction
func (r SlpGenesis) TotalSlpMsgOutputValue() (*big.Int, error) {
	total := big.NewInt(0)
	total.Add(total, new(big.Int).SetUint64(r.Qty))
	return total, nil
}

// SlpMint is an unmarshalled Mint ParseResult
type SlpMint struct {
	tokenID       []byte
	tokenType     TokenType
	MintBatonVout int
	Qty           uint64
}

// TokenType returns the TokenType per the ParserResult interface
func (r SlpMint) TokenType() TokenType {
	return r.tokenType
}

// TokenID returns the TokenID per the ParserResult interface
func (r SlpMint) TokenID() []byte {
	return r.tokenID
}

// GetVoutValue returns the output amount or boolean flag indicating
// the index is a mint baton for a given transaction output index.
// Out of range vout returns nil.
func (r SlpMint) GetVoutValue(vout int) (*big.Int, bool) {

	if r.MintBatonVout == vout {
		return nil, true
	}

	if vout == 0 {
		return nil, false
	}

	if vout == 1 {
		return new(big.Int).SetUint64(r.Qty), false
	}
	return nil, false
}

// TotalSlpMsgOutputValue computes the output amount transferred in a transaction
func (r SlpMint) TotalSlpMsgOutputValue() (*big.Int, error) {
	total := big.NewInt(0)
	total.Add(total, new(big.Int).SetUint64(r.Qty))

	return total, nil
}

// SlpSend is an unmarshalled Send ParseResult
type SlpSend struct {
	tokenID   []byte
	tokenType TokenType
	Amounts   []uint64
}

// TokenType returns the TokenType per the ParserResult interface
func (r SlpSend) TokenType() TokenType {
	return r.tokenType
}

// TokenID returns the TokenID per the ParserResult interface
func (r SlpSend) TokenID() []byte {
	return r.tokenID
}

// GetVoutValue returns the output amount or boolean flag indicating
// the index is a mint baton for a given transaction output index.
// Out of range vout returns nil.
func (r SlpSend) GetVoutValue(vout int) (*big.Int, bool) {
	if vout == 0 {
		return nil, false
	}

	if vout > len(r.Amounts) {
		return nil, false
	}

	return new(big.Int).SetUint64(r.Amounts[vout-1]), false
}

// TotalSlpMsgOutputValue computes the output amount transferred in a transaction
func (r SlpSend) TotalSlpMsgOutputValue() (*big.Int, error) {
	total := big.NewInt(0)
	for _, amt := range r.Amounts {
		total.Add(total, new(big.Int).SetUint64(amt))
	}
	return total, nil
}

// ParseSLP unmarshals an SLP message from a transaction scriptPubKey.
func ParseSLP(scriptPubKey []byte) (ParseResult, error) {
	it := 0
	itObj := scriptPubKey

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
		if littleEndian {
			return int(binary.LittleEndian.Uint64(itObj[it : it+8]))
		}
		return int(binary.BigEndian.Uint64(itObj[it : it+8]))
	}

	if len(itObj) == 0 {
		return nil, errors.New("scriptpubkey cannot be empty")
	}
	if int(itObj[it]) != OP_RETURN {
		return nil, errors.New("scriptpubkey not op_return")
	}
	if len(itObj) < 10 {
		return nil, errors.New("scriptpubkey too small")
	}

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

	bufferToBN := func() (int, error) {
		if len(itObj) == 1 {
			return extractU8(), nil
		}
		if len(itObj) == 2 {
			return extractU16(false), nil
		}
		if len(itObj) == 4 {
			return extractU32(false), nil
		}
		if len(itObj) == 8 {
			return extractU64(false), nil
		}
		return 0, errors.New("extraction of number from buffer failed")
	}

	checkValidTokenID := func(tokenID []byte) bool {
		return len(tokenID) == 32
	}

	chunks := make([][]byte, 0)
	for chunkLen := extractPushdata(); chunkLen >= 0; chunkLen = extractPushdata() {
		if it+chunkLen > len(itObj) {
			return nil, errors.New("pushdata data extraction failed")
		}

		buf := make([]byte, chunkLen)
		copy(buf, itObj[it:it+chunkLen])

		it += chunkLen
		chunks = append(chunks, buf)
		if len(chunks) == 1 {
			bchMetaTag := chunks[0]

			if len(bchMetaTag) != 4 {
				return nil, errors.New("OP_RETURN magic is wrong size")
			}

			if bchMetaTag[0] != 0x53 || bchMetaTag[1] != 0x4c || bchMetaTag[2] != 0x50 || bchMetaTag[3] != 0x00 {
				return nil, errors.New("OP_RETURN magic is not in first chunk")
			}
		}
	}

	if it != len(itObj) {
		return nil, errors.New("trailing data")
	}

	if len(chunks) == 0 {
		return nil, errors.New("chunks empty")
	}

	cit := 0

	checkNext := func() error {
		cit++

		if cit == len(chunks) {
			return errors.New("parsing ended early")
		}

		it = 0
		itObj = chunks[cit]

		return nil
	}

	if err := checkNext(); err != nil {
		return nil, err
	}

	tokenTypeBuf := itObj

	if len(tokenTypeBuf) != 1 && len(tokenTypeBuf) != 2 {
		return nil, errors.New("token_type string length must be 1 or 2")
	}

	tokenTypeInt, err := bufferToBN()
	if err != nil {
		return nil, err
	}

	tokenType := TokenType(tokenTypeInt)

	if tokenType != TokenTypeFungible01 &&
		tokenType != TokenTypeNft1Child41 &&
		tokenType != TokenTypeNft1Group81 {
		return nil, errors.New("token_type not token-type1, nft1-group, or nft1-child")
	}

	if err := checkNext(); err != nil {
		return nil, err
	}

	transactionType := string(itObj)
	if transactionType == transactionTypeGenesis {

		if len(chunks) != 10 {
			return nil, errors.New("wrong number of chunks")
		}

		if err := checkNext(); err != nil {
			return nil, err
		}

		ticker := itObj
		if err := checkNext(); err != nil {
			return nil, err
		}

		name := itObj
		if err := checkNext(); err != nil {
			return nil, err
		}

		documentURI := itObj
		if err := checkNext(); err != nil {
			return nil, err
		}

		documentHash := itObj

		if len(documentHash) != 0 && len(documentHash) != 32 {
			return nil, errors.New("documentHash string length must be 0 or 32")
		}

		if err := checkNext(); err != nil {
			return nil, err
		}

		decimalsBuf := itObj

		if len(decimalsBuf) != 1 {
			return nil, errors.New("decimals string length must be 1")
		}

		decimals, err := bufferToBN()
		if err != nil {
			return nil, err
		}

		if decimals > 9 {
			return nil, errors.New("decimals bigger than 9")
		}

		if err := checkNext(); err != nil {
			return nil, err
		}

		mintBatonVoutBuf := itObj
		mintBatonVout := 0

		if len(mintBatonVoutBuf) >= 2 {
			return nil, errors.New("mintBatonVout string length must be 0 or 1")
		}

		if len(mintBatonVoutBuf) > 0 {
			mintBatonVout, err = bufferToBN()
			if err != nil {
				return nil, err
			}

			if mintBatonVout < 2 {
				return nil, errors.New("mintBatonVout value must be at least 2")
			}
		}

		if err := checkNext(); err != nil {
			return nil, err
		}

		qtyBuf := itObj

		if len(qtyBuf) != 8 {
			return nil, errors.New("initialQty must be provided as an 8-byte buffer")
		}

		qty, err := bufferToBN()
		if err != nil {
			return nil, err
		}

		if tokenType == TokenTypeNft1Child41 {
			if decimals != 0 {
				return nil, errors.New("NFT1 child token must have divisibility set to 0 decimal places")
			}

			if mintBatonVout != 0 {
				return nil, errors.New("NFT1 child token must not have a minting baton")
			}

			if qty != 1 {
				return nil, errors.New("NFT1 child token must have quantity of 1")
			}
		}

		return &SlpGenesis{
			tokenType:     tokenType,
			Ticker:        ticker,
			Name:          name,
			DocumentURI:   documentURI,
			DocumentHash:  documentHash,
			Decimals:      decimals,
			MintBatonVout: mintBatonVout,
			Qty:           uint64(qty),
		}, nil
	} else if transactionType == transactionTypeMint {

		if tokenType == TokenTypeNft1Child41 {
			return nil, errors.New("nft1 child cannot have mint transaction type")
		}

		if len(chunks) != 6 {
			return nil, errors.New("wrong number of chunks")
		}

		if err := checkNext(); err != nil {
			return nil, err
		}

		tokenID := itObj

		if !checkValidTokenID(tokenID) {
			return nil, errors.New("tokenID invalid size")
		}

		if err := checkNext(); err != nil {
			return nil, err
		}

		mintBatonVoutBuf := itObj
		mintBatonVout := 0

		if len(mintBatonVoutBuf) >= 2 {
			return nil, errors.New("mint_baton_vout string length must be 0 or 1")
		}

		if len(mintBatonVoutBuf) > 0 {
			mintBatonVout, err = bufferToBN()
			if err != nil {
				return nil, err
			}

			if mintBatonVout < 2 {
				return nil, errors.New("mint_baton_vout must be at least 2")
			}

		}
		if err := checkNext(); err != nil {
			return nil, err
		}

		additionalQtyBuf := itObj

		if len(additionalQtyBuf) != 8 {
			return nil, errors.New("additional_qty must be provided as an 8-byte buffer")
		}

		qty, err := bufferToBN()
		if err != nil {
			return nil, err
		}

		return &SlpMint{
			tokenType:     tokenType,
			tokenID:       tokenID,
			MintBatonVout: mintBatonVout,
			Qty:           uint64(qty),
		}, nil
	} else if transactionType == transactionTypeSend {

		if len(chunks) < 4 {
			return nil, errors.New("wrong number of chunks")
		}

		if err := checkNext(); err != nil {
			return nil, err
		}

		tokenID := itObj

		if !checkValidTokenID(tokenID) {
			return nil, errors.New("tokenId invalid size")
		}

		if err := checkNext(); err != nil {
			return nil, err
		}

		amounts := make([]uint64, 0)
		for cit != len(chunks) {
			amountBuf := itObj

			if len(amountBuf) != 8 {
				return nil, errors.New("amount string size not 8 bytes")
			}

			value, err := bufferToBN()
			if err != nil {
				return nil, err
			}
			amounts = append(amounts, uint64(value))

			cit++
			if cit < len(chunks) {
				itObj = chunks[cit]
			}
			it = 0
		}

		if len(amounts) == 0 {
			return nil, errors.New("token_amounts size is 0")
		}

		if len(amounts) > 19 {
			return nil, errors.New("token_amounts size is greater than 19")
		}

		return &SlpSend{
			tokenType: tokenType,
			tokenID:   tokenID,
			Amounts:   amounts,
		}, nil
	}

	return nil, errors.New("impossible parsing result")
}

const (
	OP_0         = 0x00 // 0
	OP_PUSHDATA1 = 0x4c // 76
	OP_PUSHDATA2 = 0x4d // 77
	OP_PUSHDATA4 = 0x4e // 78
	OP_RETURN    = 0x6a // 106
)
