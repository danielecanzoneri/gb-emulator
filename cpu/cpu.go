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

	// Interrupt Master Enable
	IME        bool
	_EIDelayed bool // Set to true when EI is executed, but not yet effective

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

func New(mem Memory) *CPU {
	cpu := &CPU{SP: 0xFFFE, Mem: mem}
	return cpu
}

func (cpu *CPU) ExecuteInstruction() {
	cpu.branched = false

	if cpu.handleInterrupts() {
		return
	}

	opcode := cpu.ReadNextByte()
	cpu.logState(opcode)

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
	// LD R8 R8
	case LD_B_B_OPCODE:
		cpu.LD_B_B()
	case LD_B_C_OPCODE:
		cpu.LD_B_C()
	case LD_B_D_OPCODE:
		cpu.LD_B_D()
	case LD_B_E_OPCODE:
		cpu.LD_B_E()
	case LD_B_H_OPCODE:
		cpu.LD_B_H()
	case LD_B_L_OPCODE:
		cpu.LD_B_L()
	case LD_B_HLMEM_OPCODE:
		cpu.LD_B_HLMEM()
	case LD_B_A_OPCODE:
		cpu.LD_B_A()
	case LD_C_B_OPCODE:
		cpu.LD_C_B()
	case LD_C_C_OPCODE:
		cpu.LD_C_C()
	case LD_C_D_OPCODE:
		cpu.LD_C_D()
	case LD_C_E_OPCODE:
		cpu.LD_C_E()
	case LD_C_H_OPCODE:
		cpu.LD_C_H()
	case LD_C_L_OPCODE:
		cpu.LD_C_L()
	case LD_C_HLMEM_OPCODE:
		cpu.LD_C_HLMEM()
	case LD_C_A_OPCODE:
		cpu.LD_C_A()
	case LD_D_B_OPCODE:
		cpu.LD_D_B()
	case LD_D_C_OPCODE:
		cpu.LD_D_C()
	case LD_D_D_OPCODE:
		cpu.LD_D_D()
	case LD_D_E_OPCODE:
		cpu.LD_D_E()
	case LD_D_H_OPCODE:
		cpu.LD_D_H()
	case LD_D_L_OPCODE:
		cpu.LD_D_L()
	case LD_D_HLMEM_OPCODE:
		cpu.LD_D_HLMEM()
	case LD_D_A_OPCODE:
		cpu.LD_D_A()
	case LD_E_B_OPCODE:
		cpu.LD_E_B()
	case LD_E_C_OPCODE:
		cpu.LD_E_C()
	case LD_E_D_OPCODE:
		cpu.LD_E_D()
	case LD_E_E_OPCODE:
		cpu.LD_E_E()
	case LD_E_H_OPCODE:
		cpu.LD_E_H()
	case LD_E_L_OPCODE:
		cpu.LD_E_L()
	case LD_E_HLMEM_OPCODE:
		cpu.LD_E_HLMEM()
	case LD_E_A_OPCODE:
		cpu.LD_E_A()
	case LD_H_B_OPCODE:
		cpu.LD_H_B()
	case LD_H_C_OPCODE:
		cpu.LD_H_C()
	case LD_H_D_OPCODE:
		cpu.LD_H_D()
	case LD_H_E_OPCODE:
		cpu.LD_H_E()
	case LD_H_H_OPCODE:
		cpu.LD_H_H()
	case LD_H_L_OPCODE:
		cpu.LD_H_L()
	case LD_H_HLMEM_OPCODE:
		cpu.LD_H_HLMEM()
	case LD_H_A_OPCODE:
		cpu.LD_H_A()
	case LD_L_B_OPCODE:
		cpu.LD_L_B()
	case LD_L_C_OPCODE:
		cpu.LD_L_C()
	case LD_L_D_OPCODE:
		cpu.LD_L_D()
	case LD_L_E_OPCODE:
		cpu.LD_L_E()
	case LD_L_H_OPCODE:
		cpu.LD_L_H()
	case LD_L_L_OPCODE:
		cpu.LD_L_L()
	case LD_L_HLMEM_OPCODE:
		cpu.LD_L_HLMEM()
	case LD_L_A_OPCODE:
		cpu.LD_L_A()
	case LD_HLMEM_B_OPCODE:
		cpu.LD_HLMEM_B()
	case LD_HLMEM_C_OPCODE:
		cpu.LD_HLMEM_C()
	case LD_HLMEM_D_OPCODE:
		cpu.LD_HLMEM_D()
	case LD_HLMEM_E_OPCODE:
		cpu.LD_HLMEM_E()
	case LD_HLMEM_H_OPCODE:
		cpu.LD_HLMEM_H()
	case LD_HLMEM_L_OPCODE:
		cpu.LD_HLMEM_L()
	case LD_HLMEM_A_OPCODE:
		cpu.LD_HLMEM_A()
	case LD_A_B_OPCODE:
		cpu.LD_A_B()
	case LD_A_C_OPCODE:
		cpu.LD_A_C()
	case LD_A_D_OPCODE:
		cpu.LD_A_D()
	case LD_A_E_OPCODE:
		cpu.LD_A_E()
	case LD_A_H_OPCODE:
		cpu.LD_A_H()
	case LD_A_L_OPCODE:
		cpu.LD_A_L()
	case LD_A_HLMEM_OPCODE:
		cpu.LD_A_HLMEM()
	case LD_A_A_OPCODE:
		cpu.LD_A_A()
	// ADD A R8
	case ADD_A_B_OPCODE:
		cpu.ADD_A_B()
	case ADD_A_C_OPCODE:
		cpu.ADD_A_C()
	case ADD_A_D_OPCODE:
		cpu.ADD_A_D()
	case ADD_A_E_OPCODE:
		cpu.ADD_A_E()
	case ADD_A_H_OPCODE:
		cpu.ADD_A_H()
	case ADD_A_L_OPCODE:
		cpu.ADD_A_L()
	case ADD_A_HLMEM_OPCODE:
		cpu.ADD_A_HLMEM()
	case ADD_A_A_OPCODE:
		cpu.ADD_A_A()
	// ADC A R8
	case ADC_A_B_OPCODE:
		cpu.ADC_A_B()
	case ADC_A_C_OPCODE:
		cpu.ADC_A_C()
	case ADC_A_D_OPCODE:
		cpu.ADC_A_D()
	case ADC_A_E_OPCODE:
		cpu.ADC_A_E()
	case ADC_A_H_OPCODE:
		cpu.ADC_A_H()
	case ADC_A_L_OPCODE:
		cpu.ADC_A_L()
	case ADC_A_HLMEM_OPCODE:
		cpu.ADC_A_HLMEM()
	case ADC_A_A_OPCODE:
		cpu.ADC_A_A()
	// SUB A R8
	case SUB_A_B_OPCODE:
		cpu.SUB_A_B()
	case SUB_A_C_OPCODE:
		cpu.SUB_A_C()
	case SUB_A_D_OPCODE:
		cpu.SUB_A_D()
	case SUB_A_E_OPCODE:
		cpu.SUB_A_E()
	case SUB_A_H_OPCODE:
		cpu.SUB_A_H()
	case SUB_A_L_OPCODE:
		cpu.SUB_A_L()
	case SUB_A_HLMEM_OPCODE:
		cpu.SUB_A_HLMEM()
	case SUB_A_A_OPCODE:
		cpu.SUB_A_A()
	// SBC A R8
	case SBC_A_B_OPCODE:
		cpu.SBC_A_B()
	case SBC_A_C_OPCODE:
		cpu.SBC_A_C()
	case SBC_A_D_OPCODE:
		cpu.SBC_A_D()
	case SBC_A_E_OPCODE:
		cpu.SBC_A_E()
	case SBC_A_H_OPCODE:
		cpu.SBC_A_H()
	case SBC_A_L_OPCODE:
		cpu.SBC_A_L()
	case SBC_A_HLMEM_OPCODE:
		cpu.SBC_A_HLMEM()
	case SBC_A_A_OPCODE:
		cpu.SBC_A_A()
	// AND A R8
	case AND_A_B_OPCODE:
		cpu.AND_A_B()
	case AND_A_C_OPCODE:
		cpu.AND_A_C()
	case AND_A_D_OPCODE:
		cpu.AND_A_D()
	case AND_A_E_OPCODE:
		cpu.AND_A_E()
	case AND_A_H_OPCODE:
		cpu.AND_A_H()
	case AND_A_L_OPCODE:
		cpu.AND_A_L()
	case AND_A_HLMEM_OPCODE:
		cpu.AND_A_HLMEM()
	case AND_A_A_OPCODE:
		cpu.AND_A_A()
	// XOR A R8
	case XOR_A_B_OPCODE:
		cpu.XOR_A_B()
	case XOR_A_C_OPCODE:
		cpu.XOR_A_C()
	case XOR_A_D_OPCODE:
		cpu.XOR_A_D()
	case XOR_A_E_OPCODE:
		cpu.XOR_A_E()
	case XOR_A_H_OPCODE:
		cpu.XOR_A_H()
	case XOR_A_L_OPCODE:
		cpu.XOR_A_L()
	case XOR_A_HLMEM_OPCODE:
		cpu.XOR_A_HLMEM()
	case XOR_A_A_OPCODE:
		cpu.XOR_A_A()
	// OR A R8
	case OR_A_B_OPCODE:
		cpu.OR_A_B()
	case OR_A_C_OPCODE:
		cpu.OR_A_C()
	case OR_A_D_OPCODE:
		cpu.OR_A_D()
	case OR_A_E_OPCODE:
		cpu.OR_A_E()
	case OR_A_H_OPCODE:
		cpu.OR_A_H()
	case OR_A_L_OPCODE:
		cpu.OR_A_L()
	case OR_A_HLMEM_OPCODE:
		cpu.OR_A_HLMEM()
	case OR_A_A_OPCODE:
		cpu.OR_A_A()
	// CP A R8
	case CP_A_B_OPCODE:
		cpu.CP_A_B()
	case CP_A_C_OPCODE:
		cpu.CP_A_C()
	case CP_A_D_OPCODE:
		cpu.CP_A_D()
	case CP_A_E_OPCODE:
		cpu.CP_A_E()
	case CP_A_H_OPCODE:
		cpu.CP_A_H()
	case CP_A_L_OPCODE:
		cpu.CP_A_L()
	case CP_A_HLMEM_OPCODE:
		cpu.CP_A_HLMEM()
	case CP_A_A_OPCODE:
		cpu.CP_A_A()
	// ADD A R8
	case ADD_A_N8_OPCODE:
		cpu.ADD_A_N8()
	// ADC A R8
	case ADC_A_N8_OPCODE:
		cpu.ADC_A_N8()
	// SUB A R8
	case SUB_A_N8_OPCODE:
		cpu.SUB_A_N8()
	// SBC A R8
	case SBC_A_N8_OPCODE:
		cpu.SBC_A_N8()
	// AND A R8
	case AND_A_N8_OPCODE:
		cpu.AND_A_N8()
	// XOR A R8
	case XOR_A_N8_OPCODE:
		cpu.XOR_A_N8()
	// OR A R8
	case OR_A_N8_OPCODE:
		cpu.OR_A_N8()
	// CP A R8
	case CP_A_N8_OPCODE:
		cpu.CP_A_N8()
	// POP R16STK
	case POP_BC_OPCODE:
		cpu.POP_BC()
	case POP_DE_OPCODE:
		cpu.POP_DE()
	case POP_HL_OPCODE:
		cpu.POP_HL()
	case POP_AF_OPCODE:
		cpu.POP_AF()
	// PUSH R16STK
	case PUSH_BC_OPCODE:
		cpu.PUSH_BC()
	case PUSH_DE_OPCODE:
		cpu.PUSH_DE()
	case PUSH_HL_OPCODE:
		cpu.PUSH_HL()
	case PUSH_AF_OPCODE:
		cpu.PUSH_AF()
	// RET COND
	case RET_NZ_OPCODE:
		cpu.RET_NZ()
	case RET_Z_OPCODE:
		cpu.RET_Z()
	case RET_NC_OPCODE:
		cpu.RET_NC()
	case RET_C_OPCODE:
		cpu.RET_C()
	// RET
	case RET_OPCODE:
		cpu.RET()
	// RETI
	case RETI_OPCODE:
		cpu.RETI()
	// JP COND N16
	case JP_NZ_N16_OPCODE:
		cpu.JP_NZ_N16()
	case JP_Z_N16_OPCODE:
		cpu.JP_Z_N16()
	case JP_NC_N16_OPCODE:
		cpu.JP_NC_N16()
	case JP_C_N16_OPCODE:
		cpu.JP_C_N16()
	// JP N16
	case JP_N16_OPCODE:
		cpu.JP_N16()
	// JP HL
	case JP_HL_OPCODE:
		cpu.JP_HL()
	// CALL COND N16
	case CALL_NZ_N16_OPCODE:
		cpu.CALL_NZ_N16()
	case CALL_Z_N16_OPCODE:
		cpu.CALL_Z_N16()
	case CALL_NC_N16_OPCODE:
		cpu.CALL_NC_N16()
	case CALL_C_N16_OPCODE:
		cpu.CALL_C_N16()
	// CALL N16
	case CALL_N16_OPCODE:
		cpu.CALL_N16()
	// RST VEC
	case RST_00_OPCODE:
		cpu.RST_00()
	case RST_08_OPCODE:
		cpu.RST_08()
	case RST_10_OPCODE:
		cpu.RST_10()
	case RST_18_OPCODE:
		cpu.RST_18()
	case RST_20_OPCODE:
		cpu.RST_20()
	case RST_28_OPCODE:
		cpu.RST_28()
	case RST_30_OPCODE:
		cpu.RST_30()
	case RST_38_OPCODE:
		cpu.RST_38()
	// LDH C A
	case LDH_C_A_OPCODE:
		cpu.LDH_C_A()
	// LDH A C
	case LDH_A_C_OPCODE:
		cpu.LDH_A_C()
	// LDH N _A
	case LDH_N8_A_OPCODE:
		cpu.LDH_N8_A()
	// LDH A N8
	case LDH_A_N8_OPCODE:
		cpu.LDH_A_N8()
	// LD N16 A
	case LD_N16_A_OPCODE:
		cpu.LD_N16_A()
	// LD A N16
	case LD_A_N16_OPCODE:
		cpu.LD_A_N16()
	// ADD SP E8
	case ADD_SP_E8_OPCODE:
		cpu.ADD_SP_E8()
	// LD HL SP+E8
	case LD_HL_SP_E8_OPCODE:
		cpu.LD_HL_SP_E8()
	// LD SP HL
	case LD_SP_HL_OPCODE:
		cpu.LD_SP_HL()
	// DI
	case DI_OPCODE:
		cpu.DI()
	// EI
	case EI_OPCODE:
		cpu.EI()
	// PREFIX
	case PREFIX_OPCODE:
		cpu.prefixedOpcode()
	default:
		fmt.Printf("OPCODE 0x%02X NOT RECOGNIZED\n", opcode)
	}

	if cpu.branched {
		cpu.cycles += OPCODES_CYCLES_BRANCH[opcode]
	} else {
		cpu.cycles += OPCODES_CYCLES[opcode]
	}
}

