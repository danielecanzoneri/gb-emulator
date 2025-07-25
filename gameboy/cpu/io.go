package cpu

import (
	"github.com/danielecanzoneri/gb-emulator/util"
)

// ReadByte uses 1 M-cycle
func (cpu *CPU) ReadByte(addr uint16) uint8 {
	defer cpu.Tick(4)
	return cpu.MMU.Read(addr)
}
func (cpu *CPU) WriteByte(addr uint16, v uint8) {
	defer cpu.Tick(4)
	cpu.MMU.Write(addr, v)
}

func (cpu *CPU) ReadNextByte() uint8 {
	b := cpu.ReadByte(cpu.PC)
	if cpu.haltBug {
		// Do not increment PC
		cpu.haltBug = false
	} else {
		// Increment PC (may cause OAM bug)
		cpu.PPU.TriggerOAMBug(cpu.PC)
		cpu.PC++
	}
	return b
}
func (cpu *CPU) WriteNextByte(v uint8) {
	cpu.WriteByte(cpu.PC, v)

	// Increment PC (may cause OAM bug)
	cpu.PPU.TriggerOAMBug(cpu.PC)
	cpu.PC++
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
