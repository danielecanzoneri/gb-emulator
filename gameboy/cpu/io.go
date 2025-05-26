package cpu

import "github.com/danielecanzoneri/gb-emulator/gameboy/util"

// ReadByte uses 1 M-cycle
func (cpu *CPU) ReadByte(addr uint16) uint8 {
	defer cpu.Cycle()
	return cpu.MMU.Read(addr)
}
func (cpu *CPU) WriteByte(addr uint16, v uint8) {
	defer cpu.Cycle()
	cpu.MMU.Write(addr, v)
}

func (cpu *CPU) ReadNextByte() uint8 {
	b := cpu.ReadByte(cpu.PC)
	cpu.PC++
	return b
}
func (cpu *CPU) WriteNextByte(v uint8) {
	cpu.WriteByte(cpu.PC, v)
	cpu.PC++
}

// ReadWord uses 2 M-Cycles
func (cpu *CPU) ReadWord(addr uint16) uint16 {
	lowAddr := cpu.ReadByte(addr)
	highAddr := cpu.ReadByte(addr + 1)
	return util.CombineBytes(highAddr, lowAddr)

}
func (cpu *CPU) WriteWord(addr uint16, v uint16) {
	high, low := util.SplitWord(v)

	cpu.WriteByte(addr, low)
	cpu.WriteByte(addr+1, high)
}

func (cpu *CPU) ReadNextWord() uint16 {
	lowAddr := cpu.ReadNextByte()
	highAddr := cpu.ReadNextByte()
	return util.CombineBytes(highAddr, lowAddr)
}
func (cpu *CPU) WriteNextWord(v uint16) {
	high, low := util.SplitWord(v)

	cpu.WriteNextByte(low)
	cpu.WriteNextByte(high)
}

func (cpu *CPU) readHLmem() uint8 {
	addr := cpu.ReadHL()
	return cpu.ReadByte(addr)
}
func (cpu *CPU) writeHLmem(v uint8) {
	cpu.WriteByte(cpu.ReadHL(), v)
}
