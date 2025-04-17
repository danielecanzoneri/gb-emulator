package cpu

import "testing"

func Test_LD_R16_N16(t *testing.T) {
	var BYTE1 uint8 = 0xAC
	var BYTE2 uint8 = 0x12
	cpu := setup_CPU()

	t.Run("BC", func(t *testing.T) {
		writeTestProgram(cpu, LD_BC_N16_OPCODE, BYTE1, BYTE2)
		cpu.ExecuteInstruction()
		if !(cpu.B == BYTE2 && cpu.C == BYTE1) {
			t.Fatalf("got %X%X, expected %X%X", cpu.B, cpu.C, BYTE2, BYTE1)
		}
	})

	t.Run("DE", func(t *testing.T) {
		writeTestProgram(cpu, LD_DE_N16_OPCODE, BYTE1, BYTE2)
		cpu.ExecuteInstruction()
		if !(cpu.D == BYTE2 && cpu.E == BYTE1) {
			t.Fatalf("got %X%X, expected %X%X", cpu.D, cpu.E, BYTE2, BYTE1)
		}
	})

	t.Run("HL", func(t *testing.T) {
		writeTestProgram(cpu, LD_HL_N16_OPCODE, BYTE1, BYTE2)
		cpu.ExecuteInstruction()
		if !(cpu.H == BYTE2 && cpu.L == BYTE1) {
			t.Fatalf("got %X%X, expected %X%X", cpu.H, cpu.L, BYTE2, BYTE1)
		}
	})

	t.Run("SP", func(t *testing.T) {
		writeTestProgram(cpu, LD_SP_N16_OPCODE, BYTE1, BYTE2)
		cpu.ExecuteInstruction()
		expected := combineBytes(BYTE2, BYTE1)
		if !(cpu.SP == expected) {
			t.Fatalf("got %X%X, expected %X%X", cpu.D, cpu.E, BYTE2, BYTE1)
		}
	})
}

func Test_LD_R16MEM_A(t *testing.T) {
	cpu := setup_CPU()
	cpu.A = 0xFD

	var addr_BC uint16 = 0x10
	var addr_DE uint16 = 0x20
	var addr_HLI uint16 = 0x30
	var addr_HLD uint16 = 0x40

	t.Run("BC", func(t *testing.T) {
		cpu.writeBC(addr_BC)
		writeTestProgram(cpu, LD_BCMEM_A_OPCODE)
		cpu.ExecuteInstruction()
		if cpu.Mem.Read(addr_BC) != cpu.A {
			t.Fatalf("got %X, expected %X", cpu.Mem.Read(addr_BC), cpu.A)
		}
	})

	t.Run("DE", func(t *testing.T) {
		cpu.writeDE(addr_DE)
		writeTestProgram(cpu, LD_DEMEM_A_OPCODE)
		cpu.ExecuteInstruction()
		if cpu.Mem.Read(addr_DE) != cpu.A {
			t.Fatalf("got %X, expected %X", cpu.Mem.Read(addr_DE), cpu.A)
		}
	})

	t.Run("HLI", func(t *testing.T) {
		cpu.writeHL(addr_HLI)
		writeTestProgram(cpu, LD_HLIMEM_A_OPCODE)
		cpu.ExecuteInstruction()
		if cpu.Mem.Read(addr_HLI) != cpu.A {
			t.Fatalf("got %X, expected %X", cpu.Mem.Read(addr_HLI), cpu.A)
		}
		if cpu.readHL() != addr_HLI+1 {
			t.Error("increment failed")
		}
	})

	t.Run("HLD", func(t *testing.T) {
		cpu.writeHL(addr_HLD)
		writeTestProgram(cpu, LD_HLDMEM_A_OPCODE)
		cpu.ExecuteInstruction()
		if cpu.Mem.Read(addr_HLD) != cpu.A {
			t.Fatalf("got %X, expected %X", cpu.Mem.Read(addr_HLD), cpu.A)
		}
		if cpu.readHL() != addr_HLD-1 {
			t.Error("decrement failed")
		}
	})
}

