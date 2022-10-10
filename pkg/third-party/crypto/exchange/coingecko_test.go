package exchange

import (
	coingecko "github.com/superoo7/go-gecko/v3"
	"rederinghub.io/pkg/third-party/crypto/constants/cryptocurrency"

	"net/http"
	"testing"
	"time"
)

func Test_coingeckoImpl_GetRate(t *testing.T) {
	httpClient := &http.Client{
		Timeout: time.Second * 10,
	}
	cg := coingecko.NewClient(httpClient)

	type fields struct {
		client *coingecko.Client
	}
	type args struct {
		currency string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    float64
		wantErr bool
	}{
		{
			name: "test_coingeckoImpl_GetRate_success",
			fields: fields{
				client: cg,
			},
			args: args{
				currency: "eth",
			},
			want:    1587.1199951171875,
			wantErr: false,
		},
		{
			name: "test_coingeckoImpl_GetRate_success_usdt",
			fields: fields{
				client: cg,
			},
			args: args{
				currency: cryptocurrency.USDT,
			},
			want:    1.0,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := coingeckoImpl{
				client: tt.fields.client,
			}
			got, err := c.GetRate(tt.args.currency)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetRate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetRate() got = %v, want %v", got, tt.want)
			}
		})
	}
}
