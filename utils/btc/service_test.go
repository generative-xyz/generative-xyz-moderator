package btc

import (
	"fmt"
	"testing"
)

func TestFilterPendingUTXOs(t *testing.T) {

	utxo := []UTXOType{
		{
			Value:      4105,
			TxHash:     "52c60a810fcb18cbfa8030c2b9ecc02ddc4487cb1291fe690785ed8f13d13d3c",
			TxOutIndex: 2,
		},
		{
			Value:      999,
			TxHash:     "52c60a810fcb18cbfa8030c2b9ecc02ddc4487cb1291fe690785ed8f13d13d3c",
			TxOutIndex: 0,
		},
	}
	address := "bc1qn74ftxrvh862jcre972ulnvmve9ek50ewngwyx"
	_, _, err := FilterPendingUTXOs(utxo, address)
	fmt.Println(err)
}