func Test_LD_A_R16MEM(t *testing.T) {
	cpu := setup_CPU()

	var addr_BC uint16 = 0x10
	var byte_BC uint8 = 0x01
	var addr_DE uint16 = 0x20
	var byte_DE uint8 = 0x02
	var addr_HLI uint16 = 0x30
	var byte_HLI uint8 = 0x03
	var addr_HLD uint16 = 0x40
	var byte_HLD uint8 = 0x04

	t.Run("BC", func(t *testing.T) {
		cpu.writeBC(addr_BC)
		cpu.Mem.Write(addr_BC, byte_BC)
		writeTestProgram(cpu, LD_A_BCMEM_OPCODE)
		cpu.ExecuteInstruction()
		if cpu.A != byte_BC {
			t.Fatalf("got %X, expected %X", cpu.A, byte_BC)
		}
	})

	t.Run("DE", func(t *testing.T) {
		cpu.writeDE(addr_DE)
		cpu.Mem.Write(addr_DE, byte_DE)
		writeTestProgram(cpu, LD_A_DEMEM_OPCODE)
		cpu.ExecuteInstruction()
		if cpu.A != byte_DE {
			t.Fatalf("got %X, expected %X", cpu.A, byte_DE)
		}
	})

	t.Run("HLI", func(t *testing.T) {
		cpu.writeHL(addr_HLI)
		cpu.Mem.Write(addr_HLI, byte_HLI)
		writeTestProgram(cpu, LD_A_HLIMEM_OPCODE)
		cpu.ExecuteInstruction()
		if cpu.A != byte_HLI {
			t.Fatalf("got %X, expected %X", cpu.A, byte_HLI)
		}
		if cpu.readHL() != addr_HLI+1 {
			t.Error("increment failed")
		}
	})

	t.Run("HLD", func(t *testing.T) {
		cpu.writeHL(addr_HLD)
		cpu.Mem.Write(addr_HLD, byte_HLD)
		writeTestProgram(cpu, LD_A_HLDMEM_OPCODE)
		cpu.ExecuteInstruction()
		if cpu.A != byte_HLD {
			t.Fatalf("got %X, expected %X", cpu.A, byte_HLD)
		}
		if cpu.readHL() != addr_HLD-1 {
			t.Error("decrement failed")
		}
	})
}

func Test_LD_N16_SP(t *testing.T) {
	cpu := setup_CPU()
	cpu.SP = 0xFD53

	writeTestProgram(cpu, LD_N16_SP_OPCODE)
	cpu.ExecuteInstruction()

	read := cpu.Mem.ReadWord(cpu.PC - 2)
	if read != cpu.SP {
		t.Fatalf("CPU - LD_N16_SP failed: [N16] got %04X, expected %04X", read, cpu.SP)
	}
}

func Test_INC_R16(t *testing.T) {
	cpu := setup_CPU()

	var BC uint16 = 0x1234
	cpu.writeBC(BC)
	var DE uint16 = 0x4321
	cpu.writeDE(DE)
	var HL uint16 = 0x1111
	cpu.writeHL(HL)
	var SP uint16 = 0x2222
	cpu.SP = SP

	t.Run("BC", func(t *testing.T) {
		writeTestProgram(cpu, INC_BC_OPCODE)
		cpu.ExecuteInstruction()
		if cpu.readBC() != BC+1 {
			t.Fatalf("got %04X, expected %04X", cpu.readBC(), BC+1)
		}
	})

	t.Run("DE", func(t *testing.T) {
		writeTestProgram(cpu, INC_DE_OPCODE)
		cpu.ExecuteInstruction()
		if cpu.readDE() != DE+1 {
			t.Fatalf("got %04X, expected %04X", cpu.readDE(), DE+1)
		}
	})

	t.Run("HL", func(t *testing.T) {
		writeTestProgram(cpu, INC_HL_OPCODE)
		cpu.ExecuteInstruction()
		if cpu.readHL() != HL+1 {
			t.Fatalf("got %04X, expected %04X", cpu.readHL(), HL+1)
		}
	})

	t.Run("SP", func(t *testing.T) {
		writeTestProgram(cpu, INC_SP_OPCODE)
		cpu.ExecuteInstruction()
		if cpu.SP != SP+1 {
			t.Fatalf("got %04X, expected %04X", cpu.SP, SP+1)
		}
	})
}

func Test_DEC_R16(t *testing.T) {
	cpu := setup_CPU()

	var BC uint16 = 0x1234
	cpu.writeBC(BC)
	var DE uint16 = 0x4321
	cpu.writeDE(DE)
	var HL uint16 = 0x1111
	cpu.writeHL(HL)
	var SP uint16 = 0x2222
	cpu.SP = SP

	t.Run("BC", func(t *testing.T) {
		writeTestProgram(cpu, DEC_BC_OPCODE)
		cpu.ExecuteInstruction()
		if cpu.readBC() != BC-1 {
			t.Fatalf("got %04X, expected %04X", cpu.readBC(), BC-1)
		}
	})

	t.Run("DE", func(t *testing.T) {
		writeTestProgram(cpu, DEC_DE_OPCODE)
		cpu.ExecuteInstruction()
		if cpu.readDE() != DE-1 {
			t.Fatalf("got %04X, expected %04X", cpu.readDE(), DE-1)
		}
	})

	t.Run("HL", func(t *testing.T) {
		writeTestProgram(cpu, DEC_HL_OPCODE)
		cpu.ExecuteInstruction()
		if cpu.readHL() != HL-1 {
			t.Fatalf("got %04X, expected %04X", cpu.readHL(), HL-1)
		}
	})

	t.Run("SP", func(t *testing.T) {
		writeTestProgram(cpu, DEC_SP_OPCODE)
		cpu.ExecuteInstruction()
		if cpu.SP != SP-1 {
			t.Fatalf("got %04X, expected %04X", cpu.SP, SP-1)
		}
	})
}

