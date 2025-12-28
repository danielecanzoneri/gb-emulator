package ui

import (
	"github.com/danielecanzoneri/lucky-boy/ui/graphics"
	"log"
	"os"

	"github.com/danielecanzoneri/lucky-boy/ui/debugger"

	"github.com/danielecanzoneri/lucky-boy/gameboy"
	"github.com/ebitengine/oto/v3"
	"github.com/hajimehoshi/ebiten/v2"
)

type UI struct {
	GameBoy   *gameboy.GameBoy
	gameTitle string
	fileName  string

	// When true, stop emulation
	Paused bool

	// Audio player
	audioBuffer chan float32
	audioPlayer *oto.Player
	audioFile   *os.File // Record game audio

	// Turbo mode
	turbo bool

	debugString      string
	debugStringTimer uint

	// Color Palette
	palette theme.Palette

	// CGB color correction shader
	Shader     *ebiten.Shader
	shaderOpts *ebiten.DrawRectShaderOptions

	// Debugger
	debugger *debugger.Debugger
}

func New(useShader bool) (*UI, error) {
	ui := new(UI)

	// Create audio buffer
	ui.audioBuffer = make(chan float32, bufferSize)
	gb := gameboy.New(ui.audioBuffer, sampleRate)
	ui.GameBoy = gb

	// Debugger
	ui.debugger = debugger.New(gb)

	// Create audio player
	player, err := newAudioPlayer(ui)
	if err != nil {
		return nil, err
	}

	ui.audioPlayer = player

	// Initialize the renderer
	ui.initRenderer(useShader)

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
