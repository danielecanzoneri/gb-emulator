package mmu

func (mmu *MMU) DebugRead(addr uint16) uint8 {
	// Read Wave RAM
	if NR10Addr <= addr && addr < waveRAMStartAddr+waveRAMLength {
		return mmu.apu.DebugRead(addr)
	}
	return mmu.read(addr)
}

func (mmu *MMU) DebugGetVDMASrcAddress() uint16 {
	return mmu.vDMASrcAddress
}

func (mmu *MMU) DebugGetVDMADestAddress() uint16 {
	return 0x8000 + mmu.vDMADestAddress
}

func (mmu *MMU) DebugGetVDMALength() uint8 {
	return mmu.read(HDMA5Addr)
}
