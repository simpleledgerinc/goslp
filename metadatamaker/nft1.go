package metadatamaker

func NFT1GroupGenesis(
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
		0x81,
		ticker,
		name,
		documentUrl,
		documentHash,
		decimals,
		mintBatonVout,
		quantity,
	)
}

func NFT1GroupMint(tokenIdHex []byte, mintBatonVout *MintBatonVout, quantity uint64) ([]byte, error) {
	return CreateOpReturnMint(0x81, tokenIdHex, mintBatonVout, quantity)
}

func NFT1GroupSend(tokenIdHex []byte, slpAmounts []uint64) ([]byte, error) {
	return CreateOpReturnSend(0x81, tokenIdHex, slpAmounts)
}

func NFT1ChildGenesis(
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
		0x41,
		ticker,
		name,
		documentUrl,
		documentHash,
		decimals,
		mintBatonVout,
		quantity,
	)
}

func NFT1ChildSend(tokenIdHex []byte, slpAmounts []uint64) ([]byte, error) {
	return CreateOpReturnSend(0x41, tokenIdHex, slpAmounts)
}
