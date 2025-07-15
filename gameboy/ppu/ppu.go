package ppu

import (
	"github.com/danielecanzoneri/gb-emulator/util"
)

const (
	oamScan = 2
	drawing = 3
	hBlank  = 0
	vBlank  = 1
)

type PPU struct {
	Dots           int   // 2: 80 dots,  3: 172-289, 0: 87-204, 1: 456 * 10
	internalMode   uint8 // 2: OAM Scan, 3: Drawing, 0: HBlank, 1: VBlank
	mode3ExtraDots int

	// vRAM and OAM data
	vRAM vRAM
	oam  OAM

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

	modeTicksElapsed uint
}

func New() *PPU {
	ppu := new(PPU)
	ppu.STAT = 0x84 // Set unused bit (and LY=LYC)

	// Init buffers
	ppu.frontBuffer = new([FrameHeight][FrameWidth]uint8)
	ppu.backBuffer = new([FrameHeight][FrameWidth]uint8)

	return ppu
}

func (ppu *PPU) Tick(ticks uint) {
	if !ppu.active {
		return
	}

	ppu.Dots += int(ticks) // M-cycles -> T-states
	ppu.modeTicksElapsed += ticks

	// Delay STAT updates by 4 ticks
	if ppu.modeTicksElapsed >= 4 {
		ppu.STAT = (ppu.STAT & 0xFC) | ppu.internalMode
		switch ppu.internalMode {
		case hBlank:
			ppu.oam.readDisabled = false
			ppu.vRAM.readDisabled = false
			ppu.oam.writeDisabled = false
			ppu.vRAM.writeDisabled = false
		case oamScan:
			ppu.oam.writeDisabled = true
		case drawing:
			ppu.oam.readDisabled = true
			ppu.vRAM.readDisabled = true
			ppu.oam.writeDisabled = true
			ppu.vRAM.writeDisabled = true
		}
		ppu.checkLYLYC()
	}
	ppu.checkSTATInterruptState()

	switch ppu.internalMode {
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
			ppu.setMode(drawing)
		}

		if ppu.Dots > 456 {
			ppu.Dots -= 456
			ppu.newLine()
			if ppu.LY == 144 { // Enter VBlank period
				ppu.setMode(vBlank)

				// A STAT interrupt can also be triggered at line 144 when vblank starts.
				ppu.checkSTATInterruptState()
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

func (ppu *PPU) setMode(mode uint8) {
	// From what I gathered from mooneye tests, OAM and vRAM read behave as follows:
	// normally it is blocked 4 ticks before STAT mode changes
	// (when internal mode flag is updated) but is available again when STAT mode changes
	// (if for example mode 3 lasts 172 dots, vRAM cannot be read for 176 dots, same for OAM).
	// Things are different on the first line 0 after PPU is on. In the first 0 -> 3 mode transition,
	// both are blocked when STAT is updated.
	// Write works differently: they are always tied to the STAT register: if it is 2 or 3,
	// they are blocked for OAM (except for a cycle at the end of mode 2 in the 2 -> 3 transition)
	// and blocked for vRAM in mode 3

	ppu.internalMode = mode
	ppu.modeTicksElapsed = 0

	switch mode {
	case oamScan:
		ppu.oam.readDisabled = true
	case drawing:
		if ppu.lcdJustEnabled {
			ppu.lcdJustEnabled = false
		} else {
			ppu.oam.readDisabled = true
			ppu.vRAM.readDisabled = true
			// For the 2 -> 3 transition
			ppu.oam.writeDisabled = false
		}

		ppu.selectObjects()
		ppu.mode3ExtraDots = ppu.renderLine()
	case hBlank:
		// OAM and vRAM are re-enabled 4 ticks later
		// ppu.oam.readDisabled = false
		// ppu.vRAM.readDisabled = false
	case vBlank:
		ppu.wyCounter = 0
		ppu.RequestVBlankInterrupt()

		// Frame complete, switch buffers
		ppu.frontBuffer = ppu.backBuffer
		ppu.backBuffer = new([FrameHeight][FrameWidth]uint8)
	}
}

func (ppu *PPU) newLine() {
	ppu.LY++
	if ppu.LY > 153 {
		ppu.LY = 0
	}

	// When LY changes, LY=LYC flag is always set to 0 and then updated
	// TODO - check if this is true
	util.SetBit(&ppu.STAT, 2, 0)
}

func (ppu *PPU) checkLYLYC() {
	if !ppu.active {
		return
	}

	// Bit 2 of STAT register is set when LY = LYC
	if ppu.LY == ppu.LYC {
		util.SetBit(&ppu.STAT, 2, 1)
	} else {
		util.SetBit(&ppu.STAT, 2, 0)
	}
}

func (ppu *PPU) checkSTATInterruptState() {
	if !ppu.active {
		return
	}
	state := false

	mode := ppu.STAT & 3
	// LYC == LY int (bit 6)
	if (ppu.LY == ppu.LYC && (util.ReadBit(ppu.STAT, 6) > 0)) ||
		// Mode 2 int (bit 5)
		(mode == 2 && (util.ReadBit(ppu.STAT, 5) > 0)) ||
		// Mode 1 int (bit 4)
		(mode == 1 && (util.ReadBit(ppu.STAT, 4) > 0)) ||
		// Mode 0 int (bit 3)
		(mode == 0 && (util.ReadBit(ppu.STAT, 3) > 0)) {
		state = true
	}

	// If bit 5 (mode 2 OAM interrupt) is set, an interrupt is also triggered
	// at line 144 when vblank starts.
	if ppu.LY == 144 && util.ReadBit(ppu.STAT, 5) > 0 {
		state = true
	}

	// Detect low -> high transition
	if !ppu.STATInterruptState && state {
		ppu.RequestSTATInterrupt()
	}
	ppu.STATInterruptState = state
}
