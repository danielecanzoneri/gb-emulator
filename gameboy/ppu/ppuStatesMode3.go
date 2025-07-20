package ppu

type mode3 struct {
	penaltyDots int
}

func (st *mode3) Init(ppu *PPU) {
	ppu.oam.readDisabled = true
	ppu.vRAM.readDisabled = true
	ppu.oam.writeDisabled = true
	ppu.vRAM.writeDisabled = true

	st.penaltyDots = ppu.renderLine() // Penalty dots
	ppu.interruptMode = 3
	ppu.STAT = (ppu.STAT & 0xFC) | drawing
	ppu.checkSTATInterruptState()
}
func (st *mode3) Next(_ *PPU) ppuInternalState {
	return new(mode0)
}
func (st *mode3) Duration() int {
	return 172 + st.penaltyDots
}
