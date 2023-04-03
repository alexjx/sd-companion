package broswer

import (
	"bytes"
	"strings"

	pngstructure "github.com/dsoprea/go-png-image-structure/v2"
	"github.com/dustin/go-humanize"
	"github.com/sirupsen/logrus"
)

type Metadata map[string]string

func parseSDParams(md Metadata) Metadata {
	metadata := Metadata{}

	params := md["parameters"]

	// 1. split negtiave prompt
	s := strings.SplitN(params, "Negative prompt:", 2)
	metadata["Prompt"] = strings.TrimSpace(s[0])

	// 2. if we have a negative prompt, split the rest of the params
	if len(s) == 2 {
		// negative is the first line
		remain := strings.TrimSpace(s[1])
		s = strings.SplitN(remain, "\n", 2)
		metadata["Negative prompt"] = strings.TrimSpace(s[0])

		// 3. Options are the rest of the params are comma separated, each option has a key and a value, which are separated by a colon
		if len(s) == 2 {

			// 4. Template is special
			if strings.Contains(s[1], "Template:") {
				s = strings.SplitN(s[1], "Template:", 2)
				remain = strings.TrimSpace(s[0])
				template := strings.TrimSpace(s[1])
				metadata["Template"] = template
			} else {
				remain = s[1]
			}

			// process options
			s = strings.Split(remain, ",")
			for _, v := range s {
				s2 := strings.SplitN(v, ":", 2)
				if len(s2) == 2 {
					key := strings.TrimSpace(s2[0])
					value := strings.TrimSpace(s2[1])
					metadata[key] = value
				}
			}
		}
	}

	// insert none parameters
	for k, v := range md {
		if k != "parameters" {
			metadata[k] = v
		}
	}

	return metadata
}

func readPngMetadata(img *ImageFile) (Metadata, error) {
	pmp := pngstructure.NewPngMediaParser()

	var fSize int
	if fi, err := img.Stat(); err != nil {
		return nil, err
	} else {
		fSize = int(fi.Size())
	}

	chunks, err := pmp.Parse(img, fSize)
	if err != nil {
		return nil, err
	}

	// first parse the tEXt chunks
	md := Metadata{
		"Name": img.Name,
		"Size": humanize.Bytes(uint64(fSize)),
	}
	for _, c := range chunks.(*pngstructure.ChunkSlice).Chunks() {
		if c.Type == "tEXt" {
			comps := bytes.Split(c.Data, []byte{0})
			if len(comps)%2 == 1 {
				logrus.Error("odd number of components")
			}

			for i := 0; i+1 < len(comps); i += 2 {
				md[string(comps[i])] = string(comps[i+1])
			}
		}
	}

	// parse the stable diffusion parameters
	return parseSDParams(md), nil
}

func (b *Broswer) Metadata(p string) (Metadata, error) {
	img, err := b.Open(p)
	if err != nil {
		return nil, err
	}
	defer img.Close()

	md := Metadata{}

	ext := strings.ToLower(img.Ext())
	switch ext {
	case ".jpg", ".jpeg":
		return md, nil
	case ".png":
		return readPngMetadata(img)
	default:
		return nil, ErrUnsupportedImageFormat
	}
}
