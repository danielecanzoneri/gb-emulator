package debugger

import (
	"bytes"
	"image/color"
	"log"

	"github.com/danielecanzoneri/gb-emulator/ui/graphics"

	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font/gofont/gomono"
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

type windowInfo struct {
	Window   *widget.Window
	TitleBar *widget.Container
	Contents *widget.Container
}

func newWindow(title string, content *widget.Container, closeWindow *widget.RemoveWindowFunc) *windowInfo {
	c := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(theme.Debugger.Window.Color)),
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
	)
	c.AddChild(content)

	titleBar := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(theme.Debugger.Window.TitleBarColor)),
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(2),
			widget.GridLayoutOpts.Stretch([]bool{true, false}, []bool{true}),
		)),
	)

	titleBar.AddChild(widget.NewText(
		widget.TextOpts.Text(title, font, theme.Debugger.TitleColor),
		widget.TextOpts.Padding(theme.Debugger.Insets),
		widget.TextOpts.Position(widget.TextPositionCenter, widget.TextPositionStart),
	))

	titleBar.AddChild(widget.NewButton(
		widget.ButtonOpts.Image(theme.Debugger.Button.Image),
		widget.ButtonOpts.TextPadding(theme.Debugger.Insets),
		widget.ButtonOpts.Text("X", font, theme.Debugger.Button.TextColor),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			(*closeWindow)()
		}),
		//widget.ButtonOpts.TabOrder(99),
	))

	w := widget.NewWindow(
		widget.WindowOpts.Modal(),
		widget.WindowOpts.Contents(c),
		widget.WindowOpts.TitleBar(titleBar, 25),
	)

	return &windowInfo{
		Window:   w,
		TitleBar: titleBar,
		Contents: c,
	}
}
