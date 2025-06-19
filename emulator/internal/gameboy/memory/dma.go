package memory

const (
	dmaDuration = 0xA0
	dmaAddress  = 0xFF46
)

// DMA transfer from XX00-XX9F to FE00-FE9F where XX = 00 to DF lasting 160 cycles
func (mmu *MMU) DMA(xx uint8) {
	if xx > 0xDF {
		panic("DMA address out of range")
	}

	mmu.dmaCycles = -8 // Wait two cycles before starting dma
	mmu.dmaTransfer = true
	mmu.dmaOffset = 0
	mmu.dmaReg = xx
}
