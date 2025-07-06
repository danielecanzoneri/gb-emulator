package main

import (
	"log"

	"github.com/danielecanzoneri/gb-emulator/emulator/ui"
)

func main() {
	// Init emulator
	gui, err := ui.New()
	if err != nil {
		log.Fatal(err)
	}

	gui.LoadNewGame()

	gui.Run()
}
