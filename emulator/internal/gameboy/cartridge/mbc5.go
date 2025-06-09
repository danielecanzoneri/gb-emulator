package cartridge

import (
	"github.com/danielecanzoneri/gb-emulator/pkg/util"
	"log"
)

type MBC5 struct {
	header  *Header
	battery bool // If battery is present RAM should be stored
	rumble  bool // If rumble motor is present on cartridge

	ROMBanks uint // Up to 512
	RAMBanks uint8

	ROM []uint8
	RAM []uint8

	// Registers
	ramEnabled    bool
	romBankNumber uint  // 9 bit register
	ramBankNumber uint8 // 4 bit register
}

func (mbc *MBC5) RAMDump() []uint8 {
	if mbc.battery {
		return mbc.RAM
	}

	return nil
}

func (mbc *MBC5) Header() *Header {
	return mbc.header
}

func NewMBC5(rom []uint8, ram []uint8, header *Header, battery bool, rumble bool) *MBC5 {
	mbc := &MBC5{
		header:        header,
		battery:       battery,
		rumble:        rumble,
		ROMBanks:      header.ROMBanks,
		RAMBanks:      uint8(header.RAMBanks),
		ROM:           rom,
		romBankNumber: 1,
	}

	switch {
	case len(ram) != int(header.RAMBanks*0x2000):
		log.Println("[WARN] sav file was of a different dimension than expected, resetting to zero")
		fallthrough
	case ram == nil:
		ram = make([]uint8, header.RAMBanks*0x2000)
	}
	mbc.RAM = ram

	return mbc
}

func (mbc *MBC5) Write(addr uint16, value uint8) {
	// Set MBC5 registers
	switch {
	case addr < 0x2000:
		// Low nibble = 0xA enables RAM
		mbc.ramEnabled = value&0x0F == 0xA

	case addr < 0x3000:
		// 8 least significant bits of ROM bank number
		mbc.romBankNumber &= 0x100
		mbc.romBankNumber |= uint(value)

	case addr < 0x4000:
		// 9th bit of ROM bank number
		mbc.romBankNumber &= 0xFF
		mbc.romBankNumber |= uint(value&1) << 8

	case addr < 0x6000:
		// Only lower 4 bits are used
		mbc.ramBankNumber = value & 0xF

		if mbc.rumble {
			if util.ReadBit(value, 3) == 1 {
				// Rumble
			} else {
				// Stop rumble
			}
		}

	case addr < 0x8000:
		// Nothing

	case 0xA000 <= addr && addr < 0xC000:
		if mbc.ramEnabled {
			RAMAddress := mbc.computeRamAddress(addr)
			mbc.RAM[RAMAddress] = value
		}

	default:
		log.Printf("[WARN] MBC5 Write address is out of range: %04X\n", addr)
	}
}

func (mbc *MBC5) Read(addr uint16) uint8 {
	switch {
	case addr < 0x8000:
		cartridgeAddress := mbc.computeRomAddress(addr)
		return mbc.ROM[cartridgeAddress]

	case 0xA000 <= addr && addr < 0xC000:
		if mbc.ramEnabled {
			RAMAddress := mbc.computeRamAddress(addr)
			return mbc.RAM[RAMAddress]
		}

	default:
		log.Printf("[WARN] MBC5 Read address is out of range: %04X\n", addr)
	}

	return 0xFF
}

func (mbc *MBC5) computeRomAddress(cpuAddress uint16) uint {
	// bank number: 9 bits, cpuAddress: 14 bits
	var bankNumber uint = 0

	switch {
	case cpuAddress < 0x4000: // bank number = 0
	case cpuAddress < 0x8000:
		bankNumber = mbc.romBankNumber

	default:
		panic("should never happen")
	}

	// Bank number is masked to the required number of bits
	bankNumber %= mbc.ROMBanks

	return bankNumber<<14 | uint(cpuAddress&0x3FFF)
}

func (mbc *MBC5) computeRamAddress(cpuAddress uint16) uint {
	// bank number: 2 bits, cpuAddress: 13 bits
	switch {
	case 0xA000 <= cpuAddress && cpuAddress < 0xC000:
		cpuAddress = cpuAddress & 0x1FFF

		// Bank number is masked to the required number of bits
		bank := mbc.ramBankNumber % mbc.RAMBanks
		return uint(bank)<<13 | uint(cpuAddress)

	default:
		panic("should never happen")
	}
}
