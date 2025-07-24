package ppu

type mode3 struct {
	penaltyDots int
}

func (st *mode3) Init(ppu *PPU) {
	ppu.oam.readDisabled = true
	ppu.oam.writeDisabled = true
	ppu.vRAM.readDisabled = true
	ppu.vRAM.writeDisabled = true

	st.penaltyDots = ppu.renderLine() // Penalty dots
	ppu.interruptMode = drawing
	ppu.STAT = (ppu.STAT & 0xFC) | drawing
	ppu.checkSTATInterrupt()
}
func (st *mode3) Next(_ *PPU) ppuInternalState {
	return new(mode0)
}
func (st *mode3) Duration() int {
	// Round to M-cycle (TODO - investigate why it doesn't work otherwise)
	return 172 + st.penaltyDots & ^3
}
