package utils_test

import (
	"fmt"
	"net/http"
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

func TestConvertIpfsToHttp(t *testing.T) {
	fmt.Println(utils.ConvertIpfsToHttp("ipfs://QmaLTjaNb4bZXUqWxGcqxSVCmV4BecxLCDp3FzmKZ4fo6C"))
	uploadImage := func(urlStr string) string {
		urlStr = utils.ConvertIpfsToHttp(urlStr)
		client := http.Client{}
		r, err := client.Get(urlStr)
		if err != nil {
			return urlStr
		}
		defer r.Body.Close()
		fmt.Println(r.Header.Get("content-type"))
		// buf, err := io.ReadAll(r.Body)
		// if err != nil {
		// 	return urlStr
		// }
		// ext, err := utils.GetFileExtensionFromUrl(urlStr)
		// if err != nil {
		// 	contentType := r.Header.Get("content-type")
		// 	arr := strings.Split(contentType, "/")
		// 	if len(arr) > 0 {
		// 		ext = arr[1]
		// 	}
		// }
		// name := fmt.Sprintf("%v.%s", uuid.New().String(), ext)
		// _, err = u.GCS.UploadBaseToBucket(helpers.Base64Encode(buf), name)
		// if err != nil {
		// 	return urlStr
		// }
		return urlStr
	}
	uploadImage("ipfs://QmaLTjaNb4bZXUqWxGcqxSVCmV4BecxLCDp3FzmKZ4fo6C")
}
