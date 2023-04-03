package broswer

import (
	"bytes"
	"os"
	"testing"

	pngstructure "github.com/dsoprea/go-png-image-structure/v2"
)

func TestPngChunk(t *testing.T) {
	f := "/workspaces/image-browser/tmp2/2023-04-02/00109-2722977797.png"

	img, err := os.Open(f)
	if err != nil {
		t.Fatal(err)
	}
	defer img.Close()

	fStat, err := img.Stat()
	if err != nil {
		t.Fatal(err)
	}

	fSize := int(fStat.Size())

	pmp := pngstructure.NewPngMediaParser()
	chunks, err := pmp.Parse(img, fSize)
	if err != nil {
		t.Fatal(err)
	}

	for _, c := range chunks.(*pngstructure.ChunkSlice).Chunks() {
		if c.Type == "tEXt" {
			comps := bytes.Split(c.Data, []byte{0})
			if len(comps)%2 == 1 {
				t.Fatal("odd number of components")
			}

		}
	}
}

var metadata = map[string]string{
	"parameters": `(photograph:1.3), (photorealistic:1.5), 8K, HDR, (masterpiece:1.3), absurdres, professional lighting, intricate detailed, DSLR, pov, bokeh, 1girl, solo, extremely attractive young girl, (beautiful face:1.2), (slim body:1.2), (slim thighs:1.2), big eyes, narrow waist, medium breasts, smooth skin, small areolae, detailed breasts, collarbone, pale skin, looking at viewer, beautiful nails, nose blush, beautiful nipples, blue eyes, long hair, floating hair, blue hair, (thigh gap:1.2), lipgloss, virgin killer sweater, front open, stockings, jewelry, hair ornament, beach, sunny, outdoors, (ulzzang-6500-v1.1:0.3), <lora:fashionGirl_v52:0.2>, fashi-g, <lora:realisticVaginasGodPussy_godpussy2Innie:0.2> <lora:breastinclassBetter_v14:0.2> <lora:cuteGirlMix4_v10:0.3>
Negative prompt: (logo:2), (watermark:2), (worst quality:2), (normal quality:2), (lowres:2), (low quality:2), (ugly:2), (painting:2), (sketch:2), (illustration:2), obesity, bad anatomy, (bad-hands-5:1.2), (bad-image-v2-39000), extra ear, (extra arm), (extra hand), extra limb, disfigured, extra breasts, extra navel, (skindentation:2), (mole:2), (skin spots:2), (acnes:2), (skin blemishes:2), age spot, extra fingers, amputated finger, amputated arm, amputated leg, (bad_prompt_version2:1.2), thick thighs, (black spot:2),
Steps: 50, Sampler: DPM++ 2M Karras, CFG scale: 8, Seed: 692541119, Size: 512x768, Model hash: fc2511737a, Model: chilloutmix_NiPrunedFp32Fix, Denoising strength: 0.3, Hires upscale: 2, Hires steps: 30, Hires upscaler: 4x-UltraSharp
Template: # quality
(photograph:1.3), (photorealistic:1.5), 8K, HDR, (masterpiece:1.3), absurdres, professional lighting, intricate detailed, DSLR, pov, bokeh,

# main content
1girl, solo, extremely attractive young girl, (beautiful face:1.2), (slim body:1.2), (slim thighs:1.2), big eyes, narrow waist,  medium breasts, smooth skin, small areolae,
detailed breasts, collarbone, pale skin, looking at viewer, beautiful nails, nose blush, beautiful nipples,

# iris
__devilkkw_mod/body-1/eyes_iris_colors__,

# hair
long hair, floating hair, __devilkkw_mod/body-1/hair_color__,

# describe body
(thigh gap:1.2), lipgloss,

# cloth (if there is any)
virgin killer sweater, front open, stockings,

# decorations
jewelry, hair ornament,

# environment
{snow|beach}, sunny, outdoors,

# lora, don't enable them all
(ulzzang-6500-v1.1:{0.3|0.4|0.5}),

<lora:fashionGirl_v52:0.2>, fashi-g,

# <lora:inniesBetterVulva_v11:0.2>
<lora:realisticVaginasGodPussy_godpussy2Innie:0.2>
<lora:breastinclassBetter_v14:0.2>

<lora:cuteGirlMix4_v10:{0.2|0.3|0.4}>

{0-1$$<lora:koreanDollLikeness_v15:{0.1|0.2|0.3|0.4|0.5}>
|<lora:japaneseDollLikeness_v10:{0.1|0.2|0.3|0.4|0.5}>
|<lora:LORAChineseDoll_chinesedolllikeness1:{0.1|0.2|0.3}>
|<lora:mikuya_v15:{0.1|0.2|0.3}>
|<lora:irene_V70:{0.2|0.3|0.4}>
|<lora:liyuuLora_liyuuV1:{0.1|0.2|0.3}>
|<lora:cnGirlYcy_v10:{0.2|0.3}>
}


Negative Template: (logo:2), (watermark:2), (worst quality:2), (normal quality:2), (lowres:2), (low quality:2), (ugly:2), (painting:2), (sketch:2), (illustration:2), obesity, bad anatomy,
(bad-hands-5:1.2), (bad-image-v2-39000), extra ear, (extra arm), (extra hand), extra limb, disfigured, extra breasts, extra navel, (skindentation:2), (mole:2), (skin spots:2), (acnes:2),
(skin blemishes:2), age spot, extra fingers, amputated finger, amputated arm, amputated leg, (bad_prompt_version2:1.2), thick thighs, (black spot:2),
`,
}

func TestMetadata(t *testing.T) {
	parsed := parseSDParams(metadata)

	if parsed["Steps"] != "50" {
		t.Fatal("steps not parsed")
	}

}
