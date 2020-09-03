# goslp

Golang packages for the Simple Ledger Protocol (SLP).

### v1parser - for parsing transaction metadata

This package is used for parsing SLP metadata from the SLP transaction's input 0 scriptPubKey.

```go
// here is an example marshaled SLP message extracted from transaction's scriptPubKey.
scriptPubKey := hex.DecodeString("6a04534c500001010453454e4420c4b0d62156b3fa5c8f3436079b5394f7edc1bef5dc1cd2f9d0c4d46f82cca47908000000000000000108000000000000000408000000000000005a")

// get the unmarshaled slp message
slpMsg, err := v1parser.ParseSLP(scriptPubKey)

// do something ...

```

This usage, [here](https://github.com/simpleledgerinc/bchd/blob/slp-index/bchrpc/server.go#L1240), in BCHD gRPC server provides a good example usage of how to interact with the unmarshaled SLP metadata object.

Differential fuzzer testing has been performed with the [slp-validate.js](https://github.com/simpleledger/slp-validate) npm package, and can be reproduced following the instructions in the `./fuzz` directory.

### metadatamaker - for creating new transaction metadata

This package is used for creating marshaled SLP metadata for adding to a transaction's first output scriptPubKey.  Helper methods are provided for Type 1, NFT Group, and NFT children transactions.

**Genesis** - use CreateOpReturnGenesis, NFT1GroupGenesis, or NFT1ChildGenesis

```go
scriptPubKey, err := CreateOpReturnGenesis(
      versionType,
      ticker,
      name,
      documentURL,
      documentHash,
      decimals,
      mintBatonVout,
      quantity,
)

```



**Mint** - use CreateOpReturnMint or NFT1GroupMint

```go
scriptPubKey, err := CreateOpReturnMint(
    versionType,
    tokenID,
    vout,
    quantity,
)

```



**Send** - use CreateOpReturnSend, NFT1GroupSend, or NFT1ChildSend

```go
scriptPubKey, err := CreateOpReturnSend(
		1,
		tokenID,
		[]uint64{1, 2},
)

```

