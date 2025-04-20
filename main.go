package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/danielecanzoneri/gb-emulator/cpu"
	"github.com/danielecanzoneri/gb-emulator/memory"
	"github.com/danielecanzoneri/gb-emulator/rom"
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

	mem := &memory.Memory{}
	gb := cpu.CPU{
		Mem: mem,
	}

	// Load the ROM
	romData, err := rom.LoadROM(*romPath)
	if err != nil {
		fmt.Printf("Error loading the rom: %v\n", err)
		return
	}

	// Load ROM into memory
	for i, b := range romData {
		mem.Write(uint16(i), b)
	}

	// Gameboy doctor file
	file, err := os.OpenFile("dump.txt", os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		panic(err)
	}
	cpu.DebugFile = file
	defer file.Close()

	// Initialize cpu
	gb.Reset()

	// Simplified loop
	for {
		gb.ExecuteInstruction()
	}
}
