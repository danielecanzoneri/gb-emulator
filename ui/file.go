package ui

import (
	"errors"
	"github.com/danielecanzoneri/gb-emulator/gameboy/cartridge"
	"github.com/sqweek/dialog"
	"log"
	"os"
	"path/filepath"
)

func (ui *UI) LoadNewGame() {
	romPath, err := dialog.File().
		Filter("Game Boy ROMs", "gb", "gbc").
		Title("Choose a GameBoy ROM").
		Load()
	if err != nil {
		log.Fatal(err)
	}

	// Check if the ROM path is provided
	if romPath == "" {
		log.Fatal("Error: ROM file path is required")
	}

	// Load the ROM
	rom, err := loadROM(romPath)
	if err != nil {
		log.Fatal(err)
	}
	ui.gameBoy.Load(rom)

	ui.gameTitle = rom.Header().Title
	ui.fileName = romPath
}

func (ui *UI) Save() {
	ramDump := ui.gameBoy.Memory.Cartridge.RAMDump()
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

	ui.gameBoy.LoadBootROM(data)
	return
}

func loadROM(romPath string) (cartridge.Cartridge, error) {
	// Open the ROM file
	cartridgeData, err := os.ReadFile(romPath)
	if err != nil {
		return nil, err
	}

	// Open the SAV file
	savFile := getSavFileName(romPath)
	savData, err := os.ReadFile(savFile)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			savData = nil
		} else {
			return nil, err
		}
	}

	return cartridge.NewCartridge(cartridgeData, savData), nil
}

func getSavFileName(romPath string) string {
	// Remove gb extension
	savFile := romPath[:len(romPath)-len(filepath.Ext(romPath))]
	return savFile + ".sav"
}
