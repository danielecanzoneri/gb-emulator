package cpu

const (
	Z_FLAG_BIT = 7
	N_FLAG_BIT = 6
	H_FLAG_BIT = 5
	C_FLAG_BIT = 4
)

func (cpu *CPU) readAF() uint16 {
	return combineBytes(cpu.A, cpu.F)
}
func (cpu *CPU) readBC() uint16 {
	return combineBytes(cpu.B, cpu.C)
}
func (cpu *CPU) readDE() uint16 {
	return combineBytes(cpu.D, cpu.E)
}
func (cpu *CPU) readHL() uint16 {
	return combineBytes(cpu.H, cpu.L)
}

func (cpu *CPU) writeAF(word uint16) {
	cpu.A, cpu.F = splitWord(word)
}
func (cpu *CPU) writeBC(word uint16) {
	cpu.B, cpu.C = splitWord(word)
}
func (cpu *CPU) writeDE(word uint16) {
	cpu.D, cpu.E = splitWord(word)
}
func (cpu *CPU) writeHL(word uint16) {
	cpu.H, cpu.L = splitWord(word)
}

func (cpu *CPU) readZFlag() uint8 {
	return readBit(cpu.F, Z_FLAG_BIT)
}
func (cpu *CPU) readNFlag() uint8 {
	return readBit(cpu.F, N_FLAG_BIT)
}
func (cpu *CPU) readHFlag() uint8 {
	return readBit(cpu.F, H_FLAG_BIT)
}
func (cpu *CPU) readCFlag() uint8 {
	return readBit(cpu.F, C_FLAG_BIT)
}

func (cpu *CPU) setZFlag(value uint8) {
	setBit(&cpu.F, Z_FLAG_BIT, value)
}
func (cpu *CPU) setNFlag(value uint8) {
	setBit(&cpu.F, N_FLAG_BIT, value)
}
func (cpu *CPU) setHFlag(value uint8) {
	setBit(&cpu.F, H_FLAG_BIT, value)
}
func (cpu *CPU) setCFlag(value uint8) {
	setBit(&cpu.F, C_FLAG_BIT, value)
}
