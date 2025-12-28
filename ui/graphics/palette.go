package theme

import "image/color"

type Palette interface {
	Get(uint16) color.Color
}

var dmgPalette = [4]color.Color{
	color.RGBA{R: 198, G: 222, B: 140, A: 255},
	color.RGBA{R: 132, G: 165, B: 99, A: 255},
	color.RGBA{R: 57, G: 97, B: 57, A: 255},
	color.RGBA{R: 8, G: 24, B: 16, A: 255},
}

type DMGPalette struct{}

func (p DMGPalette) Get(c uint16) color.Color {
	return dmgPalette[c]
}

type CGBColor struct {
	// 5 bit
	r, g, b uint8
}

func (c CGBColor) RGBA() (r, g, b, a uint32) {
	r = uint32(float32(c.r) / 0x1f * 0xffff)
	g = uint32(float32(c.g) / 0x1f * 0xffff)
	b = uint32(float32(c.b) / 0x1f * 0xffff)
	a = 0xffff
	return
}

type CGBPalette struct{}

func (p CGBPalette) Get(c uint16) color.Color {
	return CGBColor{
		r: uint8(c & 0x1F),
		g: uint8((c >> 5) & 0x1F),
		b: uint8((c >> 10) & 0x1F),
	}
}
