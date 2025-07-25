package debugger

import (
	"github.com/danielecanzoneri/gb-emulator/gameboy/ppu"
	"github.com/danielecanzoneri/gb-emulator/ui/theme"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"image/color"
)

const (
	screenScale = 2
)

type screen struct {
	*widget.Graphic
}

func (s *screen) Sync(image *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(screenScale, screenScale)
	op.GeoM.Translate(float64(theme.Debugger.Padding), 0)
	s.Image.DrawImage(image, op)
}

func newScreen() *screen {
	image := ebiten.NewImage(ppu.FrameWidth*screenScale+theme.Debugger.Padding*2, ppu.FrameHeight*screenScale)
	image.Fill(color.Transparent)

	s := new(screen)
	s.Graphic = widget.NewGraphic(
		widget.GraphicOpts.Image(image),
	)
	return s
}
