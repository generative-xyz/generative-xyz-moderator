package contractaddress

import (
	"github.com/ethereum/go-ethereum/ethclient"
	"math"
	"testing"
)

func Test_ethereumImpl_GenerateContractAddress(t *testing.T) {
	tests := []struct {
		name    string
		want    string
		wantErr bool
	}{
		{
			name:    "test1",
			want:    "0x0000000000000000000000000000000000000000",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e, err := NewEthereumChain("https://goerli.infura.io/v3/1db8e9d0f92d4c6d88fa7d15da5dcefb")
			if err != nil {
				t.Errorf("NewEthereumChain() error = %v", err)
				return
			}
			got, err := e.GenerateContractAddress()
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateContractAddress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.GetWalletAddress() != tt.want {
				t.Errorf("GenerateContractAddress() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_ethereumImpl_CheckBalance(t *testing.T) {
	type fields struct {
		client *ethclient.Client
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
			name: "test1",
			args: args{
				address: "0x976D5565927cF44Ee19c346F61FCB37238B426D1",
			},
			want:    0.05,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e, err := NewEthereumChain("https://goerli.infura.io/v3/1db8e9d0f92d4c6d88fa7d15da5dcefb")
			if (err != nil) != tt.wantErr {
				t.Errorf("NewEthereumChain() error = %v, wantErr %v", err, tt.wantErr)
				return
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

func Test_ethereumImpl_Transfer(t *testing.T) {
	type fields struct {
		client *ethclient.Client
	}
	type args struct {
		fromPrivateKey string
		toAddress      string
		amount         uint64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "test1",
			args: args{
				fromPrivateKey: "319717766d44c592d9971d8c595c95caea59017755933deb850b37f0b7941dc4",
				toAddress:      "0xb0cda09aBcc2DA7760AE3862a9204401721c9bB1",
				amount:         uint64(0.005 * math.Pow10(18)),
			},
			want:    "0x0000000000000000000000000000000000000000000000000000000000000001",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e, err := NewEthereumChain("https://goerli.infura.io/v3/1db8e9d0f92d4c6d88fa7d15da5dcefb")
			if (err != nil) != tt.wantErr {
				t.Errorf("NewEthereumChain() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			got, err := e.Transfer(nil)
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
