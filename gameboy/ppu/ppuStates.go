package ppu

const (
	lineLength = 456

	mode2Length = 80
	mode3Length = 172
)

type ppuInternalState interface {
	Init(*PPU)
	Next(*PPU) ppuInternalState
	Duration() int
}

func (ppu *PPU) setState(state ppuInternalState) {
	ppu.internalState = state
	state.Init(ppu) // Here it's where the actual state switching happens
	ppu.internalStateLength += state.Duration()
}
