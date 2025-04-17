package cpu

import "fmt"

type CPU struct {
	// 8-bit registers
	A, F uint8
	B, C uint8
	D, E uint8
	H, L uint8

	// 16-bit registers
	SP uint16
	PC uint16

	// Clock cycles
	cycles int

	// Memory
	Mem Memory

	// Flag set by opcode handler to compute correct number of cycles
	branched bool
}

type Memory interface {
	Read(addr uint16) uint8
	Write(addr uint16, value uint8)
	ReadWord(addr uint16) uint16
	WriteWord(addr uint16, value uint16)
}

func (cpu *CPU) ExecuteInstruction() {
	cpu.branched = false

	opcode := cpu.ReadNextByte()
	switch opcode {
	// NOP
	case NOP_OPCODE:
		cpu.NOP()
	// LD_R16_N16
	case LD_BC_N16_OPCODE:
		cpu.LD_BC_N16()
	case LD_DE_N16_OPCODE:
		cpu.LD_DE_N16()
	case LD_HL_N16_OPCODE:
		cpu.LD_HL_N16()
	case LD_SP_N16_OPCODE:
		cpu.LD_SP_N16()
	// LD_R16MEM_A
	case LD_BCMEM_A_OPCODE:
		cpu.LD_BCMEM_A()
	case LD_DEMEM_A_OPCODE:
		cpu.LD_DEMEM_A()
	case LD_HLIMEM_A_OPCODE:
		cpu.LD_HLIMEM_A()
	case LD_HLDMEM_A_OPCODE:
		cpu.LD_HLDMEM_A()
	// LD_A_R16MEM
	case LD_A_BCMEM_OPCODE:
		cpu.LD_A_BCMEM()
	case LD_A_DEMEM_OPCODE:
		cpu.LD_A_DEMEM()
	case LD_A_HLIMEM_OPCODE:
		cpu.LD_A_HLIMEM()
	case LD_A_HLDMEM_OPCODE:
		cpu.LD_A_HLDMEM()
	// LD_N16_SP
	case LD_N16_SP_OPCODE:
		cpu.LD_N16_SP()
	// INC_R16
	case INC_BC_OPCODE:
		cpu.INC_BC()
	case INC_DE_OPCODE:
		cpu.INC_DE()
	case INC_HL_OPCODE:
		cpu.INC_HL()
	case INC_SP_OPCODE:
		cpu.INC_SP()
	// DEC_R16
	case DEC_BC_OPCODE:
		cpu.DEC_BC()
	case DEC_DE_OPCODE:
		cpu.DEC_DE()
	case DEC_HL_OPCODE:
		cpu.DEC_HL()
	case DEC_SP_OPCODE:
		cpu.DEC_SP()
	// ADD_HL_R16
	case ADD_HL_BC_OPCODE:
		cpu.ADD_HL_BC()
	case ADD_HL_DE_OPCODE:
		cpu.ADD_HL_DE()
	case ADD_HL_HL_OPCODE:
		cpu.ADD_HL_HL()
	case ADD_HL_SP_OPCODE:
		cpu.ADD_HL_SP()
	// INC_R8
	case INC_B_OPCODE:
		cpu.INC_B()
	case INC_C_OPCODE:
		cpu.INC_C()
	case INC_D_OPCODE:
		cpu.INC_D()
	case INC_E_OPCODE:
		cpu.INC_E()
	case INC_H_OPCODE:
		cpu.INC_H()
	case INC_L_OPCODE:
		cpu.INC_L()
	case INC_HLMEM_OPCODE:
		cpu.INC_HLMEM()
	case INC_A_OPCODE:
		cpu.INC_A()
	// DEC_R8
	case DEC_B_OPCODE:
		cpu.DEC_B()
	case DEC_C_OPCODE:
		cpu.DEC_C()
	case DEC_D_OPCODE:
		cpu.DEC_D()
	case DEC_E_OPCODE:
		cpu.DEC_E()
	case DEC_H_OPCODE:
		cpu.DEC_H()
	case DEC_L_OPCODE:
		cpu.DEC_L()
	case DEC_HLMEM_OPCODE:
		cpu.DEC_HLMEM()
	case DEC_A_OPCODE:
		cpu.DEC_A()
	// LD_R8_N8
	case LD_B_N8_OPCODE:
		cpu.LD_B_N8()
	case LD_C_N8_OPCODE:
		cpu.LD_C_N8()
	case LD_D_N8_OPCODE:
		cpu.LD_D_N8()
	case LD_E_N8_OPCODE:
		cpu.LD_E_N8()
	case LD_H_N8_OPCODE:
		cpu.LD_H_N8()
	case LD_L_N8_OPCODE:
		cpu.LD_L_N8()
	case LD_HLMEM_N8_OPCODE:
		cpu.LD_HLMEM_N8()
	case LD_A_N8_OPCODE:
		cpu.LD_A_N8()
	// 8-bit logic
	case RLCA_OPCODE:
		cpu.RLCA()
	case RRCA_OPCODE:
		cpu.RRCA()
	case RLA_OPCODE:
		cpu.RLA()
	case RRA_OPCODE:
		cpu.RRA()
	// 8-bit arithmetic
	case DAA_OPCODE:
		cpu.DAA()
	case CPL_OPCODE:
		cpu.CPL()
	case SCF_OPCODE:
		cpu.SCF()
	case CCF_OPCODE:
		cpu.CCF()
	// JR_E8
	case JR_E8_OPCODE:
		cpu.JR_E8()
	// JR_COND_E8
	case JR_NZ_E8_OPCODE:
		cpu.JR_NZ_E8()
	case JR_Z_E8_OPCODE:
		cpu.JR_Z_E8()
	case JR_NC_E8_OPCODE:
		cpu.JR_NC_E8()
	case JR_C_E8_OPCODE:
		cpu.JR_C_E8()
	// STOP
	case STOP_OPCODE:
		cpu.STOP()
	default:
		fmt.Printf("OPCODE 0x%02X NOT RECOGNIZED\n", opcode)
	}

	if cpu.branched {
		cpu.cycles += OPCODES_CYCLES_BRANCH[opcode]
	} else {
		cpu.cycles += OPCODES_CYCLES[opcode]
	}
}

func (cpu *CPU) ReadNextByte() uint8 {
	b := cpu.Mem.Read(cpu.PC)
	cpu.PC++
	return b
}
func (cpu *CPU) WriteNextByte(b uint8) {
	cpu.Mem.Write(cpu.PC, b)
	cpu.PC++
}

func (cpu *CPU) ReadNextWord() uint16 {
	w := cpu.Mem.ReadWord(cpu.PC)
	cpu.PC += 2
	return w
}
func (cpu *CPU) WriteNextWord(w uint16) {
	cpu.Mem.WriteWord(cpu.PC, w)
	cpu.PC += 2
}

func (cpu *CPU) String() string {
	return fmt.Sprintf("A: %02X, F: %02X, B: %02X, C: %02X, D: %02X, E: %02X, H: %02X, L: %02X, SP: %04X, PC: %04X",
		cpu.A, cpu.F, cpu.B, cpu.C, cpu.D, cpu.E, cpu.H, cpu.L, cpu.SP, cpu.PC)
}
