package cpu

import (
	"github.com/danielecanzoneri/gb-emulator/internal/util"
)

func (cpu *CPU) NOP() {
}

// LD R16 N16
func (cpu *CPU) LD_R16_N16(writeHigh, writeLow func(uint8)) {
	low := cpu.ReadNextByte()
	writeLow(low)

	high := cpu.ReadNextByte()
	writeHigh(high)
}
func (cpu *CPU) LD_BC_N16() {
	cpu.LD_R16_N16(cpu.writeB, cpu.writeC)
}
func (cpu *CPU) LD_DE_N16() {
	cpu.LD_R16_N16(cpu.writeD, cpu.writeE)
}
func (cpu *CPU) LD_HL_N16() {
	cpu.LD_R16_N16(cpu.writeH, cpu.writeL)
}
func (cpu *CPU) LD_SP_N16() {
	cpu.LD_R16_N16(cpu.writeHighSP, cpu.writeLowSP)
}

// LD [R16] A
func (cpu *CPU) LD_R16mem_A(readR16 func() uint16) {
	cpu.WriteByte(readR16(), cpu.A)
}
func (cpu *CPU) LD_BCmem_A() {
	cpu.LD_R16mem_A(cpu.readBC)
}
func (cpu *CPU) LD_DEmem_A() {
	cpu.LD_R16mem_A(cpu.readDE)
}
func (cpu *CPU) LD_HLImem_A() {
	cpu.LD_R16mem_A(cpu.readHL)
	cpu.writeHL(cpu.readHL() + 1)
}
func (cpu *CPU) LD_HLDmem_A() {
	cpu.LD_R16mem_A(cpu.readHL)
	cpu.writeHL(cpu.readHL() - 1)
}

// LD A [R16]
func (cpu *CPU) LD_A_R16mem(readR16 func() uint16) {
	cpu.A = cpu.ReadByte(readR16())
}
func (cpu *CPU) LD_A_BCmem() {
	cpu.LD_A_R16mem(cpu.readBC)
}
func (cpu *CPU) LD_A_DEmem() {
	cpu.LD_A_R16mem(cpu.readDE)
}
func (cpu *CPU) LD_A_HLImem() {
	cpu.LD_A_R16mem(cpu.readHL)
	cpu.writeHL(cpu.readHL() + 1)
}
func (cpu *CPU) LD_A_HLDmem() {
	cpu.LD_A_R16mem(cpu.readHL)
	cpu.writeHL(cpu.readHL() - 1)
}

// LD N16 SP
func (cpu *CPU) LD_N16_SP() {
	addr := cpu.ReadNextWord()
	cpu.WriteWord(addr, cpu.SP)
}

// INC R16
func (cpu *CPU) INC_R16(readR16 func() uint16, writeR16 func(uint16)) {
	writeR16(readR16() + 1)
	cpu.Cycle()
}
func (cpu *CPU) INC_BC() {
	cpu.INC_R16(cpu.readBC, cpu.writeBC)
}
func (cpu *CPU) INC_DE() {
	cpu.INC_R16(cpu.readDE, cpu.writeDE)
}
func (cpu *CPU) INC_HL() {
	cpu.INC_R16(cpu.readHL, cpu.writeHL)
}
func (cpu *CPU) INC_SP() {
	cpu.INC_R16(cpu.readSP, cpu.writeSP)
}

// DEC R16
func (cpu *CPU) DEC_R16(readR16 func() uint16, writeR16 func(uint16)) {
	writeR16(readR16() - 1)
	cpu.Cycle()
}
func (cpu *CPU) DEC_BC() {
	cpu.DEC_R16(cpu.readBC, cpu.writeBC)
}
func (cpu *CPU) DEC_DE() {
	cpu.DEC_R16(cpu.readDE, cpu.writeDE)
}
func (cpu *CPU) DEC_HL() {
	cpu.DEC_R16(cpu.readHL, cpu.writeHL)
}
func (cpu *CPU) DEC_SP() {
	cpu.DEC_R16(cpu.readSP, cpu.writeSP)
}

// ADD HL R16
func (cpu *CPU) ADD_HL_R16(readR16 func() uint16) {
	sum, carry, half_carry := util.SumWordsWithCarry(cpu.readHL(), readR16())
	cpu.setNFlag(0)
	cpu.setHFlag(half_carry)
	cpu.setCFlag(carry)
	cpu.writeHL(sum)

	cpu.Cycle()
}
func (cpu *CPU) ADD_HL_BC() {
	cpu.ADD_HL_R16(cpu.readBC)
}
func (cpu *CPU) ADD_HL_DE() {
	cpu.ADD_HL_R16(cpu.readDE)
}
func (cpu *CPU) ADD_HL_HL() {
	cpu.ADD_HL_R16(cpu.readHL)
}
func (cpu *CPU) ADD_HL_SP() {
	cpu.ADD_HL_R16(cpu.readSP)
}

// INC R8
func (cpu *CPU) INC_R8(readR8 func() uint8, writeR8 func(uint8)) {
	// Increments r8 and set correct flags
	sum, carry, half_carry := util.SumBytesWithCarry(readR8(), 1)
	cpu.setNFlag(0)
	cpu.setHFlag(half_carry)
	cpu.setZFlag(carry) // if carry it means result is 0

	writeR8(sum)
}
func (cpu *CPU) INC_B() {
	cpu.INC_R8(cpu.readB, cpu.writeB)
}
func (cpu *CPU) INC_C() {
	cpu.INC_R8(cpu.readC, cpu.writeC)
}
func (cpu *CPU) INC_D() {
	cpu.INC_R8(cpu.readD, cpu.writeD)
}
func (cpu *CPU) INC_E() {
	cpu.INC_R8(cpu.readE, cpu.writeE)
}
func (cpu *CPU) INC_H() {
	cpu.INC_R8(cpu.readH, cpu.writeH)
}
func (cpu *CPU) INC_L() {
	cpu.INC_R8(cpu.readL, cpu.writeL)
}
func (cpu *CPU) INC_HLmem() {
	cpu.INC_R8(cpu.readHLmem, cpu.writeHLmem)
}
func (cpu *CPU) INC_A() {
	cpu.INC_R8(cpu.readA, cpu.writeA)
}

