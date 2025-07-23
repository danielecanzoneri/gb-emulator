package ppu

type mode1 struct{}

func (st *mode1) Init(ppu *PPU) {
	ppu.checkLYLYC()

	if ppu.LY == 144 {
		// If bit 5 (mode 2 OAM interrupt) is set, an interrupt is also triggered
		// at line 144 when vblank starts.
		ppu.interruptMode = oamScan
		ppu.checkSTATInterruptState()

		ppu.interruptMode = vBlank
		ppu.STAT = (ppu.STAT & 0xFC) | vBlank
		ppu.checkSTATInterruptState()

		ppu.wyCounter = 0
		ppu.RequestVBlankInterrupt()

		// Frame complete, switch buffers
		ppu.frontBuffer = ppu.backBuffer
		ppu.backBuffer = new([FrameHeight][FrameWidth]uint8)
	}
}
func (st *mode1) Next(ppu *PPU) ppuInternalState {
	ppu.LY++
	ppu.Dots -= lineLength

	if ppu.LY == 154 {
		ppu.LY = 0
		return new(mode2)
	} else {
		return new(mode1)
	}
}
func (st *mode1) Duration() int { return lineLength }
