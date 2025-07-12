package main

import (
	"flag"
	"log"

	"github.com/danielecanzoneri/gb-emulator/ui"
)

var (
	startWithDebugger = flag.Bool("debug", false, "Start emulator with debugger enabled")
	bootRom           = flag.String("boot-rom", "boot/bootix_dmg.bin", "Boot ROM filename (\"None\" to skip boot ROM)")
)

func main() {
	flag.Parse()

	// Init emulator
	gui, err := ui.New()
	if err != nil {
		log.Fatal(err)
	}

	gui.LoadNewGame()

	// Load Boot ROM
	if *bootRom == "None" {
		*bootRom = ""
	}
	if err = gui.LoadBootROM(*bootRom); err != nil {
		log.Fatal(err)
	}

	if *startWithDebugger {
		gui.ToggleDebugger()
	}
	gui.Run()
}
