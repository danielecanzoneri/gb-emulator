package ppu

import (
	"github.com/danielecanzoneri/gb-emulator/util"
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
	VBKAddr  = 0xFF4F

	BGPIAddr = 0xFF68
	BGPDAddr = 0xFF69
	OBPIAddr = 0xFF6A
	OBPDAddr = 0xFF6B

	// STATMask bits 0-2 read only, bits 3-6 read/write
	STATMask = 0b01111000
)

func (ppu *PPU) Write(addr uint16, v uint8) {
	if 0x8000 <= addr && addr < 0xA000 { // vRAM
		ppu.vRAM.Write(addr, v)
		return
	} else if 0xFE00 <= addr && addr < 0xFF00 { // OAM
		ppu.GlitchedOAMAccess(addr, false)
		ppu.oam.Write(addr, v)
		return
	}

	switch addr {
	case LCDCAddr:
		//if ppu.active && ppu.Mode != vBlank && util.ReadBit(ppu.LCDC, 7) > 0 {
		//	panic("Cannot disable LCD outside VBlank period")
		//}
		ppu.LCDC = v
		// Update LCD control
		if util.ReadBit(v, 7) > 0 {
			ppu.enable()
		} else {
			ppu.disable()
		}

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
		panic("should not write LY")
	case LYCAddr:
		ppu.LYC = v
		if ppu.active { // Changing LYC while PPU is inactive doesn't update STAT
			ppu.checkSTATInterrupt()
		}
	case STATAddr:
		// Spurious STAT interrupt: http://www.devrs.com/gb/files/faqs.html#GBBugs
		// Writing anything to the STAT register while the Game Boy is either in mode 0 or 1,
		// cause bit 1 of the IF register ($ff0f) to be set.
		if !ppu.STATInterruptLine && (ppu.STAT&3 < 2) {
			ppu.RequestSTATInterrupt()
		}
		ppu.STAT = (STATMask & v) | (ppu.STAT &^ STATMask)
		ppu.checkSTATInterrupt()
	// If in CGB compatibility mode, writes to this address must update the
	// corresponding CGB palettes
	case BGPAddr:
		ppu.BGP = DMGPalette(v)

	case OBP0Addr:
		ppu.OBP[0] = DMGPalette(v)
	case OBP1Addr:
		ppu.OBP[1] = DMGPalette(v)

	case SCYAddr:
		ppu.SCY = v
	case SCXAddr:
		ppu.SCX = v
	case WYAddr:
		ppu.WY = v
	case WXAddr:
		ppu.WX = v
	case VBKAddr:
		if ppu.Cgb {
			ppu.vRAM.bankNumber = v & 1
		}
	case BGPIAddr:
		if ppu.Cgb {
			ppu.BGPI = v
		}
	case BGPDAddr:
		if ppu.Cgb {
			// address bits are 0-5
			paletteAddr := ppu.BGPI & 0x3F
			ppu.BGPalette[paletteAddr] = v

			// Auto increment of address if bit 7 of BGPI is set
			if util.ReadBit(ppu.BGPI, 7) > 0 {
				ppu.BGPI = (ppu.BGPI & 0x80) | ((paletteAddr + 1) & 0x3F)
			}
		}
	case OBPIAddr:
		if ppu.Cgb {
			ppu.OBPI = v
		}
	case OBPDAddr:
		if ppu.Cgb {
			// address bits are 0-5
			paletteAddr := ppu.OBPI & 0x3F
			ppu.OBJPalette[paletteAddr] = v

			// Auto increment of address if bit 7 of OBPI is set
			if util.ReadBit(ppu.OBPI, 7) > 0 {
				ppu.OBPI = (ppu.OBPI & 0x80) | ((paletteAddr + 1) & 0x3F)
			}
		}
	default:
		panic("PPU: unknown addr " + strconv.FormatUint(uint64(addr), 16))
	}
}

func (ppu *PPU) Read(addr uint16) uint8 {
	if 0x8000 <= addr && addr < 0xA000 { // vRAM
		return ppu.vRAM.Read(addr)
	} else if 0xFE00 <= addr && addr < 0xFF00 { // OAM
		ppu.GlitchedOAMAccess(addr, true)
		return ppu.oam.Read(addr)
	}

	switch addr {
	case LCDCAddr:
		return ppu.LCDC
	case LYAddr:
		return ppu.LY
	case LYCAddr:
		return ppu.LYC
	case STATAddr:
		return 0x80 | ppu.STAT // Bit 7 is unused
	case BGPAddr:
		return uint8(ppu.BGP)
	case OBP0Addr:
		return uint8(ppu.OBP[0])
	case OBP1Addr:
		return uint8(ppu.OBP[1])
	case SCYAddr:
		return ppu.SCY
	case SCXAddr:
		return ppu.SCX
	case WYAddr:
		return ppu.WY
	case WXAddr:
		return ppu.WX
	case VBKAddr:
		return ppu.vRAM.bankNumber | 0xFE
	case BGPIAddr:
		return ppu.BGPI
	case BGPDAddr:
		// address bits are 0-5
		paletteAddr := ppu.BGPI & 0x3F
		return ppu.BGPalette[paletteAddr]
	case OBPIAddr:
		return ppu.OBPI
	case OBPDAddr:
		// address bits are 0-5
		paletteAddr := ppu.OBPI & 0x3F
		return ppu.OBJPalette[paletteAddr]
	default:
		panic("PPU: unknown addr " + strconv.FormatUint(uint64(addr), 16))
	}
}

func (ppu *PPU) enable() {
	if ppu.active {
		return
	}
	ppu.active = true
	ppu.dots = 0
	ppu.internalStateLength = 0

	ppu.setState(new(glitchedOamScan))
	ppu.checkSTATInterrupt()

	ppu.emptyFrame()
}

func (ppu *PPU) disable() {
	if !ppu.active {
		return
	}
	ppu.active = false

	// Reset mode
	ppu.STAT &= 0xFC

	// Reset to HBlank
	ppu.LY = 0
	ppu.dots = 0
	ppu.internalStateLength = 0
	ppu.internalState = nil

	// Blank screen
	ppu.emptyFrame()
}
