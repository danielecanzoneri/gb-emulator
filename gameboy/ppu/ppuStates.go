package ppu

type ppuInternalState interface {
	Init(*PPU)
	Next(*PPU) ppuInternalState
	Duration() int
}

func (ppu *PPU) setState(state ppuInternalState) {
	ppu.InternalState = state
	state.Init(ppu)
	ppu.InternalStateLength += state.Duration()
}
