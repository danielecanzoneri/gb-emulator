package debugger

import (
	"image/color"
	"math"

	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
)

var (
	backgroundColor  = color.NRGBA{R: 24, G: 28, B: 38, A: 255}    // very dark blue-gray
	labelColor       = color.NRGBA{R: 220, G: 220, B: 230, A: 255} // light gray for text
	titleColor       = color.NRGBA{R: 255, G: 214, B: 102, A: 255} // soft pastel yellow
	buttonLabelColor = color.NRGBA{R: 230, G: 230, B: 240, A: 255} // light gray for button text

	buttonColor        = color.NRGBA{R: 104, G: 129, B: 159, A: 255} // lighter muted blue
	buttonHoverColor   = color.NRGBA{R: 134, G: 159, B: 199, A: 255} // even lighter muted blue
	buttonPressedColor = color.NRGBA{R: 164, G: 189, B: 239, A: 255} // lightest muted blue

	mainColor        = color.NRGBA{R: 44, G: 34, B: 54, A: 255} // dark muted purple
	mainHoverColor   = color.NRGBA{R: 64, G: 54, B: 74, A: 255} // lighter muted purple
	mainPressedColor = color.NRGBA{R: 84, G: 74, B: 94, A: 255} // even lighter muted purple

	breakpointColor        = color.NRGBA{R: 120, G: 60, B: 60, A: 255}   // muted dark red
	breakpointHoverColor   = color.NRGBA{R: 160, G: 80, B: 80, A: 255}   // lighter muted red
	breakpointPressedColor = color.NRGBA{R: 200, G: 100, B: 100, A: 255} // even lighter muted red

	currInstrColor        = color.NRGBA{R: 60, G: 120, B: 100, A: 255}  // muted dark teal
	currInstrHoverColor   = color.NRGBA{R: 80, G: 160, B: 130, A: 255}  // lighter muted teal
	currInstrPressedColor = color.NRGBA{R: 100, G: 200, B: 160, A: 255} // even lighter muted teal

	breakpointCurrColor        = color.NRGBA{R: 180, G: 80, B: 180, A: 255} // vibrant violet
	breakpointCurrHoverColor   = color.NRGBA{R: 200, G: 120, B: 200, A: 255}
	breakpointCurrPressedColor = color.NRGBA{R: 220, G: 160, B: 220, A: 255}

	sliderTrackColor = color.NRGBA{R: 255, G: 255, B: 255, A: 32}
)

var (
	buttonImage = &widget.ButtonImage{
		Idle:    image.NewNineSliceColor(buttonColor),
		Hover:   image.NewNineSliceColor(buttonHoverColor),
		Pressed: image.NewNineSliceColor(buttonPressedColor),
	}
	buttonTextColor = &widget.ButtonTextColor{
		Idle: buttonLabelColor,
	}
)

func blendColors(a, b color.Color) color.Color {
	r1, g1, b1, a1 := a.RGBA()
	r2, g2, b2, a2 := b.RGBA()
	return color.NRGBA64{
		R: uint16(math.Sqrt(float64(r1*r1/2 + r2*r2/2))),
		G: uint16(math.Sqrt(float64(g1*g1/2 + g2*g2/2))),
		B: uint16(math.Sqrt(float64(b1*b1/2 + b2*b2/2))),
		A: uint16((a1 + a2) / 2),
	}
}
