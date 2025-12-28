package debugger

import (
	"image"

	"github.com/danielecanzoneri/lucky-boy/gameboy"
	"github.com/ebitenui/ebitenui/input"
	"github.com/ebitenui/ebitenui/widget"
)

// WindowViewer is an interface for debugger viewers that use windows
type WindowViewer interface {
	// Window returns the window widget
	Window() *widget.Window
	// Contents returns the content container of the window
	Contents() *widget.Container
	// TitleBar returns the title bar container of the window
	TitleBar() *widget.Container
	// Sync updates the viewer with GameBoy state
	Sync(gb *gameboy.GameBoy)
	// SetCloseHandler sets the function to close the window and returns the previous handler
	SetCloseHandler(closeFunc widget.RemoveWindowFunc) widget.RemoveWindowFunc
}

// showWindow is a helper function to show/hide a window viewer
func (d *Debugger) showWindow(viewer WindowViewer) {
	if d.UI.IsWindowOpen(viewer.Window()) {
		// Close the window by calling the close handler if set
		oldHandler := viewer.SetCloseHandler(nil)
		if oldHandler != nil {
			oldHandler()
		}
		return
	}

	// Open the window
	winSize := input.GetWindowSize()

	// Get the preferred size of the content
	x1, y1 := viewer.Contents().PreferredSize()
	x2, y2 := viewer.TitleBar().PreferredSize()
	xWindow, yWindow := max(x1, x2), y1+y2
	remainingSize := winSize.Sub(image.Pt(xWindow, yWindow))

	// Set the window location at center of window
	r := image.Rect(0, 0, xWindow, yWindow).Add(remainingSize.Div(2))
	viewer.Window().SetLocation(r)

	closeWindow := d.UI.AddWindow(viewer.Window())
	viewer.SetCloseHandler(closeWindow)

	// Sync data when opening
	viewer.Sync(d.gameBoy)
}
