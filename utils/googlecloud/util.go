package googlecloud

import (
	"fmt"
	"strings"
	"time"
)

func NormalizeFileName(name string) string {
	now := time.Now().Unix()
	fileName := strings.ToLower(name)
	fileName = strings.ReplaceAll(fileName, " ", "_")
	fileName = strings.TrimSpace(fileName)
	fileName = fmt.Sprintf("%d-%s", now, fileName)
	return fileName
}
