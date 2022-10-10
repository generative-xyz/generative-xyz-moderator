package contractaddress

import (
	"context"

	"github.com/laziercoder/go-web3/eth"
	"rederinghub.io/pkg/third-party/crypto/constants/cryptocurrency"

	"reflect"
	"testing"
)

func Test_ethereumContractImpl_CheckBalance(t *testing.T) {
	type fields struct {
		contractAddress   string
		ethNativeStrategy Strategy
		coin              string
		providerURL       string
	}
	type args struct {
		address string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    float64
		wantErr bool
	}{
		{
			name: "test usdc",
			fields: fields{
				contractAddress: "0x07865c6E87B9F70255377e024ace6630C1Eaa37F",
				coin:            cryptocurrency.USDC,
				providerURL:     "https://goerli.infura.io/v3/1db8e9d0f92d4c6d88fa7d15da5dcefb",
			},
			args: args{
				address: "0x976D5565927cF44Ee19c346F61FCB37238B426D1",
			},
			want:    2,
			wantErr: false,
		},
		{
			name: "test usdt",
			fields: fields{
				contractAddress: "0x575537A9e553aD212ea46A20529fBe4326719a92",
				coin:            cryptocurrency.USDT,
				providerURL:     "https://kovan.infura.io/v3/1db8e9d0f92d4c6d88fa7d15da5dcefb",
			},
			args: args{
				address: "0x976D5565927cF44Ee19c346F61FCB37238B426D1",
			},
			want:    10,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e, err := NewEthereumContractClient(
				&TokenClientConfig{
					ContractAddress: tt.fields.contractAddress,
					Coin:            tt.fields.coin,
					ProviderURL:     tt.fields.providerURL,
				},
			)
			if err != nil {
				t.Error(err)
			}
			got, err := e.CheckBalance(tt.args.address)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckBalance() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CheckBalance() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_ethereumContractImpl_Transfer(t *testing.T) {
	fromAddress := NewContractAddress().
		WithAddress("0x976D5565927cF44Ee19c346F61FCB37238B426D1").
		WithPrivateKey("319717766d44c592d9971d8c595c95caea59017755933deb850b37f0b7941dc4")
	transferReq := &TransferRequest{
		FromAddress: fromAddress,
		Payer:       fromAddress,
		ToAddress:   "0xC0b7520e4D80d2d2D95D09381422E118f1f57FEb",
		Amount:      5000000,
	}
	type fields struct {
		contract        *eth.Contract
		providerURL     string
		coin            string
		contractAddress string
	}
	type args struct {
		req *TransferRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "test transfer usdc",
			fields: fields{
				contractAddress: "0x07865c6E87B9F70255377e024ace6630C1Eaa37F",
				coin:            cryptocurrency.USDC,
				providerURL:     "https://goerli.infura.io/v3/1db8e9d0f92d4c6d88fa7d15da5dcefb",
			},
			args: args{
				req: transferReq,
			},
			want:    "0x0000000000000000000000000000000000000000000000000000000000000001",
			wantErr: false,
		},
		{
			name: "test transfer usdt",
			fields: fields{
				contractAddress: "0x575537A9e553aD212ea46A20529fBe4326719a92",
				coin:            cryptocurrency.USDT,
				providerURL:     "https://kovan.infura.io/v3/1db8e9d0f92d4c6d88fa7d15da5dcefb",
			},
			args: args{
				req: transferReq,
			},
			want:    "0x0000000000000000000000000000000000000000000000000000000000000001",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e, err := NewEthereumContractClient(&TokenClientConfig{
				ContractAddress: tt.fields.contractAddress,
				Coin:            tt.fields.coin,
				ProviderURL:     tt.fields.providerURL,
			})
			if err != nil {
				t.Error(err)
			}

			gasPrice, err := e.SuggestGasPrice(context.Background())
			if err != nil {
				t.Error(err)
			}
			tt.args.req.GasPrice = gasPrice

			gasLimit, err := e.EstimateGas(tt.args.req)
			if err != nil {
				t.Error(err)
			}
			tt.args.req.GasLimit = gasLimit

			got, err := e.Transfer(tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Transfer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Transfer() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_ethereumContractImpl_GetTransactionReceipt(t *testing.T) {
	type fields struct {
		providerURL     string
		coin            string
		contractAddress string
	}
	type args struct {
		ctx           context.Context
		transactionID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *TransactionReceipt
		wantErr bool
	}{
		{
			name: "test get transaction receipt",
			fields: fields{
				contractAddress: "0x07865c6E87B9F70255377e024ace6630C1Eaa37F",
				coin:            cryptocurrency.USDC,
				providerURL:     "https://goerli.infura.io/v3/1db8e9d0f92d4c6d88fa7d15da5dcefb",
			},
			args: args{
				ctx:           context.Background(),
				transactionID: "0xcb32266da6a1f448cfecbe5449a02871f8c3befd024576dbf769fb11d16d5262",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e, err := NewEthereumContractClient(
				&TokenClientConfig{
					ContractAddress: tt.fields.contractAddress,
					Coin:            tt.fields.coin,
					ProviderURL:     tt.fields.providerURL,
				},
			)
			if err != nil {
				t.Error(err)
			}
			got, err := e.GetTransactionReceipt(tt.args.ctx, tt.args.transactionID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTransactionReceipt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetTransactionReceipt() got = %v, want %v", got, tt.want)
			}
		})
	}
}
