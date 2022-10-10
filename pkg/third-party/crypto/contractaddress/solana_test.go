package contractaddress

import (
	"github.com/gagliardetto/solana-go/rpc"
	"testing"
)

func Test_solonaImpl_GenerateContractAddress(t *testing.T) {
	type fields struct {
		IncognitoProxy string
		ProgramID      string
		RpcClient      *rpc.Client
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{
			name: "test generate contract address",
			fields: fields{
				IncognitoProxy: "http://localhost:8080",
				ProgramID:      "0x0000000000000000000000000000000000000001",
				RpcClient:      &rpc.Client{},
			},
			want:    "http://localhost:8080/0x0000000000000000000000000000000000000001",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := solonaImpl{
				RpcClient: tt.fields.RpcClient,
			}
			got, err := s.GenerateContractAddress()
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateContractAddress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.walletAddress != tt.want {
				t.Errorf("GenerateContractAddress() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_solonaImpl_CheckBalance(t *testing.T) {
	type fields struct {
		IncognitoProxy string
		ProgramID      string
		RpcClient      *rpc.Client
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
			name: "test check balance",
			fields: fields{
				RpcClient: rpc.New(rpc.TestNet_RPC),
			},
			args: args{
				address: "HFEu22SYhrFMVKB8X4nUAzdgaQT3z1KL2rU3iz4whNnj",
			},
			wantErr: false,
			want:    2.0,
		},
		{
			name: "test check balance main net",
			fields: fields{
				RpcClient: rpc.New(rpc.MainNetBeta_RPC),
			},
			args: args{
				address: "HFEu22SYhrFMVKB8X4nUAzdgaQT3z1KL2rU3iz4whNnj",
			},
			wantErr: false,
			want:    0,
		},
		{
			name: "test check balance main net",
			fields: fields{
				RpcClient: rpc.New(rpc.MainNetBeta_RPC),
			},
			args: args{
				address: "BM2RaPWwdvumJ9bBuqe9BsR3NELvPQVnw6FiYco2pNZy",
			},
			wantErr: false,
			want:    0,
		},
		{
			name: "test check balance test net",
			fields: fields{
				RpcClient: rpc.New(rpc.TestNet_RPC),
			},
			args: args{
				address: "BM2RaPWwdvumJ9bBuqe9BsR3NELvPQVnw6FiYco2pNZy",
			},
			wantErr: false,
			want:    0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := solonaImpl{
				RpcClient: tt.fields.RpcClient,
			}
			got, err := s.CheckBalance(tt.args.address)
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

func Test_solonaImpl_Transfer(t *testing.T) {
	type fields struct {
		IncognitoProxy string
		ProgramID      string
		RpcClient      *rpc.Client
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
			name: "test transfer",
			fields: fields{
				RpcClient: rpc.New(rpc.TestNet_RPC),
			},
			args: args{
				fromPrivateKey: "Tj43PV6KS3AWYhGRPoxouxAzEXBysXdEFGchN6ShJXqrdvAgdTmCLmuFSx48G7jBgWiP49Lnh499zQXgz7tXvpP",
				toAddress:      "BM2RaPWwdvumJ9bBuqe9BsR3NELvPQVnw6FiYco2pNZy",
				amount:         1000000000,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &solonaImpl{
				RpcClient: tt.fields.RpcClient,
			}
			got, err := s.Transfer(nil)
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
