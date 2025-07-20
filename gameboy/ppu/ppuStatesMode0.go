package ppu

import "github.com/danielecanzoneri/gb-emulator/util"

//
// States of the PPU for the first line after is enabled
//

type enabledLine0 struct{}

func (st *enabledLine0) Init(ppu *PPU) {
	// Line 0 has different timing after enabling, it starts with mode 0 and goes straight to mode 3
	// Moreover, mode 0 is shorter by 8 dots (here we subtract 4 because the other 4 dots will pass when ticking the CPU)
	ppu.Dots += 4
}
func (st *enabledLine0) Next(_ *PPU) ppuInternalState {
	return new(mode3)
}
func (st *enabledLine0) Duration() int { return 78 }

//
// Usual states of the PPU for mode 0
//

// Ticks 0,1,2 mode 0 -> 2 transition
type mode0to2 struct{}

func (st *mode0to2) Init(ppu *PPU) {
	ppu.oam.readDisabled = true
	util.SetBit(&ppu.STAT, 2, 0)

	ppu.interruptMode = hBlank
	ppu.STAT = (ppu.STAT & 0xFC) | hBlank
	ppu.checkSTATInterruptState()
}
func (st *mode0to2) Next(_ *PPU) ppuInternalState { return new(mode2Interrupt) }
func (st *mode0to2) Duration() int                { return 3 }

// Tick 3 mode 0 -> 2 transition
//
// The OAM STAT interrupt occurs 1 T-cycle before STAT actually changes, except on line 0.
type mode2Interrupt struct{}

func (st *mode2Interrupt) Init(ppu *PPU) {
	if ppu.LY != 0 {
		ppu.interruptMode = oamScan
		ppu.checkSTATInterruptState()
	}
}
func (st *mode2Interrupt) Next(_ *PPU) ppuInternalState { return new(mode2) }
func (st *mode2Interrupt) Duration() int                { return 1 }

// HBlank state until end of line
type mode0 struct {
	length int
}

func (st *mode0) Init(ppu *PPU) {
	ppu.oam.readDisabled = false
	ppu.vRAM.readDisabled = false
	ppu.oam.writeDisabled = false
	ppu.vRAM.writeDisabled = false

	// Here we have to reset the previous state length so that each line is 456 dots
	st.length = 456 - ppu.Dots - ppu.InternalStateLength
	ppu.interruptMode = hBlank
	ppu.STAT = (ppu.STAT & 0xFC) | hBlank
	ppu.checkSTATInterruptState()
}
func (st *mode0) Next(ppu *PPU) ppuInternalState {
	// To mode 1 if LY == 144, to mode 2 otherwise
	ppu.LY++
	ppu.Dots -= 456

	if ppu.LY == 144 {
		return new(mode0to1)
	} else {
		return new(mode0to2)
	}
}
func (st *mode0) Duration() int {
	return st.length
}

// Ticks 0-4 mode 0 -> 1 transition for line 144
type mode0to1 struct{}

func (st *mode0to1) Init(ppu *PPU) {
	util.SetBit(&ppu.STAT, 2, 0)
}
func (st *mode0to1) Next(_ *PPU) ppuInternalState {
	return new(line144M1)
}
func (st *mode0to1) Duration() int { return 4 }
