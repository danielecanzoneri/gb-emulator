package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

func NewRegisterViewer() fyne.CanvasObject {
	soundGrid := container.NewGridWithColumns(
		2,
		createCh1Panel(),
		createCh2Panel(),
		createCh3Panel(),
		createCh4Panel(),
		createSoundControlPanel(),
		createWaveRamPanel(),
	)

	leftContainer := container.NewVBox(
		createCPUPanel(),
		createInterruptsPanel(),
		createLCDPanel(),
		createTimerPanel(),
	)

	return container.NewHBox(
		leftContainer,
		soundGrid,
	)
}

func createCPUPanel() *Panel {
	panel := NewPanel("CPU")
	panel.AddRow("AF", "0000")
	panel.AddRow("BC", "0000")
	panel.AddRow("DE", "0000")
	panel.AddRow("HL", "0000")
	panel.AddRow("PC", "0000")
	panel.AddRow("SP", "0000")
	return panel
}

func createInterruptsPanel() *Panel {
	panel := NewPanel("Interrupts")
	panel.AddRow("FF0F IF", "00")
	panel.AddRow("FFFF IE", "00")
	panel.AddRow("IME", "disabled")
	return panel
}

func createLCDPanel() *Panel {
	panel := NewPanel("LCD")
	panel.AddRow("FF40 LCDC", "00")
	panel.AddRow("FF41 STAT", "00")
	panel.AddRow("FF42 SCY", "00")
	panel.AddRow("FF43 SCX", "00")
	panel.AddRow("FF44 LY", "00")
	panel.AddRow("FF45 LYC", "00")
	panel.AddRow("FF46 DMA", "00")
	panel.AddRow("FF47 BGP", "00")
	panel.AddRow("FF48 OBP0", "00")
	panel.AddRow("FF49 OBP1", "00")
	panel.AddRow("FF4A WY", "00")
	panel.AddRow("FF4B WX", "00")
	return panel
}

func createCh1Panel() *Panel {
	panel := NewPanel("Sound Channel 1")
	panel.AddRow("FF10 NR10", "00")
	panel.AddRow("FF11 NR11", "00")
	panel.AddRow("FF12 NR12", "00")
	panel.AddRow("FF13 NR13", "00")
	panel.AddRow("FF14 NR14", "00")
	return panel
}

func createCh2Panel() *Panel {
	panel := NewPanel("Sound Channel 2")
	panel.AddRow("FF15 NR20", "--")
	panel.AddRow("FF16 NR21", "00")
	panel.AddRow("FF17 NR22", "00")
	panel.AddRow("FF18 NR23", "00")
	panel.AddRow("FF19 NR24", "00")
	return panel
}

func createCh3Panel() *Panel {
	panel := NewPanel("Sound Channel 3")
	panel.AddRow("FF1A NR30", "00")
	panel.AddRow("FF1B NR31", "00")
	panel.AddRow("FF1C NR32", "00")
	panel.AddRow("FF1D NR33", "00")
	panel.AddRow("FF1E NR34", "00")
	return panel
}

func createCh4Panel() *Panel {
	panel := NewPanel("Sound Channel 4")
	panel.AddRow("FF1F NR40", "--")
	panel.AddRow("FF20 NR41", "00")
	panel.AddRow("FF21 NR42", "00")
	panel.AddRow("FF22 NR43", "00")
	panel.AddRow("FF23 NR44", "00")
	return panel
}

func createSoundControlPanel() *Panel {
	panel := NewPanel("Sound Control")
	panel.AddRow("FF24 NR50", "00")
	panel.AddRow("FF25 NR51", "00")
	panel.AddRow("FF26 NR52", "00")
	return panel
}

func createWaveRamPanel() *Panel {
	panel := NewPanel("Wave RAM")
	panel.AddRow("", "00 00 00 00")
	panel.AddRow("", "00 00 00 00")
	panel.AddRow("", "00 00 00 00")
	panel.AddRow("", "00 00 00 00")
	return panel
}

func createTimerPanel() *Panel {
	panel := NewPanel("Timer")
	panel.AddRow("FF04 DIV", "00")
	panel.AddRow("FF05 TIMA", "00")
	panel.AddRow("FF06 TMA", "00")
	panel.AddRow("FF07 TAC", "00")
	return panel
}
