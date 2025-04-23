package cpu

import "testing"

func TestInterruptHandlers(t *testing.T) {
	cpu := mockCPU()

	for mask, handler := range INT_HANDLERS {
		var PC, SP uint16 = 0x1234, 0xFFFE

		// Enable interrupts
		cpu.Mem.Write(IE_ADDR, mask)
		cpu.Mem.Write(IF_ADDR, mask)
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

			retAddr := uint16(cpu.Mem.Read(SP-2)) | uint16(cpu.Mem.Read(SP-1))<<8
			if retAddr != PC {
				t.Errorf("wrong return address: got %04X, expected %04X", retAddr, PC)
			}

			if cpu.IME == true {
				t.Errorf("IME not disabled")
			}

			if cpu.Mem.Read(IF_ADDR)&mask > 0 {
				t.Errorf("interrupt flag not reset")
			}
		})
	}
}

func TestRequestInterrupt(t *testing.T) {
	cpu := mockCPU()

	// Test requesting an interrupt
	cpu.requestInterrupt(VBLANK_INT_MASK)

	// Check if the interrupt flag is set
	if cpu.Mem.Read(IF_ADDR)&VBLANK_INT_MASK == 0 {
		t.Errorf("VBLANK interrupt not requested")
	}

	// Test requesting another interrupt
	cpu.requestInterrupt(TIMER_INT_MASK)

	// Check if the interrupt flag is set
	if cpu.Mem.Read(IF_ADDR)&TIMER_INT_MASK == 0 {
		t.Errorf("TIMER interrupt not requested")
	}
}

func TestHALTBug(t *testing.T) {
	t.Error("TODO")
}
