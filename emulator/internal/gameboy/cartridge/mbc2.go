package cartridge

import (
	"github.com/danielecanzoneri/gb-emulator/pkg/util"
	"log"
)

const (
	mbc2Len = 512
)

type MBC2 struct {
	header  *Header
	battery bool

	ROMBanks uint8

	ROM []uint8
	RAM [mbc2Len]uint8 // 512 half-bytes builtin RAM

	// Registers 0000-3FFF
	ramEnabled    bool
	romBankNumber uint8
}

func (mbc *MBC2) RAMDump() []uint8 {
	if mbc.battery {
		return mbc.RAM[:]
	}

	return nil
}

func (mbc *MBC2) Header() *Header {
	return mbc.header
}

func NewMBC2(rom []uint8, ram []uint8, header *Header, battery bool) *MBC2 {
	mbc := &MBC2{
		header:        header,
		battery:       battery,
		ROMBanks:      uint8(header.ROMBanks),
		ROM:           rom,
		romBankNumber: 1,
	}

	if ram != nil {
		if len(ram) != mbc2Len {
			log.Println("[WARN] sav file was of a different dimension than expected, resetting to zero")
		} else {
			copy(mbc.RAM[:], ram)
		}
	}
	return mbc
}

func (mbc *MBC2) Write(addr uint16, value uint8) {
	// Set MBC2 registers
	switch {
	case addr < 0x4000:
		// - When bit 8 is clear, the value that is written controls whether the RAM is enabled.
		//   Save RAM will be enabled if and only if the lower 4 bits of the value written here are $A.
		// - When bit 8 is set, the value that is written controls the selected ROM bank at 4000–7FFF.
		//   Specifically, the lower 4 bits of the value written to this address range specify the ROM bank number.
		//   If bank 0 is written, the resulting bank will be bank 1 instead.
		if util.ReadBit(addr, 8) == 0 {
			// Low nibble = 0xA enables RAM
			mbc.ramEnabled = value&0x0F == 0xA
		} else {
			// Only lower 5 bits are used
			value = value & 0x0F

			// Register 0 behaves as 1
			if value == 0 {
				value = 1
			}

			mbc.romBankNumber = value
		}

	case addr < 0x8000:
		// Nothing

	case 0xA000 <= addr && addr < 0xC000:
		if mbc.ramEnabled {
			RAMAddress := mbc.computeRamAddress(addr)
			mbc.RAM[RAMAddress] = 0xF0 | (value & 0x0F)
		}

	default:
		log.Printf("[WARN] MBC2 Write address is out of range: %04X\n", addr)
	}
}

func (mbc *MBC2) Read(addr uint16) uint8 {
	switch {
	case addr < 0x8000:
		cartridgeAddress := mbc.computeRomAddress(addr)
		return mbc.ROM[cartridgeAddress]

	case 0xA000 <= addr && addr < 0xC000:
		if mbc.ramEnabled {
			RAMAddress := mbc.computeRamAddress(addr)
			return 0xF0 | mbc.RAM[RAMAddress]
		}

	default:
		log.Printf("[WARN] MBC2 Read address is out of range: %04X\n", addr)
	}

	return 0xFF
}

func (mbc *MBC2) computeRomAddress(cpuAddress uint16) uint {
	var bankNumber uint8 = 0

	switch {
	case cpuAddress < 0x4000:
	case cpuAddress < 0x8000:
		bankNumber = mbc.romBankNumber
	default:
		panic("should never happen")
	}

	// Bank number is masked to the required number of bits
	bankNumber %= mbc.ROMBanks

	return uint(bankNumber)<<14 | uint(cpuAddress&0x3FFF)
}

func (mbc *MBC2) computeRamAddress(cpuAddress uint16) uint {
	// A200–BFFF — 15 “echoes” of A000–A1FF
	// Only the bottom 9 bits of the address are used to index into the internal RAM, so RAM access repeats

	return uint(cpuAddress & 0x01FF)
}
