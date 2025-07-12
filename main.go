package main

import (
	"flag"
	"log"

	"github.com/danielecanzoneri/gb-emulator/ui"
)

var startWithDebugger = flag.Bool("debug", false, "Start emulator with debugger enabled")

func main() {
	flag.Parse()

	// Init emulator
	gui, err := ui.New()
	if err != nil {
		log.Fatal(err)
	}

	gui.LoadNewGame()

	if *startWithDebugger {
		gui.ToggleDebugger()
	}
	gui.Run()
}
