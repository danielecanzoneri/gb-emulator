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
	case 0:
		return NewMBC0(romData, header)
	case 1, 2:
		return NewMBC1(romData, nil, header, false)
	case 3:
		return NewMBC1(romData, savData, header, true)
	case 5:
		return NewMBC2(romData, nil, header, false)
	case 6:
		return NewMBC2(romData, savData, header, true)
	case 0x0F, 0x10:
		return NewMBC3(romData, savData, header, true, true)
	case 0x11, 0x12:
		return NewMBC3(romData, nil, header, false, false)
	case 0x13:
		return NewMBC3(romData, savData, header, true, false)
	default:
		log.Panicf("cartridge type %02X not supported", romData[cartridgeType])
		return nil
	}
}
