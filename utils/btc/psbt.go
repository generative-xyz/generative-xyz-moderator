package btc

import (
	"strings"

	"github.com/btcsuite/btcd/btcutil/psbt"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
)

func ParsePSBTFromBase64(data string) (*psbt.Packet, error) {
	psbtTx, err := psbt.NewFromRawBytes(
		strings.NewReader(data), true,
	)
	if err != nil {
		return nil, err
	}
	// pk1, err := txscript.ParsePkScript(testPsbt.UnsignedTx.TxOut[0].PkScript)
	// if err != nil {
	// 	return err
	// }
	// address, err := pk1.Address(&chaincfg.MainNetParams)
	// if err != nil {
	// 	return err
	// }

	// testPsbtBytes, _ := json.Marshal(testPsbt)
	// fmt.Println(address.EncodeAddress())
	return psbtTx, nil
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
