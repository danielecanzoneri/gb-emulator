package ui

import (
	"log"

	"github.com/danielecanzoneri/gb-emulator/ui/debugger"

	"github.com/danielecanzoneri/gb-emulator/gameboy"
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
	ui.debugger = debugger.New(gb)

	// Create audio player
	player, err := newAudioPlayer(ui)
	if err != nil {
		return nil, err
	}

	ui.audioPlayer = player

	// Initialize the renderer
	ui.initRenderer()

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
