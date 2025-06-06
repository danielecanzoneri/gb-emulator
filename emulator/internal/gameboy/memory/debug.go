package memory

func (mmu *MMU) DebugRead(addr uint16) uint8 {
	return mmu.read(addr)
}
