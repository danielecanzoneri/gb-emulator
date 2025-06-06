package cartridge

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
)

type Cartridge interface {
	Read(uint16) uint8
	Write(uint16, uint8)

	Header() *Header
}

func LoadROM(path string) (Cartridge, error) {
	// Check if the ROM exists
	stat, err := os.Stat(path)
	if errors.Is(err, fs.ErrNotExist) {
		return nil, fmt.Errorf("ROM file does not exist: %s", path)
	}
	if errors.Is(err, fs.ErrPermission) {
		return nil, fmt.Errorf("permission denied to access ROM file: %s", path)
	}
	if stat.IsDir() {
		return nil, fmt.Errorf("ROM file is a directory: %s", path)
	}

	// Open the ROM file
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return NewCartridge(data), nil
}

func NewCartridge(data []byte) Cartridge {
	header := parseHeader(data)

	switch data[cartridgeType] {
	case 0:
		return NewMBC0(data, header)
	case 1, 2, 3:
		return NewMBC1(data, header)
	default:
		log.Panicf("cartridge type %02X not supported", data[cartridgeType])
		return nil
	}
}
