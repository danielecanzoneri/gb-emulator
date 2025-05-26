package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func (ui *UI) handleInput() {
	// Step next instruction
	if ui.debugging && inpututil.IsKeyJustPressed(ebiten.KeyF3) {
		ui.stepInstruction = true
	}

	// Ctrl+P to pause
	if inpututil.IsKeyJustPressed(ebiten.KeyP) && ebiten.IsKeyPressed(ebiten.KeyControl) {
		ui.Pause()
	}

	// ESC key to enter in debug mode
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		ui.ToggleDebugger()

		// Update window size when debugger is enabled/disabled
		width, height := ui.Layout(0, 0)
		ebiten.SetWindowSize(width, height)
	}

	if inpututil.IsKeyJustPressed(ebiten.Key1) {
		ui.gameBoy.APU.Ch1Enabled = !ui.gameBoy.APU.Ch1Enabled
		debugString := "Channel 1 "
		if ui.gameBoy.APU.Ch1Enabled {
			debugString += "enabled"
		} else {
			debugString += "muted"
		}
		ui.debugString = debugString
		ui.debugStringTimer = 60
	}

	if inpututil.IsKeyJustPressed(ebiten.Key2) {
		ui.gameBoy.APU.Ch2Enabled = !ui.gameBoy.APU.Ch2Enabled
		debugString := "Channel 2 "
		if ui.gameBoy.APU.Ch2Enabled {
			debugString += "enabled"
		} else {
			debugString += "muted"
		}
		ui.debugString = debugString
		ui.debugStringTimer = 60
	}

	if inpututil.IsKeyJustPressed(ebiten.Key3) {
		ui.gameBoy.APU.Ch3Enabled = !ui.gameBoy.APU.Ch3Enabled
		debugString := "Channel 3 "
		if ui.gameBoy.APU.Ch3Enabled {
			debugString += "enabled"
		} else {
			debugString += "muted"
		}
		ui.debugString = debugString
		ui.debugStringTimer = 60
	}

	if inpututil.IsKeyJustPressed(ebiten.Key4) {
		ui.gameBoy.APU.Ch4Enabled = !ui.gameBoy.APU.Ch4Enabled
		debugString := "Channel 4 "
		if ui.gameBoy.APU.Ch4Enabled {
			debugString += "enabled"
		} else {
			debugString += "muted"
		}
		ui.debugString = debugString
		ui.debugStringTimer = 60
	}
}
