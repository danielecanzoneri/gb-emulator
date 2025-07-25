package debugger

import (
	"github.com/ebitenui/ebitenui/input"
	"image"
)

func (d *Debugger) Toggle() {
	d.Active = !d.Active
	if d.Active {
		d.Stop()
	}
}

func (d *Debugger) CheckBreakpoint(addr uint16) bool {
	return d.disassembler.IsBreakpoint(addr)
}

// Run commands

func (d *Debugger) Step() {
	defer d.Sync()

	d.gameBoy.CPU.ExecuteInstruction()
}

func (d *Debugger) Next() {
	d.Continue()
	d.NextInstruction = true
	d.CallDepth = 0
}

func (d *Debugger) Continue() {
	d.Continued = true

	// Unselect current entry
	d.disassembler.currentInstruction = -1

	// TODO Disable control buttons
	d.disassembler.refresh()
}

func (d *Debugger) Stop() {
	defer d.Sync()

	d.Continued = false
	// TODO Enable control buttons
}

func (d *Debugger) Reset() {
	defer d.Sync()

	d.gameBoy.Reset()
	d.initHooks()
}

func (d *Debugger) initHooks() {
	callHook := func() {
		d.CallDepth++
	}
	retHook := func() {
		d.CallDepth--
	}
	d.gameBoy.CPU.SetHooks(callHook, retHook)
}

// PPU commands

func (d *Debugger) ShowOAM() {
	if d.IsWindowOpen(d.oamViewer.Window) {
		return
	}

	// Current window size
	winSize := input.GetWindowSize()

	// Get the preferred size of the content
	x1, y1 := d.oamViewer.Contents.PreferredSize()
	x2, y2 := d.oamViewer.TitleBar.PreferredSize()
	xWindow, yWindow := max(x1, x2), y1+y2
	remainingSize := winSize.Sub(image.Pt(xWindow, yWindow))

	// Set the windows location at center of window
	r := image.Rect(0, 0, xWindow, yWindow).Add(remainingSize.Div(2))
	d.oamViewer.SetLocation(r)

	closeWindow := d.AddWindow(d.oamViewer.Window)
	d.oamViewer.closeWindow = closeWindow
}
