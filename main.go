package main

import (
	"log"

	"github.com/sqweek/dialog"
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

	ui, err := NewUI(romPath)
	if err != nil {
		log.Fatal(err)
	}

	ui.Init()
	ui.Start()
}
