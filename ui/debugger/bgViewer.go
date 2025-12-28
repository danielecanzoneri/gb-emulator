package debugger

import (
	"fmt"

	"github.com/danielecanzoneri/lucky-boy/gameboy"
	"github.com/danielecanzoneri/lucky-boy/gameboy/ppu"
	"github.com/danielecanzoneri/lucky-boy/ui/graphics"
	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/widget"
)

type bgTile struct {
	*tileView

	// Row, columns
	row, col  uint16
	tileId    uint8
	attribute uint8
	address   uint16
}

func newBGTile(row, col uint16, syncData func(col, row uint16, tileId, attribute uint8)) *bgTile {
	tile := &bgTile{
		row: row,
		col: col,
	}

	// Create tile view with hover callback (must be after tile is created)
	tile.tileView = newTileView(tileScale, func() {
		syncData(row, col, tile.tileId, tile.attribute)
	})

	return tile
}

func (t *bgTile) Sync(gb *gameboy.GameBoy) {
	t.address = gb.PPU.DebugGetBGTileMapAddr() + (t.row * 32) + t.col
	t.tileId = gb.PPU.GetTileId(t.address - 0x9800)

	var systemPalette theme.Palette = theme.DMGPalette{}
	var colorPalette ppu.Palette = gb.PPU.BGP
	if gb.EmulationModel == gameboy.CGB {
		systemPalette = theme.CGBPalette{}

		bgPalette := gb.PPU.DebugGetBGPalette()
		if gb.PPU.DmgCompatibility {
			colorPalette = gb.PPU.BGP.ConvertToCGB(bgPalette[0:8])
		} else {
			t.attribute = gb.PPU.DebugGetTileMaps(1, t.address-0x9800)
			paletteId := ppu.TileAttribute(t.attribute).CGBPalette()
			colorPalette = ppu.CGBPalette(bgPalette[8*paletteId : 8*paletteId+8])
		}
	}

	// Render tile row by row using shared buffer
	var pixels [8][8]uint8
	for row := range 8 {
		pixels[row], _ = gb.PPU.GetBGWindowPixelRow(t.address, uint8(row))
	}

	// Use common rendering method
	t.renderPixels(pixels, systemPalette, colorPalette)
}

type bgViewer struct {
	// Pointer to the UI for showing the window
	ui *ebitenui.UI

	yLabel         *widget.Text
	xLabel         *widget.Text
	tileLabel      *widget.Text
	attributeLabel *widget.Text

	tiles [32][32]*bgTile

	// Window info
	windowInfo *windowInfo

	// Handler to close the window
	closeWindow widget.RemoveWindowFunc
}

func (d *Debugger) newBGViewer() *bgViewer {
	v := &bgViewer{ui: d.UI}

	// Tile data
	v.yLabel = newLabel("00", theme.Debugger.LabelColor)
	v.xLabel = newLabel("00", theme.Debugger.LabelColor)
	v.tileLabel = newLabel("00", theme.Debugger.LabelColor)
	v.attributeLabel = newLabel("00", theme.Debugger.LabelColor)

	tileContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(2),
			widget.GridLayoutOpts.Padding(theme.Debugger.Insets),
			widget.GridLayoutOpts.Spacing(theme.Debugger.Padding, theme.Debugger.Padding),
		)),
	)
	tileContainer.AddChild(
		newLabel("X", theme.Debugger.TitleColor), v.xLabel,
		newLabel("Y", theme.Debugger.TitleColor), v.yLabel,
		newLabel("Tile No", theme.Debugger.TitleColor), v.tileLabel,
		newLabel("Attribute", theme.Debugger.TitleColor), v.attributeLabel,
	)

	grid := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(32), // Display as 32x32 grid
			widget.GridLayoutOpts.Padding(theme.Debugger.Insets),
			// Define how far apart the rows and columns should be
			widget.GridLayoutOpts.Spacing(tileScale, tileScale),
		)),
	)

	syncTile := func(col, row uint16, tileId, attribute uint8) {
		v.yLabel.Label = fmt.Sprintf("%02X", col)
		v.xLabel.Label = fmt.Sprintf("%02X", row)
		v.tileLabel.Label = fmt.Sprintf("%02X", tileId)
		v.attributeLabel.Label = fmt.Sprintf("%02X", attribute)
	}

	// Initialize empty tiles
	for r := range v.tiles {
		for c := range v.tiles[r] {
			v.tiles[r][c] = newBGTile(uint16(r), uint16(c), syncTile)
			grid.AddChild(v.tiles[r][c])
		}
	}

	root := newContainer(widget.DirectionHorizontal,
		grid, tileContainer,
	)
	v.windowInfo = newWindow("BG  Viewer", root, &v.closeWindow)
	return v
}

func (v *bgViewer) Window() *widget.Window {
	return v.windowInfo.Window
}

func (v *bgViewer) Contents() *widget.Container {
	return v.windowInfo.Contents
}

func (v *bgViewer) TitleBar() *widget.Container {
	return v.windowInfo.TitleBar
}

func (v *bgViewer) SetCloseHandler(closeFunc widget.RemoveWindowFunc) widget.RemoveWindowFunc {
	old := v.closeWindow
	v.closeWindow = closeFunc
	return old
}

func (v *bgViewer) Sync(gb *gameboy.GameBoy) {
	if !v.ui.IsWindowOpen(v.Window()) {
		return
	}

	for _, row := range v.tiles {
		for _, tile := range row {
			tile.Sync(gb)
		}
	}
}
