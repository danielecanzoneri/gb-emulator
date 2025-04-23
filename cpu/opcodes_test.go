package cpu

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"
)

func Test_LD_R16_N16(t *testing.T) {
	var BYTE1 uint8 = 0xAC
	var BYTE2 uint8 = 0x12
	cpu := mockCPU()

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
	cpu := mockCPU()
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
	cpu := mockCPU()

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
	cpu := mockCPU()
	cpu.SP = 0xFD53
	var addr uint16 = 0x1234

	writeTestProgram(cpu, LD_N16_SP_OPCODE, 0x34, 0x12)
	cpu.ExecuteInstruction()

	read := cpu.Mem.ReadWord(addr)
	if read != cpu.SP {
		t.Fatalf("[N16] got %04X, expected %04X", read, cpu.SP)
	}
}

func Test_INC_R16(t *testing.T) {
	cpu := mockCPU()

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
	cpu := mockCPU()

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
	cpu := mockCPU()
	cpu.F = 0xFF // Set flags

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
	cpu := mockCPU()

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
	cpu := mockCPU()

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
	cpu := mockCPU()

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
	cpu := mockCPU()
	// Test flag reset
	cpu.F = 0xFF

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
	cpu := mockCPU()
	// Test flag reset
	cpu.F = 0xFF

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
	cpu := mockCPU()
	// Test flag reset
	cpu.F = 0xFF

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
	cpu := mockCPU()
	// Test flag reset
	cpu.F = 0xFF

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
	cpu := mockCPU()

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
	cpu := mockCPU()
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
	cpu := mockCPU()

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
	cpu := mockCPU()

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
	cpu := mockCPU()

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
	cpu := mockCPU()

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

func Test_LD_R8_R8(t *testing.T) {
	cpu := mockCPU()

	firstReg := map[string]struct {
		opcode uint8
		read   func() uint8
	}{
		"B":     {0x40, func() uint8 { return cpu.B }},
		"C":     {0x48, func() uint8 { return cpu.C }},
		"D":     {0x50, func() uint8 { return cpu.D }},
		"E":     {0x58, func() uint8 { return cpu.E }},
		"H":     {0x60, func() uint8 { return cpu.H }},
		"L":     {0x68, func() uint8 { return cpu.L }},
		"HLMEM": {0x70, func() uint8 { return cpu.Mem.Read(cpu.readHL()) }},
		"A":     {0x78, func() uint8 { return cpu.A }},
	}
	secondReg := map[string]struct {
		offset uint8
		set    func(uint8)
	}{
		"B":     {0, func(value uint8) { cpu.B = value }},
		"C":     {1, func(value uint8) { cpu.C = value }},
		"D":     {2, func(value uint8) { cpu.D = value }},
		"E":     {3, func(value uint8) { cpu.E = value }},
		"H":     {4, func(value uint8) { cpu.H = value }},
		"L":     {5, func(value uint8) { cpu.L = value }},
		"HLMEM": {6, func(value uint8) { cpu.H = 1; cpu.L = value; cpu.Mem.Write(cpu.readHL(), value) }},
		"A":     {7, func(value uint8) { cpu.A = value }},
	}

	for name1, receiver := range firstReg {
		for name2, loaded := range secondReg {
			opcode := receiver.opcode + loaded.offset
			// Skip HALT
			if opcode == HALT_OPCODE {
				continue
			}
			t.Run(name1+"<-"+name2, func(t *testing.T) {
				value := uint8(rand.Intn(0x100))
				loaded.set(value)

				writeTestProgram(cpu, opcode)
				cpu.ExecuteInstruction()
				if receiver.read() != value {
					t.Fatalf("got %02X, expected %02X", receiver.read(), value)
				}
			})
		}
	}
}

func Test_HALT(t *testing.T) {
	t.Fatalf("TODO")
}

func Test_ADD_A_R8(t *testing.T) {
	cpu := mockCPU()

	type test struct {
		A     uint8
		R8    uint8
		exp_z uint8
		exp_h uint8
		exp_c uint8
	}

	tests := map[string]test{
		"B":     {0x08, 0x08, 0, 1, 0},
		"C":     {0xF0, 0x11, 0, 0, 1},
		"D":     {0xFF, 0xFF, 0, 1, 1},
		"E":     {0xFF, 0x00, 0, 0, 0},
		"H":     {0x80, 0x08, 0, 0, 0},
		"L":     {0x57, 0xAD, 0, 1, 1},
		"HLMEM": {0x34, 0x12, 0, 0, 0},
		"A":     {0x80, 0x80, 1, 0, 1},
	}

	testCarries := func(t *testing.T, tst test) {
		if cpu.readNFlag() != 0 {
			t.Errorf("N flag: got %x, expected %x", cpu.readNFlag(), 0)
		}
		if cpu.readZFlag() != tst.exp_z {
			t.Errorf("Z flag: got %x, expected %x", cpu.readZFlag(), tst.exp_z)
		}
		if cpu.readHFlag() != tst.exp_h {
			t.Errorf("H flag: got %x, expected %x", cpu.readHFlag(), tst.exp_h)
		}
		if cpu.readCFlag() != tst.exp_c {
			t.Errorf("C flag: got %x, expected %x", cpu.readCFlag(), tst.exp_c)
		}
	}

	t.Run("B", func(t *testing.T) {
		tst := tests["B"]
		cpu.A = tst.A
		cpu.B = tst.R8
		writeTestProgram(cpu, ADD_A_B_OPCODE)
		cpu.ExecuteInstruction()
		if cpu.A != tst.A+tst.R8 {
			t.Fatalf("got %2X, expected %2X", cpu.A, tst.A+tst.R8)
		}
		testCarries(t, tst)
	})

	t.Run("C", func(t *testing.T) {
		tst := tests["C"]
		cpu.A = tst.A
		cpu.C = tst.R8
		writeTestProgram(cpu, ADD_A_C_OPCODE)
		cpu.ExecuteInstruction()
		if cpu.A != tst.A+tst.R8 {
			t.Fatalf("got %2X, expected %2X", cpu.A, tst.A+tst.R8)
		}
		testCarries(t, tst)
	})

	t.Run("D", func(t *testing.T) {
		tst := tests["D"]
		cpu.A = tst.A
		cpu.D = tst.R8
		writeTestProgram(cpu, ADD_A_D_OPCODE)
		cpu.ExecuteInstruction()
		if cpu.A != tst.A+tst.R8 {
			t.Fatalf("got %2X, expected %2X", cpu.A, tst.A+tst.R8)
		}
		testCarries(t, tst)
	})

	t.Run("E", func(t *testing.T) {
		tst := tests["E"]
		cpu.A = tst.A
		cpu.E = tst.R8
		writeTestProgram(cpu, ADD_A_E_OPCODE)
		cpu.ExecuteInstruction()
		if cpu.A != tst.A+tst.R8 {
			t.Fatalf("got %2X, expected %2X", cpu.A, tst.A+tst.R8)
		}
		testCarries(t, tst)
	})

	t.Run("H", func(t *testing.T) {
		tst := tests["H"]
		cpu.A = tst.A
		cpu.H = tst.R8
		writeTestProgram(cpu, ADD_A_H_OPCODE)
		cpu.ExecuteInstruction()
		if cpu.A != tst.A+tst.R8 {
			t.Fatalf("got %2X, expected %2X", cpu.A, tst.A+tst.R8)
		}
		testCarries(t, tst)
	})

	t.Run("L", func(t *testing.T) {
		tst := tests["L"]
		cpu.A = tst.A
		cpu.L = tst.R8
		writeTestProgram(cpu, ADD_A_L_OPCODE)
		cpu.ExecuteInstruction()
		if cpu.A != tst.A+tst.R8 {
			t.Fatalf("got %2X, expected %2X", cpu.A, tst.A+tst.R8)
		}
		testCarries(t, tst)
	})

	t.Run("A", func(t *testing.T) {
		tst := tests["A"]
		cpu.A = tst.A
		cpu.A = tst.R8
		writeTestProgram(cpu, ADD_A_A_OPCODE)
		cpu.ExecuteInstruction()
		if cpu.A != tst.A+tst.R8 {
			t.Fatalf("got %2X, expected %2X", cpu.A, tst.A+tst.R8)
		}
		testCarries(t, tst)
	})

	t.Run("HLMEM", func(t *testing.T) {
		tst := tests["HLMEM"]
		cpu.A = tst.A
		cpu.Mem.Write(cpu.readHL(), tst.R8)
		writeTestProgram(cpu, ADD_A_HLMEM_OPCODE)
		cpu.ExecuteInstruction()
		if cpu.A != tst.A+tst.R8 {
			t.Fatalf("got %2X, expected %2X", cpu.A, tst.A+tst.R8)
		}
		testCarries(t, tst)
	})
}

func Test_ADC_A_R8(t *testing.T) {
	cpu := mockCPU()

	type test struct {
		A     uint8
		R8    uint8
		carry uint8
		exp_z uint8
		exp_h uint8
		exp_c uint8
	}

	tests := map[string]test{
		"B":     {0x08, 0x08, 1, 0, 1, 0},
		"C":     {0xF0, 0x11, 0, 0, 0, 1},
		"D":     {0xFF, 0xFF, 1, 0, 1, 1},
		"E":     {0xFF, 0x00, 1, 1, 1, 1},
		"H":     {0x0F, 0x00, 1, 0, 1, 0},
		"L":     {0x57, 0xAD, 1, 0, 1, 1},
		"HLMEM": {0x34, 0x12, 0, 0, 0, 0},
		"A":     {0x80, 0x80, 0, 1, 0, 1},
	}

	testCarries := func(t *testing.T, tst test) {
		if cpu.readNFlag() != 0 {
			t.Errorf("N flag: got %x, expected %x", cpu.readNFlag(), 0)
		}
		if cpu.readZFlag() != tst.exp_z {
			t.Errorf("Z flag: got %x, expected %x", cpu.readZFlag(), tst.exp_z)
		}
		if cpu.readHFlag() != tst.exp_h {
			t.Errorf("H flag: got %x, expected %x", cpu.readHFlag(), tst.exp_h)
		}
		if cpu.readCFlag() != tst.exp_c {
			t.Errorf("C flag: got %x, expected %x", cpu.readCFlag(), tst.exp_c)
		}
	}

	t.Run("B", func(t *testing.T) {
		tst := tests["B"]
		cpu.A = tst.A
		cpu.B = tst.R8
		cpu.setCFlag(tst.carry)
		writeTestProgram(cpu, ADC_A_B_OPCODE)
		cpu.ExecuteInstruction()
		if cpu.A != tst.A+tst.R8+tst.carry {
			t.Fatalf("got %2X, expected %2X", cpu.A, tst.A+tst.R8+tst.carry)
		}
		testCarries(t, tst)
	})

	t.Run("C", func(t *testing.T) {
		tst := tests["C"]
		cpu.A = tst.A
		cpu.C = tst.R8
		cpu.setCFlag(tst.carry)
		writeTestProgram(cpu, ADC_A_C_OPCODE)
		cpu.ExecuteInstruction()
		if cpu.A != tst.A+tst.R8+tst.carry {
			t.Fatalf("got %2X, expected %2X", cpu.A, tst.A+tst.R8+tst.carry)
		}
		testCarries(t, tst)
	})

	t.Run("D", func(t *testing.T) {
		tst := tests["D"]
		cpu.A = tst.A
		cpu.D = tst.R8
		cpu.setCFlag(tst.carry)
		writeTestProgram(cpu, ADC_A_D_OPCODE)
		cpu.ExecuteInstruction()
		if cpu.A != tst.A+tst.R8+tst.carry {
			t.Fatalf("got %2X, expected %2X", cpu.A, tst.A+tst.R8+tst.carry)
		}
		testCarries(t, tst)
	})

	t.Run("E", func(t *testing.T) {
		tst := tests["E"]
		cpu.A = tst.A
		cpu.E = tst.R8
		cpu.setCFlag(tst.carry)
		writeTestProgram(cpu, ADC_A_E_OPCODE)
		cpu.ExecuteInstruction()
		if cpu.A != tst.A+tst.R8+tst.carry {
			t.Fatalf("got %2X, expected %2X", cpu.A, tst.A+tst.R8+tst.carry)
		}
		testCarries(t, tst)
	})

	t.Run("H", func(t *testing.T) {
		tst := tests["H"]
		cpu.A = tst.A
		cpu.H = tst.R8
		cpu.setCFlag(tst.carry)
		writeTestProgram(cpu, ADC_A_H_OPCODE)
		cpu.ExecuteInstruction()
		if cpu.A != tst.A+tst.R8+tst.carry {
			t.Fatalf("got %2X, expected %2X", cpu.A, tst.A+tst.R8+tst.carry)
		}
		testCarries(t, tst)
	})

	t.Run("L", func(t *testing.T) {
		tst := tests["L"]
		cpu.A = tst.A
		cpu.L = tst.R8
		cpu.setCFlag(tst.carry)
		writeTestProgram(cpu, ADC_A_L_OPCODE)
		cpu.ExecuteInstruction()
		if cpu.A != tst.A+tst.R8+tst.carry {
			t.Fatalf("got %2X, expected %2X", cpu.A, tst.A+tst.R8+tst.carry)
		}
		testCarries(t, tst)
	})

	t.Run("A", func(t *testing.T) {
		tst := tests["A"]
		cpu.A = tst.A
		cpu.A = tst.R8
		cpu.setCFlag(tst.carry)
		writeTestProgram(cpu, ADC_A_A_OPCODE)
		cpu.ExecuteInstruction()
		if cpu.A != tst.A+tst.R8+tst.carry {
			t.Fatalf("got %2X, expected %2X", cpu.A, tst.A+tst.R8+tst.carry)
		}
		testCarries(t, tst)
	})

	t.Run("HLMEM", func(t *testing.T) {
		tst := tests["HLMEM"]
		cpu.A = tst.A
		cpu.Mem.Write(cpu.readHL(), tst.R8)
		cpu.setCFlag(tst.carry)
		writeTestProgram(cpu, ADC_A_HLMEM_OPCODE)
		cpu.ExecuteInstruction()
		if cpu.A != tst.A+tst.R8+tst.carry {
			t.Fatalf("got %2X, expected %2X", cpu.A, tst.A+tst.R8+tst.carry)
		}
		testCarries(t, tst)
	})
}

func Test_SUB_A_R8(t *testing.T) {
	cpu := mockCPU()

	type test struct {
		A      uint8
		R8     uint8
		exp_z  uint8
		exp_h  uint8
		exp_c  uint8
		setR8  func(uint8)
		opcode uint8
	}

	tests := map[string]test{
		"B":     {0x08, 0x08, 1, 0, 0, func(v uint8) { cpu.B = v }, SUB_A_B_OPCODE},
		"C":     {0x0F, 0x10, 0, 0, 1, func(v uint8) { cpu.C = v }, SUB_A_C_OPCODE},
		"D":     {0x00, 0x01, 0, 1, 1, func(v uint8) { cpu.D = v }, SUB_A_D_OPCODE},
		"E":     {0xF0, 0x01, 0, 1, 0, func(v uint8) { cpu.E = v }, SUB_A_E_OPCODE},
		"H":     {0x80, 0x08, 0, 1, 0, func(v uint8) { cpu.H = v }, SUB_A_H_OPCODE},
		"L":     {0x57, 0xAD, 0, 1, 1, func(v uint8) { cpu.L = v }, SUB_A_L_OPCODE},
		"HLMEM": {0x34, 0x12, 0, 0, 0, func(v uint8) { cpu.H = v; cpu.Mem.Write(cpu.readHL(), v) }, SUB_A_HLMEM_OPCODE},
		"A":     {0x80, 0x80, 1, 0, 0, func(v uint8) { cpu.A = v }, SUB_A_A_OPCODE},
	}

	for r8, test := range tests {
		t.Run(r8, func(t *testing.T) {
			cpu.A = test.A
			test.setR8(test.R8)
			writeTestProgram(cpu, test.opcode)
			cpu.ExecuteInstruction()
			if cpu.A != test.A-test.R8 {
				t.Fatalf("got %2X, expected %2X", cpu.A, test.A-test.R8)
			}
			if cpu.readNFlag() != 1 {
				t.Errorf("N flag: got %x, expected %x", cpu.readNFlag(), 1)
			}
			if cpu.readZFlag() != test.exp_z {
				t.Errorf("Z flag: got %x, expected %x", cpu.readZFlag(), test.exp_z)
			}
			if cpu.readHFlag() != test.exp_h {
				t.Errorf("H flag: got %x, expected %x", cpu.readHFlag(), test.exp_h)
			}
			if cpu.readCFlag() != test.exp_c {
				t.Errorf("C flag: got %x, expected %x", cpu.readCFlag(), test.exp_c)
			}
		})
	}
}

func Test_SBC_A_R8(t *testing.T) {
	cpu := mockCPU()

	type test struct {
		A      uint8
		R8     uint8
		carry  uint8
		exp_z  uint8
		exp_h  uint8
		exp_c  uint8
		setR8  func(uint8)
		opcode uint8
	}

	tests := map[string]test{
		"B":     {0x08, 0x08, 0, 1, 0, 0, func(v uint8) { cpu.B = v }, SBC_A_B_OPCODE},
		"C":     {0x10, 0x0F, 1, 1, 1, 0, func(v uint8) { cpu.C = v }, SBC_A_C_OPCODE},
		"D":     {0x10, 0x0F, 0, 0, 1, 0, func(v uint8) { cpu.D = v }, SBC_A_D_OPCODE},
		"E":     {0x00, 0x00, 1, 0, 1, 1, func(v uint8) { cpu.E = v }, SBC_A_E_OPCODE},
		"H":     {0x80, 0x80, 0, 1, 0, 0, func(v uint8) { cpu.H = v }, SBC_A_H_OPCODE},
		"L":     {0x57, 0xAD, 1, 0, 1, 1, func(v uint8) { cpu.L = v }, SBC_A_L_OPCODE},
		"HLMEM": {0x34, 0x14, 1, 0, 1, 0, func(v uint8) { cpu.H = v; cpu.Mem.Write(cpu.readHL(), v) }, SBC_A_HLMEM_OPCODE},
		"A":     {0x80, 0x80, 1, 0, 1, 1, func(v uint8) { cpu.A = v }, SBC_A_A_OPCODE},
	}

	for r8, test := range tests {
		t.Run(r8, func(t *testing.T) {
			cpu.A = test.A
			test.setR8(test.R8)
			cpu.setCFlag(test.carry)
			writeTestProgram(cpu, test.opcode)
			cpu.ExecuteInstruction()
			expected := test.A - test.R8 - test.carry
			if cpu.A != expected {
				t.Fatalf("got %2X, expected %2X", cpu.A, expected)
			}
			if cpu.readNFlag() != 1 {
				t.Errorf("N flag: got %x, expected %x", cpu.readNFlag(), 1)
			}
			if cpu.readZFlag() != test.exp_z {
				t.Errorf("Z flag: got %x, expected %x", cpu.readZFlag(), test.exp_z)
			}
			if cpu.readHFlag() != test.exp_h {
				t.Errorf("H flag: got %x, expected %x", cpu.readHFlag(), test.exp_h)
			}
			if cpu.readCFlag() != test.exp_c {
				t.Errorf("C flag: got %x, expected %x", cpu.readCFlag(), test.exp_c)
			}
		})
	}
}

func Test_AND_A_R8(t *testing.T) {
	cpu := mockCPU()

	type test struct {
		A      uint8
		R8     uint8
		exp_z  uint8
		setR8  func(uint8)
		opcode uint8
	}

	tests := map[string]test{
		"B":     {0x08, 0x08, 0, func(v uint8) { cpu.B = v }, AND_A_B_OPCODE},
		"C":     {0x10, 0x0F, 1, func(v uint8) { cpu.C = v }, AND_A_C_OPCODE},
		"D":     {0x10, 0x1F, 0, func(v uint8) { cpu.D = v }, AND_A_D_OPCODE},
		"E":     {0x00, 0x00, 1, func(v uint8) { cpu.E = v }, AND_A_E_OPCODE},
		"H":     {0xAA, 0x55, 1, func(v uint8) { cpu.H = v }, AND_A_H_OPCODE},
		"L":     {0x57, 0xAD, 0, func(v uint8) { cpu.L = v }, AND_A_L_OPCODE},
		"HLMEM": {0x34, 0x14, 0, func(v uint8) { cpu.H = v; cpu.Mem.Write(cpu.readHL(), v) }, AND_A_HLMEM_OPCODE},
		"A":     {0x80, 0x80, 0, func(v uint8) { cpu.A = v }, AND_A_A_OPCODE},
	}

	for r8, test := range tests {
		t.Run(r8, func(t *testing.T) {
			cpu.A = test.A
			test.setR8(test.R8)
			writeTestProgram(cpu, test.opcode)
			cpu.ExecuteInstruction()
			expected := test.A & test.R8
			if cpu.A != expected {
				t.Fatalf("got %2X, expected %2X", cpu.A, expected)
			}
			if cpu.readZFlag() != test.exp_z {
				t.Errorf("Z flag: got %x, expected %x", cpu.readZFlag(), test.exp_z)
			}
			if cpu.readNFlag() != 0 {
				t.Errorf("N flag: got %x, expected %x", cpu.readNFlag(), 0)
			}
			if cpu.readHFlag() != 1 {
				t.Errorf("H flag: got %x, expected %x", cpu.readHFlag(), 1)
			}
			if cpu.readCFlag() != 0 {
				t.Errorf("C flag: got %x, expected %x", cpu.readCFlag(), 0)
			}
		})
	}
}

func Test_XOR_A_R8(t *testing.T) {
	cpu := mockCPU()

	type test struct {
		A      uint8
		R8     uint8
		exp_z  uint8
		setR8  func(uint8)
		opcode uint8
	}

	tests := map[string]test{
		"B":     {0x08, 0x08, 1, func(v uint8) { cpu.B = v }, XOR_A_B_OPCODE},
		"C":     {0x10, 0x0F, 0, func(v uint8) { cpu.C = v }, XOR_A_C_OPCODE},
		"D":     {0x10, 0x1F, 0, func(v uint8) { cpu.D = v }, XOR_A_D_OPCODE},
		"E":     {0x00, 0x00, 1, func(v uint8) { cpu.E = v }, XOR_A_E_OPCODE},
		"H":     {0xAA, 0x55, 0, func(v uint8) { cpu.H = v }, XOR_A_H_OPCODE},
		"L":     {0x57, 0xAD, 0, func(v uint8) { cpu.L = v }, XOR_A_L_OPCODE},
		"HLMEM": {0x34, 0x14, 0, func(v uint8) { cpu.H = v; cpu.Mem.Write(cpu.readHL(), v) }, XOR_A_HLMEM_OPCODE},
		"A":     {0x80, 0x80, 1, func(v uint8) { cpu.A = v }, XOR_A_A_OPCODE},
	}

	for r8, test := range tests {
		t.Run(r8, func(t *testing.T) {
			cpu.A = test.A
			test.setR8(test.R8)
			writeTestProgram(cpu, test.opcode)
			cpu.ExecuteInstruction()
			expected := test.A ^ test.R8
			if cpu.A != expected {
				t.Fatalf("got %2X, expected %2X", cpu.A, expected)
			}
			if cpu.readZFlag() != test.exp_z {
				t.Errorf("Z flag: got %x, expected %x", cpu.readZFlag(), test.exp_z)
			}
			if cpu.readNFlag() != 0 {
				t.Errorf("N flag: got %x, expected %x", cpu.readNFlag(), 0)
			}
			if cpu.readHFlag() != 0 {
				t.Errorf("H flag: got %x, expected %x", cpu.readHFlag(), 0)
			}
			if cpu.readCFlag() != 0 {
				t.Errorf("C flag: got %x, expected %x", cpu.readCFlag(), 0)
			}
		})
	}
}

func Test_OR_A_R8(t *testing.T) {
	cpu := mockCPU()

	type test struct {
		A      uint8
		R8     uint8
		exp_z  uint8
		setR8  func(uint8)
		opcode uint8
	}

	tests := map[string]test{
		"B":     {0x08, 0x08, 0, func(v uint8) { cpu.B = v }, OR_A_B_OPCODE},
		"C":     {0x10, 0x0F, 0, func(v uint8) { cpu.C = v }, OR_A_C_OPCODE},
		"D":     {0x10, 0x1F, 0, func(v uint8) { cpu.D = v }, OR_A_D_OPCODE},
		"E":     {0x00, 0x00, 1, func(v uint8) { cpu.E = v }, OR_A_E_OPCODE},
		"H":     {0xAA, 0x55, 0, func(v uint8) { cpu.H = v }, OR_A_H_OPCODE},
		"L":     {0x57, 0xAD, 0, func(v uint8) { cpu.L = v }, OR_A_L_OPCODE},
		"HLMEM": {0x34, 0x14, 0, func(v uint8) { cpu.H = v; cpu.Mem.Write(cpu.readHL(), v) }, OR_A_HLMEM_OPCODE},
		"A":     {0x80, 0x80, 0, func(v uint8) { cpu.A = v }, OR_A_A_OPCODE},
	}

	for r8, test := range tests {
		t.Run(r8, func(t *testing.T) {
			cpu.A = test.A
			test.setR8(test.R8)
			writeTestProgram(cpu, test.opcode)
			cpu.ExecuteInstruction()
			expected := test.A | test.R8
			if cpu.A != expected {
				t.Fatalf("got %2X, expected %2X", cpu.A, expected)
			}
			if cpu.readZFlag() != test.exp_z {
				t.Errorf("Z flag: got %x, expected %x", cpu.readZFlag(), test.exp_z)
			}
			if cpu.readNFlag() != 0 {
				t.Errorf("N flag: got %x, expected %x", cpu.readNFlag(), 0)
			}
			if cpu.readHFlag() != 0 {
				t.Errorf("H flag: got %x, expected %x", cpu.readHFlag(), 0)
			}
			if cpu.readCFlag() != 0 {
				t.Errorf("C flag: got %x, expected %x", cpu.readCFlag(), 0)
			}
		})
	}
}

func Test_CP_A_R8(t *testing.T) {
	cpu := mockCPU()

	type test struct {
		A      uint8
		R8     uint8
		exp_z  uint8
		exp_h  uint8
		exp_c  uint8
		setR8  func(uint8)
		opcode uint8
	}

	tests := map[string]test{
		"B":     {0x08, 0x08, 1, 0, 0, func(v uint8) { cpu.B = v }, CP_A_B_OPCODE},
		"C":     {0x0F, 0x10, 0, 0, 1, func(v uint8) { cpu.C = v }, CP_A_C_OPCODE},
		"D":     {0x00, 0x01, 0, 1, 1, func(v uint8) { cpu.D = v }, CP_A_D_OPCODE},
		"E":     {0xF0, 0x01, 0, 1, 0, func(v uint8) { cpu.E = v }, CP_A_E_OPCODE},
		"H":     {0x80, 0x08, 0, 1, 0, func(v uint8) { cpu.H = v }, CP_A_H_OPCODE},
		"L":     {0x57, 0xAD, 0, 1, 1, func(v uint8) { cpu.L = v }, CP_A_L_OPCODE},
		"HLMEM": {0x34, 0x12, 0, 0, 0, func(v uint8) { cpu.H = v; cpu.Mem.Write(cpu.readHL(), v) }, CP_A_HLMEM_OPCODE},
		"A":     {0x80, 0x80, 1, 0, 0, func(v uint8) { cpu.A = v }, CP_A_A_OPCODE},
	}

	for r8, test := range tests {
		t.Run(r8, func(t *testing.T) {
			cpu.A = test.A
			test.setR8(test.R8)
			writeTestProgram(cpu, test.opcode)
			cpu.ExecuteInstruction()
			if cpu.A != test.A {
				t.Fatalf("got %2X, expected %2X", cpu.A, test.A)
			}
			if cpu.readNFlag() != 1 {
				t.Errorf("N flag: got %x, expected %x", cpu.readNFlag(), 1)
			}
			if cpu.readZFlag() != test.exp_z {
				t.Errorf("Z flag: got %x, expected %x", cpu.readZFlag(), test.exp_z)
			}
			if cpu.readHFlag() != test.exp_h {
				t.Errorf("H flag: got %x, expected %x", cpu.readHFlag(), test.exp_h)
			}
			if cpu.readCFlag() != test.exp_c {
				t.Errorf("C flag: got %x, expected %x", cpu.readCFlag(), test.exp_c)
			}
		})
	}
}

func Test_ADD_A_N8(t *testing.T) {
	cpu := mockCPU()

	var A uint8 = 0x57
	var n8 uint8 = 0xAD
	var exp_z uint8 = 0
	var exp_h uint8 = 1
	var exp_c uint8 = 1

	cpu.A = A
	writeTestProgram(cpu, ADD_A_N8_OPCODE, n8)
	cpu.ExecuteInstruction()
	if cpu.A != A+n8 {
		t.Fatalf("got %2X, expected %2X", cpu.A, A+n8)
	}
	if cpu.readNFlag() != 0 {
		t.Errorf("N flag: got %x, expected %x", cpu.readNFlag(), 0)
	}
	if cpu.readZFlag() != exp_z {
		t.Errorf("Z flag: got %x, expected %x", cpu.readZFlag(), exp_z)
	}
	if cpu.readHFlag() != exp_h {
		t.Errorf("H flag: got %x, expected %x", cpu.readHFlag(), exp_h)
	}
	if cpu.readCFlag() != exp_c {
		t.Errorf("C flag: got %x, expected %x", cpu.readCFlag(), exp_c)
	}
}

func Test_ADC_A_N8(t *testing.T) {
	cpu := mockCPU()

	var A uint8 = 0x08
	var n8 uint8 = 0x08
	var carry uint8 = 1
	var exp_z uint8 = 0
	var exp_h uint8 = 1
	var exp_c uint8 = 0

	cpu.A = A
	cpu.setCFlag(carry)
	writeTestProgram(cpu, ADC_A_N8_OPCODE, n8)
	cpu.ExecuteInstruction()
	if cpu.A != A+n8+carry {
		t.Fatalf("got %2X, expected %2X", cpu.A, A+n8+carry)
	}
	if cpu.readNFlag() != 0 {
		t.Errorf("N flag: got %x, expected %x", cpu.readNFlag(), 0)
	}
	if cpu.readZFlag() != exp_z {
		t.Errorf("Z flag: got %x, expected %x", cpu.readZFlag(), exp_z)
	}
	if cpu.readHFlag() != exp_h {
		t.Errorf("H flag: got %x, expected %x", cpu.readHFlag(), exp_h)
	}
	if cpu.readCFlag() != exp_c {
		t.Errorf("C flag: got %x, expected %x", cpu.readCFlag(), exp_c)
	}
}

func Test_SUB_A_N8(t *testing.T) {
	cpu := mockCPU()

	var A uint8 = 0x80
	var n8 uint8 = 0x08
	var exp_z uint8 = 0
	var exp_h uint8 = 1
	var exp_c uint8 = 0

	cpu.A = A
	writeTestProgram(cpu, SUB_A_N8_OPCODE, n8)
	cpu.ExecuteInstruction()
	if cpu.A != A-n8 {
		t.Fatalf("got %2X, expected %2X", cpu.A, A-n8)
	}
	if cpu.readNFlag() != 1 {
		t.Errorf("N flag: got %x, expected %x", cpu.readNFlag(), 1)
	}
	if cpu.readZFlag() != exp_z {
		t.Errorf("Z flag: got %x, expected %x", cpu.readZFlag(), exp_z)
	}
	if cpu.readHFlag() != exp_h {
		t.Errorf("H flag: got %x, expected %x", cpu.readHFlag(), exp_h)
	}
	if cpu.readCFlag() != exp_c {
		t.Errorf("C flag: got %x, expected %x", cpu.readCFlag(), exp_c)
	}
}

func Test_SBC_A_N8(t *testing.T) {
	cpu := mockCPU()

	var A uint8 = 0x10
	var n8 uint8 = 0x0F
	var carry uint8 = 1
	var exp_z uint8 = 1
	var exp_h uint8 = 1
	var exp_c uint8 = 0

	cpu.A = A
	cpu.setCFlag(carry)
	writeTestProgram(cpu, SBC_A_N8_OPCODE, n8)
	cpu.ExecuteInstruction()
	if cpu.A != A-n8-carry {
		t.Fatalf("got %2X, expected %2X", cpu.A, A-n8-carry)
	}
	if cpu.readNFlag() != 1 {
		t.Errorf("N flag: got %x, expected %x", cpu.readNFlag(), 1)
	}
	if cpu.readZFlag() != exp_z {
		t.Errorf("Z flag: got %x, expected %x", cpu.readZFlag(), exp_z)
	}
	if cpu.readHFlag() != exp_h {
		t.Errorf("H flag: got %x, expected %x", cpu.readHFlag(), exp_h)
	}
	if cpu.readCFlag() != exp_c {
		t.Errorf("C flag: got %x, expected %x", cpu.readCFlag(), exp_c)
	}
}

func Test_AND_A_N8(t *testing.T) {
	cpu := mockCPU()

	var A uint8 = 0xAA
	var n8 uint8 = 0x55
	var exp_z uint8 = 1

	cpu.A = A
	writeTestProgram(cpu, AND_A_N8_OPCODE, n8)
	cpu.ExecuteInstruction()
	expected := A & n8
	if cpu.A != expected {
		t.Fatalf("got %2X, expected %2X", cpu.A, expected)
	}
	if cpu.readZFlag() != exp_z {
		t.Errorf("Z flag: got %x, expected %x", cpu.readZFlag(), exp_z)
	}
	if cpu.readNFlag() != 0 {
		t.Errorf("N flag: got %x, expected %x", cpu.readNFlag(), 0)
	}
	if cpu.readHFlag() != 1 {
		t.Errorf("H flag: got %x, expected %x", cpu.readHFlag(), 1)
	}
	if cpu.readCFlag() != 0 {
		t.Errorf("C flag: got %x, expected %x", cpu.readCFlag(), 0)
	}
}

func Test_XOR_A_N8(t *testing.T) {
	cpu := mockCPU()

	var A uint8 = 0x80
	var n8 uint8 = 0x80
	var exp_z uint8 = 1

	cpu.A = A
	writeTestProgram(cpu, XOR_A_N8_OPCODE, n8)
	cpu.ExecuteInstruction()
	expected := A ^ n8
	if cpu.A != expected {
		t.Fatalf("got %2X, expected %2X", cpu.A, expected)
	}
	if cpu.readZFlag() != exp_z {
		t.Errorf("Z flag: got %x, expected %x", cpu.readZFlag(), exp_z)
	}
	if cpu.readNFlag() != 0 {
		t.Errorf("N flag: got %x, expected %x", cpu.readNFlag(), 0)
	}
	if cpu.readHFlag() != 0 {
		t.Errorf("H flag: got %x, expected %x", cpu.readHFlag(), 0)
	}
	if cpu.readCFlag() != 0 {
		t.Errorf("C flag: got %x, expected %x", cpu.readCFlag(), 0)
	}
}

func Test_OR_A_N8(t *testing.T) {
	cpu := mockCPU()

	var A uint8 = 0x80
	var n8 uint8 = 0x80
	var exp_z uint8 = 0

	cpu.A = A
	writeTestProgram(cpu, OR_A_N8_OPCODE, n8)
	cpu.ExecuteInstruction()
	expected := A | n8
	if cpu.A != expected {
		t.Fatalf("got %2X, expected %2X", cpu.A, expected)
	}
	if cpu.readZFlag() != exp_z {
		t.Errorf("Z flag: got %x, expected %x", cpu.readZFlag(), exp_z)
	}
	if cpu.readNFlag() != 0 {
		t.Errorf("N flag: got %x, expected %x", cpu.readNFlag(), 0)
	}
	if cpu.readHFlag() != 0 {
		t.Errorf("H flag: got %x, expected %x", cpu.readHFlag(), 0)
	}
	if cpu.readCFlag() != 0 {
		t.Errorf("C flag: got %x, expected %x", cpu.readCFlag(), 0)
	}
}

func Test_CP_A_N8(t *testing.T) {
	cpu := mockCPU()

	var A uint8 = 0x10
	var n8 uint8 = 0x0F
	var exp_z uint8 = 0
	var exp_h uint8 = 1
	var exp_c uint8 = 0

	cpu.A = A
	writeTestProgram(cpu, CP_A_N8_OPCODE, n8)
	cpu.ExecuteInstruction()
	if cpu.A != A {
		t.Fatalf("got %2X, expected %2X", cpu.A, A)
	}
	if cpu.readNFlag() != 1 {
		t.Errorf("N flag: got %x, expected %x", cpu.readNFlag(), 1)
	}
	if cpu.readZFlag() != exp_z {
		t.Errorf("Z flag: got %x, expected %x", cpu.readZFlag(), exp_z)
	}
	if cpu.readHFlag() != exp_h {
		t.Errorf("H flag: got %x, expected %x", cpu.readHFlag(), exp_h)
	}
	if cpu.readCFlag() != exp_c {
		t.Errorf("C flag: got %x, expected %x", cpu.readCFlag(), exp_c)
	}
}

func Test_POP_R16STK(t *testing.T) {
	cpu := mockCPU()

	t.Run("BC", func(t *testing.T) {
		addr := uint16(0x4321)
		expected_BC := addr
		expected_SP := cpu.SP

		// Write stack pointer
		cpu.SP -= 2
		cpu.Mem.WriteWord(cpu.SP, addr)

		writeTestProgram(cpu, POP_BC_OPCODE)
		cpu.ExecuteInstruction()

		if cpu.readBC() != uint16(expected_BC) {
			t.Errorf("got BC=%04X, expected %04X", cpu.readBC(), expected_BC)
		}
		if cpu.SP != uint16(expected_SP) {
			t.Errorf("got SP=%04X, expected %04X", cpu.SP, expected_SP)
		}
	})

	t.Run("DE", func(t *testing.T) {
		addr := uint16(0x1234)
		expected_DE := addr
		expected_SP := cpu.SP

		// Write stack pointer
		cpu.SP -= 2
		cpu.Mem.WriteWord(cpu.SP, addr)

		writeTestProgram(cpu, POP_DE_OPCODE)
		cpu.ExecuteInstruction()

		if cpu.readDE() != uint16(expected_DE) {
			t.Errorf("got DE=%04X, expected %04X", cpu.readDE(), expected_DE)
		}
		if cpu.SP != uint16(expected_SP) {
			t.Errorf("got SP=%04X, expected %04X", cpu.SP, expected_SP)
		}
	})

	t.Run("HL", func(t *testing.T) {
		addr := uint16(0x1111)
		expected_HL := addr
		expected_SP := cpu.SP

		// Write stack pointer
		cpu.SP -= 2
		cpu.Mem.WriteWord(cpu.SP, addr)

		writeTestProgram(cpu, POP_HL_OPCODE)
		cpu.ExecuteInstruction()

		if cpu.readHL() != uint16(expected_HL) {
			t.Errorf("got HL=%04X, expected %04X", cpu.readHL(), expected_HL)
		}
		if cpu.SP != uint16(expected_SP) {
			t.Errorf("got SP=%04X, expected %04X", cpu.SP, expected_SP)
		}
	})

	t.Run("AF", func(t *testing.T) {
		addr := uint16(0x56F0)
		expected_AF := addr
		expected_SP := cpu.SP

		// Write stack pointer
		cpu.SP -= 2
		cpu.Mem.WriteWord(cpu.SP, addr)

		writeTestProgram(cpu, POP_AF_OPCODE)
		cpu.ExecuteInstruction()

		if cpu.readAF() != uint16(expected_AF) {
			t.Errorf("got AF=%04X, expected %04X", cpu.readAF(), expected_AF)
		}
		if cpu.SP != uint16(expected_SP) {
			t.Errorf("got SP=%04X, expected %04X", cpu.SP, expected_SP)
		}

		// Test registers
		got, expected := cpu.readZFlag(), readBit(uint8(addr), Z_FLAG_BIT)
		if got != expected {
			t.Errorf("Z flag: got %x, expected %x", got, expected)
		}
		got, expected = cpu.readNFlag(), readBit(uint8(addr), N_FLAG_BIT)
		if got != expected {
			t.Errorf("N flag: got %x, expected %x", got, expected)
		}
		got, expected = cpu.readHFlag(), readBit(uint8(addr), H_FLAG_BIT)
		if got != expected {
			t.Errorf("H flag: got %x, expected %x", got, expected)
		}
		got, expected = cpu.readCFlag(), readBit(uint8(addr), C_FLAG_BIT)
		if got != expected {
			t.Errorf("C flag: got %x, expected %x", got, expected)
		}
	})
}

func Test_PUSH_R16STK(t *testing.T) {
	cpu := mockCPU()

	t.Run("BC", func(t *testing.T) {
		addr := uint16(0x4321)
		expected_SP := cpu.SP - 2
		cpu.writeBC(addr)

		writeTestProgram(cpu, PUSH_BC_OPCODE)
		cpu.ExecuteInstruction()

		if cpu.Mem.ReadWord(cpu.SP) != addr {
			t.Errorf("got [SP]=%04X, expected %04X", cpu.Mem.ReadWord(cpu.SP), addr)
		}
		if cpu.SP != uint16(expected_SP) {
			t.Errorf("got SP=%04X, expected %04X", cpu.SP, expected_SP)
		}
	})

	t.Run("DE", func(t *testing.T) {
		addr := uint16(0x1234)
		expected_SP := cpu.SP - 2
		cpu.writeDE(addr)

		writeTestProgram(cpu, PUSH_DE_OPCODE)
		cpu.ExecuteInstruction()

		if cpu.Mem.ReadWord(cpu.SP) != addr {
			t.Errorf("got [SP]=%04X, expected %04X", cpu.Mem.ReadWord(cpu.SP), addr)
		}
		if cpu.SP != uint16(expected_SP) {
			t.Errorf("got SP=%04X, expected %04X", cpu.SP, expected_SP)
		}
	})

	t.Run("HL", func(t *testing.T) {
		addr := uint16(0x1111)
		expected_SP := cpu.SP - 2
		cpu.writeHL(addr)

		writeTestProgram(cpu, PUSH_HL_OPCODE)
		cpu.ExecuteInstruction()

		if cpu.Mem.ReadWord(cpu.SP) != addr {
			t.Errorf("got [SP]=%04X, expected %04X", cpu.Mem.ReadWord(cpu.SP), addr)
		}
		if cpu.SP != uint16(expected_SP) {
			t.Errorf("got SP=%04X, expected %04X", cpu.SP, expected_SP)
		}
	})

	t.Run("AF", func(t *testing.T) {
		addr := uint16(0x2220)
		expected_SP := cpu.SP - 2
		cpu.writeAF(addr)

		writeTestProgram(cpu, PUSH_AF_OPCODE)
		cpu.ExecuteInstruction()

		if cpu.Mem.ReadWord(cpu.SP) != addr {
			t.Errorf("got [SP]=%04X, expected %04X", cpu.Mem.ReadWord(cpu.SP), addr)
		}
		if cpu.SP != uint16(expected_SP) {
			t.Errorf("got SP=%04X, expected %04X", cpu.SP, expected_SP)
		}
	})
}

func Test_RET_COND(t *testing.T) {
	cpu := mockCPU()

	// Prepare flags for each condition
	bool_to_int := map[bool]uint8{false: 0, true: 1}
	conditions := map[string]func(bool){
		"Z":  func(cc bool) { cpu.setZFlag(bool_to_int[cc]) },
		"NZ": func(cc bool) { cpu.setZFlag(bool_to_int[!cc]) },
		"C":  func(cc bool) { cpu.setCFlag(bool_to_int[cc]) },
		"NC": func(cc bool) { cpu.setCFlag(bool_to_int[!cc]) },
	}
	opcodes := map[string]uint8{
		"Z":  RET_Z_OPCODE,
		"NZ": RET_NZ_OPCODE,
		"C":  RET_C_OPCODE,
		"NC": RET_NC_OPCODE,
	}

	for cond, set_flag := range conditions {
		addr := uint16(rand.Intn(0x10000))
		// Write stack pointer
		cpu.SP -= 2
		cpu.Mem.WriteWord(cpu.SP, addr)

		t.Run(cond+"_unmet", func(t *testing.T) {
			set_flag(false)

			expected_PC := int(cpu.PC) + OPCODES_BYTES[opcodes[cond]]
			expected_cycles := cpu.cycles + OPCODES_CYCLES[opcodes[cond]]
			expected_SP := cpu.SP

			writeTestProgram(cpu, opcodes[cond])
			cpu.ExecuteInstruction()

			if cpu.PC != uint16(expected_PC) {
				t.Errorf("unmet: got PC=%04X, expected %04X", cpu.PC, expected_PC)
			}
			if cpu.SP != uint16(expected_SP) {
				t.Errorf("unmet: got SP=%04X, expected %04X", cpu.SP, expected_SP)
			}
			if cpu.cycles != expected_cycles {
				t.Errorf("unmet: got cycles=%v, expected %v", cpu.cycles, expected_cycles)
			}
		})

		t.Run(cond+"_met", func(t *testing.T) {
			set_flag(true)

			expected_PC := addr
			expected_cycles := cpu.cycles + OPCODES_CYCLES_BRANCH[opcodes[cond]]
			expected_SP := cpu.SP + 2

			writeTestProgram(cpu, opcodes[cond])
			cpu.ExecuteInstruction()

			if cpu.PC != uint16(expected_PC) {
				t.Errorf("met: got PC=%04X, expected %04X", cpu.PC, expected_PC)
			}
			if cpu.SP != uint16(expected_SP) {
				t.Errorf("met: got SP=%04X, expected %04X", cpu.SP, expected_SP)
			}
			if cpu.cycles != expected_cycles {
				t.Errorf("met: got cycles=%v, expected %v", cpu.cycles, expected_cycles)
			}
		})
	}
}

func Test_RET(t *testing.T) {
	cpu := mockCPU()

	addr := uint16(0x4321)
	expected_PC := addr
	expected_SP := cpu.SP

	// Write stack pointer
	cpu.SP -= 2
	cpu.Mem.WriteWord(cpu.SP, addr)

	writeTestProgram(cpu, RET_OPCODE)
	cpu.ExecuteInstruction()

	if cpu.PC != uint16(expected_PC) {
		t.Errorf("got PC=%04X, expected %04X", cpu.PC, expected_PC)
	}
	if cpu.SP != uint16(expected_SP) {
		t.Errorf("got SP=%04X, expected %04X", cpu.SP, expected_SP)
	}
}

func Test_RETI(t *testing.T) {
	cpu := mockCPU()
	cpu.IME = false

	addr := uint16(0x4321)
	expected_PC := addr
	expected_SP := cpu.SP

	// Write stack pointer
	cpu.SP -= 2
	cpu.Mem.WriteWord(cpu.SP, addr)

	writeTestProgram(cpu, RETI_OPCODE)
	cpu.ExecuteInstruction()

	if cpu.PC != uint16(expected_PC) {
		t.Errorf("got PC=%04X, expected %04X", cpu.PC, expected_PC)
	}
	if cpu.SP != uint16(expected_SP) {
		t.Errorf("got SP=%04X, expected %04X", cpu.SP, expected_SP)
	}
	if cpu.IME != true {
		t.Error("Interrupts not enabled")
	}
}

func Test_JP_COND_N16(t *testing.T) {
	cpu := mockCPU()

	// Prepare flags for each condition
	bool_to_int := map[bool]uint8{false: 0, true: 1}
	conditions := map[string]func(bool){
		"Z":  func(cc bool) { cpu.setZFlag(bool_to_int[cc]) },
		"NZ": func(cc bool) { cpu.setZFlag(bool_to_int[!cc]) },
		"C":  func(cc bool) { cpu.setCFlag(bool_to_int[cc]) },
		"NC": func(cc bool) { cpu.setCFlag(bool_to_int[!cc]) },
	}
	opcodes := map[string]uint8{
		"Z":  JP_Z_N16_OPCODE,
		"NZ": JP_NZ_N16_OPCODE,
		"C":  JP_C_N16_OPCODE,
		"NC": JP_NC_N16_OPCODE,
	}

	for cond, set_flag := range conditions {
		addr := uint16(rand.Intn(0x10000))

		t.Run(cond+"_unmet", func(t *testing.T) {
			set_flag(false)

			expected_PC := int(cpu.PC) + OPCODES_BYTES[opcodes[cond]]
			expected_cycles := cpu.cycles + OPCODES_CYCLES[opcodes[cond]]

			writeTestProgram(cpu, opcodes[cond], uint8(addr), uint8(addr>>8))
			cpu.ExecuteInstruction()

			if cpu.PC != uint16(expected_PC) {
				t.Errorf("unmet: got PC=%04X, expected %04X", cpu.PC, expected_PC)
			}
			if cpu.cycles != expected_cycles {
				t.Errorf("unmet: got cycles=%v, expected %v", cpu.cycles, expected_cycles)
			}
		})

		t.Run(cond+"_met", func(t *testing.T) {
			set_flag(true)

			expected_PC := addr
			expected_cycles := cpu.cycles + OPCODES_CYCLES_BRANCH[opcodes[cond]]

			writeTestProgram(cpu, opcodes[cond], uint8(addr), uint8(addr>>8))
			cpu.ExecuteInstruction()

			if cpu.PC != uint16(expected_PC) {
				t.Errorf("met: got PC=%04X, expected %04X", cpu.PC, expected_PC)
			}
			if cpu.cycles != expected_cycles {
				t.Errorf("met: got cycles=%v, expected %v", cpu.cycles, expected_cycles)
			}
		})
	}
}

func Test_JP_N16(t *testing.T) {
	cpu := mockCPU()

	addr := uint16(rand.Intn(0x10000))

	expected_PC := addr

	writeTestProgram(cpu, JP_N16_OPCODE, uint8(addr), uint8(addr>>8))
	cpu.ExecuteInstruction()

	if cpu.PC != uint16(expected_PC) {
		t.Errorf("got PC=%04X, expected %04X", cpu.PC, expected_PC)
	}
}

func Test_JP_HL(t *testing.T) {
	cpu := mockCPU()

	addr := uint16(rand.Intn(0x10000))

	expected_PC := addr

	writeTestProgram(cpu, JP_HL_OPCODE)
	cpu.writeHL(addr)
	cpu.ExecuteInstruction()

	if cpu.PC != uint16(expected_PC) {
		t.Errorf("got PC=%04X, expected %04X", cpu.PC, expected_PC)
	}
}

func Test_CALL_COND_N16(t *testing.T) {
	cpu := mockCPU()

	// Prepare flags for each condition
	bool_to_int := map[bool]uint8{false: 0, true: 1}
	conditions := map[string]func(bool){
		"Z":  func(cc bool) { cpu.setZFlag(bool_to_int[cc]) },
		"NZ": func(cc bool) { cpu.setZFlag(bool_to_int[!cc]) },
		"C":  func(cc bool) { cpu.setCFlag(bool_to_int[cc]) },
		"NC": func(cc bool) { cpu.setCFlag(bool_to_int[!cc]) },
	}
	opcodes := map[string]uint8{
		"Z":  CALL_Z_N16_OPCODE,
		"NZ": CALL_NZ_N16_OPCODE,
		"C":  CALL_C_N16_OPCODE,
		"NC": CALL_NC_N16_OPCODE,
	}

	for cond, set_flag := range conditions {
		addr := uint16(rand.Intn(0x10000))

		t.Run(cond+"_unmet", func(t *testing.T) {
			set_flag(false)

			expected_PC := cpu.PC + uint16(OPCODES_BYTES[opcodes[cond]])
			expected_SP := cpu.SP
			expected_stack := cpu.Mem.ReadWord(cpu.SP)

			writeTestProgram(cpu, opcodes[cond], uint8(addr), uint8(addr>>8))
			cpu.ExecuteInstruction()

			if cpu.PC != uint16(expected_PC) {
				t.Errorf("got PC=%04X, expected %04X", cpu.PC, expected_PC)
			}
			if cpu.SP != uint16(expected_SP) {
				t.Errorf("got SP=%04X, expected %04X", cpu.SP, expected_SP)
			}
			if cpu.Mem.ReadWord(cpu.SP) != uint16(expected_stack) {
				t.Errorf("got stack=%04X, expected %04X", cpu.Mem.ReadWord(cpu.SP), expected_stack)
			}
		})

		t.Run(cond+"_met", func(t *testing.T) {
			set_flag(true)

			expected_PC := addr
			expected_SP := cpu.SP - 2
			expected_stack := cpu.PC + uint16(OPCODES_BYTES[opcodes[cond]])

			writeTestProgram(cpu, opcodes[cond], uint8(addr), uint8(addr>>8))
			cpu.ExecuteInstruction()

			if cpu.PC != uint16(expected_PC) {
				t.Errorf("got PC=%04X, expected %04X", cpu.PC, expected_PC)
			}
			if cpu.SP != uint16(expected_SP) {
				t.Errorf("got SP=%04X, expected %04X", cpu.SP, expected_SP)
			}
			if cpu.Mem.ReadWord(cpu.SP) != uint16(expected_stack) {
				t.Errorf("got stack=%04X, expected %04X", cpu.Mem.ReadWord(cpu.SP), expected_stack)
			}
		})
	}
}

func Test_CALL_N16(t *testing.T) {
	cpu := mockCPU()

	addr := uint16(rand.Intn(0x10000))

	expected_PC := addr
	expected_SP := cpu.SP - 2
	expected_stack := cpu.PC + uint16(OPCODES_BYTES[CALL_N16_OPCODE])

	writeTestProgram(cpu, CALL_N16_OPCODE, uint8(addr), uint8(addr>>8))
	cpu.ExecuteInstruction()

	if cpu.PC != uint16(expected_PC) {
		t.Errorf("got PC=%04X, expected %04X", cpu.PC, expected_PC)
	}
	if cpu.SP != uint16(expected_SP) {
		t.Errorf("got SP=%04X, expected %04X", cpu.SP, expected_SP)
	}
	if cpu.Mem.ReadWord(cpu.SP) != uint16(expected_stack) {
		t.Errorf("got stack=%04X, expected %04X", cpu.Mem.ReadWord(cpu.SP), expected_stack)
	}
}

func Test_TST_VEC_N16(t *testing.T) {
	cpu := mockCPU()

	opcodes := map[string]uint8{
		"00": RST_00_OPCODE,
		"08": RST_08_OPCODE,
		"10": RST_10_OPCODE,
		"18": RST_18_OPCODE,
		"20": RST_20_OPCODE,
		"28": RST_28_OPCODE,
		"30": RST_30_OPCODE,
		"38": RST_38_OPCODE,
	}

	for vec, opcode := range opcodes {
		addr64, _ := strconv.ParseUint(vec, 16, 8)
		addr := uint8(addr64)

		t.Run(vec, func(t *testing.T) {
			expected_PC := addr
			expected_SP := cpu.SP - 2

			writeTestProgram(cpu, opcode)
			cpu.ExecuteInstruction()

			if cpu.PC != uint16(expected_PC) {
				t.Errorf("got PC=%04X, expected %04X", cpu.PC, expected_PC)
			}
			if cpu.SP != uint16(expected_SP) {
				t.Errorf("got SP=%04X, expected %04X", cpu.SP, expected_SP)
			}
		})
	}
}

func Test_LDH_C_A(t *testing.T) {
	cpu := mockCPU()
	var A, C uint8 = 0xD1, 0x12

	cpu.A = A
	cpu.C = C
	writeTestProgram(cpu, LDH_C_A_OPCODE)
	cpu.ExecuteInstruction()

	got, expected := cpu.Mem.Read(0xFF00+uint16(cpu.C)), A
	if got != expected {
		t.Fatalf("got %02X, expected %02X", got, expected)
	}
}

func Test_LDH_A_C(t *testing.T) {
	cpu := mockCPU()
	var value, C uint8 = 0xD1, 0x12

	cpu.Mem.Write(0xFF00+uint16(C), value)
	cpu.C = C
	writeTestProgram(cpu, LDH_A_C_OPCODE)
	cpu.ExecuteInstruction()

	got, expected := cpu.A, value
	if got != expected {
		t.Fatalf("got %02X, expected %02X", got, expected)
	}
}

func Test_LDH_N8_A(t *testing.T) {
	cpu := mockCPU()
	var A, offset uint8 = 0xD1, 0x12

	cpu.A = A
	writeTestProgram(cpu, LDH_N8_A_OPCODE, offset)
	cpu.ExecuteInstruction()

	got, expected := cpu.Mem.Read(0xFF00+uint16(offset)), A
	if got != expected {
		t.Fatalf("got %02X, expected %02X", got, expected)
	}
}

func Test_LDH_A_N8(t *testing.T) {
	cpu := mockCPU()
	var value, offset uint8 = 0xD1, 0x12

	cpu.Mem.Write(0xFF00+uint16(offset), value)
	writeTestProgram(cpu, LDH_A_N8_OPCODE, offset)
	cpu.ExecuteInstruction()

	got, expected := cpu.A, value
	if got != expected {
		t.Fatalf("got %02X, expected %02X", got, expected)
	}
}

func Test_LD_N16_A(t *testing.T) {
	cpu := mockCPU()
	var A uint8 = 0xD1
	var addr uint16 = 0xABCD

	cpu.A = A
	writeTestProgram(cpu, LD_N16_A_OPCODE, uint8(addr), uint8(addr>>8))
	cpu.ExecuteInstruction()

	got, expected := cpu.Mem.Read(addr), A
	if got != expected {
		t.Fatalf("got %02X, expected %02X", got, expected)
	}
}

func Test_LD_A_N16(t *testing.T) {
	cpu := mockCPU()
	var value uint8 = 0xD1
	var addr uint16 = 0xDCBA

	cpu.Mem.Write(addr, value)
	writeTestProgram(cpu, LD_A_N16_OPCODE, uint8(addr), uint8(addr>>8))
	cpu.ExecuteInstruction()

	got, expected := cpu.A, value
	if got != expected {
		t.Fatalf("got %02X, expected %02X", got, expected)
	}
}

func Test_ADD_SP_E8(t *testing.T) {
	cpu := mockCPU()

	tests := map[string]struct {
		SP    uint16
		e8    int8
		exp_H uint8
		exp_C uint8
	}{
		"neg-with-carry":    {SP: 0xFFFE, e8: -1, exp_H: 1, exp_C: 1},
		"neg-without-carry": {SP: 0xFF00, e8: -1, exp_H: 0, exp_C: 0},
		"pos-with-carry":    {SP: 0xF0FF, e8: 0x0F, exp_H: 1, exp_C: 1},
		"pos-without-carry": {SP: 0xFF00, e8: 0x0F, exp_H: 0, exp_C: 0},
		"overflow":          {SP: 0xFFFF, e8: 1, exp_H: 1, exp_C: 1},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			cpu.SP = test.SP
			writeTestProgram(cpu, ADD_SP_E8_OPCODE, uint8(test.e8))

			cpu.ExecuteInstruction()
			expected_SP := int(test.SP) + int(test.e8)
			if cpu.SP != uint16(expected_SP) {
				t.Fatalf("got %04X, expected %04X", cpu.SP, expected_SP)
			}
			if cpu.readZFlag() != 0 {
				t.Fatalf("Z flag: got %x, expected %x", cpu.readZFlag(), 0)
			}
			if cpu.readNFlag() != 0 {
				t.Fatalf("N flag: got %x, expected %x", cpu.readNFlag(), 0)
			}
			if cpu.readHFlag() != test.exp_H {
				t.Fatalf("H flag: got %x, expected %x", cpu.readHFlag(), test.exp_H)
			}
			if cpu.readCFlag() != test.exp_C {
				t.Fatalf("C flag: got %x, expected %x", cpu.readCFlag(), test.exp_C)
			}
		})
	}
}

