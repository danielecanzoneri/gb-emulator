package memory

import (
	"github.com/danielecanzoneri/gb-emulator/internal/ppu"
	"github.com/danielecanzoneri/gb-emulator/internal/timer"
	"strconv"
)

const Size = 0x10000 // 64KB

type MMU struct {
	Data          [Size]uint8
	CartridgeData []uint8

	// Memory Bank Controller
	mbc MBC

	Timer *timer.Timer
	PPU   *ppu.PPU
}

func (mmu *MMU) Read(addr uint16) uint8 {
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
		return mmu.Read(addr - 0x2000)
	case 0xFE00 <= addr && addr < 0xFEA0: // OAM
		return mmu.PPU.ReadOAM(addr)
	case 0xFEA0 <= addr && addr < 0xFF00:
		panic("Can't read reserved memory: " + strconv.FormatUint(uint64(addr), 16))
	case 0xFF00 <= addr && addr < 0xFF80: // I/O registers
		return mmu.readIO(addr)
	default:
		return mmu.Data[addr]
	}
}

func (mmu *MMU) Write(addr uint16, value uint8) {
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
		panic("Can't write reserved memory: " + strconv.FormatUint(uint64(addr), 16))
	case 0xFF00 <= addr && addr < 0xFF80: // I/O registers
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
