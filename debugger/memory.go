package debugger

import (
	"fmt"
	"image"
	"image/color"
	"strings"

	ebitenimage "github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/input"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// MemoryViewer represents the memory visualization widget
type MemoryViewer struct {
	mem         MemoryDebugger
	container   *widget.Container
	scrollCont  *widget.ScrollContainer
	labels      []*widget.Label
	startAddr   uint16
	visibleRows int
	face        *text.GoTextFace
}

// initMemoryViewer initializes the memory viewer widget
func (d *Debugger) initMemoryViewer() *widget.Container {
	// Create container for the memory viewer
	container := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Spacing(2),
		)),
	)

	// Add title
	titleColor := &widget.LabelColor{
		Idle: color.RGBA{R: 255, G: 255, B: 200, A: 255},
	}
	titleLabel := widget.NewLabel(
		widget.LabelOpts.Text("Memory", d.face, titleColor),
	)
	container.AddChild(titleLabel)

	// Create memory viewer
	d.memViewer = NewMemoryViewer(d.mem, d.face)
	container.AddChild(d.memViewer.GetWidget())

	return container
}

// updateMemoryViewer updates the memory viewer content
func (d *Debugger) updateMemoryViewer() {
	if d.memViewer != nil {
		d.memViewer.Update()
	}
}

// NewMemoryViewer creates a new memory viewer widget
func NewMemoryViewer(mem MemoryDebugger, face *text.GoTextFace) *MemoryViewer {
	mv := &MemoryViewer{
		mem:         mem,
		face:        face,
		visibleRows: 16,
		labels:      make([]*widget.Label, 16),
	}

	// Create the container for memory lines
	mv.container = widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Spacing(2),
		)),
	)

	// Create the scroll container
	mv.scrollCont = widget.NewScrollContainer(
		widget.ScrollContainerOpts.Content(mv.container),
		widget.ScrollContainerOpts.StretchContentWidth(),
		widget.ScrollContainerOpts.Image(&widget.ScrollContainerImage{
			Idle: ebitenimage.NewNineSliceColor(color.RGBA{60, 60, 60, 255}),
			Mask: ebitenimage.NewNineSliceColor(color.RGBA{60, 60, 60, 255}),
		}),
	)

	// Create labels for each row
	textColor := &widget.LabelColor{
		Idle: color.White,
	}

	for i := 0; i < mv.visibleRows; i++ {
		label := widget.NewLabel(
			widget.LabelOpts.Text("", face, textColor),
		)
		mv.labels[i] = label
		mv.container.AddChild(label)
	}

	// Initial update
	mv.Update()

	return mv
}

// GetWidget returns the root widget of the memory viewer
func (mv *MemoryViewer) GetWidget() widget.PreferredSizeLocateableWidget {
	return mv.scrollCont
}

// Update updates the memory viewer's content and handles input
func (mv *MemoryViewer) Update() {
	// Handle mouse wheel scrolling
	x, y := input.CursorPosition()
	rect := mv.scrollCont.GetWidget().Rect
	if image.Pt(x, y).In(rect) {
		_, scrollY := input.Wheel()
		if scrollY > 0 {
			mv.ScrollUp()
		} else if scrollY < 0 {
			mv.ScrollDown()
		}
		fmt.Println("scrollY", scrollY)
	}

	// Update each row's content
	for i := 0; i < mv.visibleRows; i++ {
		addr := mv.startAddr + uint16(i*16)
		var hexBytes strings.Builder
		var asciiChars strings.Builder

		// Build hex representation
		for j := 0; j < 16; j++ {
			value := mv.mem.DebugRead(addr + uint16(j))
			hexBytes.WriteString(fmt.Sprintf("%02X ", value))

			// Add ASCII representation
			if value >= 32 && value <= 126 {
				asciiChars.WriteRune(rune(value))
			} else {
				asciiChars.WriteRune('.')
			}
		}

		// Update label text
		mv.labels[i].Label = fmt.Sprintf("%04X: %-48s %s", addr, hexBytes.String(), asciiChars.String())
	}
}

// ScrollUp scrolls the memory view up by one row
func (mv *MemoryViewer) ScrollUp() {
	if mv.startAddr >= 16 {
		mv.startAddr -= 16
		mv.Update()
	}
}

// ScrollDown scrolls the memory view down by one row
func (mv *MemoryViewer) ScrollDown() {
	if mv.startAddr <= 0xFFFF-uint16(mv.visibleRows*16) {
		mv.startAddr += 16
		mv.Update()
	}
}
