package btc

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"testing"

	"github.com/btcsuite/btcd/btcutil/psbt"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/wire"
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

func TestParseTx(t *testing.T) {
	type args struct {
		data string
	}
	tests := []struct {
		name    string
		args    args
		want    *wire.MsgTx
		wantErr bool
	}{
		{name: "asdasd", args: args{data: "010000000001038bd0d0c54e8543f718c1ee6ee7db26ae3ce22fd6be1deec57c66d37b6c97b4db0200000000fdffffff9be83d568d617ea931857f0f90065137f6a7e67a20cd0c80f3957484762e16810000000000fdffffffe07074e1bd2cc630886cf71eaf1f5977b030d32cd953919be2bd33a6d12184340000000000fdffffff03f84d000000000000225120c920e06060005c98739fa4ea58e9fd1859e6affef1b3edbef65257175fa780af9525000000000000225120c920e06060005c98739fa4ea58e9fd1859e6affef1b3edbef65257175fa780af5344000000000000225120c920e06060005c98739fa4ea58e9fd1859e6affef1b3edbef65257175fa780af01406edcb25b1b166731d294a6a26074640271758d20ee9e1aeea86820816102445e79855458c4f4142c7cd01297a78bb6f03a0d62cd8a7a5f4a7845f1363244282d0140e78694fab5996c184b62d481cef8b0fc2bb730c052d7cae860e90219f6610ffd93cf321dd9ad76c7304d3f70e15a5dbaf6f29f44928d941889e48b124cfee16a01408193373f41253fa5d37ffd0c592c9394ffc7f8fc49bff39317e1e829900d240bb496644d0fc05aa884e40093cfb28ebfb54b6d09733a474e3a0ebd179d64f69e00000000"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseTx(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseTx() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			result, err := json.MarshalIndent(got, "", " ")
			log.Println("result", string(result))
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseTx() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParsePSBTFromHex(t *testing.T) {
	type args struct {
		data string
	}
	tests := []struct {
		name    string
		args    args
		want    *psbt.Packet
		wantErr bool
	}{
		{name: "sdfsdf", args: args{data: "70736274ff0100fd060101000000038bd0d0c54e8543f718c1ee6ee7db26ae3ce22fd6be1deec57c66d37b6c97b4db0200000000fdffffff9be83d568d617ea931857f0f90065137f6a7e67a20cd0c80f3957484762e16810000000000fdffffffe07074e1bd2cc630886cf71eaf1f5977b030d32cd953919be2bd33a6d12184340000000000fdffffff03f84d000000000000225120c920e06060005c98739fa4ea58e9fd1859e6affef1b3edbef65257175fa780af9525000000000000225120c920e06060005c98739fa4ea58e9fd1859e6affef1b3edbef65257175fa780af5344000000000000225120c920e06060005c98739fa4ea58e9fd1859e6affef1b3edbef65257175fa780af000000000001012bfa11000000000000225120c920e06060005c98739fa4ea58e9fd1859e6affef1b3edbef65257175fa780af2116c6b1f8a432d0bf501bce5d9ec28767ab9f60c4a0de1abcd1c8bbaeba4269d9ee05007ad34a0c011720c6b1f8a432d0bf501bce5d9ec28767ab9f60c4a0de1abcd1c8bbaeba4269d9ee0001012bf401000000000000225120c920e06060005c98739fa4ea58e9fd1859e6affef1b3edbef65257175fa780af2116c6b1f8a432d0bf501bce5d9ec28767ab9f60c4a0de1abcd1c8bbaeba4269d9ee05007ad34a0c011720c6b1f8a432d0bf501bce5d9ec28767ab9f60c4a0de1abcd1c8bbaeba4269d9ee0001012b50c3000000000000225120c920e06060005c98739fa4ea58e9fd1859e6affef1b3edbef65257175fa780af2116c6b1f8a432d0bf501bce5d9ec28767ab9f60c4a0de1abcd1c8bbaeba4269d9ee05007ad34a0c011720c6b1f8a432d0bf501bce5d9ec28767ab9f60c4a0de1abcd1c8bbaeba4269d9ee00010520c6b1f8a432d0bf501bce5d9ec28767ab9f60c4a0de1abcd1c8bbaeba4269d9ee2107c6b1f8a432d0bf501bce5d9ec28767ab9f60c4a0de1abcd1c8bbaeba4269d9ee05007ad34a0c00010520c6b1f8a432d0bf501bce5d9ec28767ab9f60c4a0de1abcd1c8bbaeba4269d9ee2107c6b1f8a432d0bf501bce5d9ec28767ab9f60c4a0de1abcd1c8bbaeba4269d9ee05007ad34a0c00010520c6b1f8a432d0bf501bce5d9ec28767ab9f60c4a0de1abcd1c8bbaeba4269d9ee2107c6b1f8a432d0bf501bce5d9ec28767ab9f60c4a0de1abcd1c8bbaeba4269d9ee05007ad34a0c00"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParsePSBTFromHex(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParsePSBTFromHex() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// b64, _ := got.B64Encode()
			// log.Printf("result %v\n", b64)
			// got2, err := ParsePSBTFromBase64(b64)
			// if (err != nil) != tt.wantErr {
			// 	t.Errorf("ParsePSBTFromHex() error = %v, wantErr %v", err, tt.wantErr)
			// 	return
			// }

			result, err := json.MarshalIndent(got, "", " ")
			log.Println("result", string(result))
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParsePSBTFromHex() = %v, want %v", got, tt.want)
			}
		})
	}
}
