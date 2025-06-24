package ui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/danielecanzoneri/gb-emulator/pkg/debug"
	"image/color"
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

	parentEntry *disassemblyEntry

	address            uint16
	name               string
	bytes              []uint8
	isBreakpoint       bool
	currentInstruction bool

	// Prevent breakpoint selection when inactive
	active bool

	onTapped func(uint16, bool)
}

func newDisassemblyEntry(address uint16, name string) *disassemblyEntry {
	entry := &disassemblyEntry{
		address: address,
		name:    name,
		active:  true,
	}
	entry.ExtendBaseWidget(entry)
	return entry
}

func (e *disassemblyEntry) Tapped(_ *fyne.PointEvent) {
	if e.active {
		e.isBreakpoint = !e.isBreakpoint
		e.parentEntry.isBreakpoint = e.isBreakpoint
		if e.onTapped != nil {
			e.onTapped(e.address, e.isBreakpoint)
		}

		e.Refresh()
	}
}

func (e *disassemblyEntry) CreateRenderer() fyne.WidgetRenderer {
	th := e.Theme()
	v := fyne.CurrentApp().Settings().ThemeVariant()

	// Text placeholder
	text := canvas.NewText("0000    NOP    ; No operation", th.Color(theme.ColorNameForeground, v))
	text.TextStyle = fyne.TextStyle{Monospace: true}

	breakBackground := canvas.NewRectangle(color.Transparent)
	currentBackground := canvas.NewRectangle(color.Transparent)

	// Create the breakpoint circle
	circle := canvas.NewCircle(color.RGBA{R: 255, G: 0, B: 0, A: 255})

	return &disassemblyEntryRenderer{
		entry:             e,
		text:              text,
		breakBackground:   breakBackground,
		currentBackground: currentBackground,
		circle:            circle,
	}
}

type disassemblyEntryRenderer struct {
	entry             *disassemblyEntry
	text              *canvas.Text
	breakBackground   *canvas.Rectangle
	currentBackground *canvas.Rectangle
	circle            *canvas.Circle
}

func (r *disassemblyEntryRenderer) Destroy() {}

func (r *disassemblyEntryRenderer) Layout(size fyne.Size) {
	// Position and size the breakpoint icon
	circleSize := float32(breakpointIconSize)
	circleMargin := (size.Height - circleSize) / 2
	r.circle.Resize(fyne.NewSize(circleSize, circleSize))
	r.circle.Move(fyne.NewPos(circleMargin, circleMargin))

	// Position "current instruction" background
	r.currentBackground.Resize(fyne.NewSize(size.Width, size.Height))
	r.currentBackground.Move(fyne.NewPos(0, 0))

	// Position text and background after the circle with some padding
	circlePadding := circleSize + 2*circleMargin
	r.breakBackground.Resize(fyne.NewSize(size.Width-circlePadding, size.Height))
	r.breakBackground.Move(fyne.NewPos(circlePadding, 0))
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
	return []fyne.CanvasObject{r.breakBackground, r.currentBackground, r.circle, r.text}
}

func (r *disassemblyEntryRenderer) Refresh() {
	// Update circle visibility and background colors based on breakpoint status
	if r.entry.isBreakpoint {
		r.breakBackground.FillColor = color.RGBA{R: 255, G: 0, B: 0, A: 191}
		r.circle.FillColor = color.RGBA{R: 255, G: 0, B: 0, A: 255}
	} else {
		r.breakBackground.FillColor = color.Transparent
		r.circle.FillColor = color.Transparent
	}
	r.breakBackground.Refresh()
	r.circle.Refresh()

	if r.entry.currentInstruction {
		r.currentBackground.FillColor = color.RGBA{R: 0, G: 127, B: 255, A: 63}
	} else {
		r.currentBackground.FillColor = color.Transparent
	}
	r.currentBackground.Refresh()

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

	// Keep track to unselect when stepping
	previousEntry *disassemblyEntry
}

func newDisassembler(onEntryTapped func(uint16, bool)) *disassembler {
	dl := &disassembler{
		entries: make([]*disassemblyEntry, 0x10000),
	}
	dl.ExtendBaseWidget(dl)

	// Initialize the list with dummy data
	for i := range 0x10000 {
		dl.entries[i] = newDisassemblyEntry(uint16(i), fmt.Sprintf("%04X    NOP    ; No operation", i))
		dl.entries[i].onTapped = onEntryTapped
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
			currentEntry.parentEntry = entry

			currentEntry.address = entry.address
			currentEntry.name = entry.name
			currentEntry.currentInstruction = entry.currentInstruction
			currentEntry.isBreakpoint = entry.isBreakpoint
			currentEntry.bytes = entry.bytes
			currentEntry.active = entry.active
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
		if 0x104 <= addr && addr < 0x150 { // Header memory
			dl.entries[counter].name = "Cart Header"
			dl.entries[counter].address = uint16(addr)
			dl.entries[counter].bytes = []uint8{state.Memory[addr]}
			counter++
			addr++
			continue
		}

		if uint16(addr) == state.PC {
			// Clear previous instruction background
			if dl.previousEntry != nil {
				dl.previousEntry.currentInstruction = false
			}
			scrollTo = counter

			dl.entries[counter].currentInstruction = true
			dl.previousEntry = dl.entries[counter]
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

func (dl *disassembler) Enable() {
	for i := range dl.length {
		dl.entries[i].active = true
	}
}

func (dl *disassembler) Disable() {
	for i := range dl.length {
		dl.entries[i].active = false
	}
}
