package ppu

import "github.com/danielecanzoneri/gb-emulator/internal/util"

type PPU struct {
	Mode uint8 // 2: OAM Scan, 3: Drawing, 0: HBlank, 1: VBlank
	Dots uint  // 2: 80 dots,  3: 172-289, 0: 87-204, 1: 456 * 10

	LCDC uint8
	STAT uint8
	SCY  uint8
	SCX  uint8
	LY   uint8
	LYC  uint8
	DMA  uint8
	BGP  uint8
	OBP0 uint8
	OBP1 uint8
	WY   uint8
	WX   uint8

	// LCD control
	active               bool  // Bit 7
	windowTileMapArea    uint8 // Bit 6 (0 = 9800–9BFF; 1 = 9C00–9FFF)
	windowEnabled        bool  // Bit 5
	bgWindowTileDataArea uint8 // Bit 4 (8800–97FF; 1 = 8000–8FFF)
	bgTileMapArea        uint8 // Bit 3 (0 = 9800–9BFF; 1 = 9C00–9FFF)
	objSize              uint8 // Bit 2 (0 = 8x8; 1 = 8x16)
	objEnabled           bool  // Bit 1
	bgWindowEnabled      bool  // Bit 0

	// Shared STAT interrupt line (STAT interrupt is triggered after low -> high transition)
	STATInterruptState bool

	// Callbacks to request interrupt
	RequestVBlankInterrupt func()
	RequestSTATInterrupt   func()
}

const (
	oamScan = 2
	drawing = 3
	hBlank  = 1
	vBlank  = 0
)

func (ppu *PPU) Step(cycles uint) {
	ppu.Dots += cycles * 4 // M-cycles -> T-states

	switch ppu.Mode {
	case oamScan:
		if ppu.Dots >= 80 {
			ppu.Dots -= 80
			ppu.setMode(drawing)
		}
	case drawing: // TODO
	case hBlank:
		if ppu.Dots >= 456 {
			ppu.Dots -= 456
			ppu.newLine()
			if ppu.LY == 144 { // Enter VBlank period
				ppu.RequestVBlankInterrupt()
				ppu.setMode(vBlank)
			} else {
				ppu.setMode(oamScan)
			}
		}
	case vBlank:
		if ppu.Dots >= 456 {
			ppu.Dots -= 456
			ppu.newLine()

			if ppu.LY > 153 {
				ppu.LY = 0
				ppu.setMode(oamScan)
			}
		}
	}
}

func (ppu *PPU) setMode(mode uint8) {
	ppu.checkSTATInterruptState()
	ppu.Mode = mode
	ppu.STAT = (ppu.STAT & 0xFC) | mode
}

func (ppu *PPU) newLine() {
	ppu.LY++

	// Bit 2 of STAT register is set when LY = LYC
	if ppu.LY == ppu.LYC {
		util.SetBit(&ppu.STAT, 2, 1)
	} else {
		util.SetBit(&ppu.STAT, 2, 0)
	}
	ppu.checkSTATInterruptState()
}

func (ppu *PPU) checkSTATInterruptState() {
	state := false

	// LYC == LY int (bit 6)
	if (ppu.LY == ppu.LYC && (util.ReadBit(ppu.STAT, 6) > 0)) ||
		// Mode 2 int (bit 5)
		(ppu.Mode == 2 && (util.ReadBit(ppu.STAT, 5) > 0)) ||
		// Mode 1 int (bit 4)
		(ppu.Mode == 1 && (util.ReadBit(ppu.STAT, 4) > 0)) ||
		// Mode 0 int (bit 3)
		(ppu.Mode == 2 && (util.ReadBit(ppu.STAT, 3) > 0)) {
		state = true
	}

	// Detect low -> high transition
	if !ppu.STATInterruptState && state {
		ppu.RequestSTATInterrupt()
	}
	ppu.STATInterruptState = state
}
