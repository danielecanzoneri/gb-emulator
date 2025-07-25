package ppu

type mode0 struct {
	length int
}

func (st *mode0) Init(ppu *PPU) {
	ppu.OAM.readDisabled = false
	ppu.OAM.writeDisabled = false
	ppu.vRAM.readDisabled = false
	ppu.vRAM.writeDisabled = false

	// Here we have to reset the previous state length so that each line is 456 dots
	st.length = lineLength - ppu.Dots - ppu.InternalStateLength

	ppu.interruptMode = hBlank
	ppu.STAT = (ppu.STAT & 0xFC) | hBlank
	ppu.checkSTATInterrupt()
}
func (st *mode0) Next(ppu *PPU) ppuInternalState {
	// To mode 1 if LY == 144, to mode 2 otherwise
	ppu.LY++
	ppu.Dots -= lineLength

	if ppu.LY == 144 {
		return new(mode1Start)
	} else {
		return new(mode0ToMode2)
	}
}
func (st *mode0) Duration() int {
	return st.length
}
