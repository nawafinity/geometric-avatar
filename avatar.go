package geometric_avatar

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/fogleman/gg"
	"image/color"
	"math/big"
)

type Avatar struct {
	hash     string
	options  map[string]interface{}
	defaults map[string]interface{}
}

func GenerateAvatar(email string, options map[string]interface{}) *Avatar {
	hash := CreateHashFromString(email)

	defaults := map[string]interface{}{
		"background": color.RGBA{238, 238, 238, 255},
		"hash":       CreateHashFromString(""),
		"margin":     0.08,
		"size":       128.0,
	}

	for key, value := range options {
		defaults[key] = value
	}

	_, ok := defaults["size"].(float64)
	if !ok {
		_ = 128.0
	}

	_, ok = defaults["margin"].(float64)
	if !ok {
		_ = 0.08
	}

	imageGen := &Avatar{
		hash:     hash,
		options:  defaults,
		defaults: defaults,
	}

	return imageGen
}
func (ig *Avatar) rectangle(dc *gg.Context, x, y, w, h float64, color color.Color) {
	dc.SetColor(color)
	dc.DrawRectangle(x, y, w, h)
	dc.Fill()
}

func (ig *Avatar) hslToRGB(h, s, l float64) color.RGBA {
	var r, g, b float64

	if s == 0 {
		r = l
		g = l
		b = l
	} else {
		var hue2rgb = func(p, q, t float64) float64 {
			if t < 0 {
				t += 1
			}
			if t > 1 {
				t -= 1
			}
			if t < 1/6.0 {
				return p + (q-p)*6*t
			}
			if t < 1/2.0 {
				return q
			}
			if t < 2/3.0 {
				return p + (q-p)*(2/3.0-t)*6
			}
			return p
		}

		q := l + s - l*s
		p := 2*l - q

		r = hue2rgb(p, q, h+1/3.0)
		g = hue2rgb(p, q, h)
		b = hue2rgb(p, q, h-1/3.0)
	}

	return color.RGBA{uint8(r * 255), uint8(g * 255), uint8(b * 255), 255}
}

func (ig *Avatar) Render() *gg.Context {
	size := ig.options["size"].(float64)
	margin := size * ig.options["margin"].(float64)
	cell := (size - (margin * 2)) / 5

	dc := gg.NewContext(int(size), int(size))
	dc.SetColor(ig.options["background"].(color.Color))
	dc.Clear()

	hashInt, _ := new(big.Int).SetString(ig.hash[len(ig.hash)-7:], 16)
	hue := float64(hashInt.Uint64()) / 0xfffffff
	fgColor := ig.hslToRGB(hue, 0.5, 0.7)
	bgColor := ig.options["background"].(color.Color)

	for i := 0; i < 15; i++ {
		var color color.Color
		if int(ig.hash[i]-'0')%2 == 0 {
			color = bgColor
		} else {
			color = fgColor
		}

		if i < 5 {
			ig.rectangle(dc, 2*cell+margin, float64(i)*cell+margin, cell, cell, color)
		} else if i < 10 {
			ig.rectangle(dc, 1*cell+margin, float64(i-5)*cell+margin, cell, cell, color)
			ig.rectangle(dc, 3*cell+margin, float64(i-5)*cell+margin, cell, cell, color)
		} else {
			ig.rectangle(dc, 0*cell+margin, float64(i-10)*cell+margin, cell, cell, color)
			ig.rectangle(dc, 4*cell+margin, float64(i-10)*cell+margin, cell, cell, color)
		}
	}

	return dc
}

func CreateHashFromString(str string) string {
	hasher := md5.New()
	hasher.Write([]byte(str))
	return hex.EncodeToString(hasher.Sum(nil))
}
