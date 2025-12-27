package cpu

func (cpu *CPU) SkipDMGBoot() {
	cpu.writeAF(0x01B0)
	cpu.writeBC(0x0013)
	cpu.writeDE(0x00D8)
	cpu.writeHL(0x014D)
	cpu.SP = 0xFFFE
	cpu.PC = 0x100
}

func (cpu *CPU) SkipCGBBoot(bRegister uint8) {
	cpu.writeAF(0x1180)
	cpu.writeB(bRegister)
	cpu.writeC(0)
	cpu.writeDE(0x0008)
	
	// If the B register is $43 or $58 (on CGB) -> HL = $991A
	// Otherwise                                -> HL = $007C
	if bRegister == 0x43 || bRegister == 0x58 {
		cpu.writeHL(0x991A)
	} else {
		cpu.writeHL(0x007C)
	}
	cpu.SP = 0xFFFE
	cpu.PC = 0x100
}
