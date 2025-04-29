package cpu

import (
	"fmt"
	"github.com/danielecanzoneri/gb-emulator/internal/util"
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
		expected := util.CombineBytes(BYTE2, BYTE1)
		if !(cpu.SP == expected) {
			t.Fatalf("got %X%X, expected %X%X", cpu.D, cpu.E, BYTE2, BYTE1)
		}
	})
}

func Test_LD_R16MEM_A(t *testing.T) {
	cpu := mockCPU()
	cpu.A = 0xFD

	var addrBC uint16 = 0xA010
	var addrDE uint16 = 0xA020
	var addrHLI uint16 = 0xA030
	var addrHLD uint16 = 0xA040

	t.Run("BC", func(t *testing.T) {
		cpu.writeBC(addrBC)
		writeTestProgram(cpu, LD_BCmem_A_OPCODE)
		cpu.ExecuteInstruction()
		if cpu.MMU.Read(addrBC) != cpu.A {
			t.Fatalf("got %X, expected %X", cpu.MMU.Read(addrBC), cpu.A)
		}
	})

	t.Run("DE", func(t *testing.T) {
		cpu.writeDE(addrDE)
		writeTestProgram(cpu, LD_DEmem_A_OPCODE)
		cpu.ExecuteInstruction()
		if cpu.MMU.Read(addrDE) != cpu.A {
			t.Fatalf("got %X, expected %X", cpu.MMU.Read(addrDE), cpu.A)
		}
	})

	t.Run("HLI", func(t *testing.T) {
		cpu.writeHL(addrHLI)
		writeTestProgram(cpu, LD_HLImem_A_OPCODE)
		cpu.ExecuteInstruction()
		if cpu.MMU.Read(addrHLI) != cpu.A {
			t.Fatalf("got %X, expected %X", cpu.MMU.Read(addrHLI), cpu.A)
		}
		if cpu.readHL() != addrHLI+1 {
			t.Error("increment failed")
		}
	})

	t.Run("HLD", func(t *testing.T) {
		cpu.writeHL(addrHLD)
		writeTestProgram(cpu, LD_HLDmem_A_OPCODE)
		cpu.ExecuteInstruction()
		if cpu.MMU.Read(addrHLD) != cpu.A {
			t.Fatalf("got %X, expected %X", cpu.MMU.Read(addrHLD), cpu.A)
		}
		if cpu.readHL() != addrHLD-1 {
			t.Error("decrement failed")
		}
	})
}

func Test_LD_A_R16MEM(t *testing.T) {
	cpu := mockCPU()

	var addrBC uint16 = 0xA010
	var byteBC uint8 = 0x01
	var addrDE uint16 = 0xA020
	var byteDE uint8 = 0x02
	var addrHLI uint16 = 0xA030
	var byteHLI uint8 = 0x03
	var addrHLD uint16 = 0xA040
	var byteHLD uint8 = 0x04

	t.Run("BC", func(t *testing.T) {
		cpu.writeBC(addrBC)
		cpu.MMU.Write(addrBC, byteBC)
		writeTestProgram(cpu, LD_A_BCMEM_OPCODE)
		cpu.ExecuteInstruction()
		if cpu.A != byteBC {
			t.Fatalf("got %X, expected %X", cpu.A, byteBC)
		}
	})

	t.Run("DE", func(t *testing.T) {
		cpu.writeDE(addrDE)
		cpu.MMU.Write(addrDE, byteDE)
		writeTestProgram(cpu, LD_A_DEMEM_OPCODE)
		cpu.ExecuteInstruction()
		if cpu.A != byteDE {
			t.Fatalf("got %X, expected %X", cpu.A, byteDE)
		}
	})

	t.Run("HLI", func(t *testing.T) {
		cpu.writeHL(addrHLI)
		cpu.MMU.Write(addrHLI, byteHLI)
		writeTestProgram(cpu, LD_A_HLIMEM_OPCODE)
		cpu.ExecuteInstruction()
		if cpu.A != byteHLI {
			t.Fatalf("got %X, expected %X", cpu.A, byteHLI)
		}
		if cpu.readHL() != addrHLI+1 {
			t.Error("increment failed")
		}
	})

	t.Run("HLD", func(t *testing.T) {
		cpu.writeHL(addrHLD)
		cpu.MMU.Write(addrHLD, byteHLD)
		writeTestProgram(cpu, LD_A_HLDMEM_OPCODE)
		cpu.ExecuteInstruction()
		if cpu.A != byteHLD {
			t.Fatalf("got %X, expected %X", cpu.A, byteHLD)
		}
		if cpu.readHL() != addrHLD-1 {
			t.Error("decrement failed")
		}
	})
}

func Test_LD_N16_SP(t *testing.T) {
	cpu := mockCPU()
	cpu.SP = 0xFD53
	var addr uint16 = 0xA034
	highAddr, lowAddr := util.SplitWord(addr)

	writeTestProgram(cpu, LD_N16_SP_OPCODE, lowAddr, highAddr)
	cpu.ExecuteInstruction()

	read := cpu.MMU.ReadWord(addr)
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

	type testAdd struct {
		hl   uint16
		r16  uint16
		sum  uint16
		expH uint8
		expC uint8
	}

	tests := map[string]testAdd{
		"standard":   {hl: 0x00FF, r16: 0x00FF, sum: 0x01FE, expH: 0, expC: 0},
		"half-carry": {hl: 0x0800, r16: 0x0800, sum: 0x1000, expH: 1, expC: 0},
		"carry":      {hl: 0xFFFF, r16: 0xFFFF, sum: 0xFFFE, expH: 1, expC: 1},
	}

	testR16 := func(test testAdd, t *testing.T) {
		cpu.writeHL(test.hl)
		cpu.ExecuteInstruction()

		if cpu.readHL() != test.sum {
			t.Fatalf("wrong sum: got %04X, expected %04X", cpu.readHL(), test.sum)
		}
		if cpu.readNFlag() != 0 {
			t.Error("wrong N flag: should be 0")
		}
		if cpu.readHFlag() != test.expH {
			t.Errorf("wrong H flag: got %d, expected %d", cpu.readHFlag(), test.expH)
		}
		if cpu.readCFlag() != test.expC {
			t.Errorf("wrong C flag: got %d, expected %d", cpu.readCFlag(), test.expC)
		}
	}

	for name, test := range tests {
		t.Run("BC_"+name, func(t *testing.T) {
			cpu.writeBC(test.r16)
			writeTestProgram(cpu, ADD_HL_BC_OPCODE)
			testR16(test, t)
		})
		t.Run("DE_"+name, func(t *testing.T) {
			cpu.writeDE(test.r16)
			writeTestProgram(cpu, ADD_HL_DE_OPCODE)
			testR16(test, t)
		})
		t.Run("HL_"+name, func(t *testing.T) {
			cpu.writeHL(test.r16)
			writeTestProgram(cpu, ADD_HL_HL_OPCODE)
			testR16(test, t)
		})
		t.Run("SP_"+name, func(t *testing.T) {
			cpu.SP = test.r16
			writeTestProgram(cpu, ADD_HL_SP_OPCODE)
			testR16(test, t)
		})
		cpu.F = 0
	}
}

func Test_INC_R8(t *testing.T) {
	cpu := mockCPU()

	var B uint8 = 0x01
	var ZFlagB, HFlagB uint8 = 0, 0
	var C uint8 = 0x02
	var ZFlagC, HFlagC uint8 = 0, 0
	var D uint8 = 0x03
	var ZFlagD, HFlagD uint8 = 0, 0
	var E uint8 = 0x04
	var ZFlagE, HFlagE uint8 = 0, 0
	var H uint8 = 0x05
	var ZFlagH, HFlagH uint8 = 0, 0
	var L uint8 = 0x06
	var ZFlagL, HFlagL uint8 = 0, 0
	var A uint8 = 0x0F
	var ZFlagA, HFlagA uint8 = 0, 1

	cpu.A = A
	cpu.B = B
	cpu.C = C
	cpu.D = D
	cpu.E = E
	cpu.H = H
	cpu.L = L

	testCarries := func(t *testing.T, expZ, expH uint8) {
		if cpu.readNFlag() != 0 {
			t.Errorf("N flag: got %x, expected %x", cpu.readNFlag(), 0)
		}
		if cpu.readZFlag() != expZ {
			t.Errorf("Z flag: got %x, expected %x", cpu.readZFlag(), expZ)
		}
		if cpu.readHFlag() != expH {
			t.Errorf("H flag: got %x, expected %x", cpu.readHFlag(), expH)
		}
	}

	t.Run("B", func(t *testing.T) {
		writeTestProgram(cpu, INC_B_OPCODE)
		cpu.ExecuteInstruction()
		if cpu.B != B+1 {
			t.Fatalf("got %2X, expected %2X", cpu.B, B+1)
		}
		testCarries(t, ZFlagB, HFlagB)
	})

	t.Run("C", func(t *testing.T) {
		writeTestProgram(cpu, INC_C_OPCODE)
		cpu.ExecuteInstruction()
		if cpu.C != C+1 {
			t.Fatalf("got %2X, expected %2X", cpu.C, C+1)
		}
		testCarries(t, ZFlagC, HFlagC)
	})

	t.Run("D", func(t *testing.T) {
		writeTestProgram(cpu, INC_D_OPCODE)
		cpu.ExecuteInstruction()
		if cpu.D != D+1 {
			t.Fatalf("got %2X, expected %2X", cpu.D, D+1)
		}
		testCarries(t, ZFlagD, HFlagD)
	})

	t.Run("E", func(t *testing.T) {
		writeTestProgram(cpu, INC_E_OPCODE)
		cpu.ExecuteInstruction()
		if cpu.E != E+1 {
			t.Fatalf("got %2X, expected %2X", cpu.E, E+1)
		}
		testCarries(t, ZFlagE, HFlagE)
	})

	t.Run("H", func(t *testing.T) {
		writeTestProgram(cpu, INC_H_OPCODE)
		cpu.ExecuteInstruction()
		if cpu.H != H+1 {
			t.Fatalf("got %2X, expected %2X", cpu.H, H+1)
		}
		testCarries(t, ZFlagH, HFlagH)
	})

	t.Run("L", func(t *testing.T) {
		writeTestProgram(cpu, INC_L_OPCODE)
		cpu.ExecuteInstruction()
		if cpu.L != L+1 {
			t.Fatalf("got %2X, expected %2X", cpu.L, L+1)
		}
		testCarries(t, ZFlagL, HFlagL)
	})

	t.Run("A", func(t *testing.T) {
		writeTestProgram(cpu, INC_A_OPCODE)
		cpu.ExecuteInstruction()
		if cpu.A != A+1 {
			t.Fatalf("got %2X, expected %2X", cpu.A, A+1)
		}
		testCarries(t, ZFlagA, HFlagA)
	})

	var addrHL uint16 = 0xA050
	var value uint8 = 0xFF
	var ZFlagHL, HFlagHL uint8 = 1, 1

	cpu.writeHL(addrHL)
	cpu.MMU.Write(addrHL, value)

	t.Run("HL_MEM", func(t *testing.T) {
		writeTestProgram(cpu, INC_HLMEM_OPCODE)
		cpu.ExecuteInstruction()
		if cpu.MMU.Read(addrHL) != value+1 {
			t.Fatalf("got %2X, expected %2X", cpu.MMU.Read(addrHL), value+1)
		}
		testCarries(t, ZFlagHL, HFlagHL)
	})
}

