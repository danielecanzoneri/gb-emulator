package cpu

// NOP
func (cpu *CPU) NOP() {
}

// LD_R16_N16
func (cpu *CPU) LD_BC_N16() {
	cpu.C = cpu.ReadNextByte()
	cpu.B = cpu.ReadNextByte()
}
func (cpu *CPU) LD_DE_N16() {
	cpu.E = cpu.ReadNextByte()
	cpu.D = cpu.ReadNextByte()
}
func (cpu *CPU) LD_HL_N16() {
	cpu.L = cpu.ReadNextByte()
	cpu.H = cpu.ReadNextByte()
}
func (cpu *CPU) LD_SP_N16() {
	cpu.SP = cpu.ReadNextWord()
}

// LD_R16MEM_A
func (cpu *CPU) LD_BCMEM_A() {
	cpu.Mem.Write(cpu.readBC(), cpu.A)
}
func (cpu *CPU) LD_DEMEM_A() {
	cpu.Mem.Write(cpu.readDE(), cpu.A)
}
func (cpu *CPU) LD_HLIMEM_A() {
	cpu.Mem.Write(cpu.readHL(), cpu.A)
	cpu.writeHL(cpu.readHL() + 1)
}
func (cpu *CPU) LD_HLDMEM_A() {
	cpu.Mem.Write(cpu.readHL(), cpu.A)
	cpu.writeHL(cpu.readHL() - 1)
}

// LD_A_R16MEM
func (cpu *CPU) LD_A_BCMEM() {
	cpu.A = cpu.Mem.Read(cpu.readBC())
}
func (cpu *CPU) LD_A_DEMEM() {
	cpu.A = cpu.Mem.Read(cpu.readDE())
}
func (cpu *CPU) LD_A_HLIMEM() {
	cpu.A = cpu.Mem.Read(cpu.readHL())
	cpu.writeHL(cpu.readHL() + 1)
}
func (cpu *CPU) LD_A_HLDMEM() {
	cpu.A = cpu.Mem.Read(cpu.readHL())
	cpu.writeHL(cpu.readHL() - 1)
}

// LD_N16_SP
func (cpu *CPU) LD_N16_SP() {
	cpu.WriteNextWord(cpu.SP)
}

// INC_R16
func (cpu *CPU) INC_BC() {
	cpu.writeBC(cpu.readBC() + 1)
}
func (cpu *CPU) INC_DE() {
	cpu.writeDE(cpu.readDE() + 1)
}
func (cpu *CPU) INC_HL() {
	cpu.writeHL(cpu.readHL() + 1)
}
func (cpu *CPU) INC_SP() {
	cpu.SP++
}

// DEC_R16
func (cpu *CPU) DEC_BC() {
	cpu.writeBC(cpu.readBC() - 1)
}
func (cpu *CPU) DEC_DE() {
	cpu.writeDE(cpu.readDE() - 1)
}
func (cpu *CPU) DEC_HL() {
	cpu.writeHL(cpu.readHL() - 1)
}
func (cpu *CPU) DEC_SP() {
	cpu.SP--
}

// ADD_HL_R16
func (cpu *CPU) ADD_HL_R16(r16 uint16) {
	sum, carry, half_carry := sumWordsWithCarry(cpu.readHL(), r16)
	cpu.setHFlag(half_carry)
	cpu.setCFlag(carry)
	cpu.writeHL(sum)
}
func (cpu *CPU) ADD_HL_BC() {
	cpu.ADD_HL_R16(cpu.readBC())
}
func (cpu *CPU) ADD_HL_DE() {
	cpu.ADD_HL_R16(cpu.readDE())
}
func (cpu *CPU) ADD_HL_HL() {
	cpu.ADD_HL_R16(cpu.readHL())
}
func (cpu *CPU) ADD_HL_SP() {
	cpu.ADD_HL_R16(cpu.SP)
}