// DEC R8
func (cpu *CPU) DEC_R8(readR8 func() uint8, writeR8 func(uint8)) {
	// Decrements r8 and set correct flags
	sub, _, half_carry := util.SubBytesWithCarry(readR8(), 1)
	cpu.setNFlag(1)
	cpu.setHFlag(half_carry)
	cpu.setZFlag(util.IsByteZeroUint8(sub))

	writeR8(sub)
}
func (cpu *CPU) DEC_B() {
	cpu.DEC_R8(cpu.readB, cpu.writeB)
}
func (cpu *CPU) DEC_C() {
	cpu.DEC_R8(cpu.readC, cpu.writeC)
}
func (cpu *CPU) DEC_D() {
	cpu.DEC_R8(cpu.readD, cpu.writeD)
}
func (cpu *CPU) DEC_E() {
	cpu.DEC_R8(cpu.readE, cpu.writeE)
}
func (cpu *CPU) DEC_H() {
	cpu.DEC_R8(cpu.readH, cpu.writeH)
}
func (cpu *CPU) DEC_L() {
	cpu.DEC_R8(cpu.readL, cpu.writeL)
}
func (cpu *CPU) DEC_HLmem() {
	cpu.DEC_R8(cpu.readHLmem, cpu.writeHLmem)
}
func (cpu *CPU) DEC_A() {
	cpu.DEC_R8(cpu.readA, cpu.writeA)
}

// LD R8 N8
func (cpu *CPU) LD_R8_N8(writeR8 func(uint8)) {
	v := cpu.ReadNextByte()
	writeR8(v)
}
func (cpu *CPU) LD_B_N8() {
	cpu.LD_R8_N8(cpu.writeB)
}
func (cpu *CPU) LD_C_N8() {
	cpu.LD_R8_N8(cpu.writeC)
}
func (cpu *CPU) LD_D_N8() {
	cpu.LD_R8_N8(cpu.writeD)
}
func (cpu *CPU) LD_E_N8() {
	cpu.LD_R8_N8(cpu.writeE)
}
func (cpu *CPU) LD_H_N8() {
	cpu.LD_R8_N8(cpu.writeH)
}
func (cpu *CPU) LD_L_N8() {
	cpu.LD_R8_N8(cpu.writeL)
}
func (cpu *CPU) LD_HLmem_N8() {
	cpu.LD_R8_N8(cpu.writeHLmem)
}
func (cpu *CPU) LD_A_N8() {
	cpu.LD_R8_N8(cpu.writeA)
}

// 8-bit logic
func (cpu *CPU) RLCA() {
	cFlag := cpu.A >> 7
	cpu.writeA((cpu.A << 1) | cFlag)
	cpu.setZFlag(0)
	cpu.setNFlag(0)
	cpu.setHFlag(0)
	cpu.setCFlag(cFlag)
}
func (cpu *CPU) RRCA() {
	cFlag := cpu.A & 0x01
	cpu.writeA((cpu.A >> 1) | (cFlag << 7))
	cpu.setZFlag(0)
	cpu.setNFlag(0)
	cpu.setHFlag(0)
	cpu.setCFlag(cFlag)
}
func (cpu *CPU) RLA() {
	newA := (cpu.A << 1) | cpu.readCFlag()
	cFlag := cpu.A >> 7
	cpu.writeA(newA)
	cpu.setZFlag(0)
	cpu.setNFlag(0)
	cpu.setHFlag(0)
	cpu.setCFlag(cFlag)
}
func (cpu *CPU) RRA() {
	newA := (cpu.A >> 1) | (cpu.readCFlag() << 7)
	cFlag := cpu.A & 0x01
	cpu.writeA(newA)
	cpu.setZFlag(0)
	cpu.setNFlag(0)
	cpu.setHFlag(0)
	cpu.setCFlag(cFlag)
}

// 8-bit arithmetic
func (cpu *CPU) DAA() {
	var adj uint8 = 0

	if cpu.readNFlag() == 1 { // Subtraction
		if cpu.readHFlag() == 1 {
			adj += 0x6
		}
		if cpu.readCFlag() == 1 {
			adj += 0x60
		}
		cpu.A -= adj
	} else { // Addition
		if cpu.readHFlag() == 1 || (cpu.A&0x0F) > 0x9 {
			adj += 0x6
		}
		if cpu.readCFlag() == 1 || cpu.A > 0x99 {
			adj += 0x60
			cpu.setCFlag(1)
		} else {
			cpu.setCFlag(0)
		}
		cpu.A += adj
	}
	cpu.setHFlag(0)
	cpu.setZFlag(util.IsByteZeroUint8(cpu.A))
}
func (cpu *CPU) CPL() {
	cpu.writeA(^cpu.A)
	cpu.setNFlag(1)
	cpu.setHFlag(1)
}
func (cpu *CPU) SCF() {
	cpu.setNFlag(0)
	cpu.setHFlag(0)
	cpu.setCFlag(1)
}
func (cpu *CPU) CCF() {
	cpu.setNFlag(0)
	cpu.setHFlag(0)
	cpu.setCFlag(1 - cpu.readCFlag())
}