func Test_LD_HL_SP_E8(t *testing.T) {
	cpu := mockCPU()

	var SP uint16 = 0xFFFE
	var e8 int8 = -1
	var exp_H, exp_C uint8 = 1, 1

	cpu.SP = SP
	writeTestProgram(cpu, LD_HL_SP_E8_OPCODE, uint8(e8))

	cpu.ExecuteInstruction()
	expected_HL := int(SP) + int(e8)
	if cpu.readHL() != uint16(expected_HL) {
		t.Fatalf("got %04X, expected %04X", cpu.readHL(), expected_HL)
	}
	if cpu.readZFlag() != 0 {
		t.Fatalf("Z flag: got %x, expected %x", cpu.readZFlag(), 0)
	}
	if cpu.readNFlag() != 0 {
		t.Fatalf("N flag: got %x, expected %x", cpu.readNFlag(), 0)
	}
	if cpu.readHFlag() != exp_H {
		t.Fatalf("H flag: got %x, expected %x", cpu.readHFlag(), exp_H)
	}
	if cpu.readCFlag() != exp_C {
		t.Fatalf("C flag: got %x, expected %x", cpu.readCFlag(), exp_C)
	}
}

func Test_LD_SP_HL(t *testing.T) {
	cpu := mockCPU()

	var HL uint16 = 0x4321

	cpu.writeHL(HL)
	writeTestProgram(cpu, LD_SP_HL_OPCODE)
	cpu.ExecuteInstruction()

	if cpu.SP != HL {
		t.Fatalf("got %04X, expected %04X", cpu.SP, HL)
	}
}

