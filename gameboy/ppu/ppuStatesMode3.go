package ppu

type mode3 struct {
	penaltyDots int
}

func (st *mode3) Init(ppu *PPU) {
	st.penaltyDots = ppu.renderLine() // Penalty dots
	ppu.interruptMode = drawing
	ppu.STAT = (ppu.STAT & 0xFC) | drawing
	ppu.checkSTATInterruptState()
}
func (st *mode3) Next(_ *PPU) ppuInternalState {
	return new(mode0)
}
func (st *mode3) Duration() int {
	return 172 + st.penaltyDots
}
