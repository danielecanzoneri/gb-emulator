package debugger

import (
	"bytes"
	"fmt"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font/gofont/gomono"
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

var font text.Face

// disassemblyEntry represents a single line in the disassembler
type disassemblerEntry struct {
	address      uint16
	name         string
	bytes        []uint8
	isBreakpoint bool
}

type disassembler struct {
	widget  *widget.Container
	entries []*disassemblerEntry

	// Current entry highlighted (to unselect when stepping)
	selected *disassemblerEntry

	// Entries to show
	first  int
	length int
}

func newDisassembler() *disassembler {
	d := &disassembler{
		entries: make([]*disassemblerEntry, 0x10000),
		length:  16,
	}

	// Initialize the disassembler with dummy data
	for i := range 0x10000 {
		d.entries[i] = &disassemblerEntry{
			address: uint16(i),
			name:    fmt.Sprintf("%04X    NOP    ; No operation", i),
		}
	}
	d.selected = d.entries[0]

	font = loadFont(16)

	container := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Spacing(1), // Add a small margin between entries
		)),
	)
	d.widget = container

	// Populate the container with buttons
	for i := 0; i < d.length; i++ {
		entry := createEntryWidget(d.entries[i])
		container.AddChild(entry)
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
				d.scroll(offset)
			}),
		)
		navigateButtons.AddChild(button)
	}
	navigateButtons.AddChild(widget.NewContainer(widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.MinSize(10, 0))))
	for i := 1; i <= 3; i++ {
		txt := strings.Repeat("↑", i)
		offset := 1 << (4 * (i - 1)) // 1: 1, 2: 32, 3: 1024
		button := widget.NewButton(
			widget.ButtonOpts.Image(buttonImage),               // Background
			widget.ButtonOpts.Text(txt, font, buttonTextColor), // Font and text
			widget.ButtonOpts.TextPadding(buttonTextPadding),

			// Click handler
			widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
				d.scroll(-offset)
			}),
		)
		navigateButtons.AddChild(button)
	}

	containerWidth, _ := container.PreferredSize()
	centerButtons := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.MinSize(containerWidth, 0),
		),
	)
	centerButtons.AddChild(navigateButtons)
	container.AddChild(centerButtons)

	return d
}

func createEntryWidget(entry *disassemblerEntry) widget.PreferredSizeLocateableWidget {
	return widget.NewButton(
		widget.ButtonOpts.Image(buttonImage),                      // Background
		widget.ButtonOpts.Text(entry.name, font, buttonTextColor), // Font and text
		widget.ButtonOpts.TextPadding(buttonTextPadding),

		// Click handler
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			entry.isBreakpoint = !entry.isBreakpoint
			updateEntry(args.Button, entry)

			log.Printf("Pressed %04X\n", entry.address)
		}),
	)
}

// updateEntry changes button rendering based on entry state
func updateEntry(button *widget.Button, entry *disassemblerEntry) {
	button.Text().Label = entry.name // Update label

	// Update color
	if entry.isBreakpoint {
		button.Image = buttonImageBreakpoint
	} else {
		button.Image = buttonImage
	}
}

func (d *disassembler) update() {
	// Update all rows
	rows := d.widget.Children()[:d.length]
	for i, r := range rows {
		button := r.(*widget.Button)
		updateEntry(button, d.entries[d.first+i])
	}
}

func (d *disassembler) scroll(offset int) {
	d.first += offset
	d.first = max(d.first, 0)                       // Reset to 0 if too low
	d.first = min(d.first, len(d.entries)-d.length) // Reset to maximum if too high

	d.update()
}

func loadFont(size float64) text.Face {
	s, err := text.NewGoTextFaceSource(bytes.NewReader(gomono.TTF))
	if err != nil {
		log.Fatal(err)
	}

	return &text.GoTextFace{
		Source: s,
		Size:   size,
	}
}
