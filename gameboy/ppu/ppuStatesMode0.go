package ppu

import "github.com/danielecanzoneri/gb-emulator/util"

//
// States of the PPU for the first line after is enabled
//

type line0startingMode0 struct{}

func (st *line0startingMode0) Init(ppu *PPU) {
	// Line 0 has different timing after enabling, it starts with mode 0 and goes straight to mode 3
	// Moreover, mode 0 is shorter by 8 dots (here we subtract 4 because the other 4 dots will pass when ticking the CPU)
	ppu.Dots += 4
}
func (st *line0startingMode0) Next(_ *PPU) ppuInternalState {
	return new(mode3)
}
func (st *line0startingMode0) Duration() int { return 78 }

//
// Usual states of the PPU for mode 0
//

// Ticks 0-4 mode 0 -> 2 transition
type mode0to2 struct{}

func (st *mode0to2) Init(ppu *PPU) {
	ppu.oam.readDisabled = true
	ppu.oam.writeDisabled = true

	ppu.setMode(hBlank)
}
func (st *mode0to2) Next(_ *PPU) ppuInternalState {
	return new(mode2)
}
func (st *mode0to2) Duration() int { return 4 }

// HBlank state until end of line
type mode0 struct {
	length int
}

func (st *mode0) Init(ppu *PPU) {
	ppu.oam.readDisabled = false
	ppu.vRAM.readDisabled = false
	ppu.oam.writeDisabled = false
	ppu.vRAM.writeDisabled = false

	st.length = 456 - ppu.Dots
	ppu.setMode(hBlank)
}
func (st *mode0) Next(ppu *PPU) ppuInternalState {
	// To mode 1 if LY == 144, to mode 2 otherwise
	ppu.LY++
	util.SetBit(&ppu.STAT, 2, 0)
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

func (st *mode0to1) Init(_ *PPU) {
}
func (st *mode0to1) Next(ppu *PPU) ppuInternalState {
	ppu.checkLYLYC()
	ppu.setMode(vBlank)

	ppu.wyCounter = 0
	ppu.RequestVBlankInterrupt()

	// Frame complete, switch buffers
	ppu.frontBuffer = ppu.backBuffer
	ppu.backBuffer = new([FrameHeight][FrameWidth]uint8)

	// A STAT interrupt can also be triggered at line 144 when vblank starts.
	ppu.checkSTATInterruptState()

	return new(mode1)
}
func (st *mode0to1) Duration() int { return 4 }
