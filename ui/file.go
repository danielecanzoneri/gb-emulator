package ui

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/danielecanzoneri/gb-emulator/gameboy/cartridge"
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

	ui.gameTitle = rom.Header().Title
	ui.fileName = romPath

	return nil
}

func getSavFileName(romPath string) string {
	// Remove gb extension
	savFile := romPath[:len(romPath)-len(filepath.Ext(romPath))]
	return savFile + ".sav"
}

// Create an audio file containing the game audio.
// The file will have the same name of the game, followed by a timestamp and .dat extension.
func createAudioFile(romName string) (*os.File, error) {
	baseName := romName[:len(romName)-len(filepath.Ext(romName))]
	timestamp := time.Now().Format("20060102_150405")
	audioFileName := fmt.Sprintf("%s_%s.dat", baseName, timestamp)

	// Create the file
	file, err := os.Create(audioFileName)
	if err != nil {
		return nil, err
	}

	return file, nil
}
