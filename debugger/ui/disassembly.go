package ui

import (
	"fmt"
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
	text         string
	isBreakpoint bool

	onTapped func(uint16)
}

func newDisassemblyEntry(address uint16, text string) *disassemblyEntry {
	entry := &disassemblyEntry{
		address: address,
		text:    text,
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
	// Update circle visibility and background colorsbased on breakpoint status
	if r.entry.isBreakpoint {
		r.background.FillColor = color.RGBA{R: 255, G: 0, B: 0, A: 64}
		r.circle.FillColor = color.RGBA{R: 255, G: 0, B: 0, A: 255}
	} else {
		r.background.FillColor = color.Transparent
		r.circle.FillColor = color.Transparent
	}
	r.background.Refresh()
	r.circle.Refresh()

	r.text.Text = r.entry.text
	r.text.Refresh()
}

type Disassembler struct {
	widget.List
	entries []*disassemblyEntry
}

func NewDisassembler() *Disassembler {
	dl := &Disassembler{
		entries: make([]*disassemblyEntry, 0x10000),
	}

	// Initialize the list with dummy data
	for i := range 0x10000 {
		dl.entries[i] = newDisassemblyEntry(uint16(i), fmt.Sprintf("%04X    NOP    ; No operation", i))
		dl.entries[i].onTapped = func(addr uint16) {
			dl.entries[i].isBreakpoint = !dl.entries[i].isBreakpoint
			fmt.Println(addr)
		}
	}

	dl.List = widget.List{
		Length: func() int {
			return len(dl.entries)
		},
		CreateItem: func() fyne.CanvasObject {
			return newDisassemblyEntry(0, "")
		},
		UpdateItem: func(id widget.ListItemID, item fyne.CanvasObject) {
			entry := dl.entries[id]
			currentEntry := item.(*disassemblyEntry)

			currentEntry.address = entry.address
			currentEntry.text = entry.text
			currentEntry.isBreakpoint = entry.isBreakpoint
			currentEntry.onTapped = entry.onTapped
			currentEntry.Refresh()
		},
	}

	return dl
}
