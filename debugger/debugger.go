package debugger

import (
	"bytes"
	"image/color"
	"log"

	"golang.org/x/image/font/gofont/gomono"

	"github.com/ebitenui/ebitenui"
	ebitenimage "github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type MemoryDebugger interface {
	DebugRead(uint16) uint8
}

type CPUDebugger interface {
	ReadAF() uint16
	ReadBC() uint16
	ReadDE() uint16
	ReadHL() uint16

	ReadPC() uint16
	ReadSP() uint16

	InterruptEnabled() bool
}

// Debugger is a component that displays Game Boy I/O registers
type Debugger struct {
	mem     MemoryDebugger
	cpu     CPUDebugger
	visible bool
	face    *text.GoTextFace

	ui       *ebitenui.UI
	rootCont *widget.Container

	// Gameboy screen placeholder
	GameboyScreen *widget.Container

	// Memory viewer
	MemViewer *MemoryViewer

	// Disassembly viewer
	DisViewer *DisassemblyViewer

	// updateHandlers is a map of functions that are called to update the contents of the panels
	updateHandlers map[string]func()
}

// NewDebugger creates a new I/O registers debugger
func NewDebugger(mem MemoryDebugger, cpu CPUDebugger) *Debugger {
	d := &Debugger{
		mem:            mem,
		cpu:            cpu,
		visible:        false,
		updateHandlers: make(map[string]func()),
	}

	// Load the font
	fontData := bytes.NewReader(gomono.TTF)
	s, err := text.NewGoTextFaceSource(fontData)
	if err != nil {
		log.Fatal(err)
	}
	d.face = &text.GoTextFace{
		Source: s,
		Size:   16,
	}

	d.initUI()

	return d
}

// ToggleVisibility toggles the visibility of the debug panel
func (d *Debugger) ToggleVisibility() {
	d.visible = !d.visible
}

// IsVisible returns true if the debug panel is visible
func (d *Debugger) IsVisible() bool {
	return d.visible
}

func (d *Debugger) initUI() {
	screenRow := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionHorizontal),
			widget.RowLayoutOpts.Spacing(10),
		)),
	)

	// Add disassembly viewer
	d.DisViewer = NewDisassemblyViewer(d.mem, d.face)
	screenRow.AddChild(d.DisViewer)
	d.updateHandlers["DisassemblyViewer"] = d.DisViewer.Update

	// Add the gameboy screen
	d.GameboyScreen = widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.MinSize(480, 432),
		),
	)
	screenRow.AddChild(d.GameboyScreen)

	mainCol := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Spacing(10),
		)),
	)
	mainCol.AddChild(screenRow)

	// Add memory viewer
	d.MemViewer = NewMemoryViewer(d.mem, d.face)
	mainCol.AddChild(d.MemViewer)
	d.updateHandlers["MemoryViewer"] = d.MemViewer.Update

	bg := ebitenimage.NewNineSliceColor(color.RGBA{40, 40, 40, 230})
	// Root container
	d.rootCont = widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionHorizontal),
			widget.RowLayoutOpts.Spacing(10),
			widget.RowLayoutOpts.Padding(widget.NewInsetsSimple(10)),
		)),
		widget.ContainerOpts.BackgroundImage(bg),
	)

	// Left column - CPU registers, interrupts, LCD
	leftCol := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Spacing(10),
		)),
	)

	// Create all panels
	d.createPanel(leftCol, "CPU", d.updateCPU)
	d.createPanel(leftCol, "Interrupts", d.updateInterrupts)
	d.createPanel(leftCol, "LCD", d.updateLCD)

	// Right columns
	rightCol := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Spacing(10),
		)),
	)

	// Sound channels
	soundRow := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionHorizontal),
			widget.RowLayoutOpts.Spacing(10),
		)),
	)

	ch12Col := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Spacing(10),
		)),
	)

	d.createPanel(ch12Col, "Square (ch1)", d.updateCh1)
	d.createPanel(ch12Col, "Square (ch2)", d.updateCh2)

	ch34Col := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Spacing(10),
		)),
	)

	d.createPanel(ch34Col, "Wave (ch3)", d.updateCh3)
	d.createPanel(ch34Col, "Noise (ch4)", d.updateCh4)

	soundRow.AddChild(ch12Col)
	soundRow.AddChild(ch34Col)

	// Sound control and wave ram
	soundMiscRow := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionHorizontal),
			widget.RowLayoutOpts.Spacing(10),
		)),
	)

	d.createPanel(soundMiscRow, "Sound Control", d.updateSoundControl)
	d.createPanel(soundMiscRow, "WaveRam", d.updateWaveRam)

	// Opcode and timer
	opcodeTimerRow := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionHorizontal),
			widget.RowLayoutOpts.Spacing(10),
		)),
	)

	d.createPanel(opcodeTimerRow, "Opcode", d.updateOpcode)
	d.createPanel(opcodeTimerRow, "Timer", d.updateTimer)

	// Add sound rows to right column
	rightCol.AddChild(soundRow)
	rightCol.AddChild(soundMiscRow)
	rightCol.AddChild(opcodeTimerRow)

	// Add columns to root
	d.rootCont.AddChild(mainCol)
	d.rootCont.AddChild(leftCol)
	d.rootCont.AddChild(rightCol)

	// Create UI
	d.ui = &ebitenui.UI{
		Container: d.rootCont,
	}

	// Initial update of panel contents
	d.updatePanelContents()
}

func (d *Debugger) createPanel(parent *widget.Container, title string, updateContent func() string) {
	// Panel background
	panelBg := ebitenimage.NewNineSliceColor(color.RGBA{60, 60, 60, 255})

	// Container for the panel with title and content
	panelContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Spacing(5),
			widget.RowLayoutOpts.Padding(widget.NewInsetsSimple(8)),
		)),
		widget.ContainerOpts.BackgroundImage(panelBg),
	)

	// Create a custom label color for the title
	titleColor := &widget.LabelColor{
		Idle: color.RGBA{R: 255, G: 255, B: 200, A: 255},
	}

	// Title label
	titleLabel := widget.NewLabel(
		widget.LabelOpts.Text(title, d.face, titleColor),
	)

	// Content label color
	contentColor := &widget.LabelColor{
		Idle: color.White,
	}

	// Content label
	contentLabel := widget.NewLabel(
		widget.LabelOpts.Text(updateContent(), d.face, contentColor),
	)

	// Store reference to the text label for updating later
	d.updateHandlers[title] = func() {
		contentLabel.Label = updateContent()
	}

	// Add widgets to panel
	panelContainer.AddChild(titleLabel)
	panelContainer.AddChild(contentLabel)

	// Add panel to parent
	parent.AddChild(panelContainer)
}

func (d *Debugger) updatePanelContents() {
	for _, updateHandler := range d.updateHandlers {
		updateHandler()
	}
}
