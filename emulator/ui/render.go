package ui

import (
	"image/color"

	"github.com/danielecanzoneri/gb-emulator/emulator/internal/gameboy/ppu"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	Scale = 3
)

var (
	// Original palette
	GameBoyScreenPalette = [4]color.Color{
		color.RGBA{R: 198, G: 222, B: 140, A: 255},
		color.RGBA{R: 132, G: 165, B: 99, A: 255},
		color.RGBA{R: 57, G: 97, B: 57, A: 255},
		color.RGBA{R: 8, G: 24, B: 16, A: 255},
	}
	frameImage *ebiten.Image
)

func initRenderer() {
	// Create a single image for the entire frame
	frameImage = ebiten.NewImage(ppu.FrameWidth, ppu.FrameHeight)
}

// Inherit Ebiten Game interface

func (ui *UI) Update() error {
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
			frameImage.Set(x, y, GameBoyScreenPalette[colorId])
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
