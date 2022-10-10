package web3

import (
	"reflect"
	"testing"

	"github.com/laziercoder/go-web3/eth"
	"rederinghub.io/pkg/third-party/crypto/web3/nftdata"
)

func Test_clientImpl_BalanceOf(t *testing.T) {
	type fields struct {
		contract *eth.Contract
	}
	type args struct {
		address string
		chainID int
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
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewClientERC721("https://matic-mumbai.chainstacklabs.com", "0xae0c96bbd7733a1c7843af27e0683c74e182a3a7")
			got, err := c.BalanceOf(&BalanceRequest{
				Address: tt.args.address,
				ChainID: tt.args.chainID,
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
