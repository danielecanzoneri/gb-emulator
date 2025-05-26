package main

import (
	"image/color"

	"github.com/danielecanzoneri/gb-emulator/gameboy/ppu"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	Scale = 3
)

var (
	colors = [4]color.Color{
		color.Gray{Y: 0xFF}, // White
		color.Gray{Y: 0xAA}, // Light gray
		color.Gray{Y: 0x55}, // Dark gray
		color.Gray{Y: 0},    // Black
	}
	pixels [4]*ebiten.Image
)

func initRenderer() {
	for i := range pixels {
		// Create a Scale x Scale image of the corresponding color
		square := ebiten.NewImage(Scale, Scale)
		square.Fill(colors[i])

		pixels[i] = square
	}
}

// Inherit Ebiten Game interface

func (ui *UI) Update() error {
	ui.handleInput()

	// Update the debugger
	ui.debugger.Update()

	// Game updates are called in the audio callback function
	return nil
}

func (ui *UI) Draw(screen *ebiten.Image) {
	// Background color
	screen.Fill(color.RGBA{R: 40, G: 40, B: 40, A: 220})

	// Draw the debugger
	ui.debugger.Draw(screen, Scale*ppu.FrameWidth, 0)

	// Draw all 144 x 160 pixels
	for y := range ppu.FrameHeight {
		for x := range ppu.FrameWidth {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(Scale*x), float64(Scale*y))

			colorId := ui.gameBoy.PPU.Framebuffer[y][x]
			screen.DrawImage(pixels[colorId], op)
		}
	}

	if ui.debugStringTimer > 0 {
		ebitenutil.DebugPrint(screen, ui.debugString)
		ui.debugStringTimer--
	}
}

func (ui *UI) Layout(_, _ int) (int, int) {
	// Adjust the layout based on whether the debugger is visible
	if ui.debugger.IsVisible() {
		debugWidth, debugHeight := ui.debugger.Layout(0, 0)
		return Scale*ppu.FrameWidth + debugWidth, max(Scale*ppu.FrameHeight, debugHeight)
	}
	return Scale * ppu.FrameWidth, Scale * ppu.FrameHeight
}
