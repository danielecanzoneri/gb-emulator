package debugger

import (
	"github.com/danielecanzoneri/gb-emulator/emulator/internal/gameboy/ppu"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
)

const screenScale = 2

type screen struct {
	widget *widget.Graphic
}

func (s *screen) Sync(image *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(screenScale, screenScale)
	s.widget.Image.DrawImage(image, op)
}

func newScreen() *screen {
	image := ebiten.NewImage(ppu.FrameWidth*screenScale, ppu.FrameHeight*screenScale)
	w := widget.NewGraphic(
		widget.GraphicOpts.Image(image),
	)

	return &screen{
		widget: w,
	}
}
