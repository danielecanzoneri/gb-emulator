package cartridge

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
)

type Rom struct {
	// Struct containing cartridge information
	Header *Header

	Data []byte
}

func LoadROM(path string) (*Rom, error) {
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

	header, err := parseHeader(data)
	if err != nil {
		return nil, err
	}

	return &Rom{Header: header, Data: data}, nil
}
