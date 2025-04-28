package main

import (
	"github.com/sqweek/dialog"

	gameboy "github.com/danielecanzoneri/gb-emulator/internal"
	"github.com/danielecanzoneri/gb-emulator/internal/cartridge"
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
		return
	}
	gb.Load(rom)

	// Initialize cpu
	gb.Reset()

	gb.Run()
}
