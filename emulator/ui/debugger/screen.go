package debugger

import (
	"github.com/danielecanzoneri/gb-emulator/emulator/internal/gameboy/ppu"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
)

const screenScale = 2

type screen struct {
	*widget.Graphic
}

func (s *screen) Sync(image *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(screenScale, screenScale)
	s.Image.DrawImage(image, op)
}

func newScreen() *screen {
	image := ebiten.NewImage(ppu.FrameWidth*screenScale, ppu.FrameHeight*screenScale)
	s := new(screen)
	s.Graphic = widget.NewGraphic(
		widget.GraphicOpts.Image(image),
	)
	return s
}
