package debugger

import (
	"fmt"
	"github.com/danielecanzoneri/gb-emulator/gameboy"
	"github.com/danielecanzoneri/gb-emulator/ui/theme"
	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"image/color"
)

type tileData struct {
	*widget.Container

	// Bank, address
	bank    uint8
	address uint16

	// Image displaying the object
	sprite      *ebiten.Image
	drawOptions *ebiten.DrawImageOptions
	graphic     *widget.Graphic
}

type syncTileDataFunc func(bank uint8, address uint16)

func newTileData(bank uint8, address uint16, syncData syncTileDataFunc) *tileData {
	tile := &tileData{bank: bank, address: address}

	// Object image
	tile.sprite = ebiten.NewImage(8, 8)
	tile.sprite.Fill(color.Transparent)
	tile.drawOptions = &ebiten.DrawImageOptions{}
	tile.drawOptions.GeoM.Scale(tileScale, tileScale)

	scaledSprite := ebiten.NewImage(8*tileScale, 8*tileScale)
	scaledSprite.DrawImage(tile.sprite, tile.drawOptions)
	tile.graphic = widget.NewGraphic(
		widget.GraphicOpts.Image(scaledSprite),
		widget.GraphicOpts.WidgetOpts(
			widget.WidgetOpts.CursorEnterHandler(func(args *widget.WidgetCursorEnterEventArgs) {
				syncData(tile.bank, tile.address)
			}),
		),
	)

	tile.Container = widget.NewContainer(widget.ContainerOpts.Layout(
		widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionHorizontal),
			widget.RowLayoutOpts.Spacing(2),
		),
	))
	tile.Container.AddChild(tile.graphic)

	return tile
}

func (t *tileData) Sync(gb *gameboy.GameBoy) {
	tileOffset := (t.address - 0x8000) >> 4

	t.sprite.Fill(color.Transparent)

	// Update image
	for row := range 8 {
		tile := gb.PPU.DebugGetTileData(t.bank, tileOffset)
		pixels := tile.GetRow(0, uint8(row))

		for col := range 8 {
			t.sprite.Set(col, row, theme.DMGPalette{}.Get(uint16(pixels[col])))
		}
	}
	t.graphic.Image.DrawImage(t.sprite, t.drawOptions)
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

	for _, row := range v.tiles {
		for _, tile := range row {
			tile.Sync(gb)
		}
	}
}