// JR COND E8
func (cpu *CPU) JR_COND_E8(checkCondition func() bool) {
	e8 := cpu.ReadNextByte()

	if checkCondition() {
		cpu.PC = uint16(int(cpu.PC) + int(int8(e8)))
		cpu.Cycle()
	}
}
func (cpu *CPU) JR_Z_E8() {
	cpu.JR_COND_E8(func() bool { return cpu.readZFlag() == 1 })
}
func (cpu *CPU) JR_NZ_E8() {
	cpu.JR_COND_E8(func() bool { return cpu.readZFlag() == 0 })
}
func (cpu *CPU) JR_C_E8() {
	cpu.JR_COND_E8(func() bool { return cpu.readCFlag() == 1 })
}
func (cpu *CPU) JR_NC_E8() {
	cpu.JR_COND_E8(func() bool { return cpu.readCFlag() == 0 })
}

// JR E8
func (cpu *CPU) JR_E8() {
	cpu.JR_COND_E8(func() bool { return true })
}

// STOP
func (cpu *CPU) STOP() {
	cpu.PC++
}

// LD R8 R8
func (cpu *CPU) LD_R8_R8(readR8 func() uint8, writeR8 func(uint8)) {
	writeR8(readR8())
}

func (cpu *CPU) LD_B_B() {
	cpu.LD_R8_R8(cpu.readB, cpu.writeB)
}
func (cpu *CPU) LD_B_C() {
	cpu.LD_R8_R8(cpu.readC, cpu.writeB)
}
func (cpu *CPU) LD_B_D() {
	cpu.LD_R8_R8(cpu.readD, cpu.writeB)
}
func (cpu *CPU) LD_B_E() {
	cpu.LD_R8_R8(cpu.readE, cpu.writeB)
}
func (cpu *CPU) LD_B_H() {
	cpu.LD_R8_R8(cpu.readH, cpu.writeB)
}
func (cpu *CPU) LD_B_L() {
	cpu.LD_R8_R8(cpu.readL, cpu.writeB)
}
func (cpu *CPU) LD_B_HLmem() {
	cpu.LD_R8_R8(cpu.readHLmem, cpu.writeB)
}
func (cpu *CPU) LD_B_A() {
	cpu.LD_R8_R8(cpu.readA, cpu.writeB)
}

func (cpu *CPU) LD_C_B() {
	cpu.LD_R8_R8(cpu.readB, cpu.writeC)
}
func (cpu *CPU) LD_C_C() {
	cpu.LD_R8_R8(cpu.readC, cpu.writeC)
}
func (cpu *CPU) LD_C_D() {
	cpu.LD_R8_R8(cpu.readD, cpu.writeC)
}
func (cpu *CPU) LD_C_E() {
	cpu.LD_R8_R8(cpu.readE, cpu.writeC)
}
func (cpu *CPU) LD_C_H() {
	cpu.LD_R8_R8(cpu.readH, cpu.writeC)
}
func (cpu *CPU) LD_C_L() {
	cpu.LD_R8_R8(cpu.readL, cpu.writeC)
}
func (cpu *CPU) LD_C_HLmem() {
	cpu.LD_R8_R8(cpu.readHLmem, cpu.writeC)
}
func (cpu *CPU) LD_C_A() {
	cpu.LD_R8_R8(cpu.readA, cpu.writeC)
}

func (cpu *CPU) LD_D_B() {
	cpu.LD_R8_R8(cpu.readB, cpu.writeD)
}
func (cpu *CPU) LD_D_C() {
	cpu.LD_R8_R8(cpu.readC, cpu.writeD)
}
func (cpu *CPU) LD_D_D() {
	cpu.LD_R8_R8(cpu.readD, cpu.writeD)
}
func (cpu *CPU) LD_D_E() {
	cpu.LD_R8_R8(cpu.readE, cpu.writeD)
}
func (cpu *CPU) LD_D_H() {
	cpu.LD_R8_R8(cpu.readH, cpu.writeD)
}
func (cpu *CPU) LD_D_L() {
	cpu.LD_R8_R8(cpu.readL, cpu.writeD)
}
func (cpu *CPU) LD_D_HLmem() {
	cpu.LD_R8_R8(cpu.readHLmem, cpu.writeD)
}
func (cpu *CPU) LD_D_A() {
	cpu.LD_R8_R8(cpu.readA, cpu.writeD)
}

func (cpu *CPU) LD_E_B() {
	cpu.LD_R8_R8(cpu.readB, cpu.writeE)
}
func (cpu *CPU) LD_E_C() {
	cpu.LD_R8_R8(cpu.readC, cpu.writeE)
}
func (cpu *CPU) LD_E_D() {
	cpu.LD_R8_R8(cpu.readD, cpu.writeE)
}
func (cpu *CPU) LD_E_E() {
	cpu.LD_R8_R8(cpu.readE, cpu.writeE)
}
func (cpu *CPU) LD_E_H() {
	cpu.LD_R8_R8(cpu.readH, cpu.writeE)
}
func (cpu *CPU) LD_E_L() {
	cpu.LD_R8_R8(cpu.readL, cpu.writeE)
}
func (cpu *CPU) LD_E_HLmem() {
	cpu.LD_R8_R8(cpu.readHLmem, cpu.writeE)
}
func (cpu *CPU) LD_E_A() {
	cpu.LD_R8_R8(cpu.readA, cpu.writeE)
}

