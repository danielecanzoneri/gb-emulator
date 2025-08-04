package cpu

import "testing"

func TestInterruptHandlers(t *testing.T) {
	cpu := mockCPU()

	for mask, handler := range interruptsHandler {
		var PC, SP uint16 = 0x1234, 0xFFFE

		// Enable interrupts
		cpu.mmu.Write(ieAddr, mask)
		cpu.mmu.Write(ifAddr, mask)
		cpu.PC = PC
		cpu.SP = SP

		t.Run(string(mask)+"/IME=false", func(t *testing.T) {
			cpu.IME = false
			cpu.handleInterrupts()

			if cpu.PC != PC {
				t.Errorf("PC: got %04X, expected %04X", cpu.PC, PC)
			}
		})

		t.Run(string(mask)+"/IME=true", func(t *testing.T) {
			cpu.IME = true
			cpu.handleInterrupts()

			if cpu.PC != handler {
				t.Errorf("PC: got %04X, expected %04X", cpu.PC, handler)
			}

			retAddr := uint16(cpu.mmu.Read(SP-2)) | uint16(cpu.mmu.Read(SP-1))<<8
			if retAddr != PC {
				t.Errorf("wrong return address: got %04X, expected %04X", retAddr, PC)
			}

			if cpu.IME == true {
				t.Errorf("IME not disabled")
			}

			if cpu.mmu.Read(ifAddr)&mask > 0 {
				t.Errorf("interrupt flag not reset")
			}
		})
	}
}

func TestRequestInterrupt(t *testing.T) {
	cpu := mockCPU()

	// Test requesting an interrupt
	cpu.RequestInterrupt(VBlankInterruptMask)

	// Check if the interrupt flag is set
	if cpu.mmu.Read(ifAddr)&VBlankInterruptMask == 0 {
		t.Errorf("VBLANK interrupt not requested")
	}

	// Test requesting another interrupt
	cpu.RequestInterrupt(TimerInterruptMask)

	// Check if the interrupt flag is set
	if cpu.mmu.Read(ifAddr)&TimerInterruptMask == 0 {
		t.Errorf("TIMER interrupt not requested")
	}
}

func TestHALTBug(t *testing.T) {
	cpu := mockCPU()

	// When a halt instruction is executed with IME = 0 and [IE] & [IF] != 0 (interrupt pending),
	// the halt instruction ends immediately, but pc fails to be normally incremented.
	//
	// This causes the byte after the halt to be read a second time
	// (and this behaviour can repeat if said byte executes another halt instruction).
	t.Run("normal", func(t *testing.T) {
		PC := cpu.PC

		cpu.IME = false
		cpu.halted = false
		cpu.mmu.Write(ifAddr, 1)
		cpu.mmu.Write(ieAddr, 1)
		writeTestProgram(cpu, HALT_OPCODE, NOP_OPCODE, NOP_OPCODE)
		cpu.ExecuteInstruction() // Should immediately exit HALT mod
		cpu.ExecuteInstruction() // Should fail to increment PC

		if cpu.PC != PC+1 {
			t.Errorf("PC was wrongly incremented, got PC=%02X, expected %02X", cpu.PC, PC+1)
		}
	})
}
