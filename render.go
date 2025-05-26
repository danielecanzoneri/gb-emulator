package main

import (
	"image/color"

	"github.com/danielecanzoneri/gb-emulator/internal/ppu"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
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

func RenderInit() {
	for i := range pixels {
		// Create a Scale x Scale image of the corresponding color
		square := ebiten.NewImage(Scale, Scale)
		square.Fill(colors[i])

		pixels[i] = square
	}
}

func (gb *GameBoy) handleKeys() {
	// Step next instruction
	if gb.debugging && inpututil.IsKeyJustPressed(ebiten.KeyF3) {
		gb.stepInstruction = true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyP) && ebiten.IsKeyPressed(ebiten.KeyControl) {
		gb.Pause()
	}

	// ESC key to enter in debug mode
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		gb.ToggleDebugger()

		// Update window size when debugger is enabled/disabled
		width, height := gb.Layout(0, 0)
		ebiten.SetWindowSize(width, height)
	}

	if inpututil.IsKeyJustPressed(ebiten.Key1) {
		gb.APU.Ch1Enabled = !gb.APU.Ch1Enabled
		debugString := "Channel 1 "
		if gb.APU.Ch1Enabled {
			debugString += "enabled"
		} else {
			debugString += "muted"
		}
		gb.debugString = debugString
		gb.debugStringTimer = 60
	}

	if inpututil.IsKeyJustPressed(ebiten.Key2) {
		gb.APU.Ch2Enabled = !gb.APU.Ch2Enabled
		debugString := "Channel 2 "
		if gb.APU.Ch2Enabled {
			debugString += "enabled"
		} else {
			debugString += "muted"
		}
		gb.debugString = debugString
		gb.debugStringTimer = 60
	}

	if inpututil.IsKeyJustPressed(ebiten.Key3) {
		gb.APU.Ch3Enabled = !gb.APU.Ch3Enabled
		debugString := "Channel 3 "
		if gb.APU.Ch3Enabled {
			debugString += "enabled"
		} else {
			debugString += "muted"
		}
		gb.debugString = debugString
		gb.debugStringTimer = 60
	}

	if inpututil.IsKeyJustPressed(ebiten.Key4) {
		gb.APU.Ch4Enabled = !gb.APU.Ch4Enabled
		debugString := "Channel 4 "
		if gb.APU.Ch4Enabled {
			debugString += "enabled"
		} else {
			debugString += "muted"
		}
		gb.debugString = debugString
		gb.debugStringTimer = 60
	}
}

// Inherit Ebiten Game interface

func (gb *GameBoy) Update() error {
	gb.handleKeys()

	// Update the debugger
	gb.debugger.Update()

	// Game updates are called in the audio callback function
	return nil
}

func (gb *GameBoy) Draw(screen *ebiten.Image) {
	// Background color
	screen.Fill(color.RGBA{R: 40, G: 40, B: 40, A: 220})

	// Draw the debugger
	gb.debugger.Draw(screen, Scale*ppu.FrameWidth, 0)

	// Draw all 144 x 160 pixels
	for y := range ppu.FrameHeight {
		for x := range ppu.FrameWidth {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(Scale*x), float64(Scale*y))

			colorId := gb.PPU.Framebuffer[y][x]
			screen.DrawImage(pixels[colorId], op)
		}
	}

	if gb.debugStringTimer > 0 {
		ebitenutil.DebugPrint(screen, gb.debugString)
		gb.debugStringTimer--
	}
}

func (gb *GameBoy) Layout(_, _ int) (int, int) {
	// Adjust the layout based on whether the debugger is visible
	if gb.debugger.IsVisible() {
		debugWidth, debugHeight := gb.debugger.Layout(0, 0)
		return Scale*ppu.FrameWidth + debugWidth, max(Scale*ppu.FrameHeight, debugHeight)
	}
	return Scale * ppu.FrameWidth, Scale * ppu.FrameHeight
}
