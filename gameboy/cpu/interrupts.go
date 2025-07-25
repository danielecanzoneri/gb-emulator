package cpu

const (
	ifAddr = 0xFF0F
	ieAddr = 0xFFFF

	vblankMask    = 0b1
	vblankHandler = 0x40

	statMask    = 0b10
	statHandler = 0x48

	timerMask    = 0b100
	timerHandler = 0x50

	serialMask    = 0b1000
	serialHandler = 0x58

	joypadMask    = 0b10000
	joypadHandler = 0x60
)

var interruptsHandler = map[uint8]uint16{
	vblankMask: vblankHandler,
	statMask:   statHandler,
	timerMask:  timerHandler,
	serialMask: serialHandler,
	joypadMask: joypadHandler,
}

var (
	RequestTimerInterruptFunc = func(cpu *CPU) func() {
		return func() {
			cpu.MMU.Write(ifAddr, cpu.MMU.Read(ifAddr)|timerMask)
		}
	}
	RequestVBlankInterruptFunc = func(cpu *CPU) func() {
		return func() {
			cpu.MMU.Write(ifAddr, cpu.MMU.Read(ifAddr)|vblankMask)
		}
	}
	RequestSTATInterruptFunc = func(cpu *CPU) func() {
		return func() {
			cpu.MMU.Write(ifAddr, cpu.MMU.Read(ifAddr)|statMask)
		}
	}
)

func (cpu *CPU) handleInterrupts() {
	defer cpu.handleIME()

	// Check for pending interrupts
	triggered := cpu.MMU.Read(ieAddr) & cpu.MMU.Read(ifAddr) & 0x1F

	// Awake if halted
	if cpu.halted && triggered > 0 {
		cpu.halted = false
	}

	if !cpu.IME || triggered == 0 {
		return
	}

	// Serve interrupts with priority
	for i := range 5 {
		mask := uint8(1 << i)
		if triggered&mask > 0 {
			// If interrupt is cancelled, serve other
			if cpu.serveInterrupt(mask) {
				return
			}
		}
	}
	return
}

func (cpu *CPU) requestInterrupt(interruptMask uint8) {
	IF := cpu.MMU.Read(ifAddr)
	IF |= interruptMask
	cpu.MMU.Write(ifAddr, IF)
}

// serveInterrupt return false if interrupt is cancelled, true otherwise
func (cpu *CPU) serveInterrupt(interruptMask uint8) bool {
	cpu.interruptMaskRequested = interruptMask
	cpu.IME = false
	if cpu.callHook != nil {
		cpu.callHook()
	}

	// 2 NOP cycles (one is executed in cpu.PUSH_STACK)
	cpu.Tick(4)
	cpu.PUSH_STACK(cpu.PC)

	defer cpu.Tick(4) // Internal (set PC)

	// Check if interrupt is still requested, otherwise set PC to 0000
	if !cpu.interruptCancelled {
		cpu.PC = interruptsHandler[interruptMask]

		IF := cpu.MMU.Read(ifAddr)
		IF &= ^interruptMask
		cpu.MMU.Write(ifAddr, IF)
		return true
	} else {
		cpu.PC = 0
		return false
	}
}

func (cpu *CPU) handleIME() {
	if cpu._EIDelayed {
		cpu._EIDelayed = false
		cpu.IME = true
	}
}