func Test_ADD_HL_R16(t *testing.T) {
	cpu := setup_CPU()

	type test_add struct {
		hl    uint16
		r16   uint16
		sum   uint16
		exp_h uint8
		exp_c uint8
	}

	tests := map[string]test_add{
		"standard":   {hl: 0x00FF, r16: 0x00FF, sum: 0x01FE, exp_h: 0, exp_c: 0},
		"half-carry": {hl: 0x0800, r16: 0x0800, sum: 0x1000, exp_h: 1, exp_c: 0},
		"carry":      {hl: 0xFFFF, r16: 0xFFFF, sum: 0xFFFE, exp_h: 1, exp_c: 1},
	}

	test_R16 := func(test test_add, t *testing.T) {
		cpu.writeHL(test.hl)
		cpu.ExecuteInstruction()

		t.Log(cpu)
		if cpu.readHL() != test.sum {
			t.Fatalf("wrong sum: got %04X, expected %04X", cpu.readHL(), test.sum)
		}
		if cpu.readNFlag() != 0 {
			t.Error("wrong N flag: should be 0")
		}
		if cpu.readHFlag() != test.exp_h {
			t.Errorf("wrong H flag: got %d, expected %d", cpu.readHFlag(), test.exp_h)
		}
		if cpu.readCFlag() != test.exp_c {
			t.Errorf("wrong C flag: got %d, expected %d", cpu.readCFlag(), test.exp_c)
		}
	}

	for name, test := range tests {
		t.Run("BC_"+name, func(t *testing.T) {
			cpu.writeBC(test.r16)
			writeTestProgram(cpu, ADD_HL_BC_OPCODE)
			test_R16(test, t)
		})
		t.Run("DE_"+name, func(t *testing.T) {
			cpu.writeDE(test.r16)
			writeTestProgram(cpu, ADD_HL_DE_OPCODE)
			test_R16(test, t)
		})
		t.Run("HL_"+name, func(t *testing.T) {
			cpu.writeHL(test.r16)
			writeTestProgram(cpu, ADD_HL_HL_OPCODE)
			test_R16(test, t)
		})
		t.Run("SP_"+name, func(t *testing.T) {
			cpu.SP = test.r16
			writeTestProgram(cpu, ADD_HL_SP_OPCODE)
			test_R16(test, t)
		})
		cpu.F = 0
	}
}

