package debugger

import (
	"fmt"
	"github.com/danielecanzoneri/gb-emulator/gameboy/ppu"

	"github.com/danielecanzoneri/gb-emulator/gameboy"
	"github.com/danielecanzoneri/gb-emulator/ui/theme"
	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/widget"
)

var (
	basicCGBPalette = ppu.CGBPalette{
		0xFE, 0xFF, // Black
		0x94, 0x52, // Dark Gray
		0x4A, 0x29, // Light Gray
		0x00, 0x00, // White
	}
	basicDMGPalette = ppu.DMGPalette(0xE4)
)

type tileData struct {
	*tileView

	// Bank, address
	bank    uint8
	address uint16
}

func newTileData(bank uint8, address uint16, syncTile func(bank uint8, address uint16)) *tileData {
	tile := &tileData{
		tileView: newTileView(tileScale, func() {
			syncTile(bank, address)
		}),
		bank:    bank,
		address: address,
	}

	return tile
}

func (t *tileData) Sync(gb *gameboy.GameBoy) {
	tileOffset := (t.address - 0x8000) >> 4

	var systemPalette theme.Palette = theme.DMGPalette{}
	var colorPalette ppu.Palette = basicDMGPalette
	if gb.EmulationModel == gameboy.CGB {
		systemPalette = theme.CGBPalette{}
		colorPalette = basicCGBPalette
	}

	// Render tile using shared buffer
	tile := gb.PPU.DebugGetTileData(t.bank, tileOffset)
	var pixels [8][8]uint8
	for row := range 8 {
		pixels[row] = tile.GetRow(0, uint8(row))
	}

	// Use common rendering method (DMG palette only for tiles viewer)
	t.renderPixels(pixels, systemPalette, colorPalette)
}

type tilesViewer struct {
	// Pointer to the UI for showing the window
	ui *ebitenui.UI

	bankLabel    *widget.Text
	addressLabel *widget.Text

	tiles [2][384]*tileData

	// Window info
	windowInfo *windowInfo

	// Handler to close the window
	closeWindow widget.RemoveWindowFunc
}

func (d *Debugger) newTilesViewer() *tilesViewer {
	v := &tilesViewer{ui: d.UI}

	// Tile data
	v.bankLabel = newLabel("0", theme.Debugger.LabelColor)
	v.addressLabel = newLabel("0000", theme.Debugger.LabelColor)

	tileContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(2),
			widget.GridLayoutOpts.Padding(theme.Debugger.Insets),
			widget.GridLayoutOpts.Spacing(theme.Debugger.Padding, theme.Debugger.Padding),
		)),
	)
	tileContainer.AddChild(
		newLabel("Bank", theme.Debugger.TitleColor), v.bankLabel,
		newLabel("Address", theme.Debugger.TitleColor), v.addressLabel,
	)

	root := newContainer(widget.DirectionHorizontal)

	syncTile := func(bank uint8, address uint16) {
		v.bankLabel.Label = fmt.Sprintf("%d", bank)
		v.addressLabel.Label = fmt.Sprintf("%04X", address)
	}

	for bank := range 2 {
		grid := widget.NewContainer(
			widget.ContainerOpts.Layout(widget.NewGridLayout(
				widget.GridLayoutOpts.Columns(16), // Display as 16x24 grid
				widget.GridLayoutOpts.Padding(theme.Debugger.Insets),
				// Define how far apart the rows and columns should be
				widget.GridLayoutOpts.Spacing(tileScale, tileScale),
			)),
		)

		// Initialize empty tiles
		for i := range v.tiles[bank] {
			addr := uint16(i<<4) + 0x8000
			v.tiles[bank][i] = newTileData(uint8(bank), addr, syncTile)
			grid.AddChild(v.tiles[bank][i])
		}

		root.AddChild(grid)
	}

	root.AddChild(tileContainer)

	v.windowInfo = newWindow("Tiles Viewer", root, &v.closeWindow)
	return v
}

func (v *tilesViewer) Window() *widget.Window {
	return v.windowInfo.Window
}

func (v *tilesViewer) Contents() *widget.Container {
	return v.windowInfo.Contents
}

func (v *tilesViewer) TitleBar() *widget.Container {
	return v.windowInfo.TitleBar
}

func (v *tilesViewer) SetCloseHandler(closeFunc widget.RemoveWindowFunc) widget.RemoveWindowFunc {
	old := v.closeWindow
	v.closeWindow = closeFunc
	return old
}

func (v *tilesViewer) Sync(gb *gameboy.GameBoy) {
	if !v.ui.IsWindowOpen(v.Window()) {
		return
	}

	for _, tile := range v.tiles[0] {
		tile.Sync(gb)
	}
	for _, tile := range v.tiles[1] {
		if gb.PPU.Cgb && !gb.PPU.DmgCompatibility {
			tile.Sync(gb)
		} else {
			tile.clear()
		}
	}
}
