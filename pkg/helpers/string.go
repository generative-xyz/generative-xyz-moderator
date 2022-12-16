package helpers

import (
	"crypto/md5"
	"fmt"
)

func GenerateMd5String(input string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(input)))
}