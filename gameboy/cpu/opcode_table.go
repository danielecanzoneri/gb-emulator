package cpu

func (cpu *CPU) initOpcodeTable() {
	cpu.opcodesTable = [256]func(){
		/*  x0                    x1                   x2                    x3                   x4                    x5                   x6                    x7
		    x8                    x9                   xA                    xB                   xC                    xD                   xE                    xF                */
		cpu.NOP /*        */, cpu.LD_BC_N16 /* */, cpu.LD_BCmem_A /* */, cpu.INC_BC /*    */, cpu.INC_B /*      */, cpu.DEC_B /*     */, cpu.LD_B_N8 /*    */, cpu.RLCA, /*       0x */
		cpu.LD_N16_SP /*  */, cpu.ADD_HL_BC /* */, cpu.LD_A_BCmem /* */, cpu.DEC_BC /*    */, cpu.INC_C /*      */, cpu.DEC_C /*     */, cpu.LD_C_N8 /*    */, cpu.RRCA,
		cpu.STOP /*       */, cpu.LD_DE_N16 /* */, cpu.LD_DEmem_A /* */, cpu.INC_DE /*    */, cpu.INC_D /*      */, cpu.DEC_D /*     */, cpu.LD_D_N8 /*    */, cpu.RLA, /*        1x */
		cpu.JR_E8 /*      */, cpu.ADD_HL_DE /* */, cpu.LD_A_DEmem /* */, cpu.DEC_DE /*    */, cpu.INC_E /*      */, cpu.DEC_E /*     */, cpu.LD_E_N8 /*    */, cpu.RRA,
		cpu.JR_NZ_E8 /*   */, cpu.LD_HL_N16 /* */, cpu.LD_HLImem_A /**/, cpu.INC_HL /*    */, cpu.INC_H /*      */, cpu.DEC_H /*     */, cpu.LD_H_N8 /*    */, cpu.DAA, /*        2x */
		cpu.JR_Z_E8 /*    */, cpu.ADD_HL_HL /* */, cpu.LD_A_HLImem /**/, cpu.DEC_HL /*    */, cpu.INC_L /*      */, cpu.DEC_L /*     */, cpu.LD_L_N8 /*    */, cpu.CPL,
		cpu.JR_NC_E8 /*   */, cpu.LD_SP_N16 /* */, cpu.LD_HLDmem_A /**/, cpu.INC_SP /*    */, cpu.INC_HLmem /*  */, cpu.DEC_HLmem /* */, cpu.LD_HLmem_N8 /**/, cpu.SCF, /*        3x */
		cpu.JR_C_E8 /*    */, cpu.ADD_HL_SP /* */, cpu.LD_A_HLDmem /**/, cpu.DEC_SP /*    */, cpu.INC_A /*      */, cpu.DEC_A /*     */, cpu.LD_A_N8 /*    */, cpu.CCF,
		cpu.LD_B_B /*     */, cpu.LD_B_C /*    */, cpu.LD_B_D /*     */, cpu.LD_B_E /*    */, cpu.LD_B_H /*     */, cpu.LD_B_L /*    */, cpu.LD_B_HLmem /* */, cpu.LD_B_A, /*     4x */
		cpu.LD_C_B /*     */, cpu.LD_C_C /*    */, cpu.LD_C_D /*     */, cpu.LD_C_E /*    */, cpu.LD_C_H /*     */, cpu.LD_C_L /*    */, cpu.LD_C_HLmem /* */, cpu.LD_C_A,
		cpu.LD_D_B /*     */, cpu.LD_D_C /*    */, cpu.LD_D_D /*     */, cpu.LD_D_E /*    */, cpu.LD_D_H /*     */, cpu.LD_D_L /*    */, cpu.LD_D_HLmem /* */, cpu.LD_D_A, /*     5x */
		cpu.LD_E_B /*     */, cpu.LD_E_C /*    */, cpu.LD_E_D /*     */, cpu.LD_E_E /*    */, cpu.LD_E_H /*     */, cpu.LD_E_L /*    */, cpu.LD_E_HLmem /* */, cpu.LD_E_A,
		cpu.LD_H_B /*     */, cpu.LD_H_C /*    */, cpu.LD_H_D /*     */, cpu.LD_H_E /*    */, cpu.LD_H_H /*     */, cpu.LD_H_L /*    */, cpu.LD_H_HLmem /* */, cpu.LD_H_A, /*     6x */
		cpu.LD_L_B /*     */, cpu.LD_L_C /*    */, cpu.LD_L_D /*     */, cpu.LD_L_E /*    */, cpu.LD_L_H /*     */, cpu.LD_L_L /*    */, cpu.LD_L_HLmem /* */, cpu.LD_L_A,
		cpu.LD_HLmem_B /* */, cpu.LD_HLmem_C /**/, cpu.LD_HLmem_D /* */, cpu.LD_HLmem_E /**/, cpu.LD_HLmem_H /* */, cpu.LD_HLmem_L /**/, cpu.HALT /*       */, cpu.LD_HLmem_A, /* 7x */
		cpu.LD_A_B /*     */, cpu.LD_A_C /*    */, cpu.LD_A_D /*     */, cpu.LD_A_E /*    */, cpu.LD_A_H /*     */, cpu.LD_A_L /*    */, cpu.LD_A_HLmem /* */, cpu.LD_A_A,
		cpu.ADD_A_B /*    */, cpu.ADD_A_C /*   */, cpu.ADD_A_D /*    */, cpu.ADD_A_E /*   */, cpu.ADD_A_H /*    */, cpu.ADD_A_L /*   */, cpu.ADD_A_HLmem /**/, cpu.ADD_A_A, /*    8x */
		cpu.ADC_A_B /*    */, cpu.ADC_A_C /*   */, cpu.ADC_A_D /*    */, cpu.ADC_A_E /*   */, cpu.ADC_A_H /*    */, cpu.ADC_A_L /*   */, cpu.ADC_A_HLmem /**/, cpu.ADC_A_A,
		cpu.SUB_A_B /*    */, cpu.SUB_A_C /*   */, cpu.SUB_A_D /*    */, cpu.SUB_A_E /*   */, cpu.SUB_A_H /*    */, cpu.SUB_A_L /*   */, cpu.SUB_A_HLmem /**/, cpu.SUB_A_A, /*    9x */
		cpu.SBC_A_B /*    */, cpu.SBC_A_C /*   */, cpu.SBC_A_D /*    */, cpu.SBC_A_E /*   */, cpu.SBC_A_H /*    */, cpu.SBC_A_L /*   */, cpu.SBC_A_HLmem /**/, cpu.SBC_A_A,
		cpu.AND_A_B /*    */, cpu.AND_A_C /*   */, cpu.AND_A_D /*    */, cpu.AND_A_E /*   */, cpu.AND_A_H /*    */, cpu.AND_A_L /*   */, cpu.AND_A_HLmem /**/, cpu.AND_A_A, /*    Ax */
		cpu.XOR_A_B /*    */, cpu.XOR_A_C /*   */, cpu.XOR_A_D /*    */, cpu.XOR_A_E /*   */, cpu.XOR_A_H /*    */, cpu.XOR_A_L /*   */, cpu.XOR_A_HLmem /**/, cpu.XOR_A_A,
		cpu.OR_A_B /*     */, cpu.OR_A_C /*    */, cpu.OR_A_D /*     */, cpu.OR_A_E /*    */, cpu.OR_A_H /*     */, cpu.OR_A_L /*    */, cpu.OR_A_HLmem /* */, cpu.OR_A_A, /*     Bx */
		cpu.CP_A_B /*     */, cpu.CP_A_C /*    */, cpu.CP_A_D /*     */, cpu.CP_A_E /*    */, cpu.CP_A_H /*     */, cpu.CP_A_L /*    */, cpu.CP_A_HLmem /* */, cpu.CP_A_A,
		cpu.RET_NZ /*     */, cpu.POP_BC /*    */, cpu.JP_NZ_N16 /*  */, cpu.JP_N16 /*    */, cpu.CALL_NZ_N16 /**/, cpu.PUSH_BC /*   */, cpu.ADD_A_N8 /*   */, cpu.RST_00, /*     Cx */
		cpu.RET_Z /*      */, cpu.RET /*       */, cpu.JP_Z_N16 /*   */, cpu.PREFIX /*    */, cpu.CALL_Z_N16 /* */, cpu.CALL_N16 /*  */, cpu.ADC_A_N8 /*   */, cpu.RST_08,
		cpu.RET_NC /*     */, cpu.POP_DE /*    */, cpu.JP_NC_N16 /*  */, cpu.INVALID /*   */, cpu.CALL_NC_N16 /**/, cpu.PUSH_DE /*   */, cpu.SUB_A_N8 /*   */, cpu.RST_10, /*     Dx */
		cpu.RET_C /*      */, cpu.RETI /*      */, cpu.JP_C_N16 /*   */, cpu.INVALID /*   */, cpu.CALL_C_N16 /* */, cpu.INVALID /*   */, cpu.SBC_A_N8 /*   */, cpu.RST_18,
		cpu.LDH_N8_A /*   */, cpu.POP_HL /*    */, cpu.LDH_C_A /*    */, cpu.INVALID /*   */, cpu.INVALID /*    */, cpu.PUSH_HL /*   */, cpu.AND_A_N8 /*   */, cpu.RST_20, /*     Ex */
		cpu.ADD_SP_E8 /*  */, cpu.JP_HL /*     */, cpu.LD_N16_A /*   */, cpu.INVALID /*   */, cpu.INVALID /*    */, cpu.INVALID /*   */, cpu.XOR_A_N8 /*   */, cpu.RST_28,
		cpu.LDH_A_N8 /*   */, cpu.POP_AF /*    */, cpu.LDH_A_C /*    */, cpu.DI /*        */, cpu.INVALID /*    */, cpu.PUSH_AF /*   */, cpu.OR_A_N8 /*    */, cpu.RST_30, /*     Fx */
		cpu.LD_HL_SP_E8 /**/, cpu.LD_SP_HL /*  */, cpu.LD_A_N16 /*   */, cpu.EI /*        */, cpu.INVALID /*    */, cpu.INVALID /*   */, cpu.CP_A_N8 /*    */, cpu.RST_38,
	}

	cpu.prefixedOpcodesTable = [32]func(opcode uint8){
		/*  x0                                             x8                                          */
		cpu.RLC_R8 /*                              */, cpu.RRC_R8, /*                               0x */
		cpu.RL_R8 /*                               */, cpu.RR_R8, /*                                1x */
		cpu.SLA_R8 /*                              */, cpu.SRA_R8, /*                               2x */
		cpu.SWAP_R8 /*                             */, cpu.SRL_R8, /*                               3x */
		func(op uint8) { cpu.BIT_B3_R8(0, op) } /* */, func(op uint8) { cpu.BIT_B3_R8(1, op) }, /*  4x */
		func(op uint8) { cpu.BIT_B3_R8(2, op) } /* */, func(op uint8) { cpu.BIT_B3_R8(3, op) }, /*  5x */
		func(op uint8) { cpu.BIT_B3_R8(4, op) } /* */, func(op uint8) { cpu.BIT_B3_R8(5, op) }, /*  6x */
		func(op uint8) { cpu.BIT_B3_R8(6, op) } /* */, func(op uint8) { cpu.BIT_B3_R8(7, op) }, /*  7x */
		func(op uint8) { cpu.RES_B3_R8(0, op) } /* */, func(op uint8) { cpu.RES_B3_R8(1, op) }, /*  8x */
		func(op uint8) { cpu.RES_B3_R8(2, op) } /* */, func(op uint8) { cpu.RES_B3_R8(3, op) }, /*  9x */
		func(op uint8) { cpu.RES_B3_R8(4, op) } /* */, func(op uint8) { cpu.RES_B3_R8(5, op) }, /*  Ax */
		func(op uint8) { cpu.RES_B3_R8(6, op) } /* */, func(op uint8) { cpu.RES_B3_R8(7, op) }, /*  Bx */
		func(op uint8) { cpu.SET_B3_R8(0, op) } /* */, func(op uint8) { cpu.SET_B3_R8(1, op) }, /*  Cx */
		func(op uint8) { cpu.SET_B3_R8(2, op) } /* */, func(op uint8) { cpu.SET_B3_R8(3, op) }, /*  Dx */
		func(op uint8) { cpu.SET_B3_R8(4, op) } /* */, func(op uint8) { cpu.SET_B3_R8(5, op) }, /*  Ex */
		func(op uint8) { cpu.SET_B3_R8(6, op) } /* */, func(op uint8) { cpu.SET_B3_R8(7, op) }, /*  Fx */
	}

}
