package debugger

import (
	"fmt"
	"image"
	"image/color"
	"strings"

	ebitenimage "github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// MemoryViewer represents the memory visualization widget
type MemoryViewer struct {
	mem         MemoryDebugger
	container   *widget.Container
	scrollCont  *widget.ScrollContainer
	labels      []*widget.Label
	startRow    int
	visibleRows int
	face        *text.GoTextFace

	slider *widget.Slider // Slider for the memory viewer
}

// initMemoryViewer initializes the memory viewer widget
func (d *Debugger) initMemoryViewer() *widget.Container {
	// Create container for the memory viewer
	container := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionHorizontal),
			widget.RowLayoutOpts.Spacing(4),
		)),
	)

	// Add title and memory view in a vertical container
	viewContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Spacing(2),
		)),
	)

	titleColor := &widget.LabelColor{
		Idle: color.RGBA{R: 255, G: 255, B: 200, A: 255},
	}
	titleLabel := widget.NewLabel(
		widget.LabelOpts.Text("Memory", d.face, titleColor),
	)
	viewContainer.AddChild(titleLabel)

	// Create memory viewer
	d.MemViewer = NewMemoryViewer(d.mem, d.face)
	viewContainer.AddChild(d.MemViewer.GetWidget())

	// Add the view container and slider
	container.AddChild(viewContainer)
	container.AddChild(d.MemViewer.GetSlider())

	return container
}

// updateMemoryViewer updates the memory viewer content
func (d *Debugger) updateMemoryViewer() {
	if d.MemViewer != nil {
		d.MemViewer.Update()
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

	// Create slider
	mv.slider = widget.NewSlider(
		// Set the slider orientation - n/s vs e/w
		widget.SliderOpts.Direction(widget.DirectionVertical),
		// Set the minimum and maximum value for the slider
		widget.SliderOpts.MinMax(0, 4096-mv.visibleRows),
		// Set the current value of the slider, without triggering a change event
		widget.SliderOpts.InitialCurrent(0),
		widget.SliderOpts.WidgetOpts(
			// Set the Widget to stretch vertically
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionEnd,
				Stretch:  true,
			}),
			// Set minimum width for the slider
			widget.WidgetOpts.MinSize(15, 0),
		),
		widget.SliderOpts.Images(
			// Set the track images
			&widget.SliderTrackImage{
				Idle:  ebitenimage.NewNineSliceColor(color.NRGBA{100, 100, 100, 255}),
				Hover: ebitenimage.NewNineSliceColor(color.NRGBA{100, 100, 100, 255}),
			},
			// Set the handle images
			&widget.ButtonImage{
				Idle:    ebitenimage.NewNineSliceColor(color.NRGBA{255, 100, 100, 255}),
				Hover:   ebitenimage.NewNineSliceColor(color.NRGBA{255, 100, 100, 255}),
				Pressed: ebitenimage.NewNineSliceColor(color.NRGBA{255, 100, 100, 255}),
			},
		),
		// Set the size of the handle
		widget.SliderOpts.FixedHandleSize(15),
		// Set the offset to display the track
		widget.SliderOpts.TrackOffset(0),
		// Set the size to move the handle
		widget.SliderOpts.PageSizeFunc(func() int {
			return 16 // Move 16 rows at a time
		}),
	)

	// Initial update
	mv.Update()

	return mv
}

// GetWidget returns the root widget of the memory viewer
func (mv *MemoryViewer) GetWidget() widget.PreferredSizeLocateableWidget {
	return mv.scrollCont
}

// GetSlider returns the slider widget
func (mv *MemoryViewer) GetSlider() widget.PreferredSizeLocateableWidget {
	return mv.slider
}

// Update updates the memory viewer's content and handles input
func (mv *MemoryViewer) Update() {
	// Update the start row based on the slider value here and not in a specific handler to avoid race conditions
	mv.startRow = mv.slider.Current

	// Update each row's content
	for i := 0; i < mv.visibleRows; i++ {
		addr := 16 * (mv.startRow + i)
		var hexBytes strings.Builder
		var asciiChars strings.Builder

		// Build hex representation
		for j := 0; j < 16; j++ {
			value := mv.mem.DebugRead(uint16(addr + j))
			hexBytes.WriteString(fmt.Sprintf("%02X ", value))

			// Add an additional space between the first 8 bytes and the last 8 bytes
			if j == 7 {
				hexBytes.WriteString(" ")
			}

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

// Scroll scrolls the memory view up/down
func (mv *MemoryViewer) Scroll(xCursor, yCursor int, yWheel float64) {
	rect := mv.container.GetWidget().Rect
	if yWheel != 0 && image.Pt(xCursor, yCursor).In(rect) {
		newStartRow := mv.startRow + int(yWheel*yWheel*yWheel)
		if newStartRow < 0 {
			newStartRow = 0
		} else if newStartRow > 4096-16 {
			newStartRow = 4096 - 16
		}
		mv.startRow = newStartRow
		mv.slider.Current = newStartRow
	}
}
