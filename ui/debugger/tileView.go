package debugger

import (
	"image/color"

	"github.com/danielecanzoneri/gb-emulator/gameboy/ppu"
	"github.com/danielecanzoneri/gb-emulator/ui/theme"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	tileSize     = 8
	tileScale    = 2
	objectsScale = 10
)

// sharedTileRenderer provides shared rendering resources to avoid allocating thousands of images
var sharedTileRenderer = struct {
	// Single shared buffer for rendering tiles (8x8)
	// All tiles render to this first, then scale to their final image
	renderBuffer *ebiten.Image

	// Shared placeholder images for different scales (allocated once)
	placeholder2x  *ebiten.Image // For tileScale
	placeholder10x *ebiten.Image // For objectsScale
}{
	renderBuffer:   ebiten.NewImage(tileSize, tileSize),
	placeholder2x:  ebiten.NewImage(tileSize*tileScale, tileSize*tileScale),
	placeholder10x: ebiten.NewImage(tileSize*objectsScale, tileSize*objectsScale),
}

// tileView is a common structure for displaying 8x8 tiles in debugger viewers
// It uses a shared render buffer to minimize memory allocation
// Images are allocated lazily only when needed
type tileView struct {
	*widget.Container

	// Final scaled image (allocated lazily)
	scaledImage *ebiten.Image
	graphic     *widget.Graphic

	// Scale factor and draw options for this tile
	scale       int
	drawOptions *ebiten.DrawImageOptions

	// Optional callback for hover events
	onHover func()

	// Lazy allocation flag
	initialized bool
}

// newTileView creates a new tile view with the specified scale
// The scaled image is allocated lazily when first rendered
func newTileView(scale int, onHover func()) *tileView {
	tv := &tileView{
		scale:       scale,
		onHover:     onHover,
		drawOptions: &ebiten.DrawImageOptions{},
		initialized: false,
	}

	// Configure draw options for scaling (each tile has its own)
	tv.drawOptions.GeoM.Reset()
	tv.drawOptions.GeoM.Scale(float64(scale), float64(scale))

	// Use shared placeholder image (will be replaced on first render)
	var placeholder *ebiten.Image
	if scale == tileScale {
		placeholder = sharedTileRenderer.placeholder2x
	} else if scale == objectsScale {
		placeholder = sharedTileRenderer.placeholder10x
	} else {
		// Fallback: create placeholder (shouldn't happen with current scales)
		placeholder = ebiten.NewImage(tileSize*scale, tileSize*scale)
		placeholder.Fill(color.Transparent)
	}

	// Create graphic widget
	graphicOpts := []widget.GraphicOpt{
		widget.GraphicOpts.Image(placeholder),
	}
	if onHover != nil {
		graphicOpts = append(graphicOpts, widget.GraphicOpts.WidgetOpts(
			widget.WidgetOpts.CursorEnterHandler(func(args *widget.WidgetCursorEnterEventArgs) {
				onHover()
			}),
		))
	}
	tv.graphic = widget.NewGraphic(graphicOpts...)

	tv.Container = widget.NewContainer(widget.ContainerOpts.Layout(
		widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionHorizontal),
			widget.RowLayoutOpts.Spacing(2),
		),
	))
	tv.Container.AddChild(tv.graphic)

	return tv
}

// ensureInitialized allocates the scaled image if not already done
func (tv *tileView) ensureInitialized() {
	if !tv.initialized {
		tv.scaledImage = ebiten.NewImage(tileSize*tv.scale, tileSize*tv.scale)
		tv.scaledImage.Fill(color.Transparent)
		tv.graphic.Image = tv.scaledImage
		tv.initialized = true
	}
}

// renderPixels renders pixel data to the tile using the shared buffer
// pixels is an array of 8 rows, each with 8 pixel values
func (tv *tileView) renderPixels(pixels [8][8]uint8, systemPalette theme.Palette, colorPalette ppu.Palette) {
	tv.ensureInitialized()
	buf := sharedTileRenderer.renderBuffer

	// Render to shared buffer
	for row := 0; row < 8; row++ {
		for col := 0; col < 8; col++ {
			pixelValue := pixels[row][col]
			c := systemPalette.Get(colorPalette.GetColor(pixelValue))
			buf.Set(col, row, c)
		}
	}

	// Scale from shared buffer to final image
	tv.scaledImage.Clear()
	tv.scaledImage.DrawImage(buf, tv.drawOptions)
}
