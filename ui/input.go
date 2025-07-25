package ui

import (
	"github.com/danielecanzoneri/gb-emulator/ui/debugger"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func (ui *UI) handleInput() {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		ui.ToggleDebugger()
	}

	// Ctrl+L to load a new game
	if inpututil.IsKeyJustPressed(ebiten.KeyL) && ebiten.IsKeyPressed(ebiten.KeyControl) {
		// Stop running
		ui.Paused = true

		// Save game before switching
		ui.Save()

		ui.LoadNewGame()
		ui.gameBoy.Reset()

		// Start running
		ui.Paused = false
	}

	ui.handleAudioToggle()

	// Handle debugger input
	for _, handler := range debugger.InputHandlers {
		handler()
	}
}

func (ui *UI) handleAudioToggle() {
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

func (ui *UI) ToggleDebugger() {
	ui.debugger.Toggle()

	// Resize window
	newWidth, newHeight := ui.Layout(0, 0)
	ebiten.SetWindowSize(newWidth, newHeight)
}
