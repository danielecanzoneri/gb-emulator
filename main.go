package main

import (
	gameboy "github.com/danielecanzoneri/gb-emulator/internal"
	"github.com/danielecanzoneri/gb-emulator/internal/cartridge"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sqweek/dialog"
	"log"
)

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

	gb := gameboy.Init()

	// Load the ROM
	rom, err := cartridge.LoadROM(romPath)
	if err != nil {
		log.Fatalf("Error loading the cartridge: %v", err)
	}
	gb.Load(rom)

	// Initialize cpu
	gb.Reset()

	ebiten.SetWindowTitle(romPath)
	ebiten.SetWindowSize(gb.Layout(0, 0))

	gameboy.RenderInit()
	if err := ebiten.RunGame(gb); err != nil {
		log.Fatal(err)
	}
}
