package cpu

import (
	"github.com/danielecanzoneri/gb-emulator/internal/memory"
)

func mockCPU() *CPU {
	mem := &memory.MMU{}
	return &CPU{MMU: mem}
}

func writeTestProgram(cpu *CPU, data ...byte) {
	for i, b := range data {
		cpu.MMU.Write(uint16(i)+cpu.PC, b)
	}
}
