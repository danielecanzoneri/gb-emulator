package debugger

import (
	"image/color"

	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
)

var (
	// General Background and Text
	backgroundColor = color.NRGBA{R: 20, G: 22, B: 30, A: 255}    // nearly black-blue
	labelColor      = color.NRGBA{R: 220, G: 220, B: 230, A: 255} // light gray for readability
	titleColor      = color.NRGBA{R: 255, G: 200, B: 80, A: 255}  // soft warm yellow

	// Buttons
	buttonColor        = color.NRGBA{R: 60, G: 90, B: 130, A: 255}   // desaturated blue
	buttonHoverColor   = color.NRGBA{R: 80, G: 110, B: 160, A: 255}  // brighter on hover
	buttonPressedColor = color.NRGBA{R: 100, G: 130, B: 190, A: 255} // even brighter
	buttonLabelColor   = color.NRGBA{R: 235, G: 235, B: 245, A: 255} // high contrast

	// Main UI Panels / Containers
	mainColor        = color.NRGBA{R: 35, G: 30, B: 45, A: 255} // deep purple-gray
	mainHoverColor   = color.NRGBA{R: 50, G: 45, B: 65, A: 255} // slightly lighter
	mainPressedColor = color.NRGBA{R: 65, G: 60, B: 85, A: 255} // more lifted

	// Breakpoints
	breakpointColor        = color.NRGBA{R: 140, G: 40, B: 40, A: 255} // dark red
	breakpointHoverColor   = color.NRGBA{R: 170, G: 60, B: 60, A: 255} // hover pop
	breakpointPressedColor = color.NRGBA{R: 200, G: 80, B: 80, A: 255} // more vivid

	// Current Instruction
	currInstrColor        = color.NRGBA{R: 40, G: 120, B: 110, A: 255} // dark teal
	currInstrHoverColor   = color.NRGBA{R: 60, G: 150, B: 140, A: 255} // hover highlight
	currInstrPressedColor = color.NRGBA{R: 80, G: 180, B: 170, A: 255} // active highlight

	// Breakpoint + Current Instruction
	breakpointCurrColor        = color.NRGBA{R: 160, G: 60, B: 160, A: 255} // magenta
	breakpointCurrHoverColor   = color.NRGBA{R: 180, G: 100, B: 180, A: 255}
	breakpointCurrPressedColor = color.NRGBA{R: 200, G: 140, B: 200, A: 255}

	// Toolbar
	toolbarBackgroundColor   = color.NRGBA{R: 18, G: 20, B: 30, A: 255} // matches general bg
	toolbarTextColor         = color.NRGBA{R: 230, G: 230, B: 240, A: 255}
	toolbarMenuColor         = color.NRGBA{R: 30, G: 30, B: 40, A: 255}
	toolbarMenuHoverColor    = color.NRGBA{R: 40, G: 40, B: 50, A: 255}
	toolbarMenuPressedColor  = color.NRGBA{R: 50, G: 50, B: 60, A: 255}
	toolbarEntryColor        = toolbarMenuColor
	toolbarEntryHoverColor   = toolbarMenuHoverColor
	toolbarEntryPressedColor = toolbarMenuPressedColor

	// Slider
	sliderTrackColor = color.NRGBA{R: 255, G: 255, B: 255, A: 48} // semi-transparent white

	// Floating window
	windowColor   = color.NRGBA{R: 28, G: 30, B: 42, A: 255}
	titleBarColor = color.NRGBA{R: 22, G: 24, B: 34, A: 255}
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
	padding = 5
	insets  = widget.NewInsetsSimple(padding)
)
