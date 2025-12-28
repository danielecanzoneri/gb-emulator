package mmu

import (
	"github.com/danielecanzoneri/lucky-boy/gameboy/audio"
	"github.com/danielecanzoneri/lucky-boy/gameboy/cartridge"
	"github.com/danielecanzoneri/lucky-boy/gameboy/joypad"
	"github.com/danielecanzoneri/lucky-boy/gameboy/ppu"
	"github.com/danielecanzoneri/lucky-boy/gameboy/serial"
	"github.com/danielecanzoneri/lucky-boy/gameboy/timer"
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

	// vRAM DMA transfer
	vDMAActive      bool // If true, CPU will halt to complete vDMA
	vDMATicks       int
	vDMAHBlank      bool // True if vDMA is active but only in HBlank
	vDMASrcAddress  uint16
	vDMADestAddress uint16
	vDMALength      uint8
	IsCPUHalted     func() bool

	// CGB speed switch
	PrepareSpeedSwitch bool
	DoubleSpeed        bool

	// In double speed mode, DMA is at same speed of CPU, VDMA is at normal speed
	speedFactor int // (0: normal, 1: double)

	// Boot ROM
	BootRomDisabled bool
	BootRom         []uint8

	// CGB flag
	cgb bool
}

func New(ppu *ppu.PPU, apu *audio.APU, timer *timer.Timer, jp *joypad.Joypad, serialPort *serial.Port, cgb bool) *MMU {
	return &MMU{
		ppu:         ppu,
		apu:         apu,
		timer:       timer,
		joypad:      jp,
		serial:      serialPort,
		cgb:         cgb,
		speedFactor: 0,
	}
}

func (mmu *MMU) Tick(ticks int) {
	// In double speed mode, only DMA to OAM will be faster
	if mmu.dmaTransfer {
		mmu.dmaTicks += ticks << mmu.speedFactor
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
		mmu.delayDmaTicks -= ticks << mmu.speedFactor
		if mmu.delayDmaTicks <= 0 {
			mmu.dmaTransfer = true
			mmu.dmaTicks = 0
			mmu.dmaOffset = 0
		}
	}

	// vRAM DMA
	if mmu.vDMAActive {
		mmu.vDMATicks += ticks
		if mmu.vDMATicks >= 8*4 {
			mmu.vDMATicks -= 8 * 4

			// Actually transfer data and set up next transfer
			mmu.vDMATransfer()
		}
	}
}

func (mmu *MMU) SwitchSpeed(doubleSpeed bool) {
	mmu.PrepareSpeedSwitch = false
	mmu.DoubleSpeed = doubleSpeed
	if doubleSpeed {
		mmu.speedFactor = 1
	} else {
		mmu.speedFactor = 0
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
	if !mmu.BootRomDisabled {
		if addr < 0x100 {
			return mmu.BootRom[addr]
		}

		// In CGB mode the boot ROM is actually split in two parts, a $0000-00FF one, and a $0200-08FF one.
		if mmu.cgb && 0x200 <= addr && addr < 0x900 {
			return mmu.BootRom[addr]
		}
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
