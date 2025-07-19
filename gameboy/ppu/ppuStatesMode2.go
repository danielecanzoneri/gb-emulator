package ppu

type mode2 struct{}

func (st *mode2) Init(ppu *PPU) {
	ppu.checkLYLYC()
	ppu.oam.writeDisabled = true

	ppu.setMode(oamScan)
	ppu.searchOAM()
}
func (st *mode2) Next(_ *PPU) ppuInternalState {
	return new(mode2to3)
}
func (st *mode2) Duration() int { return 76 }

type mode2to3 struct{}

func (st *mode2to3) Init(ppu *PPU) {
	ppu.vRAM.readDisabled = true
	ppu.oam.writeDisabled = false
}
func (st *mode2to3) Next(_ *PPU) ppuInternalState {
	return new(mode3)
}
func (st *mode2to3) Duration() int { return 4 }
