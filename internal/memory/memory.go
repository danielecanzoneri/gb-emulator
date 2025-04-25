package memory

import (
	"github.com/danielecanzoneri/gb-emulator/internal/ppu"
	"github.com/danielecanzoneri/gb-emulator/internal/timer"
)

const Size = 0x10000 // 64KB

type MMU struct {
	Data [Size]uint8

	Timer *timer.Timer
	PPU   *ppu.PPU
}

func (mmu *MMU) Read(addr uint16) uint8 {
	switch {
	case vRAM <= addr && addr < eRAM:
		return mmu.PPU.ReadVRAM(addr)
	case echoRAM <= addr && addr < OAM:
		return mmu.Read(addr - 0x2000)
	case OAM <= addr && addr < reservedMemory:
		return mmu.PPU.ReadOAM(addr)
	case reservedMemory <= addr && addr < ioRegisters:
		panic("Can't read reserved memory")
	case ioRegisters <= addr && addr < hRAM:
		return mmu.readIO(addr)
	default:
		return mmu.Data[addr]
	}
}

func (mmu *MMU) Write(addr uint16, value uint8) {
	switch {
	case vRAM <= addr && addr < eRAM:
		mmu.PPU.WriteVRAM(addr, value)
	case echoRAM <= addr && addr < OAM:
		mmu.Write(addr-0x2000, value)
	// OAM inaccessible during PPU mode 2 and 3
	case OAM <= addr && addr < reservedMemory:
		mmu.PPU.WriteOAM(addr, value)
	case reservedMemory <= addr && addr < ioRegisters:
		panic("Can't write reserved memory")
	case ioRegisters <= addr && addr < hRAM:
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
