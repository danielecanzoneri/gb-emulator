package mmu

import (
	"github.com/danielecanzoneri/gb-emulator/gameboy/audio"
	"github.com/danielecanzoneri/gb-emulator/gameboy/cartridge"
	"github.com/danielecanzoneri/gb-emulator/gameboy/joypad"
	"github.com/danielecanzoneri/gb-emulator/gameboy/ppu"
	"github.com/danielecanzoneri/gb-emulator/gameboy/serial"
	"github.com/danielecanzoneri/gb-emulator/gameboy/timer"
)

type MMU struct {
	wRAM [0x8000]uint8 // Work RAM (8 0x1000 banks for CGB, 2 for DMG)
	hRAM [0x7F]uint8   // High RAM

	// Cartridge (with MBC and data)
	Cartridge cartridge.Cartridge

	ppu    *ppu.PPU
	apu    *audio.APU
	timer  *timer.Timer
	joypad *joypad.Joypad
	serial *serial.Port

	// wRAM bank register
	vbk uint8

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

	// CGB flag
	cgb bool
}

func New(ppu *ppu.PPU, apu *audio.APU, timer *timer.Timer, jp *joypad.Joypad, serialPort *serial.Port, cgb bool) *MMU {
	return &MMU{
		ppu:    ppu,
		apu:    apu,
		timer:  timer,
		joypad: jp,
		serial: serialPort,
		cgb:    cgb,
	}
}

func (mmu *MMU) Tick(ticks int) {
	if mmu.dmaTransfer {
		mmu.dmaTicks += ticks
		for mmu.dmaTicks >= 4 {
			mmu.dmaTicks -= 4

			addr := uint16(mmu.read(dmaAddress)) << 8
			mmu.dmaValue = mmu.read(addr + mmu.dmaOffset)

			mmu.ppu.DMAWrite(mmu.dmaOffset, mmu.dmaValue)
			mmu.dmaOffset++

			if mmu.dmaOffset == dmaDuration {
				mmu.dmaTransfer = false
			}
		}
	}

	if mmu.delayDmaTicks > 0 {
		mmu.delayDmaTicks -= ticks
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
		return mmu.ppu.Read(addr)
	case addr < 0xC000:
		return mmu.Cartridge.Read(addr)
	case addr < 0xD000: // wRAM bank 0
		return mmu.wRAM[addr-0xC000]
	case addr < 0xE000: // wRAM (bank 1-7)
		baseAddr := addr - 0xD000

		var bank uint16 = 1
		if mmu.cgb && mmu.vbk > 0 {
			bank = uint16(mmu.vbk)
		}
		return mmu.wRAM[0x1000*bank+baseAddr]
	case addr < 0xFE00: // Echo RAM
		return mmu.read(addr - 0x2000)
	case addr < 0xFF00: // OAM
		return mmu.ppu.Read(addr)
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
		mmu.ppu.Write(addr, value)
	case addr < 0xC000:
		mmu.Cartridge.Write(addr, value)
	case addr < 0xD000: // wRAM bank 0
		mmu.wRAM[addr-0xC000] = value
	case addr < 0xE000: // wRAM (bank 1-7)
		baseAddr := addr - 0xD000

		var bank uint16 = 1
		if mmu.cgb && mmu.vbk > 0 {
			bank = uint16(mmu.vbk)
		}
		mmu.wRAM[0x1000*bank+baseAddr] = value
	case addr < 0xFE00: // Echo RAM
		mmu.Write(addr-0x2000, value)
	case addr < 0xFF00: // OAM
		mmu.ppu.Write(addr, value)
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
