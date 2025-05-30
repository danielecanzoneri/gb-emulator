package main

import (
	"log"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"github.com/danielecanzoneri/gb-emulator/debugger/internal/client"
	"github.com/danielecanzoneri/gb-emulator/debugger/ui"
)

func main() {
	debuggerClient := client.New("localhost", 8080)

	// Connect to the emulator
	if err := debuggerClient.Connect(); err != nil {
		log.Printf("Failed to connect to emulator: %v\n", err)
	}

	myApp := app.New()
	myApp.Settings().SetTheme(new(ui.GameBoyTheme))
	window := myApp.NewWindow("GameBoy Disassembler")

	// Create the disassembler
	disassembler := ui.NewDisassembler()

	// Create the memory viewer list
	memoryViewer := ui.NewMemoryViewer()

	// Create the register viewer
	registerViewer := ui.NewRegisterViewer()

	// Create a split container with the disassembler on the left (40% of space)
	// and the register/memory viewers on the right (60% of space)
	split := container.NewHBox(
		disassembler,
		container.NewVBox(
			registerViewer,
			memoryViewer,
		),
	)

	window.SetContent(split)
	window.SetFixedSize(true)
	window.Show()

	// When the window is closed, disconnect from the emulator
	window.SetOnClosed(func() {
		if err := debuggerClient.Disconnect(); err != nil {
			log.Printf("Error disconnecting from emulator: %v\n", err)
		}
	})

	myApp.Run()
}
