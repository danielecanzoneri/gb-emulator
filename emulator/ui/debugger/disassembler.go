package debugger

import (
	"fmt"
	"github.com/danielecanzoneri/gb-emulator/emulator/internal/gameboy"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"image/color"
	"log"
	"strings"
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
	buttonTextColor = &widget.ButtonTextColor{
		Idle: color.Black,
	}
	buttonTextPadding = widget.Insets{
		Left: 2, Right: 2, Top: 2, Bottom: 2,
	}
)

// disassemblyEntry represents a single line in the disassembler
type disassemblerEntry struct {
	address      uint16
	name         string
	bytes        []uint8
	isBreakpoint bool
}

type disassembler struct {
	*widget.Container

	entries      []*disassemblerEntry
	totalEntries int

	// Current entry highlighted (to unselect when stepping)
	selected *disassemblerEntry

	// Entries to show
	first  int
	length int
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

	d.totalEntries = counter
	d.scrollTo(scrollTo - d.length/2) // Selected instruction always at center

	d.refresh()
}

func newDisassembler() *disassembler {
	d := &disassembler{
		entries: make([]*disassemblerEntry, 0x10000),
		length:  24,
	}

	// Initialize the disassembler with dummy data
	for i := range 0x10000 {
		d.entries[i] = &disassemblerEntry{
			address: uint16(i),
			name:    fmt.Sprintf("%04X    NOP    ; No operation", i),
		}
	}
	d.selected = d.entries[0]

	d.Container = widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Spacing(1), // Add a small margin between entries
		)),
	)

	// Populate the container with buttons
	for i := 0; i < d.length; i++ {
		entry := d.createRow(i)
		d.AddChild(entry)
	}

	// Add a row containing buttons for fast scrolling
	navigateButtons := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Spacing(5),
		)),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
			}),
		),
	)
	for i := 1; i <= 3; i++ {
		txt := strings.Repeat("↓", i)
		offset := 1 << (5 * (i - 1)) // 1: 1, 2: 32, 3: 1024
		button := widget.NewButton(
			widget.ButtonOpts.Image(buttonImage),               // Background
			widget.ButtonOpts.Text(txt, font, buttonTextColor), // Font and text
			widget.ButtonOpts.TextPadding(buttonTextPadding),

			// Click handler
			widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
				d.scrollTo(d.first + offset)
			}),
		)
		navigateButtons.AddChild(button)
	}
	navigateButtons.AddChild(widget.NewContainer(widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.MinSize(10, 0))))
	for i := 1; i <= 3; i++ {
		txt := strings.Repeat("↑", i)
		offset := 1 << (5 * (i - 1)) // 1: 1, 2: 32, 3: 1024
		button := widget.NewButton(
			widget.ButtonOpts.Image(buttonImage),               // Background
			widget.ButtonOpts.Text(txt, font, buttonTextColor), // Font and text
			widget.ButtonOpts.TextPadding(buttonTextPadding),

			// Click handler
			widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
				d.scrollTo(d.first - offset)
			}),
		)
		navigateButtons.AddChild(button)
	}

	containerWidth, _ := d.PreferredSize()
	centerButtons := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.MinSize(containerWidth, 0),
		),
	)
	centerButtons.AddChild(navigateButtons)
	d.AddChild(centerButtons)

	return d
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
			entry.isBreakpoint = !entry.isBreakpoint
			refreshEntry(args.Button, entry)

			log.Printf("Pressed %04X\n", entry.address)
		}),
	)

	// Fix widget min width so that even if buttons are smaller it doesn't resize
	button.GetWidget().MinWidth = 300

	return button
}

// refreshEntry changes button rendering based on entry state
func refreshEntry(button *widget.Button, entry *disassemblerEntry) {
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
	if entry.isBreakpoint {
		button.Image = buttonImageBreakpoint
	} else {
		button.Image = buttonImage
	}
}

// refresh disassembler rows
func (d *disassembler) refresh() {
	// Update all rows
	rows := d.Children()[:d.length]
	for i, r := range rows {
		button := r.(*widget.Button)
		refreshEntry(button, d.entries[d.first+i])
	}
}

func (d *disassembler) scrollTo(newOffset int) {
	d.first = newOffset
	d.first = max(d.first, 0)                       // Reset to 0 if too low
	d.first = min(d.first, d.totalEntries-d.length) // Reset to maximum if too high

	d.refresh()
}