func (cpu *CPU) LD_H_B() {
	cpu.LD_R8_R8(cpu.readB, cpu.writeH)
}
func (cpu *CPU) LD_H_C() {
	cpu.LD_R8_R8(cpu.readC, cpu.writeH)
}
func (cpu *CPU) LD_H_D() {
	cpu.LD_R8_R8(cpu.readD, cpu.writeH)
}
func (cpu *CPU) LD_H_E() {
	cpu.LD_R8_R8(cpu.readE, cpu.writeH)
}
func (cpu *CPU) LD_H_H() {
	cpu.LD_R8_R8(cpu.readH, cpu.writeH)
}
func (cpu *CPU) LD_H_L() {
	cpu.LD_R8_R8(cpu.readL, cpu.writeH)
}
func (cpu *CPU) LD_H_HLmem() {
	cpu.LD_R8_R8(cpu.readHLmem, cpu.writeH)
}
func (cpu *CPU) LD_H_A() {
	cpu.LD_R8_R8(cpu.readA, cpu.writeH)
}

func (cpu *CPU) LD_L_B() {
	cpu.LD_R8_R8(cpu.readB, cpu.writeL)
}
func (cpu *CPU) LD_L_C() {
	cpu.LD_R8_R8(cpu.readC, cpu.writeL)
}
func (cpu *CPU) LD_L_D() {
	cpu.LD_R8_R8(cpu.readD, cpu.writeL)
}
func (cpu *CPU) LD_L_E() {
	cpu.LD_R8_R8(cpu.readE, cpu.writeL)
}
func (cpu *CPU) LD_L_H() {
	cpu.LD_R8_R8(cpu.readH, cpu.writeL)
}
func (cpu *CPU) LD_L_L() {
	cpu.LD_R8_R8(cpu.readL, cpu.writeL)
}
func (cpu *CPU) LD_L_HLmem() {
	cpu.LD_R8_R8(cpu.readHLmem, cpu.writeL)
}
func (cpu *CPU) LD_L_A() {
	cpu.LD_R8_R8(cpu.readA, cpu.writeL)
}

func (cpu *CPU) LD_HLmem_B() {
	cpu.LD_R8_R8(cpu.readB, cpu.writeHLmem)
}
func (cpu *CPU) LD_HLmem_C() {
	cpu.LD_R8_R8(cpu.readC, cpu.writeHLmem)
}
func (cpu *CPU) LD_HLmem_D() {
	cpu.LD_R8_R8(cpu.readD, cpu.writeHLmem)
}
func (cpu *CPU) LD_HLmem_E() {
	cpu.LD_R8_R8(cpu.readE, cpu.writeHLmem)
}
func (cpu *CPU) LD_HLmem_H() {
	cpu.LD_R8_R8(cpu.readH, cpu.writeHLmem)
}
func (cpu *CPU) LD_HLmem_L() {
	cpu.LD_R8_R8(cpu.readL, cpu.writeHLmem)
}
func (cpu *CPU) LD_HLmem_A() {
	cpu.LD_R8_R8(cpu.readA, cpu.writeHLmem)
}

func (cpu *CPU) LD_A_B() {
	cpu.LD_R8_R8(cpu.readB, cpu.writeA)
}
func (cpu *CPU) LD_A_C() {
	cpu.LD_R8_R8(cpu.readC, cpu.writeA)
}
func (cpu *CPU) LD_A_D() {
	cpu.LD_R8_R8(cpu.readD, cpu.writeA)
}
func (cpu *CPU) LD_A_E() {
	cpu.LD_R8_R8(cpu.readE, cpu.writeA)
}
func (cpu *CPU) LD_A_H() {
	cpu.LD_R8_R8(cpu.readH, cpu.writeA)
}
func (cpu *CPU) LD_A_L() {
	cpu.LD_R8_R8(cpu.readL, cpu.writeA)
}
func (cpu *CPU) LD_A_HLmem() {
	cpu.LD_R8_R8(cpu.readHLmem, cpu.writeA)
}
func (cpu *CPU) LD_A_A() {
	cpu.LD_R8_R8(cpu.readA, cpu.writeA)
}

// HALT
func (cpu *CPU) HALT() {
	cpu.halted = true
}

// ADD A R8
func (cpu *CPU) ADD_A_R8(readR8 func() uint8) {
	sum, carry, halfCarry := util.SumBytesWithCarry(cpu.A, readR8())

	cpu.setZFlag(util.IsByteZeroUint8(sum))
	cpu.setNFlag(0)
	cpu.setHFlag(halfCarry)
	cpu.setCFlag(carry)

	cpu.writeA(sum)
}
func (cpu *CPU) ADD_A_B() {
	cpu.ADD_A_R8(cpu.readB)
}
func (cpu *CPU) ADD_A_C() {
	cpu.ADD_A_R8(cpu.readC)
}
func (cpu *CPU) ADD_A_D() {
	cpu.ADD_A_R8(cpu.readD)
}
func (cpu *CPU) ADD_A_E() {
	cpu.ADD_A_R8(cpu.readE)
}
func (cpu *CPU) ADD_A_H() {
	cpu.ADD_A_R8(cpu.readH)
}
func (cpu *CPU) ADD_A_L() {
	cpu.ADD_A_R8(cpu.readL)
}
func (cpu *CPU) ADD_A_HLmem() {
	cpu.ADD_A_R8(cpu.readHLmem)
}
func (cpu *CPU) ADD_A_A() {
	cpu.ADD_A_R8(cpu.readA)
}

