package ppu

type mode2 struct{}

func (st *mode2) Init(ppu *PPU) {
	ppu.checkLYLYC()

	ppu.interruptMode = oamScan
	ppu.STAT = (ppu.STAT & 0xFC) | oamScan
	ppu.checkSTATInterruptState()
	ppu.searchOAM()
}
func (st *mode2) Next(_ *PPU) ppuInternalState {
	return new(mode3)
}
func (st *mode2) Duration() int { return mode2Length }
