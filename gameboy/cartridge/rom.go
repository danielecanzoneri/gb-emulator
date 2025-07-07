package cartridge

import (
	"log"
)

type Cartridge interface {
	Read(uint16) uint8
	Write(uint16, uint8)

	RAMDump() []uint8
	Header() *Header
}

func NewCartridge(romData []uint8, savData []uint8) Cartridge {
	header := parseHeader(romData)

	switch romData[cartridgeType] {
	case 0: // ROM ONLY
		return NewMBC0(romData, header)
	case 1: // MBC1
		return NewMBC1(romData, false, nil, header, false)
	case 2: // MBC1 + RAM
		return NewMBC1(romData, true, nil, header, false)
	case 3: // MBC1 + RAM + BATTERY
		return NewMBC1(romData, true, savData, header, true)
	case 5: // MBC2
		return NewMBC2(romData, nil, header, false)
	case 6: // MBC2 + BATTERY
		return NewMBC2(romData, savData, header, true)
	case 0x0F: // MBC3 + TIMER + BATTERY
		return NewMBC3(romData, false, savData, header, true, true)
	case 0x10: // MBC3 + TIMER + RAM + BATTERY
		return NewMBC3(romData, true, savData, header, true, true)
	case 0x11: // MBC3
		return NewMBC3(romData, false, nil, header, false, false)
	case 0x12: // MBC3 + RAM
		return NewMBC3(romData, true, nil, header, false, false)
	case 0x13: // MBC3 + RAM + BATTERY
		return NewMBC3(romData, true, savData, header, true, false)
	case 0x19: // MBC5
		return NewMBC5(romData, false, nil, header, false, false)
	case 0x1A: // MBC5 + RAM
		return NewMBC5(romData, true, nil, header, false, false)
	case 0x1B: // MBC5 + RAM + BATTERY
		return NewMBC5(romData, true, savData, header, true, false)
	case 0x1C: // MBC5 + RUMBLE
		return NewMBC5(romData, false, nil, header, false, true)
	case 0x1D: // MBC5 + RUMBLE + RAM
		return NewMBC5(romData, true, nil, header, false, true)
	case 0x1E: // MBC5 + RUMBLE + RAM + BATTERY
		return NewMBC5(romData, true, savData, header, true, true)
	default:
		log.Panicf("cartridge type %02X not supported", romData[cartridgeType])
		return nil
	}
}
