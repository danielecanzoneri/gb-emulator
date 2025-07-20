package ppu

import "github.com/danielecanzoneri/gb-emulator/util"

//
// VBlank mode for lines 144-152
//

// Ticks 4-456 for line 144
type line144M1 struct{}

func (st *line144M1) Init(ppu *PPU) {
	ppu.checkLYLYC()

	// If bit 5 (mode 2 OAM interrupt) is set, an interrupt is also triggered
	// at line 144 when vblank starts.
	ppu.interruptMode = oamScan
	ppu.checkSTATInterruptState()

	ppu.interruptMode = vBlank
	ppu.STAT = (ppu.STAT & 0xFC) | vBlank
	ppu.checkSTATInterruptState()

	ppu.wyCounter = 0
	ppu.RequestVBlankInterrupt()

	// Frame complete, switch buffers
	ppu.frontBuffer = ppu.backBuffer
	ppu.backBuffer = new([FrameHeight][FrameWidth]uint8)

	// A STAT interrupt can also be triggered at line 144 when vblank starts.
	ppu.checkSTATInterruptState()
}
func (st *line144M1) Next(ppu *PPU) ppuInternalState {
	ppu.LY++
	util.SetBit(&ppu.STAT, 2, 0)
	ppu.Dots -= 456

	return new(startM1)
}
func (st *line144M1) Duration() int { return 456 - 4 }

// Ticks 0-4 for lines 145-152
type startM1 struct{}

func (st *startM1) Init(ppu *PPU) {
	util.SetBit(&ppu.STAT, 2, 0)
}
func (st *startM1) Next(_ *PPU) ppuInternalState {
	return new(mode1)
}
func (st *startM1) Duration() int { return 4 }

// Ticks 4-456 for lines 145-152
type mode1 struct{}

func (st *mode1) Init(ppu *PPU) {
	ppu.checkLYLYC()
}
func (st *mode1) Next(ppu *PPU) ppuInternalState {
	ppu.LY++
	ppu.Dots -= 456

	if ppu.LY == 153 {
		return new(startLine153)
	} else {
		return new(startM1)
	}
}
func (st *mode1) Duration() int { return 456 - 4 }

//
// VBlank mode for line 153
//

// Ticks 0-4 for line 153
type startLine153 struct{}

func (st *startLine153) Init(ppu *PPU) {
	util.SetBit(&ppu.STAT, 2, 0)
}
func (st *startLine153) Next(_ *PPU) ppuInternalState {
	return new(line153LY0)
}
func (st *startLine153) Duration() int { return 4 }

// Ticks 4-8 for line 153
type line153LY0 struct{}

func (st *line153LY0) Init(ppu *PPU) {
	ppu.LY = 0
	ppu.checkLYLYC()
}
func (st *line153LY0) Next(_ *PPU) ppuInternalState {
	return new(line153LYLYCUnset)
}
func (st *line153LY0) Duration() int { return 4 }

// Ticks 8-12 for line 153
type line153LYLYCUnset struct{}

func (st *line153LYLYCUnset) Init(ppu *PPU) {
	util.SetBit(&ppu.STAT, 2, 0)
}
func (st *line153LYLYCUnset) Next(_ *PPU) ppuInternalState {
	return new(line153LYLYC0)
}
func (st *line153LYLYCUnset) Duration() int { return 4 }

// Ticks 12-456 for line 153
type line153LYLYC0 struct{}

func (st *line153LYLYC0) Init(ppu *PPU) {
	ppu.checkLYLYC()
}
func (st *line153LYLYC0) Next(ppu *PPU) ppuInternalState {
	ppu.Dots -= 456
	return new(mode0to2)
}
func (st *line153LYLYC0) Duration() int { return 456 - 12 }
