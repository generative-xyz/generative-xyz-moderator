package helpers

import (
	"crypto/md5"
	b64 "encoding/base64"
	"fmt"
)

func GenerateMd5String(input string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(input)))
}

func Base64Decode(base64Str string) ([]byte, error) {
	sDec, err := b64.StdEncoding.DecodeString(base64Str)
	if err != nil {
		return nil, err
	}
    return sDec, nil
}

func Base64Eecode(data []byte) string {
	sDec := b64.StdEncoding.EncodeToString(data)
    return sDec
}