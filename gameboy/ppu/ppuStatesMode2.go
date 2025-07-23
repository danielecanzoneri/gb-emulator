package ppu

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

// Normal mode 2
type mode2 struct{}

func (st *mode2) Init(ppu *PPU) {
	ppu.oam.readDisabled = true
	ppu.oam.writeDisabled = true

	ppu.interruptMode = oamScan
	ppu.STAT = (ppu.STAT & 0xFC) | oamScan
	ppu.checkSTATInterrupt()
	ppu.searchOAM()
}
func (st *mode2) Next(_ *PPU) ppuInternalState {
	return new(mode3)
}
func (st *mode2) Duration() int { return mode2Length }
