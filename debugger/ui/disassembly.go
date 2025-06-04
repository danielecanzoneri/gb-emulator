package ui

import (
	"fmt"
	"github.com/danielecanzoneri/gb-emulator/pkg/debug"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

const (
	disassemblyEntryVerticalPadding = 2
	disassemblyEntryRightPadding    = 16
	disassemblyEntryTextLeftPadding = 2

	breakpointIconSize = 10
)

// disassemblyEntry represents a single line in the disassembler
type disassemblyEntry struct {
	widget.BaseWidget

	address      uint16
	name         string
	bytes        []uint8
	isBreakpoint bool

	onTapped func(uint16)
}

func newDisassemblyEntry(address uint16, name string) *disassemblyEntry {
	entry := &disassemblyEntry{
		address: address,
		name:    name,
	}
	entry.ExtendBaseWidget(entry)
	return entry
}

func (e *disassemblyEntry) Tapped(_ *fyne.PointEvent) {
	e.isBreakpoint = !e.isBreakpoint
	if e.onTapped != nil {
		e.onTapped(e.address)
	}

	e.Refresh()
}

func (e *disassemblyEntry) CreateRenderer() fyne.WidgetRenderer {
	th := e.Theme()
	v := fyne.CurrentApp().Settings().ThemeVariant()

	// Text placeholder
	text := canvas.NewText("0000    NOP    ; No operation", th.Color(theme.ColorNameForeground, v))
	text.TextStyle = fyne.TextStyle{Monospace: true}
	background := canvas.NewRectangle(color.Transparent)

	// Create the breakpoint circle
	circle := canvas.NewCircle(color.RGBA{R: 255, G: 0, B: 0, A: 255})

	return &disassemblyEntryRenderer{
		entry:      e,
		text:       text,
		background: background,
		circle:     circle,
	}
}

type disassemblyEntryRenderer struct {
	entry      *disassemblyEntry
	text       *canvas.Text
	background *canvas.Rectangle
	circle     *canvas.Circle
}

func (r *disassemblyEntryRenderer) Destroy() {}

func (r *disassemblyEntryRenderer) Layout(size fyne.Size) {
	// Position and size the breakpoint icon
	circleSize := float32(breakpointIconSize)
	circleMargin := (size.Height - circleSize) / 2
	r.circle.Resize(fyne.NewSize(circleSize, circleSize))
	r.circle.Move(fyne.NewPos(circleMargin, circleMargin))

	// Position text and background after the circle with some padding
	circlePadding := circleSize + 2*circleMargin
	r.background.Resize(fyne.NewSize(size.Width-circlePadding, size.Height))
	r.background.Move(fyne.NewPos(circlePadding, 0))
	r.text.Resize(fyne.NewSize(size.Width-circlePadding-disassemblyEntryTextLeftPadding, size.Height))
	r.text.Move(fyne.NewPos(circlePadding+disassemblyEntryTextLeftPadding, 0))
}

func (r *disassemblyEntryRenderer) MinSize() fyne.Size {
	width, height := r.text.MinSize().Width, r.text.MinSize().Height

	// Add padding
	height += 2 * disassemblyEntryVerticalPadding
	width += disassemblyEntryRightPadding + disassemblyEntryTextLeftPadding

	// Add left padding for the disassembly icon
	circleSize := float32(breakpointIconSize)
	circleMargin := (height - circleSize) / 2
	width += circleSize + 2*circleMargin

	return fyne.NewSize(width, height)
}

func (r *disassemblyEntryRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.background, r.circle, r.text}
}

func (r *disassemblyEntryRenderer) Refresh() {
	// Update circle visibility and background colors based on breakpoint status
	if r.entry.isBreakpoint {
		r.background.FillColor = color.RGBA{R: 255, G: 0, B: 0, A: 64}
		r.circle.FillColor = color.RGBA{R: 255, G: 0, B: 0, A: 255}
	} else {
		r.background.FillColor = color.Transparent
		r.circle.FillColor = color.Transparent
	}
	r.background.Refresh()
	r.circle.Refresh()

	// Format the bytes part
	bytesStr := ""
	for _, b := range r.entry.bytes {
		bytesStr += fmt.Sprintf("%02X ", b)
	}

	// Add padding to align the instruction name
	for len(bytesStr) < 9 { // 3 chars per byte, up to 3 bytes
		bytesStr += "   "
	}
	r.text.Text = fmt.Sprintf("%04X: %s  %s", r.entry.address, bytesStr, r.entry.name)
	r.text.Refresh()
}

type disassembler struct {
	widget.BaseWidget

	list *widget.List

	length  int
	entries []*disassemblyEntry
}

func newDisassembler() *disassembler {
	dl := &disassembler{
		entries: make([]*disassemblyEntry, 0x10000),
	}
	dl.ExtendBaseWidget(dl)

	// Initialize the list with dummy data
	for i := range 0x10000 {
		dl.entries[i] = newDisassemblyEntry(uint16(i), fmt.Sprintf("%04X    NOP    ; No operation", i))
		dl.entries[i].onTapped = func(addr uint16) {
			dl.entries[i].isBreakpoint = !dl.entries[i].isBreakpoint
			fmt.Println(addr)
		}
	}

	dl.list = &widget.List{
		Length: func() int {
			return dl.length
		},
		CreateItem: func() fyne.CanvasObject {
			return newDisassemblyEntry(0, "")
		},
		UpdateItem: func(id widget.ListItemID, item fyne.CanvasObject) {
			entry := dl.entries[id]
			currentEntry := item.(*disassemblyEntry)

			currentEntry.address = entry.address
			currentEntry.name = entry.name
			currentEntry.isBreakpoint = entry.isBreakpoint
			currentEntry.bytes = entry.bytes
			currentEntry.onTapped = entry.onTapped
			currentEntry.Refresh()
		},
	}

	return dl
}

func (dl *disassembler) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(dl.list)
}

// Update scans through memory and marks which addresses contain
// executable code vs data bytes that are part of multibyte instructions.
// This is used to properly display the disassembly view.
func (dl *disassembler) Update(state *debug.GameBoyState) {
	counter := 0
	var scrollTo int

	for addr := 0; addr < 0x10000; {
		if uint16(addr) == state.PC {
			scrollTo = counter
		}

		name, length, bytes := getOpcodeInfo(state, uint16(addr))
		dl.entries[counter].name = name
		dl.entries[counter].address = uint16(addr)
		dl.entries[counter].bytes = bytes
		counter++

		addr += length
	}

	dl.length = counter

	// Scroll to current PC
	fyne.Do(func() {
		dl.list.ScrollTo(scrollTo)
		dl.Refresh()
	})
}
