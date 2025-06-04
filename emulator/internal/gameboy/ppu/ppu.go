package ppu

import (
	"github.com/danielecanzoneri/gb-emulator/pkg/util"
)

type PPU struct {
	Mode           uint8 // 2: OAM Scan, 3: Drawing, 0: HBlank, 1: VBlank
	Dots           uint  // 2: 80 dots,  3: 172-289, 0: 87-204, 1: 456 * 10
	mode3ExtraDots uint

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

	// On line 0 after LCD is enabled mode 2 (OAM scan) is replaced by mode 0 (HBlank)
	lcdJustEnabled bool

	// Shared STAT interrupt line (STAT interrupt is triggered after low -> high transition)
	STATInterruptState bool

	// Callbacks to request interrupt
	RequestVBlankInterrupt func()
	RequestSTATInterrupt   func()
}

func (ppu *PPU) Reset() {
	ppu.Mode = 0
	ppu.Dots = 0
	ppu.mode3ExtraDots = 0
	ppu.vRAM.Data = [vRAMSize]uint8{}
	ppu.OAM.Data = [OAMSize]uint8{}
	ppu.objsLY = [objsLimit]*Object{}
	ppu.numObjs = 0
	ppu.Framebuffer = [FrameHeight][FrameWidth]uint8{}
	ppu.LCDC = 0
	ppu.STAT = 0
	ppu.SCY = 0
	ppu.SCX = 0
	ppu.LY = 0
	ppu.LYC = 0
	ppu.BGP = 0
	ppu.OBP0 = 0
	ppu.OBP1 = 0
	ppu.WY = 0
	ppu.WX = 0
	ppu.wyCounter = 0
	ppu.active = false
	ppu.windowTileMapAddr = 0
	ppu.windowEnabled = false
	ppu.bgWindowTileDataArea = 0
	ppu.bgTileMapAddr = 0
	ppu.obj8x16Size = false
	ppu.objEnabled = false
	ppu.bgWindowEnabled = false
	ppu.lcdJustEnabled = false
	ppu.STATInterruptState = false
}

const (
	oamScan = 2
	drawing = 3
	hBlank  = 0
	vBlank  = 1
)

func (ppu *PPU) Cycle() {
	if !ppu.active {
		return
	}

	ppu.Dots += 4 // M-cycles -> T-states

	switch ppu.Mode {
	case oamScan:
		if ppu.Dots >= 80 {
			ppu.setMode(drawing)
		}
	case drawing:
		if ppu.Dots >= 172+80+ppu.mode3ExtraDots {
			ppu.setMode(hBlank)
		}
	case hBlank:
		// On line 0 after LCD is enabled mode 2 (OAM scan) is replaced by mode 0 (HBlank)
		if ppu.lcdJustEnabled && ppu.Dots >= 80 {
			// When enabling LCD line 0 has different timings because the PPU is late by 2 T-cycles
			//ppu.Dots -= 2

			ppu.lcdJustEnabled = false
			ppu.setMode(drawing)
		}

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

			if ppu.LY == 0 {
				ppu.setMode(oamScan)
			}
		}
	}
}

func New() *PPU {
	ppu := new(PPU)
	ppu.STAT = 0x80 // Set unused bit
	ppu.setMode(oamScan)
	ppu.checkLYLYC()

	return ppu
}

func (ppu *PPU) setMode(mode uint8) {
	ppu.Mode = mode
	ppu.STAT = (ppu.STAT & 0xFC) | mode

	ppu.checkSTATInterruptState()

	switch mode {
	case oamScan:
	case drawing:
		ppu.selectObjects()
		ppu.mode3ExtraDots = ppu.drawLine()
	case hBlank:
	case vBlank:
		ppu.wyCounter = 0
		// ppu.FrameComplete = true
		ppu.RequestVBlankInterrupt()
	}
}

func (ppu *PPU) newLine() {
	ppu.LY++
	if ppu.LY > 153 {
		ppu.LY = 0
	}

	ppu.checkLYLYC()
}

func (ppu *PPU) checkLYLYC() {
	if !ppu.active {
		return
	}

	if util.ReadBit(ppu.STAT, 6) > 0 {
		// Bit 2 of STAT register is set when LY = LYC
		if ppu.LY == ppu.LYC {
			util.SetBit(&ppu.STAT, 2, 1)
		} else {
			util.SetBit(&ppu.STAT, 2, 0)
		}
		ppu.checkSTATInterruptState()
	}
}

func (ppu *PPU) checkSTATInterruptState() {
	if !ppu.active {
		return
	}
	state := false

	// LYC == LY int (bit 6)
	if (ppu.LY == ppu.LYC && (util.ReadBit(ppu.STAT, 6) > 0)) ||
		// Mode 2 int (bit 5)
		(ppu.Mode == 2 && (util.ReadBit(ppu.STAT, 5) > 0)) ||
		// Mode 1 int (bit 4)
		(ppu.Mode == 1 && (util.ReadBit(ppu.STAT, 4) > 0)) ||
		// Mode 0 int (bit 3)
		(ppu.Mode == 0 && (util.ReadBit(ppu.STAT, 3) > 0)) {
		state = true
	}

	// Detect low -> high transition
	if !ppu.STATInterruptState && state {
		ppu.RequestSTATInterrupt()
	}
	ppu.STATInterruptState = state
}
