package debugger

import (
	"fmt"
	"image/color"

	ebitenimage "github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// DisassemblyViewer represents the opcode visualization widget
type DisassemblyViewer struct {
	mem         MemoryDebugger
	container   *widget.Container
	scrollCont  *widget.ScrollContainer
	labels      []*widget.Label
	startRow    int
	visibleRows int
	face        *text.GoTextFace
	slider      *widget.Slider
	breakpoints map[uint16]bool // Track breakpoints by address
}

// NewDisassemblyViewer creates a new disassembly viewer widget
func NewDisassemblyViewer(mem MemoryDebugger, face *text.GoTextFace) *DisassemblyViewer {
	dv := &DisassemblyViewer{
		mem:         mem,
		face:        face,
		visibleRows: 24,
		labels:      make([]*widget.Label, 24),
		breakpoints: make(map[uint16]bool),
	}

	// Create the container for opcode lines
	dv.container = widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Spacing(2),
		)),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.MinSize(240, 480),
		),
	)

	// Create the scroll container
	dv.scrollCont = widget.NewScrollContainer(
		widget.ScrollContainerOpts.Content(dv.container),
		widget.ScrollContainerOpts.StretchContentWidth(),
		widget.ScrollContainerOpts.Image(&widget.ScrollContainerImage{
			Idle: ebitenimage.NewNineSliceColor(color.RGBA{60, 60, 60, 255}),
			Mask: ebitenimage.NewNineSliceColor(color.RGBA{60, 60, 60, 255}),
		}),
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
		dv.container.AddChild(label)
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
			widget.WidgetOpts.MinSize(15, 0),
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

	// Initial update
	dv.Update()

	return dv
}

// GetWidget returns the root widget of the disassembly viewer
func (dv *DisassemblyViewer) GetWidget() widget.PreferredSizeLocateableWidget {
	return dv.scrollCont
}

// GetSlider returns the slider widget
func (dv *DisassemblyViewer) GetSlider() widget.PreferredSizeLocateableWidget {
	return dv.slider
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
}
