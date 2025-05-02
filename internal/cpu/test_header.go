package cpu

import (
	"github.com/danielecanzoneri/gb-emulator/internal/cartridge"
	"github.com/danielecanzoneri/gb-emulator/internal/memory"
	"github.com/danielecanzoneri/gb-emulator/internal/ppu"
	"github.com/danielecanzoneri/gb-emulator/internal/timer"
)

func mockCPU() *CPU {
	p := &ppu.PPU{}
	mem := &memory.MMU{PPU: p, CartridgeData: make([]uint8, 0x8000), Timer: &timer.Timer{}}
	mem.SetMBC(&cartridge.Header{ROMBanks: 1})
	return &CPU{SP: 0xFFFE, MMU: mem}
}

func writeTestProgram(cpu *CPU, data ...byte) {
	for i, b := range data {
		addr := uint16(i) + cpu.PC
		if addr < 0x8000 {
			cpu.MMU.CartridgeData[addr] = b
		} else {
			cpu.MMU.Write(addr, b)
		}
	}
}
