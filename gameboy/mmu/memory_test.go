package mmu

import (
	"github.com/danielecanzoneri/gb-emulator/gameboy/cartridge"
	"testing"
)

func TestMemoryReadWrite(t *testing.T) {
	mem := &MMU{Cartridge: cartridge.NewMBC1(make([]uint8, 0x8000), true, nil, &cartridge.Header{ROMBanks: 1, RAMBanks: 1}, false)}
	mem.Write(0, 0x0A) // Enable RAM

	// Test writing and reading a byte
	addr := uint16(0xA034)
	value := uint8(0xAB)
	mem.Write(addr, value)

	readValue := mem.Read(addr)
	if readValue != value {
		t.Errorf("Expected %X, got %X", value, readValue)
	}

	// Test writing and reading a word
	wordAddr := uint16(0xB078)
	wordValue := uint16(0xCDEF)
	mem.WriteWord(wordAddr, wordValue)

	readWordValue := mem.ReadWord(wordAddr)
	if readWordValue != wordValue {
		t.Errorf("Expected %X, got %X", wordValue, readWordValue)
	}
}
