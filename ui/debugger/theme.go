package debugger

import (
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"golang.org/x/image/colornames"
	"image/color"
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
