package ppu

type ppuInternalState interface {
	Init(*PPU)
	Next(*PPU) ppuInternalState
	Duration() int
}

func (ppu *PPU) setState(state ppuInternalState) {
	ppu.internalState = state
	state.Init(ppu)
}
