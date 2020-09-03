package metadatamaker

// TokenType1Genesis creates serialized Genesis op_return
func TokenType1Genesis(
	ticker []byte,
	name []byte,
	documentURL []byte,
	documentHash []byte,
	decimals int,
	mintBatonVout *MintBatonVout,
	quantity uint64,
) ([]byte, error) {
	return CreateOpReturnGenesis(
		0x01,
		ticker,
		name,
		documentURL,
		documentHash,
		decimals,
		mintBatonVout,
		quantity,
	)
}

// TokenType1Mint creates serialized Mint op_return
func TokenType1Mint(
	tokenIDHex []byte,
	mintBatonVout *MintBatonVout,
	quantity uint64) ([]byte, error) {
	return CreateOpReturnMint(0x01, tokenIDHex, mintBatonVout, quantity)
}

// TokenType1Send creates serialized Send op_return
func TokenType1Send(
	tokenIDHex []byte,
	slpAmounts []uint64) ([]byte, error) {
	return CreateOpReturnSend(0x01, tokenIDHex, slpAmounts)
}
