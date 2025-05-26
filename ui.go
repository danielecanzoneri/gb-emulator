package main

import (
	"github.com/danielecanzoneri/gb-emulator/debugger"
	"github.com/danielecanzoneri/gb-emulator/gameboy"
	"github.com/ebitengine/oto/v3"
	"github.com/hajimehoshi/ebiten/v2"
	"log"
)

type UI struct {
	gameBoy   *gameboy.GameBoy
	gameTitle string

	// Audio player
	audioBuffer chan float32
	audioPlayer *oto.Player

	// True if we are in debug mode
	debugging       bool
	debugger        *debugger.Debugger
	stepInstruction bool // Set when in debug we must execute one instruction

	paused bool

	debugString      string
	debugStringTimer uint
}

func NewUI(romPath string) (*UI, error) {
	ui := &UI{}

	// Create audio buffer
	ui.audioBuffer = make(chan float32, bufferSize)
	gb := gameboy.New(ui.audioBuffer, sampleRate)
	ui.gameBoy = gb

	// Load the ROM
	gameTitle, err := gb.Load(romPath)
	if err != nil {
		return nil, err
	}

	ui.gameTitle = gameTitle

	// Create Debugger
	ui.debugger = debugger.NewDebugger(gb.Memory, gb.CPU)

	// Create audio player
	err = ui.initAudioPlayer()
	if err != nil {
		return nil, err
	}

	return ui, nil
}

func (ui *UI) Init() {
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
}

func (ui *UI) Start() {
	ui.audioPlayer.Play()

	if err := ebiten.RunGame(ui); err != nil {
		log.Fatal(err)
	}
}

func (ui *UI) Pause() {
	ui.paused = !ui.paused

	if ui.paused {
		ebiten.SetWindowTitle(ui.gameTitle + " (paused)")
	} else {
		ebiten.SetWindowTitle(ui.gameTitle)
	}
}

// ToggleDebugger enables/disables visualization of I/O registers
func (ui *UI) ToggleDebugger() {
	ui.debugger.ToggleVisibility()
	ui.debugging = !ui.debugging
}
