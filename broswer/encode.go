package broswer

import (
	"image"
	"image/jpeg"
	"image/png"

	"github.com/anthonynsimon/bild/transform"
	"github.com/sirupsen/logrus"
	"go.uber.org/zap/buffer"
)

// EncodedImage repsent a jpg image, optizmied for fast transfer
type EncodedImage struct {
	// Data is the encoded image data
	Data []byte
	// Width is the width of the image
	Width int
	// Height is the height of the image
	Height int
}

func encodeIntoJpeg(im image.Image, width int, height int) (*EncodedImage, error) {
	// resize respect to the ratio
	originalWidth := im.Bounds().Dx()
	originalHeight := im.Bounds().Dy()

	if width == 0 && height != 0 {
		width = originalWidth * height / originalHeight
	} else if height == 0 && width != 0 {
		height = originalHeight * width / originalWidth
	} else if width == 0 && height == 0 {
		width = originalWidth
		height = originalHeight
	}

	resized := transform.Resize(im, width, height, transform.NearestNeighbor)

	buf := &buffer.Buffer{}
	if err := jpeg.Encode(buf, resized, &jpeg.Options{Quality: 80}); err != nil {
		return nil, err
	}

	return &EncodedImage{
		Data:   buf.Bytes(),
		Width:  width,
		Height: height,
	}, nil
}

// encodeFromPng encode a png image to jpg
func encodeFromPng(img *ImageFile, width, height int) (*EncodedImage, error) {
	// decode from png
	im, err := png.Decode(img)
	if err != nil {
		return nil, err
	}

	return encodeIntoJpeg(im, width, height)
}

// encodeFromJpg encode a jpg image to jpg
func encodeFromJpg(img *ImageFile, width, height int) (*EncodedImage, error) {
	// decode from jpg
	im, err := jpeg.Decode(img)
	if err != nil {
		return nil, err
	}

	return encodeIntoJpeg(im, width, height)
}

func (b *Broswer) Encoded(p string, width, height int) (*EncodedImage, error) {
	logrus.Debugf("encode image: %s into %d x %d", p, width, height)
	// load image
	img, err := b.Open(p)
	if err != nil {
		return nil, err
	}
	defer img.Close()

	ext := img.Ext()
	switch ext {
	case ".jpg", ".jpeg":
		return encodeFromJpg(img, width, height)
	case ".png":
		return encodeFromPng(img, width, height)
	default:
		return nil, ErrUnsupportedImageFormat
	}
}
