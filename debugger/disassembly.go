package debugger

import (
	"fmt"
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"

	ebitenimage "github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

const rowsDisplayed = 32

// DisassemblyViewer represents the opcode visualization widget
type DisassemblyViewer struct {
	labels []*widget.Label
	face   *text.GoTextFace

	mem         MemoryDebugger
	startRow    int
	visibleRows int
	breakpoints map[uint16]bool // Track breakpoints by address

	rootContent *widget.Container
	scrollArea  *widget.Container
	slider      *widget.Slider
}

// initRootContainer initializes the disassembly viewer widget
func (dv *DisassemblyViewer) initRootContainer() {
	// Create container for the disassembly viewer
	dv.rootContent = widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Spacing(4),
		)),
	)

	// Add disassembly view and slider in a horizontal container
	viewContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionHorizontal),
			widget.RowLayoutOpts.Spacing(2),
		)),
	)

	// Add the scroll area and the slider
	viewContainer.AddChild(dv.scrollArea)
	viewContainer.AddChild(dv.slider)

	// Add the title and the view container
	titleColor := &widget.LabelColor{
		Idle: color.RGBA{R: 255, G: 255, B: 200, A: 255},
	}
	titleLabel := widget.NewLabel(
		widget.LabelOpts.Text("Disassembly", dv.face, titleColor),
	)
	dv.rootContent.AddChild(titleLabel)
	dv.rootContent.AddChild(viewContainer)
}

// NewDisassemblyViewer creates a new disassembly viewer widget
func NewDisassemblyViewer(mem MemoryDebugger, face *text.GoTextFace) *DisassemblyViewer {
	dv := &DisassemblyViewer{
		mem:         mem,
		face:        face,
		visibleRows: rowsDisplayed,
		labels:      make([]*widget.Label, rowsDisplayed),
		breakpoints: make(map[uint16]bool),
	}

	// Create the container for opcode lines
	dv.scrollArea = widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Spacing(2),
		)),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.MinSize(240, 750),
		),
	)

	// Create labels for each row
	textColor := &widget.LabelColor{
		Idle: color.White,
	}

	for i := 0; i < dv.visibleRows; i++ {
		label := widget.NewLabel(
			widget.LabelOpts.Text("", face, textColor),
		)
		dv.labels[i] = label
		dv.scrollArea.AddChild(label)
	}

	// Create slider
	dv.slider = widget.NewSlider(
		widget.SliderOpts.Direction(widget.DirectionVertical),
		widget.SliderOpts.MinMax(0, 0xFFFF-dv.visibleRows), // Full address range
		widget.SliderOpts.InitialCurrent(0),
		widget.SliderOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionEnd,
				Stretch:  true,
			}),
			widget.WidgetOpts.MinSize(15, 750),
		),
		widget.SliderOpts.Images(
			&widget.SliderTrackImage{
				Idle:  ebitenimage.NewNineSliceColor(color.NRGBA{100, 100, 100, 255}),
				Hover: ebitenimage.NewNineSliceColor(color.NRGBA{100, 100, 100, 255}),
			},
			&widget.ButtonImage{
				Idle:    ebitenimage.NewNineSliceColor(color.NRGBA{255, 100, 100, 255}),
				Hover:   ebitenimage.NewNineSliceColor(color.NRGBA{255, 100, 100, 255}),
				Pressed: ebitenimage.NewNineSliceColor(color.NRGBA{255, 100, 100, 255}),
			},
		),
		widget.SliderOpts.FixedHandleSize(15),
		widget.SliderOpts.TrackOffset(0),
		widget.SliderOpts.PageSizeFunc(func() int {
			return 24
		}),
	)

	dv.initRootContainer()

	// Initial update
	dv.Update()

	return dv
}

func (dv *DisassemblyViewer) GetWidget() *widget.Widget {
	return dv.rootContent.GetWidget()
}

func (dv *DisassemblyViewer) PreferredSize() (int, int) {
	return dv.rootContent.PreferredSize()
}

func (dv *DisassemblyViewer) SetLocation(rect image.Rectangle) {
	dv.rootContent.SetLocation(rect)
}

func (dv *DisassemblyViewer) Render(screen *ebiten.Image) {
	dv.rootContent.Render(screen)
}

// Update updates the disassembly viewer's content
func (dv *DisassemblyViewer) Update() {
	dv.startRow = dv.slider.Current

	// Update each row's content
	for i := 0; i < dv.visibleRows; i++ {
		addr := uint16(dv.startRow + i)
		opcode := dv.mem.DebugRead(addr)
		opcodeName := opcodes[opcode]

		// Create breakpoint indicator
		breakpointIndicator := " "
		if dv.HasBreakpoint(addr) {
			breakpointIndicator = "●" // Unicode bullet as breakpoint indicator
		}

		// Format the line with address, opcode value, and opcode name
		dv.labels[i].Label = fmt.Sprintf("%04X %s %02X  %s", addr, breakpointIndicator, opcode, opcodeName)
	}

	dv.rootContent.Update()
}

// ToggleBreakpoint toggles a breakpoint at the specified address
func (dv *DisassemblyViewer) ToggleBreakpoint(addr uint16) {
	if dv.breakpoints[addr] {
		delete(dv.breakpoints, addr)
	} else {
		dv.breakpoints[addr] = true
	}
	dv.Update() // Refresh display to show/hide breakpoint indicator
}

// HasBreakpoint checks if there's a breakpoint at the specified address
func (dv *DisassemblyViewer) HasBreakpoint(addr uint16) bool {
	return dv.breakpoints[addr]
}
