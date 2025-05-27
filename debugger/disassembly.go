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

	// Track which memory addresses contain executable code vs data
	// Updated when debugger first shown and during instruction execution
	codeAddresses [0x10000]bool

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

	// Set all code addresses to true
	for i := 0; i < 0x10000; i++ {
		dv.codeAddresses[i] = true
	}

	// Create the container for opcode lines
	dv.scrollArea = widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Spacing(2),
		)),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.MinSize(300, 750),
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
		widget.SliderOpts.MinMax(0, 0x10000-dv.visibleRows), // Full address range
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

// getOpcodeInfo returns information about the opcode at the given address
func (dv *DisassemblyViewer) getOpcodeInfo(addr uint16) (name string, length int, bytes []uint8) {
	opcode := dv.mem.DebugRead(addr)

	opcodeInfo := opcodesInfo[opcode]
	length = opcodeInfo.length
	switch length {
	case 0:
		fallthrough
	case 1:
		name = opcodeInfo.name
		bytes = []uint8{opcode}
		return
	case 2:
		data1 := dv.mem.DebugRead(addr + 1)
		name = opcodeInfo.format(opcodeInfo.name, data1)
		bytes = []uint8{opcode, data1}
		return
	case 3:
		data1 := dv.mem.DebugRead(addr + 1)
		data2 := dv.mem.DebugRead(addr + 2)
		name = opcodeInfo.format(opcodeInfo.name, data1, data2)
		bytes = []uint8{opcode, data1, data2}
		return
	}
	return
}

// formatInstruction formats an instruction with its bytes and arguments
func (dv *DisassemblyViewer) formatInstruction(addr uint16) (string, int) {
	name, length, bytes := dv.getOpcodeInfo(addr)

	// Format the bytes part
	bytesStr := ""
	for _, b := range bytes {
		bytesStr += fmt.Sprintf("%02X ", b)
	}

	// Add padding to align the instruction name
	for len(bytesStr) < 9 { // 3 chars per byte, up to 3 bytes
		bytesStr += "   "
	}

	return fmt.Sprintf("%s  %s", bytesStr, name), length
}

// UpdateCodeAddresses scans through memory and marks which addresses contain
// executable code vs data bytes that are part of multi-byte instructions.
// This is used to properly display the disassembly view.
func (dv *DisassemblyViewer) UpdateCodeAddresses() {
	for addr := 0; addr < 0x10000; addr++ {
		dv.codeAddresses[addr] = true

		opcode := dv.mem.DebugRead(uint16(addr))
		opcodeLen := opcodesInfo[opcode].length
		switch opcodeLen {
		case 2:
			if addr+1 < 0x10000 {
				dv.codeAddresses[addr+1] = false
			}
			addr++
		case 3:
			if addr+1 < 0x10000 {
				dv.codeAddresses[addr+1] = false
			}
			if addr+2 < 0x10000 {
				dv.codeAddresses[addr+2] = false
			}
			addr += 2
		}
	}

	// Set max slider value to the last displayable code address
	dv.slider.Max = int(dv.lastDisplayableCodeAddress())
}

func (dv *DisassemblyViewer) firstValidCodeAddressBefore(addr uint16) uint16 {
	for !dv.codeAddresses[addr] {
		addr--
	}
	return addr
}

func (dv *DisassemblyViewer) firstValidCodeAddressAfter(addr uint16) uint16 {
	for !dv.codeAddresses[addr] {
		addr++
	}
	return addr
}

func (dv *DisassemblyViewer) lastDisplayableCodeAddress() uint16 {
	// We count 0xFFFF as valid address
	codeAddresses := 1
	addr := uint16(0xFFFF)
	for codeAddresses < dv.visibleRows {
		addr--
		if dv.codeAddresses[addr] {
			codeAddresses++
		}
	}
	return addr
}

// Update updates the disassembly viewer's content
func (dv *DisassemblyViewer) Update() {
	// Find the first displayable code address
	for !dv.codeAddresses[dv.slider.Current] {
		dv.slider.Current++
	}
	dv.startRow = dv.slider.Current
	currentAddr := uint16(dv.startRow)

	// Update each row's content
	for i := 0; i < dv.visibleRows; i++ {
		// Create breakpoint indicator
		breakpointIndicator := " "
		if dv.HasBreakpoint(currentAddr) {
			breakpointIndicator = "●" // Unicode bullet as breakpoint indicator
		}

		// Format the instruction with its bytes
		instruction, length := dv.formatInstruction(currentAddr)

		// Format the line with address and instruction
		dv.labels[i].Label = fmt.Sprintf("%04X %s %s", currentAddr, breakpointIndicator, instruction)

		// Move to next instruction
		currentAddr += uint16(length)
	}

	dv.rootContent.Update()
}

// Scroll scrolls the disassembly view up/down
func (dv *DisassemblyViewer) Scroll(xCursor, yCursor int, yWheel float64) {
	rect := dv.scrollArea.GetWidget().Rect
	if yWheel != 0 && image.Pt(xCursor, yCursor).In(rect) {
		newStartRow := dv.startRow - int(yWheel*yWheel*yWheel)
		if newStartRow < 0 {
			newStartRow = 0
		} else if newStartRow > dv.slider.Max {
			newStartRow = dv.slider.Max
		}

		// If we scroll up, we need to find the first valid code address before the current start row
		// If we scroll down, we need to find the first valid code address after the current start row
		if yWheel > 0 {
			newStartRow = int(dv.firstValidCodeAddressBefore(uint16(newStartRow)))
		} else {
			newStartRow = int(dv.firstValidCodeAddressAfter(uint16(newStartRow)))
		}
		dv.startRow = newStartRow
		dv.slider.Current = newStartRow
	}
}

// ToggleBreakpoint toggles a breakpoint at the specified address
func (dv *DisassemblyViewer) ToggleBreakpoint(addr uint16) {
	if dv.breakpoints[addr] {
		delete(dv.breakpoints, addr)
	} else {
		dv.breakpoints[addr] = true
	}
}

// HasBreakpoint checks if there's a breakpoint at the specified address
func (dv *DisassemblyViewer) HasBreakpoint(addr uint16) bool {
	return dv.breakpoints[addr]
}
