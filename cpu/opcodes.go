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
