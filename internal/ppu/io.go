package ppu

import (
	"github.com/danielecanzoneri/gb-emulator/internal/util"
	"strconv"
)

const (
	LCDCAddr = 0xFF40
	STATAddr = 0xFF41
	SCYAddr  = 0xFF42
	SCXAddr  = 0xFF43
	LYAddr   = 0xFF44
	LYCAddr  = 0xFF45
	DMAAddr  = 0xFF46
	BGPAddr  = 0xFF47
	OBP0Addr = 0xFF48
	OBP1Addr = 0xFF49
	WYAddr   = 0xFF4A
	WXAddr   = 0xFF4B

	// STATMask bits 0-2 read only, bits 3-6 read/write
	STATMask = 0b01111000
)

func (ppu *PPU) Write(addr uint16, v uint8) {
	switch addr {
	case LCDCAddr:
		//
		if ppu.active && ppu.Mode != vBlank && util.ReadBit(ppu.LCDC, 7) > 0 {
			panic("Cannot disable LCD outside of VBlank period")
		}
		ppu.LCDC = v
		// Update LCD control
		ppu.active = util.ReadBit(v, 7) > 0
		ppu.windowTileMapArea = util.ReadBit(v, 6)
		ppu.windowEnabled = util.ReadBit(v, 5) > 0
		ppu.bgWindowTileDataArea = util.ReadBit(v, 4)
		ppu.bgTileMapArea = util.ReadBit(v, 3)
		ppu.objSize = util.ReadBit(v, 2)
		ppu.objEnabled = util.ReadBit(v, 1) > 0
		ppu.bgWindowEnabled = util.ReadBit(v, 0) > 0
	case LYAddr:
		panic("cannot write LY address")
	case LYCAddr:
		ppu.LYC = v
	case STATAddr:
		ppu.STAT = (STATMask & v) | (ppu.STAT &^ STATMask)
		ppu.checkSTATInterruptState()
	default:
		panic("PPU: unknown addr " + strconv.Itoa(int(addr)))
	}
}

func (ppu *PPU) Read(addr uint16) uint8 {
	switch addr {
	case LCDCAddr:
		return ppu.LCDC
	case LYAddr:
		return ppu.LY
	case LYCAddr:
		return ppu.LYC
	case STATAddr:
		return ppu.STAT
	default:
		panic("PPU: unknown addr " + strconv.Itoa(int(addr)))
	}
}
