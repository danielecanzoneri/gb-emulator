package debugger

import (
	"fmt"
	"github.com/danielecanzoneri/gb-emulator/ui/theme"
	"image/color"

	"github.com/danielecanzoneri/gb-emulator/gameboy"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
)

var (
	entryImage = &widget.ButtonImage{
		Idle:    image.NewNineSliceColor(theme.Debugger.Main.Color),
		Hover:   image.NewNineSliceColor(theme.Debugger.Main.HoverColor),
		Pressed: image.NewNineSliceColor(theme.Debugger.Main.PressedColor),
	}
	entryBreakpointImage = &widget.ButtonImage{
		Idle:    image.NewNineSliceColor(theme.Debugger.Disassembler.BreakpointColor),
		Hover:   image.NewNineSliceColor(theme.Debugger.Disassembler.BreakpointHoverColor),
		Pressed: image.NewNineSliceColor(theme.Debugger.Disassembler.BreakpointPressedColor),
	}
	entryCurrentImage = &widget.ButtonImage{
		Idle:    image.NewNineSliceColor(theme.Debugger.Disassembler.CurrInstrColor),
		Hover:   image.NewNineSliceColor(theme.Debugger.Disassembler.CurrInstrHoverColor),
		Pressed: image.NewNineSliceColor(theme.Debugger.Disassembler.CurrInstrPressedColor),
	}
	entryBreakpointAndCurrentImage = &widget.ButtonImage{
		Idle:    image.NewNineSliceColor(theme.Debugger.Disassembler.BreakpointCurrColor),
		Hover:   image.NewNineSliceColor(theme.Debugger.Disassembler.BreakpointCurrHoverColor),
		Pressed: image.NewNineSliceColor(theme.Debugger.Disassembler.BreakpointCurrPressedColor),
	}
)

// disassemblyEntry represents a single line in the disassembler
type disassemblerEntry struct {
	address uint16
	name    string
	bytes   []uint8
}

type disassembler struct {
	*widget.Container

	slider *widget.Slider

	entries      []*disassemblerEntry
	totalEntries int // Number of actual entries
	rowsWidget   []widget.PreferredSizeLocateableWidget

	// Address of entry to highlight (-1 if no entry has to be highlighted)
	currentInstruction int

	// Entries to show
	first  int
	length int

	// Map with all breakpoints
	breakpoints map[uint16]struct{}
}

// Sync scans through memory and marks which addresses contain
// executable code vs data bytes that are part of multibyte instructions.
// This is used to properly display the disassembly view.
func (d *disassembler) Sync(gb *gameboy.GameBoy) {
	d.currentInstruction = int(gb.CPU.PC)

	counter := 0
	var scrollTo int

	for addr := 0; addr < 0x10000; {
		if 0x104 <= addr && addr < 0x150 { // Header memory
			d.entries[counter].name = "Cart Header"
			d.entries[counter].address = uint16(addr)
			d.entries[counter].bytes = []uint8{gb.Memory.DebugRead(uint16(addr))}
			counter++
			addr++
			continue
		}

		if addr == d.currentInstruction {
			// Entry to be highlighted
			scrollTo = counter
		}

		name, length, b := getOpcodeInfo(gb, uint16(addr))
		d.entries[counter].name = name
		d.entries[counter].address = uint16(addr)
		d.entries[counter].bytes = b
		counter++

		addr += length
	}

	// Update number of entries
	d.totalEntries = counter
	d.slider.Max = counter - d.length

	// Scroll to correct offset
	scrollTo -= d.length / 2 // Selected instruction always at center
	d.scrollTo(scrollTo)
	d.slider.Current = scrollTo

	d.refresh()
}

