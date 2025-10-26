package main

import (
	"flag"
	"github.com/danielecanzoneri/gb-emulator/ui"
	"log"
)

const (
	socketPort = "4321"
)

var (
	startWithDebugger = flag.Bool("debug", false, "Start emulator with debugger enabled")
	bootRom           = flag.String("boot-rom", "", "Boot ROM filename")
	romPath           = flag.String("rom", "", "ROM filename")
	recordAudio       = flag.Bool("record", false, "Record game audio (2 channels uncompressed 32-bit float little endian")
	serial            = flag.String("serial", "", "Serial role (master or slave)")
	systemModel       = flag.String("model", "auto", "GameBoy model (auto, dmg, cgb)")
)

func main() {
	flag.Parse()

	// Init emulator
	gui, err := ui.New()
	if err != nil {
		log.Fatal(err)
	}

	if *romPath == "" {
		*romPath, err = gui.AskRomPath()
		if err != nil {
			log.Fatal(err)
		}
	}

	if err = gui.SetModel(*systemModel); err != nil {
		log.Fatal(err)
	}

	err = gui.LoadROM(*romPath)
	if err != nil {
		log.Fatal(err)
	}

	// Load Boot ROM
	if err = gui.LoadBootROM(*bootRom); err != nil {
		log.Fatal(err)
	}

	// Serial port data exchange
	switch *serial {
	case "master":
		gui.Listen(socketPort)
	case "slave":
		gui.Connect(socketPort)
	case "":
	default:
		log.Printf("Invalid serial role %q", *serial)
	}

	if *recordAudio {
		if filename, err := gui.RecordAudio(); err != nil {
			log.Println("Could not record game audio:", err)
		} else {
			log.Println("Recording game audio to ", filename)
		}
	}
	if *startWithDebugger {
		gui.ToggleDebugger()
	}
	gui.Run()
}
