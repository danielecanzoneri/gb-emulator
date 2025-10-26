package ppu

import (
	"github.com/danielecanzoneri/gb-emulator/util"
)

type PPU struct {
	Dots int // Dots elapsed rendering this line

	// Internal (machine state and length)
	InternalState       ppuInternalState
	InternalStateLength int   // When it reaches 0, switch to next state
	interruptMode       uint8 // Mode used for the STAT interrupt line

	// vRAM and OAM data
	vRAM vRAM
	OAM  OAM

	// objects on current line
	objsLY  [objsLimit]*Object
	numObjs int

	// Double buffering to avoid screen tearing
	frontBuffer *[FrameHeight][FrameWidth]uint16
	backBuffer  *[FrameHeight][FrameWidth]uint16

	LCDC uint8 // LCD control
	STAT uint8 // STAT interrupt
	SCY  uint8 // y - background scrolling
	SCX  uint8 // x - background scrolling
	LY   uint8 // line counter
	LYC  uint8 // line counter check
	BGP  DMGPalette
	OBP  [2]DMGPalette
	WY   uint8
	WX   uint8

	// CGB only registers
	cgb        bool
	BGPI       uint8 // Background palette index
	OBPI       uint8 // Object palette index
	BGPalette  [64]uint8
	OBJPalette [64]uint8

	// Window Y count
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
	STATInterruptLine bool

	// Callbacks to request interrupt
	RequestVBlankInterrupt func()
	RequestSTATInterrupt   func()

	modeTicksElapsed uint
}

func New(cgb bool) *PPU {
	ppu := new(PPU)
	ppu.STAT = 0x84 // Set unused bit (and LY=LYC)
	ppu.OBP[0] = 0xFF
	ppu.OBP[1] = 0xFF

	// Init buffers
	ppu.frontBuffer = new([FrameHeight][FrameWidth]uint16)
	ppu.backBuffer = new([FrameHeight][FrameWidth]uint16)

	ppu.cgb = cgb

	return ppu
}

func (ppu *PPU) Tick(ticks int) {
	if !ppu.active {
		return
	}

	// OAM bug
	if ppu.OAM.buggedRead && ppu.OAM.buggedWrite {
		ppu.triggerOAMBugWriteAndRead()
	} else if ppu.OAM.buggedRead {
		ppu.triggerOAMBugRead()
	} else if ppu.OAM.buggedWrite {
		ppu.triggerOAMBugWrite()
	}
	ppu.OAM.buggedRead = false
	ppu.OAM.buggedWrite = false

	// From what I gathered from mooneye tests, OAM and vRAM read behave as follows:
	// normally it is blocked 4 ticks before STAT mode changes
	// (when internal mode flag is updated) but is available again when STAT mode changes
	// (if for example mode 3 lasts 172 dots, vRAM cannot be read for 176 dots, same for OAM).
	// Things are different on the first line 0 after PPU is on. In the first 0 -> 3 mode transition,
	// both are blocked when STAT is updated.
	// Write works differently: they are always tied to the STAT register: if it is 2 or 3,
	// they are blocked for OAM (except for a cycle at the end of mode 2 in the 2 -> 3 transition)
	// and blocked for vRAM in mode 3
	//
	// When incrementing LY, LY==LYC flag on STAT is set to 0 and then it is updated 4 ticks later

	ppu.InternalStateLength -= ticks
	ppu.Dots += ticks

	// Switch PPU internal state
	for ppu.InternalStateLength <= 0 {
		// Given the PPU state, we compute the next state
		nextState := ppu.InternalState.Next(ppu)
		ppu.setState(nextState)
	}
}

func (ppu *PPU) checkSTATInterrupt() {
	if !ppu.active {
		return
	}
	state := false

	// Bit 2 of STAT register is set when LY = LYC
	if ppu.LY == ppu.LYC {
		util.SetBit(&ppu.STAT, 2, 1)
	} else {
		util.SetBit(&ppu.STAT, 2, 0)
	}

	// LYC == LY int (bit 6)
	if (ppu.LY == ppu.LYC && (util.ReadBit(ppu.STAT, 6) > 0)) ||
		// Mode 2 int (bit 5)
		(ppu.interruptMode == 2 && (util.ReadBit(ppu.STAT, 5) > 0)) ||
		// Mode 1 int (bit 4)
		(ppu.interruptMode == 1 && (util.ReadBit(ppu.STAT, 4) > 0)) ||
		// Mode 0 int (bit 3)
		(ppu.interruptMode == 0 && (util.ReadBit(ppu.STAT, 3) > 0)) {
		state = true
	}

	// Detect low -> high transition
	if !ppu.STATInterruptLine && state {
		ppu.RequestSTATInterrupt()
	}
	ppu.STATInterruptLine = state
}
