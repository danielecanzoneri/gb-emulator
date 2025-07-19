package ppu

import "github.com/danielecanzoneri/gb-emulator/util"

//
// VBlank mode for lines 144-152
//

// Ticks 4-456
type mode1 struct{}

func (st *mode1) Init(_ *PPU) {}
func (st *mode1) Next(ppu *PPU) ppuInternalState {
	ppu.LY++
	util.SetBit(&ppu.STAT, 2, 0)
	ppu.Dots -= 456

	if ppu.LY == 153 {
		return new(line153Start)
	} else {
		return new(startingMode1)
	}
}
func (st *mode1) Duration() int { return 456 - 4 }

// Ticks 0-4 for lines 145-152
type startingMode1 struct{}

func (st *startingMode1) Init(_ *PPU) {}
func (st *startingMode1) Next(ppu *PPU) ppuInternalState {
	ppu.checkLYLYC()
	return new(mode1)
}
func (st *startingMode1) Duration() int { return 4 }

//
// VBlank mode for line 153
//

// Ticks 0-4 for line 153
type line153Start struct{}

func (st *line153Start) Init(_ *PPU) {}
func (st *line153Start) Next(_ *PPU) ppuInternalState {
	return new(line153LY0)
}
func (st *line153Start) Duration() int { return 4 }

// Ticks 4-8 for line 153
type line153LY0 struct{}

func (st *line153LY0) Init(ppu *PPU) {
	ppu.checkLYLYC()

	ppu.LY = 0
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
func (st *line153LYLYC0) Next(_ *PPU) ppuInternalState {
	return new(mode0to2)
}
func (st *line153LYLYC0) Duration() int { return 456 - 12 }
