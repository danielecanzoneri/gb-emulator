package ui

import (
	"github.com/danielecanzoneri/gb-emulator/emulator/internal/gameboy"
	"github.com/ebitengine/oto/v3"
	"github.com/hajimehoshi/ebiten/v2"
)

type UI interface {
	ebiten.Game

	Load(string) error
}

type ui struct {
	gameBoy   *gameboy.GameBoy
	gameTitle string

	// Audio player
	audioBuffer chan float32
	audioPlayer *oto.Player

	debugString      string
	debugStringTimer uint
}

func New() (UI, error) {
	ui := &ui{}

	// Create audio buffer
	ui.audioBuffer = make(chan float32, bufferSize)
	gb := gameboy.New(ui.audioBuffer, sampleRate)
	ui.gameBoy = gb

	// Create audio player
	player, err := newAudioPlayer(ui)
	if err != nil {
		return nil, err
	}

	ui.audioPlayer = player

	// Since game boy is 59.7 FPS but ebiten updates at 60 FPS there are
	// some frames where nothing is drawn. This avoids screen flickering
	ebiten.SetScreenClearedEveryFrame(false)

	// Initialize the renderer
	initRenderer()

	// Initial window size without the debug panel
	screenWidth, screenHeight := ui.Layout(0, 0)
	ebiten.SetWindowSize(screenWidth, screenHeight)

	// Initialize CPU
	ui.gameBoy.Reset()

	return ui, nil
}

func (ui *ui) Load(romPath string) error {
	// Load the ROM
	gameTitle, err := ui.gameBoy.Load(romPath)
	if err != nil {
		return err
	}

	ui.gameTitle = gameTitle
	return nil
}
