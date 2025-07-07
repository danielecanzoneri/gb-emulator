package debugger

import (
	"github.com/danielecanzoneri/gb-emulator/emulator/internal/gameboy"
	"image/color"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
)

var backgroundColor = color.NRGBA{R: 0x13, G: 0x1A, B: 0x22, A: 0xFF}

type Debugger struct {
	ui *ebitenui.UI

	// Widgets
	disassembler    *disassembler
	screen          *screen
	memoryViewer    *memoryViewer
	registersViewer *registersViewer
}

func New() *Debugger {
	// Misc
	font = loadFont(16)

	d := new(Debugger)
	d.disassembler = newDisassembler()
	d.screen = newScreen()
	d.memoryViewer = newMemoryViewer()
	d.registersViewer = newRegisterViewer()

	root := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(backgroundColor)),
		widget.ContainerOpts.Layout(widget.NewRowLayout()),
	)

	// Add widgets to the root container
	root.AddChild(
		d.disassembler,
		newContainer(widget.DirectionVertical,
			newContainer(widget.DirectionHorizontal,
				d.screen, d.registersViewer,
			),
			d.memoryViewer,
		),
	)

	d.ui = &ebitenui.UI{
		Container: root,
	}
	return d
}

// Sync state between game boy and debugger
func (d *Debugger) Sync(gb *gameboy.GameBoy) {
	d.disassembler.Sync(gb)
	d.memoryViewer.Sync(gb)
	d.registersViewer.Sync(gb)
}

func (d *Debugger) Update() error {
	d.ui.Update()
	return nil
}

func (d *Debugger) Draw(screen *ebiten.Image, frame *ebiten.Image) {
	d.screen.Sync(frame)
	d.ui.Draw(screen)
}

func (d *Debugger) Layout(_, _ int) (int, int) {
	return d.ui.Container.PreferredSize()
}
