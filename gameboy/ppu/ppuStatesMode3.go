package ppu

type mode3 struct {
	penaltyDots int
}

func (st *mode3) Init(ppu *PPU) {
	ppu.OAM.readDisabled = true
	ppu.OAM.writeDisabled = true
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
	return 172 + st.penaltyDots
}
