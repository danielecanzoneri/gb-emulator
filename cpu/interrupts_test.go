package cpu

import "testing"

func TestInterruptHandlers(t *testing.T) {
	cpu := setup_CPU()

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