// ADC A R8
func (cpu *CPU) ADC_A_R8(readR8 func() uint8) {
	sum, carry1, halfCarry1 := util.SumBytesWithCarry(cpu.A, readR8())
	sum, carry2, halfCarry2 := util.SumBytesWithCarry(sum, cpu.readCFlag())

	cpu.setZFlag(util.IsByteZeroUint8(sum))
	cpu.setNFlag(0)
	cpu.setHFlag(halfCarry1 | halfCarry2)
	cpu.setCFlag(carry1 | carry2)

	cpu.writeA(sum)
}
func (cpu *CPU) ADC_A_B() {
	cpu.ADC_A_R8(cpu.readB)
}
func (cpu *CPU) ADC_A_C() {
	cpu.ADC_A_R8(cpu.readC)
}
func (cpu *CPU) ADC_A_D() {
	cpu.ADC_A_R8(cpu.readD)
}
func (cpu *CPU) ADC_A_E() {
	cpu.ADC_A_R8(cpu.readE)
}
func (cpu *CPU) ADC_A_H() {
	cpu.ADC_A_R8(cpu.readH)
}
func (cpu *CPU) ADC_A_L() {
	cpu.ADC_A_R8(cpu.readL)
}
func (cpu *CPU) ADC_A_HLmem() {
	cpu.ADC_A_R8(cpu.readHLmem)
}
func (cpu *CPU) ADC_A_A() {
	cpu.ADC_A_R8(cpu.readA)
}

// SUB A R8
func (cpu *CPU) SUB_A_R8(readR8 func() uint8) {
	sub, carry, halfCarry := util.SubBytesWithCarry(cpu.A, readR8())

	cpu.setZFlag(util.IsByteZeroUint8(sub))
	cpu.setNFlag(1)
	cpu.setHFlag(halfCarry)
	cpu.setCFlag(carry)

	cpu.writeA(sub)
}
func (cpu *CPU) SUB_A_B() {
	cpu.SUB_A_R8(cpu.readB)
}
func (cpu *CPU) SUB_A_C() {
	cpu.SUB_A_R8(cpu.readC)
}
func (cpu *CPU) SUB_A_D() {
	cpu.SUB_A_R8(cpu.readD)
}
func (cpu *CPU) SUB_A_E() {
	cpu.SUB_A_R8(cpu.readE)
}
func (cpu *CPU) SUB_A_H() {
	cpu.SUB_A_R8(cpu.readH)
}
func (cpu *CPU) SUB_A_L() {
	cpu.SUB_A_R8(cpu.readL)
}
func (cpu *CPU) SUB_A_HLmem() {
	cpu.SUB_A_R8(cpu.readHLmem)
}
func (cpu *CPU) SUB_A_A() {
	cpu.SUB_A_R8(cpu.readA)
}

// SBC A R8
func (cpu *CPU) SBC_A_R8(readR8 func() uint8) {
	sub, carry1, halfCarry1 := util.SubBytesWithCarry(cpu.A, readR8())
	sub, carry2, halfCarry2 := util.SubBytesWithCarry(sub, cpu.readCFlag())

	cpu.setZFlag(util.IsByteZeroUint8(sub))
	cpu.setNFlag(1)
	cpu.setHFlag(halfCarry1 | halfCarry2)
	cpu.setCFlag(carry1 | carry2)

	cpu.writeA(sub)
}
func (cpu *CPU) SBC_A_B() {
	cpu.SBC_A_R8(cpu.readB)
}
func (cpu *CPU) SBC_A_C() {
	cpu.SBC_A_R8(cpu.readC)
}
func (cpu *CPU) SBC_A_D() {
	cpu.SBC_A_R8(cpu.readD)
}
func (cpu *CPU) SBC_A_E() {
	cpu.SBC_A_R8(cpu.readE)
}
func (cpu *CPU) SBC_A_H() {
	cpu.SBC_A_R8(cpu.readH)
}
func (cpu *CPU) SBC_A_L() {
	cpu.SBC_A_R8(cpu.readL)
}
func (cpu *CPU) SBC_A_HLmem() {
	cpu.SBC_A_R8(cpu.readHLmem)
}
func (cpu *CPU) SBC_A_A() {
	cpu.SBC_A_R8(cpu.readA)
}

// AND A R8
func (cpu *CPU) AND_A_R8(readR8 func() uint8) {
	cpu.writeA(cpu.A & readR8())
	cpu.setZFlag(util.IsByteZeroUint8(cpu.A))
	cpu.setNFlag(0)
	cpu.setHFlag(1)
	cpu.setCFlag(0)
}
func (cpu *CPU) AND_A_B() {
	cpu.AND_A_R8(cpu.readB)
}
func (cpu *CPU) AND_A_C() {
	cpu.AND_A_R8(cpu.readC)
}
func (cpu *CPU) AND_A_D() {
	cpu.AND_A_R8(cpu.readD)
}
func (cpu *CPU) AND_A_E() {
	cpu.AND_A_R8(cpu.readE)
}
func (cpu *CPU) AND_A_H() {
	cpu.AND_A_R8(cpu.readH)
}
func (cpu *CPU) AND_A_L() {
	cpu.AND_A_R8(cpu.readL)
}
func (cpu *CPU) AND_A_HLmem() {
	cpu.AND_A_R8(cpu.readHLmem)
}
func (cpu *CPU) AND_A_A() {
	cpu.AND_A_R8(cpu.readA)
}

// XOR A R8
func (cpu *CPU) XOR_A_R8(readR8 func() uint8) {
	cpu.writeA(cpu.A ^ readR8())
	cpu.setZFlag(util.IsByteZeroUint8(cpu.A))
	cpu.setNFlag(0)
	cpu.setHFlag(0)
	cpu.setCFlag(0)
}
func (cpu *CPU) XOR_A_B() {
	cpu.XOR_A_R8(cpu.readB)
}
func (cpu *CPU) XOR_A_C() {
	cpu.XOR_A_R8(cpu.readC)
}
func (cpu *CPU) XOR_A_D() {
	cpu.XOR_A_R8(cpu.readD)
}
func (cpu *CPU) XOR_A_E() {
	cpu.XOR_A_R8(cpu.readE)
}
func (cpu *CPU) XOR_A_H() {
	cpu.XOR_A_R8(cpu.readH)
}
func (cpu *CPU) XOR_A_L() {
	cpu.XOR_A_R8(cpu.readL)
}
func (cpu *CPU) XOR_A_HLmem() {
	cpu.XOR_A_R8(cpu.readHLmem)
}
func (cpu *CPU) XOR_A_A() {
	cpu.XOR_A_R8(cpu.readA)
}

