package debugger

import (
	"fmt"
	"github.com/danielecanzoneri/gb-emulator/gameboy"
	"github.com/danielecanzoneri/gb-emulator/gameboy/ppu"
	"github.com/danielecanzoneri/gb-emulator/ui/theme"
	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"image/color"
)

const (
	tileScale = 2
)

type bgTile struct {
	*widget.Container

	// Row, columns
	row, col  uint16
	tileId    uint8
	attribute uint8

	// Image displaying the object
	sprite      *ebiten.Image
	drawOptions *ebiten.DrawImageOptions
	graphic     *widget.Graphic

	address uint16
}

type syncTileFunc func(col, row uint16, tileId, attribute uint8)

func newBGTile(row, col uint16, syncData syncTileFunc) *bgTile {
	tile := &bgTile{row: row, col: col}

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
				syncData(tile.row, tile.col, tile.tileId, tile.attribute)
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

func (t *bgTile) Sync(gb *gameboy.GameBoy) {
	t.address = gb.PPU.BGTileMapAddr + (t.row * 32) + t.col
	t.tileId = gb.PPU.GetTileId(t.address - 0x9800)

	// Update image
	for row := range 8 {
		pixels, attr := gb.PPU.GetBGWindowPixelRow(t.address, uint8(row))
		for col := range 8 {
			if gb.EmulationModel == gameboy.CGB {
				paletteId := attr.CGBPalette()
				p := ppu.CGBPalette(gb.PPU.OBJPalette[8*paletteId : 8*paletteId+8])
				t.sprite.Set(col, row, theme.CGBPalette{}.Get(p.GetColor(pixels[col])))
			} else {
				t.sprite.Set(col, row, theme.DMGPalette{}.Get(uint16(pixels[col])))
			}
		}
		t.attribute = uint8(attr)
	}
	t.graphic.Image.DrawImage(t.sprite, t.drawOptions)
}

type bgViewer struct {
	*widget.Window

	// Pointer to the UI for showing the window
	ui *ebitenui.UI

	yLabel         *widget.Text
	xLabel         *widget.Text
	tileLabel      *widget.Text
	attributeLabel *widget.Text

	tiles [32][32]*bgTile

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
	v.Window = newWindow("BG  Viewer", root, &v.closeWindow)
	return v
}

func (v *bgViewer) Sync(gb *gameboy.GameBoy) {
	if !v.ui.IsWindowOpen(v.Window) {
		return
	}

	for _, row := range v.tiles {
		for _, tile := range row {
			tile.Sync(gb)
		}
	}
}
