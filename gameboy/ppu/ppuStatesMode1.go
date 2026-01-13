package ppu

import "github.com/danielecanzoneri/lucky-boy/util"

// First 4 ticks
type vBlankStart struct{}

func (st *vBlankStart) Init(ppu *PPU) {
	util.SetBit(&ppu.STAT, 2, 0)
}
func (st *vBlankStart) Next(_ *PPU) ppuInternalState {
	return new(vBlank)
}
func (st *vBlankStart) Duration() int { return 4 }
func (st *vBlankStart) Name() string  { return "VBlank Start" }

// Remaining ticks for vBlank state
type vBlank struct{}

func (st *vBlank) Init(ppu *PPU) {
	if ppu.LY == 144 {
		// If bit 5 (mode 2 OAM interrupt) is set, an interrupt is also triggered
		// at line 144 when vblank starts.
		ppu.interruptMode = 2
		ppu.checkSTATInterrupt()

		ppu.interruptMode = 1
		ppu.STAT = (ppu.STAT & 0xFC) | 1
		ppu.checkSTATInterrupt()

		ppu.wyCounter = 0
		ppu.RequestVBlankInterrupt()

		// Frame complete, switch buffers
		ppu.swapBuffers()

		if ppu.VBlankCallback != nil {
			ppu.VBlankCallback()
		}
	}
}
func (st *vBlank) Next(ppu *PPU) ppuInternalState {
	ppu.LY++
	ppu.dots -= lineLength

	if ppu.LY == 154 {
		ppu.LY = 0
		return new(oamScan)
	} else {
		return new(vBlankStart)
	}
}
func (st *vBlank) Duration() int { return lineLength - 4 }
func (st *vBlank) Name() string  { return "VBlank" }
