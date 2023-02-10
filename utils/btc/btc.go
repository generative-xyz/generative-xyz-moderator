package btc

import (
	"github.com/incognito-core-libs/btcd/btcec"
	"github.com/incognito-core-libs/btcd/chaincfg"
	"github.com/incognito-core-libs/btcutil"
	"github.com/pkg/errors"
)

// gen a segwit address:
func GenerateAddressSegwit() (privKey, pubKey, addressSegwit string, err error) {

	secret, err := btcec.NewPrivateKey(btcec.S256())
	if err != nil {
		err = errors.Wrap(err, "c.GenerateAddressSegwit")
		return
	}

	wif, err := btcutil.NewWIF(secret, &chaincfg.MainNetParams, true)
	if err != nil {
		err = errors.Wrap(err, "c.GenerateAddressSegwit")
		return
	}

	privKey = wif.String()

	witnessProg := btcutil.Hash160(wif.PrivKey.PubKey().SerializeCompressed())
	addressWitnessPubKeyHash, err := btcutil.NewAddressWitnessPubKeyHash(witnessProg, &chaincfg.MainNetParams)

	if err != nil {
		err = errors.Wrap(err, "btcutil.NewAddressWitnessPubKeyHash")
		return
	}

	addressSegwit = addressWitnessPubKeyHash.EncodeAddress()

	return
}
