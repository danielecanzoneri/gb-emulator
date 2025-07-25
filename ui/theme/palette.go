package theme

import "image/color"

type Palette [4]color.Color

var GameBoyPalette = Palette{
	color.RGBA{R: 198, G: 222, B: 140, A: 255},
	color.RGBA{R: 132, G: 165, B: 99, A: 255},
	color.RGBA{R: 57, G: 97, B: 57, A: 255},
	color.RGBA{R: 8, G: 24, B: 16, A: 255},
}