// OR A R8
func (cpu *CPU) OR_A_R8(readR8 func() uint8) {
	cpu.writeA(cpu.A | readR8())
	cpu.setZFlag(util.IsByteZeroUint8(cpu.A))
	cpu.setNFlag(0)
	cpu.setHFlag(0)
	cpu.setCFlag(0)
}
func (cpu *CPU) OR_A_B() {
	cpu.OR_A_R8(cpu.readB)
}
func (cpu *CPU) OR_A_C() {
	cpu.OR_A_R8(cpu.readC)
}
func (cpu *CPU) OR_A_D() {
	cpu.OR_A_R8(cpu.readD)
}
func (cpu *CPU) OR_A_E() {
	cpu.OR_A_R8(cpu.readE)
}
func (cpu *CPU) OR_A_H() {
	cpu.OR_A_R8(cpu.readH)
}
func (cpu *CPU) OR_A_L() {
	cpu.OR_A_R8(cpu.readL)
}
func (cpu *CPU) OR_A_HLmem() {
	cpu.OR_A_R8(cpu.readHLmem)
}
func (cpu *CPU) OR_A_A() {
	cpu.OR_A_R8(cpu.readA)
}

// CP A R8
func (cpu *CPU) CP_A_R8(readR8 func() uint8) {
	sub, carry, halfCarry := util.SubBytesWithCarry(cpu.A, readR8())

	cpu.setZFlag(util.IsByteZeroUint8(sub))
	cpu.setNFlag(1)
	cpu.setHFlag(halfCarry)
	cpu.setCFlag(carry)
}
func (cpu *CPU) CP_A_B() {
	cpu.CP_A_R8(cpu.readB)
}
func (cpu *CPU) CP_A_C() {
	cpu.CP_A_R8(cpu.readC)
}
func (cpu *CPU) CP_A_D() {
	cpu.CP_A_R8(cpu.readD)
}
func (cpu *CPU) CP_A_E() {
	cpu.CP_A_R8(cpu.readE)
}
func (cpu *CPU) CP_A_H() {
	cpu.CP_A_R8(cpu.readH)
}
func (cpu *CPU) CP_A_L() {
	cpu.CP_A_R8(cpu.readL)
}
func (cpu *CPU) CP_A_HLmem() {
	cpu.CP_A_R8(cpu.readHLmem)
}
func (cpu *CPU) CP_A_A() {
	cpu.CP_A_R8(cpu.readA)
}

// ADD A R8
func (cpu *CPU) ADD_A_N8() {
	cpu.ADD_A_R8(cpu.ReadNextByte)
}

// ADC A R8
func (cpu *CPU) ADC_A_N8() {
	cpu.ADC_A_R8(cpu.ReadNextByte)
}

// SUB A R8
func (cpu *CPU) SUB_A_N8() {
	cpu.SUB_A_R8(cpu.ReadNextByte)
}

// SBC A R8
func (cpu *CPU) SBC_A_N8() {
	cpu.SBC_A_R8(cpu.ReadNextByte)
}

// AND A R8
func (cpu *CPU) AND_A_N8() {
	cpu.AND_A_R8(cpu.ReadNextByte)
}

// XOR A R8
func (cpu *CPU) XOR_A_N8() {
	cpu.XOR_A_R8(cpu.ReadNextByte)
}

// OR A R8
func (cpu *CPU) OR_A_N8() {
	cpu.OR_A_R8(cpu.ReadNextByte)
}

// CP A R8
func (cpu *CPU) CP_A_N8() {
	cpu.CP_A_R8(cpu.ReadNextByte)
}

// POP R16stk
func (cpu *CPU) POP_STACK() uint16 {
	cpu.SP += 2
	return cpu.ReadWord(cpu.SP - 2)
}
func (cpu *CPU) POP_BC() {
	cpu.writeBC(cpu.POP_STACK())
}
func (cpu *CPU) POP_DE() {
	cpu.writeDE(cpu.POP_STACK())
}
func (cpu *CPU) POP_HL() {
	cpu.writeHL(cpu.POP_STACK())
}
func (cpu *CPU) POP_AF() {
	cpu.writeAF(cpu.POP_STACK())
}

// PUSH R16STK
func (cpu *CPU) PUSH_STACK(v uint16) {
	cpu.SP -= 2
	cpu.Cycle() // Internal
	cpu.WriteWord(cpu.SP, v)
}
func (cpu *CPU) PUSH_BC() {
	cpu.PUSH_STACK(cpu.readBC())
}
func (cpu *CPU) PUSH_DE() {
	cpu.PUSH_STACK(cpu.readDE())
}
func (cpu *CPU) PUSH_HL() {
	cpu.PUSH_STACK(cpu.readHL())
}
func (cpu *CPU) PUSH_AF() {
	cpu.PUSH_STACK(cpu.readAF())
}