func Test_DI(t *testing.T) {
	cpu := mockCPU()
	cpu.IME = true

	writeTestProgram(cpu, DI_OPCODE)

	cpu.ExecuteInstruction()
	if cpu.IME == true {
		t.Fatal("IME was not set to false")
	}
}

func Test_EI(t *testing.T) {
	cpu := mockCPU()
	cpu.IME = false

	writeTestProgram(cpu, EI_OPCODE, NOP_OPCODE)

	cpu.ExecuteInstruction()
	if cpu.IME == true {
		t.Error("EI instruction was not delayed")
	}

	cpu.ExecuteInstruction()
	if cpu.IME == false {
		t.Fatal("IME was not set")
	}
}

var R8_offset = map[string]uint8{"B": 0, "C": 1, "D": 2, "E": 3, "H": 4, "L": 5, "HLMEM": 6, "A": 7}
var R8_name = [8]string{"B", "C", "D", "E", "H", "L", "HLMEM", "A"}

func Test_RLC_R8(t *testing.T) {
	cpu := mockCPU()
	// Test flag reset
	cpu.F = 0xFF

	r8_name := "B"
	r8 := R8_offset[r8_name]

	tests := map[string]struct {
		R8     uint8
		exp_R8 uint8
		exp_c  uint8
		exp_z  uint8
	}{
		"carry":    {R8: 0b10101010, exp_R8: 0b01010101, exp_c: 1, exp_z: 0},
		"no-carry": {R8: 0b01110001, exp_R8: 0b11100010, exp_c: 0, exp_z: 0},
		"zero":     {R8: 0, exp_R8: 0, exp_c: 0, exp_z: 1},
	}

	for name, test := range tests {
		t.Run(r8_name+"/"+name, func(t *testing.T) {
			cpu.writeR8(r8, test.R8)
			writeTestProgram(cpu, PREFIX_OPCODE, RLC_R8_OPCODE+r8)
			cpu.ExecuteInstruction()

			if cpu.readR8(r8) != test.exp_R8 {
				t.Fatalf("got %08b, expected %08b", cpu.readR8(r8), test.exp_R8)
			}
			// Test flags
			if cpu.readZFlag() != test.exp_z {
				t.Errorf("Z flag: got %x, expected %x", cpu.readZFlag(), test.exp_z)
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

func Test_RRC_R8(t *testing.T) {
	cpu := mockCPU()
	// Test flag reset
	cpu.F = 0xFF

	r8_name := "C"
	r8 := R8_offset[r8_name]

	tests := map[string]struct {
		R8     uint8
		exp_R8 uint8
		exp_c  uint8
		exp_z  uint8
	}{
		"carry":    {R8: 0b10101010, exp_R8: 0b01010101, exp_c: 0, exp_z: 0},
		"no-carry": {R8: 0b01110001, exp_R8: 0b10111000, exp_c: 1, exp_z: 0},
		"zero":     {R8: 0, exp_R8: 0, exp_c: 0, exp_z: 1},
	}

	for name, test := range tests {
		t.Run(r8_name+"/"+name, func(t *testing.T) {
			cpu.writeR8(r8, test.R8)
			writeTestProgram(cpu, PREFIX_OPCODE, RRC_R8_OPCODE+r8)
			cpu.ExecuteInstruction()

			if cpu.readR8(r8) != test.exp_R8 {
				t.Fatalf("got %08b, expected %08b", cpu.readR8(r8), test.exp_R8)
			}
			// Test flags
			if cpu.readZFlag() != test.exp_z {
				t.Errorf("Z flag: got %x, expected %x", cpu.readZFlag(), test.exp_z)
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

func Test_RL_R8(t *testing.T) {
	cpu := mockCPU()
	// Test flag reset
	cpu.F = 0xFF

	r8_name := "D"
	r8 := R8_offset[r8_name]

	tests := map[string]struct {
		R8     uint8
		carry  uint8
		exp_R8 uint8
		exp_c  uint8
		exp_z  uint8
	}{
		"carry_set":   {R8: 0b01011010, carry: 1, exp_R8: 0b10110101, exp_c: 0, exp_z: 0},
		"carry_unset": {R8: 0b11011010, carry: 0, exp_R8: 0b10110100, exp_c: 1, exp_z: 0},
		"carry":       {R8: 0b11011010, carry: 1, exp_R8: 0b10110101, exp_c: 1, exp_z: 0},
		"no-carry":    {R8: 0b01011010, carry: 0, exp_R8: 0b10110100, exp_c: 0, exp_z: 0},
		"zero":        {R8: 0b10000000, carry: 0, exp_R8: 0, exp_c: 1, exp_z: 1},
	}

	for name, test := range tests {
		t.Run(r8_name+"/"+name, func(t *testing.T) {
			cpu.writeR8(r8, test.R8)
			cpu.setCFlag(test.carry)
			writeTestProgram(cpu, PREFIX_OPCODE, RL_R8_OPCODE+r8)
			cpu.ExecuteInstruction()

			if cpu.readR8(r8) != test.exp_R8 {
				t.Fatalf("got %08b, expected %08b", cpu.readR8(r8), test.exp_R8)
			}
			// Test flags
			if cpu.readZFlag() != test.exp_z {
				t.Errorf("Z flag: got %x, expected %x", cpu.readZFlag(), test.exp_z)
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

func Test_RR_R8(t *testing.T) {
	cpu := mockCPU()
	// Test flag reset
	cpu.F = 0xFF

	r8_name := "E"
	r8 := R8_offset[r8_name]

	tests := map[string]struct {
		R8     uint8
		carry  uint8
		exp_R8 uint8
		exp_c  uint8
		exp_z  uint8
	}{
		"carry_set":   {R8: 0b01011010, carry: 1, exp_R8: 0b10101101, exp_c: 0, exp_z: 0},
		"carry_unset": {R8: 0b11011011, carry: 0, exp_R8: 0b01101101, exp_c: 1, exp_z: 0},
		"carry":       {R8: 0b11011011, carry: 1, exp_R8: 0b11101101, exp_c: 1, exp_z: 0},
		"no-carry":    {R8: 0b01011010, carry: 0, exp_R8: 0b00101101, exp_c: 0, exp_z: 0},
		"zero":        {R8: 0b00000001, carry: 0, exp_R8: 0, exp_c: 1, exp_z: 1},
	}

	for name, test := range tests {
		t.Run(r8_name+"/"+name, func(t *testing.T) {
			cpu.writeR8(r8, test.R8)
			cpu.setCFlag(test.carry)
			writeTestProgram(cpu, PREFIX_OPCODE, RR_R8_OPCODE+r8)
			cpu.ExecuteInstruction()

			if cpu.readR8(r8) != test.exp_R8 {
				t.Fatalf("got %08b, expected %08b", cpu.readR8(r8), test.exp_R8)
			}
			// Test flags
			if cpu.readZFlag() != test.exp_z {
				t.Errorf("Z flag: got %x, expected %x", cpu.readZFlag(), test.exp_z)
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

func Test_SLA_R8(t *testing.T) {
	cpu := mockCPU()
	r8_name := "H"
	r8 := R8_offset[r8_name]

	tests := map[string]struct {
		R8     uint8
		exp_R8 uint8
		exp_c  uint8
		exp_z  uint8
	}{
		"carry":    {R8: 0b11011011, exp_R8: 0b10110110, exp_c: 1, exp_z: 0},
		"no-carry": {R8: 0b01011010, exp_R8: 0b10110100, exp_c: 0, exp_z: 0},
		"zero":     {R8: 0b10000000, exp_R8: 0, exp_c: 1, exp_z: 1},
	}

	for name, test := range tests {
		t.Run(r8_name+"/"+name, func(t *testing.T) {
			cpu.writeR8(r8, test.R8)
			writeTestProgram(cpu, PREFIX_OPCODE, SLA_R8_OPCODE+r8)
			cpu.ExecuteInstruction()

			if cpu.readR8(r8) != test.exp_R8 {
				t.Fatalf("got %08b, expected %08b", cpu.readR8(r8), test.exp_R8)
			}
			// Test flags
			if cpu.readZFlag() != test.exp_z {
				t.Errorf("Z flag: got %x, expected %x", cpu.readZFlag(), test.exp_z)
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

func Test_SRA_R8(t *testing.T) {
	cpu := mockCPU()
	r8_name := "L"
	r8 := R8_offset[r8_name]

	tests := map[string]struct {
		R8     uint8
		exp_R8 uint8
		exp_c  uint8
		exp_z  uint8
	}{
		"carry":    {R8: 0b11011011, exp_R8: 0b11101101, exp_c: 1, exp_z: 0},
		"no-carry": {R8: 0b01011010, exp_R8: 0b00101101, exp_c: 0, exp_z: 0},
		"zero":     {R8: 0b00000001, exp_R8: 0, exp_c: 1, exp_z: 1},
	}

	for name, test := range tests {
		t.Run(r8_name+"/"+name, func(t *testing.T) {
			cpu.writeR8(r8, test.R8)
			writeTestProgram(cpu, PREFIX_OPCODE, SRA_R8_OPCODE+r8)
			cpu.ExecuteInstruction()

			if cpu.readR8(r8) != test.exp_R8 {
				t.Fatalf("got %08b, expected %08b", cpu.readR8(r8), test.exp_R8)
			}
			// Test flags
			if cpu.readZFlag() != test.exp_z {
				t.Errorf("Z flag: got %x, expected %x", cpu.readZFlag(), test.exp_z)
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

func Test_SWAP_R8(t *testing.T) {
	cpu := mockCPU()
	r8_name := "HLMEM"
	r8 := R8_offset[r8_name]

	tests := map[string]struct {
		R8     uint8
		exp_R8 uint8
		exp_z  uint8
	}{
		"non-zero": {R8: 0xAF, exp_R8: 0xFA, exp_z: 0},
		"zero":     {R8: 0, exp_R8: 0, exp_z: 1},
	}

	for name, test := range tests {
		t.Run(r8_name+"/"+name, func(t *testing.T) {
			addr := uint16(0xF000) + uint16(test.R8)
			cpu.writeHL(addr)
			cpu.writeR8(r8, test.R8)
			writeTestProgram(cpu, PREFIX_OPCODE, SWAP_R8_OPCODE+r8)
			cpu.ExecuteInstruction()

			if cpu.readR8(r8) != test.exp_R8 {
				t.Fatalf("got %2X, expected %2X", cpu.readR8(r8), test.exp_R8)
			}
			// Test flags
			if cpu.readZFlag() != test.exp_z {
				t.Errorf("Z flag: got %x, expected %x", cpu.readZFlag(), test.exp_z)
			}
			if cpu.readNFlag() != 0 {
				t.Error("N flag: got 1, expected 0")
			}
			if cpu.readHFlag() != 0 {
				t.Error("H flag: got 1, expected 0")
			}
			if cpu.readCFlag() != 0 {
				t.Errorf("C flag: got 1, expected 0")
			}
		})
	}
}

func Test_SRL_R8(t *testing.T) {
	cpu := mockCPU()
	r8_name := "A"
	r8 := R8_offset[r8_name]

	tests := map[string]struct {
		R8     uint8
		exp_R8 uint8
		exp_c  uint8
		exp_z  uint8
	}{
		"carry":    {R8: 0b11011011, exp_R8: 0b01101101, exp_c: 1, exp_z: 0},
		"no-carry": {R8: 0b01011010, exp_R8: 0b00101101, exp_c: 0, exp_z: 0},
		"zero":     {R8: 0b00000001, exp_R8: 0, exp_c: 1, exp_z: 1},
	}

	for name, test := range tests {
		t.Run(r8_name+"/"+name, func(t *testing.T) {
			cpu.writeR8(r8, test.R8)
			writeTestProgram(cpu, PREFIX_OPCODE, SRL_R8_OPCODE+r8)
			cpu.ExecuteInstruction()

			if cpu.readR8(r8) != test.exp_R8 {
				t.Fatalf("got %08b, expected %08b", cpu.readR8(r8), test.exp_R8)
			}
			// Test flags
			if cpu.readZFlag() != test.exp_z {
				t.Errorf("Z flag: got %x, expected %x", cpu.readZFlag(), test.exp_z)
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

func Test_BIT_B3_R8(t *testing.T) {
	cpu := mockCPU()

	tests := [8]struct {
		opcode uint8
		R8     uint8
		exp_z  uint8
	}{
		{opcode: BIT_0_R8_OPCODE + 0, R8: 0b11011010, exp_z: 1},
		{opcode: BIT_1_R8_OPCODE + 1, R8: 0b01011010, exp_z: 0},
		{opcode: BIT_2_R8_OPCODE + 2, R8: 0b00000101, exp_z: 0},
		{opcode: BIT_3_R8_OPCODE + 3, R8: 0b11010011, exp_z: 1},
		{opcode: BIT_4_R8_OPCODE + 4, R8: 0b01001010, exp_z: 1},
		{opcode: BIT_5_R8_OPCODE + 5, R8: 0b00100001, exp_z: 0},
		{opcode: BIT_6_R8_OPCODE + 6, R8: 0b10011011, exp_z: 1},
		{opcode: BIT_7_R8_OPCODE + 7, R8: 0b11011010, exp_z: 0},
	}

	for i, test := range tests {
		t.Run("B"+fmt.Sprint(i)+"/"+R8_name[i], func(t *testing.T) {
			cpu.writeR8(test.opcode, test.R8)
			writeTestProgram(cpu, PREFIX_OPCODE, test.opcode)
			cpu.ExecuteInstruction()

			if cpu.readZFlag() != test.exp_z {
				t.Errorf("Z flag: got %x, expected %x", cpu.readZFlag(), test.exp_z)
			}
			if cpu.readNFlag() != 0 {
				t.Error("N flag: got 1, expected 0")
			}
			if cpu.readHFlag() != 1 {
				t.Error("H flag: got 0, expected 1")
			}
		})
	}
}

func Test_RES_B3_R8(t *testing.T) {
	cpu := mockCPU()

	tests := [8]struct {
		opcode uint8
		exp_R8 uint8
	}{
		{opcode: RES_0_R8_OPCODE + 0, exp_R8: 0b11111110},
		{opcode: RES_1_R8_OPCODE + 1, exp_R8: 0b11111101},
		{opcode: RES_2_R8_OPCODE + 2, exp_R8: 0b11111011},
		{opcode: RES_3_R8_OPCODE + 3, exp_R8: 0b11110111},
		{opcode: RES_4_R8_OPCODE + 4, exp_R8: 0b11101111},
		{opcode: RES_5_R8_OPCODE + 5, exp_R8: 0b11011111},
		{opcode: RES_6_R8_OPCODE + 6, exp_R8: 0b10111111},
		{opcode: RES_7_R8_OPCODE + 7, exp_R8: 0b01111111},
	}

	for i, test := range tests {
		t.Run("B"+fmt.Sprint(i)+"/"+R8_name[i], func(t *testing.T) {
			cpu.writeR8(test.opcode, 0xFF)
			writeTestProgram(cpu, PREFIX_OPCODE, test.opcode)
			cpu.ExecuteInstruction()

			if cpu.readR8(test.opcode) != test.exp_R8 {
				t.Fatalf("got %08b, expected %08b", cpu.readR8(test.opcode), test.exp_R8)
			}
		})
	}
}

func Test_SET_B3_R8(t *testing.T) {
	cpu := mockCPU()

	tests := [8]struct {
		opcode uint8
		exp_R8 uint8
	}{
		{opcode: SET_0_R8_OPCODE + 0, exp_R8: 0b00000001},
		{opcode: SET_1_R8_OPCODE + 1, exp_R8: 0b00000010},
		{opcode: SET_2_R8_OPCODE + 2, exp_R8: 0b00000100},
		{opcode: SET_3_R8_OPCODE + 3, exp_R8: 0b00001000},
		{opcode: SET_4_R8_OPCODE + 4, exp_R8: 0b00010000},
		{opcode: SET_5_R8_OPCODE + 5, exp_R8: 0b00100000},
		{opcode: SET_6_R8_OPCODE + 6, exp_R8: 0b01000000},
		{opcode: SET_7_R8_OPCODE + 7, exp_R8: 0b10000000},
	}

	for i, test := range tests {
		t.Run("B"+fmt.Sprint(i)+"/"+R8_name[i], func(t *testing.T) {
			cpu.writeR8(test.opcode, 0)
			writeTestProgram(cpu, PREFIX_OPCODE, test.opcode)
			cpu.ExecuteInstruction()

			if cpu.readR8(test.opcode) != test.exp_R8 {
				t.Fatalf("got %08b, expected %08b", cpu.readR8(test.opcode), test.exp_R8)
			}
		})
	}
}
