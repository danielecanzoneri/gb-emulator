package mmu

func (mmu *MMU) DebugRead(addr uint16) uint8 {
	// Read Wave RAM
	if NR10Addr <= addr && addr < waveRAMStartAddr+waveRAMLength {
		return mmu.apu.DebugRead(addr)
	}
	return mmu.read(addr)
}
