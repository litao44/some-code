package main

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/sha3"
	"github.com/ethereum/go-ethereum/rlp"
)

var headerJson = `
{
    "difficulty": "0x2",
    "extraData": "0xd883010812846765746888676f312e31302e34856c696e757800000000000000793ac8714338541c7be17cf5bb7aaa1259e33b19bfdae8b5d1891f426efe331d60da5d7be4b90061db3c94549a28dae7ec3c1f5857b34aaec0aa3a2f4c2ae2ed01",
    "gasLimit": "0x7a1200",
    "gasUsed": "0xdc758",
    "logsBloom": "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
    "miner": "0x0000000000000000000000000000000000000000",
    "mixHash": "0x0000000000000000000000000000000000000000000000000000000000000000",
    "nonce": "0x0000000000000000",
    "number": "0x662c",
    "parentHash": "0xdb0c172429b1250fc450ad6da4f8b7b34858fa4a1f0c894289b7843c2bb18746",
    "receiptsRoot": "0xbc98ce155b5c3d240fc8e4ec19b0c301b6acf1e99c02d83cdada2c4672ef5241",
    "sha3Uncles": "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
    "stateRoot": "0xcecce0065d7bfb2ec6394f92d63d4dcccf5b49f8b64ad98ff88419a9e56a821a",
    "timestamp": "0x5bf3e712",
    "transactionsRoot": "0xb50cbfff9b9919373d82cb258c6d39b50e8ebbf5c5b44a8dbf03eb0f1a281e2e"
}
`

var (
	extraVanity         = 32 // Fixed number of extra-data prefix bytes reserved for signer vanity
	extraSeal           = 65 // Fixed number of extra-data suffix bytes reserved for signer seal
	errMissingSignature = errors.New("extra-data 65 byte signature suffix missing")
)

func main() {
	var header types.Header
	err := json.Unmarshal([]byte(headerJson), &header)
	if err != nil {
		panic(err)
	}

	hash, err := ecrecover(&header)
	if err != nil {
		panic(err)
	}

	fmt.Println(hash.String())
}

// sigHash returns the hash which is used as input for the proof-of-authority
// signing. It is the hash of the entire header apart from the 65 byte signature
// contained at the end of the extra data.
//
// Note, the method requires the extra data to be at least 65 bytes, otherwise it
// panics. This is done to avoid accidentally using both forms (signature present
// or not), which could be abused to produce different hashes for the same header.
func sigHash(header *types.Header) (hash common.Hash) {
	hasher := sha3.NewKeccak256()

	rlp.Encode(hasher, []interface{}{
		header.ParentHash,
		header.UncleHash,
		header.Coinbase,
		header.Root,
		header.TxHash,
		header.ReceiptHash,
		header.Bloom,
		header.Difficulty,
		header.Number,
		header.GasLimit,
		header.GasUsed,
		header.Time,
		header.Extra[:len(header.Extra)-65], // Yes, this will panic if extra is too short
		header.MixDigest,
		header.Nonce,
	})
	hasher.Sum(hash[:0])
	return hash
}

// ecrecover extracts the Ethereum account address from a signed header.
func ecrecover(header *types.Header) (common.Address, error) {
	// Retrieve the signature from the header extra-data
	if len(header.Extra) < extraSeal {
		return common.Address{}, errMissingSignature
	}
	signature := header.Extra[len(header.Extra)-extraSeal:]

	// Recover the public key and the Ethereum address
	pubkey, err := crypto.Ecrecover(sigHash(header).Bytes(), signature)
	if err != nil {
		return common.Address{}, err
	}
	var signer common.Address
	copy(signer[:], crypto.Keccak256(pubkey[1:])[12:])

	return signer, nil
}
