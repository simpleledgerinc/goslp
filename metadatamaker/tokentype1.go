package metadatamaker

func TokenType1Genesis(
	versionType int,
	ticker []byte,
	name []byte,
	documentUrl []byte,
	documentHash []byte,
	decimals int,
	mintBatonVout *MintBatonVout,
	quantity uint64,
) ([]byte, error) {
	return CreateOpReturnGenesis(
		0x01,
		ticker,
		name,
		documentUrl,
		documentHash,
		decimals,
		mintBatonVout,
		quantity,
	)
}

func TokenType1Mint(tokenIdHex []byte, mintBatonVout *MintBatonVout, quantity uint64) ([]byte, error) {
	return CreateOpReturnMint(0x01, tokenIdHex, mintBatonVout, quantity)
}

func TokenType1Send(tokenIdHex []byte, slpAmounts []uint64) ([]byte, error) {
	return CreateOpReturnSend(0x01, tokenIdHex, slpAmounts)
}
