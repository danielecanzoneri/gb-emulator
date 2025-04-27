package ppu

import (
	"github.com/danielecanzoneri/gb-emulator/internal/util"
)

type PPU struct {
	Mode uint8 // 2: OAM Scan, 3: Drawing, 0: HBlank, 1: VBlank
	Dots uint  // 2: 80 dots,  3: 172-289, 0: 87-204, 1: 456 * 10

	// vRAM and OAM data
	vRAM vRAM
	OAM  OAM

	// objects on current line
	objsLY  [objsLimit]*Object
	numObjs int

	// frameBuffer contains data to be displayed
	Framebuffer [FrameHeight][FrameWidth]uint8

	LCDC uint8 // LCD control
	STAT uint8 // STAT interrupt
	SCY  uint8 // y - background scrolling
	SCX  uint8 // x - background scrolling
	LY   uint8 // line counter
	LYC  uint8 // line counter check
	BGP  Palette
	OBP0 Palette
	OBP1 Palette
	WY   uint8
	WX   uint8

	// Window Y counter
	wyCounter uint8

	// LCD control
	active               bool   // Bit 7
	windowTileMapAddr    uint16 // Bit 6 (0 = 9800–9BFF; 1 = 9C00–9FFF)
	windowEnabled        bool   // Bit 5
	bgWindowTileDataArea uint8  // Bit 4 (0 = 8800–97FF; 1 = 8000–8FFF)
	bgTileMapAddr        uint16 // Bit 3 (0 = 9800–9BFF; 1 = 9C00–9FFF)
	obj8x16Size          bool   // Bit 2 (false = 8x8; true = 8x16)
	objEnabled           bool   // Bit 1
	bgWindowEnabled      bool   // Bit 0

	// Shared STAT interrupt line (STAT interrupt is triggered after low -> high transition)
	STATInterruptState bool

	// Callbacks to request interrupt
	RequestVBlankInterrupt func()
	RequestSTATInterrupt   func()

	// FrameComplete signal when frame is ready to be rendered
	FrameComplete bool
}

const (
	oamScan = 2
	drawing = 3
	hBlank  = 0
	vBlank  = 1
)

func (ppu *PPU) Step(cycles uint) {
	if !ppu.active {
		return
	}

	ppu.Dots += cycles * 4 // M-cycles -> T-states

	switch ppu.Mode {
	case oamScan:
		if ppu.Dots >= 80 {
			ppu.setMode(drawing)
		}
	case drawing:
		if ppu.Dots >= 172+80 {
			// TODO - Correctly compute Mode 3 length
			ppu.setMode(hBlank)
		}
	case hBlank:
		if ppu.Dots >= 456 {
			ppu.Dots -= 456
			ppu.newLine()
			if ppu.LY == 144 { // Enter VBlank period
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
	ppu.Mode = mode
	ppu.STAT = (ppu.STAT & 0xFC) | mode

	ppu.checkSTATInterruptState()

	switch mode {
	case oamScan:
	case drawing:
		ppu.selectObjects()
	case hBlank:
		ppu.drawLine()
	case vBlank:
		ppu.wyCounter = 0
		ppu.FrameComplete = true
		ppu.RequestVBlankInterrupt()
	}
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
