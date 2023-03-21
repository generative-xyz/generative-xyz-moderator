package usecase

import (
	"fmt"

	"github.com/btcsuite/btcd/btcutil/psbt"
	"github.com/btcsuite/btcd/wire"
	"rederinghub.io/utils/btc"
)

func retrievePreviousTxFromPSBT(psbtData *psbt.Packet) ([]string, error) {
	result := []string{}
	for _, input := range psbtData.UnsignedTx.TxIn {
		result = append(result, input.PreviousOutPoint.Hash.String())
	}
	return result, nil
}

func retrieveAllVinFromPSBT(psbtData *psbt.Packet) ([]string, error) {
	result := []string{}
	for _, input := range psbtData.UnsignedTx.TxIn {
		result = append(result, fmt.Sprintf("%v:%v", input.PreviousOutPoint.Hash.String(), input.PreviousOutPoint.Index))
	}
	return result, nil
}

func extractAllOutputFromPSBT(psbtData *psbt.Packet) (map[string][]*wire.TxOut, error) {
	result := make(map[string][]*wire.TxOut)
	for _, output := range psbtData.UnsignedTx.TxOut {
		address, err := btc.GetAddressFromPKScript(output.PkScript)
		if err != nil {
			return nil, err
		}
		result[address] = append(result[address], output)
	}
	return result, nil
}
