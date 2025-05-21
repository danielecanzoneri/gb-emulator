package main

import (
	"log"

	gameboy "github.com/danielecanzoneri/gb-emulator/internal"
	"github.com/danielecanzoneri/gb-emulator/internal/cartridge"
	"github.com/ebitengine/oto/v3"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sqweek/dialog"
)

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

	gb, player := gameboy.Init()

	// Load the ROM
	rom, err := cartridge.LoadROM(romPath)
	if err != nil {
		log.Fatalf("Error loading the cartridge: %v", err)
	}
	gb.Load(rom)

	// Initialize cpu
	gb.Reset()

	// Keep a reference to the audio player
	audioPlayer = player
	audioPlayer.Play()

	// Since game boy is 59.7 FPS but ebiten updates at 60 FPS there are
	// some frames where nothing is drawn. This avoids screen flickering
	ebiten.SetScreenClearedEveryFrame(false)
	ebiten.SetWindowSize(gb.Layout(0, 0))

	gameboy.RenderInit()
	if err := ebiten.RunGame(gb); err != nil {
		log.Fatal(err)
	}
}
