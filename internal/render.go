package gameboy

import (
	"github.com/danielecanzoneri/gb-emulator/internal/ppu"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"image/color"
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

// Inherit Ebiten Game interface

func (gb *GameBoy) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		gb.Pause()
	}
	// Game updates are called in the audio callback function
	return nil
}

func (gb *GameBoy) Draw(screen *ebiten.Image) {
	if !gb.PPU.FrameComplete {
		return
	}

	gb.PPU.FrameComplete = false

	if gb.PPU.EmptyFrame {
		screen.Fill(color.White)
		gb.PPU.EmptyFrame = false
		return
	}

	// Draw all 144 x 160 pixels
	for y := range ppu.FrameHeight {
		for x := range ppu.FrameWidth {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(Scale*x), float64(Scale*y))

			colorId := gb.PPU.Framebuffer[y][x]
			screen.DrawImage(pixels[colorId], op)
		}
	}
}

func (gb *GameBoy) Layout(_, _ int) (int, int) {
	return Scale * ppu.FrameWidth, Scale * ppu.FrameHeight
}
