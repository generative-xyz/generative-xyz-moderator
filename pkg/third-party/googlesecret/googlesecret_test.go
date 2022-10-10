package googlesecret

import (
	"context"
	"google.golang.org/api/option"
	"testing"
)

func TestGoogleSecret_GetSecret(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "test1",
			args: args{
				name: "projects/342393054229/secrets/autonomous-crypto-payment-public-key-dev/versions/latest",
			},
			want:    "",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g, err := NewClient(
				context.Background(),
				option.WithCredentialsFile("/Volumes/data/autonomous_project/auto-backend-microservices/bff-service/setting/google-secret-key.json"),
			)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.GetSecret() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			got, err := g.GetSecret(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetSecret() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetSecret() got = %v, want %v", got, tt.want)
			}
		})
	}
}