// RET COND
func (cpu *CPU) RET_COND(checkCondition func() bool) {
	cpu.Cycle() // Internal (branch decision)
	if checkCondition() {
		cpu.PC = cpu.POP_STACK()
		cpu.Cycle() // Internal (set PC)
	}
}
func (cpu *CPU) RET_Z() {
	cpu.RET_COND(func() bool { return cpu.readZFlag() == 1 })
}
func (cpu *CPU) RET_NZ() {
	cpu.RET_COND(func() bool { return cpu.readZFlag() == 0 })
}
func (cpu *CPU) RET_C() {
	cpu.RET_COND(func() bool { return cpu.readCFlag() == 1 })
}
func (cpu *CPU) RET_NC() {
	cpu.RET_COND(func() bool { return cpu.readCFlag() == 0 })
}

// RET
func (cpu *CPU) RET() {
	cpu.PC = cpu.POP_STACK()
	cpu.Cycle() // Internal (set PC)
}

// RETI
func (cpu *CPU) RETI() {
	cpu.PC = cpu.POP_STACK()
	cpu.Cycle() // Internal (set PC)
	cpu.IME = true
}

// JP COND
func (cpu *CPU) JP_COND_N16(checkCondition func() bool) {
	addr := cpu.ReadNextWord()
	if checkCondition() {
		cpu.PC = addr
		cpu.Cycle() // Internal (set PC)
	}
}
func (cpu *CPU) JP_Z_N16() {
	cpu.JP_COND_N16(func() bool { return cpu.readZFlag() == 1 })
}
func (cpu *CPU) JP_NZ_N16() {
	cpu.JP_COND_N16(func() bool { return cpu.readZFlag() == 0 })
}
func (cpu *CPU) JP_C_N16() {
	cpu.JP_COND_N16(func() bool { return cpu.readCFlag() == 1 })
}
func (cpu *CPU) JP_NC_N16() {
	cpu.JP_COND_N16(func() bool { return cpu.readCFlag() == 0 })
}

// JP N16
func (cpu *CPU) JP_N16() {
	cpu.JP_COND_N16(func() bool { return true })
}

// JP HL
func (cpu *CPU) JP_HL() {
	cpu.PC = cpu.readHL()
}

// CALL COND N16
func (cpu *CPU) CALL_COND_N16(checkCondition func() bool) {
	addr := cpu.ReadNextWord()
	if checkCondition() {
		cpu.PUSH_STACK(cpu.PC)
		cpu.PC = addr
	}
}
func (cpu *CPU) CALL_Z_N16() {
	cpu.CALL_COND_N16(func() bool { return cpu.readZFlag() == 1 })
}
func (cpu *CPU) CALL_NZ_N16() {
	cpu.CALL_COND_N16(func() bool { return cpu.readZFlag() == 0 })
}
func (cpu *CPU) CALL_C_N16() {
	cpu.CALL_COND_N16(func() bool { return cpu.readCFlag() == 1 })
}
func (cpu *CPU) CALL_NC_N16() {
	cpu.CALL_COND_N16(func() bool { return cpu.readCFlag() == 0 })
}

// CALL N16
func (cpu *CPU) CALL_N16() {
	cpu.CALL_COND_N16(func() bool { return true })
}

// RST VEC
func (cpu *CPU) RST_VEC(addr uint16) {
	cpu.PUSH_STACK(cpu.PC)
	cpu.PC = addr
}
func (cpu *CPU) RST_00() {
	cpu.RST_VEC(0x00)
}
func (cpu *CPU) RST_08() {
	cpu.RST_VEC(0x08)
}
func (cpu *CPU) RST_10() {
	cpu.RST_VEC(0x10)
}
func (cpu *CPU) RST_18() {
	cpu.RST_VEC(0x18)
}
func (cpu *CPU) RST_20() {
	cpu.RST_VEC(0x20)
}
func (cpu *CPU) RST_28() {
	cpu.RST_VEC(0x28)
}
func (cpu *CPU) RST_30() {
	cpu.RST_VEC(0x30)
}
func (cpu *CPU) RST_38() {
	cpu.RST_VEC(0x38)
}

// LDH_C_A
func (cpu *CPU) LDH_C_A() {
	cpu.WriteByte(0xFF00+uint16(cpu.C), cpu.A)
}
func (cpu *CPU) LDH_A_C() {
	cpu.writeA(cpu.ReadByte(0xFF00 + uint16(cpu.C)))
}

// LDH_N8_A
func (cpu *CPU) LDH_N8_A() {
	offset := cpu.ReadNextByte()
	cpu.WriteByte(0xFF00+uint16(offset), cpu.A)
}
func (cpu *CPU) LDH_A_N8() {
	offset := cpu.ReadNextByte()
	cpu.writeA(cpu.ReadByte(0xFF00 + uint16(offset)))
}

// LDH_N8_A
func (cpu *CPU) LD_N16_A() {
	addr := cpu.ReadNextWord()
	cpu.WriteByte(addr, cpu.A)
}
func (cpu *CPU) LD_A_N16() {
	addr := cpu.ReadNextWord()
	cpu.writeA(cpu.ReadByte(addr))
}

// ADD SP E8
func (cpu *CPU) SUM_SP_E8() uint16 {
	e8 := cpu.ReadNextByte()

	_, carry, halfCarry := util.SumBytesWithCarry(uint8(cpu.SP), e8)
	cpu.setZFlag(0)
	cpu.setNFlag(0)
	cpu.setHFlag(halfCarry)
	cpu.setCFlag(carry)

	return uint16(int(cpu.SP) + int(int8(e8)))
}
func (cpu *CPU) ADD_SP_E8() {
	high, low := util.SplitWord(cpu.SUM_SP_E8())

	cpu.writeLowSP(low)
	cpu.Cycle()
	cpu.writeHighSP(high)
	cpu.Cycle()
}

// LD HL SP+E8
func (cpu *CPU) LD_HL_SP_E8() {
	sum := cpu.SUM_SP_E8()
	cpu.writeHL(sum)
	cpu.Cycle() // Internal
}

