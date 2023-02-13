//go:build integration

package delegate

import (
	"context"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/assert"
)

func mustInitService(t *testing.T) *Service {
	ethClient, err := ethclient.Dial("https://goerli.infura.io/v3/98d0983dd54d4076b370506bc98e3485")
	assert.NoError(t, err)
	svc, err := NewService(ethClient)
	return svc
}

func TestService_GetDelegatesForAll(t *testing.T) {
	type args struct {
		vaultAddress string
	}
	tests := []struct {
		name        string
		args        args
		expectedLen int
		wantErr     bool
	}{
		{
			name: "test delegate for all return value",
			args: args{
				vaultAddress: "0xc4b1effe79ed9b9d9eba5d85efb875b188693fcc",
			},
			expectedLen: 1,
			wantErr:     false,
		},
		{
			name: "test delegate for all not return value",
			args: args{
				vaultAddress: "0x9153f96f654b5a912f047016ac4b9c356ca62072",
			},
			expectedLen: 0,
			wantErr:     false,
		},
	}

	ctx, cancelFn := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFn()

	svc := mustInitService(t)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			delegations, err := svc.GetDelegationsByDelegate(ctx, tt.args.vaultAddress)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetDelegatesForAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equalf(t, tt.expectedLen, len(delegations), tt.name)
		})
	}
}
