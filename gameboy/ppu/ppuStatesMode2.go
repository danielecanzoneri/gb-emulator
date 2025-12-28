package ppu

import "github.com/danielecanzoneri/gb-emulator/util"

// When PPU is enabled:
//   - line 0 starts with mode 0 and goes straight to mode 3
//   - line 0 has different timings because the PPU is late by 2 T-cycles
type glitchedOamScan struct{}

func (st *glitchedOamScan) Init(ppu *PPU) {
	// This line is 8 ticks shorter (4 ticks already passed when enabling PPU)
	ppu.dots += 4

	ppu.interruptMode = 0xFF
	ppu.STAT = (ppu.STAT & 0xFC) | 0
	ppu.checkSTATInterrupt()
}
func (st *glitchedOamScan) Next(_ *PPU) ppuInternalState {
	return new(drawing)
}
func (st *glitchedOamScan) Duration() int { return mode2Length }

// ------- Normal mode 2 -------

// Normal mode 2
type oamScan struct {
	// OAM is divided in 20 rows of 8 bytes, every M-cycle a different row is read
	rowAccessed uint8
}

func (st *oamScan) Init(ppu *PPU) {
	if st.rowAccessed == 0 { // Still hBlank
		util.SetBit(&ppu.STAT, 2, 0)
		ppu.oam.readDisabled = true

	} else if st.rowAccessed == 1 {
		ppu.oam.readDisabled = true
		ppu.oam.writeDisabled = true

		ppu.interruptMode = 2
		ppu.STAT = (ppu.STAT & 0xFC) | 2
		ppu.checkSTATInterrupt()
		ppu.searchOAM()
	}
}
func (st *oamScan) Next(_ *PPU) ppuInternalState {
	// If this was the last row, next mode will transition to drawing state
	if st.rowAccessed == 19 {
		return new(oamScanToDrawing)
	}

	st.rowAccessed++
	return st
}
func (st *oamScan) Duration() int { return 4 }
