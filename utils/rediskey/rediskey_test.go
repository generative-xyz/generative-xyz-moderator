package rediskey_test

import (
	"reflect"
	"testing"

	"rederinghub.io/utils/rediskey"
)

func TestBeauty(t *testing.T) {
	type args struct {
		prefix     string
		uuid       int
		params     []interface{}
		paramsHash []interface{}
		structHash struct {
			Test string
		}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "No uuid",
			args: args{
				prefix: "ahihi",
			},
			want: "ahihi",
		},
		{
			name: "With uuid",
			args: args{
				prefix: "ahihi",
				uuid:   12,
			},
			want: "ahihi_12",
		},
		{
			name: "With struct hash",
			args: args{
				prefix: "ahihi_struct",
				structHash: struct {
					Test string
				}{
					Test: "123",
				},
			},
			want: "ahihi_struct_3719508699856307341",
		},
		{
			name: "With params",
			args: args{
				prefix: "ahihi",
				params: []interface{}{12},
			},
			want: "ahihi_12",
		},
		{
			name: "With params hash",
			args: args{
				prefix:     "ahihi",
				paramsHash: []interface{}{12},
			},
			want: "ahihi_c20ad4d76fe97759aa27a0c99bff6710",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := rediskey.Beauty(tt.args.prefix)
			if tt.args.uuid > 0 {
				got = got.WithUUID(tt.args.uuid)
			}

			if len(tt.args.params) > 0 {
				got = got.WithParams(tt.args.params...)
			}

			if tt.args.structHash.Test != "" {
				got = got.WithStructHash(tt.args.structHash, nil)
			}
			if len(tt.args.paramsHash) > 0 {
				got = got.WithParamsHash(tt.args.paramsHash...)
			}

			if !reflect.DeepEqual(got.String(), tt.want) {
				t.Errorf("Beauty() = %v, want %v", got, tt.want)
			}
		})
	}
}
