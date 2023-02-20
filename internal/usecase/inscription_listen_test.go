package usecase

import (
	"testing"

	"github.com/opentracing/opentracing-go"
)

func Test_getLatestOrdServiceBlockCount(t *testing.T) {
	type args struct {
		rootSpan  opentracing.Span
		ordServer string
	}
	tests := []struct {
		name    string
		args    args
		want    int64
		wantErr bool
	}{
		{
			name: "test",
			args: args{
				rootSpan:  nil,
				ordServer: "https://ordinals-explorer.generative.xyz",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getLatestOrdServiceBlockCount(tt.args.rootSpan, tt.args.ordServer)
			if (err != nil) != tt.wantErr {
				t.Errorf("getLatestOrdServiceBlockCount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getLatestOrdServiceBlockCount() = %v, want %v", got, tt.want)
			}
		})
	}
}
