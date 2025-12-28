package debugger

import (
	"fmt"

	"github.com/danielecanzoneri/gb-emulator/gameboy"
	"github.com/danielecanzoneri/gb-emulator/gameboy/ppu"
	"github.com/danielecanzoneri/gb-emulator/ui/graphics"
	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/widget"
)

type oamViewerObject struct {
	*tileView

	// Index (from 0 to 39) in OAM data
	index int

	// Labels for object data
	yLabel         *widget.Text
	xLabel         *widget.Text
	tileLabel      *widget.Text
	attributeLabel *widget.Text
}

func newOamViewerObject(index int) *oamViewerObject {
	obj := &oamViewerObject{index: index}

	// Create tile view with larger scale for OAM objects
	obj.tileView = newTileView(objectsScale, nil)

	// Object data labels
	obj.yLabel = newLabel("00", theme.Debugger.LabelColor)
	obj.xLabel = newLabel("00", theme.Debugger.LabelColor)
	obj.tileLabel = newLabel("00", theme.Debugger.LabelColor)
	obj.attributeLabel = newLabel("00", theme.Debugger.LabelColor)

	dataContainer := newContainer(widget.DirectionVertical,
		obj.yLabel, obj.xLabel, obj.tileLabel, obj.attributeLabel,
	)

	// Add labels to container (tileView already has graphic)
	obj.Container.AddChild(dataContainer)

	return obj
}

func (obj *oamViewerObject) Sync(gb *gameboy.GameBoy) {
	oamObj := gb.PPU.DebugGetOAMObject(obj.index)
	if oamObj == nil {
		return
	}

	// Update data labels
	obj.yLabel.Label = fmt.Sprintf("%02X", oamObj.Read(0))
	obj.xLabel.Label = fmt.Sprintf("%02X", oamObj.Read(1))
	obj.tileLabel.Label = fmt.Sprintf("%02X", oamObj.Read(2))
	obj.attributeLabel.Label = fmt.Sprintf("%02X", oamObj.Read(3))

	var systemPalette theme.Palette = theme.DMGPalette{}
	paletteId := ppu.TileAttribute(oamObj.Read(3)).DMGPalette()
	var colorPalette ppu.Palette = gb.PPU.OBP[paletteId]
	if gb.EmulationModel == gameboy.CGB {
		systemPalette = theme.CGBPalette{}

		objPalette := gb.PPU.DebugGetOBJPalette()
		if gb.PPU.DmgCompatibility {
			colorPalette = gb.PPU.OBP[paletteId].ConvertToCGB(objPalette[8*paletteId : 8*paletteId+8])
		} else {
			paletteId = ppu.TileAttribute(oamObj.Read(3)).CGBPalette()
			colorPalette = ppu.CGBPalette(objPalette[8*paletteId : 8*paletteId+8])
		}
	}

	// Render tile row by row using shared buffer
	var pixels [8][8]uint8
	for row := range 8 {
		pixels[row] = gb.PPU.GetObjectRow(oamObj, uint8(row))
	}

	// Use common rendering method
	obj.renderPixels(pixels, systemPalette, colorPalette)
}

type oamViewer struct {
	// Pointer to the UI for showing the window
	ui *ebitenui.UI

	objects [40]*oamViewerObject

	// Window info
	windowInfo *windowInfo

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

	o.windowInfo = newWindow("OAM Viewer", root, &o.closeWindow)
	return o
}

func (o *oamViewer) Window() *widget.Window {
	return o.windowInfo.Window
}

func (o *oamViewer) Contents() *widget.Container {
	return o.windowInfo.Contents
}

func (o *oamViewer) TitleBar() *widget.Container {
	return o.windowInfo.TitleBar
}

func (o *oamViewer) SetCloseHandler(closeFunc widget.RemoveWindowFunc) widget.RemoveWindowFunc {
	old := o.closeWindow
	o.closeWindow = closeFunc
	return old
}

func (o *oamViewer) Sync(gb *gameboy.GameBoy) {
	if !o.ui.IsWindowOpen(o.Window()) {
		return
	}

	for _, obj := range o.objects {
		obj.Sync(gb)
	}
}
