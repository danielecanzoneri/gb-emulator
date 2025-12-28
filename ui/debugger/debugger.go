package debugger

import (
	"github.com/danielecanzoneri/lucky-boy/gameboy"
	"github.com/danielecanzoneri/lucky-boy/ui/graphics"
	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
)

type Debugger struct {
	*ebitenui.UI

	// Widgets
	toolbar         *toolbar
	disassembler    *disassembler
	screen          *screen
	memoryViewer    *memoryViewer
	registersViewer *registersViewer

	oamViewer   *oamViewer
	bgViewer    *bgViewer
	tilesViewer *tilesViewer

	// State
	gameBoy *gameboy.GameBoy
	Active  bool
	Running bool // True when debugger is active and we are stepping until breakpoint

	// Run until next instruction
	NextInstruction bool
	CallDepth       int
}

func New(gb *gameboy.GameBoy) *Debugger {
	// Misc
	font = loadFont(16)

	// Main container
	root := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(theme.Debugger.BackgroundColor)),
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
		)),
	)

	d := &Debugger{
		UI:      &ebitenui.UI{Container: root},
		gameBoy: gb,
	}

	// Create widgets
	d.toolbar = d.newToolbar()
	d.disassembler = newDisassembler()
	d.screen = newScreen()
	d.memoryViewer = newMemoryViewer()
	d.registersViewer = newRegisterViewer()

	d.oamViewer = d.newOamViewer()
	d.bgViewer = d.newBGViewer()
	d.tilesViewer = d.newTilesViewer()

	// Add widgets to the root container
	main := newContainer(widget.DirectionHorizontal,
		d.disassembler,
		newContainer(widget.DirectionVertical,
			newContainer(widget.DirectionHorizontal,
				d.screen, d.registersViewer,
			),
			d.memoryViewer,
		),
	)
	root.AddChild(d.toolbar, main)
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
	d.UI.Update()
	return nil
}

func (d *Debugger) Draw(screen *ebiten.Image, frame *ebiten.Image) {
	d.screen.Sync(frame)
	d.UI.Draw(screen)
}

func (d *Debugger) Layout(_, _ int) (int, int) {
	return d.UI.Container.PreferredSize()
}
