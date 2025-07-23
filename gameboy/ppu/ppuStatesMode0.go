package ppu

type mode0 struct {
	length int
}

func (st *mode0) Init(ppu *PPU) {
	// Here we have to reset the previous state length so that each line is 456 dots
	st.length = lineLength - ppu.Dots - ppu.InternalStateLength

	ppu.interruptMode = hBlank
	ppu.STAT = (ppu.STAT & 0xFC) | hBlank
	ppu.checkSTATInterruptState()
}
func (st *mode0) Next(ppu *PPU) ppuInternalState {
	// To mode 1 if LY == 144, to mode 2 otherwise
	ppu.LY++
	ppu.Dots -= lineLength

	if ppu.LY == 144 {
		return new(mode1)
	} else {
		return new(mode2)
	}
}
func (st *mode0) Duration() int {
	return st.length
}
