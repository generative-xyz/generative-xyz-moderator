package helpers

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/libsv/go-bk/bec"
	"github.com/libsv/go-bt/v2/bscript"
	"golang.org/x/crypto/ripemd160"
	"hash"
)

// Calculate the hash of hasher over buf.
func calcHash(buf []byte, hasher hash.Hash) []byte {
	hasher.Write(buf)
	return hasher.Sum(nil)
}

// Hash160 calculates the hash ripemd160(sha256(b)).
func Hash160(buf []byte) []byte {
	return calcHash(calcHash(buf, sha256.New()), ripemd160.New())
}

// GetAddressFromPubKey gets a bscript.Address from a bec.PublicKey
func GetAddressFromPubKey(publicKey *btcec.PublicKey, compressed bool) (*bscript.Address, error) {
	if publicKey == nil {
		return nil, fmt.Errorf("publicKey cannot be nil")
	} else if publicKey.X == nil {
		return nil, fmt.Errorf("publicKey.X cannot be nil")
	}

	if !compressed {
		// go-bt/v2/bscript does not have a function that exports the uncompressed address
		// https://github.com/libsv/go-bt/blob/master/bscript/address.go#L98
		hash := Hash160(publicKey.SerializeUncompressed())
		bb := make([]byte, 1)
		//nolint: makezero // we need to set up the array with 1
		bb = append(bb, hash...)
		return &bscript.Address{
			AddressString: bscript.Base58EncodeMissingChecksum(bb),
			PublicKeyHash: hex.EncodeToString(hash),
		}, nil
	}

	temp, err := json.Marshal(publicKey.ToECDSA())
	if err != nil {
		return nil, err
	}
	var result bec.PublicKey
	err = json.Unmarshal(temp, &result)
	if err != nil {
		return nil, err
	}
	return bscript.NewAddressFromPublicKey(&result, true)
}
