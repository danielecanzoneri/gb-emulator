package cpu

import "testing"

func TestOpcodesTiming(t *testing.T) {
	for opcode := range 0x100 {
		// Test opcodes that don't branch
		if _, ok := OPCODES_CYCLES_BRANCH[uint8(opcode)]; ok {
			continue
		}
		// Test existing opcodes
		if OPCODES_CYCLES[opcode] == 0 || opcode == PREFIX_OPCODE {
			continue
		}

		cpu := mockCPU()
		writeTestProgram(cpu, uint8(opcode))
		cycles := cpu.ExecuteInstruction()

		if cycles != OPCODES_CYCLES[opcode] {
			t.Errorf("Wrong timing for opcode %02X: got %d, expected %d", opcode, cycles, OPCODES_CYCLES[opcode])
		}
	}
}

func TestPrefixedOpcodesTiming(t *testing.T) {
	for opcode := range 0x100 {
		cpu := mockCPU()
		writeTestProgram(cpu, PREFIX_OPCODE, uint8(opcode))
		cycles := cpu.ExecuteInstruction()

		expectedCycles := OPCODES_CYCLES[PREFIX_OPCODE] + PREFIX_OPCODES_CYCLES[opcode]
		if cycles != expectedCycles {
			t.Errorf("Wrong timing for opcode %02X: got %d, expected %d", opcode, cycles, expectedCycles)
		}
	}
}

func TestOpcodesWithBranchingTiming(t *testing.T) {
	cpu := mockCPU()

	bool_to_int := map[bool]uint8{false: 0, true: 1}
	conditions := [4]func(bool){
		func(cc bool) { cpu.setZFlag(bool_to_int[!cc]) }, // NZ
		func(cc bool) { cpu.setZFlag(bool_to_int[cc]) },  // Z
		func(cc bool) { cpu.setCFlag(bool_to_int[!cc]) }, // NC
		func(cc bool) { cpu.setCFlag(bool_to_int[cc]) },  // C
	}

	var opcodesBase = [4]uint8{0x20, 0xC0, 0xC2, 0xC4}
	for _, opcode := range opcodesBase {
		for offset, setCondition := range conditions {
			op := opcode + uint8(offset*8)

			// No branch
			setCondition(false)
			writeTestProgram(cpu, op)

			cycles := cpu.ExecuteInstruction()
			if cycles != OPCODES_CYCLES[op] {
				t.Errorf("opcode %02X no branch: got %d cycles, expected %d", op, cycles, OPCODES_CYCLES[op])
			}

			// Branch
			setCondition(true)
			writeTestProgram(cpu, op)

			cycles = cpu.ExecuteInstruction()
			if cycles != OPCODES_CYCLES_BRANCH[op] {
				t.Errorf("opcode %02X branch: got %d cycles, expected %d", op, cycles, OPCODES_CYCLES_BRANCH[op])
			}
		}
	}
}
