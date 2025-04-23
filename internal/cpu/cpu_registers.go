package cpu

import (
	"github.com/danielecanzoneri/gb-emulator/internal/util"
)

const (
	ZFlagBit = 7
	NFlagBit = 6
	HFlagBit = 5
	CFlagBit = 4
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
	return util.CombineBytes(cpu.A, cpu.F)
}
func (cpu *CPU) readBC() uint16 {
	return util.CombineBytes(cpu.B, cpu.C)
}
func (cpu *CPU) readDE() uint16 {
	return util.CombineBytes(cpu.D, cpu.E)
}
func (cpu *CPU) readHL() uint16 {
	return util.CombineBytes(cpu.H, cpu.L)
}

func (cpu *CPU) writeAF(word uint16) {
	h, l := util.SplitWord(word)
	cpu.writeA(h)
	cpu.writeF(l)
}
func (cpu *CPU) writeBC(word uint16) {
	h, l := util.SplitWord(word)
	cpu.writeB(h)
	cpu.writeC(l)
}
func (cpu *CPU) writeDE(word uint16) {
	h, l := util.SplitWord(word)
	cpu.writeD(h)
	cpu.writeE(l)
}
func (cpu *CPU) writeHL(word uint16) {
	h, l := util.SplitWord(word)
	cpu.writeH(h)
	cpu.writeL(l)
}

func (cpu *CPU) readZFlag() uint8 {
	return util.ReadBit(cpu.F, ZFlagBit)
}
func (cpu *CPU) readNFlag() uint8 {
	return util.ReadBit(cpu.F, NFlagBit)
}
func (cpu *CPU) readHFlag() uint8 {
	return util.ReadBit(cpu.F, HFlagBit)
}
func (cpu *CPU) readCFlag() uint8 {
	return util.ReadBit(cpu.F, CFlagBit)
}

func (cpu *CPU) setZFlag(value uint8) {
	util.SetBit(&cpu.F, ZFlagBit, value)
}
func (cpu *CPU) setNFlag(value uint8) {
	util.SetBit(&cpu.F, NFlagBit, value)
}
func (cpu *CPU) setHFlag(value uint8) {
	util.SetBit(&cpu.F, HFlagBit, value)
}
func (cpu *CPU) setCFlag(value uint8) {
	util.SetBit(&cpu.F, CFlagBit, value)
}
