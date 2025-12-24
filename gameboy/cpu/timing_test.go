package cpu

import "testing"

type CycleCounter struct {
	ticks uint
}

func (c *CycleCounter) Tick(ticks uint) {
	c.ticks += ticks
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
		cpu.AddTicker(counter)
		writeTestProgram(cpu, uint8(opcode))
		cpu.ExecuteInstruction()

		if counter.ticks != OPCODES_CYCLES[opcode]*4 {
			t.Errorf("Wrong timing for opcode %02X: got %d, expected %d", opcode, counter.ticks, OPCODES_CYCLES[opcode]*4)
		}
	}
}

func TestPrefixedOpcodesTiming(t *testing.T) {
	for opcode := range 0x100 {
		cpu := mockCPU()
		counter := &CycleCounter{}
		cpu.AddTicker(counter)
		writeTestProgram(cpu, PREFIX_OPCODE, uint8(opcode))
		cpu.ExecuteInstruction()

		expectedCycles := OPCODES_CYCLES[PREFIX_OPCODE] + PREFIX_OPCODES_CYCLES[opcode]
		if counter.ticks != expectedCycles*4 {
			t.Errorf("Wrong timing for opcode CB %02X: got %d, expected %d", opcode, counter.ticks, expectedCycles*4)
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
			cpu.AddTicker(counter)

			op := opcode + uint8(offset*8)

			// No branch
			setCondition(false)
			writeTestProgram(cpu, op)

			cpu.ExecuteInstruction()
			if counter.ticks != OPCODES_CYCLES[op]*4 {
				t.Errorf("opcode %02X no branch: got %d Cycles, expected %d", op, counter.ticks, OPCODES_CYCLES[op]*4)
			}

			// Branch
			setCondition(true)
			writeTestProgram(cpu, op)

			counter.ticks = 0
			cpu.ExecuteInstruction()
			if counter.ticks != OPCODES_CYCLES_BRANCH[op]*4 {
				t.Errorf("opcode %02X branch: got %d Cycles, expected %d", op, counter.ticks, OPCODES_CYCLES_BRANCH[op]*4)
			}
		}
	}
}
