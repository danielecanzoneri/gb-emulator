package debugger

import (
	"fmt"
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

	// Index (from 0 to 39) in OAM data
	index int

	// Image displaying the object
	sprite      *ebiten.Image
	drawOptions *ebiten.DrawImageOptions
	graphic     *widget.Graphic

	xLabel         *widget.Text
	yLabel         *widget.Text
	tileLabel      *widget.Text
	attributeLabel *widget.Text
}

func newOamViewerObject(index int) *oamViewerObject {
	obj := &oamViewerObject{index: index}

	// Object image
	obj.sprite = ebiten.NewImage(8, 8)
	obj.sprite.Fill(theme.GameBoyPalette[0])
	obj.drawOptions = &ebiten.DrawImageOptions{}
	obj.drawOptions.GeoM.Scale(objectsScale, objectsScale)

	scaledSprite := ebiten.NewImage(8*objectsScale, 8*objectsScale)
	scaledSprite.DrawImage(obj.sprite, obj.drawOptions)
	obj.graphic = widget.NewGraphic(
		widget.GraphicOpts.Image(scaledSprite),
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
	obj.Container.AddChild(obj.graphic, dataContainer)

	return obj
}

func (obj *oamViewerObject) Sync(gb *gameboy.GameBoy) {
	oamObj := &gb.PPU.OAM.Data[obj.index]

	// Update data
	obj.xLabel.Label = fmt.Sprintf("%02X", oamObj.Read(0))
	obj.yLabel.Label = fmt.Sprintf("%02X", oamObj.Read(1))
	obj.tileLabel.Label = fmt.Sprintf("%02X", oamObj.Read(2))
	obj.attributeLabel.Label = fmt.Sprintf("%02X", oamObj.Read(3))

	// Update image
	for row := range 8 {
		pixels := gb.PPU.GetObjectRow(oamObj, uint8(row))
		for col := range 8 {
			obj.sprite.Set(col, row, theme.GameBoyPalette[pixels[col]])
		}
	}
	obj.graphic.Image.DrawImage(obj.sprite, obj.drawOptions)
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
		o.objects[i] = newOamViewerObject(i)
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