func Test_INC_R8(t *testing.T) {
	cpu := setup_CPU()

	var B uint8 = 0x01
	var zflag_B, hflag_B uint8 = 0, 0
	var C uint8 = 0x02
	var zflag_C, hflag_C uint8 = 0, 0
	var D uint8 = 0x03
	var zflag_D, hflag_D uint8 = 0, 0
	var E uint8 = 0x04
	var zflag_E, hflag_E uint8 = 0, 0
	var H uint8 = 0x05
	var zflag_H, hflag_H uint8 = 0, 0
	var L uint8 = 0x06
	var zflag_L, hflag_L uint8 = 0, 0
	var A uint8 = 0x0F
	var zflag_A, hflag_A uint8 = 0, 1

	cpu.A = A
	cpu.B = B
	cpu.C = C
	cpu.D = D
	cpu.E = E
	cpu.H = H
	cpu.L = L

	testCarries := func(t *testing.T, exp_z, exp_h uint8) {
		if cpu.readNFlag() != 0 {
			t.Errorf("N flag: got %x, expected %x", cpu.readNFlag(), 0)
		}
		if cpu.readZFlag() != exp_z {
			t.Errorf("Z flag: got %x, expected %x", cpu.readZFlag(), exp_z)
		}
		if cpu.readHFlag() != exp_h {
			t.Errorf("H flag: got %x, expected %x", cpu.readHFlag(), exp_h)
		}
	}

	t.Run("B", func(t *testing.T) {
		writeTestProgram(cpu, INC_B_OPCODE)
		cpu.ExecuteInstruction()
		if cpu.B != B+1 {
			t.Fatalf("got %2X, expected %2X", cpu.B, B+1)
		}
		testCarries(t, zflag_B, hflag_B)
	})

	t.Run("C", func(t *testing.T) {
		writeTestProgram(cpu, INC_C_OPCODE)
		cpu.ExecuteInstruction()
		if cpu.C != C+1 {
			t.Fatalf("got %2X, expected %2X", cpu.C, C+1)
		}
		testCarries(t, zflag_C, hflag_C)
	})

	t.Run("D", func(t *testing.T) {
		writeTestProgram(cpu, INC_D_OPCODE)
		cpu.ExecuteInstruction()
		if cpu.D != D+1 {
			t.Fatalf("got %2X, expected %2X", cpu.D, D+1)
		}
		testCarries(t, zflag_D, hflag_D)
	})

	t.Run("E", func(t *testing.T) {
		writeTestProgram(cpu, INC_E_OPCODE)
		cpu.ExecuteInstruction()
		if cpu.E != E+1 {
			t.Fatalf("got %2X, expected %2X", cpu.E, E+1)
		}
		testCarries(t, zflag_E, hflag_E)
	})

	t.Run("H", func(t *testing.T) {
		writeTestProgram(cpu, INC_H_OPCODE)
		cpu.ExecuteInstruction()
		if cpu.H != H+1 {
			t.Fatalf("got %2X, expected %2X", cpu.H, H+1)
		}
		testCarries(t, zflag_H, hflag_H)
	})

	t.Run("L", func(t *testing.T) {
		writeTestProgram(cpu, INC_L_OPCODE)
		cpu.ExecuteInstruction()
		if cpu.L != L+1 {
			t.Fatalf("got %2X, expected %2X", cpu.L, L+1)
		}
		testCarries(t, zflag_L, hflag_L)
	})

	t.Run("A", func(t *testing.T) {
		writeTestProgram(cpu, INC_A_OPCODE)
		cpu.ExecuteInstruction()
		if cpu.A != A+1 {
			t.Fatalf("got %2X, expected %2X", cpu.A, A+1)
		}
		testCarries(t, zflag_A, hflag_A)
	})

	var addr_HL uint16 = 0x50
	var value uint8 = 0xFF
	var zflag_HL, hflag_HL uint8 = 1, 1

	cpu.writeHL(addr_HL)
	cpu.Mem.Write(addr_HL, value)

	t.Run("HLMEM", func(t *testing.T) {
		writeTestProgram(cpu, INC_HLMEM_OPCODE)
		cpu.ExecuteInstruction()
		if cpu.Mem.Read(addr_HL) != value+1 {
			t.Fatalf("got %2X, expected %2X", cpu.Mem.Read(addr_HL), value+1)
		}
		testCarries(t, zflag_HL, hflag_HL)
	})
}

func Test_DEC_R8(t *testing.T) {
	cpu := setup_CPU()

	var B uint8 = 0x01
	var zflag_B, hflag_B uint8 = 1, 0
	var C uint8 = 0x02
	var zflag_C, hflag_C uint8 = 0, 0
	var D uint8 = 0x03
	var zflag_D, hflag_D uint8 = 0, 0
	var E uint8 = 0x04
	var zflag_E, hflag_E uint8 = 0, 0
	var H uint8 = 0x05
	var zflag_H, hflag_H uint8 = 0, 0
	var L uint8 = 0x06
	var zflag_L, hflag_L uint8 = 0, 0
	var A uint8 = 0x10
	var zflag_A, hflag_A uint8 = 0, 1

	cpu.A = A
	cpu.B = B
	cpu.C = C
	cpu.D = D
	cpu.E = E
	cpu.H = H
	cpu.L = L

	testCarries := func(t *testing.T, exp_z, exp_h uint8) {
		if cpu.readNFlag() != 1 {
			t.Errorf("N flag: got %x, expected %x", cpu.readNFlag(), 1)
		}
		if cpu.readZFlag() != exp_z {
			t.Errorf("Z flag: got %x, expected %x", cpu.readZFlag(), exp_z)
		}
		if cpu.readHFlag() != exp_h {
			t.Errorf("H flag: got %x, expected %x", cpu.readHFlag(), exp_h)
		}
	}

	t.Run("B", func(t *testing.T) {
		writeTestProgram(cpu, DEC_B_OPCODE)
		cpu.ExecuteInstruction()
		if cpu.B != B-1 {
			t.Fatalf("got %2X, expected %2X", cpu.B, B-1)
		}
		testCarries(t, zflag_B, hflag_B)
	})

	t.Run("C", func(t *testing.T) {
		writeTestProgram(cpu, DEC_C_OPCODE)
		cpu.ExecuteInstruction()
		if cpu.C != C-1 {
			t.Fatalf("got %2X, expected %2X", cpu.C, C-1)
		}
		testCarries(t, zflag_C, hflag_C)
	})

	t.Run("D", func(t *testing.T) {
		writeTestProgram(cpu, DEC_D_OPCODE)
		cpu.ExecuteInstruction()
		if cpu.D != D-1 {
			t.Fatalf("got %2X, expected %2X", cpu.D, D-1)
		}
		testCarries(t, zflag_D, hflag_D)
	})

	t.Run("E", func(t *testing.T) {
		writeTestProgram(cpu, DEC_E_OPCODE)
		cpu.ExecuteInstruction()
		if cpu.E != E-1 {
			t.Fatalf("got %2X, expected %2X", cpu.E, E-1)
		}
		testCarries(t, zflag_E, hflag_E)
	})

	t.Run("H", func(t *testing.T) {
		writeTestProgram(cpu, DEC_H_OPCODE)
		cpu.ExecuteInstruction()
		if cpu.H != H-1 {
			t.Fatalf("got %2X, expected %2X", cpu.H, H-1)
		}
		testCarries(t, zflag_H, hflag_H)
	})

	t.Run("L", func(t *testing.T) {
		writeTestProgram(cpu, DEC_L_OPCODE)
		cpu.ExecuteInstruction()
		if cpu.L != L-1 {
			t.Fatalf("got %2X, expected %2X", cpu.L, L-1)
		}
		testCarries(t, zflag_L, hflag_L)
	})

	t.Run("A", func(t *testing.T) {
		writeTestProgram(cpu, DEC_A_OPCODE)
		cpu.ExecuteInstruction()
		if cpu.A != A-1 {
			t.Fatalf("got %2X, expected %2X", cpu.A, A-1)
		}
		testCarries(t, zflag_A, hflag_A)
	})

	var addr_HL uint16 = 0x50
	var value uint8 = 0x00
	var zflag_HL, hflag_HL uint8 = 0, 1

	cpu.writeHL(addr_HL)
	cpu.Mem.Write(addr_HL, value)

	t.Run("HLMEM", func(t *testing.T) {
		writeTestProgram(cpu, DEC_HLMEM_OPCODE)
		cpu.ExecuteInstruction()
		if cpu.Mem.Read(addr_HL) != value-1 {
			t.Fatalf("got %2X, expected %2X", cpu.Mem.Read(addr_HL), value-1)
		}
		testCarries(t, zflag_HL, hflag_HL)
	})
}

