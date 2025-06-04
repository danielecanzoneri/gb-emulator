package ui

import (
	"github.com/danielecanzoneri/gb-emulator/emulator/internal/gameboy"
	"github.com/danielecanzoneri/gb-emulator/emulator/internal/server"
	"github.com/ebitengine/oto/v3"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sqweek/dialog"
	"log"
	"os/exec"
)

type UI struct {
	gameBoy   *gameboy.GameBoy
	gameTitle string

	// Audio player
	audioBuffer chan float32
	audioPlayer *oto.Player

	debugString      string
	debugStringTimer uint

	// Debugger
	debuggerCmd *exec.Cmd

	DebugState *server.Server
}

func New(s *server.Server) (*UI, error) {
	ui := new(UI)

	// Create audio buffer
	ui.audioBuffer = make(chan float32, bufferSize)
	gb := gameboy.New(ui.audioBuffer, sampleRate)
	ui.gameBoy = gb

	// Debugger
	ui.DebugState = s
	s.SetDebugger(
		gb.CPU, gb.Memory,
	)
	s.OnStep = ui.gameBoy.CPU.ExecuteInstruction

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

	return ui, nil
}

func (ui *UI) LoadNewGame() {
	romPath, err := dialog.File().
		Filter("Game Boy ROMs", "gb", "bin").
		Title("Choose a GameBoy ROM").
		Load()
	if err != nil {
		log.Fatal(err)
	}

	// Check if the ROM path is provided
	if romPath == "" {
		log.Fatal("Error: ROM file path is required")
	}

	// Load the ROM
	gameTitle, err := ui.gameBoy.Load(romPath)
	if err != nil {
		log.Fatal(err)
	}

	ui.gameTitle = gameTitle
}
