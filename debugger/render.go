package debugger

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

// Update updates the debugger state
func (d *Debugger) Update() {
	if !d.visible {
		return
	}

	d.updatePanelContents()
	d.ui.Update()
}

// Draw renders the debug panel to the screen
func (d *Debugger) Draw(screen *ebiten.Image) {
	if !d.visible {
		return
	}

	d.ui.Draw(screen)
}

func (d *Debugger) Layout(_, _ int) (int, int) {
	w, h := d.rootCont.PreferredSize()
	return w, h
}

func (d *Debugger) GameboyScreenPosition() image.Rectangle {
	rect := d.GameboyScreen.GetWidget().Rect
	if !d.visible {
		// Translate to the origin
		return rect.Sub(rect.Min)
	}
	return rect
}
