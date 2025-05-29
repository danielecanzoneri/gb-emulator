package main

import (
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"github.com/danielecanzoneri/gb-emulator/debugger/ui"
)

func main() {
	// TODO: Inizializza il client (WebSocket/TCP)
	// TODO: Crea l'applicazione Fyne
	// TODO: Configura l'UI del debugger
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
	myApp.Run()
}
