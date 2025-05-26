package memory

func (mmu *MMU) DebugRead(addr uint16) uint8 {
	switch {
	// MBC addresses
	case addr < 0x8000:
		cartridgeAddress := mmu.mbc.computeROMAddress(addr)
		return mmu.CartridgeData[cartridgeAddress]
	case addr < 0xA000: // vRAM
		return mmu.PPU.ReadVRAM(addr)
	case addr < 0xC000:
		RAMAddress := mmu.mbc.computeRAMAddress(addr)
		return mmu.mbc.RAM[RAMAddress]
	// case addr < 0xE000 // wRAM
	case 0xE000 <= addr && addr < 0xFE00: // Echo RAM
		return mmu.read(addr - 0x2000)
	case 0xFE00 <= addr && addr < 0xFEA0: // OAM
		return mmu.PPU.ReadOAM(addr)
	case 0xFEA0 <= addr && addr < 0xFF00:
		//panic("Can't read reserved memory: " + strconv.FormatUint(uint64(addr), 16))
		return 0
	case 0xFF30 <= addr && addr < 0xFF40: // Wave Ram
		return mmu.APU.DebugRead(addr)
	case 0xFF00 <= addr && addr < 0xFF80 || addr == 0xFFFF: // I/O registers
		return mmu.readIO(addr)
	default:
		return mmu.Data[addr]
	}
}
