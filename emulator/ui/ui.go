package ui

import (
	"github.com/danielecanzoneri/gb-emulator/emulator/ui/debugger"
	"log"

	"github.com/danielecanzoneri/gb-emulator/emulator/internal/gameboy"
	"github.com/ebitengine/oto/v3"
	"github.com/hajimehoshi/ebiten/v2"
)

type UI struct {
	gameBoy   *gameboy.GameBoy
	gameTitle string
	fileName  string

	// Audio player
	audioBuffer chan float32
	audioPlayer *oto.Player

	debugString      string
	debugStringTimer uint

	// Debugger
	debugger *debugger.Debugger
}

func New() (*UI, error) {
	ui := new(UI)

	// Create audio buffer
	ui.audioBuffer = make(chan float32, bufferSize)
	gb := gameboy.New(ui.audioBuffer, sampleRate)
	ui.gameBoy = gb

	// Debugger
	ui.debugger = debugger.New()

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

	// Save when closing
	ebiten.SetWindowClosingHandled(true)

	return ui, nil
}

func (ui *UI) Run() {
	// Start audio player the first time
	ui.audioPlayer.Play()

	// Start the game loop
	if err := ebiten.RunGame(ui); err != nil {
		log.Fatal(err)
	}
}
