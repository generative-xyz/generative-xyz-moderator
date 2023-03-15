package btc

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"testing"

	"github.com/btcsuite/btcd/btcutil/psbt"
	"github.com/btcsuite/btcd/chaincfg"
)

func TestCreateZeroValueOutputs(t *testing.T) {
	utxo := []UTXOType{
		{
			Value:      72260,
			TxHash:     "765d9649127338ee6bbde5ca4ac932a7cecb32476cd35cf52156a4b9d758028d",
			TxOutIndex: 1,
		},
		// {
		// 	Value: 0,
		// 	Hash:  "52c60a810fcb18cbfa8030c2b9ecc02ddc4487cb1291fe690785ed8f13d13d3c",
		// 	Index: 0,
		// },
	}

	privateKey := ""
	address := "tb1q0qjpqrgz54xsseymjrql6xs7p9qm0uj53v6rw9"

	BTCChainConf = &chaincfg.TestNet3Params

	// txId, txHex, fee, err := CreateZeroValueOutputs(privateKey, address, 2, utxo, 10)
	// fmt.Println("txId, txHex, fee, err: ", txId, txHex, fee, err)

	txId, txHex, fee, err := CreateTx(privateKey, address, utxo, []PaymentInfo{{Address: address, Amount: 0}}, 5)
	fmt.Println("txId, txHex, fee, err: ", txId, txHex, fee, err)
}

func TestCreateTxBuy(t *testing.T) {
	utxo := []UTXOType{
		{
			Value:      10000,
			TxHash:     "b32a81a0469758e53322ae6db20f86794c0741194d75da99e124dd84e9fd4d52",
			TxOutIndex: 0,
		},
	}

	privateKey := ""
	address := "bc1qn74ftxrvh862jcre972ulnvmve9ek50ewngwyx"
	receiverInsc := "bc1pj2t2szx6rqzcyv63t3xepgdnhuj2zd3kfggrqmd9qwlg3vsx37fqywwhyx"

	// sellerPSBT := "cHNidP8BALICAAAAAskA0un4AlG92RovYKKKNN4SPfXO4jfI1D/YNcAtUtsbAAAAAAD/////yQDS6fgCUb3ZGi9gooo03hI99c7iN8jUP9g1wC1S2xsBAAAAAP////8CMwgAAAAAAAAiUSAbHPlP4PCuwkZG5/JCgkbZ7BEpyHDarPp4zWlO70j4TOgDAAAAAAAAIlEgDBzmg1gYh2X/tSP4NO8AyefnLXBhmmiW1R8Y17zmONgAAAAAAAEBK9IEAAAAAAAAIlEgGxz5T+DwrsJGRufyQoJG2ewRKchw2qz6eM1pTu9I+EwBCEMBQewmfEmb4ocyHjnPLaG68IctXH7Bf9Nr8W4OXYaB+We9K0PwehdWbpQG+k+oTHw13w+wNyi8+6c9SxpkkakdsqCDAAEBK+cDAAAAAAAAIlEgGxz5T+DwrsJGRufyQoJG2ewRKchw2qz6eM1pTu9I+EwBCEMBQdRWgt/6a3bS9+M1J7IqtZaEkcgQFjeMGyHBnE1V4Lj46iawHQ/0I/RnSC8sjEFbwOYKONQlnYkL5X5+j2XvvueDAAAA"

	sellerPSBT := "cHNidP8BALICAAAAAg65q1KIyITDG2o3p/bCApKiS147G5Qx9eYngboqpDvjAAAAAAD/////9bxqxBm95ZeJvOcSW10URv2aVjNwJEyEOqb54urf4cQAAAAAAP////8C0QcAAAAAAAAiUSCSlqgI2hgFgjNRXE2QobO/JKE2NkoQMG2lA76IsgaPkugDAAAAAAAAIlEgGxz5T+DwrsJGRufyQoJG2ewRKchw2qz6eM1pTu9I+EwAAAAAAAEBK+IKAAAAAAAAIlEgkpaoCNoYBYIzUVxNkKGzvyShNjZKEDBtpQO+iLIGj5IBCEMBQZAJPfy9FaCUvArehHclAmCFLgeYPXUD+4RGFW++6Fmxw1vpcZtjnrCmiXtBDCxSlQ6+M7LiJqDOQSJfoOICCzmDAAEBK+gDAAAAAAAAIlEgkpaoCNoYBYIzUVxNkKGzvyShNjZKEDBtpQO+iLIGj5IBCEMBQRPmV73v+yKAMRCT/rBqFR0Fxn+/LeQwEnKBdQUZOnIDzBuuWDxyvu/NzgsoeuO196iHjx3engB9nfhUSVoArHaDAAAA"
	price := uint64(2001)
	feeRate := uint(8)
	maxFee := EstimateTxFee(5, 4, feeRate) + EstimateTxFee(1, 2, feeRate)
	fmt.Printf("maxFee: %v\n", maxFee)

	// txId, txHex, fee, err := CreateZeroValueOutputs(privateKey, address, 2, utxo, 10)
	// fmt.Println("txId, txHex, fee, err: ", txId, txHex, fee, err)

	resp, err := CreatePSBTToBuyInscription(sellerPSBT, privateKey, address, receiverInsc, price, utxo, uint64(feeRate), maxFee)
	fmt.Println("resp, err: ", resp, err)

	fmt.Printf("Split tx: %v\n", resp.SplitTxHex)
	fmt.Printf("Split txID: %v\n", resp.SplitTxID)
	fmt.Printf("split tx fee: %v\n", resp.SplitTxFee)
	fmt.Printf("Buy tx: %v\n", resp.TxHex)
	fmt.Printf("Buy txID: %v\n", resp.TxID)
	fmt.Printf("tx fee: %v\n", resp.BuyTxFee)
}

func TestParsePSBTFromBase64(t *testing.T) {
	type args struct {
		data string
	}
	tests := []struct {
		name    string
		args    args
		want    *psbt.Packet
		wantErr bool
	}{
		{name: "asdasd", args: args{
			data: "cHNidP8BAF4CAAAAAV250MCfjiwSJas0NO8cvPiGPTPSWOiWr6N7jzousaeOAAAAAAD/////AegDAAAAAAAAIlEgGxz5T+DwrsJGRufyQoJG2ewRKchw2qz6eM1pTu9I+EwAAAAAAAEBK3YkAAAAAAAAIlEgGxz5T+DwrsJGRufyQoJG2ewRKchw2qz6eM1pTu9I+EwBCEMBQa7Bm4mP/Qo36HlXPe7s88tC+LUzqF4dzIZvsnKhJ9JwjK+PCPKb3pVr5dlDBIQxh0L+PLArhb+nl3ScKJ6qRwaDAAA=",
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParsePSBTFromBase64(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParsePSBTFromBase64() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			gotBytes, _ := json.MarshalIndent(got, "", " ")
			log.Println("got", string(gotBytes))
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParsePSBTFromBase64() = %v, want %v", got, tt.want)
			}
		})
	}
}
