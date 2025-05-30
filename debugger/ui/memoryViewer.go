package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

const (
	memoryRowRightPadding = 16
)

// disassemblyEntry represents a single line in the disassembler
type memoryRow struct {
	widget.BaseWidget

	baseAddress uint16
	data        [16]uint8
}

func newMemoryRow(baseAddress uint16) *memoryRow {
	row := &memoryRow{
		baseAddress: baseAddress,
		data:        [16]uint8{},
	}
	row.ExtendBaseWidget(row)
	return row
}

func (r *memoryRow) CreateRenderer() fyne.WidgetRenderer {
	th := r.Theme()
	v := fyne.CurrentApp().Settings().ThemeVariant()

	// Text placeholder
	text := canvas.NewText("0000  00 00 00 00 00 00 00 00|00 00 00 00 00 00 00 00", th.Color(theme.ColorNameForeground, v))
	text.TextStyle = fyne.TextStyle{Monospace: true}
	background := canvas.NewRectangle(th.Color(theme.ColorNameBackground, v))

	return &memoryRowRenderer{
		entry:      r,
		text:       text,
		background: background,
	}
}

type memoryRowRenderer struct {
	entry      *memoryRow
	text       *canvas.Text
	background *canvas.Rectangle
}

func (r *memoryRowRenderer) Destroy() {}

func (r *memoryRowRenderer) Layout(size fyne.Size) {
	r.background.Resize(size)
	r.background.Move(fyne.NewPos(0, 0))
	r.text.Resize(size)
	r.text.Move(fyne.NewPos(0, 0))
}

func (r *memoryRowRenderer) MinSize() fyne.Size {
	// Right padding for the scrollbar
	return fyne.NewSize(r.text.MinSize().Width+memoryRowRightPadding, r.text.MinSize().Height)
}

func (r *memoryRowRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.background, r.text}
}

func (r *memoryRowRenderer) Refresh() {
	r.text.Text = fmt.Sprintf("%04X  %02X %02X %02X %02X %02X %02X %02X %02X|%02X %02X %02X %02X %02X %02X %02X %02X",
		r.entry.baseAddress,
		r.entry.data[0], r.entry.data[1], r.entry.data[2], r.entry.data[3],
		r.entry.data[4], r.entry.data[5], r.entry.data[6], r.entry.data[7],
		r.entry.data[8], r.entry.data[9], r.entry.data[10], r.entry.data[11],
		r.entry.data[12], r.entry.data[13], r.entry.data[14], r.entry.data[15])
	r.text.Refresh()
}

type memoryViewer struct {
	widget.List
	rows    int
	entries []*memoryRow
}

func newMemoryViewer() *memoryViewer {
	mv := &memoryViewer{
		rows:    0x10000 / 16,
		entries: make([]*memoryRow, 0x10000/16),
	}

	// Initialize the list with dummy data
	for i := range mv.rows {
		mv.entries[i] = newMemoryRow(uint16(i * 16))
	}
	mv.entries[1].data[0] = 0x12

	mv.List = widget.List{
		Length: func() int {
			return len(mv.entries)
		},
		CreateItem: func() fyne.CanvasObject {
			return newMemoryRow(0)
		},
		UpdateItem: func(id widget.ListItemID, item fyne.CanvasObject) {
			entry := mv.entries[id]
			currentEntry := item.(*memoryRow)

			currentEntry.baseAddress = entry.baseAddress
			currentEntry.data = entry.data
			currentEntry.Refresh()
		},
	}

	return mv
}

// MinSize returns the minimum size for the memory viewer widget.
// It calculates this by taking the minimum size of a single row and
// multiplying the height by 16 to show a reasonable number of rows.
func (mv *memoryViewer) MinSize() fyne.Size {
	baseSize := mv.entries[0].MinSize()
	height := baseSize.Height * 16

	return fyne.NewSize(baseSize.Width, height)
}
