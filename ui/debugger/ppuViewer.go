package debugger

import (
	"github.com/danielecanzoneri/gb-emulator/gameboy"
	"github.com/danielecanzoneri/gb-emulator/ui/theme"
	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	objectsScale = 10
)

type oamViewerObject struct {
	*widget.Container

	// Image displaying the object
	sprite *ebiten.Image

	xLabel         *widget.Text
	yLabel         *widget.Text
	tileLabel      *widget.Text
	attributeLabel *widget.Text
}

func newOamViewerObject() *oamViewerObject {
	obj := new(oamViewerObject)

	// Object image
	obj.sprite = ebiten.NewImage(8*objectsScale, 8*objectsScale)
	obj.sprite.Fill(theme.GameBoyPalette[0])

	sprite := widget.NewGraphic(
		widget.GraphicOpts.Image(obj.sprite),
	)

	// Object data
	obj.xLabel = newLabel("00", theme.Debugger.LabelColor)
	obj.yLabel = newLabel("00", theme.Debugger.LabelColor)
	obj.tileLabel = newLabel("00", theme.Debugger.LabelColor)
	obj.attributeLabel = newLabel("00", theme.Debugger.LabelColor)

	dataContainer := newContainer(widget.DirectionVertical,
		obj.xLabel, obj.yLabel, obj.tileLabel, obj.attributeLabel,
	)

	obj.Container = widget.NewContainer(widget.ContainerOpts.Layout(
		widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionHorizontal),
			widget.RowLayoutOpts.Spacing(2),
		),
	))
	obj.Container.AddChild(sprite, dataContainer)

	return obj
}

func (obj *oamViewerObject) Sync(gb *gameboy.GameBoy) {
	// TODO
}

type oamViewer struct {
	*widget.Window

	// Pointer to the UI for showing the window
	ui *ebitenui.UI

	objects [40]*oamViewerObject

	// Handler to close the window
	closeWindow widget.RemoveWindowFunc
}

func (d *Debugger) newOamViewer() *oamViewer {
	o := &oamViewer{ui: d.UI}

	root := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(8), // Display as 5x8 grid
			widget.GridLayoutOpts.Padding(theme.Debugger.Insets),
			//Define how far apart the rows and columns should be
			widget.GridLayoutOpts.Spacing(theme.Debugger.Padding, theme.Debugger.Padding*2),
		)),
	)

	// Initialize empty objects
	for i := range o.objects {
		o.objects[i] = newOamViewerObject()
		root.AddChild(o.objects[i])
	}

	o.Window = newWindow("OAM Viewer", root, &o.closeWindow)
	return o
}

func (o *oamViewer) Sync(gb *gameboy.GameBoy) {
	if !o.ui.IsWindowOpen(o.Window) {
		return
	}

	for _, obj := range o.objects {
		obj.Sync(gb)
	}
}
