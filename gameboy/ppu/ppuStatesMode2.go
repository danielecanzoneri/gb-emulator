package ppu

import "github.com/danielecanzoneri/gb-emulator/util"

// When PPU is enabled:
//   - line 0 starts with mode 0 and goes straight to mode 3
//   - line 0 has different timings because the PPU is late by 2 T-cycles
type glitchedMode2 struct{}

func (st *glitchedMode2) Init(ppu *PPU) {
	// This line is 8 ticks shorter (4 ticks already passed when enabling PPU)
	ppu.Dots += 4

	ppu.interruptMode = 0xFF
	ppu.STAT = (ppu.STAT & 0xFC) | hBlank
	ppu.checkSTATInterrupt()
}
func (st *glitchedMode2) Next(_ *PPU) ppuInternalState {
	return new(mode3)
}
func (st *glitchedMode2) Duration() int { return mode2Length }

// ------- Normal mode 2 -------

// Mode 2 first 4 ticks
type mode0ToMode2 struct {
}

func (st *mode0ToMode2) Init(ppu *PPU) {
	util.SetBit(&ppu.STAT, 2, 0)

	ppu.OAM.readDisabled = true
}
func (st *mode0ToMode2) Next(_ *PPU) ppuInternalState {
	return new(mode2)
}
func (st *mode0ToMode2) Duration() int { return 4 }

// Normal mode 2
type mode2 struct{}

func (st *mode2) Init(ppu *PPU) {
	ppu.OAM.readDisabled = true
	ppu.OAM.writeDisabled = true

	ppu.interruptMode = oamScan
	ppu.STAT = (ppu.STAT & 0xFC) | oamScan
	ppu.checkSTATInterrupt()
	ppu.searchOAM()
}
func (st *mode2) Next(_ *PPU) ppuInternalState {
	return new(mode2ToMode3)
}
func (st *mode2) Duration() int { return mode2Length - 4 }

type mode2ToMode3 struct{}

func (st *mode2ToMode3) Init(ppu *PPU) {
	ppu.OAM.writeDisabled = false
	ppu.vRAM.readDisabled = true
}
func (st *mode2ToMode3) Next(_ *PPU) ppuInternalState {
	return new(mode3)
}
func (st *mode2ToMode3) Duration() int { return 4 }
