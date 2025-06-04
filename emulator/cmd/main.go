package main

import (
	"github.com/danielecanzoneri/gb-emulator/emulator/internal/server"
	"log"

	"github.com/danielecanzoneri/gb-emulator/emulator/ui"
	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	debugServer := new(server.Server)

	// Init emulator
	gui, err := ui.New(debugServer)
	if err != nil {
		log.Fatal(err)
	}

	gui.LoadNewGame()

	// Start debugging server
	debugServer.Start("8080")
	defer debugServer.Close()

	// Start the game loop
	if err := ebiten.RunGame(gui); err != nil {
		log.Fatal(err)
	}
}
