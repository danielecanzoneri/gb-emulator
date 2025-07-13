package cpu

func (cpu *CPU) SkipBoot() {
	cpu.writeAF(0x01B0)
	cpu.writeBC(0x0013)
	cpu.writeDE(0x00D8)
	cpu.writeHL(0x014D)
	cpu.SP = 0xFFFE
	cpu.PC = 0x100
}
