package ui

import (
	"github.com/danielecanzoneri/gb-emulator/gameboy/ppu"
	"github.com/danielecanzoneri/gb-emulator/ui/theme"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	Scale = 3
)

var (
	// Original palette
	frameImage *ebiten.Image
)

func (ui *UI) initRenderer() {
	// Since game boy is 59.7 FPS but ebiten updates at 60 FPS there are
	// some frames where nothing is drawn. This avoids screen flickering
	ebiten.SetScreenClearedEveryFrame(false)

	// Create a single image for the entire frame
	frameImage = ebiten.NewImage(ppu.FrameWidth, ppu.FrameHeight)

	// Initial window size without the debug panel
	screenWidth, screenHeight := ui.Layout(0, 0)
	ebiten.SetWindowSize(screenWidth, screenHeight)

	// Save when closing
	ebiten.SetWindowClosingHandled(true)
}

// Inherit Ebiten Game interface

func (ui *UI) Update() error {
	// If window is unfocused, stop the game
	if !ebiten.IsFocused() {
		ui.Paused = true
	} else {
		ui.Paused = false
	}

	// If closing, save game
	if ebiten.IsWindowBeingClosed() {
		ui.Save()
		return ebiten.Termination
	}

	ui.handleInput()

	if ui.debugger.Active {
		ebiten.SetWindowTitle(ui.gameTitle + " (debugging)")
		return ui.debugger.Update()
	} else {
		ebiten.SetWindowTitle(ui.gameTitle)
		// Game updates are called in the audio callback function
		return nil
	}
}

func (ui *UI) Draw(screen *ebiten.Image) {
	// Update the frame image with the current frame in the PPU
	frameBuffer := ui.gameBoy.PPU.GetFrame()
	for y := range ppu.FrameHeight {
		for x := range ppu.FrameWidth {
			colorId := frameBuffer[y][x]
			frameImage.Set(x, y, theme.GameBoyPalette[colorId])
		}
	}

	if ui.debugger.Active {
		ui.debugger.Draw(screen, frameImage)
		return
	}

	// Draw the entire frame at once with scaling
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(Scale, Scale)
	screen.DrawImage(frameImage, op)

	if ui.debugStringTimer > 0 {
		ebitenutil.DebugPrint(screen, ui.debugString)
		ui.debugStringTimer--
	}
}

func (ui *UI) Layout(_, _ int) (int, int) {
	// Adjust the layout based on whether the debugger is visible
	if ui.debugger.Active {
		return ui.debugger.Layout(0, 0)
	} else {
		return Scale * ppu.FrameWidth, Scale * ppu.FrameHeight
	}
}
