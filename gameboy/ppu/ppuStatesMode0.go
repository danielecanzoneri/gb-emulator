package ppu

type hBlank struct {
	length int
}

func (st *hBlank) Init(ppu *PPU) {
	ppu.OAM.readDisabled = false
	ppu.OAM.writeDisabled = false
	ppu.vRAM.readDisabled = false
	ppu.vRAM.writeDisabled = false

	// Here we have to reset the previous state length so that each line is 456 dots
	st.length = lineLength - ppu.Dots - ppu.InternalStateLength

	ppu.interruptMode = 0
	ppu.STAT = (ppu.STAT & 0xFC) | 0
	ppu.checkSTATInterrupt()
}
func (st *hBlank) Next(ppu *PPU) ppuInternalState {
	// To mode 1 if LY == 144, to mode 2 otherwise
	ppu.LY++
	ppu.Dots -= lineLength

	if ppu.LY == 144 {
		return new(vBlankStart)
	} else {
		return new(oamScan)
	}
}
func (st *hBlank) Duration() int {
	return st.length
}
