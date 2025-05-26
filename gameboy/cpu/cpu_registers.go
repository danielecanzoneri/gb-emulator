package cpu

import (
	"github.com/danielecanzoneri/gb-emulator/gameboy/util"
)

const (
	ZFlagBit = 7
	NFlagBit = 6
	HFlagBit = 5
	CFlagBit = 4
)

func (cpu *CPU) readA() uint8 {
	return cpu.A
}
func (cpu *CPU) readB() uint8 {
	return cpu.B
}
func (cpu *CPU) readC() uint8 {
	return cpu.C
}
func (cpu *CPU) readD() uint8 {
	return cpu.D
}
func (cpu *CPU) readE() uint8 {
	return cpu.E
}
func (cpu *CPU) readH() uint8 {
	return cpu.H
}
func (cpu *CPU) readL() uint8 {
	return cpu.L
}
func (cpu *CPU) readF() uint8 {
	return cpu.F
}
func (cpu *CPU) readHighSP() uint8 {
	return uint8(cpu.SP >> 8)
}
func (cpu *CPU) readLowSP() uint8 {
	return uint8(cpu.SP & 0xFF)
}

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
func (cpu *CPU) writeHighSP(v uint8) {
	cpu.SP = (cpu.SP & 0x00FF) | (uint16(v) << 8)
}
func (cpu *CPU) writeLowSP(v uint8) {
	cpu.SP = (cpu.SP & 0xFF00) | uint16(v)
}

func (cpu *CPU) ReadAF() uint16 {
	return util.CombineBytes(cpu.A, cpu.F)
}
func (cpu *CPU) ReadBC() uint16 {
	return util.CombineBytes(cpu.B, cpu.C)
}
func (cpu *CPU) ReadDE() uint16 {
	return util.CombineBytes(cpu.D, cpu.E)
}
func (cpu *CPU) ReadHL() uint16 {
	return util.CombineBytes(cpu.H, cpu.L)
}
func (cpu *CPU) ReadSP() uint16 {
	return cpu.SP
}
func (cpu *CPU) ReadPC() uint16 {
	return cpu.PC
}
func (cpu *CPU) InterruptEnabled() bool {
	return cpu.IME
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
func (cpu *CPU) writeSP(word uint16) {
	cpu.SP = word
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
