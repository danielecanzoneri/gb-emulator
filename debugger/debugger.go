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

	if d.visible {
		d.DisViewer.UpdateCodeAddresses()
	}
}

// IsVisible returns true if the debug panel is visible
func (d *Debugger) IsVisible() bool {
	return d.visible
}

func (d *Debugger) initUI() {
	// Create main container with horizontal layout
	mainRow := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionHorizontal),
			widget.RowLayoutOpts.Spacing(10),
		)),
	)

	// Left side - Disassembly viewer
	d.DisViewer = NewDisassemblyViewer(d.mem, d.face)
	mainRow.AddChild(d.DisViewer)
	d.updateHandlers["DisassemblyViewer"] = d.DisViewer.Update

	// Right side container
	rightSide := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Spacing(10),
		)),
	)

	// Top row container for game screen and panels
	topRow := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionHorizontal),
			widget.RowLayoutOpts.Spacing(10),
		)),
	)

	// Add the gameboy screen
	d.GameboyScreen = widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.MinSize(480, 432),
		),
	)
	topRow.AddChild(d.GameboyScreen)

	// Right panels container
	rightPanels := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionHorizontal),
			widget.RowLayoutOpts.Spacing(10),
		)),
	)

	// Left panels (CPU and LCD)
	leftPanels := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Spacing(10),
		)),
	)

	// Create CPU and LCD panels
	d.createPanel(leftPanels, "CPU", d.updateCPU)
	d.createPanel(leftPanels, "LCD", d.updateLCD)

	// Sound panels container
	soundPanels := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionHorizontal),
			widget.RowLayoutOpts.Spacing(10),
		)),
	)

	// First column of sound panels (ch1, ch3, sound control)
	soundCol1 := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Spacing(10),
		)),
	)

	d.createPanel(soundCol1, "Square (ch1)", d.updateCh1)
	d.createPanel(soundCol1, "Wave (ch3)", d.updateCh3)
	d.createPanel(soundCol1, "Sound Control", d.updateSoundControl)

	// Second column of sound panels (ch2, ch4, wave ram)
	soundCol2 := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Spacing(10),
		)),
	)

	d.createPanel(soundCol2, "Square (ch2)", d.updateCh2)
	d.createPanel(soundCol2, "Noise (ch4)", d.updateCh4)
	d.createPanel(soundCol2, "WaveRam", d.updateWaveRam)

	// Add sound columns to sound panels
	soundPanels.AddChild(soundCol1)
	soundPanels.AddChild(soundCol2)

	// Add left panels and sound panels to right panels container
	rightPanels.AddChild(leftPanels)
	rightPanels.AddChild(soundPanels)

	// Add panels to top row
	topRow.AddChild(rightPanels)

	// Add top row to right side
	rightSide.AddChild(topRow)

	// Add memory viewer at the bottom of right side
	d.MemViewer = NewMemoryViewer(d.mem, d.face)
	rightSide.AddChild(d.MemViewer)
	d.updateHandlers["MemoryViewer"] = d.MemViewer.Update

	// Add both main columns to main row
	mainRow.AddChild(rightSide)

	// Background for the root container
	bg := ebitenimage.NewNineSliceColor(color.RGBA{40, 40, 40, 230})

	// Root container with padding
	d.rootCont = widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Padding(widget.NewInsetsSimple(10)),
		)),
		widget.ContainerOpts.BackgroundImage(bg),
	)

	d.rootCont.AddChild(mainRow)

	// TODO - modify
	d.createPanel(soundCol2, "Timer", d.updateTimer)

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
