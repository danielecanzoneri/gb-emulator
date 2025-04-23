package cpu

const (
	IF_ADDR = 0xFF0F
	IE_ADDR = 0xFFFF

	VBLANK_INT_MASK    = 0b1
	VBLANK_INT_HANDLER = 0x40

	STAT_INT_MASK    = 0b10
	STAT_INT_HANDLER = 0x48

	TIMER_INT_MASK    = 0b100
	TIMER_INT_HANDLER = 0x50

	SERIAL_INT_MASK    = 0b1000
	SERIAL_INT_HANDLER = 0x58

	JOYPAD_INT_MASK    = 0b10000
	JOYPAD_INT_HANDLER = 0x60

	INT_HANDLER_CYCLES = 5
)

var INT_HANDLERS = map[uint8]uint16{
	VBLANK_INT_MASK: VBLANK_INT_HANDLER,
	STAT_INT_MASK:   STAT_INT_HANDLER,
	TIMER_INT_MASK:  TIMER_INT_HANDLER,
	SERIAL_INT_MASK: SERIAL_INT_HANDLER,
	JOYPAD_INT_MASK: JOYPAD_INT_HANDLER,
}

func (cpu *CPU) handleInterrupts() bool {
	defer cpu.handleIME()

	// Check for pending interrupts
	triggered := cpu.Mem.Read(IE_ADDR) & cpu.Mem.Read(IF_ADDR)

	// Awake if halted
	if cpu.halted && triggered > 0 {
		cpu.halted = false
	}

	if !cpu.IME || triggered == 0 {
		return false
	}

	// Serve interrupts with priority
	for i := range 5 {
		mask := uint8(1 << i)
		if triggered&mask > 0 {
			cpu.serveInterrupt(mask)
			cpu.cycles += INT_HANDLER_CYCLES
			return true
		}
	}
	return false
}

func (cpu *CPU) requestInterrupt(interruptMask uint8) {
	IF := cpu.Mem.Read(IF_ADDR)
	IF |= interruptMask
	cpu.Mem.Write(IF_ADDR, IF)
}

func (cpu *CPU) serveInterrupt(interruptMask uint8) {
	IF := cpu.Mem.Read(IF_ADDR)
	IF &= ^interruptMask
	cpu.Mem.Write(IF_ADDR, IF)
	cpu.IME = false

	cpu.PUSH_STACK(cpu.PC)
	cpu.PC = INT_HANDLERS[interruptMask]
}

func (cpu *CPU) handleIME() {
	if cpu._EIDelayed {
		cpu._EIDelayed = false
		cpu.IME = true
	}
}