// LD SP HL
func (cpu *CPU) LD_SP_HL() {
	cpu.SP = cpu.readHL()
	cpu.Cycle() // Internal
}

// DI
func (cpu *CPU) DI() {
	cpu.IME = false
}

// EI
func (cpu *CPU) EI() {
	cpu._EIDelayed = true
}

func (cpu *CPU) readR8(opcode uint8) uint8 {
	switch opcode & 0x07 {
	case 0:
		return cpu.readB()
	case 1:
		return cpu.readC()
	case 2:
		return cpu.readD()
	case 3:
		return cpu.readE()
	case 4:
		return cpu.readH()
	case 5:
		return cpu.readL()
	case 6:
		return cpu.readHLmem()
	case 7:
		return cpu.readA()
	}
	return 0xFF
}
func (cpu *CPU) writeR8(opcode uint8, value uint8) {
	switch opcode & 0x07 {
	case 0:
		cpu.writeB(value)
	case 1:
		cpu.writeC(value)
	case 2:
		cpu.writeD(value)
	case 3:
		cpu.writeE(value)
	case 4:
		cpu.writeH(value)
	case 5:
		cpu.writeL(value)
	case 6:
		cpu.writeHLmem(value)
	case 7:
		cpu.writeA(value)
	}
}

// RLC R8
func (cpu *CPU) RLC_R8(opcode uint8) {
	r8 := cpu.readR8(opcode)
	cFlag := r8 >> 7
	newR8 := (r8 << 1) | cFlag

	cpu.setZFlag(util.IsByteZeroUint8(newR8))
	cpu.setNFlag(0)
	cpu.setHFlag(0)
	cpu.setCFlag(cFlag)
	cpu.writeR8(opcode, newR8)
}

// RRC R8
func (cpu *CPU) RRC_R8(opcode uint8) {
	r8 := cpu.readR8(opcode)
	cFlag := r8 & 1
	newR8 := (r8 >> 1) | (cFlag << 7)

	cpu.setZFlag(util.IsByteZeroUint8(newR8))
	cpu.setNFlag(0)
	cpu.setHFlag(0)
	cpu.setCFlag(cFlag)
	cpu.writeR8(opcode, newR8)
}

// RL R8
func (cpu *CPU) RL_R8(opcode uint8) {
	r8 := cpu.readR8(opcode)
	cFlag := r8 >> 7
	newR8 := (r8 << 1) | cpu.readCFlag()

	cpu.setZFlag(util.IsByteZeroUint8(newR8))
	cpu.setNFlag(0)
	cpu.setHFlag(0)
	cpu.setCFlag(cFlag)
	cpu.writeR8(opcode, newR8)
}

// RR R8
func (cpu *CPU) RR_R8(opcode uint8) {
	r8 := cpu.readR8(opcode)
	cFlag := r8 & 1
	newR8 := (r8 >> 1) | (cpu.readCFlag() << 7)

	cpu.setZFlag(util.IsByteZeroUint8(newR8))
	cpu.setNFlag(0)
	cpu.setHFlag(0)
	cpu.setCFlag(cFlag)
	cpu.writeR8(opcode, newR8)
}

// SLA R8
func (cpu *CPU) SLA_R8(opcode uint8) {
	r8 := cpu.readR8(opcode)
	cFlag := r8 >> 7
	newR8 := r8 << 1

	cpu.setZFlag(util.IsByteZeroUint8(newR8))
	cpu.setNFlag(0)
	cpu.setHFlag(0)
	cpu.setCFlag(cFlag)
	cpu.writeR8(opcode, newR8)
}

// SRA R8
func (cpu *CPU) SRA_R8(opcode uint8) {
	r8 := cpu.readR8(opcode)
	cFlag := r8 & 1
	newR8 := (r8 >> 1) | (r8 & 0x80)

	cpu.setZFlag(util.IsByteZeroUint8(newR8))
	cpu.setNFlag(0)
	cpu.setHFlag(0)
	cpu.setCFlag(cFlag)
	cpu.writeR8(opcode, newR8)
}

// SWAP R8
func (cpu *CPU) SWAP_R8(opcode uint8) {
	r8 := cpu.readR8(opcode)
	newR8 := ((r8 & 0x0F) << 4) | ((r8 & 0xF0) >> 4)

	cpu.setZFlag(util.IsByteZeroUint8(newR8))
	cpu.setNFlag(0)
	cpu.setHFlag(0)
	cpu.setCFlag(0)
	cpu.writeR8(opcode, newR8)
}

// SRL R8
func (cpu *CPU) SRL_R8(opcode uint8) {
	r8 := cpu.readR8(opcode)
	cFlag := r8 & 1
	newR8 := r8 >> 1

	cpu.setZFlag(util.IsByteZeroUint8(newR8))
	cpu.setNFlag(0)
	cpu.setHFlag(0)
	cpu.setCFlag(cFlag)
	cpu.writeR8(opcode, newR8)
}

func (cpu *CPU) BIT_B3_R8(bit uint8, opcode uint8) {
	isSet := util.ReadBit(cpu.readR8(opcode), bit)

	cpu.setZFlag(1 - isSet)
	cpu.setNFlag(0)
	cpu.setHFlag(1)
}

func (cpu *CPU) RES_B3_R8(bit uint8, opcode uint8) {
	b := cpu.readR8(opcode)
	util.SetBit(&b, bit, 0)
	cpu.writeR8(opcode, b)
}

func (cpu *CPU) SET_B3_R8(bit uint8, opcode uint8) {
	b := cpu.readR8(opcode)
	util.SetBit(&b, bit, 1)
	cpu.writeR8(opcode, b)
}
