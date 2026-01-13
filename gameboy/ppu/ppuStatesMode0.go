package ppu

type hBlank struct {
	length int
}

func (st *hBlank) Init(ppu *PPU) {
	ppu.oam.readDisabled = false
	ppu.oam.writeDisabled = false
	ppu.vRAM.readDisabled = false
	ppu.vRAM.writeDisabled = false

	// Here we have to reset the previous state length so that each line is 456 dots
	st.length = lineLength - ppu.dots - ppu.internalStateLength

	ppu.interruptMode = 0
	ppu.STAT = (ppu.STAT & 0xFC) | 0
	ppu.checkSTATInterrupt()

	if ppu.HBlankCallback != nil {
		ppu.HBlankCallback()
	}
}
func (st *hBlank) Next(ppu *PPU) ppuInternalState {
	// To mode 1 if LY == 144, to mode 2 otherwise
	ppu.LY++
	ppu.dots -= lineLength

	if ppu.LY == 144 {
		return new(vBlankStart)
	} else {
		return new(oamScan)
	}
}
func (st *hBlank) Duration() int {
	return st.length
}
func (st *hBlank) Name() string { return "HBlank" }
