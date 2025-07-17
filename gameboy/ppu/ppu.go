package ppu

import (
	"github.com/danielecanzoneri/gb-emulator/util"
)

type ppuState int

const (
	line0startingMode0 ppuState = iota // On line 0 after PPU is enabled, mode 2 is replaced by mode 0
	mode0to2           ppuState = iota
	mode2              ppuState = iota
	mode2to3           ppuState = iota
	mode3              ppuState = iota
	mode0              ppuState = iota
	mode1              ppuState = iota
)

const (
	oamScan = 2
	drawing = 3
	hBlank  = 0
	vBlank  = 1
)

type PPU struct {
	Dots int // Dots elapsed rendering this line

	// Internal (machine state and length)
	state       ppuState
	stateLength int // When it reaches 0, switch to next state

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

	// From what I gathered from mooneye tests, OAM and vRAM read behave as follows:
	// normally it is blocked 4 ticks before STAT mode changes
	// (when internal mode flag is updated) but is available again when STAT mode changes
	// (if for example mode 3 lasts 172 dots, vRAM cannot be read for 176 dots, same for OAM).
	// Things are different on the first line 0 after PPU is on. In the first 0 -> 3 mode transition,
	// both are blocked when STAT is updated.
	// Write works differently: they are always tied to the STAT register: if it is 2 or 3,
	// they are blocked for OAM (except for a cycle at the end of mode 2 in the 2 -> 3 transition)
	// and blocked for vRAM in mode 3

	ticksThisMode := min(int(ticks), ppu.stateLength)
	ppu.Dots += ticksThisMode
	ppu.stateLength -= ticksThisMode

	// Change internal mode and update state
	if ppu.stateLength <= 0 {
		switch ppu.state {
		case line0startingMode0, mode2to3:
			ppu.oam.readDisabled = true
			ppu.vRAM.readDisabled = true
			ppu.oam.writeDisabled = true
			ppu.vRAM.writeDisabled = true

			ppu.stateLength = 172 + ppu.renderLine() // Penalty dots
			ppu.state = mode3
			ppu.setMode(drawing)

		case mode0to2:
			ppu.oam.writeDisabled = true

			ppu.stateLength = 76
			ppu.state = mode2
			ppu.setMode(oamScan)
			ppu.searchOAM()

		case mode2:
			ppu.vRAM.readDisabled = true
			ppu.oam.writeDisabled = false

			ppu.stateLength = 4
			ppu.state = mode2to3

		case mode3:
			ppu.oam.readDisabled = false
			ppu.vRAM.readDisabled = false
			ppu.oam.writeDisabled = false
			ppu.vRAM.writeDisabled = false

			ppu.stateLength = 456 - ppu.Dots
			ppu.state = mode0
			ppu.setMode(hBlank)

		case mode0: // To mode 1 if LY == 144, to mode 2 otherwise
			ppu.LY++
			ppu.checkLYLYC()
			ppu.Dots = 0

			if ppu.LY == 144 {
				// Enter VBlank period
				ppu.stateLength = 456
				ppu.state = mode1
				ppu.setMode(vBlank)

				ppu.wyCounter = 0
				ppu.RequestVBlankInterrupt()

				// Frame complete, switch buffers
				ppu.frontBuffer = ppu.backBuffer
				ppu.backBuffer = new([FrameHeight][FrameWidth]uint8)

				// A STAT interrupt can also be triggered at line 144 when vblank starts.
				ppu.checkSTATInterruptState()
			} else {
				ppu.oam.readDisabled = true
				ppu.stateLength = 4
				ppu.state = mode0to2
			}

		case mode1:
			ppu.LY++
			ppu.Dots = 0

			if ppu.LY > 153 {
				ppu.LY = 0
				ppu.oam.readDisabled = true
				ppu.oam.writeDisabled = true
				ppu.stateLength = 80
				ppu.state = mode2
				ppu.setMode(oamScan)
			} else {
				ppu.stateLength = 456
			}

			ppu.checkLYLYC()
		}
	}

	if ticksThisMode < int(ticks) {
		ppu.Tick(ticks - uint(ticksThisMode))
	}
}

func (ppu *PPU) setMode(mode uint8) {
	ppu.STAT = (ppu.STAT & 0xFC) | mode
	ppu.checkSTATInterruptState()
}

func (ppu *PPU) checkLYLYC() {
	// Bit 2 of STAT register is set when LY = LYC
	if ppu.LY == ppu.LYC {
		util.SetBit(&ppu.STAT, 2, 1)
	} else {
		util.SetBit(&ppu.STAT, 2, 0)
	}
	ppu.checkSTATInterruptState()
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
