package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"github.com/danielecanzoneri/gb-emulator/debugger/internal/client"
	"github.com/danielecanzoneri/gb-emulator/pkg/debug"
)

type UI struct {
	app    fyne.App
	window fyne.Window

	// Components
	disassembler    *disassembler
	memoryViewer    *memoryViewer
	registersViewer *registersViewer

	debugger *client.Client
}

func New(debugger *client.Client) *UI {
	ui := &UI{
		debugger: debugger,
	}

	ui.app = app.New()
	ui.app.Settings().SetTheme(new(gameBoyTheme))
	ui.window = ui.app.NewWindow("GameBoy Disassembler")

	// Create the disassembler
	ui.disassembler = newDisassembler()

	// Create the memory viewer list
	ui.memoryViewer = newMemoryViewer()

	// Create the register viewer
	ui.registersViewer = newRegisterViewer()

	// Create a split container with the disassembler on the left (40% of space)
	// and the register/memory viewers on the right (60% of space)
	split := container.NewHBox(
		ui.disassembler,
		container.NewVBox(
			ui.registersViewer,
			ui.memoryViewer,
		),
	)

	ui.window.SetContent(split)
	ui.window.SetFixedSize(true)

	return ui
}

func (ui *UI) Update(state *debug.GameBoyState) {
	ui.registersViewer.Update(state)
}

func (ui *UI) Run() {
	ui.window.Show()
	ui.app.Run()
}
