package theme

import (
	"github.com/ebitenui/ebitenui/image"
	"image/color"

	"github.com/ebitenui/ebitenui/widget"
)

type DebuggerTheme struct {
	BackgroundColor color.NRGBA
	LabelColor      color.NRGBA
	TitleColor      color.NRGBA

	Button struct {
		Image     *widget.ButtonImage
		TextColor *widget.ButtonTextColor
	}

	Main struct {
		Color        color.NRGBA
		HoverColor   color.NRGBA
		PressedColor color.NRGBA
	}

	Disassembler struct {
		BreakpointColor            color.NRGBA
		BreakpointHoverColor       color.NRGBA
		BreakpointPressedColor     color.NRGBA
		CurrInstrColor             color.NRGBA
		CurrInstrHoverColor        color.NRGBA
		CurrInstrPressedColor      color.NRGBA
		BreakpointCurrColor        color.NRGBA
		BreakpointCurrHoverColor   color.NRGBA
		BreakpointCurrPressedColor color.NRGBA
	}

	Toolbar struct {
		BackgroundColor   color.NRGBA
		TextColor         color.NRGBA
		MenuColor         color.NRGBA
		MenuHoverColor    color.NRGBA
		MenuPressedColor  color.NRGBA
		EntryColor        color.NRGBA
		EntryHoverColor   color.NRGBA
		EntryPressedColor color.NRGBA
	}

	Slider struct {
		TrackColor color.NRGBA
	}

	Window struct {
		Color         color.NRGBA
		TitleBarColor color.NRGBA
	}

	Padding int
	Insets  *widget.Insets
}

var Debugger = DebuggerTheme{
	BackgroundColor: color.NRGBA{R: 20, G: 22, B: 30, A: 255},
	LabelColor:      color.NRGBA{R: 220, G: 220, B: 230, A: 255},
	TitleColor:      color.NRGBA{R: 255, G: 200, B: 80, A: 255},
	Button: struct {
		Image     *widget.ButtonImage
		TextColor *widget.ButtonTextColor
	}{
		Image: &widget.ButtonImage{
			Idle:    image.NewNineSliceColor(color.NRGBA{R: 60, G: 90, B: 130, A: 255}),
			Hover:   image.NewNineSliceColor(color.NRGBA{R: 80, G: 110, B: 160, A: 255}),
			Pressed: image.NewNineSliceColor(color.NRGBA{R: 100, G: 130, B: 190, A: 255}),
		},
		TextColor: &widget.ButtonTextColor{
			Idle: color.NRGBA{R: 235, G: 235, B: 245, A: 255},
		},
	},
	Main: struct {
		Color        color.NRGBA
		HoverColor   color.NRGBA
		PressedColor color.NRGBA
	}{
		Color:        color.NRGBA{R: 35, G: 30, B: 45, A: 255},
		HoverColor:   color.NRGBA{R: 50, G: 45, B: 65, A: 255},
		PressedColor: color.NRGBA{R: 65, G: 60, B: 85, A: 255},
	},
	Disassembler: struct {
		BreakpointColor            color.NRGBA
		BreakpointHoverColor       color.NRGBA
		BreakpointPressedColor     color.NRGBA
		CurrInstrColor             color.NRGBA
		CurrInstrHoverColor        color.NRGBA
		CurrInstrPressedColor      color.NRGBA
		BreakpointCurrColor        color.NRGBA
		BreakpointCurrHoverColor   color.NRGBA
		BreakpointCurrPressedColor color.NRGBA
	}{
		BreakpointColor:            color.NRGBA{R: 140, G: 40, B: 40, A: 255},
		BreakpointHoverColor:       color.NRGBA{R: 170, G: 60, B: 60, A: 255},
		BreakpointPressedColor:     color.NRGBA{R: 200, G: 80, B: 80, A: 255},
		CurrInstrColor:             color.NRGBA{R: 40, G: 120, B: 110, A: 255},
		CurrInstrHoverColor:        color.NRGBA{R: 60, G: 150, B: 140, A: 255},
		CurrInstrPressedColor:      color.NRGBA{R: 80, G: 180, B: 170, A: 255},
		BreakpointCurrColor:        color.NRGBA{R: 160, G: 60, B: 160, A: 255},
		BreakpointCurrHoverColor:   color.NRGBA{R: 180, G: 100, B: 180, A: 255},
		BreakpointCurrPressedColor: color.NRGBA{R: 200, G: 140, B: 200, A: 255},
	},
	Toolbar: struct {
		BackgroundColor   color.NRGBA
		TextColor         color.NRGBA
		MenuColor         color.NRGBA
		MenuHoverColor    color.NRGBA
		MenuPressedColor  color.NRGBA
		EntryColor        color.NRGBA
		EntryHoverColor   color.NRGBA
		EntryPressedColor color.NRGBA
	}{
		BackgroundColor:   color.NRGBA{R: 18, G: 20, B: 30, A: 255},
		TextColor:         color.NRGBA{R: 230, G: 230, B: 240, A: 255},
		MenuColor:         color.NRGBA{R: 30, G: 30, B: 40, A: 255},
		MenuHoverColor:    color.NRGBA{R: 40, G: 40, B: 50, A: 255},
		MenuPressedColor:  color.NRGBA{R: 50, G: 50, B: 60, A: 255},
		EntryColor:        color.NRGBA{R: 30, G: 30, B: 40, A: 255},
		EntryHoverColor:   color.NRGBA{R: 40, G: 40, B: 50, A: 255},
		EntryPressedColor: color.NRGBA{R: 50, G: 50, B: 60, A: 255},
	},
	Slider: struct {
		TrackColor color.NRGBA
	}{
		TrackColor: color.NRGBA{R: 255, G: 255, B: 255, A: 48},
	},
	Window: struct {
		Color         color.NRGBA
		TitleBarColor color.NRGBA
	}{
		Color:         color.NRGBA{R: 28, G: 30, B: 42, A: 255},
		TitleBarColor: color.NRGBA{R: 22, G: 24, B: 34, A: 255},
	},
	Padding: 5,
	Insets:  widget.NewInsetsSimple(5),
}
