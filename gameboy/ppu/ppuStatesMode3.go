package ppu

type drawing struct {
	penaltyDots int
}

func (st *drawing) Init(ppu *PPU) {
	ppu.OAM.readDisabled = true
	ppu.OAM.writeDisabled = true
	ppu.vRAM.readDisabled = true
	ppu.vRAM.writeDisabled = true

	st.penaltyDots = ppu.renderLine() // Penalty dots
	ppu.interruptMode = 3
	ppu.STAT = (ppu.STAT & 0xFC) | 3
	ppu.checkSTATInterrupt()
}
func (st *drawing) Next(_ *PPU) ppuInternalState {
	return new(hBlank)
}
func (st *drawing) Duration() int {
	return 172 + st.penaltyDots
}
