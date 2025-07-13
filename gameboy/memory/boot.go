package memory

func (mmu *MMU) SkipBoot() {
	// Memory (from examining MMU after boot ROM)
	mmu.hRAM[0x7A] = 0x39
	mmu.hRAM[0x7B] = 0x01
	mmu.hRAM[0x7C] = 0x2E
	mmu.ifReg = 0xE1
}
