package debugger

import (
	"fmt"
	"github.com/danielecanzoneri/lucky-boy/ui/graphics"
	"image/color"

	"github.com/danielecanzoneri/lucky-boy/gameboy"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"golang.org/x/image/colornames"
)

// disassemblyEntry represents a single line in the disassembler
type memoryRow struct {
	baseAddress uint16
	data        [16]uint8
}

type memoryViewer struct {
	*widget.Container

	slider *widget.Slider

	entries    []*memoryRow
	rowsWidget []widget.PreferredSizeLocateableWidget

	// Entries to show
	first  int
	length int
}

func newMemoryViewer() *memoryViewer {
	mv := &memoryViewer{
		entries: make([]*memoryRow, 0x10000/16),
		length:  16,
	}
	mv.rowsWidget = make([]widget.PreferredSizeLocateableWidget, mv.length)

	// Initialize the rows
	for i := range mv.entries {
		mv.entries[i] = &memoryRow{
			baseAddress: uint16(i * 16),
		}
	}

	entryList := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Spacing(1), // Add a small margin between entries
		)),
	)
	// Populate the container with the rows
	for i := 0; i < mv.length; i++ {
		entry := mv.createRow()
		entryList.AddChild(entry)
		mv.rowsWidget[i] = entry
	}

	// Slider
	mv.slider = widget.NewSlider(
		widget.SliderOpts.Images(&widget.SliderTrackImage{
			Idle: image.NewNineSliceColor(theme.Debugger.Slider.TrackColor),
		}, theme.Debugger.Button.Image),
		widget.SliderOpts.MinHandleSize(15), // Width of handle
		widget.SliderOpts.Direction(widget.DirectionVertical),
		widget.SliderOpts.MinMax(0, len(mv.entries)-mv.length),
		widget.SliderOpts.PageSizeFunc(func() int { return mv.length / 2 }),
		widget.SliderOpts.ChangedHandler(func(args *widget.SliderChangedEventArgs) {
			mv.scrollTo(args.Slider.Current)
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
		widget.ScrollContainerOpts.Padding(widget.NewInsetsSimple(theme.Debugger.Padding)),
	)
	scrollContainer.GetWidget().ScrolledEvent.AddHandler(func(args any) {
		if a, ok := args.(*widget.WidgetScrolledEventArgs); ok {
			amount := computeRowsToScroll(a.Y)
			mv.scrollTo(mv.first + amount)
		}
	})

	mv.Container = widget.NewContainer(
		widget.ContainerOpts.Layout(
			widget.NewRowLayout(widget.RowLayoutOpts.Direction(widget.DirectionHorizontal)),
		),
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(theme.Debugger.Main.Color)),
	)
	mv.Container.AddChild(scrollContainer, mv.slider)

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
	dummyText := "0000  00 00 00 00 00 00 00 00 | 00 00 00 00 00 00 00 00 | ................"
	label := widget.NewText(
		widget.TextOpts.Text(dummyText, &font, colornames.White), // Font and text
	)
	return label
}

func (mv *memoryViewer) refresh() {
	// Update all rows
	for i, r := range mv.rowsWidget {
		label := r.(*widget.Text)
		entry := mv.entries[mv.first+i]
		ascii := make([]byte, 16)
		for i, b := range entry.data {
			if b >= 32 && b <= 126 {
				ascii[i] = b
			} else {
				ascii[i] = '.'
			}
		}
		label.Label = fmt.Sprintf("%04X  %02X %02X %02X %02X %02X %02X %02X %02X | %02X %02X %02X %02X %02X %02X %02X %02X | %s",
			entry.baseAddress,
			entry.data[0], entry.data[1], entry.data[2], entry.data[3],
			entry.data[4], entry.data[5], entry.data[6], entry.data[7],
			entry.data[8], entry.data[9], entry.data[10], entry.data[11],
			entry.data[12], entry.data[13], entry.data[14], entry.data[15],
			string(ascii),
		)
	}
}

func (mv *memoryViewer) scrollTo(newOffset int) {
	mv.first = newOffset
	mv.first = max(mv.first, 0)                         // Reset to 0 if too low
	mv.first = min(mv.first, len(mv.entries)-mv.length) // Reset to maximum if too high

	mv.slider.Current = mv.first
	mv.refresh()
}
