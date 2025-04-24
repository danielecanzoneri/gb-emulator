package memory

import (
	"github.com/danielecanzoneri/gb-emulator/internal/ppu"
	"github.com/danielecanzoneri/gb-emulator/internal/timer"
)

const Size = 0x10000 // 64KB

type MMU struct {
	data [Size]uint8

	Timer *timer.Timer
	PPU   *ppu.PPU
}

func (mmu *MMU) Read(addr uint16) uint8 {
	switch {
	case ioRegisters <= addr && addr < hRAM:
		return mmu.readIO(addr)
	default:
		return mmu.data[addr]
	}
}

func (mmu *MMU) Write(addr uint16, value uint8) {
	switch {
	case ioRegisters <= addr && addr < hRAM:
		mmu.writeIO(addr, value)
	default:
		mmu.data[addr] = value
	}
}

func (mmu *MMU) ReadWord(addr uint16) uint16 {
	return uint16(mmu.Read(addr)) | (uint16(mmu.Read(addr+1)) << 8)
}

func (mmu *MMU) WriteWord(addr uint16, value uint16) {
	mmu.Write(addr, uint8(value))
	mmu.Write(addr+1, uint8(value>>8))
}
