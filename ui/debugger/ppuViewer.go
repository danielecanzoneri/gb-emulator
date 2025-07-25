package debugger

import (
	"github.com/danielecanzoneri/gb-emulator/gameboy"
	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/widget"
)

type oamViewer struct {
	*widget.Window

	// Pointer to the UI for showing the window
	ui *ebitenui.UI

	// Handler to close the window
	closeWindow widget.RemoveWindowFunc
}

func (d *Debugger) newOamViewer() *oamViewer {
	o := &oamViewer{ui: d.UI}

	root := widget.NewContainer(widget.ContainerOpts.WidgetOpts(
		widget.WidgetOpts.MinSize(400, 400),
	))

	o.Window = newWindow("OAM Viewer", root, &o.closeWindow)
	return o
}

func (o *oamViewer) Sync(gb *gameboy.GameBoy) {
	if !o.ui.IsWindowOpen(o.Window) {
		return
	}
}
