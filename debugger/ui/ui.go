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

	active         bool
	stepButton     *widget.Button
	continueButton *widget.Button

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
	ui.disassembler = newDisassembler(func(addr uint16, breakpointState bool) {
		debugger.Breakpoint(addr, breakpointState)
	})

	// Create the memory viewer list
	ui.memoryViewer = newMemoryViewer()

	// Create the register viewer
	ui.registersViewer = newRegisterViewer()

	// Debug buttons
	ui.stepButton = widget.NewButtonWithIcon(
		"Step",
		theme.Icon(theme.IconNameNavigateNext),
		func() {
			ui.debugger.Step()
		},
	)
	ui.continueButton = widget.NewButtonWithIcon(
		"Continue",
		theme.Icon(theme.IconNameMediaSkipNext),
		func() {
			ui.debugger.Continue()

			// Remove highlight from current instruction
			ui.disassembler.previousEntry.currentInstruction = false
			fyne.Do(ui.disassembler.Refresh)

			ui.SetActive(false)
		},
	)
	buttons := container.NewHBox(ui.stepButton, ui.continueButton)

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

func (ui *UI) BreakpointHit() {
	ui.SetActive(true)
}

func (ui *UI) SetActive(active bool) {
	ui.active = active
	if active {
		ui.stepButton.Enable()
		ui.continueButton.Enable()
		ui.disassembler.Enable()
	} else {
		ui.stepButton.Disable()
		ui.continueButton.Disable()
		ui.disassembler.Disable()
	}
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