func Test_DEC_R8(t *testing.T) {
	cpu := mockCPU()

	var B uint8 = 0x01
	var ZFlagB, HFlagB uint8 = 1, 0
	var C uint8 = 0x02
	var ZFlagC, HFlagC uint8 = 0, 0
	var D uint8 = 0x03
	var ZFlagD, HFlagD uint8 = 0, 0
	var E uint8 = 0x04
	var ZFlagE, HFlagE uint8 = 0, 0
	var H uint8 = 0x05
	var ZFlagH, HFlagH uint8 = 0, 0
	var L uint8 = 0x06
	var ZFlagL, HFlagL uint8 = 0, 0
	var A uint8 = 0x10
	var ZFlagA, HFlagA uint8 = 0, 1

	cpu.A = A
	cpu.B = B
	cpu.C = C
	cpu.D = D
	cpu.E = E
	cpu.H = H
	cpu.L = L

	testCarries := func(t *testing.T, expZ, expH uint8) {
		if cpu.readNFlag() != 1 {
			t.Errorf("N flag: got %x, expected %x", cpu.readNFlag(), 1)
		}
		if cpu.readZFlag() != expZ {
			t.Errorf("Z flag: got %x, expected %x", cpu.readZFlag(), expZ)
		}
		if cpu.readHFlag() != expH {
			t.Errorf("H flag: got %x, expected %x", cpu.readHFlag(), expH)
		}
	}

	t.Run("B", func(t *testing.T) {
		writeTestProgram(cpu, DEC_B_OPCODE)
		cpu.ExecuteInstruction()
		if cpu.B != B-1 {
			t.Fatalf("got %2X, expected %2X", cpu.B, B-1)
		}
		testCarries(t, ZFlagB, HFlagB)
	})

	t.Run("C", func(t *testing.T) {
		writeTestProgram(cpu, DEC_C_OPCODE)
		cpu.ExecuteInstruction()
		if cpu.C != C-1 {
			t.Fatalf("got %2X, expected %2X", cpu.C, C-1)
		}
		testCarries(t, ZFlagC, HFlagC)
	})

	t.Run("D", func(t *testing.T) {
		writeTestProgram(cpu, DEC_D_OPCODE)
		cpu.ExecuteInstruction()
		if cpu.D != D-1 {
			t.Fatalf("got %2X, expected %2X", cpu.D, D-1)
		}
		testCarries(t, ZFlagD, HFlagD)
	})

	t.Run("E", func(t *testing.T) {
		writeTestProgram(cpu, DEC_E_OPCODE)
		cpu.ExecuteInstruction()
		if cpu.E != E-1 {
			t.Fatalf("got %2X, expected %2X", cpu.E, E-1)
		}
		testCarries(t, ZFlagE, HFlagE)
	})

	t.Run("H", func(t *testing.T) {
		writeTestProgram(cpu, DEC_H_OPCODE)
		cpu.ExecuteInstruction()
		if cpu.H != H-1 {
			t.Fatalf("got %2X, expected %2X", cpu.H, H-1)
		}
		testCarries(t, ZFlagH, HFlagH)
	})

	t.Run("L", func(t *testing.T) {
		writeTestProgram(cpu, DEC_L_OPCODE)
		cpu.ExecuteInstruction()
		if cpu.L != L-1 {
			t.Fatalf("got %2X, expected %2X", cpu.L, L-1)
		}
		testCarries(t, ZFlagL, HFlagL)
	})

	t.Run("A", func(t *testing.T) {
		writeTestProgram(cpu, DEC_A_OPCODE)
		cpu.ExecuteInstruction()
		if cpu.A != A-1 {
			t.Fatalf("got %2X, expected %2X", cpu.A, A-1)
		}
		testCarries(t, ZFlagA, HFlagA)
	})

	var addrHL uint16 = 0xA050
	var value uint8 = 0x00
	var ZFlagHL, HFlagHL uint8 = 0, 1

	cpu.writeHL(addrHL)
	cpu.MMU.Write(addrHL, value)

	t.Run("HL_MEM", func(t *testing.T) {
		writeTestProgram(cpu, DEC_HLMEM_OPCODE)
		cpu.ExecuteInstruction()
		if cpu.MMU.Read(addrHL) != value-1 {
			t.Fatalf("got %2X, expected %2X", cpu.MMU.Read(addrHL), value-1)
		}
		testCarries(t, ZFlagHL, HFlagHL)
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

	t.Run("HL_MEM", func(t *testing.T) {
		var addr uint16 = 0xA0F0
		cpu.writeHL(addr)
		writeTestProgram(cpu, LD_HLMEM_N8_OPCODE, value)
		cpu.ExecuteInstruction()
		if cpu.MMU.Read(addr) != value {
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
		A    uint8
		expA uint8
		expC uint8
	}{
		"carry":    {A: 0b10101010, expA: 0b01010101, expC: 1},
		"no-carry": {A: 0b01110001, expA: 0b11100010, expC: 0},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			cpu.A = test.A
			writeTestProgram(cpu, RLCA_OPCODE)
			cpu.ExecuteInstruction()

			if cpu.A != test.expA {
				t.Fatalf("got %08b, expected %08b", cpu.A, test.expA)
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
			if cpu.readCFlag() != test.expC {
				t.Errorf("C flag: got %x, expected %x", cpu.readCFlag(), test.expC)
			}
		})
	}
}

func Test_RRCA(t *testing.T) {
	cpu := mockCPU()
	// Test flag reset
	cpu.F = 0xFF

	tests := map[string]struct {
		A    uint8
		expA uint8
		expC uint8
	}{
		"no-carry": {A: 0b10101010, expA: 0b01010101, expC: 0},
		"carry":    {A: 0b01110001, expA: 0b10111000, expC: 1},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			cpu.A = test.A
			writeTestProgram(cpu, RRCA_OPCODE)
			cpu.ExecuteInstruction()

			if cpu.A != test.expA {
				t.Fatalf("got %08b, expected %08b", cpu.A, test.expA)
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
			if cpu.readCFlag() != test.expC {
				t.Errorf("C flag: got %x, expected %x", cpu.readCFlag(), test.expC)
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
		expA  uint8
		expC  uint8
	}{
		"carry_set":   {A: 0b01011010, carry: 1, expA: 0b10110101, expC: 0},
		"carry_unset": {A: 0b11011010, carry: 0, expA: 0b10110100, expC: 1},
		"carry":       {A: 0b11011010, carry: 1, expA: 0b10110101, expC: 1},
		"no-carry":    {A: 0b01011010, carry: 0, expA: 0b10110100, expC: 0},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			cpu.A = test.A
			cpu.setCFlag(test.carry)
			writeTestProgram(cpu, RLA_OPCODE)
			cpu.ExecuteInstruction()

			if cpu.A != test.expA {
				t.Fatalf("got %08b, expected %08b", cpu.A, test.expA)
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
			if cpu.readCFlag() != test.expC {
				t.Errorf("C flag: got %x, expected %x", cpu.readCFlag(), test.expC)
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
		expA  uint8
		expC  uint8
	}{
		"carry_set":   {A: 0b01011010, carry: 1, expA: 0b10101101, expC: 0},
		"carry_unset": {A: 0b11011011, carry: 0, expA: 0b01101101, expC: 1},
		"carry":       {A: 0b11011011, carry: 1, expA: 0b11101101, expC: 1},
		"no-carry":    {A: 0b01011010, carry: 0, expA: 0b00101101, expC: 0},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			cpu.A = test.A
			cpu.setCFlag(test.carry)
			writeTestProgram(cpu, RRA_OPCODE)
			cpu.ExecuteInstruction()

			if cpu.A != test.expA {
				t.Fatalf("got %08b, expected %08b", cpu.A, test.expA)
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
			if cpu.readCFlag() != test.expC {
				t.Errorf("C flag: got %x, expected %x", cpu.readCFlag(), test.expC)
			}
		})
	}
}

func Test_DAA(t *testing.T) {
	cpu := mockCPU()

	tests := map[string]struct {
		n1    uint8
		n2    uint8
		isSub uint8
	}{
		"sum-HFlag":            {n1: 19, n2: 19, isSub: 0},
		"sum-low_adj":          {n1: 15, n2: 19, isSub: 0},
		"sum-HFlag,high_adj":   {n1: 69, n2: 59, isSub: 0},
		"sum-low_adj,high_adj": {n1: 65, n2: 59, isSub: 0},
		"sum-carry":            {n1: 99, n2: 99, isSub: 0},
		"sum-ZFlag":            {n1: 50, n2: 50, isSub: 0},
		"sub-HFlag":            {n1: 10, n2: 1, isSub: 1},
		"sub-CFlag":            {n1: 100, n2: 10, isSub: 1},
		"sub-HFlag,CFlag":      {n1: 100, n2: 1, isSub: 1},
		"sub-ZFlag":            {n1: 100, n2: 100, isSub: 1},
		"sub-negative":         {n1: 50, n2: 60, isSub: 1},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			bcdN1 := util.ByteToBCD(test.n1)
			bcdN2 := util.ByteToBCD(test.n2)
			var result, carry, halfCarry uint8
			if test.isSub == 0 {
				result, carry, halfCarry = util.SumBytesWithCarry(bcdN1, bcdN2)
			} else {
				result, carry, halfCarry = util.SubBytesWithCarry(bcdN1, bcdN2)
			}
			cpu.A = result
			cpu.setHFlag(halfCarry)
			cpu.setCFlag(carry)
			cpu.setNFlag(test.isSub)

			writeTestProgram(cpu, DAA_OPCODE)
			var cFlag, expA uint8
			if test.isSub == 0 {
				cFlag = (test.n1 + test.n2) / 100
				expA = util.ByteToBCD((test.n1 + test.n2) % 100)
			} else {
				if test.n1 < test.n2 {
					cFlag = 1
					expA = util.ByteToBCD(100 - (test.n2 - test.n1))
				} else {
					cFlag = 0
					expA = util.ByteToBCD(test.n1 - test.n2)
				}
				cpu.setNFlag(1)
			}

			cpu.ExecuteInstruction()

			if cpu.readCFlag() != cFlag {
				t.Errorf("C flag: got %x, expected %x", cpu.readCFlag(), cFlag)
			}
			if cpu.readHFlag() != 0 {
				t.Errorf("H flag: got 1, expected 0")
			}
			if cpu.readZFlag() != util.IsByteZeroUint8(expA) {
				t.Errorf("Z flag: got %x, expected %x", cpu.readZFlag(), util.IsByteZeroUint8(expA))
			}
			if cpu.A != expA {
				t.Fatalf("got %02X, expected %02X", cpu.A, expA)
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
		F    uint8
		expC uint8
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
			if cpu.readCFlag() != test.expC {
				t.Errorf("C flag: got %x, expected %x", cpu.readCFlag(), test.expC)
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
			expectedPC := int(cpu.PC) + OPCODES_BYTES[JR_E8_OPCODE] + int(test.e8)

			writeTestProgram(cpu, JR_E8_OPCODE, uint8(test.e8))
			cpu.ExecuteInstruction()

			if cpu.PC != uint16(expectedPC) {
				t.Fatalf("got %04X, expected %04X", cpu.PC, expectedPC)
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
	boolToInt := map[bool]uint8{false: 0, true: 1}
	conditions := map[string]func(bool){
		"Z":  func(cc bool) { cpu.setZFlag(boolToInt[cc]) },
		"NZ": func(cc bool) { cpu.setZFlag(boolToInt[!cc]) },
		"C":  func(cc bool) { cpu.setCFlag(boolToInt[cc]) },
		"NC": func(cc bool) { cpu.setCFlag(boolToInt[!cc]) },
	}
	opcodes := map[string]uint8{
		"Z":  JR_Z_E8_OPCODE,
		"NZ": JR_NZ_E8_OPCODE,
		"C":  JR_C_E8_OPCODE,
		"NC": JR_NC_E8_OPCODE,
	}

	for cond, setFlag := range conditions {
		for name, test := range tests {
			t.Run(cond+"/"+name, func(t *testing.T) {
				// Condition not met
				expectedPC := int(cpu.PC) + OPCODES_BYTES[opcodes[cond]]
				setFlag(false)

				writeTestProgram(cpu, opcodes[cond], uint8(test.e8))
				cpu.ExecuteInstruction()

				if cpu.PC != uint16(expectedPC) {
					t.Errorf("condition not met: got %04X, expected %04X", cpu.PC, expectedPC)
				}

				// Condition met
				expectedPC = int(cpu.PC) + OPCODES_BYTES[opcodes[cond]] + int(test.e8)
				setFlag(true)

				writeTestProgram(cpu, opcodes[cond], uint8(test.e8))
				cpu.ExecuteInstruction()

				if cpu.PC != uint16(expectedPC) {
					t.Errorf("condition met: got %04X, expected %04X", cpu.PC, expectedPC)
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
		"HLmem": {0x70, func() uint8 { return cpu.MMU.Read(cpu.readHL()) }},
		"A":     {0x78, func() uint8 { return cpu.A }},
	}
	secondReg := map[string]struct {
		offset uint8
		set    func(uint8) uint8
	}{
		"B": {0, func(value uint8) uint8 { cpu.B = value; return value }},
		"C": {1, func(value uint8) uint8 { cpu.C = value; return value }},
		"D": {2, func(value uint8) uint8 { cpu.D = value; return value }},
		"E": {3, func(value uint8) uint8 { cpu.E = value; return value }},
		"H": {4, func(value uint8) uint8 { cpu.H = 0xA0 + value&0x0F; return cpu.H }},
		"L": {5, func(value uint8) uint8 { cpu.L = value; return value }},
		"HLmem": {6, func(value uint8) uint8 {
			cpu.H = 0xA0
			cpu.L = value
			cpu.MMU.Write(cpu.readHL(), value)
			return value
		}},
		"A": {7, func(value uint8) uint8 { cpu.A = value; return value }},
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
				value = loaded.set(value)

				writeTestProgram(cpu, opcode)
				cpu.ExecuteInstruction()
				if receiver.read() != value {
					t.Logf("%04X", cpu.readHL())
					t.Fatalf("got %02X, expected %02X", receiver.read(), value)
				}
			})
		}
	}
}

func Test_HALT(t *testing.T) {
	cpu := mockCPU()

	t.Run("halted", func(t *testing.T) {
		PC := cpu.PC

		cpu.IME = true
		cpu.halted = false
		cpu.MMU.Write(ifAddr, 0)
		cpu.MMU.Write(ieAddr, 0)
		writeTestProgram(cpu, HALT_OPCODE, NOP_OPCODE, NOP_OPCODE)
		cpu.ExecuteInstruction()
		cpu.ExecuteInstruction()
		cpu.ExecuteInstruction()

		if cpu.PC != PC+1 {
			t.Errorf("CPU was not halted, got PC=%02X, expected %02X", cpu.PC, PC+1)
		}
	})

	t.Run("IME set", func(t *testing.T) {
		PC := cpu.PC

		cpu.IME = true
		cpu.halted = false
		cpu.MMU.Write(ifAddr, 0)
		cpu.MMU.Write(ieAddr, 0)
		writeTestProgram(cpu, HALT_OPCODE, NOP_OPCODE)
		cpu.ExecuteInstruction()
		cpu.ExecuteInstruction()

		if cpu.PC != PC+1 {
			t.Errorf("CPU was not halted, got PC=%02X, expected %02X", cpu.PC, PC+1)
		}

		cpu.MMU.Write(ifAddr, timerMask)
		cpu.MMU.Write(ieAddr, timerMask)
		cpu.ExecuteInstruction()

		if cpu.PC != timerHandler {
			t.Errorf("CPU was not woken up by interrupt, got PC=%02X, expected %02X", cpu.PC, timerHandler)
		}
	})

	t.Run("IME not set", func(t *testing.T) {
		PC := cpu.PC

		cpu.IME = false
		cpu.halted = false
		cpu.MMU.Write(ifAddr, 0)
		cpu.MMU.Write(ieAddr, 0)
		writeTestProgram(cpu, HALT_OPCODE, NOP_OPCODE)
		cpu.ExecuteInstruction()
		cpu.ExecuteInstruction()

		if cpu.PC != PC+1 {
			t.Errorf("CPU was not halted, got PC=%02X, expected %02X", cpu.PC, PC+1)
		}

		cpu.MMU.Write(ifAddr, timerMask)
		cpu.MMU.Write(ieAddr, timerMask)
		cpu.ExecuteInstruction()

		if cpu.PC != PC+2 {
			t.Errorf("CPU was not woken up by interrupt, got PC=%02X, expected %02X", cpu.PC, PC+2)
		}
	})
}

func Test_ADD_A_R8(t *testing.T) {
	cpu := mockCPU()

	type test struct {
		A    uint8
		R8   uint8
		expZ uint8
		expH uint8
		expC uint8
	}

	tests := map[string]test{
		"B":      {0x08, 0x08, 0, 1, 0},
		"C":      {0xF0, 0x11, 0, 0, 1},
		"D":      {0xFF, 0xFF, 0, 1, 1},
		"E":      {0xFF, 0x00, 0, 0, 0},
		"H":      {0x80, 0x08, 0, 0, 0},
		"L":      {0x57, 0xAD, 0, 1, 1},
		"HL_MEM": {0x34, 0x12, 0, 0, 0},
		"A":      {0x80, 0x80, 1, 0, 1},
	}

	testCarries := func(t *testing.T, tst test) {
		if cpu.readNFlag() != 0 {
			t.Errorf("N flag: got %x, expected %x", cpu.readNFlag(), 0)
		}
		if cpu.readZFlag() != tst.expZ {
			t.Errorf("Z flag: got %x, expected %x", cpu.readZFlag(), tst.expZ)
		}
		if cpu.readHFlag() != tst.expH {
			t.Errorf("H flag: got %x, expected %x", cpu.readHFlag(), tst.expH)
		}
		if cpu.readCFlag() != tst.expC {
			t.Errorf("C flag: got %x, expected %x", cpu.readCFlag(), tst.expC)
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

	t.Run("HL_MEM", func(t *testing.T) {
		tst := tests["HL_MEM"]
		cpu.A = tst.A
		var addrHL uint16 = 0xA043

		cpu.writeHL(addrHL)
		cpu.MMU.Write(addrHL, tst.R8)
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
		expZ  uint8
		expH  uint8
		expC  uint8
	}

	tests := map[string]test{
		"B":      {0x08, 0x08, 1, 0, 1, 0},
		"C":      {0xF0, 0x11, 0, 0, 0, 1},
		"D":      {0xFF, 0xFF, 1, 0, 1, 1},
		"E":      {0xFF, 0x00, 1, 1, 1, 1},
		"H":      {0x0F, 0x00, 1, 0, 1, 0},
		"L":      {0x57, 0xAD, 1, 0, 1, 1},
		"HL_MEM": {0x34, 0x12, 0, 0, 0, 0},
		"A":      {0x80, 0x80, 0, 1, 0, 1},
	}

	testCarries := func(t *testing.T, tst test) {
		if cpu.readNFlag() != 0 {
			t.Errorf("N flag: got %x, expected %x", cpu.readNFlag(), 0)
		}
		if cpu.readZFlag() != tst.expZ {
			t.Errorf("Z flag: got %x, expected %x", cpu.readZFlag(), tst.expZ)
		}
		if cpu.readHFlag() != tst.expH {
			t.Errorf("H flag: got %x, expected %x", cpu.readHFlag(), tst.expH)
		}
		if cpu.readCFlag() != tst.expC {
			t.Errorf("C flag: got %x, expected %x", cpu.readCFlag(), tst.expC)
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

	t.Run("HL_MEM", func(t *testing.T) {
		tst := tests["HL_MEM"]
		cpu.A = tst.A
		var addrHL uint16 = 0xA043

		cpu.writeHL(addrHL)
		cpu.MMU.Write(addrHL, tst.R8)
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
		expZ   uint8
		expH   uint8
		expC   uint8
		setR8  func(uint8)
		opcode uint8
	}

	tests := map[string]test{
		"B":      {0x08, 0x08, 1, 0, 0, func(v uint8) { cpu.B = v }, SUB_A_B_OPCODE},
		"C":      {0x0F, 0x10, 0, 0, 1, func(v uint8) { cpu.C = v }, SUB_A_C_OPCODE},
		"D":      {0x00, 0x01, 0, 1, 1, func(v uint8) { cpu.D = v }, SUB_A_D_OPCODE},
		"E":      {0xF0, 0x01, 0, 1, 0, func(v uint8) { cpu.E = v }, SUB_A_E_OPCODE},
		"H":      {0x80, 0x08, 0, 1, 0, func(v uint8) { cpu.H = v }, SUB_A_H_OPCODE},
		"L":      {0x57, 0xAD, 0, 1, 1, func(v uint8) { cpu.L = v }, SUB_A_L_OPCODE},
		"HL_MEM": {0x34, 0x12, 0, 0, 0, func(v uint8) { cpu.H = 0xA0; cpu.L = v; cpu.MMU.Write(cpu.readHL(), v) }, SUB_A_HLMEM_OPCODE},
		"A":      {0x80, 0x80, 1, 0, 0, func(v uint8) { cpu.A = v }, SUB_A_A_OPCODE},
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
			if cpu.readZFlag() != test.expZ {
				t.Errorf("Z flag: got %x, expected %x", cpu.readZFlag(), test.expZ)
			}
			if cpu.readHFlag() != test.expH {
				t.Errorf("H flag: got %x, expected %x", cpu.readHFlag(), test.expH)
			}
			if cpu.readCFlag() != test.expC {
				t.Errorf("C flag: got %x, expected %x", cpu.readCFlag(), test.expC)
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
		expZ   uint8
		expH   uint8
		expC   uint8
		setR8  func(uint8)
		opcode uint8
	}

	tests := map[string]test{
		"B":      {0x08, 0x08, 0, 1, 0, 0, func(v uint8) { cpu.B = v }, SBC_A_B_OPCODE},
		"C":      {0x10, 0x0F, 1, 1, 1, 0, func(v uint8) { cpu.C = v }, SBC_A_C_OPCODE},
		"D":      {0x10, 0x0F, 0, 0, 1, 0, func(v uint8) { cpu.D = v }, SBC_A_D_OPCODE},
		"E":      {0x00, 0x00, 1, 0, 1, 1, func(v uint8) { cpu.E = v }, SBC_A_E_OPCODE},
		"H":      {0x80, 0x80, 0, 1, 0, 0, func(v uint8) { cpu.H = v }, SBC_A_H_OPCODE},
		"L":      {0x57, 0xAD, 1, 0, 1, 1, func(v uint8) { cpu.L = v }, SBC_A_L_OPCODE},
		"HL_MEM": {0x34, 0x14, 1, 0, 1, 0, func(v uint8) { cpu.H = 0xA0; cpu.H = 0xA0; cpu.L = v; cpu.MMU.Write(cpu.readHL(), v) }, SBC_A_HLMEM_OPCODE},
		"A":      {0x80, 0x80, 1, 0, 1, 1, func(v uint8) { cpu.A = v }, SBC_A_A_OPCODE},
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
			if cpu.readZFlag() != test.expZ {
				t.Errorf("Z flag: got %x, expected %x", cpu.readZFlag(), test.expZ)
			}
			if cpu.readHFlag() != test.expH {
				t.Errorf("H flag: got %x, expected %x", cpu.readHFlag(), test.expH)
			}
			if cpu.readCFlag() != test.expC {
				t.Errorf("C flag: got %x, expected %x", cpu.readCFlag(), test.expC)
			}
		})
	}
}

func Test_AND_A_R8(t *testing.T) {
	cpu := mockCPU()

	type test struct {
		A      uint8
		R8     uint8
		expZ   uint8
		setR8  func(uint8)
		opcode uint8
	}

	tests := map[string]test{
		"B":      {0x08, 0x08, 0, func(v uint8) { cpu.B = v }, AND_A_B_OPCODE},
		"C":      {0x10, 0x0F, 1, func(v uint8) { cpu.C = v }, AND_A_C_OPCODE},
		"D":      {0x10, 0x1F, 0, func(v uint8) { cpu.D = v }, AND_A_D_OPCODE},
		"E":      {0x00, 0x00, 1, func(v uint8) { cpu.E = v }, AND_A_E_OPCODE},
		"H":      {0xAA, 0x55, 1, func(v uint8) { cpu.H = v }, AND_A_H_OPCODE},
		"L":      {0x57, 0xAD, 0, func(v uint8) { cpu.L = v }, AND_A_L_OPCODE},
		"HL_MEM": {0x34, 0x14, 0, func(v uint8) { cpu.H = 0xA0; cpu.L = v; cpu.MMU.Write(cpu.readHL(), v) }, AND_A_HLMEM_OPCODE},
		"A":      {0x80, 0x80, 0, func(v uint8) { cpu.A = v }, AND_A_A_OPCODE},
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
			if cpu.readZFlag() != test.expZ {
				t.Errorf("Z flag: got %x, expected %x", cpu.readZFlag(), test.expZ)
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
		expZ   uint8
		setR8  func(uint8)
		opcode uint8
	}

	tests := map[string]test{
		"B":      {0x08, 0x08, 1, func(v uint8) { cpu.B = v }, XOR_A_B_OPCODE},
		"C":      {0x10, 0x0F, 0, func(v uint8) { cpu.C = v }, XOR_A_C_OPCODE},
		"D":      {0x10, 0x1F, 0, func(v uint8) { cpu.D = v }, XOR_A_D_OPCODE},
		"E":      {0x00, 0x00, 1, func(v uint8) { cpu.E = v }, XOR_A_E_OPCODE},
		"H":      {0xAA, 0x55, 0, func(v uint8) { cpu.H = v }, XOR_A_H_OPCODE},
		"L":      {0x57, 0xAD, 0, func(v uint8) { cpu.L = v }, XOR_A_L_OPCODE},
		"HL_MEM": {0x34, 0x14, 0, func(v uint8) { cpu.H = 0xA0; cpu.L = v; cpu.MMU.Write(cpu.readHL(), v) }, XOR_A_HLMEM_OPCODE},
		"A":      {0x80, 0x80, 1, func(v uint8) { cpu.A = v }, XOR_A_A_OPCODE},
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
			if cpu.readZFlag() != test.expZ {
				t.Errorf("Z flag: got %x, expected %x", cpu.readZFlag(), test.expZ)
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
		expZ   uint8
		setR8  func(uint8)
		opcode uint8
	}

	tests := map[string]test{
		"B":      {0x08, 0x08, 0, func(v uint8) { cpu.B = v }, OR_A_B_OPCODE},
		"C":      {0x10, 0x0F, 0, func(v uint8) { cpu.C = v }, OR_A_C_OPCODE},
		"D":      {0x10, 0x1F, 0, func(v uint8) { cpu.D = v }, OR_A_D_OPCODE},
		"E":      {0x00, 0x00, 1, func(v uint8) { cpu.E = v }, OR_A_E_OPCODE},
		"H":      {0xAA, 0x55, 0, func(v uint8) { cpu.H = v }, OR_A_H_OPCODE},
		"L":      {0x57, 0xAD, 0, func(v uint8) { cpu.L = v }, OR_A_L_OPCODE},
		"HL_MEM": {0x34, 0x14, 0, func(v uint8) { cpu.H = 0xA0; cpu.L = v; cpu.MMU.Write(cpu.readHL(), v) }, OR_A_HLMEM_OPCODE},
		"A":      {0x80, 0x80, 0, func(v uint8) { cpu.A = v }, OR_A_A_OPCODE},
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
			if cpu.readZFlag() != test.expZ {
				t.Errorf("Z flag: got %x, expected %x", cpu.readZFlag(), test.expZ)
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
		expZ   uint8
		expH   uint8
		expC   uint8
		setR8  func(uint8)
		opcode uint8
	}

	tests := map[string]test{
		"B":      {0x08, 0x08, 1, 0, 0, func(v uint8) { cpu.B = v }, CP_A_B_OPCODE},
		"C":      {0x0F, 0x10, 0, 0, 1, func(v uint8) { cpu.C = v }, CP_A_C_OPCODE},
		"D":      {0x00, 0x01, 0, 1, 1, func(v uint8) { cpu.D = v }, CP_A_D_OPCODE},
		"E":      {0xF0, 0x01, 0, 1, 0, func(v uint8) { cpu.E = v }, CP_A_E_OPCODE},
		"H":      {0x80, 0x08, 0, 1, 0, func(v uint8) { cpu.H = v }, CP_A_H_OPCODE},
		"L":      {0x57, 0xAD, 0, 1, 1, func(v uint8) { cpu.L = v }, CP_A_L_OPCODE},
		"HL_MEM": {0x34, 0x12, 0, 0, 0, func(v uint8) { cpu.H = 0xA0; cpu.L = v; cpu.MMU.Write(cpu.readHL(), v) }, CP_A_HLMEM_OPCODE},
		"A":      {0x80, 0x80, 1, 0, 0, func(v uint8) { cpu.A = v }, CP_A_A_OPCODE},
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
			if cpu.readZFlag() != test.expZ {
				t.Errorf("Z flag: got %x, expected %x", cpu.readZFlag(), test.expZ)
			}
			if cpu.readHFlag() != test.expH {
				t.Errorf("H flag: got %x, expected %x", cpu.readHFlag(), test.expH)
			}
			if cpu.readCFlag() != test.expC {
				t.Errorf("C flag: got %x, expected %x", cpu.readCFlag(), test.expC)
			}
		})
	}
}

func Test_ADD_A_N8(t *testing.T) {
	cpu := mockCPU()

	var A uint8 = 0x57
	var n8 uint8 = 0xAD
	var expZ uint8 = 0
	var expH uint8 = 1
	var expC uint8 = 1

	cpu.A = A
	writeTestProgram(cpu, ADD_A_N8_OPCODE, n8)
	cpu.ExecuteInstruction()
	if cpu.A != A+n8 {
		t.Fatalf("got %2X, expected %2X", cpu.A, A+n8)
	}
	if cpu.readNFlag() != 0 {
		t.Errorf("N flag: got %x, expected %x", cpu.readNFlag(), 0)
	}
	if cpu.readZFlag() != expZ {
		t.Errorf("Z flag: got %x, expected %x", cpu.readZFlag(), expZ)
	}
	if cpu.readHFlag() != expH {
		t.Errorf("H flag: got %x, expected %x", cpu.readHFlag(), expH)
	}
	if cpu.readCFlag() != expC {
		t.Errorf("C flag: got %x, expected %x", cpu.readCFlag(), expC)
	}
}

func Test_ADC_A_N8(t *testing.T) {
	cpu := mockCPU()

	var A uint8 = 0x08
	var n8 uint8 = 0x08
	var carry uint8 = 1
	var expZ uint8 = 0
	var expH uint8 = 1
	var expC uint8 = 0

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
	if cpu.readZFlag() != expZ {
		t.Errorf("Z flag: got %x, expected %x", cpu.readZFlag(), expZ)
	}
	if cpu.readHFlag() != expH {
		t.Errorf("H flag: got %x, expected %x", cpu.readHFlag(), expH)
	}
	if cpu.readCFlag() != expC {
		t.Errorf("C flag: got %x, expected %x", cpu.readCFlag(), expC)
	}
}

func Test_SUB_A_N8(t *testing.T) {
	cpu := mockCPU()

	var A uint8 = 0x80
	var n8 uint8 = 0x08
	var expZ uint8 = 0
	var expH uint8 = 1
	var expC uint8 = 0

	cpu.A = A
	writeTestProgram(cpu, SUB_A_N8_OPCODE, n8)
	cpu.ExecuteInstruction()
	if cpu.A != A-n8 {
		t.Fatalf("got %2X, expected %2X", cpu.A, A-n8)
	}
	if cpu.readNFlag() != 1 {
		t.Errorf("N flag: got %x, expected %x", cpu.readNFlag(), 1)
	}
	if cpu.readZFlag() != expZ {
		t.Errorf("Z flag: got %x, expected %x", cpu.readZFlag(), expZ)
	}
	if cpu.readHFlag() != expH {
		t.Errorf("H flag: got %x, expected %x", cpu.readHFlag(), expH)
	}
	if cpu.readCFlag() != expC {
		t.Errorf("C flag: got %x, expected %x", cpu.readCFlag(), expC)
	}
}

func Test_SBC_A_N8(t *testing.T) {
	cpu := mockCPU()

	var A uint8 = 0x10
	var n8 uint8 = 0x0F
	var carry uint8 = 1
	var expZ uint8 = 1
	var expH uint8 = 1
	var expC uint8 = 0

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
	if cpu.readZFlag() != expZ {
		t.Errorf("Z flag: got %x, expected %x", cpu.readZFlag(), expZ)
	}
	if cpu.readHFlag() != expH {
		t.Errorf("H flag: got %x, expected %x", cpu.readHFlag(), expH)
	}
	if cpu.readCFlag() != expC {
		t.Errorf("C flag: got %x, expected %x", cpu.readCFlag(), expC)
	}
}

func Test_AND_A_N8(t *testing.T) {
	cpu := mockCPU()

	var A uint8 = 0xAA
	var n8 uint8 = 0x55
	var expZ uint8 = 1

	cpu.A = A
	writeTestProgram(cpu, AND_A_N8_OPCODE, n8)
	cpu.ExecuteInstruction()
	expected := A & n8
	if cpu.A != expected {
		t.Fatalf("got %2X, expected %2X", cpu.A, expected)
	}
	if cpu.readZFlag() != expZ {
		t.Errorf("Z flag: got %x, expected %x", cpu.readZFlag(), expZ)
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
	var expZ uint8 = 1

	cpu.A = A
	writeTestProgram(cpu, XOR_A_N8_OPCODE, n8)
	cpu.ExecuteInstruction()
	expected := A ^ n8
	if cpu.A != expected {
		t.Fatalf("got %2X, expected %2X", cpu.A, expected)
	}
	if cpu.readZFlag() != expZ {
		t.Errorf("Z flag: got %x, expected %x", cpu.readZFlag(), expZ)
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
	var expZ uint8 = 0

	cpu.A = A
	writeTestProgram(cpu, OR_A_N8_OPCODE, n8)
	cpu.ExecuteInstruction()
	expected := A | n8
	if cpu.A != expected {
		t.Fatalf("got %2X, expected %2X", cpu.A, expected)
	}
	if cpu.readZFlag() != expZ {
		t.Errorf("Z flag: got %x, expected %x", cpu.readZFlag(), expZ)
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
	var expZ uint8 = 0
	var expH uint8 = 1
	var expC uint8 = 0

	cpu.A = A
	writeTestProgram(cpu, CP_A_N8_OPCODE, n8)
	cpu.ExecuteInstruction()
	if cpu.A != A {
		t.Fatalf("got %2X, expected %2X", cpu.A, A)
	}
	if cpu.readNFlag() != 1 {
		t.Errorf("N flag: got %x, expected %x", cpu.readNFlag(), 1)
	}
	if cpu.readZFlag() != expZ {
		t.Errorf("Z flag: got %x, expected %x", cpu.readZFlag(), expZ)
	}
	if cpu.readHFlag() != expH {
		t.Errorf("H flag: got %x, expected %x", cpu.readHFlag(), expH)
	}
	if cpu.readCFlag() != expC {
		t.Errorf("C flag: got %x, expected %x", cpu.readCFlag(), expC)
	}
}

func Test_POP_R16STK(t *testing.T) {
	cpu := mockCPU()

	t.Run("BC", func(t *testing.T) {
		addr := uint16(0x4321)
		expectedBC := addr
		expectedSP := cpu.SP

		// write stack pointer
		cpu.SP -= 2
		cpu.MMU.WriteWord(cpu.SP, addr)

		writeTestProgram(cpu, POP_BC_OPCODE)
		cpu.ExecuteInstruction()

		if cpu.readBC() != expectedBC {
			t.Errorf("got BC=%04X, expected %04X", cpu.readBC(), expectedBC)
		}
		if cpu.SP != expectedSP {
			t.Errorf("got SP=%04X, expected %04X", cpu.SP, expectedSP)
		}
	})

	t.Run("DE", func(t *testing.T) {
		addr := uint16(0x1234)
		expectedDE := addr
		expectedSP := cpu.SP

		// write stack pointer
		cpu.SP -= 2
		cpu.MMU.WriteWord(cpu.SP, addr)

		writeTestProgram(cpu, POP_DE_OPCODE)
		cpu.ExecuteInstruction()

		if cpu.readDE() != expectedDE {
			t.Errorf("got DE=%04X, expected %04X", cpu.readDE(), expectedDE)
		}
		if cpu.SP != expectedSP {
			t.Errorf("got SP=%04X, expected %04X", cpu.SP, expectedSP)
		}
	})

	t.Run("HL", func(t *testing.T) {
		addr := uint16(0x1111)
		expectedHL := addr
		expectedSP := cpu.SP

		// write stack pointer
		cpu.SP -= 2
		cpu.MMU.WriteWord(cpu.SP, addr)

		writeTestProgram(cpu, POP_HL_OPCODE)
		cpu.ExecuteInstruction()

		if cpu.readHL() != expectedHL {
			t.Errorf("got HL=%04X, expected %04X", cpu.readHL(), expectedHL)
		}
		if cpu.SP != expectedSP {
			t.Errorf("got SP=%04X, expected %04X", cpu.SP, expectedSP)
		}
	})

	t.Run("AF", func(t *testing.T) {
		addr := uint16(0x56F0)
		expectedAF := addr
		expectedSP := cpu.SP

		// write stack pointer
		cpu.SP -= 2
		cpu.MMU.WriteWord(cpu.SP, addr)

		writeTestProgram(cpu, POP_AF_OPCODE)
		cpu.ExecuteInstruction()

		if cpu.readAF() != expectedAF {
			t.Errorf("got AF=%04X, expected %04X", cpu.readAF(), expectedAF)
		}
		if cpu.SP != expectedSP {
			t.Errorf("got SP=%04X, expected %04X", cpu.SP, expectedSP)
		}

		// Test registers
		got, expected := cpu.readZFlag(), util.ReadBit(uint8(addr), ZFlagBit)
		if got != expected {
			t.Errorf("Z flag: got %x, expected %x", got, expected)
		}
		got, expected = cpu.readNFlag(), util.ReadBit(uint8(addr), NFlagBit)
		if got != expected {
			t.Errorf("N flag: got %x, expected %x", got, expected)
		}
		got, expected = cpu.readHFlag(), util.ReadBit(uint8(addr), HFlagBit)
		if got != expected {
			t.Errorf("H flag: got %x, expected %x", got, expected)
		}
		got, expected = cpu.readCFlag(), util.ReadBit(uint8(addr), CFlagBit)
		if got != expected {
			t.Errorf("C flag: got %x, expected %x", got, expected)
		}
	})
}

func Test_PUSH_R16STK(t *testing.T) {
	cpu := mockCPU()

	t.Run("BC", func(t *testing.T) {
		addr := uint16(0x4321)
		expectedSP := cpu.SP - 2
		cpu.writeBC(addr)

		writeTestProgram(cpu, PUSH_BC_OPCODE)
		cpu.ExecuteInstruction()

		if cpu.MMU.ReadWord(cpu.SP) != addr {
			t.Errorf("got [SP]=%04X, expected %04X", cpu.MMU.ReadWord(cpu.SP), addr)
		}
		if cpu.SP != expectedSP {
			t.Errorf("got SP=%04X, expected %04X", cpu.SP, expectedSP)
		}
	})

	t.Run("DE", func(t *testing.T) {
		addr := uint16(0x1234)
		expectedSP := cpu.SP - 2
		cpu.writeDE(addr)

		writeTestProgram(cpu, PUSH_DE_OPCODE)
		cpu.ExecuteInstruction()

		if cpu.MMU.ReadWord(cpu.SP) != addr {
			t.Errorf("got [SP]=%04X, expected %04X", cpu.MMU.ReadWord(cpu.SP), addr)
		}
		if cpu.SP != expectedSP {
			t.Errorf("got SP=%04X, expected %04X", cpu.SP, expectedSP)
		}
	})

	t.Run("HL", func(t *testing.T) {
		addr := uint16(0x1111)
		expectedSP := cpu.SP - 2
		cpu.writeHL(addr)

		writeTestProgram(cpu, PUSH_HL_OPCODE)
		cpu.ExecuteInstruction()

		if cpu.MMU.ReadWord(cpu.SP) != addr {
			t.Errorf("got [SP]=%04X, expected %04X", cpu.MMU.ReadWord(cpu.SP), addr)
		}
		if cpu.SP != expectedSP {
			t.Errorf("got SP=%04X, expected %04X", cpu.SP, expectedSP)
		}
	})

	t.Run("AF", func(t *testing.T) {
		addr := uint16(0x2220)
		expectedSP := cpu.SP - 2
		cpu.writeAF(addr)

		writeTestProgram(cpu, PUSH_AF_OPCODE)
		cpu.ExecuteInstruction()

		if cpu.MMU.ReadWord(cpu.SP) != addr {
			t.Errorf("got [SP]=%04X, expected %04X", cpu.MMU.ReadWord(cpu.SP), addr)
		}
		if cpu.SP != expectedSP {
			t.Errorf("got SP=%04X, expected %04X", cpu.SP, expectedSP)
		}
	})
}

func Test_RET_COND(t *testing.T) {
	cpu := mockCPU()

	// Prepare flags for each condition
	boolToInt := map[bool]uint8{false: 0, true: 1}
	conditions := map[string]func(bool){
		"Z":  func(cc bool) { cpu.setZFlag(boolToInt[cc]) },
		"NZ": func(cc bool) { cpu.setZFlag(boolToInt[!cc]) },
		"C":  func(cc bool) { cpu.setCFlag(boolToInt[cc]) },
		"NC": func(cc bool) { cpu.setCFlag(boolToInt[!cc]) },
	}
	opcodes := map[string]uint8{
		"Z":  RET_Z_OPCODE,
		"NZ": RET_NZ_OPCODE,
		"C":  RET_C_OPCODE,
		"NC": RET_NC_OPCODE,
	}

	for cond, setFlag := range conditions {
		var addr uint16 = 0x1000
		// write stack pointer
		cpu.SP -= 2
		cpu.MMU.WriteWord(cpu.SP, addr)

		t.Run(cond+"_unmet", func(t *testing.T) {
			setFlag(false)

			expectedPC := int(cpu.PC) + OPCODES_BYTES[opcodes[cond]]
			expectedSP := cpu.SP

			writeTestProgram(cpu, opcodes[cond])
			cpu.ExecuteInstruction()

			if cpu.PC != uint16(expectedPC) {
				t.Errorf("unmet: got PC=%04X, expected %04X", cpu.PC, expectedPC)
			}
			if cpu.SP != expectedSP {
				t.Errorf("unmet: got SP=%04X, expected %04X", cpu.SP, expectedSP)
			}
		})

		t.Run(cond+"_met", func(t *testing.T) {
			setFlag(true)

			expectedPC := addr
			expectedSP := cpu.SP + 2

			writeTestProgram(cpu, opcodes[cond])
			cpu.ExecuteInstruction()

			if cpu.PC != expectedPC {
				t.Errorf("met: got PC=%04X, expected %04X", cpu.PC, expectedPC)
			}
			if cpu.SP != expectedSP {
				t.Errorf("met: got SP=%04X, expected %04X", cpu.SP, expectedSP)
			}
		})
	}
}

func Test_RET(t *testing.T) {
	cpu := mockCPU()

	addr := uint16(0x4321)
	expectedPC := addr
	expectedSP := cpu.SP

	// write stack pointer
	cpu.SP -= 2
	cpu.MMU.WriteWord(cpu.SP, addr)

	writeTestProgram(cpu, RET_OPCODE)
	cpu.ExecuteInstruction()

	if cpu.PC != expectedPC {
		t.Errorf("got PC=%04X, expected %04X", cpu.PC, expectedPC)
	}
	if cpu.SP != expectedSP {
		t.Errorf("got SP=%04X, expected %04X", cpu.SP, expectedSP)
	}
}

func Test_RETI(t *testing.T) {
	cpu := mockCPU()
	cpu.IME = false

	addr := uint16(0x4321)
	expectedPC := addr
	expectedSP := cpu.SP

	// write stack pointer
	cpu.SP -= 2
	cpu.MMU.WriteWord(cpu.SP, addr)

	writeTestProgram(cpu, RETI_OPCODE)
	cpu.ExecuteInstruction()

	if cpu.PC != expectedPC {
		t.Errorf("got PC=%04X, expected %04X", cpu.PC, expectedPC)
	}
	if cpu.SP != expectedSP {
		t.Errorf("got SP=%04X, expected %04X", cpu.SP, expectedSP)
	}
	if cpu.IME != true {
		t.Error("Interrupts not enabled")
	}
}

func Test_JP_COND_N16(t *testing.T) {
	cpu := mockCPU()

	// Prepare flags for each condition
	boolToInt := map[bool]uint8{false: 0, true: 1}
	conditions := map[string]func(bool){
		"Z":  func(cc bool) { cpu.setZFlag(boolToInt[cc]) },
		"NZ": func(cc bool) { cpu.setZFlag(boolToInt[!cc]) },
		"C":  func(cc bool) { cpu.setCFlag(boolToInt[cc]) },
		"NC": func(cc bool) { cpu.setCFlag(boolToInt[!cc]) },
	}
	opcodes := map[string]uint8{
		"Z":  JP_Z_N16_OPCODE,
		"NZ": JP_NZ_N16_OPCODE,
		"C":  JP_C_N16_OPCODE,
		"NC": JP_NC_N16_OPCODE,
	}

	for cond, setFlag := range conditions {
		var addr uint16 = 0x1000

		t.Run(cond+"_unmet", func(t *testing.T) {
			setFlag(false)
			addr += 2

			expectedPC := int(cpu.PC) + OPCODES_BYTES[opcodes[cond]]

			writeTestProgram(cpu, opcodes[cond], uint8(addr), uint8(addr>>8))
			cpu.ExecuteInstruction()

			if cpu.PC != uint16(expectedPC) {
				t.Errorf("unmet: got PC=%04X, expected %04X", cpu.PC, expectedPC)
			}
		})

		t.Run(cond+"_met", func(t *testing.T) {
			setFlag(true)
			addr += 2

			expectedPC := addr

			writeTestProgram(cpu, opcodes[cond], uint8(addr), uint8(addr>>8))
			cpu.ExecuteInstruction()

			if cpu.PC != expectedPC {
				t.Errorf("met: got PC=%04X, expected %04X", cpu.PC, expectedPC)
			}
		})
	}
}

func Test_JP_N16(t *testing.T) {
	cpu := mockCPU()

	var addr uint16 = 0x4321
	expectedPC := addr

	writeTestProgram(cpu, JP_N16_OPCODE, uint8(addr), uint8(addr>>8))
	cpu.ExecuteInstruction()

	if cpu.PC != expectedPC {
		t.Errorf("got PC=%04X, expected %04X", cpu.PC, expectedPC)
	}
}

func Test_JP_HL(t *testing.T) {
	cpu := mockCPU()

	var addr uint16 = 0x1234
	expectedPC := addr

	writeTestProgram(cpu, JP_HL_OPCODE)
	cpu.writeHL(addr)
	cpu.ExecuteInstruction()

	if cpu.PC != expectedPC {
		t.Errorf("got PC=%04X, expected %04X", cpu.PC, expectedPC)
	}
}

func Test_CALL_COND_N16(t *testing.T) {
	cpu := mockCPU()

	// Prepare flags for each condition
	boolToInt := map[bool]uint8{false: 0, true: 1}
	conditions := map[string]func(bool){
		"Z":  func(cc bool) { cpu.setZFlag(boolToInt[cc]) },
		"NZ": func(cc bool) { cpu.setZFlag(boolToInt[!cc]) },
		"C":  func(cc bool) { cpu.setCFlag(boolToInt[cc]) },
		"NC": func(cc bool) { cpu.setCFlag(boolToInt[!cc]) },
	}
	opcodes := map[string]uint8{
		"Z":  CALL_Z_N16_OPCODE,
		"NZ": CALL_NZ_N16_OPCODE,
		"C":  CALL_C_N16_OPCODE,
		"NC": CALL_NC_N16_OPCODE,
	}

	for cond, setFlag := range conditions {
		var addr uint16 = 0x1111

		t.Run(cond+"_unmet", func(t *testing.T) {
			setFlag(false)
			addr += 2

			expectedPC := cpu.PC + uint16(OPCODES_BYTES[opcodes[cond]])
			expectedSP := cpu.SP
			expectedStack := cpu.MMU.ReadWord(cpu.SP)

			writeTestProgram(cpu, opcodes[cond], uint8(addr), uint8(addr>>8))
			cpu.ExecuteInstruction()

			if cpu.PC != expectedPC {
				t.Errorf("got PC=%04X, expected %04X", cpu.PC, expectedPC)
			}
			if cpu.SP != expectedSP {
				t.Errorf("got SP=%04X, expected %04X", cpu.SP, expectedSP)
			}
			if cpu.MMU.ReadWord(cpu.SP) != expectedStack {
				t.Errorf("got stack=%04X, expected %04X", cpu.MMU.ReadWord(cpu.SP), expectedStack)
			}
		})

		t.Run(cond+"_met", func(t *testing.T) {
			setFlag(true)
			addr += 2

			expectedPC := addr
			expectedSP := cpu.SP - 2
			expectedStack := cpu.PC + uint16(OPCODES_BYTES[opcodes[cond]])

			writeTestProgram(cpu, opcodes[cond], uint8(addr), uint8(addr>>8))
			cpu.ExecuteInstruction()

			if cpu.PC != expectedPC {
				t.Errorf("got PC=%04X, expected %04X", cpu.PC, expectedPC)
			}
			if cpu.SP != expectedSP {
				t.Errorf("got SP=%04X, expected %04X", cpu.SP, expectedSP)
			}
			if cpu.MMU.ReadWord(cpu.SP) != expectedStack {
				t.Errorf("got stack=%04X, expected %04X", cpu.MMU.ReadWord(cpu.SP), expectedStack)
			}
		})
	}
}

func Test_CALL_N16(t *testing.T) {
	cpu := mockCPU()

	var addr uint16 = 0x1111

	expectedPC := addr
	expectedSP := cpu.SP - 2
	expectedStack := cpu.PC + uint16(OPCODES_BYTES[CALL_N16_OPCODE])

	writeTestProgram(cpu, CALL_N16_OPCODE, uint8(addr), uint8(addr>>8))
	cpu.ExecuteInstruction()

	if cpu.PC != expectedPC {
		t.Errorf("got PC=%04X, expected %04X", cpu.PC, expectedPC)
	}
	if cpu.SP != expectedSP {
		t.Errorf("got SP=%04X, expected %04X", cpu.SP, expectedSP)
	}
	if cpu.MMU.ReadWord(cpu.SP) != expectedStack {
		t.Errorf("got stack=%04X, expected %04X", cpu.MMU.ReadWord(cpu.SP), expectedStack)
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
			expectedPC := addr
			expectedSP := cpu.SP - 2

			writeTestProgram(cpu, opcode)
			cpu.ExecuteInstruction()

			if cpu.PC != uint16(expectedPC) {
				t.Errorf("got PC=%04X, expected %04X", cpu.PC, expectedPC)
			}
			if cpu.SP != expectedSP {
				t.Errorf("got SP=%04X, expected %04X", cpu.SP, expectedSP)
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

	got, expected := cpu.MMU.Read(0xFF00+uint16(cpu.C)), A
	if got != expected {
		t.Fatalf("got %02X, expected %02X", got, expected)
	}
}

func Test_LDH_A_C(t *testing.T) {
	cpu := mockCPU()
	var value, C uint8 = 0xD1, 0x12

	cpu.MMU.Write(0xFF00+uint16(C), value)
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

	got, expected := cpu.MMU.Read(0xFF00+uint16(offset)), A
	if got != expected {
		t.Fatalf("got %02X, expected %02X", got, expected)
	}
}

func Test_LDH_A_N8(t *testing.T) {
	cpu := mockCPU()
	var value, offset uint8 = 0xD1, 0x12

	cpu.MMU.Write(0xFF00+uint16(offset), value)
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

	got, expected := cpu.MMU.Read(addr), A
	if got != expected {
		t.Fatalf("got %02X, expected %02X", got, expected)
	}
}

func Test_LD_A_N16(t *testing.T) {
	cpu := mockCPU()
	var value uint8 = 0xD1
	var addr uint16 = 0xDCBA

	cpu.MMU.Write(addr, value)
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
		SP   uint16
		e8   int8
		expH uint8
		expC uint8
	}{
		"neg-with-carry":    {SP: 0xFFFE, e8: -1, expH: 1, expC: 1},
		"neg-without-carry": {SP: 0xFF00, e8: -1, expH: 0, expC: 0},
		"pos-with-carry":    {SP: 0xF0FF, e8: 0x0F, expH: 1, expC: 1},
		"pos-without-carry": {SP: 0xFF00, e8: 0x0F, expH: 0, expC: 0},
		"overflow":          {SP: 0xFFFF, e8: 1, expH: 1, expC: 1},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			cpu.SP = test.SP
			writeTestProgram(cpu, ADD_SP_E8_OPCODE, uint8(test.e8))

			cpu.ExecuteInstruction()
			expectedSP := int(test.SP) + int(test.e8)
			if cpu.SP != uint16(expectedSP) {
				t.Fatalf("got %04X, expected %04X", cpu.SP, expectedSP)
			}
			if cpu.readZFlag() != 0 {
				t.Fatalf("Z flag: got %x, expected %x", cpu.readZFlag(), 0)
			}
			if cpu.readNFlag() != 0 {
				t.Fatalf("N flag: got %x, expected %x", cpu.readNFlag(), 0)
			}
			if cpu.readHFlag() != test.expH {
				t.Fatalf("H flag: got %x, expected %x", cpu.readHFlag(), test.expH)
			}
			if cpu.readCFlag() != test.expC {
				t.Fatalf("C flag: got %x, expected %x", cpu.readCFlag(), test.expC)
			}
		})
	}
}

func Test_LD_HL_SP_E8(t *testing.T) {
	cpu := mockCPU()

	var SP uint16 = 0xFFFE
	var e8 int8 = -1
	var expH, expC uint8 = 1, 1

	cpu.SP = SP
	writeTestProgram(cpu, LD_HL_SP_E8_OPCODE, uint8(e8))

	cpu.ExecuteInstruction()
	expectedHL := int(SP) + int(e8)
	if cpu.readHL() != uint16(expectedHL) {
		t.Fatalf("got %04X, expected %04X", cpu.readHL(), expectedHL)
	}
	if cpu.readZFlag() != 0 {
		t.Fatalf("Z flag: got %x, expected %x", cpu.readZFlag(), 0)
	}
	if cpu.readNFlag() != 0 {
		t.Fatalf("N flag: got %x, expected %x", cpu.readNFlag(), 0)
	}
	if cpu.readHFlag() != expH {
		t.Fatalf("H flag: got %x, expected %x", cpu.readHFlag(), expH)
	}
	if cpu.readCFlag() != expC {
		t.Fatalf("C flag: got %x, expected %x", cpu.readCFlag(), expC)
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

var r8Offset = map[string]uint8{"B": 0, "C": 1, "D": 2, "E": 3, "H": 4, "L": 5, "HL_MEM": 6, "A": 7}
var r8Name = [8]string{"B", "C", "D", "E", "H", "L", "HL_MEM", "A"}

func Test_RLC_R8(t *testing.T) {
	cpu := mockCPU()
	// Test flag reset
	cpu.F = 0xFF

	r8 := "B"
	r8Id := r8Offset[r8]

	tests := map[string]struct {
		R8    uint8
		expR8 uint8
		expC  uint8
		expZ  uint8
	}{
		"carry":    {R8: 0b10101010, expR8: 0b01010101, expC: 1, expZ: 0},
		"no-carry": {R8: 0b01110001, expR8: 0b11100010, expC: 0, expZ: 0},
		"zero":     {R8: 0, expR8: 0, expC: 0, expZ: 1},
	}

	for name, test := range tests {
		t.Run(r8+"/"+name, func(t *testing.T) {
			cpu.writeR8(r8Id, test.R8)
			writeTestProgram(cpu, PREFIX_OPCODE, RLC_R8_OPCODE+r8Id)
			cpu.ExecuteInstruction()

			if cpu.readR8(r8Id) != test.expR8 {
				t.Fatalf("got %08b, expected %08b", cpu.readR8(r8Id), test.expR8)
			}
			// Test flags
			if cpu.readZFlag() != test.expZ {
				t.Errorf("Z flag: got %x, expected %x", cpu.readZFlag(), test.expZ)
			}
			if cpu.readNFlag() != 0 {
				t.Error("N flag: got 1, expected 0")
			}
			if cpu.readHFlag() != 0 {
				t.Error("H flag: got 1, expected 0")
			}
			if cpu.readCFlag() != test.expC {
				t.Errorf("C flag: got %x, expected %x", cpu.readCFlag(), test.expC)
			}
		})
	}
}

func Test_RRC_R8(t *testing.T) {
	cpu := mockCPU()
	// Test flag reset
	cpu.F = 0xFF

	r8 := "C"
	r8Id := r8Offset[r8]

	tests := map[string]struct {
		R8    uint8
		expR8 uint8
		expC  uint8
		expZ  uint8
	}{
		"carry":    {R8: 0b10101010, expR8: 0b01010101, expC: 0, expZ: 0},
		"no-carry": {R8: 0b01110001, expR8: 0b10111000, expC: 1, expZ: 0},
		"zero":     {R8: 0, expR8: 0, expC: 0, expZ: 1},
	}

	for name, test := range tests {
		t.Run(r8+"/"+name, func(t *testing.T) {
			cpu.writeR8(r8Id, test.R8)
			writeTestProgram(cpu, PREFIX_OPCODE, RRC_R8_OPCODE+r8Id)
			cpu.ExecuteInstruction()

			if cpu.readR8(r8Id) != test.expR8 {
				t.Fatalf("got %08b, expected %08b", cpu.readR8(r8Id), test.expR8)
			}
			// Test flags
			if cpu.readZFlag() != test.expZ {
				t.Errorf("Z flag: got %x, expected %x", cpu.readZFlag(), test.expZ)
			}
			if cpu.readNFlag() != 0 {
				t.Error("N flag: got 1, expected 0")
			}
			if cpu.readHFlag() != 0 {
				t.Error("H flag: got 1, expected 0")
			}
			if cpu.readCFlag() != test.expC {
				t.Errorf("C flag: got %x, expected %x", cpu.readCFlag(), test.expC)
			}
		})
	}
}

func Test_RL_R8(t *testing.T) {
	cpu := mockCPU()
	// Test flag reset
	cpu.F = 0xFF

	r8 := "D"
	r8Id := r8Offset[r8]

	tests := map[string]struct {
		R8    uint8
		carry uint8
		expR8 uint8
		expC  uint8
		expZ  uint8
	}{
		"carry_set":   {R8: 0b01011010, carry: 1, expR8: 0b10110101, expC: 0, expZ: 0},
		"carry_unset": {R8: 0b11011010, carry: 0, expR8: 0b10110100, expC: 1, expZ: 0},
		"carry":       {R8: 0b11011010, carry: 1, expR8: 0b10110101, expC: 1, expZ: 0},
		"no-carry":    {R8: 0b01011010, carry: 0, expR8: 0b10110100, expC: 0, expZ: 0},
		"zero":        {R8: 0b10000000, carry: 0, expR8: 0, expC: 1, expZ: 1},
	}

	for name, test := range tests {
		t.Run(r8+"/"+name, func(t *testing.T) {
			cpu.writeR8(r8Id, test.R8)
			cpu.setCFlag(test.carry)
			writeTestProgram(cpu, PREFIX_OPCODE, RL_R8_OPCODE+r8Id)
			cpu.ExecuteInstruction()

			if cpu.readR8(r8Id) != test.expR8 {
				t.Fatalf("got %08b, expected %08b", cpu.readR8(r8Id), test.expR8)
			}
			// Test flags
			if cpu.readZFlag() != test.expZ {
				t.Errorf("Z flag: got %x, expected %x", cpu.readZFlag(), test.expZ)
			}
			if cpu.readNFlag() != 0 {
				t.Error("N flag: got 1, expected 0")
			}
			if cpu.readHFlag() != 0 {
				t.Error("H flag: got 1, expected 0")
			}
			if cpu.readCFlag() != test.expC {
				t.Errorf("C flag: got %x, expected %x", cpu.readCFlag(), test.expC)
			}
		})
	}
}

func Test_RR_R8(t *testing.T) {
	cpu := mockCPU()
	// Test flag reset
	cpu.F = 0xFF

	r8 := "E"
	r8Id := r8Offset[r8]

	tests := map[string]struct {
		R8    uint8
		carry uint8
		expR8 uint8
		expC  uint8
		expZ  uint8
	}{
		"carry_set":   {R8: 0b01011010, carry: 1, expR8: 0b10101101, expC: 0, expZ: 0},
		"carry_unset": {R8: 0b11011011, carry: 0, expR8: 0b01101101, expC: 1, expZ: 0},
		"carry":       {R8: 0b11011011, carry: 1, expR8: 0b11101101, expC: 1, expZ: 0},
		"no-carry":    {R8: 0b01011010, carry: 0, expR8: 0b00101101, expC: 0, expZ: 0},
		"zero":        {R8: 0b00000001, carry: 0, expR8: 0, expC: 1, expZ: 1},
	}

	for name, test := range tests {
		t.Run(r8+"/"+name, func(t *testing.T) {
			cpu.writeR8(r8Id, test.R8)
			cpu.setCFlag(test.carry)
			writeTestProgram(cpu, PREFIX_OPCODE, RR_R8_OPCODE+r8Id)
			cpu.ExecuteInstruction()

			if cpu.readR8(r8Id) != test.expR8 {
				t.Fatalf("got %08b, expected %08b", cpu.readR8(r8Id), test.expR8)
			}
			// Test flags
			if cpu.readZFlag() != test.expZ {
				t.Errorf("Z flag: got %x, expected %x", cpu.readZFlag(), test.expZ)
			}
			if cpu.readNFlag() != 0 {
				t.Error("N flag: got 1, expected 0")
			}
			if cpu.readHFlag() != 0 {
				t.Error("H flag: got 1, expected 0")
			}
			if cpu.readCFlag() != test.expC {
				t.Errorf("C flag: got %x, expected %x", cpu.readCFlag(), test.expC)
			}
		})
	}
}

func Test_SLA_R8(t *testing.T) {
	cpu := mockCPU()
	r8 := "H"
	r8Id := r8Offset[r8]

	tests := map[string]struct {
		R8    uint8
		expR8 uint8
		expC  uint8
		expZ  uint8
	}{
		"carry":    {R8: 0b11011011, expR8: 0b10110110, expC: 1, expZ: 0},
		"no-carry": {R8: 0b01011010, expR8: 0b10110100, expC: 0, expZ: 0},
		"zero":     {R8: 0b10000000, expR8: 0, expC: 1, expZ: 1},
	}

	for name, test := range tests {
		t.Run(r8+"/"+name, func(t *testing.T) {
			cpu.writeR8(r8Id, test.R8)
			writeTestProgram(cpu, PREFIX_OPCODE, SLA_R8_OPCODE+r8Id)
			cpu.ExecuteInstruction()

			if cpu.readR8(r8Id) != test.expR8 {
				t.Fatalf("got %08b, expected %08b", cpu.readR8(r8Id), test.expR8)
			}
			// Test flags
			if cpu.readZFlag() != test.expZ {
				t.Errorf("Z flag: got %x, expected %x", cpu.readZFlag(), test.expZ)
			}
			if cpu.readNFlag() != 0 {
				t.Error("N flag: got 1, expected 0")
			}
			if cpu.readHFlag() != 0 {
				t.Error("H flag: got 1, expected 0")
			}
			if cpu.readCFlag() != test.expC {
				t.Errorf("C flag: got %x, expected %x", cpu.readCFlag(), test.expC)
			}
		})
	}
}

func Test_SRA_R8(t *testing.T) {
	cpu := mockCPU()
	r8 := "L"
	r8Id := r8Offset[r8]

	tests := map[string]struct {
		R8    uint8
		expR8 uint8
		expC  uint8
		expZ  uint8
	}{
		"carry":    {R8: 0b11011011, expR8: 0b11101101, expC: 1, expZ: 0},
		"no-carry": {R8: 0b01011010, expR8: 0b00101101, expC: 0, expZ: 0},
		"zero":     {R8: 0b00000001, expR8: 0, expC: 1, expZ: 1},
	}

	for name, test := range tests {
		t.Run(r8+"/"+name, func(t *testing.T) {
			cpu.writeR8(r8Id, test.R8)
			writeTestProgram(cpu, PREFIX_OPCODE, SRA_R8_OPCODE+r8Id)
			cpu.ExecuteInstruction()

			if cpu.readR8(r8Id) != test.expR8 {
				t.Fatalf("got %08b, expected %08b", cpu.readR8(r8Id), test.expR8)
			}
			// Test flags
			if cpu.readZFlag() != test.expZ {
				t.Errorf("Z flag: got %x, expected %x", cpu.readZFlag(), test.expZ)
			}
			if cpu.readNFlag() != 0 {
				t.Error("N flag: got 1, expected 0")
			}
			if cpu.readHFlag() != 0 {
				t.Error("H flag: got 1, expected 0")
			}
			if cpu.readCFlag() != test.expC {
				t.Errorf("C flag: got %x, expected %x", cpu.readCFlag(), test.expC)
			}
		})
	}
}

func Test_SWAP_R8(t *testing.T) {
	cpu := mockCPU()
	r8 := "HL_MEM"
	r8Id := r8Offset[r8]

	tests := map[string]struct {
		R8    uint8
		expR8 uint8
		expZ  uint8
	}{
		"non-zero": {R8: 0xAF, expR8: 0xFA, expZ: 0},
		"zero":     {R8: 0, expR8: 0, expZ: 1},
	}

	for name, test := range tests {
		t.Run(r8+"/"+name, func(t *testing.T) {
			addr := uint16(0xF000) + uint16(test.R8)
			cpu.writeHL(addr)
			cpu.writeR8(r8Id, test.R8)
			writeTestProgram(cpu, PREFIX_OPCODE, SWAP_R8_OPCODE+r8Id)
			cpu.ExecuteInstruction()

			if cpu.readR8(r8Id) != test.expR8 {
				t.Fatalf("got %2X, expected %2X", cpu.readR8(r8Id), test.expR8)
			}
			// Test flags
			if cpu.readZFlag() != test.expZ {
				t.Errorf("Z flag: got %x, expected %x", cpu.readZFlag(), test.expZ)
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
	r8 := "A"
	r8Id := r8Offset[r8]

	tests := map[string]struct {
		R8    uint8
		expR8 uint8
		expC  uint8
		expZ  uint8
	}{
		"carry":    {R8: 0b11011011, expR8: 0b01101101, expC: 1, expZ: 0},
		"no-carry": {R8: 0b01011010, expR8: 0b00101101, expC: 0, expZ: 0},
		"zero":     {R8: 0b00000001, expR8: 0, expC: 1, expZ: 1},
	}

	for name, test := range tests {
		t.Run(r8+"/"+name, func(t *testing.T) {
			cpu.writeR8(r8Id, test.R8)
			writeTestProgram(cpu, PREFIX_OPCODE, SRL_R8_OPCODE+r8Id)
			cpu.ExecuteInstruction()

			if cpu.readR8(r8Id) != test.expR8 {
				t.Fatalf("got %08b, expected %08b", cpu.readR8(r8Id), test.expR8)
			}
			// Test flags
			if cpu.readZFlag() != test.expZ {
				t.Errorf("Z flag: got %x, expected %x", cpu.readZFlag(), test.expZ)
			}
			if cpu.readNFlag() != 0 {
				t.Error("N flag: got 1, expected 0")
			}
			if cpu.readHFlag() != 0 {
				t.Error("H flag: got 1, expected 0")
			}
			if cpu.readCFlag() != test.expC {
				t.Errorf("C flag: got %x, expected %x", cpu.readCFlag(), test.expC)
			}
		})
	}
}

func Test_BIT_B3_R8(t *testing.T) {
	cpu := mockCPU()

	tests := [8]struct {
		opcode uint8
		R8     uint8
		expZ   uint8
	}{
		{opcode: BIT_0_R8_OPCODE + 0, R8: 0b11011010, expZ: 1},
		{opcode: BIT_1_R8_OPCODE + 1, R8: 0b01011010, expZ: 0},
		{opcode: BIT_2_R8_OPCODE + 2, R8: 0b00000101, expZ: 0},
		{opcode: BIT_3_R8_OPCODE + 3, R8: 0b11010011, expZ: 1},
		{opcode: BIT_4_R8_OPCODE + 4, R8: 0b01001010, expZ: 1},
		{opcode: BIT_5_R8_OPCODE + 5, R8: 0b00100001, expZ: 0},
		{opcode: BIT_6_R8_OPCODE + 6, R8: 0b10011011, expZ: 1},
		{opcode: BIT_7_R8_OPCODE + 7, R8: 0b11011010, expZ: 0},
	}

	for i, test := range tests {
		t.Run("B"+fmt.Sprint(i)+"/"+r8Name[i], func(t *testing.T) {
			cpu.H = 0xA0
			cpu.writeR8(test.opcode, test.R8)
			writeTestProgram(cpu, PREFIX_OPCODE, test.opcode)
			cpu.ExecuteInstruction()

			if cpu.readZFlag() != test.expZ {
				t.Errorf("Z flag: got %x, expected %x", cpu.readZFlag(), test.expZ)
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
		expR8  uint8
	}{
		{opcode: RES_0_R8_OPCODE + 0, expR8: 0b11111110},
		{opcode: RES_1_R8_OPCODE + 1, expR8: 0b11111101},
		{opcode: RES_2_R8_OPCODE + 2, expR8: 0b11111011},
		{opcode: RES_3_R8_OPCODE + 3, expR8: 0b11110111},
		{opcode: RES_4_R8_OPCODE + 4, expR8: 0b11101111},
		{opcode: RES_5_R8_OPCODE + 5, expR8: 0b11011111},
		{opcode: RES_6_R8_OPCODE + 6, expR8: 0b10111111},
		{opcode: RES_7_R8_OPCODE + 7, expR8: 0b01111111},
	}

	for i, test := range tests {
		t.Run("B"+fmt.Sprint(i)+"/"+r8Name[i], func(t *testing.T) {
			cpu.H = 0xA0
			cpu.writeR8(test.opcode, 0xFF)
			writeTestProgram(cpu, PREFIX_OPCODE, test.opcode)
			cpu.ExecuteInstruction()

			if cpu.readR8(test.opcode) != test.expR8 {
				t.Fatalf("got %08b, expected %08b", cpu.readR8(test.opcode), test.expR8)
			}
		})
	}
}

func Test_SET_B3_R8(t *testing.T) {
	cpu := mockCPU()

	tests := [8]struct {
		opcode uint8
		expR8  uint8
	}{
		{opcode: SET_0_R8_OPCODE + 0, expR8: 0b00000001},
		{opcode: SET_1_R8_OPCODE + 1, expR8: 0b00000010},
		{opcode: SET_2_R8_OPCODE + 2, expR8: 0b00000100},
		{opcode: SET_3_R8_OPCODE + 3, expR8: 0b00001000},
		{opcode: SET_4_R8_OPCODE + 4, expR8: 0b00010000},
		{opcode: SET_5_R8_OPCODE + 5, expR8: 0b00100000},
		{opcode: SET_6_R8_OPCODE + 6, expR8: 0b01000000},
		{opcode: SET_7_R8_OPCODE + 7, expR8: 0b10000000},
	}

	for i, test := range tests {
		t.Run("B"+fmt.Sprint(i)+"/"+r8Name[i], func(t *testing.T) {
			cpu.H = 0xA0
			cpu.writeR8(test.opcode, 0)
			writeTestProgram(cpu, PREFIX_OPCODE, test.opcode)
			cpu.ExecuteInstruction()

			if cpu.readR8(test.opcode) != test.expR8 {
				t.Fatalf("got %08b, expected %08b", cpu.readR8(test.opcode), test.expR8)
			}
		})
	}
}
