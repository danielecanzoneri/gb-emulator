package debugger

import (
	"fmt"
	"github.com/danielecanzoneri/gb-emulator/emulator/internal/gameboy"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"image/color"
)

var (
	buttonImage = &widget.ButtonImage{
		Idle:    image.NewNineSliceColor(color.NRGBA{R: 170, G: 170, B: 180, A: 255}),
		Hover:   image.NewNineSliceColor(color.NRGBA{R: 130, G: 130, B: 150, A: 255}),
		Pressed: image.NewNineSliceColor(color.NRGBA{R: 90, G: 90, B: 120, A: 255}),
	}
	buttonImageBreakpoint = &widget.ButtonImage{
		Idle:    image.NewNineSliceColor(color.NRGBA{R: 255, G: 170, B: 180, A: 255}),
		Hover:   image.NewNineSliceColor(color.NRGBA{R: 255, G: 130, B: 150, A: 255}),
		Pressed: image.NewNineSliceColor(color.NRGBA{R: 255, G: 90, B: 120, A: 255}),
	}
	buttonImageCurrent = &widget.ButtonImage{
		Idle:    image.NewNineSliceColor(color.NRGBA{R: 170, G: 170, B: 255, A: 255}),
		Hover:   image.NewNineSliceColor(color.NRGBA{R: 130, G: 130, B: 255, A: 255}),
		Pressed: image.NewNineSliceColor(color.NRGBA{R: 90, G: 90, B: 255, A: 255}),
	}
	buttonTextColor = &widget.ButtonTextColor{
		Idle: color.Black,
	}
	buttonTextPadding = widget.Insets{
		Left: 2, Right: 2, Top: 2, Bottom: 2,
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

	// Current entry highlighted (to unselect when stepping)
	selected *disassemblerEntry

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

		if uint16(addr) == gb.CPU.PC {
			// Entry to be highlighted
			d.selected = d.entries[counter]
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

func (d *Debugger) newDisassembler() *disassembler {
	dis := &disassembler{
		entries:      make([]*disassemblerEntry, 0x10000),
		totalEntries: 0x10000,
		length:       24,
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
	dis.selected = dis.entries[0]

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
			Idle: image.NewNineSliceColor(color.NRGBA{255, 255, 255, 32}),
		}, buttonImage),
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
			p := -int(a.Y)
			dis.scrollTo(dis.first + p)
		}
	})

	// Step, continue buttons
	controlButtons := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Spacing(10)),
		),
	)
	controlButtons.AddChild(
		dis.createControlButton("Step", d.Step),
		dis.createControlButton("Continue", d.Continue),
		dis.createControlButton("Stop", d.Stop),
		dis.createControlButton("Reset", d.Reset),
	)

	dis.Container = newContainer(widget.DirectionVertical,
		newContainer(widget.DirectionHorizontal,
			scrollContainer, dis.slider,
		),
		controlButtons,
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
		widget.ButtonOpts.Image(buttonImage),                // Background
		widget.ButtonOpts.Text(name, font, buttonTextColor), // Font and text
		// widget.ButtonOpts.TextPadding(buttonTextPadding),
		// widget.ButtonOpts.TextPosition(0, 0),

		// Click handler
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			f()
		}),
	)
}

func (d *disassembler) createRow(rowId int) widget.PreferredSizeLocateableWidget {
	button := widget.NewButton(
		widget.ButtonOpts.Image(buttonImage),              // Background
		widget.ButtonOpts.Text("", font, buttonTextColor), // Font and text
		widget.ButtonOpts.TextPadding(buttonTextPadding),
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
	if d.selected == entry {
		button.Image = buttonImageCurrent
	} else {
		if d.IsBreakpoint(entry.address) {
			button.Image = buttonImageBreakpoint
		} else {
			button.Image = buttonImage
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
