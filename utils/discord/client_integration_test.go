//go:build integration

package discordclient

import (
	"context"
	"testing"
	"time"
)

func TestClient_SendMessage(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancelFn()
	type fields struct {
		WebhookURL string
	}
	type args struct {
		message Message
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "test send simple discord message",
			fields: fields{},
			args: args{
				message: Message{
					Username: "Test Bot",
					Content:  "Test data",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{}
			if err := c.SendMessage(ctx, "https://discord.com/api/webhooks/1075257578910666784/9RCSGaDGeAgLn59wWk6L4SB-P6A9wUxWy0qO_uN-EYJb3BCnn25kp8ID42ixYCkllyMf", tt.args.message); (err != nil) != tt.wantErr {
				t.Errorf("SendMessage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
