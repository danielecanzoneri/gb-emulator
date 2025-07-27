package memory

import (
	"github.com/danielecanzoneri/gb-emulator/gameboy/audio"
	"github.com/danielecanzoneri/gb-emulator/gameboy/cartridge"
	"github.com/danielecanzoneri/gb-emulator/gameboy/joypad"
	"github.com/danielecanzoneri/gb-emulator/gameboy/memory/serial"
	"github.com/danielecanzoneri/gb-emulator/gameboy/ppu"
	"github.com/danielecanzoneri/gb-emulator/gameboy/timer"
)

type MMU struct {
	wRAM [0x2000]uint8 // Work RAM
	hRAM [0x7F]uint8   // High RAM

	// Cartridge (with MBC and data)
	Cartridge cartridge.Cartridge

	Serial serial.Port
	Timer  *timer.Timer
	PPU    *ppu.PPU
	Joypad *joypad.Joypad
	APU    *audio.APU

	// I/O registers
	dmaReg uint8
	ifReg  uint8
	ieReg  uint8

	// DMA cycles
	delayDmaTicks int
	dmaTicks      int
	dmaTransfer   bool
	dmaOffset     uint16
	dmaValue      uint8

	// Boot ROM
	BootRomDisabled bool
	BootRom         []uint8
}

func (mmu *MMU) Tick(ticks uint) {
	if mmu.dmaTransfer {
		mmu.dmaTicks += int(ticks)
		for mmu.dmaTicks >= 4 {
			mmu.dmaTicks -= 4

			addr := uint16(mmu.read(dmaAddress)) << 8
			mmu.dmaValue = mmu.read(addr + mmu.dmaOffset)

			mmu.PPU.DMAWrite(mmu.dmaOffset, mmu.dmaValue)
			mmu.dmaOffset++

			if mmu.dmaOffset == dmaDuration {
				mmu.dmaTransfer = false
			}
		}
	}

	if mmu.delayDmaTicks > 0 {
		mmu.delayDmaTicks -= int(ticks)
		if mmu.delayDmaTicks <= 0 {
			mmu.dmaTransfer = true
			mmu.dmaTicks = 0
			mmu.dmaOffset = 0
		}
	}
}

func (mmu *MMU) Read(addr uint16) uint8 {
	// During DMA, HRAM can still be accessed otherwise return what DMA is reading
	//if mmu.dmaTransfer && !(0xFF00 <= addr && addr < 0xFFFF) {
	//	return mmu.dmaValue
	//}
	// OAM is inaccessible during DMA
	if mmu.dmaTransfer && 0xFE00 <= addr && addr < 0xFEA0 {
		return 0xFF
	}

	return mmu.read(addr)
}

func (mmu *MMU) read(addr uint16) uint8 {
	// Map boot ROM over memory
	if !mmu.BootRomDisabled && addr < 0x100 {
		return mmu.BootRom[addr]
	}

	switch {
	// MBC addresses
	case addr < 0x8000:
		return mmu.Cartridge.Read(addr)
	case addr < 0xA000: // vRAM
		return mmu.PPU.Read(addr)
	case addr < 0xC000:
		return mmu.Cartridge.Read(addr)
	case addr < 0xE000: // wRAM
		return mmu.wRAM[addr-0xC000]
	case addr < 0xFE00: // Echo RAM
		return mmu.read(addr - 0x2000)
	case addr < 0xFF00: // OAM
		return mmu.PPU.Read(addr)
	case addr < 0xFF80 || addr == 0xFFFF: // I/O registers
		return mmu.readIO(addr)
	default: // hRAM
		return mmu.hRAM[addr-0xFF80]
	}
}

func (mmu *MMU) Write(addr uint16, value uint8) {
	// During DMA, HRAM can still be accessed
	// if mmu.dmaTransfer && !(0xFF80 <= addr && addr < 0xFFFF) {
	// 	return
	// }
	// OAM is inaccessible during DMA
	if mmu.dmaTransfer && 0xFE00 <= addr && addr < 0xFEA0 {
		return
	}

	mmu.write(addr, value)
}

func (mmu *MMU) write(addr uint16, value uint8) {
	switch {
	// MBC addresses
	case addr < 0x8000:
		mmu.Cartridge.Write(addr, value)
	case addr < 0xA000: // vRAM
		mmu.PPU.Write(addr, value)
	case addr < 0xC000:
		mmu.Cartridge.Write(addr, value)
	case addr < 0xE000: // wRAM
		mmu.wRAM[addr-0xC000] = value
	case addr < 0xFE00: // Echo RAM
		mmu.Write(addr-0x2000, value)
	case addr < 0xFF00: // OAM
		mmu.PPU.Write(addr, value)
	case addr < 0xFF80 || addr == 0xFFFF: // I/O registers
		mmu.writeIO(addr, value)
	default: // hRAM
		mmu.hRAM[addr-0xFF80] = value
	}
}

func (mmu *MMU) ReadWord(addr uint16) uint16 {
	return uint16(mmu.Read(addr)) | (uint16(mmu.Read(addr+1)) << 8)
}

func (mmu *MMU) WriteWord(addr uint16, value uint16) {
	mmu.Write(addr, uint8(value))
	mmu.Write(addr+1, uint8(value>>8))
}
