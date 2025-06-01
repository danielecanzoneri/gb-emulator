package main

import (
	"github.com/danielecanzoneri/gb-emulator/pkg/debug"
	"log"

	"github.com/danielecanzoneri/gb-emulator/debugger/internal/client"
	"github.com/danielecanzoneri/gb-emulator/debugger/ui"
)

func main() {
	debuggerClient := client.New("localhost", 8080)

	gui := ui.New(debuggerClient)
	debuggerClient.StateConsumer = func(s *debug.GameBoyState) { gui.Update(s) }

	// Connect to the emulator
	if err := debuggerClient.Connect(); err != nil {
		log.Printf("Failed to connect to emulator: %v\n", err)
	}
	gui.Run()

	// When the app is closed, disconnect from the emulator
	if err := debuggerClient.Disconnect(); err != nil {
		log.Printf("Error disconnecting from emulator: %v\n", err)
	}
}