func newDisassembler() *disassembler {
	dis := &disassembler{
		entries:      make([]*disassemblerEntry, 0x10000),
		totalEntries: 0x10000,
		length:       32,
		breakpoints:  make(map[uint16]struct{}),
	}
	dis.rowsWidget = make([]widget.PreferredSizeLocateableWidget, dis.length)

	// Initialize the disassembler with dummy data
	for i := range 0x10000 {
		dis.entries[i] = &disassemblerEntry{
			address: uint16(i),
			name:    fmt.Sprintf("%04X    NOP    ; No operation", i),
		}
	}

	entryList := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Spacing(1), // Add a small margin between entries
		)),
	)

	// Populate the container with buttons
	for i := 0; i < dis.length; i++ {
		entry := dis.createRow(i)
		entryList.AddChild(entry)
		dis.rowsWidget[i] = entry
	}

	// Slider
	dis.slider = widget.NewSlider(
		widget.SliderOpts.Images(&widget.SliderTrackImage{
			Idle: image.NewNineSliceColor(theme.Debugger.Slider.TrackColor),
		}, theme.Debugger.Button.Image),
		widget.SliderOpts.MinHandleSize(15), // Width of handle
		widget.SliderOpts.Direction(widget.DirectionVertical),
		widget.SliderOpts.MinMax(0, dis.totalEntries-dis.length),
		widget.SliderOpts.PageSizeFunc(func() int {
			return dis.length / 2
		}),
		widget.SliderOpts.ChangedHandler(func(args *widget.SliderChangedEventArgs) {
			dis.scrollTo(args.Slider.Current)
		}),
		widget.SliderOpts.WidgetOpts(
			// Stretch to container height
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{Stretch: true}),
			// Set slider height to non-zero value for correct layout computation
			widget.WidgetOpts.MinSize(0, 1),
		),
	)

	// Allow scrolling with mouse wheel
	scrollContainer := widget.NewScrollContainer(
		widget.ScrollContainerOpts.Content(entryList),
		// Image is required (set to transparent)
		widget.ScrollContainerOpts.Image(&widget.ScrollContainerImage{
			Idle: image.NewNineSliceColor(color.RGBA{}),
			Mask: image.NewNineSliceColor(color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}),
		}),
	)
	scrollContainer.GetWidget().ScrolledEvent.AddHandler(func(args any) {
		if a, ok := args.(*widget.WidgetScrolledEventArgs); ok {
			amount := computeRowsToScroll(a.Y)
			dis.scrollTo(dis.first + amount)
		}
	})

	dis.Container = newContainer(widget.DirectionHorizontal,
		scrollContainer, dis.slider,
	)
	return dis
}

func (d *disassembler) ToggleBreakpoint(addr uint16) {
	if d.IsBreakpoint(addr) {
		delete(d.breakpoints, addr)
	} else {
		d.breakpoints[addr] = struct{}{}
	}
}

func (d *disassembler) IsBreakpoint(addr uint16) bool {
	_, ok := d.breakpoints[addr]
	return ok
}

func (d *disassembler) createControlButton(name string, f func()) *widget.Button {
	return widget.NewButton(
		widget.ButtonOpts.Image(theme.Debugger.Button.Image),                // Background
		widget.ButtonOpts.Text(name, font, theme.Debugger.Button.TextColor), // Font and text

		// Click handler
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			f()
		}),
	)
}

func (d *disassembler) createRow(rowId int) widget.PreferredSizeLocateableWidget {
	button := widget.NewButton(
		widget.ButtonOpts.Image(entryImage),                               // Background
		widget.ButtonOpts.Text("", font, theme.Debugger.Button.TextColor), // Font and text
		widget.ButtonOpts.TextPosition(0, 0),

		// Click handler
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			entry := d.entries[d.first+rowId]
			d.ToggleBreakpoint(entry.address)
			d.refreshEntry(rowId)
		}),
	)

	// Fix widget min width so that even if buttons are smaller it doesn't resize
	button.GetWidget().MinWidth = 300

	return button
}

// refreshEntry changes button rendering based on entry state
func (d *disassembler) refreshEntry(entryId int) {
	entry := d.entries[d.first+entryId]
	button := d.rowsWidget[entryId].(*widget.Button)

	// Update label
	bytesStr := ""
	for _, b := range entry.bytes {
		bytesStr += fmt.Sprintf("%02X ", b)
	}

	// Add padding to align the instruction name
	for len(bytesStr) < 9 { // 3 chars per byte, up to 3 bytes
		bytesStr += "   "
	}
	button.Text().Label = fmt.Sprintf("%04X: %s  %s", entry.address, bytesStr, entry.name)

	// Update color
	isCurr := int(entry.address) == d.currentInstruction
	isBreakpoint := d.IsBreakpoint(entry.address)

	if isCurr {
		if isBreakpoint {
			button.Image = entryBreakpointAndCurrentImage
		} else {
			button.Image = entryCurrentImage
		}
	} else {
		if isBreakpoint {
			button.Image = entryBreakpointImage
		} else {
			button.Image = entryImage
		}
	}
}

// refresh disassembler rows
func (d *disassembler) refresh() {
	// Update all rows
	for i := range d.rowsWidget {
		d.refreshEntry(i)
	}
}

func (d *disassembler) scrollTo(newOffset int) {
	d.first = newOffset
	d.first = max(d.first, 0)                       // Reset to 0 if too low
	d.first = min(d.first, d.totalEntries-d.length) // Reset to maximum if too high

	d.slider.Current = d.first
	d.refresh()
}
