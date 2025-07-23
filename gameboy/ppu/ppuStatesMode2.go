package ppu

// When PPU is enabled, mode 2 of the first line is replaced with mode 0
type glitchedMode2 struct{}

func (st *glitchedMode2) Init(ppu *PPU) {
	ppu.interruptMode = 0xFF
	ppu.STAT = (ppu.STAT & 0xFC) | hBlank
	ppu.checkSTATInterrupt()
}
func (st *glitchedMode2) Next(_ *PPU) ppuInternalState {
	return new(mode3)
}
func (st *glitchedMode2) Duration() int { return mode2Length }

// Normal mode 2
type mode2 struct{}

func (st *mode2) Init(ppu *PPU) {
	ppu.oam.readDisabled = true
	ppu.oam.writeDisabled = true

	ppu.interruptMode = oamScan
	ppu.STAT = (ppu.STAT & 0xFC) | oamScan
	ppu.checkSTATInterrupt()
	ppu.searchOAM()
}
func (st *mode2) Next(_ *PPU) ppuInternalState {
	return new(mode3)
}
func (st *mode2) Duration() int { return mode2Length }
