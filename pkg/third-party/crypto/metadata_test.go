package crypto

import (
	"encoding/json"
	"testing"
	"time"

	"rederinghub.io/pkg/third-party/crypto/constants/cryptocurrency"
	"rederinghub.io/pkg/types"
)

func TestMetadata_WithTotalAmount(t *testing.T) {
	type fields struct {
		Rate           float64
		Currency       string
		TotalAmount    *types.Price
		CurrencySymbol string
		USDAmount      float64
		RefreshTime    time.Time
	}
	type args struct {
		cartTotalAmount float64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Metadata
	}{
		{
			name: "test with total amount",
			fields: fields{
				Rate:           1622.0,
				Currency:       cryptocurrency.Ethereum,
				CurrencySymbol: cryptocurrency.CurrencySymbolByCurrency[cryptocurrency.Ethereum],
			},
			args: args{
				cartTotalAmount: 100.0,
			},
			want: &Metadata{
				Rate:     1622.0,
				Currency: cryptocurrency.Ethereum,
				TotalAmount: types.NewPrice(
					100/1622.0,
					types.WithPrecision(cryptocurrency.DefaultRoundPlaces),
					types.WithCryptoPrecision(cryptocurrency.EthereumRoundPlaces),
				),
				CurrencySymbol: cryptocurrency.CurrencySymbolByCurrency[cryptocurrency.Ethereum],
				USDAmount:      100.0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Metadata{
				Rate:           tt.fields.Rate,
				Currency:       tt.fields.Currency,
				TotalAmount:    tt.fields.TotalAmount,
				CurrencySymbol: tt.fields.CurrencySymbol,
				USDAmount:      tt.fields.USDAmount,
				RefreshTime:    tt.fields.RefreshTime,
			}

			wantTotalAmount := tt.want.TotalAmount.ToCryptoAmount()
			got := c.WithTotalAmount(tt.args.cartTotalAmount)
			gotTotalAmount := got.TotalAmount.ToCryptoAmount()
			if gotTotalAmount != wantTotalAmount {
				t.Errorf("WithTotalAmount() = %+v, want %+v", gotTotalAmount, wantTotalAmount)
			}

			wantJsonBytes, _ := json.Marshal(tt.want.TotalAmount)
			gotJsonBytes, _ := json.Marshal(got.TotalAmount)
			if string(gotJsonBytes) != string(wantJsonBytes) {
				t.Errorf("WithTotalAmount() = %+v, want %+v", string(gotJsonBytes), string(wantJsonBytes))
			}
		})
	}
}