func Test_LD_R8_N8(t *testing.T) {
	cpu := setup_CPU()

	var value uint8 = 0xD1

	t.Run("B", func(t *testing.T) {
		writeTestProgram(cpu, LD_B_N8_OPCODE, value)
		cpu.ExecuteInstruction()
		if cpu.B != value {
			t.Fatalf("got %2X, expected %2X", cpu.B, value)
		}
	})

	t.Run("C", func(t *testing.T) {
		writeTestProgram(cpu, LD_C_N8_OPCODE, value)
		cpu.ExecuteInstruction()
		if cpu.C != value {
			t.Fatalf("got %2X, expected %2X", cpu.C, value)
		}
	})

	t.Run("D", func(t *testing.T) {
		writeTestProgram(cpu, LD_D_N8_OPCODE, value)
		cpu.ExecuteInstruction()
		if cpu.D != value {
			t.Fatalf("got %2X, expected %2X", cpu.D, value)
		}
	})

	t.Run("E", func(t *testing.T) {
		writeTestProgram(cpu, LD_E_N8_OPCODE, value)
		cpu.ExecuteInstruction()
		if cpu.E != value {
			t.Fatalf("got %2X, expected %2X", cpu.E, value)
		}
	})

	t.Run("H", func(t *testing.T) {
		writeTestProgram(cpu, LD_H_N8_OPCODE, value)
		cpu.ExecuteInstruction()
		if cpu.H != value {
			t.Fatalf("got %2X, expected %2X", cpu.H, value)
		}
	})

	t.Run("L", func(t *testing.T) {
		writeTestProgram(cpu, LD_L_N8_OPCODE, value)
		cpu.ExecuteInstruction()
		if cpu.L != value {
			t.Fatalf("got %2X, expected %2X", cpu.L, value)
		}
	})

	t.Run("HLMEM", func(t *testing.T) {
		var addr uint16 = 0xF0
		cpu.writeHL(addr)
		writeTestProgram(cpu, LD_HLMEM_N8_OPCODE, value)
		cpu.ExecuteInstruction()
		if cpu.Mem.Read(addr) != value {
			t.Fatalf("got %2X, expected %2X", cpu.B, value)
		}
	})

	t.Run("A", func(t *testing.T) {
		writeTestProgram(cpu, LD_A_N8_OPCODE, value)
		cpu.ExecuteInstruction()
		if cpu.A != value {
			t.Fatalf("got %2X, expected %2X", cpu.A, value)
		}
	})
}

