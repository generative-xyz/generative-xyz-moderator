package utils_test

import (
	"testing"

	"github.com/test-go/testify/assert"
	"rederinghub.io/utils"
)

func TestGetFileExtensionFromUrl(t *testing.T) {
	tests := []struct {
		fullUrl, want string
		wantErr       bool
	}{
		{
			fullUrl: "https://assets.somedomain.com/assets/files/mypicture.jpg?width=1000&height=600",
			want:    "jpg",
			wantErr: false,
		},
		{
			fullUrl: "https://cryptopunks.app/cryptopunks/cryptopunk000.png",
			want:    "png",
			wantErr: false,
		}, {
			fullUrl: "https://google.com",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.fullUrl, func(t *testing.T) {
			ext, err := utils.GetFileExtensionFromUrl(tt.fullUrl)
			if tt.wantErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, ext, tt.want)
			}
		})
	}
}
