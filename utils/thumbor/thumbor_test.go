package thumbor

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestThumbor_BuildUrl(t *testing.T) {
	tb := ProvideThumbor(Config{
		ServerUrl: "https://thumbor.generative.xyz",
		SecretKey: "",
	})

	imageUrl := "https://soulgenesis.art/api/images/1104/11cTdjUhh4h477-stage1.jpg"
	thumborUrl, err := tb.SetImagePath(imageUrl).SmartCrop(true).Format("jpg").Compress(90).BuildUrl()

	log.Println("thumborUrl", thumborUrl)

	result := "https://thumbor.generative.xyz/eP9spWNpQM8IYCVLFnwFV2KIWwo=/1000x0/smart/https://soulgenesis.art/api/images/1104/11cTdjUhh4h477-stage1.jpg"
	assert.NoError(t, err, "no error")
	assert.Equal(t, thumborUrl, result, "equal result")
}
