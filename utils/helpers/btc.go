package helpers

import (
	"crypto/sha256"
	"encoding/base64"
	"hash"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/ecdsa"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"golang.org/x/crypto/ripemd160"
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
func GetAddressFromPubKey(publicKey *btcec.PublicKey, compressed bool) (*btcutil.AddressPubKeyHash, error) {
	temp, err := btcutil.NewAddressPubKeyHash(btcutil.Hash160(publicKey.SerializeCompressed()), &chaincfg.MainNetParams)
	if err != nil {
		return nil, err
	}
	return temp, nil
}

// PubKeyFromSignature gets a publickey for a signature and tells you whether is was compressed
func PubKeyFromSignature(sig, msg string, prefix string) (pubKey *btcec.PublicKey, wasCompressed bool, err error) {
	var decodedSig []byte
	if decodedSig, err = base64.StdEncoding.DecodeString(sig); err != nil {
		return nil, false, err
	}

	temp, err := MagicHash(msg, prefix)
	if err != nil {
		return nil, false, err
	}
	k, c, err := ecdsa.RecoverCompact(decodedSig, temp[:])
	return k, c, err
}

func MagicHash(msg, messagePrefix string) (chainhash.Hash, error) {
	if messagePrefix == "" {
		messagePrefix = "\u0018Bitcoin Signed Message:\n"
	}

	bytes := append([]byte(messagePrefix), []byte(msg)...)
	return chainhash.DoubleHashH(bytes), nil
}

// func VerifyETHSignature(sig, msg, address string) (bool, error) {

// 	decodedSig, err := base64.StdEncoding.DecodeString(sig)
// 	if err != nil {
// 		return false, err
// 	}
// 	decodedSig, err :=
// 	if err != nil {
// 		return false, err
// 	}

// 	crypto.Ecrecover()

// }
