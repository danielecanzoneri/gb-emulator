package mmu

import "github.com/danielecanzoneri/gb-emulator/gameboy/cartridge"

func (mmu *MMU) DisableBootROM() {
	mmu.BootRomDisabled = true

	// Set CGB compatibility mode
	if mmu.cgb && mmu.Cartridge.Header().CgbMode == cartridge.DmgOnly {
		mmu.ppu.DmgCompatibility = true
	}
}

func (mmu *MMU) SkipBoot() {
	// Memory (from examining MMU after boot ROM)
	mmu.hRAM[0x7A] = 0x39
	mmu.hRAM[0x7B] = 0x01
	mmu.hRAM[0x7C] = 0x2E
	mmu.ifReg = 0xE1
}
