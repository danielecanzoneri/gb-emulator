package cpu

const (
	NOP_OPCODE = 0x00

	LD_BC_N16_OPCODE = 0x01
	LD_DE_N16_OPCODE = 0x11
	LD_HL_N16_OPCODE = 0x21
	LD_SP_N16_OPCODE = 0x31

	LD_BCMEM_A_OPCODE  = 0x02
	LD_DEMEM_A_OPCODE  = 0x12
	LD_HLIMEM_A_OPCODE = 0x22
	LD_HLDMEM_A_OPCODE = 0x32

	LD_A_BCMEM_OPCODE  = 0x0A
	LD_A_DEMEM_OPCODE  = 0x1A
	LD_A_HLIMEM_OPCODE = 0x2A
	LD_A_HLDMEM_OPCODE = 0x3A

	LD_N16_SP_OPCODE = 0x08

	INC_BC_OPCODE = 0x03
	INC_DE_OPCODE = 0x13
	INC_HL_OPCODE = 0x23
	INC_SP_OPCODE = 0x33

	DEC_BC_OPCODE = 0x0B
	DEC_DE_OPCODE = 0x1B
	DEC_HL_OPCODE = 0x2B
	DEC_SP_OPCODE = 0x3B

	ADD_HL_BC_OPCODE = 0x09
	ADD_HL_DE_OPCODE = 0x19
	ADD_HL_HL_OPCODE = 0x29
	ADD_HL_SP_OPCODE = 0x39

	INC_B_OPCODE     = 0x04
	INC_C_OPCODE     = 0x0C
	INC_D_OPCODE     = 0x14
	INC_E_OPCODE     = 0x1C
	INC_H_OPCODE     = 0x24
	INC_L_OPCODE     = 0x2C
	INC_HLMEM_OPCODE = 0x34
	INC_A_OPCODE     = 0x3C

	DEC_B_OPCODE     = 0x05
	DEC_C_OPCODE     = 0x0D
	DEC_D_OPCODE     = 0x15
	DEC_E_OPCODE     = 0x1D
	DEC_H_OPCODE     = 0x25
	DEC_L_OPCODE     = 0x2D
	DEC_HLMEM_OPCODE = 0x35
	DEC_A_OPCODE     = 0x3D

	LD_B_N8_OPCODE     = 0x06
	LD_C_N8_OPCODE     = 0x0E
	LD_D_N8_OPCODE     = 0x16
	LD_E_N8_OPCODE     = 0x1E
	LD_H_N8_OPCODE     = 0x26
	LD_L_N8_OPCODE     = 0x2E
	LD_HLMEM_N8_OPCODE = 0x36
	LD_A_N8_OPCODE     = 0x3E

	RLCA_OPCODE = 0x07
	RRCA_OPCODE = 0x0F
	RLA_OPCODE  = 0x17
	RRA_OPCODE  = 0x1F
	DAA_OPCODE  = 0x27
	CPL_OPCODE  = 0x2F
	SCF_OPCODE  = 0x37
	CCF_OPCODE  = 0x3F

	JR_E8_OPCODE = 0x18

	JR_NZ_E8_OPCODE = 0x20
	JR_Z_E8_OPCODE  = 0x28
	JR_NC_E8_OPCODE = 0x30
	JR_C_E8_OPCODE  = 0x38

	STOP_OPCODE = 0x10
)

var OPCODES_CYCLES = [256]int{
	4, 12, 8, 8, 4, 4, 8, 4, 20, 8, 8, 8, 4, 4, 8, 4,
	4, 12, 8, 8, 4, 4, 8, 4, 12, 8, 8, 8, 4, 4, 8, 4,
	8, 12, 8, 8, 4, 4, 8, 4, 8, 8, 8, 8, 4, 4, 8, 4,
	8, 12, 8, 8, 12, 12, 12, 4, 8, 8, 8, 8, 4, 4, 8, 4,
	4, 4, 4, 4, 4, 4, 8, 4, 4, 4, 4, 4, 4, 4, 8, 4,
	4, 4, 4, 4, 4, 4, 8, 4, 4, 4, 4, 4, 4, 4, 8, 4,
	4, 4, 4, 4, 4, 4, 8, 4, 4, 4, 4, 4, 4, 4, 8, 4,
	8, 8, 8, 8, 8, 8, 4, 8, 4, 4, 4, 4, 4, 4, 8, 4,
	4, 4, 4, 4, 4, 4, 8, 4, 4, 4, 4, 4, 4, 4, 8, 4,
	4, 4, 4, 4, 4, 4, 8, 4, 4, 4, 4, 4, 4, 4, 8, 4,
	4, 4, 4, 4, 4, 4, 8, 4, 4, 4, 4, 4, 4, 4, 8, 4,
	4, 4, 4, 4, 4, 4, 8, 4, 4, 4, 4, 4, 4, 4, 8, 4,
	8, 12, 12, 16, 12, 16, 8, 16, 8, 16, 12, 4, 12, 24, 8, 16,
	8, 12, 12, 0, 12, 16, 8, 16, 8, 16, 12, 0, 12, 0, 8, 16,
	2, 12, 8, 0, 0, 16, 8, 16, 16, 4, 16, 0, 0, 0, 8, 16,
	2, 12, 8, 4, 0, 16, 8, 16, 12, 8, 16, 4, 0, 0, 8, 16,
}

var OPCODES_CYCLES_BRANCH = map[uint8]int{
	0x20: 12, 0x28: 12, 0x30: 12, 0x38: 12,
	0xC0: 20, 0xC8: 20, 0xD0: 20, 0xD8: 20,
	0xC2: 16, 0xCA: 16, 0xD2: 16, 0xDA: 16,
	0xC4: 24, 0xCC: 24, 0xD4: 24, 0xDC: 24,
}

var OPCODES_BYTES = [256]int{
	1, 3, 1, 1, 1, 1, 2, 1, 3, 1, 1, 1, 1, 1, 2, 1,
	2, 3, 1, 1, 1, 1, 2, 1, 2, 1, 1, 1, 1, 1, 2, 1,
	2, 3, 1, 1, 1, 1, 2, 1, 2, 1, 1, 1, 1, 1, 2, 1,
	2, 3, 1, 1, 1, 1, 2, 1, 2, 1, 1, 1, 1, 1, 2, 1,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 3, 3, 3, 1, 2, 1, 1, 1, 3, 1, 3, 3, 2, 1,
	1, 1, 3, 0, 3, 1, 2, 1, 1, 1, 3, 0, 3, 0, 2, 1,
	2, 1, 1, 0, 0, 1, 2, 1, 2, 1, 3, 0, 0, 0, 2, 1,
	2, 1, 1, 1, 0, 1, 2, 1, 2, 1, 3, 1, 0, 0, 2, 1,
}