func (cpu *CPU) prefixedOpcode() {
	cpu.branched = false

	opcode := cpu.ReadNextByte()
	switch opcode & 0b11111000 {
	case RLC_R8_OPCODE:
		cpu.RLC_R8(opcode)
	case RRC_R8_OPCODE:
		cpu.RRC_R8(opcode)
	case RL_R8_OPCODE:
		cpu.RL_R8(opcode)
	case RR_R8_OPCODE:
		cpu.RR_R8(opcode)
	case SLA_R8_OPCODE:
		cpu.SLA_R8(opcode)
	case SRA_R8_OPCODE:
		cpu.SRA_R8(opcode)
	case SWAP_R8_OPCODE:
		cpu.SWAP_R8(opcode)
	case SRL_R8_OPCODE:
		cpu.SRL_R8(opcode)
	case BIT_0_R8_OPCODE:
		cpu.BIT_B3_R8(0, opcode)
	case BIT_1_R8_OPCODE:
		cpu.BIT_B3_R8(1, opcode)
	case BIT_2_R8_OPCODE:
		cpu.BIT_B3_R8(2, opcode)
	case BIT_3_R8_OPCODE:
		cpu.BIT_B3_R8(3, opcode)
	case BIT_4_R8_OPCODE:
		cpu.BIT_B3_R8(4, opcode)
	case BIT_5_R8_OPCODE:
		cpu.BIT_B3_R8(5, opcode)
	case BIT_6_R8_OPCODE:
		cpu.BIT_B3_R8(6, opcode)
	case BIT_7_R8_OPCODE:
		cpu.BIT_B3_R8(7, opcode)
	case RES_0_R8_OPCODE:
		cpu.RES_B3_R8(0, opcode)
	case RES_1_R8_OPCODE:
		cpu.RES_B3_R8(1, opcode)
	case RES_2_R8_OPCODE:
		cpu.RES_B3_R8(2, opcode)
	case RES_3_R8_OPCODE:
		cpu.RES_B3_R8(3, opcode)
	case RES_4_R8_OPCODE:
		cpu.RES_B3_R8(4, opcode)
	case RES_5_R8_OPCODE:
		cpu.RES_B3_R8(5, opcode)
	case RES_6_R8_OPCODE:
		cpu.RES_B3_R8(6, opcode)
	case RES_7_R8_OPCODE:
		cpu.RES_B3_R8(7, opcode)
	case SET_0_R8_OPCODE:
		cpu.SET_B3_R8(0, opcode)
	case SET_1_R8_OPCODE:
		cpu.SET_B3_R8(1, opcode)
	case SET_2_R8_OPCODE:
		cpu.SET_B3_R8(2, opcode)
	case SET_3_R8_OPCODE:
		cpu.SET_B3_R8(3, opcode)
	case SET_4_R8_OPCODE:
		cpu.SET_B3_R8(4, opcode)
	case SET_5_R8_OPCODE:
		cpu.SET_B3_R8(5, opcode)
	case SET_6_R8_OPCODE:
		cpu.SET_B3_R8(6, opcode)
	case SET_7_R8_OPCODE:
		cpu.SET_B3_R8(7, opcode)
	}

	cpu.cycles += PREFIX_OPCODES_CYCLES[opcode]
}

func (cpu *CPU) Reset() {
	cpu.A = 0x01
	cpu.F = 0xB0
	cpu.B = 0x00
	cpu.C = 0x13
	cpu.D = 0x00
	cpu.E = 0xD8
	cpu.H = 0x01
	cpu.L = 0x4D
	cpu.SP = 0xFFFE
	cpu.PC = 0x0100
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
