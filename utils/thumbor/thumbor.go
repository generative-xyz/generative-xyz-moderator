package thumbor

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
)

type CropValue struct {
	Top    int
	Left   int
	Right  int
	Bottom int
}

type Thumbor struct {
	secretKey string
	serverUrl string

	imagePath            string
	width                int
	height               int
	smart                bool
	fitInFlag            bool
	withFlipHorizontally bool
	withFlipVertically   bool
	halignValue          string
	valignValue          string
	cropValues           *CropValue
	meta                 bool
	filtersCalls         []string
}

func ProvideThumbor(cnf Config) *Thumbor {
	return &Thumbor{
		secretKey: cnf.SecretKey,
		serverUrl: cnf.ServerUrl,

		//
		imagePath:            "",
		width:                -1,
		height:               -1,
		smart:                false,
		fitInFlag:            false,
		withFlipHorizontally: false,
		withFlipVertically:   false,
		halignValue:          "",
		valignValue:          "",
		cropValues:           nil,
		meta:                 false,
		filtersCalls:         []string{},
	}
}

func (t *Thumbor) Reset() {
	t.imagePath = ""
	t.width = -1
	t.height = -1
	t.smart = false
	t.fitInFlag = false
	t.withFlipHorizontally = false
	t.withFlipVertically = false
	t.halignValue = ""
	t.valignValue = ""
	t.cropValues = nil
	t.meta = false
	t.filtersCalls = []string{}
}

// getOperationPath Converts operation array to string
func (t *Thumbor) getOperationPath() (string, error) {
	parts, err := t.urlParts()
	if err != nil {
		return "", err
	}

	if len(parts) == 0 {
		return "", nil
	}

	return fmt.Sprintf("%s/", strings.Join(parts, "/")), nil
}

// urlParts Build operation array
func (t *Thumbor) urlParts() ([]string, error) {
	if t.imagePath == "" {
		return nil, errors.New("the image url can't be null or empty")
	}

	parts := []string{}

	if t.meta {
		parts = append(parts, "meta")
	}

	if t.cropValues != nil {
		parts = append(parts, fmt.Sprintf("%dx%d:%dx%d", t.cropValues.Left, t.cropValues.Top, t.cropValues.Right, t.cropValues.Bottom))
	}

	if t.fitInFlag {
		parts = append(parts, "fit-in")
	}

	if t.width != -1 || t.height != -1 || t.withFlipHorizontally || t.withFlipVertically {
		sizeString := ""
		if t.withFlipHorizontally {
			sizeString = fmt.Sprintf("%s-", sizeString)
		}
		sizeString = fmt.Sprintf("%s%dx", sizeString, t.width)
		if t.withFlipVertically {
			sizeString = fmt.Sprintf("%s-", sizeString)
		}
		sizeString = fmt.Sprintf("%s%d", sizeString, t.height)
		parts = append(parts, sizeString)
	}

	if t.halignValue != "" {
		parts = append(parts, t.halignValue)
	}

	if t.valignValue != "" {
		parts = append(parts, t.valignValue)
	}

	if t.smart {
		parts = append(parts, "smart")
	}

	if len(t.filtersCalls) > 0 {
		parts = append(parts, fmt.Sprintf("filters:%s", strings.Join(t.filtersCalls, ":")))
	}

	return parts, nil
}

func (t *Thumbor) SetImagePath(path string) *Thumbor {
	t.imagePath = path
	if path[0:1] == "/" {
		pathLen := len(path)
		t.imagePath = path[1:pathLen]
	}
	return t
}

func (t *Thumbor) Resize(width int, height int) *Thumbor {
	t.width = width
	t.height = height
	t.fitInFlag = false
	return t
}

func (t *Thumbor) Format(format string) *Thumbor {
	t.Filter(fmt.Sprintf("format(%s)", format))
	return t
}

func (t *Thumbor) Compress(quality int) *Thumbor {
	t.Filter(fmt.Sprintf("quality(%d)", quality))
	return t
}

func (t *Thumbor) SmartCrop(smart bool) *Thumbor {
	t.smart = smart
	return t
}

func (t *Thumbor) FitIn(width int, height int) *Thumbor {
	t.width = width
	t.height = height
	t.fitInFlag = true
	return t
}

func (t *Thumbor) FlipHorizontally() *Thumbor {
	t.withFlipHorizontally = true
	return t
}

func (t *Thumbor) FlipVertically() *Thumbor {
	t.withFlipVertically = true
	return t
}

func (t *Thumbor) HAlign(value string) *Thumbor {
	if value == "left" || value == "right" || value == "center" {
		t.halignValue = value
	}
	return t
}

func (t *Thumbor) VAlign(value string) *Thumbor {
	if value == "top" || value == "middle" || value == "bottom" {
		t.valignValue = value
	}
	return t
}

func (t *Thumbor) MetadataOnly(value bool) *Thumbor {
	t.meta = value
	return t
}

func (t *Thumbor) Filter(value string) *Thumbor {
	t.filtersCalls = append(t.filtersCalls, value)
	return t
}

func (t *Thumbor) Crop(top int, left int, right int, bottom int) *Thumbor {
	t.cropValues = &CropValue{
		Top:    top,
		Left:   left,
		Right:  right,
		Bottom: bottom,
	}
	return t
}

func (t *Thumbor) BuildUrl() (string, error) {
	operation, err := t.getOperationPath()
	if err != nil {
		return "", err
	}
	if t.secretKey != "" {
		key := hmac.New(sha1.New, []byte(t.secretKey))
		key.Write([]byte(fmt.Sprintf("%s%s", operation, t.imagePath)))
		keyB64 := base64.StdEncoding.EncodeToString(key.Sum(nil))
		keyB64 = strings.Replace(keyB64, "+", "-", -1)
		keyB64 = strings.Replace(keyB64, "/", "_", -1)

		return fmt.Sprintf("%s/%s/%s%s", t.serverUrl, keyB64, operation, t.imagePath), nil
	} else {
		return fmt.Sprintf("%s/unsafe/%s%s", t.serverUrl, operation, t.imagePath), nil
	}
}
