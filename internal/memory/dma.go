package memory

// DMA transfer from XX00-XX9F to FE00-FE9F where XX = 00 to DF
// TODO check timing
func (mmu *MMU) DMA(xx uint8) {
	if xx > 0xDF {
		panic("DMA address out of range")
	}

	addr := uint16(xx) << 8
	for i := range 0xA0 {
		mmu.PPU.DMAWrite(i, mmu.Read(addr+uint16(i)))
	}
}
