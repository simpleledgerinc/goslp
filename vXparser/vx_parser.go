package vxparser

import (
	"encoding/binary"
	"errors"
)

// TokenType is an uint16 representing the slp version type
type TokenType uint16

const (
	// TokenTypeOpGroup version type used for ParseResult.TokenType
	TokenTypeOpGroup TokenType = 0xff
	// transactionTypeGenesis transaction type used during parsing
	transactionTypeGenesis string = "GENESIS"
)

// AuthorityFlagDesc is a string representation for an authority bit flag
type AuthorityFlagDesc string

// GroupFlagDesc is a string representation of a group flag
type GroupFlagDesc string

const (
	// IsAuthorityFlag flag is used to control whether
	// or not the value associated with the OP_GROUP
	// value in scriptPubKey is a token quantity or
	// authority flag(s)
	//
	// From the spec:
	// This is an authority UTXO, not a “normal”
	// quantity holding UTXO
	//
	IsAuthorityFlag = 1 << 63

	// MintAuthorityStr ...
	MintAuthorityStr = "MintAuthority"

	// MintAuthorityFlag ...
	MintAuthorityFlag = 1 << 62

	// MeltAuthorityStr ...
	MeltAuthorityStr = "MeltAuthority"

	// MeltAuthorityFlag ...
	MeltAuthorityFlag = 1 << 61

	// BatonAuthorityStr ...
	BatonAuthorityStr = "BatonAuthority"

	// BatonAuthorityFlag ...
	BatonAuthorityFlag = 1 << 60

	// RescriptAuthorityStr ...
	RescriptAuthorityStr = "RescriptAuthority"

	// RescriptAuthorityFlag ...
	RescriptAuthorityFlag = 1 << 59

	// SubgroupAuthorityStr ...
	SubgroupAuthorityStr = "SubgroupAuthority"

	// SubgroupAuthorityFlag ...
	SubgroupAuthorityFlag = 1 << 58

	// GroupIsCovenantStr ...
	GroupIsCovenantStr = "Covenant"

	// GroupIsCovenantFlag ...
	GroupIsCovenantFlag = 1

	// GroupIsBchStr ...
	GroupIsBchStr = "HoldsBch"

	// GroupIsBchFlag ...
	GroupIsBchFlag = 2

	// ReservedGroupFlags ...
	ReservedGroupFlags = 0xfffd
)

var (
	// ErrUnsupportedSlpVersion is an error that indicates the parsed slp metadata is
	// an unsupported version
	ErrUnsupportedSlpVersion = errors.New("token_type is not op_group")

	// ActiveAuthorityFlags ...
	ActiveAuthorityFlags = uint64(MintAuthorityFlag |
		MeltAuthorityFlag |
		BatonAuthorityFlag |
		RescriptAuthorityFlag |
		SubgroupAuthorityFlag)

	// AllAuthorityFlags ...
	AllAuthorityFlags = uint64(0xffff) << (64 - 16)

	// ReservedAuthorityFlags ...
	ReservedAuthorityFlags = AllAuthorityFlags ^ uint64(ActiveAuthorityFlags)
)

// ParseResult returns the parsed result.
type ParseResult interface {
	TokenType() TokenType
}

// SlpGenesis is an unmarshalled Genesis ParseResult
type SlpGenesis struct {
	tokenType                               TokenType
	Ticker, Name, DocumentURI, DocumentHash []byte
	Decimals                                int
}

// TokenType returns the TokenType per the ParserResult interface
func (r SlpGenesis) TokenType() TokenType {
	return TokenTypeOpGroup
}

// GroupOutput is an unmarshalled OP_GROUP ParseResult
type GroupOutput struct {
	tokenID          []byte
	quantityOrFlags  uint64
	groupFlags       uint16
	scriptPubKeyTail []byte
}

// TokenType returns the TokenType per the ParserResult interface
func (r GroupOutput) TokenType() TokenType {
	return TokenTypeOpGroup
}

// TokenID returns the TokenID per the ParserResult interface
func (r GroupOutput) TokenID() []byte {
	return r.tokenID
}

// IsAuthority returns a boolean indicating whether or not this
// group output is an authority output
func (r GroupOutput) IsAuthority() bool {
	return r.quantityOrFlags&uint64(IsAuthorityFlag) == 1
}

// Amount of the output
func (r GroupOutput) Amount() (*uint64, error) {
	return nil, errors.New("unimplemented")
}

// IsMintAuthority ...
func (r GroupOutput) IsMintAuthority() (bool, error) {
	if !r.IsAuthority() {
		return false, errors.New("not an authority output")
	}
	return r.quantityOrFlags&MintAuthorityFlag == 1, nil
}