// INC_R8
func (cpu *CPU) INC_R8(r8 uint8) uint8 {
	// Increments r8 and set correct flags
	sum, carry, half_carry := sumBytesWithCarry(r8, 1)
	cpu.setNFlag(0)
	cpu.setHFlag(half_carry)
	cpu.setZFlag(carry) // if carry it means result is 0
	return sum
}
func (cpu *CPU) INC_B() {
	cpu.B = cpu.INC_R8(cpu.B)
}
func (cpu *CPU) INC_C() {
	cpu.C = cpu.INC_R8(cpu.C)
}
func (cpu *CPU) INC_D() {
	cpu.D = cpu.INC_R8(cpu.D)
}
func (cpu *CPU) INC_E() {
	cpu.E = cpu.INC_R8(cpu.E)
}
func (cpu *CPU) INC_H() {
	cpu.H = cpu.INC_R8(cpu.H)
}
func (cpu *CPU) INC_L() {
	cpu.L = cpu.INC_R8(cpu.L)
}
func (cpu *CPU) INC_HLMEM() {
	addr := cpu.readHL()
	hl_mem := cpu.Mem.Read(addr)
	inc := cpu.INC_R8(hl_mem)
	cpu.Mem.Write(addr, inc)
}
func (cpu *CPU) INC_A() {
	cpu.A = cpu.INC_R8(cpu.A)
}

// DEC_R8
func (cpu *CPU) DEC_R8(r8 uint8) uint8 {
	// Decrements r8 and set correct flags
	sub, _, half_carry := subBytesWithCarry(r8, 1)
	cpu.setNFlag(1)
	cpu.setHFlag(half_carry)
	cpu.setZFlag(isByteZeroUint8(sub))
	return sub
}
func (cpu *CPU) DEC_B() {
	cpu.B = cpu.DEC_R8(cpu.B)
}
func (cpu *CPU) DEC_C() {
	cpu.C = cpu.DEC_R8(cpu.C)
}
func (cpu *CPU) DEC_D() {
	cpu.D = cpu.DEC_R8(cpu.D)
}
func (cpu *CPU) DEC_E() {
	cpu.E = cpu.DEC_R8(cpu.E)
}
func (cpu *CPU) DEC_H() {
	cpu.H = cpu.DEC_R8(cpu.H)
}
func (cpu *CPU) DEC_L() {
	cpu.L = cpu.DEC_R8(cpu.L)
}
func (cpu *CPU) DEC_HLMEM() {
	addr := cpu.readHL()
	hl_mem := cpu.Mem.Read(addr)
	dec := cpu.DEC_R8(hl_mem)
	cpu.Mem.Write(addr, dec)
}
func (cpu *CPU) DEC_A() {
	cpu.A = cpu.DEC_R8(cpu.A)
}

// LD_R8_N8
func (cpu *CPU) LD_B_N8() {
	cpu.B = cpu.ReadNextByte()
}
func (cpu *CPU) LD_C_N8() {
	cpu.C = cpu.ReadNextByte()
}
func (cpu *CPU) LD_D_N8() {
	cpu.D = cpu.ReadNextByte()
}
func (cpu *CPU) LD_E_N8() {
	cpu.E = cpu.ReadNextByte()
}
func (cpu *CPU) LD_H_N8() {
	cpu.H = cpu.ReadNextByte()
}
func (cpu *CPU) LD_L_N8() {
	cpu.L = cpu.ReadNextByte()
}
func (cpu *CPU) LD_HLMEM_N8() {
	addr := cpu.readHL()
	cpu.Mem.Write(addr, cpu.ReadNextByte())
}
func (cpu *CPU) LD_A_N8() {
	cpu.A = cpu.ReadNextByte()
}

