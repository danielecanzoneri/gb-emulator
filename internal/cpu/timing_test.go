package cpu

import "testing"

type CycleCounter struct {
	Cycles uint
}

func (c *CycleCounter) Cycle() {
	c.Cycles++
}

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
		counter := &CycleCounter{}
		cpu.AddCyclable(counter)
		writeTestProgram(cpu, uint8(opcode))
		cpu.ExecuteInstruction()

		if counter.Cycles != OPCODES_CYCLES[opcode] {
			t.Errorf("Wrong timing for opcode %02X: got %d, expected %d", opcode, counter.Cycles, OPCODES_CYCLES[opcode])
		}
	}
}

func TestPrefixedOpcodesTiming(t *testing.T) {
	for opcode := range 0x100 {
		cpu := mockCPU()
		counter := &CycleCounter{}
		cpu.AddCyclable(counter)
		writeTestProgram(cpu, PREFIX_OPCODE, uint8(opcode))
		cpu.ExecuteInstruction()

		expectedCycles := OPCODES_CYCLES[PREFIX_OPCODE] + PREFIX_OPCODES_CYCLES[opcode]
		if counter.Cycles != expectedCycles {
			t.Errorf("Wrong timing for opcode CB %02X: got %d, expected %d", opcode, counter.Cycles, expectedCycles)
		}
	}
}

func TestOpcodesWithBranchingTiming(t *testing.T) {
	cpu := mockCPU()

	boolToInt := map[bool]uint8{false: 0, true: 1}
	conditions := [4]func(bool){
		func(cc bool) { cpu.setZFlag(boolToInt[!cc]) }, // NZ
		func(cc bool) { cpu.setZFlag(boolToInt[cc]) },  // Z
		func(cc bool) { cpu.setCFlag(boolToInt[!cc]) }, // NC
		func(cc bool) { cpu.setCFlag(boolToInt[cc]) },  // C
	}

	var opcodesBase = [4]uint8{0x20, 0xC0, 0xC2, 0xC4}
	for _, opcode := range opcodesBase {
		for offset, setCondition := range conditions {
			counter := &CycleCounter{}
			cpu.AddCyclable(counter)

			op := opcode + uint8(offset*8)

			// No branch
			setCondition(false)
			writeTestProgram(cpu, op)

			cpu.ExecuteInstruction()
			if counter.Cycles != OPCODES_CYCLES[op] {
				t.Errorf("opcode %02X no branch: got %d Cycles, expected %d", op, counter.Cycles, OPCODES_CYCLES[op])
			}

			// Branch
			setCondition(true)
			writeTestProgram(cpu, op)

			counter.Cycles = 0
			cpu.ExecuteInstruction()
			if counter.Cycles != OPCODES_CYCLES_BRANCH[op] {
				t.Errorf("opcode %02X branch: got %d Cycles, expected %d", op, counter.Cycles, OPCODES_CYCLES_BRANCH[op])
			}
		}
	}
}
