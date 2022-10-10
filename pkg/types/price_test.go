package types

import (
	"reflect"
	"strings"
	"testing"
)

func TestNewPrice(t *testing.T) {
	type args struct {
		value float64
		opts  []Option
	}
	tests := []struct {
		name         string
		args         args
		want         *Price
		wantStrValue string
	}{
		{
			name: "12.00",
			args: args{
				value: 12.00,
				opts:  []Option{WithPrecision(2), WithNeedBeautyZeroPrecision(true)},
			},
			want: NewPrice(12.00, WithPrecision(0), WithNeedBeautyZeroPrecision(true)),
		},
		{
			name: "12.00",
			args: args{
				value: 12.00,
				opts:  []Option{WithNeedBeautyZeroPrecision(false)},
			},
			want: NewPrice(12.00, WithNeedBeautyZeroPrecision(false)),
		},
		{
			name: "12.12",
			args: args{
				value: 12.12,
				opts:  []Option{WithNeedBeautyZeroPrecision(false)},
			},
			want: NewPrice(12.12, WithNeedBeautyZeroPrecision(false)),
		},
		{
			name: "121234.12",
			args: args{
				value: 121234.12,
				opts:  []Option{},
			},
			want: NewPrice(121234.12, WithNeedBeautyZeroPrecision(false)),
		},
		{
			name: "121234.123",
			args: args{
				value: 121234.123,
				opts:  []Option{WithPrecision(3), WithNeedBeautyZeroPrecision(false)},
			},
			want: NewPrice(121234.123, WithPrecision(3), WithNeedBeautyZeroPrecision(false)),
		},
		{
			name: "Big number",
			args: args{
				value: 14.607726642223438,
				opts:  []Option{WithPrecision(6), WithCryptoPrecision(18)},
			},
			want: NewPrice(14.607726642223438, WithCryptoPrecision(18), WithPrecision(6)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewPrice(tt.args.value, tt.args.opts...)
			if !reflect.DeepEqual(got.float64Value, tt.want.float64Value) {
				t.Errorf("NewPrice() = %v, want %v", got, tt.want)
			}

			cryptoAmount := got.ToCryptoAmount()
			if strings.Contains(cryptoAmount, "-") {
				t.Errorf("Convert to cryptoAmount failed: %s", cryptoAmount)
			}

			gotJsonBytes, err := got.MarshalJSON()
			if err != nil {
				t.Error(err)
			}

			wantJsonBytes, err := tt.want.MarshalJSON()
			if err != nil {
				t.Error(err)
			}

			if string(wantJsonBytes) != string(gotJsonBytes) {
				t.Errorf("NewPrice() = %v, want %v", string(gotJsonBytes), string(wantJsonBytes))
			}

			gotPriceAfterUnmarshal := ZeroPriceP
			err = gotPriceAfterUnmarshal.UnmarshalJSON(gotJsonBytes)
			if err != nil {
				t.Error(err)
			}
			if !reflect.DeepEqual(gotPriceAfterUnmarshal, tt.want) {
				t.Errorf("NewPrice() = %v, want %v", gotPriceAfterUnmarshal, tt.want)
			}
		})
	}
}
