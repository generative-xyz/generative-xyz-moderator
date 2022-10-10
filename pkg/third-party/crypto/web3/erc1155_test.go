package web3

import (
	"reflect"
	"testing"

	"github.com/laziercoder/go-web3/eth"
	"rederinghub.io/pkg/third-party/crypto/web3/nftdata"
)

func Test_clientERC1155Impl_BalanceOf(t *testing.T) {
	type fields struct {
		contract *eth.Contract
	}
	type args struct {
		address string
		chainID int
		tokenID int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *nftdata.NFTCustomerInfo
		wantErr bool
	}{
		{
			name: "solana address success",
			fields: fields{
				contract: &eth.Contract{},
			},
			args: args{
				address: "HgYwra6a2e6d3gizk4CdakH9R1uGuxejJMKvJedPPJ7d",
			},
		},
		{
			name: "rinkeby address success",
			fields: fields{
				contract: &eth.Contract{},
			},
			args: args{
				address: "0x976D5565927cF44Ee19c346F61FCB37238B426D1",
				chainID: 1,
				tokenID: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewClientERC1155("https://rinkeby.infura.io/v3/1db8e9d0f92d4c6d88fa7d15da5dcefb", "0x54db584a3fc0d4e78c82f06eacde76f93b28ae15")
			got, err := c.BalanceOf(&BalanceRequest{
				Address: tt.args.address,
				ChainID: tt.args.chainID,
				TokenID: tt.args.tokenID,
			})
			if (err != nil) != tt.wantErr {
				t.Errorf("BalanceOf() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BalanceOf() got = %v, want %v", got, tt.want)
			}
		})
	}
}
