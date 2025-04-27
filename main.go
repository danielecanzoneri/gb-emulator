package main

import (
	"flag"
	"fmt"
	gameboy "github.com/danielecanzoneri/gb-emulator/internal"
	"github.com/danielecanzoneri/gb-emulator/internal/cartridge"
	"github.com/danielecanzoneri/gb-emulator/internal/cpu"
)

// Define a flag for the ROM file path
var romPath = flag.String("rom", "", "Path to the ROM file")

// Define a flag for the debug mode
var debugMode = flag.Bool("debug", false, "Enable debug mode")

func main() {
	flag.Parse()

	// Check if the ROM path is provided
	if *romPath == "" {
		fmt.Println("Error: ROM file path is required")
		flag.Usage()
		return
	}

	// Enable debug mode if specified
	cpu.Debug = *debugMode

	gb := gameboy.Init()

	// Load the ROM
	rom, err := cartridge.LoadROM(*romPath)
	if err != nil {
		fmt.Printf("Error loading the cartridge: %v\n", err)
		return
	}
	gb.Load(rom)

	// Initialize cpu
	gb.Reset()

	gb.Run()
}
