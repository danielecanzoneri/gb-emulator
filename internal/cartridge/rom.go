package cartridge

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
)

func LoadROM(path string) ([]byte, error) {
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
	return os.ReadFile(path)
}
