package btc

import (
	"encoding/hex"
	"strings"

	"github.com/btcsuite/btcd/btcutil/psbt"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
)

func ParsePSBTFromBase64(data string) (*psbt.Packet, error) {
	psbtTx, err := psbt.NewFromRawBytes(
		strings.NewReader(data), true,
	)
	if err != nil {
		return nil, err
	}
	return psbtTx, nil
}

func ParseTx(data string) (*wire.MsgTx, error) {
	dataBytes, err := hex.DecodeString(data)
	if err != nil {
		return nil, err
	}
	var tx wire.MsgTx
	err = tx.Deserialize(strings.NewReader(string(dataBytes)))
	if err != nil {
		return nil, err
	}
	return &tx, nil
}

func GetAddressFromPKScript(script []byte) (string, error) {
	pk1, err := txscript.ParsePkScript(script)
	if err != nil {
		return "", err
	}
	address, err := pk1.Address(&chaincfg.MainNetParams)
	if err != nil {
		return "", err
	}
	return address.EncodeAddress(), nil
}
