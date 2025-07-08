package debugger

import (
	"github.com/danielecanzoneri/gb-emulator/gameboy"
	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
)

type Debugger struct {
	ui *ebitenui.UI

	// Widgets
	disassembler    *disassembler
	screen          *screen
	memoryViewer    *memoryViewer
	registersViewer *registersViewer

	// State
	gameBoy   *gameboy.GameBoy
	Active    bool
	Continued bool // True when debugger is active and we are stepping until breakpoint
}

func New(gb *gameboy.GameBoy) *Debugger {
	// Misc
	font = loadFont(16)

	d := &Debugger{
		gameBoy: gb,
	}
	d.disassembler = d.newDisassembler()
	d.screen = newScreen()
	d.memoryViewer = newMemoryViewer()
	d.registersViewer = newRegisterViewer()

	root := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(mainColor)),
		widget.ContainerOpts.Layout(widget.NewRowLayout()),
	)
	d.ui = &ebitenui.UI{
		Container: root,
	}

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
	return d
}

// Sync state between game boy and debugger
func (d *Debugger) Sync() {
	d.disassembler.Sync(d.gameBoy)
	d.memoryViewer.Sync(d.gameBoy)
	d.registersViewer.Sync(d.gameBoy)
}

func (d *Debugger) Update() error {
	d.registersViewer.Sync(d.gameBoy)
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
