package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
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

	// Debug buttons
	stepButton := widget.NewButtonWithIcon(
		"Step",
		theme.Icon(theme.IconNameNavigateNext),
		func() {
			ui.debugger.Step()
		},
	)
	buttons := container.NewHBox(stepButton)

	// Create a container with the disassembler on the left
	// and the register/memory viewers on the right
	split := container.NewHBox(
		ui.disassembler,
		container.NewVBox(
			ui.registersViewer,
			buttons,
			ui.memoryViewer,
		),
	)

	ui.window.SetContent(split)
	ui.window.SetFixedSize(true)

	return ui
}

func (ui *UI) Update(state *debug.GameBoyState) {
	ui.disassembler.Update(state)
	ui.memoryViewer.Update(state)
	ui.registersViewer.Update(state)
}

func (ui *UI) Run() {
	ui.window.Show()
	ui.app.Run()
}
