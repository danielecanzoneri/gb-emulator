package debugger

import (
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"golang.org/x/image/colornames"
	"image/color"
	"math"
)

var (
	mainColor        = color.NRGBA{R: 0x13, G: 0x1A, B: 0x22, A: 0xFF}
	labelColor       = colornames.White
	titleColor       = colornames.Yellow
	buttonLabelColor = colornames.Black

	buttonColor        = color.NRGBA{R: 170, G: 170, B: 180, A: 255}
	buttonHoverColor   = color.NRGBA{R: 130, G: 130, B: 150, A: 255}
	buttonPressedColor = color.NRGBA{R: 90, G: 90, B: 120, A: 255}

	breakpointColor        = color.NRGBA{R: 255, G: 170, B: 180, A: 255}
	breakpointHoverColor   = color.NRGBA{R: 255, G: 130, B: 150, A: 255}
	breakpointPressedColor = color.NRGBA{R: 255, G: 90, B: 120, A: 255}

	currInstrColor        = color.NRGBA{R: 170, G: 170, B: 255, A: 255}
	currInstrHoverColor   = color.NRGBA{R: 130, G: 130, B: 255, A: 255}
	currInstrPressedColor = color.NRGBA{R: 90, G: 90, B: 255, A: 255}

	breakpointCurrColor        = blendColors(breakpointColor, currInstrColor)
	breakpointCurrHoverColor   = blendColors(breakpointHoverColor, currInstrHoverColor)
	breakpointCurrPressedColor = blendColors(breakpointPressedColor, currInstrPressedColor)
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
