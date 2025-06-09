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
	colors = [4]color.Color{
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

	titleSuffix := ""
	if ui.DebugState.IsActive() {
		titleSuffix = " (debugging)"
	}
	ebiten.SetWindowTitle(ui.gameTitle + titleSuffix)

	ui.handleInput()

	// Game updates are called in the audio callback function
	return nil
}

func (ui *UI) Draw(screen *ebiten.Image) {
	// Update the frame image with the current framebuffer
	for y := range ppu.FrameHeight {
		for x := range ppu.FrameWidth {
			colorId := ui.gameBoy.PPU.Framebuffer[y][x]
			frameImage.Set(x, y, colors[colorId])
		}
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
	return Scale * ppu.FrameWidth, Scale * ppu.FrameHeight
}
