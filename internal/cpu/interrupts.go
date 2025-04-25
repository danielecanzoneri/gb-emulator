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

	handlerCycles = 5
)

var interruptsHandler = map[uint8]uint16{
	vblankMask: vblankHandler,
	statMask:   statHandler,
	timerMask:  timerHandler,
	serialMask: serialHandler,
	joypadMask: joypadHandler,
}

var RequestTimerInterruptFunc = func(cpu *CPU) func() {
	return func() {
		cpu.MMU.Write(ifAddr, cpu.MMU.Read(ifAddr)|timerMask)
	}
}

func (cpu *CPU) handleInterrupts() uint {
	defer cpu.handleIME()

	// Check for pending interrupts
	triggered := cpu.MMU.Read(ieAddr) & cpu.MMU.Read(ifAddr)

	// Awake if halted
	if cpu.halted && triggered > 0 {
		cpu.halted = false
	}

	if !cpu.IME || triggered == 0 {
		return 0
	}

	// Serve interrupts with priority
	for i := range 5 {
		mask := uint8(1 << i)
		if triggered&mask > 0 {
			cpu.serveInterrupt(mask)
			return handlerCycles
		}
	}
	return 0
}

func (cpu *CPU) requestInterrupt(interruptMask uint8) {
	IF := cpu.MMU.Read(ifAddr)
	IF |= interruptMask
	cpu.MMU.Write(ifAddr, IF)
}

func (cpu *CPU) serveInterrupt(interruptMask uint8) {
	IF := cpu.MMU.Read(ifAddr)
	IF &= ^interruptMask
	cpu.MMU.Write(ifAddr, IF)
	cpu.IME = false

	cpu.PUSH_STACK(cpu.PC)
	cpu.PC = interruptsHandler[interruptMask]
}

func (cpu *CPU) handleIME() {
	if cpu._EIDelayed {
		cpu._EIDelayed = false
		cpu.IME = true
	}
}
