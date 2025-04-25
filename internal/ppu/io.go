package ppu

import (
	"github.com/danielecanzoneri/gb-emulator/internal/util"
	"strconv"
)

const (
	vRAMStartAddr = 0x8000
	OAMStartAddr  = 0xFE00

	LCDCAddr = 0xFF40
	STATAddr = 0xFF41
	SCYAddr  = 0xFF42
	SCXAddr  = 0xFF43
	LYAddr   = 0xFF44
	LYCAddr  = 0xFF45
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
		if util.ReadBit(v, 6) == 0 {
			ppu.windowTileMapAddr = 0x9800
		} else {
			ppu.windowTileMapAddr = 0x9C00
		}
		ppu.windowEnabled = util.ReadBit(v, 5) > 0
		ppu.bgWindowTileDataArea = util.ReadBit(v, 4)
		if util.ReadBit(v, 3) == 0 {
			ppu.bgTileMapAddr = 0x9800
		} else {
			ppu.bgTileMapAddr = 0x9C00
		}
		ppu.obj8x16Size = util.ReadBit(v, 2) > 0
		ppu.objEnabled = util.ReadBit(v, 1) > 0
		ppu.bgWindowEnabled = util.ReadBit(v, 0) > 0
	case LYAddr:
		panic("cannot write LY address")
	case LYCAddr:
		ppu.LYC = v
	case STATAddr:
		ppu.STAT = (STATMask & v) | (ppu.STAT &^ STATMask)
		ppu.checkSTATInterruptState()
	case BGPAddr:
		ppu.BGP = Palette(v)
	case OBP0Addr:
		ppu.OBP0 = Palette(v)
	case OBP1Addr:
		ppu.OBP1 = Palette(v)
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
	case BGPAddr:
		return uint8(ppu.BGP)
	case OBP0Addr:
		return uint8(ppu.OBP0)
	case OBP1Addr:
		return uint8(ppu.OBP1)
	default:
		panic("PPU: unknown addr " + strconv.Itoa(int(addr)))
	}
}

// ReadVRAM prevents vRAM reads during PPU mode 3
func (ppu *PPU) ReadVRAM(addr uint16) uint8 {
	if ppu.Mode == 3 {
		return 0xFF
	}
	return ppu.vRAM.read(addr - vRAMStartAddr)
}

// WriteVRAM prevents vRAM writes during PPU mode 3
func (ppu *PPU) WriteVRAM(addr uint16, value uint8) {
	if ppu.Mode == 3 {
		return
	}
	ppu.vRAM.write(addr-vRAMStartAddr, value)
}

// ReadOAM prevents OAM reads during PPU mode 2 and 3
func (ppu *PPU) ReadOAM(addr uint16) uint8 {
	if ppu.Mode == 2 || ppu.Mode == 3 {
		return 0xFF
	}
	return ppu.OAM.read(addr - OAMStartAddr)
}

// WriteOAM prevents OAM writes during PPU mode 2 and 3
func (ppu *PPU) WriteOAM(addr uint16, value uint8) {
	if ppu.Mode == 2 || ppu.Mode == 3 {
		return
	}
	ppu.OAM.write(addr-OAMStartAddr, value)
}

func (ppu *PPU) DMAWrite(index int, value uint8) {
	ppu.OAM.data[index] = value
}
