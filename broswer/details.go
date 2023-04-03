package broswer

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"

	pngstructure "github.com/dsoprea/go-png-image-structure/v2"
	"github.com/dustin/go-humanize"
	"github.com/sirupsen/logrus"
)

var (
	LoraRegex = regexp.MustCompile(`<lora:([^:]+):([\d\.]+)>`)
)

type LoraInfo struct {
	// Name is the name of lora
	Name string `json:"name,omitempty"`
	// Weight is the weight of lora
	Weight float32 `json:"weight,omitempty"`
}

type SDParams struct {
	// Prompt is the prompt
	Prompt string `json:"prompt,omitempty"`
	// NegativePrompt is the negative prompt
	NegativePrompt string `json:"negative_prompt,omitempty"`
	// Options is the options
	Options   map[string]string `json:"options,omitempty"`
	OptionStr string            `json:"option_str,omitempty"`
	// Template
	Template string `json:"template,omitempty"`

	// Extra info
	Loras []LoraInfo `json:"loras,omitempty"`
}

type Metadata struct {
	// Name is the name of the file
	Name string `json:"name,omitempty"`
	// Size is the size of the file in human readable format
	Size string `json:"size,omitempty"`

	// Width is the width of the image
	Width int `json:"width,omitempty"`
	// Height is the height of the image
	Height int `json:"height,omitempty"`

	// Metadata is the raw metadata
	Metadata map[string]string `json:"metadata,omitempty"`

	// SDParams is the params for SD
	SDParams *SDParams `json:"sd_params,omitempty"`
}

func parseSDParams(md *Metadata, rawMd map[string]string) {
	params := rawMd["parameters"]
	if params == "" {
		return
	}

	sd := &SDParams{
		Options: make(map[string]string),
	}

	// 1. split negtiave prompt
	s := strings.SplitN(params, "Negative prompt:", 2)
	sd.Prompt = strings.TrimSpace(s[0])

	// 2. if we have a negative prompt, split the rest of the params
	if len(s) == 2 {
		// negative is the first line
		remain := strings.TrimSpace(s[1])
		s = strings.SplitN(remain, "\n", 2)
		sd.NegativePrompt = strings.TrimSpace(s[0])

		// 3. Options are the rest of the params are comma separated, each option has a key and a value, which are separated by a colon
		if len(s) == 2 {

			// 4. Template is special
			if strings.Contains(s[1], "Template:") {
				s = strings.SplitN(s[1], "Template:", 2)
				remain = strings.TrimSpace(s[0])
				template := strings.TrimSpace(s[1])
				sd.Template = template
			} else {
				remain = s[1]
			}

			// process options
			sd.OptionStr = strings.TrimSpace(remain)
			s = strings.Split(remain, ",")
			for _, v := range s {
				s2 := strings.SplitN(v, ":", 2)
				if len(s2) == 2 {
					key := strings.TrimSpace(s2[0])
					value := strings.TrimSpace(s2[1])
					sd.Options[key] = value
				}
			}
		}
	}

	// 5. Parsing Loras
	loras := LoraRegex.FindAllStringSubmatch(sd.Prompt, -1)
	for _, loraMatch := range loras {
		lora := LoraInfo{
			Name: loraMatch[1],
		}
		fmt.Sscanf(loraMatch[2], "%f", &lora.Weight)
		sd.Loras = append(sd.Loras, lora)
	}

	md.SDParams = sd
}

func readPngMetadata(img *ImageFile, size int, md *Metadata) error {
	pmp := pngstructure.NewPngMediaParser()
	chunks, err := pmp.Parse(img, size)
	if err != nil {
		return err
	}

	// first parse the tEXt chunks
	rawMetadata := make(map[string]string)
	for _, c := range chunks.(*pngstructure.ChunkSlice).Chunks() {
		if c.Type == "tEXt" {
			comps := bytes.Split(c.Data, []byte{0})
			if len(comps)%2 == 1 {
				logrus.Error("odd number of components")
			}

			for i := 0; i+1 < len(comps); i += 2 {
				rawMetadata[string(comps[i])] = string(comps[i+1])
			}
		}
	}
	md.Metadata = rawMetadata

	// parse the stable diffusion parameters
	parseSDParams(md, rawMetadata)
	delete(rawMetadata, "parameters")

	return nil
}

func (b *Broswer) Metadata(p string) (*Metadata, error) {
	img, err := b.Open(p)
	if err != nil {
		return nil, err
	}
	defer img.Close()

	stat, err := img.Stat()
	if err != nil {
		return nil, err
	}
	fSize := int(stat.Size())

	md := &Metadata{
		Name: img.Name,
		Size: humanize.Bytes(uint64(stat.Size())),
	}

	ext := strings.ToLower(img.Ext())
	switch ext {
	case ".jpg", ".jpeg":
	case ".png":
		if err := readPngMetadata(img, fSize, md); err != nil {
			return nil, err
		}
	default:
		return nil, ErrUnsupportedImageFormat
	}

	return md, nil
}
