package memory

import (
	"github.com/danielecanzoneri/gb-emulator/emulator/internal/gameboy/audio"
	"github.com/danielecanzoneri/gb-emulator/emulator/internal/gameboy/joypad"
	"github.com/danielecanzoneri/gb-emulator/emulator/internal/gameboy/ppu"
	"github.com/danielecanzoneri/gb-emulator/emulator/internal/gameboy/timer"
)

const Size = 0x10000 // 64KB

type MMU struct {
	Data          [Size]uint8
	CartridgeData []uint8

	// Memory Bank Controller
	mbc MBC

	Timer  *timer.Timer
	PPU    *ppu.PPU
	Joypad *joypad.Joypad
	APU    *audio.APU

	// I/O registers
	dmaReg uint8
	ifReg  uint8
	ieReg  uint8

	// DMA cycles
	dmaStart      bool
	dmaWaitCycles uint8 // Wait two cycles before starting dma
	dmaTransfer   bool
	dmaOffset     uint16
	dmaValue      uint8
}

func (mmu *MMU) Reset() {
	mmu.Write(0xFF0F, 0xE1) // IF
	mmu.Write(0xFF40, 0x91) // LCDC
	mmu.Write(0xFF41, 0x81) // STAT
	mmu.Write(0xFF47, 0xFC) // BGP
}

func (mmu *MMU) Cycle() {
	if mmu.dmaTransfer {
		addr := uint16(mmu.read(dmaAddress)) << 8
		mmu.dmaValue = mmu.read(addr + mmu.dmaOffset)

		mmu.PPU.OAM.Data[mmu.dmaOffset] = mmu.dmaValue
		mmu.dmaOffset++

		if mmu.dmaOffset == dmaDuration {
			mmu.dmaTransfer = false
		}
	}

	// Start dma one cycle later
	if mmu.dmaStart {
		mmu.dmaWaitCycles--
		if mmu.dmaWaitCycles == 0 {
			mmu.dmaStart = false
			mmu.dmaTransfer = true
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
	case 0xFF00 <= addr && addr < 0xFF80 || addr == 0xFFFF: // I/O registers
		return mmu.readIO(addr)
	default:
		return mmu.Data[addr]
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
	case addr < 0x2000:
		mmu.mbc.enableRAM(value)
	case addr < 0x4000:
		mmu.mbc.SetROMBank(value)
	case addr < 0x6000:
		mmu.mbc.SetRAMBank(value)
	case addr < 0x8000:
		mmu.mbc.SetMode(value)
	case addr < 0xA000: // vRAM
		mmu.PPU.WriteVRAM(addr, value)
	case addr < 0xC000:
		RAMAddress := mmu.mbc.computeRAMAddress(addr)
		mmu.mbc.RAM[RAMAddress] = value
	// case addr < 0xE000 // wRAM
	case 0xE000 <= addr && addr < 0xFE00: // Echo RAM
		mmu.Write(addr-0x2000, value)
	case 0xFE00 <= addr && addr < 0xFEA0: // OAM
		mmu.PPU.WriteOAM(addr, value)
	case 0xFEA0 <= addr && addr < 0xFF00:
		//panic("Can't write reserved memory: " + strconv.FormatUint(uint64(addr), 16))
	case 0xFF00 <= addr && addr < 0xFF80 || addr == 0xFFFF: // I/O registers
		mmu.writeIO(addr, value)
	default:
		mmu.Data[addr] = value
	}
}

func (mmu *MMU) ReadWord(addr uint16) uint16 {
	return uint16(mmu.Read(addr)) | (uint16(mmu.Read(addr+1)) << 8)
}

func (mmu *MMU) WriteWord(addr uint16, value uint16) {
	mmu.Write(addr, uint8(value))
	mmu.Write(addr+1, uint8(value>>8))
}
