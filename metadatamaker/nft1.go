package metadatamaker

// NFT1GroupGenesis creates serialized NFT Group genesis op_return message
func NFT1GroupGenesis(
	ticker []byte,
	name []byte,
	documentURL []byte,
	documentHash []byte,
	decimals int,
	mintBatonVout *MintBatonVout,
	quantity uint64,
) ([]byte, error) {
	return CreateOpReturnGenesis(
		0x81,
		ticker,
		name,
		documentURL,
		documentHash,
		decimals,
		mintBatonVout,
		quantity,
	)
}

// NFT1GroupMint creates serialized Mint op_return message
func NFT1GroupMint(
	tokenIDHex []byte,
	mintBatonVout *MintBatonVout,
	quantity uint64) ([]byte, error) {
	return CreateOpReturnMint(0x81, tokenIDHex, mintBatonVout, quantity)
}

// NFT1GroupSend creates serialized Send op_return message
func NFT1GroupSend(
	tokenIDHex []byte,
	slpAmounts []uint64) ([]byte, error) {
	return CreateOpReturnSend(0x81, tokenIDHex, slpAmounts)
}

// NFT1ChildGenesis creates serialized NFT Genesis op_return message
func NFT1ChildGenesis(
	ticker []byte,
	name []byte,
	documentURL []byte,
	documentHash []byte,
	decimals int,
	quantity uint64,
) ([]byte, error) {
	return CreateOpReturnGenesis(
		0x41,
		ticker,
		name,
		documentURL,
		documentHash,
		decimals,
		nil,
		quantity,
	)
}

// NFT1ChildSend creates serialized Send op_return message
func NFT1ChildSend(
	tokenIDHex []byte,
	slpAmounts []uint64) ([]byte, error) {
	return CreateOpReturnSend(0x41, tokenIDHex, slpAmounts)
}
