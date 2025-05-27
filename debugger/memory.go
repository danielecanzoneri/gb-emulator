package debugger

import (
	"fmt"
	"image"
	"image/color"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"

	ebitenimage "github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// MemoryViewer represents the memory visualization widget
// It is a wrapper of a widget.Container
type MemoryViewer struct {
	labels []*widget.Label
	face   *text.GoTextFace

	mem         MemoryDebugger
	startRow    int
	visibleRows int

	rootContent *widget.Container
	scrollArea  *widget.Container // Area containing the rows
	slider      *widget.Slider    // Slider for the memory viewer
}

// initRootContainer initializes the memory viewer widget
func (mv *MemoryViewer) initRootContainer() {
	// Create container for the memory viewer
	mv.rootContent = widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Spacing(4),
		)),
	)

	// Add memory view and slider in a horizontal container
	viewContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionHorizontal),
			widget.RowLayoutOpts.Spacing(2),
		)),
	)

	// Add the scroll area and the slider
	viewContainer.AddChild(mv.scrollArea)
	viewContainer.AddChild(mv.slider)

	// Add the title and the view container
	titleColor := &widget.LabelColor{
		Idle: color.RGBA{R: 255, G: 255, B: 200, A: 255},
	}
	titleLabel := widget.NewLabel(
		widget.LabelOpts.Text("Memory", mv.face, titleColor),
	)
	mv.rootContent.AddChild(titleLabel)
	mv.rootContent.AddChild(viewContainer)
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
	mv.scrollArea = widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Spacing(2),
		)),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.MinSize(800, 300),
		),
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
		mv.scrollArea.AddChild(label)
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
			widget.WidgetOpts.MinSize(15, 300),
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

	mv.initRootContainer()

	// Initial update
	mv.Update()

	return mv
}

func (mv *MemoryViewer) GetWidget() *widget.Widget {
	return mv.rootContent.GetWidget()
}

func (mv *MemoryViewer) PreferredSize() (int, int) {
	return mv.rootContent.PreferredSize()
}

func (mv *MemoryViewer) SetLocation(rect image.Rectangle) {
	mv.rootContent.SetLocation(rect)
}

func (mv *MemoryViewer) Render(screen *ebiten.Image) {
	mv.rootContent.Render(screen)
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

			// Add a space between the first 8 bytes and the last 8 bytes
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

	mv.rootContent.Update()
}

// Scroll scrolls the memory view up/down
func (mv *MemoryViewer) Scroll(xCursor, yCursor int, yWheel float64) {
	rect := mv.scrollArea.GetWidget().Rect
	if yWheel != 0 && image.Pt(xCursor, yCursor).In(rect) {
		newStartRow := mv.startRow - int(yWheel*yWheel*yWheel)
		if newStartRow < 0 {
			newStartRow = 0
		} else if newStartRow > 4096-16 {
			newStartRow = 4096 - 16
		}
		mv.startRow = newStartRow
		mv.slider.Current = newStartRow
	}
}