// IsMeltAuthority ...
func (r GroupOutput) IsMeltAuthority() (bool, error) {
	if !r.IsAuthority() {
		return false, errors.New("not an authority output")
	}
	return r.quantityOrFlags&MeltAuthorityFlag == 1, nil
}

// IsSubGroupAuthority ...
func (r GroupOutput) IsSubGroupAuthority() (bool, error) {
	if !r.IsAuthority() {
		return false, errors.New("not an authority output")
	}
	return r.quantityOrFlags&SubgroupAuthorityFlag == 1, nil
}

// IsRescriptAuthority ...
func (r GroupOutput) IsRescriptAuthority() (bool, error) {
	if !r.IsAuthority() {
		return false, errors.New("not an authority output")
	}
	return r.quantityOrFlags&RescriptAuthorityFlag == 1, nil
}

// AuthorityFlags returns authority flags associated with this output. An error
// is returned if this is an amount output.
func (r GroupOutput) AuthorityFlags() ([]AuthorityFlagDesc, error) {
	flags := []AuthorityFlagDesc{}
	if !r.IsAuthority() {
		return nil, errors.New("output is not an authority")
	}

	// check unsupported flags
	if r.quantityOrFlags&ReservedAuthorityFlags > 0 {
		return nil, errors.New("authority output contains un-supported flags")
	}

	// check Mint authority
	if hasFlag, _ := r.IsMintAuthority(); hasFlag {
		flags = append(flags, MintAuthorityStr)
	}

	// check Melt authority
	if hasFlag, _ := r.IsMeltAuthority(); hasFlag {
		flags = append(flags, MeltAuthorityStr)
	}

	// check Subgroup authority
	if hasFlag, _ := r.IsSubGroupAuthority(); hasFlag {
		flags = append(flags, SubgroupAuthorityStr)
	}

	// check Rescript authority
	if hasFlag, _ := r.IsRescriptAuthority(); hasFlag {
		// TODO: check group flags for covenant flag
		flags = append(flags, RescriptAuthorityStr)
	}

	if len(flags) == 0 {
		return nil, errors.New("authority output has no flags set")
	}

	return flags, nil
}

// IsGroupBch ...
func (r GroupOutput) IsGroupBch() bool {
	return r.groupFlags&GroupIsBchFlag == 1
}

// IsGroupCovenant ...
func (r GroupOutput) IsGroupCovenant() bool {
	return r.groupFlags&GroupIsCovenantFlag == 1
}

// GroupFlags returns group flags associated with the token id involved in
// this output.
func (r GroupOutput) GroupFlags() ([]GroupFlagDesc, error) {
	flags := []GroupFlagDesc{}

	if r.groupFlags&ReservedGroupFlags > 0 {
		return nil, errors.New("group token id contains un-supported flags")
	}

	if r.IsGroupBch() {
		flags = append(flags, GroupIsBchStr)
	}

	if r.IsGroupCovenant() {
		flags = append(flags, GroupIsCovenantStr)
	}

	return flags, nil
}

// ParseGroupOutput unmarshals an OP_GROUP output scriptPubKey.
func ParseGroupOutput(scriptPubKey []byte) (ParseResult, error) {
	return nil, errors.New("not implemented")

	// checkValidTokenID := func(tokenID []byte) bool {
	// 	return len(tokenID) == 32
	// }

	// TODO: use bchutil to parse the scriptPubKey and unmarshal into GroupOutput
}

// ParseSLP unmarshals an SLP message from an output scriptPubKey.
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
		return ParseGroupOutput(scriptPubKey)
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

	if tokenType != TokenTypeOpGroup {
		return nil, ErrUnsupportedSlpVersion
	}

	if err := checkNext(); err != nil {
		return nil, err
	}

	transactionType := string(itObj)

	switch transactionType {
	case transactionTypeGenesis:
		if len(chunks) != 8 {
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

		return &SlpGenesis{
			tokenType:    tokenType,
			Ticker:       ticker,
			Name:         name,
			DocumentURI:  documentURI,
			DocumentHash: documentHash,
			Decimals:     decimals,
		}, nil
	}

	return nil, errors.New("unrecognized transaction type")
}

const (
	OP_0         = 0x00 // 0
	OP_PUSHDATA1 = 0x4c // 76
	OP_PUSHDATA2 = 0x4d // 77
	OP_PUSHDATA4 = 0x4e // 78
	OP_RETURN    = 0x6a // 106
)