// 8-bit logic
func (cpu *CPU) RLCA() {
	C_flag := cpu.A >> 7
	cpu.A = (cpu.A << 1) | C_flag
	cpu.setCFlag(C_flag)
}
func (cpu *CPU) RRCA() {
	C_flag := cpu.A & 0x01
	cpu.A = (cpu.A >> 1) | (C_flag << 7)
	cpu.setCFlag(C_flag)
}
func (cpu *CPU) RLA() {
	new_A := (cpu.A << 1) | cpu.readCFlag()
	C_flag := cpu.A >> 7
	cpu.A = new_A
	cpu.setCFlag(C_flag)
}
func (cpu *CPU) RRA() {
	new_A := (cpu.A >> 1) | (cpu.readCFlag() << 7)
	C_flag := cpu.A & 0x01
	cpu.A = new_A
	cpu.setCFlag(C_flag)
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
	cpu.setZFlag(isByteZeroUint8(cpu.A))
}
func (cpu *CPU) CPL() {
	cpu.A = ^cpu.A
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

// JR_E8
func (cpu *CPU) JR_E8() {
	e8 := cpu.ReadNextByte()
	cpu.PC = uint16(int(cpu.PC) + int(int8(e8)))
}

// JR_COND_E8
func (cpu *CPU) JR_Z_E8() {
	e8 := cpu.ReadNextByte()
	if cpu.readZFlag() == 1 {
		cpu.PC = uint16(int(cpu.PC) + int(int8(e8)))
		cpu.branched = true
	}
}
func (cpu *CPU) JR_NZ_E8() {
	e8 := cpu.ReadNextByte()
	if cpu.readZFlag() == 0 {
		cpu.PC = uint16(int(cpu.PC) + int(int8(e8)))
		cpu.branched = true
	}
}
func (cpu *CPU) JR_C_E8() {
	e8 := cpu.ReadNextByte()
	if cpu.readCFlag() == 1 {
		cpu.PC = uint16(int(cpu.PC) + int(int8(e8)))
		cpu.branched = true
	}
}
func (cpu *CPU) JR_NC_E8() {
	e8 := cpu.ReadNextByte()
	if cpu.readCFlag() == 0 {
		cpu.PC = uint16(int(cpu.PC) + int(int8(e8)))
		cpu.branched = true
	}
}

// STOP
func (cpu *CPU) STOP() {
	cpu.PC++
}

// LD R8 R8
func (cpu *CPU) LD_B_B() {
}
func (cpu *CPU) LD_B_C() {
	cpu.B = cpu.C
}
func (cpu *CPU) LD_B_D() {
	cpu.B = cpu.D
}
func (cpu *CPU) LD_B_E() {
	cpu.B = cpu.E
}
func (cpu *CPU) LD_B_H() {
	cpu.B = cpu.H
}
func (cpu *CPU) LD_B_L() {
	cpu.B = cpu.L
}
func (cpu *CPU) LD_B_HLMEM() {
	cpu.B = cpu.Mem.Read(cpu.readHL())
}
func (cpu *CPU) LD_B_A() {
	cpu.B = cpu.A
}

func (cpu *CPU) LD_C_B() {
	cpu.C = cpu.B
}
func (cpu *CPU) LD_C_C() {
}
func (cpu *CPU) LD_C_D() {
	cpu.C = cpu.D
}
func (cpu *CPU) LD_C_E() {
	cpu.C = cpu.E
}
func (cpu *CPU) LD_C_H() {
	cpu.C = cpu.H
}
func (cpu *CPU) LD_C_L() {
	cpu.C = cpu.L
}
func (cpu *CPU) LD_C_HLMEM() {
	cpu.C = cpu.Mem.Read(cpu.readHL())
}
func (cpu *CPU) LD_C_A() {
	cpu.C = cpu.A
}

func (cpu *CPU) LD_D_B() {
	cpu.D = cpu.B
}
func (cpu *CPU) LD_D_C() {
	cpu.D = cpu.C
}
func (cpu *CPU) LD_D_D() {
}
func (cpu *CPU) LD_D_E() {
	cpu.D = cpu.E
}
func (cpu *CPU) LD_D_H() {
	cpu.D = cpu.H
}
func (cpu *CPU) LD_D_L() {
	cpu.D = cpu.L
}
func (cpu *CPU) LD_D_HLMEM() {
	cpu.D = cpu.Mem.Read(cpu.readHL())
}
func (cpu *CPU) LD_D_A() {
	cpu.D = cpu.A
}

func (cpu *CPU) LD_E_B() {
	cpu.E = cpu.B
}
func (cpu *CPU) LD_E_C() {
	cpu.E = cpu.C
}
func (cpu *CPU) LD_E_D() {
	cpu.E = cpu.D
}
func (cpu *CPU) LD_E_E() {
}
func (cpu *CPU) LD_E_H() {
	cpu.E = cpu.H
}
func (cpu *CPU) LD_E_L() {
	cpu.E = cpu.L
}
func (cpu *CPU) LD_E_HLMEM() {
	cpu.E = cpu.Mem.Read(cpu.readHL())
}
func (cpu *CPU) LD_E_A() {
	cpu.E = cpu.A
}

func (cpu *CPU) LD_H_B() {
	cpu.H = cpu.B
}
func (cpu *CPU) LD_H_C() {
	cpu.H = cpu.C
}
func (cpu *CPU) LD_H_D() {
	cpu.H = cpu.D
}
func (cpu *CPU) LD_H_E() {
	cpu.H = cpu.E
}
func (cpu *CPU) LD_H_H() {
}
func (cpu *CPU) LD_H_L() {
	cpu.H = cpu.L
}
func (cpu *CPU) LD_H_HLMEM() {
	cpu.H = cpu.Mem.Read(cpu.readHL())
}
func (cpu *CPU) LD_H_A() {
	cpu.H = cpu.A
}

func (cpu *CPU) LD_L_B() {
	cpu.L = cpu.B
}
func (cpu *CPU) LD_L_C() {
	cpu.L = cpu.C
}
func (cpu *CPU) LD_L_D() {
	cpu.L = cpu.D
}
func (cpu *CPU) LD_L_E() {
	cpu.L = cpu.E
}
func (cpu *CPU) LD_L_H() {
	cpu.L = cpu.H
}
func (cpu *CPU) LD_L_L() {
}
func (cpu *CPU) LD_L_HLMEM() {
	cpu.L = cpu.Mem.Read(cpu.readHL())
}
func (cpu *CPU) LD_L_A() {
	cpu.L = cpu.A
}

func (cpu *CPU) LD_HLMEM_B() {
	cpu.Mem.Write(cpu.readHL(), cpu.B)
}
func (cpu *CPU) LD_HLMEM_C() {
	cpu.Mem.Write(cpu.readHL(), cpu.C)
}
func (cpu *CPU) LD_HLMEM_D() {
	cpu.Mem.Write(cpu.readHL(), cpu.D)
}
func (cpu *CPU) LD_HLMEM_E() {
	cpu.Mem.Write(cpu.readHL(), cpu.E)
}
func (cpu *CPU) LD_HLMEM_H() {
	cpu.Mem.Write(cpu.readHL(), cpu.H)
}
func (cpu *CPU) LD_HLMEM_L() {
	cpu.Mem.Write(cpu.readHL(), cpu.L)
}
func (cpu *CPU) LD_HLMEM_A() {
	cpu.Mem.Write(cpu.readHL(), cpu.A)
}

func (cpu *CPU) LD_A_B() {
	cpu.A = cpu.B
}
func (cpu *CPU) LD_A_C() {
	cpu.A = cpu.C
}
func (cpu *CPU) LD_A_D() {
	cpu.A = cpu.D
}
func (cpu *CPU) LD_A_E() {
	cpu.A = cpu.E
}
func (cpu *CPU) LD_A_H() {
	cpu.A = cpu.H
}
func (cpu *CPU) LD_A_L() {
	cpu.A = cpu.L
}
func (cpu *CPU) LD_A_HLMEM() {
	cpu.A = cpu.Mem.Read(cpu.readHL())
}
func (cpu *CPU) LD_A_A() {
}

// HALT
func (cpu *CPU) HALT() {
}

// ADD A R8
func (cpu *CPU) ADD_A_R8(r8 uint8) {
	sum, carry, halfCarry := sumBytesWithCarry(cpu.A, r8)

	cpu.setZFlag(isByteZeroUint8(sum))
	cpu.setNFlag(0)
	cpu.setHFlag(halfCarry)
	cpu.setCFlag(carry)
	cpu.A = sum
}
func (cpu *CPU) ADD_A_B() {
	cpu.ADD_A_R8(cpu.B)
}
func (cpu *CPU) ADD_A_C() {
	cpu.ADD_A_R8(cpu.C)
}
func (cpu *CPU) ADD_A_D() {
	cpu.ADD_A_R8(cpu.D)
}
func (cpu *CPU) ADD_A_E() {
	cpu.ADD_A_R8(cpu.E)
}
func (cpu *CPU) ADD_A_H() {
	cpu.ADD_A_R8(cpu.H)
}
func (cpu *CPU) ADD_A_L() {
	cpu.ADD_A_R8(cpu.L)
}
func (cpu *CPU) ADD_A_HLMEM() {
	cpu.ADD_A_R8(cpu.Mem.Read(cpu.readHL()))
}
func (cpu *CPU) ADD_A_A() {
	cpu.ADD_A_R8(cpu.A)
}

// ADC A R8
func (cpu *CPU) ADC_A_R8(r8 uint8) {
	sum, carry1, halfCarry1 := sumBytesWithCarry(cpu.A, r8)
	sum, carry2, halfCarry2 := sumBytesWithCarry(sum, cpu.readCFlag())

	cpu.setZFlag(isByteZeroUint8(sum))
	cpu.setNFlag(0)
	cpu.setHFlag(halfCarry1 | halfCarry2)
	cpu.setCFlag(carry1 | carry2)
	cpu.A = sum
}
func (cpu *CPU) ADC_A_B() {
	cpu.ADC_A_R8(cpu.B)
}
func (cpu *CPU) ADC_A_C() {
	cpu.ADC_A_R8(cpu.C)
}
func (cpu *CPU) ADC_A_D() {
	cpu.ADC_A_R8(cpu.D)
}
func (cpu *CPU) ADC_A_E() {
	cpu.ADC_A_R8(cpu.E)
}
func (cpu *CPU) ADC_A_H() {
	cpu.ADC_A_R8(cpu.H)
}
func (cpu *CPU) ADC_A_L() {
	cpu.ADC_A_R8(cpu.L)
}
func (cpu *CPU) ADC_A_HLMEM() {
	cpu.ADC_A_R8(cpu.Mem.Read(cpu.readHL()))
}
func (cpu *CPU) ADC_A_A() {
	cpu.ADC_A_R8(cpu.A)
}

// SUB A R8
func (cpu *CPU) SUB_A_R8(r8 uint8) {
	sub, carry, halfCarry := subBytesWithCarry(cpu.A, r8)

	cpu.setZFlag(isByteZeroUint8(sub))
	cpu.setNFlag(1)
	cpu.setHFlag(halfCarry)
	cpu.setCFlag(carry)
	cpu.A = sub
}
func (cpu *CPU) SUB_A_B() {
	cpu.SUB_A_R8(cpu.B)
}
func (cpu *CPU) SUB_A_C() {
	cpu.SUB_A_R8(cpu.C)
}
func (cpu *CPU) SUB_A_D() {
	cpu.SUB_A_R8(cpu.D)
}
func (cpu *CPU) SUB_A_E() {
	cpu.SUB_A_R8(cpu.E)
}
func (cpu *CPU) SUB_A_H() {
	cpu.SUB_A_R8(cpu.H)
}
func (cpu *CPU) SUB_A_L() {
	cpu.SUB_A_R8(cpu.L)
}
func (cpu *CPU) SUB_A_HLMEM() {
	cpu.SUB_A_R8(cpu.Mem.Read(cpu.readHL()))
}
func (cpu *CPU) SUB_A_A() {
	cpu.SUB_A_R8(cpu.A)
}

// SBC A R8
func (cpu *CPU) SBC_A_R8(r8 uint8) {
	sub, carry1, halfCarry1 := subBytesWithCarry(cpu.A, r8)
	sub, carry2, halfCarry2 := subBytesWithCarry(sub, cpu.readCFlag())

	cpu.setZFlag(isByteZeroUint8(sub))
	cpu.setNFlag(1)
	cpu.setHFlag(halfCarry1 | halfCarry2)
	cpu.setCFlag(carry1 | carry2)
	cpu.A = sub
}
func (cpu *CPU) SBC_A_B() {
	cpu.SBC_A_R8(cpu.B)
}
func (cpu *CPU) SBC_A_C() {
	cpu.SBC_A_R8(cpu.C)
}
func (cpu *CPU) SBC_A_D() {
	cpu.SBC_A_R8(cpu.D)
}
func (cpu *CPU) SBC_A_E() {
	cpu.SBC_A_R8(cpu.E)
}
func (cpu *CPU) SBC_A_H() {
	cpu.SBC_A_R8(cpu.H)
}
func (cpu *CPU) SBC_A_L() {
	cpu.SBC_A_R8(cpu.L)
}
func (cpu *CPU) SBC_A_HLMEM() {
	cpu.SBC_A_R8(cpu.Mem.Read(cpu.readHL()))
}
func (cpu *CPU) SBC_A_A() {
	cpu.SBC_A_R8(cpu.A)
}

// AND A R8
func (cpu *CPU) AND_A_R8(r8 uint8) {
	cpu.A = cpu.A & r8
	cpu.setZFlag(isByteZeroUint8(cpu.A))
	cpu.setNFlag(0)
	cpu.setHFlag(1)
	cpu.setCFlag(0)
}
func (cpu *CPU) AND_A_B() {
	cpu.AND_A_R8(cpu.B)
}
func (cpu *CPU) AND_A_C() {
	cpu.AND_A_R8(cpu.C)
}
func (cpu *CPU) AND_A_D() {
	cpu.AND_A_R8(cpu.D)
}
func (cpu *CPU) AND_A_E() {
	cpu.AND_A_R8(cpu.E)
}
func (cpu *CPU) AND_A_H() {
	cpu.AND_A_R8(cpu.H)
}
func (cpu *CPU) AND_A_L() {
	cpu.AND_A_R8(cpu.L)
}
func (cpu *CPU) AND_A_HLMEM() {
	cpu.AND_A_R8(cpu.Mem.Read(cpu.readHL()))
}
func (cpu *CPU) AND_A_A() {
	cpu.AND_A_R8(cpu.A)
}

// XOR A R8
func (cpu *CPU) XOR_A_R8(r8 uint8) {
	cpu.A = cpu.A ^ r8
	cpu.setZFlag(isByteZeroUint8(cpu.A))
	cpu.setNFlag(0)
	cpu.setHFlag(0)
	cpu.setCFlag(0)
}
func (cpu *CPU) XOR_A_B() {
	cpu.XOR_A_R8(cpu.B)
}
func (cpu *CPU) XOR_A_C() {
	cpu.XOR_A_R8(cpu.C)
}
func (cpu *CPU) XOR_A_D() {
	cpu.XOR_A_R8(cpu.D)
}
func (cpu *CPU) XOR_A_E() {
	cpu.XOR_A_R8(cpu.E)
}
func (cpu *CPU) XOR_A_H() {
	cpu.XOR_A_R8(cpu.H)
}
func (cpu *CPU) XOR_A_L() {
	cpu.XOR_A_R8(cpu.L)
}
func (cpu *CPU) XOR_A_HLMEM() {
	cpu.XOR_A_R8(cpu.Mem.Read(cpu.readHL()))
}
func (cpu *CPU) XOR_A_A() {
	cpu.XOR_A_R8(cpu.A)
}

// OR A R8
func (cpu *CPU) OR_A_R8(r8 uint8) {
	cpu.A = cpu.A | r8
	cpu.setZFlag(isByteZeroUint8(cpu.A))
	cpu.setNFlag(0)
	cpu.setHFlag(0)
	cpu.setCFlag(0)
}
func (cpu *CPU) OR_A_B() {
	cpu.OR_A_R8(cpu.B)
}
func (cpu *CPU) OR_A_C() {
	cpu.OR_A_R8(cpu.C)
}
func (cpu *CPU) OR_A_D() {
	cpu.OR_A_R8(cpu.D)
}
func (cpu *CPU) OR_A_E() {
	cpu.OR_A_R8(cpu.E)
}
func (cpu *CPU) OR_A_H() {
	cpu.OR_A_R8(cpu.H)
}
func (cpu *CPU) OR_A_L() {
	cpu.OR_A_R8(cpu.L)
}
func (cpu *CPU) OR_A_HLMEM() {
	cpu.OR_A_R8(cpu.Mem.Read(cpu.readHL()))
}
func (cpu *CPU) OR_A_A() {
	cpu.OR_A_R8(cpu.A)
}

// CP A R8
func (cpu *CPU) CP_A_R8(r8 uint8) {
	sub, carry, halfCarry := subBytesWithCarry(cpu.A, r8)

	cpu.setZFlag(isByteZeroUint8(sub))
	cpu.setNFlag(1)
	cpu.setHFlag(halfCarry)
	cpu.setCFlag(carry)
}
func (cpu *CPU) CP_A_B() {
	cpu.CP_A_R8(cpu.B)
}
func (cpu *CPU) CP_A_C() {
	cpu.CP_A_R8(cpu.C)
}
func (cpu *CPU) CP_A_D() {
	cpu.CP_A_R8(cpu.D)
}
func (cpu *CPU) CP_A_E() {
	cpu.CP_A_R8(cpu.E)
}
func (cpu *CPU) CP_A_H() {
	cpu.CP_A_R8(cpu.H)
}
func (cpu *CPU) CP_A_L() {
	cpu.CP_A_R8(cpu.L)
}
func (cpu *CPU) CP_A_HLMEM() {
	cpu.CP_A_R8(cpu.Mem.Read(cpu.readHL()))
}
func (cpu *CPU) CP_A_A() {
	cpu.CP_A_R8(cpu.A)
}

// ADD A R8
func (cpu *CPU) ADD_A_N8() {
	n8 := cpu.ReadNextByte()
	cpu.ADD_A_R8(n8)
}

// ADC A R8
func (cpu *CPU) ADC_A_N8() {
	n8 := cpu.ReadNextByte()
	cpu.ADC_A_R8(n8)
}

// SUB A R8
func (cpu *CPU) SUB_A_N8() {
	n8 := cpu.ReadNextByte()
	cpu.SUB_A_R8(n8)
}

// SBC A R8
func (cpu *CPU) SBC_A_N8() {
	n8 := cpu.ReadNextByte()
	cpu.SBC_A_R8(n8)
}

// AND A R8
func (cpu *CPU) AND_A_N8() {
	n8 := cpu.ReadNextByte()
	cpu.AND_A_R8(n8)
}

// XOR A R8
func (cpu *CPU) XOR_A_N8() {
	n8 := cpu.ReadNextByte()
	cpu.XOR_A_R8(n8)
}

// OR A R8
func (cpu *CPU) OR_A_N8() {
	n8 := cpu.ReadNextByte()
	cpu.OR_A_R8(n8)
}

// CP A R8
func (cpu *CPU) CP_A_N8() {
	n8 := cpu.ReadNextByte()
	cpu.CP_A_R8(n8)
}
