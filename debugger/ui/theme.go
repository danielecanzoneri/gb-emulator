package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type gameBoyTheme struct {
}

func (m *gameBoyTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
	case theme.ColorYellow: // Panels title
		return color.RGBA{R: 255, G: 255, B: 191, A: 255}
	default:
		return theme.DefaultTheme().Color(name, variant)
	}
}

func (m *gameBoyTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (m *gameBoyTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (m *gameBoyTheme) Size(name fyne.ThemeSizeName) float32 {
	if name == theme.SizeNamePadding {
		return 0
	}

	return theme.DefaultTheme().Size(name)
}
