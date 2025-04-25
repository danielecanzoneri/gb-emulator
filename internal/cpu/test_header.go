package cpu

import (
	"github.com/danielecanzoneri/gb-emulator/internal/memory"
	"github.com/danielecanzoneri/gb-emulator/internal/ppu"
)

func mockCPU() *CPU {
	p := &ppu.PPU{}
	mem := &memory.MMU{PPU: p}
	return &CPU{SP: 0xFFFE, MMU: mem}
}

func writeTestProgram(cpu *CPU, data ...byte) {
	for i, b := range data {
		cpu.MMU.Write(uint16(i)+cpu.PC, b)
	}
}