func Test_RLCA(t *testing.T) {
	cpu := setup_CPU()

	tests := map[string]struct {
		A     uint8
		exp_A uint8
		exp_c uint8
	}{
		"carry":    {A: 0b10101010, exp_A: 0b01010101, exp_c: 1},
		"no-carry": {A: 0b01110001, exp_A: 0b11100010, exp_c: 0},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			cpu.A = test.A
			writeTestProgram(cpu, RLCA_OPCODE)
			cpu.ExecuteInstruction()

			if cpu.A != test.exp_A {
				t.Fatalf("got %08b, expected %08b", cpu.A, test.exp_A)
			}
			// Test flags
			if cpu.readZFlag() != 0 {
				t.Error("Z flag: got 1, expected 0")
			}
			if cpu.readNFlag() != 0 {
				t.Error("N flag: got 1, expected 0")
			}
			if cpu.readHFlag() != 0 {
				t.Error("H flag: got 1, expected 0")
			}
			if cpu.readCFlag() != test.exp_c {
				t.Errorf("C flag: got %x, expected %x", cpu.readCFlag(), test.exp_c)
			}
		})
	}
}

func Test_RRCA(t *testing.T) {
	cpu := setup_CPU()

	tests := map[string]struct {
		A     uint8
		exp_A uint8
		exp_c uint8
	}{
		"no-carry": {A: 0b10101010, exp_A: 0b01010101, exp_c: 0},
		"carry":    {A: 0b01110001, exp_A: 0b10111000, exp_c: 1},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			cpu.A = test.A
			writeTestProgram(cpu, RRCA_OPCODE)
			cpu.ExecuteInstruction()

			if cpu.A != test.exp_A {
				t.Fatalf("got %08b, expected %08b", cpu.A, test.exp_A)
			}
			// Test flags
			if cpu.readZFlag() != 0 {
				t.Error("Z flag: got 1, expected 0")
			}
			if cpu.readNFlag() != 0 {
				t.Error("N flag: got 1, expected 0")
			}
			if cpu.readHFlag() != 0 {
				t.Error("H flag: got 1, expected 0")
			}
			if cpu.readCFlag() != test.exp_c {
				t.Errorf("C flag: got %x, expected %x", cpu.readCFlag(), test.exp_c)
			}
		})
	}
}

func Test_RLA(t *testing.T) {
	cpu := setup_CPU()

	tests := map[string]struct {
		A     uint8
		carry uint8
		exp_A uint8
		exp_c uint8
	}{
		"carry_set":   {A: 0b01011010, carry: 1, exp_A: 0b10110101, exp_c: 0},
		"carry_unset": {A: 0b11011010, carry: 0, exp_A: 0b10110100, exp_c: 1},
		"carry":       {A: 0b11011010, carry: 1, exp_A: 0b10110101, exp_c: 1},
		"no-carry":    {A: 0b01011010, carry: 0, exp_A: 0b10110100, exp_c: 0},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			cpu.A = test.A
			cpu.setCFlag(test.carry)
			writeTestProgram(cpu, RLA_OPCODE)
			cpu.ExecuteInstruction()

			if cpu.A != test.exp_A {
				t.Fatalf("got %08b, expected %08b", cpu.A, test.exp_A)
			}
			// Test flags
			if cpu.readZFlag() != 0 {
				t.Error("Z flag: got 1, expected 0")
			}
			if cpu.readNFlag() != 0 {
				t.Error("N flag: got 1, expected 0")
			}
			if cpu.readHFlag() != 0 {
				t.Error("H flag: got 1, expected 0")
			}
			if cpu.readCFlag() != test.exp_c {
				t.Errorf("C flag: got %x, expected %x", cpu.readCFlag(), test.exp_c)
			}
		})
	}
}

func Test_RRA(t *testing.T) {
	cpu := setup_CPU()

	tests := map[string]struct {
		A     uint8
		carry uint8
		exp_A uint8
		exp_c uint8
	}{
		"carry_set":   {A: 0b01011010, carry: 1, exp_A: 0b10101101, exp_c: 0},
		"carry_unset": {A: 0b11011011, carry: 0, exp_A: 0b01101101, exp_c: 1},
		"carry":       {A: 0b11011011, carry: 1, exp_A: 0b11101101, exp_c: 1},
		"no-carry":    {A: 0b01011010, carry: 0, exp_A: 0b00101101, exp_c: 0},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			cpu.A = test.A
			cpu.setCFlag(test.carry)
			writeTestProgram(cpu, RRA_OPCODE)
			cpu.ExecuteInstruction()

			if cpu.A != test.exp_A {
				t.Fatalf("got %08b, expected %08b", cpu.A, test.exp_A)
			}
			// Test flags
			if cpu.readZFlag() != 0 {
				t.Error("Z flag: got 1, expected 0")
			}
			if cpu.readNFlag() != 0 {
				t.Error("N flag: got 1, expected 0")
			}
			if cpu.readHFlag() != 0 {
				t.Error("H flag: got 1, expected 0")
			}
			if cpu.readCFlag() != test.exp_c {
				t.Errorf("C flag: got %x, expected %x", cpu.readCFlag(), test.exp_c)
			}
		})
	}
}

