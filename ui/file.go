package ui

import (
	"errors"
	"fmt"
	"github.com/danielecanzoneri/lucky-boy/gameboy"
	"github.com/danielecanzoneri/lucky-boy/ui/graphics"
	"log"
	"os"
	"path/filepath"

	"github.com/danielecanzoneri/lucky-boy/gameboy/cartridge"
	"github.com/sqweek/dialog"
)

func (ui *UI) Save() {
	ramDump := ui.GameBoy.Memory.Cartridge.RAMDump()
	if ramDump == nil {
		return
	}

	savFile := getSavFileName(ui.fileName)
	err := os.WriteFile(savFile, ramDump, 0644)
	if err != nil {
		log.Println("error writing game save:", err)
	}
}

func (ui *UI) LoadBootROM(bootRom string) (err error) {
	var data []uint8

	if bootRom != "" {
		// Open the ROM file
		data, err = os.ReadFile(bootRom)
	}

	ui.GameBoy.LoadBootROM(data)
	return
}

func (ui *UI) AskRomPath() (string, error) {
	romPath, err := dialog.File().
		Filter("Game Boy ROMs", "gb", "gbc").
		Title("Choose a GameBoy ROM").
		Load()
	if err != nil {
		return "", err
	}

	// Check if the ROM path is provided
	if romPath == "" {
		return "", errors.New("ROM file path is required")
	}

	return romPath, nil
}

func (ui *UI) LoadROM(romPath string) error {
	// Open the ROM file
	cartridgeData, err := os.ReadFile(romPath)
	if err != nil {
		return err
	}

	// Open the SAV file
	savFile := getSavFileName(romPath)
	savData, err := os.ReadFile(savFile)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			savData = nil
		} else {
			return err
		}
	}

	rom := cartridge.NewCartridge(cartridgeData, savData)
	ui.GameBoy.Load(rom)

	if ui.GameBoy.EmulationModel == gameboy.DMG {
		ui.palette = theme.DMGPalette{}
	} else {
		ui.palette = theme.CGBPalette{}
	}

	ui.gameTitle = rom.Header().Title
	ui.fileName = romPath

	return nil
}

func (ui *UI) SetModel(model string) error {
	switch model {
	case "auto":
		ui.GameBoy.Model = gameboy.Auto
	case "dmg":
		ui.GameBoy.Model = gameboy.DMG
	case "cgb":
		ui.GameBoy.Model = gameboy.CGB
	default:
		return fmt.Errorf("invalid model type: %s", model)
	}

	return nil
}

func getSavFileName(romPath string) string {
	// Remove gb extension
	savFile := romPath[:len(romPath)-len(filepath.Ext(romPath))]
	return savFile + ".sav"
}
