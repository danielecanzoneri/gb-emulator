package debugger

import (
	"bytes"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font/gofont/gomono"
	"image/color"
	"log"
)

var font text.Face

// Monospace font
func loadFont(size float64) text.Face {
	s, err := text.NewGoTextFaceSource(bytes.NewReader(gomono.TTF))
	if err != nil {
		log.Fatal(err)
	}

	return &text.GoTextFace{
		Source: s,
		Size:   size,
	}
}

func newContainer(direction widget.Direction, children ...widget.PreferredSizeLocateableWidget) *widget.Container {
	c := widget.NewContainer(widget.ContainerOpts.Layout(
		widget.NewRowLayout(widget.RowLayoutOpts.Direction(
			direction,
		)),
	))
	c.AddChild(children...)
	return c
}

func newLabel(text string, color color.Color) *widget.Text {
	return widget.NewText(
		widget.TextOpts.Text(text, font, color),
	)
}