func Test_DAA(t *testing.T) {
	cpu := setup_CPU()

	tests := map[string]struct {
		n1     uint8
		n2     uint8
		is_sub uint8
	}{
		"sum-hflag":            {n1: 19, n2: 19, is_sub: 0},
		"sum-low_adj":          {n1: 15, n2: 19, is_sub: 0},
		"sum-hflag,high_adj":   {n1: 69, n2: 59, is_sub: 0},
		"sum-low_adj,high_adj": {n1: 65, n2: 59, is_sub: 0},
		"sum-carry":            {n1: 99, n2: 99, is_sub: 0},
		"sum-zflag":            {n1: 50, n2: 50, is_sub: 0},
		"sub-hflag":            {n1: 10, n2: 1, is_sub: 1},
		"sub-cflag":            {n1: 100, n2: 10, is_sub: 1},
		"sub-hflag,cflag":      {n1: 100, n2: 1, is_sub: 1},
		"sub-zflag":            {n1: 100, n2: 100, is_sub: 1},
		"sub-negative":         {n1: 50, n2: 60, is_sub: 1},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			bcd_n1 := byteToBCD(test.n1)
			bcd_n2 := byteToBCD(test.n2)
			var result, carry, half_carry uint8
			if test.is_sub == 0 {
				result, carry, half_carry = sumBytesWithCarry(bcd_n1, bcd_n2)
			} else {
				result, carry, half_carry = subBytesWithCarry(bcd_n1, bcd_n2)
			}
			cpu.A = result
			cpu.setHFlag(half_carry)
			cpu.setCFlag(carry)
			cpu.setNFlag(test.is_sub)

			writeTestProgram(cpu, DAA_OPCODE)
			var c_flag, exp_A uint8
			if test.is_sub == 0 {
				c_flag = (test.n1 + test.n2) / 100
				exp_A = byteToBCD((test.n1 + test.n2) % 100)
			} else {
				if test.n1 < test.n2 {
					c_flag = 1
					exp_A = byteToBCD(100 - (test.n2 - test.n1))
				} else {
					c_flag = 0
					exp_A = byteToBCD(test.n1 - test.n2)
				}
				cpu.setNFlag(1)
			}

			cpu.ExecuteInstruction()

			if cpu.readCFlag() != c_flag {
				t.Errorf("C flag: got %x, expected %x", cpu.readCFlag(), c_flag)
			}
			if cpu.readHFlag() != 0 {
				t.Errorf("H flag: got 1, expected 0")
			}
			if cpu.readZFlag() != isByteZeroUint8(exp_A) {
				t.Errorf("Z flag: got %x, expected %x", cpu.readZFlag(), isByteZeroUint8(exp_A))
			}
			if cpu.A != exp_A {
				t.Fatalf("got %02X, expected %02X", cpu.A, exp_A)
			}
		})
	}
}

func Test_CPL(t *testing.T) {
	cpu := setup_CPU()
	var value uint8 = 0b10010111

	cpu.A = value
	writeTestProgram(cpu, CPL_OPCODE)
	cpu.ExecuteInstruction()

	if cpu.A != ^value {
		t.Fatalf("got %08b, expected %08b", cpu.A, ^value)
	}
	// Test registers
	if cpu.readNFlag() != 1 {
		t.Error("N flag: got 0, expected 1")
	}
	if cpu.readHFlag() != 1 {
		t.Error("H flag: got 0, expected 1")
	}
}

func Test_SCF(t *testing.T) {
	cpu := setup_CPU()

	// Test cpu.F = 0xFF, 0x00
	tests := map[string]struct {
		F uint8
	}{
		"SCF_Keep": {0xFF},
		"SCF_Set":  {0x00},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			cpu.F = test.F
			writeTestProgram(cpu, SCF_OPCODE)
			cpu.ExecuteInstruction()

			if cpu.readNFlag() != 0 {
				t.Error("N flag: got 1, expected 0")
			}
			if cpu.readHFlag() != 0 {
				t.Error("H flag: got 1, expected 0")
			}
			if cpu.readCFlag() != 1 {
				t.Error("C flag: got 0, expected 1")
			}
		})
	}
}

func Test_CCF(t *testing.T) {
	cpu := setup_CPU()

	// Test cpu.F = 0xFF, 0x00
	tests := map[string]struct {
		F     uint8
		exp_c uint8
	}{
		"set":   {0xFF, 0},
		"unset": {0x00, 1},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			cpu.F = test.F
			writeTestProgram(cpu, CCF_OPCODE)
			cpu.ExecuteInstruction()

			if cpu.readNFlag() != 0 {
				t.Error("N flag: got 1, expected 0")
			}
			if cpu.readHFlag() != 0 {
				t.Error("H flag: got 1, expected 0")
			}
			if cpu.readCFlag() != test.exp_c {
				t.Errorf("C flag: got %x, expected %x", cpu.readCFlag(), test.exp_c)
			}
		})
	}
}

