package main

import (
	"github.com/danielecanzoneri/gb-emulator/emulator/internal/server"
	"log"

	"github.com/danielecanzoneri/gb-emulator/emulator/ui"
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

	gui.Run()
}
