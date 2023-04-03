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

func encodeIntoJpeg(im image.Image, h, w, q int) (*EncodedImage, error) {
	// resize respect to the ratio
	originalWidth := im.Bounds().Dx()
	originalHeight := im.Bounds().Dy()
	if w > originalWidth {
		w = originalWidth
	}
	if h > originalHeight {
		h = originalHeight
	}
	if w == 0 && h != 0 {
		w = originalWidth * h / originalHeight
	} else if h == 0 && w != 0 {
		h = originalHeight * w / originalWidth
	} else if w == 0 && h == 0 {
		w = originalWidth
		h = originalHeight
	}

	resized := transform.Resize(im, w, h, transform.Linear)

	buf := &buffer.Buffer{}
	if err := jpeg.Encode(buf, resized, &jpeg.Options{Quality: q}); err != nil {
		return nil, err
	}

	return &EncodedImage{
		Data:   buf.Bytes(),
		Width:  w,
		Height: h,
	}, nil
}

// encodeFromPng encode a png image to jpg
func encodeFromPng(img *ImageFile, w, h, q int) (*EncodedImage, error) {
	// decode from png
	im, err := png.Decode(img)
	if err != nil {
		return nil, err
	}

	return encodeIntoJpeg(im, w, h, q)
}

// encodeFromJpg encode a jpg image to jpg
func encodeFromJpg(img *ImageFile, w, h, q int) (*EncodedImage, error) {
	// decode from jpg
	im, err := jpeg.Decode(img)
	if err != nil {
		return nil, err
	}

	return encodeIntoJpeg(im, w, h, q)
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
		return encodeFromJpg(img, width, height, b.quality)
	case ".png":
		return encodeFromPng(img, width, height, b.quality)
	default:
		return nil, ErrUnsupportedImageFormat
	}
}
