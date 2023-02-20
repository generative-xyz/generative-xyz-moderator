package helpers

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/ecdsa"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
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
	// TODO hiennguyen
	return nil, nil

	/*if publicKey == nil {
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

	ecdsaPub := publicKey.ToECDSA()
	//temp, err := json.Marshal(ecdsaPub)
	//if err != nil {
	//	return nil, err
	//}
	result := bec.PublicKey{
		Curve: ecdsaPub.Curve, X: ecdsaPub.X, Y: ecdsaPub.Y,
	}
	err = json.Unmarshal(temp, &result)
	if err != nil {
		return nil, err
	}
	return bscript.NewAddressFromPublicKey(&result, true)*/
}

// PubKeyFromSignature gets a publickey for a signature and tells you whether is was compressed
func PubKeyFromSignature(sig, msg string, prefix string) (pubKey *btcec.PublicKey, wasCompressed bool, err error) {
	// TODO hiennguyen
	var decodedSig []byte
	if decodedSig, err = base64.StdEncoding.DecodeString(sig); err != nil {
		return nil, false, err
	}

	temp, err := MagicHash(msg, prefix)
	if err != nil {
		return nil, false, err
	}
	aaa := []byte{
		89,
		217,
		162,
		121,
		96,
		82,
		197,
		186,
		9,
		206,
		127,
		206,
		247,
		56,
		104,
		197,
		70,
		22,
		253,
		55,
		13,
		145,
		1,
		144,
		223,
		46,
		73,
		71,
		6,
		111,
		174,
		40,
	}
	_ = temp
	k, c, err := ecdsa.RecoverCompact(decodedSig, aaa[:])
	return k, c, err
}

func MagicHash(msg, messagePrefix string) (*chainhash.Hash, error) {
	// TODO hiennguyen
	if messagePrefix == "" {
		messagePrefix = `\u0018Bitcoin Signed Message:\n`
	}

	// Validate the signature - this just shows that it was valid at all
	// we will compare it with the key next
	var buf bytes.Buffer
	if err := wire.WriteVarString(&buf, 0, messagePrefix); err != nil {
		return nil, err
	}
	if err := wire.WriteVarString(&buf, 0, msg); err != nil {
		return nil, err
	}

	// Create the hash
	expectedMessageHash := chainhash.HashH(buf.Bytes())
	return &expectedMessageHash, nil
}
