package ui

import (
	"github.com/danielecanzoneri/lucky-boy/ui/debugger"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"log"
)

func (ui *UI) handleInput() {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		ui.ToggleDebugger()
	}

	// Turbo (play at max speed)
	ui.turbo = ebiten.IsKeyPressed(ebiten.KeySpace)

	// Ctrl+L to load a new game
	if inpututil.IsKeyJustPressed(ebiten.KeyL) && ebiten.IsKeyPressed(ebiten.KeyControl) {
		// Stop running
		ui.Paused = true

		// Save game before switching
		ui.Save()

		romPath, err := ui.AskRomPath()
		if err == nil {
			err = ui.LoadROM(romPath)
		}
		if err != nil {
			log.Fatal(err)
		}

		ui.GameBoy.Reset()

		// Start running
		ui.Paused = false
	}

	ui.handleAudioToggle()

	// Handle debugger input
	if ui.debugger.Active {
		for _, handler := range debugger.InputHandlers {
			handler()
		}
	}
}

func (ui *UI) handleAudioToggle() {
	if inpututil.IsKeyJustPressed(ebiten.Key1) {
		ui.GameBoy.APU.Ch1Enabled = !ui.GameBoy.APU.Ch1Enabled
		debugString := "Channel 1 "
		if ui.GameBoy.APU.Ch1Enabled {
			debugString += "enabled"
		} else {
			debugString += "muted"
		}
		ui.debugString = debugString
		ui.debugStringTimer = 60
	}

	if inpututil.IsKeyJustPressed(ebiten.Key2) {
		ui.GameBoy.APU.Ch2Enabled = !ui.GameBoy.APU.Ch2Enabled
		debugString := "Channel 2 "
		if ui.GameBoy.APU.Ch2Enabled {
			debugString += "enabled"
		} else {
			debugString += "muted"
		}
		ui.debugString = debugString
		ui.debugStringTimer = 60
	}

	if inpututil.IsKeyJustPressed(ebiten.Key3) {
		ui.GameBoy.APU.Ch3Enabled = !ui.GameBoy.APU.Ch3Enabled
		debugString := "Channel 3 "
		if ui.GameBoy.APU.Ch3Enabled {
			debugString += "enabled"
		} else {
			debugString += "muted"
		}
		ui.debugString = debugString
		ui.debugStringTimer = 60
	}

	if inpututil.IsKeyJustPressed(ebiten.Key4) {
		ui.GameBoy.APU.Ch4Enabled = !ui.GameBoy.APU.Ch4Enabled
		debugString := "Channel 4 "
		if ui.GameBoy.APU.Ch4Enabled {
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
