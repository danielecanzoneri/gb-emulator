package ppu

import "github.com/danielecanzoneri/gb-emulator/util"

// First 4 ticks
type mode1Start struct{}

func (st *mode1Start) Init(ppu *PPU) {
	util.SetBit(&ppu.STAT, 2, 0)
}
func (st *mode1Start) Next(ppu *PPU) ppuInternalState {
	return new(mode1)
}
func (st *mode1Start) Duration() int { return 4 }

// Remaining ticks for vBlank state
type mode1 struct{}

func (st *mode1) Init(ppu *PPU) {
	if ppu.LY == 144 {
		// If bit 5 (mode 2 OAM interrupt) is set, an interrupt is also triggered
		// at line 144 when vblank starts.
		ppu.interruptMode = oamScan
		ppu.checkSTATInterrupt()

		ppu.interruptMode = vBlank
		ppu.STAT = (ppu.STAT & 0xFC) | vBlank
		ppu.checkSTATInterrupt()

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
		return new(mode0ToMode2)
	} else {
		return new(mode1Start)
	}
}
func (st *mode1) Duration() int { return lineLength - 4 }
