package fileutil

import (
	"bytes"
	"errors"
	"image"
	"strings"

	"github.com/disintegration/imaging"
	"rederinghub.io/utils/thumbor"
)

const (
	MaxImageByteSize = 400 * 1024 // 400KB
	MaxImagePx       = 1024       // px
	JpegQuality      = 90         // %
)

func ResizeImage(imgSrc []byte, imageType string, maxImageByteSize int) ([]byte, error) {
	byteSize := len(imgSrc)
	if byteSize <= maxImageByteSize {
		return imgSrc, nil
	}
	img, err := imaging.Decode(bytes.NewReader(imgSrc))
	if err != nil {
		return nil, err
	}
	var imgByte []byte
	switch strings.ToLower(imageType) {
	case "png":
		imgByte, err = resize(img, imaging.PNG)
	case "jpeg", "jpg":
		imgByte, err = resize(img, imaging.JPEG, imaging.JPEGQuality(JpegQuality))
	case "gif":
		imgByte, err = resize(img, imaging.GIF)
	default:
		return nil, errors.New("image not support")
	}
	if err != nil {
		return nil, err
	}
	return imgByte, nil
}
func resize(img image.Image, format imaging.Format, opts ...imaging.EncodeOption) ([]byte, error) {
	resized := scaleDown(img)
	var buf bytes.Buffer
	if err := imaging.Encode(&buf, resized, format, opts...); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
func scaleDown(src image.Image) image.Image {
	rect := src.Bounds()

	x := rect.Dx()
	y := rect.Dy()

	switch {
	case MaxImagePx > x && MaxImagePx > y:
		return src
	case x > y:
		return imaging.Resize(src, MaxImagePx, 0, imaging.MitchellNetravali)
	default:
		return imaging.Resize(src, 0, MaxImagePx, imaging.MitchellNetravali)
	}
}

func ImageCompress(imageUrl string, quality int, key string) (string, error) {

	tb := thumbor.ProvideThumbor(thumbor.Config{
		ServerUrl: "https://thumbor.generative.xyz",
		SecretKey: key,
	})

	// imageUrl := "https://soulgenesis.art/api/images/1104/11cTdjUhh4h477-stage1.jpg"
	return tb.SetImagePath(imageUrl).Compress(10).BuildUrl()

}
