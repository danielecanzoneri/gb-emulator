package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func (ui *ui) handleInput() {
	ui.handleAudioToggle()
}

func (ui *ui) handleAudioToggle() {
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