func Test_JR_E8(t *testing.T) {
	cpu := setup_CPU()

	tests := map[string]struct {
		e8 int8
	}{
		"zero_jump": {e8: 0},
		"pos_jump":  {e8: 2},
		"neg_jump":  {e8: -2},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			expected_PC := int(cpu.PC) + OPCODES_BYTES[JR_E8_OPCODE] + int(test.e8)

			writeTestProgram(cpu, JR_E8_OPCODE, uint8(test.e8))
			cpu.ExecuteInstruction()

			if cpu.PC != uint16(expected_PC) {
				t.Fatalf("got %04X, expected %04X", cpu.PC, expected_PC)
			}
		})
	}
}

func Test_JR_COND_E8(t *testing.T) {
	cpu := setup_CPU()

	tests := map[string]struct {
		e8 int8
	}{
		"zero_jump": {e8: 0},
		"pos_jump":  {e8: 4},
		"neg_jump":  {e8: -4},
	}
	// Prepare flags for each condition
	bool_to_int := map[bool]uint8{false: 0, true: 1}
	conditions := map[string]func(bool){
		"Z":  func(cc bool) { cpu.setZFlag(bool_to_int[cc]) },
		"NZ": func(cc bool) { cpu.setZFlag(bool_to_int[!cc]) },
		"C":  func(cc bool) { cpu.setCFlag(bool_to_int[cc]) },
		"NC": func(cc bool) { cpu.setCFlag(bool_to_int[!cc]) },
	}
	opcodes := map[string]uint8{
		"Z":  JR_Z_E8_OPCODE,
		"NZ": JR_NZ_E8_OPCODE,
		"C":  JR_C_E8_OPCODE,
		"NC": JR_NC_E8_OPCODE,
	}

	for cond, set_flag := range conditions {
		for name, test := range tests {
			t.Run(cond+"/"+name, func(t *testing.T) {
				// Condition not met
				expected_PC := int(cpu.PC) + OPCODES_BYTES[opcodes[cond]]
				expected_cycles := cpu.cycles + OPCODES_CYCLES[opcodes[cond]]
				set_flag(false)

				writeTestProgram(cpu, opcodes[cond], uint8(test.e8))
				cpu.ExecuteInstruction()

				if cpu.PC != uint16(expected_PC) {
					t.Errorf("condition not met: got %04X, expected %04X", cpu.PC, expected_PC)
				}
				if cpu.cycles != expected_cycles {
					t.Errorf("condition not met: got %v cycles, expected %v", cpu.cycles, expected_cycles)
				}

				// Condition met
				expected_PC = int(cpu.PC) + OPCODES_BYTES[opcodes[cond]] + int(test.e8)
				expected_cycles = cpu.cycles + OPCODES_CYCLES_BRANCH[opcodes[cond]]
				set_flag(true)

				writeTestProgram(cpu, opcodes[cond], uint8(test.e8))
				cpu.ExecuteInstruction()

				if cpu.PC != uint16(expected_PC) {
					t.Errorf("condition met: got %04X, expected %04X", cpu.PC, expected_PC)
				}
				if cpu.cycles != expected_cycles {
					t.Errorf("condition not met: got %v cycles, expected %v", cpu.cycles, expected_cycles)
				}
			})
		}
	}
}

func Test_STOP(t *testing.T) {
	// TODO - Has to be thoroughly tested
	cpu := setup_CPU()

	expected_PC := cpu.PC + STOP_BYTES

	writeTestProgram(cpu, STOP_OPCODE)
	cpu.ExecuteInstruction()

	if cpu.PC != expected_PC {
		t.Fatalf("PC: got %04X, expected %04X", cpu.PC, expected_PC)
	}
}

// Mock implementation of Memory interface
type MockMemory struct {
	data [0x100]byte
}

func (m *MockMemory) Read(addr uint16) uint8 {
	return m.data[addr]
}
func (m *MockMemory) Write(addr uint16, value uint8) {
	m.data[addr] = value
}
func (m *MockMemory) ReadWord(addr uint16) uint16 {
	return uint16(m.data[addr]) | (uint16(m.data[addr+1]) << 8)
}
func (m *MockMemory) WriteWord(addr uint16, value uint16) {
	m.data[addr] = uint8(value)
	m.data[addr+1] = uint8(value >> 8)
}

func setup_CPU() *CPU {
	return &CPU{Mem: &MockMemory{}}
}

func writeTestProgram(cpu *CPU, data ...byte) {
	for i, b := range data {
		cpu.Mem.Write(uint16(i)+cpu.PC, b)
	}
}
