package fileutil

import (
	"bytes"
	"image"

	"github.com/disintegration/imaging"
)

const (
	MaxImageByteSize = 1024 * 1024 // 1MB
	MaxImagePx       = 1024        // px
	JpegQuality      = 90          // %
)

func ResizeImage(img image.Image, format imaging.Format, opts ...imaging.EncodeOption) ([]byte, error) {
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
