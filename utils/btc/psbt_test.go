package btc

import (
	"fmt"
	"testing"
)

func TestCreateZeroValueOutputs(t *testing.T) {
	utxo := []QuickNodeUTXO{
		{
			Value: 73030,
			Hash:  "52c60a810fcb18cbfa8030c2b9ecc02ddc4487cb1291fe690785ed8f13d13d3c",
			Index: 2,
		},
		// {
		// 	Value: 0,
		// 	Hash:  "52c60a810fcb18cbfa8030c2b9ecc02ddc4487cb1291fe690785ed8f13d13d3c",
		// 	Index: 0,
		// },
	}

	privateKey := ""
	address := ""

	// txId, txHex, fee, err := CreateZeroValueOutputs(privateKey, address, 2, utxo, 10)
	// fmt.Println("txId, txHex, fee, err: ", txId, txHex, fee, err)

	txId, txHex, fee, err := CreateTx(privateKey, address, utxo, []PaymentInfo{{Address: address, Amount: 0}}, 5)
	fmt.Println("txId, txHex, fee, err: ", txId, txHex, fee, err)
}
