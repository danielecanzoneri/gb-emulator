package main

import (
	"log"

	"github.com/danielecanzoneri/gb-emulator/debugger"

	"github.com/danielecanzoneri/gb-emulator/gameboy/cartridge"
	"github.com/ebitengine/oto/v3"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sqweek/dialog"
)

// Global reference to the audio player to keep it active
var audioPlayer *oto.Player

func main() {
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

	gb, player := Init()

	// Load the ROM
	rom, err := cartridge.LoadROM(romPath)
	if err != nil {
		log.Fatalf("Error loading the cartridge: %v", err)
	}
	gb.Load(rom)

	// Create Debugger
	gb.debugger = debugger.NewDebugger(gb.Memory, gb.CPU)

	// Initialize CPU
	gb.Reset()

	// Keep a reference to the audio player
	audioPlayer = player
	audioPlayer.Play()

	// Since game boy is 59.7 FPS but ebiten updates at 60 FPS there are
	// some frames where nothing is drawn. This avoids screen flickering
	ebiten.SetScreenClearedEveryFrame(false)

	// Initialize the renderer
	RenderInit()

	// Initial window size without the debug panel
	ebiten.SetWindowSize(gb.Layout(0, 0))

	if err := ebiten.RunGame(gb); err != nil {
		log.Fatal(err)
	}
}
