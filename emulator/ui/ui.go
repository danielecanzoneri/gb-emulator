package ui

import (
	"github.com/danielecanzoneri/gb-emulator/emulator/internal/debugger"
	"github.com/danielecanzoneri/gb-emulator/emulator/internal/gameboy"
	"github.com/ebitengine/oto/v3"
	"github.com/hajimehoshi/ebiten/v2"
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

	DebugState *debugger.Debugger
}

func New() (*UI, error) {
	ui := new(UI)

	// Create audio buffer
	ui.audioBuffer = make(chan float32, bufferSize)
	gb := gameboy.New(ui.audioBuffer, sampleRate)
	ui.gameBoy = gb

	// Debugger
	ui.DebugState = debugger.NewDebugger(gb.CPU, gb.Memory)

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

func (ui *UI) Load(romPath string) error {
	// Load the ROM
	gameTitle, err := ui.gameBoy.Load(romPath)
	if err != nil {
		return err
	}

	ui.gameTitle = gameTitle
	return nil
}
