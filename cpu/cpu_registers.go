package cpu

const (
	Z_FLAG_BIT = 7
	N_FLAG_BIT = 6
	H_FLAG_BIT = 5
	C_FLAG_BIT = 4
)

func (cpu *CPU) writeA(v uint8) {
	cpu.A = v
}
func (cpu *CPU) writeB(v uint8) {
	cpu.B = v
}
func (cpu *CPU) writeC(v uint8) {
	cpu.C = v
}
func (cpu *CPU) writeD(v uint8) {
	cpu.D = v
}
func (cpu *CPU) writeE(v uint8) {
	cpu.E = v
}
func (cpu *CPU) writeH(v uint8) {
	cpu.H = v
}
func (cpu *CPU) writeL(v uint8) {
	cpu.L = v
}
func (cpu *CPU) writeF(v uint8) {
	cpu.F = v & 0xF0
}

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
	h, l := splitWord(word)
	cpu.writeA(h)
	cpu.writeF(l)
}
func (cpu *CPU) writeBC(word uint16) {
	h, l := splitWord(word)
	cpu.writeB(h)
	cpu.writeC(l)
}
func (cpu *CPU) writeDE(word uint16) {
	h, l := splitWord(word)
	cpu.writeD(h)
	cpu.writeE(l)
}
func (cpu *CPU) writeHL(word uint16) {
	h, l := splitWord(word)
	cpu.writeH(h)
	cpu.writeL(l)
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
