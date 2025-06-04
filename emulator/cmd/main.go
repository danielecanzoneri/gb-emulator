package main

import (
	"github.com/danielecanzoneri/gb-emulator/emulator/internal/server"
	"log"

	"github.com/danielecanzoneri/gb-emulator/emulator/ui"
	"github.com/hajimehoshi/ebiten/v2"
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

	debugServer := new(server.Server)

	// Init emulator
	gui, err := ui.New(debugServer)
	if err != nil {
		log.Fatal(err)
	}

	err = gui.Load(romPath)
	if err != nil {
		log.Fatal(err)
	}

	// Start debugging server
	debugServer.Start("8080")
	defer debugServer.Close()

	// Start the game loop
	if err := ebiten.RunGame(gui); err != nil {
		log.Fatal(err)
	}
}
