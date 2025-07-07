package memory

func (mmu *MMU) DebugRead(addr uint16) uint8 {
	// Read Wave RAM
	if NR10Addr <= addr && addr < waveRAMStartAddr+waveRAMLength {
		return mmu.APU.DebugRead(addr)
	}
	return mmu.read(addr)
}
