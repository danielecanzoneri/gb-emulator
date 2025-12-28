package ppu

// OAM scan to Drawing transition
type oamScanToDrawing struct{}

func (st *oamScanToDrawing) Init(ppu *PPU) {
	ppu.oam.writeDisabled = false
	ppu.vRAM.readDisabled = true
}
func (st *oamScanToDrawing) Next(_ *PPU) ppuInternalState {
	return new(drawing)
}
func (st *oamScanToDrawing) Duration() int { return 4 }

type drawing struct {
	penaltyDots int
}

func (st *drawing) Init(ppu *PPU) {
	ppu.oam.readDisabled = true
	ppu.oam.writeDisabled = true
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
