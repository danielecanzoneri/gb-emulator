package debugger

import (
	"github.com/danielecanzoneri/gb-emulator/emulator/internal/gameboy"
	"image/color"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
)

type Debugger struct {
	ui *ebitenui.UI

	disassembler *disassembler
}

func New() *Debugger {
	d := new(Debugger)
	dis := newDisassembler()
	d.disassembler = dis

	// construct a new container that serves as the root of the UI hierarchy
	rootContainer := widget.NewContainer(
		// the container will use a plain color as its background
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{R: 0x13, G: 0x1a, B: 0x22, A: 0xff})),

		// the container will use an anchor layout to lay out its single child widget
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
	)

	// Add list to the root container
	rootContainer.AddChild(dis.widget)
	d.ui = &ebitenui.UI{
		Container: rootContainer,
	}
	return d
}

// Sync state between game boy and debugger
func (d *Debugger) Sync(gb *gameboy.GameBoy) {
	d.disassembler.Sync(gb)
}

func (d *Debugger) Update() error {
	d.ui.Update()
	// _, err := d.ui.Update()
	// return err
	return nil
}

func (d *Debugger) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0x40, 0x40, 0x80, 0xff})
	d.ui.Draw(screen)
}

func (d *Debugger) Layout(_, _ int) (int, int) {
	return 640, 480
}
