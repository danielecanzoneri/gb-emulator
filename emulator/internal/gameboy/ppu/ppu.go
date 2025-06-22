package ppu

import (
	"github.com/danielecanzoneri/gb-emulator/pkg/util"
)

type PPU struct {
	Mode           uint8 // 2: OAM Scan, 3: Drawing, 0: HBlank, 1: VBlank
	Dots           int   // 2: 80 dots,  3: 172-289, 0: 87-204, 1: 456 * 10
	mode3ExtraDots int

	// vRAM and OAM data
	vRAM vRAM
	OAM  OAM

	// objects on current line
	objsLY  [objsLimit]*Object
	numObjs int

	// Double buffering to avoid screen tearing
	frontBuffer *[FrameHeight][FrameWidth]uint8
	backBuffer  *[FrameHeight][FrameWidth]uint8

	LCDC uint8 // LCD control
	STAT uint8 // STAT interrupt
	SCY  uint8 // y - background scrolling
	SCX  uint8 // x - background scrolling
	LY   uint8 // line counter
	LYC  uint8 // line counter check
	BGP  Palette
	OBP  [2]Palette
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

	// Delay STAT updates by 4 ticks
	delayedSTATUpdates []delayedSTATUpdate
}

type delayedSTATUpdate struct {
	delay  int
	update func()
}

func (ppu *PPU) Reset() {
	ppu.Mode = 0
	ppu.Dots = 0
	ppu.mode3ExtraDots = 0
	ppu.objsLY = [objsLimit]*Object{}
	ppu.numObjs = 0
	ppu.frontBuffer = &[FrameHeight][FrameWidth]uint8{}
	ppu.backBuffer = &[FrameHeight][FrameWidth]uint8{}
	ppu.LCDC = 0
	ppu.STAT = 0x80 // Unused bit
	ppu.SCY = 0
	ppu.SCX = 0
	ppu.LY = 0
	ppu.LYC = 0
	ppu.BGP = 0
	ppu.OBP[0] = 0
	ppu.OBP[1] = 0
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

func (ppu *PPU) Tick(ticks uint) {
	if !ppu.active {
		return
	}

	// Handle delayed STAT updates
	var remaining []delayedSTATUpdate
	for _, d := range ppu.delayedSTATUpdates {
		d.delay -= int(ticks)
		if d.delay <= 0 {
			d.update()
			ppu.checkSTATInterruptState()
		} else {
			remaining = append(remaining, d)
		}
	}
	ppu.delayedSTATUpdates = remaining

	ppu.Dots += int(ticks) // M-cycles -> T-states

	switch ppu.Mode {
	case oamScan:
		if ppu.Dots > 80 {
			ppu.setMode(drawing)
		}
	case drawing:
		if ppu.Dots > 172+80+ppu.mode3ExtraDots {
			ppu.setMode(hBlank)
		}
	case hBlank:
		// On line 0 after LCD is enabled mode 2 (OAM scan) is replaced by mode 0 (HBlank)
		if ppu.lcdJustEnabled && ppu.Dots > 80 {
			ppu.lcdJustEnabled = false
			ppu.setMode(drawing)
		}

		if ppu.Dots > 456 {
			ppu.Dots -= 456
			ppu.newLine()
			if ppu.LY == 144 { // Enter VBlank period
				ppu.setMode(vBlank)
			} else {
				ppu.setMode(oamScan)
			}
		}
	case vBlank:
		if ppu.Dots > 456 {
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

	// Init buffers
	ppu.frontBuffer = new([FrameHeight][FrameWidth]uint8)
	ppu.backBuffer = new([FrameHeight][FrameWidth]uint8)

	return ppu
}

func (ppu *PPU) setMode(mode uint8) {
	ppu.Mode = mode

	switch mode {
	case oamScan:
	case drawing:
		ppu.selectObjects()
		ppu.mode3ExtraDots = ppu.drawLine()
	case hBlank:
	case vBlank:
		ppu.wyCounter = 0
		ppu.RequestVBlankInterrupt()

		// Frame complete, switch buffers
		ppu.frontBuffer = ppu.backBuffer
		ppu.backBuffer = new([FrameHeight][FrameWidth]uint8)
	}

	ppu.setSTATUpdateHandler(func() {
		ppu.STAT = (ppu.STAT & 0xFC) | mode
	})
}

func (ppu *PPU) newLine() {
	ppu.LY++
	if ppu.LY > 153 {
		ppu.LY = 0
	}

	// When LY changes, LY=LYC flag is always set to 0 and then updated
	// TODO - check if this is true
	util.SetBit(&ppu.STAT, 2, 0)

	ppu.checkLYLYC()
}

func (ppu *PPU) setSTATUpdateHandler(statUpdate func()) {
	// STAT register is updated 4 (or 1 is enough?) T-cycles later
	// (see https://github.com/Gekkio/mooneye-test-suite/blob/main/acceptance/ppu/lcdon_timing-GS.s)

	ppu.delayedSTATUpdates = append(ppu.delayedSTATUpdates, delayedSTATUpdate{
		delay:  4,
		update: statUpdate,
	})
}

func (ppu *PPU) checkLYLYC() {
	if !ppu.active {
		return
	}

	// Bit 2 of STAT register is set when LY = LYC
	if ppu.LY == ppu.LYC {
		ppu.setSTATUpdateHandler(func() {
			util.SetBit(&ppu.STAT, 2, 1)
		})
	} else {
		ppu.setSTATUpdateHandler(func() {
			util.SetBit(&ppu.STAT, 2, 0)
		})
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
