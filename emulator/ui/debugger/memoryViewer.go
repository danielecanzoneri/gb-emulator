package debugger

import (
	"fmt"
	"github.com/danielecanzoneri/gb-emulator/emulator/internal/gameboy"
	"github.com/ebitenui/ebitenui/widget"
	"strings"
)

// disassemblyEntry represents a single line in the disassembler
type memoryRow struct {
	baseAddress uint16
	data        [16]uint8
}

type memoryViewer struct {
	*widget.Container

	entries []*memoryRow

	// Entries to show
	first  int
	length int
}

func newMemoryViewer() *memoryViewer {
	mv := &memoryViewer{
		entries: make([]*memoryRow, 0x10000/16),
		length:  16,
	}

	// Initialize the rows
	for i := range mv.entries {
		mv.entries[i] = &memoryRow{
			baseAddress: uint16(i * 16),
		}
	}

	mv.Container = widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Spacing(1), // Add a small margin between entries
		)),
	)

	// Populate the container with the rows
	for i := 0; i < mv.length; i++ {
		entry := mv.createRow()
		mv.AddChild(entry)
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
		offset := 1 << (4 * (i - 1)) // 1: 1, 2: 16, 3: 256
		button := widget.NewButton(
			widget.ButtonOpts.Image(buttonImage),               // Background
			widget.ButtonOpts.Text(txt, font, buttonTextColor), // Font and text
			widget.ButtonOpts.TextPadding(buttonTextPadding),

			// Click handler
			widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
				mv.scrollTo(mv.first + offset)
			}),
		)
		navigateButtons.AddChild(button)
	}
	navigateButtons.AddChild(widget.NewContainer(widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.MinSize(10, 0))))
	for i := 1; i <= 3; i++ {
		txt := strings.Repeat("↑", i)
		offset := 1 << (4 * (i - 1)) // 1: 1, 2: 16, 3: 256
		button := widget.NewButton(
			widget.ButtonOpts.Image(buttonImage),               // Background
			widget.ButtonOpts.Text(txt, font, buttonTextColor), // Font and text
			widget.ButtonOpts.TextPadding(buttonTextPadding),

			// Click handler
			widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
				mv.scrollTo(mv.first - offset)
			}),
		)
		navigateButtons.AddChild(button)
	}

	containerWidth, _ := mv.PreferredSize()
	centerButtons := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.MinSize(containerWidth, 0),
		),
	)
	centerButtons.AddChild(navigateButtons)
	mv.AddChild(centerButtons)

	return mv
}

// Sync data from memory
func (mv *memoryViewer) Sync(gb *gameboy.GameBoy) {
	for _, entry := range mv.entries {
		for i := range 16 {
			entry.data[i] = gb.Memory.DebugRead(entry.baseAddress + uint16(i))
		}
	}

	mv.refresh()
}

func (mv *memoryViewer) createRow() widget.PreferredSizeLocateableWidget {
	dummyText := "0000  00 00 00 00 00 00 00 00 | 00 00 00 00 00 00 00 00"
	label := widget.NewText(
		widget.TextOpts.Text(dummyText, font, buttonTextColor.Idle), // Font and text
		widget.TextOpts.Insets(buttonTextPadding),
	)
	c := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
		widget.ContainerOpts.BackgroundImage(buttonImage.Idle), // Same as idle button
	)
	c.AddChild(label)
	return c
}

func (mv *memoryViewer) refresh() {
	// Update all rows
	rows := mv.Children()[:mv.length]
	for i, r := range rows {
		label := r.(*widget.Container).Children()[0].(*widget.Text)
		entry := mv.entries[mv.first+i]
		label.Label = fmt.Sprintf("%04X  %02X %02X %02X %02X %02X %02X %02X %02X | %02X %02X %02X %02X %02X %02X %02X %02X",
			entry.baseAddress,
			entry.data[0], entry.data[1], entry.data[2], entry.data[3],
			entry.data[4], entry.data[5], entry.data[6], entry.data[7],
			entry.data[8], entry.data[9], entry.data[10], entry.data[11],
			entry.data[12], entry.data[13], entry.data[14], entry.data[15],
		)
	}
}

func (mv *memoryViewer) scrollTo(newOffset int) {
	mv.first = newOffset
	mv.first = max(mv.first, 0)                         // Reset to 0 if too low
	mv.first = min(mv.first, len(mv.entries)-mv.length) // Reset to maximum if too high

	mv.refresh()
}
