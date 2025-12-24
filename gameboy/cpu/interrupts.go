package cpu

const (
	ifAddr = 0xFF0F
	ieAddr = 0xFFFF

	VBlankInterruptMask = 0b1
	vblankHandler       = 0x40

	STATInterruptMask = 0b10
	statHandler       = 0x48

	TimerInterruptMask = 0b100
	timerHandler       = 0x50

	SerialInterruptMask = 0b1000
	serialHandler       = 0x58

	JoypadInterruptMask = 0b10000
	joypadHandler       = 0x60
)

var interruptsHandler = map[uint8]uint16{
	VBlankInterruptMask: vblankHandler,
	STATInterruptMask:   statHandler,
	TimerInterruptMask:  timerHandler,
	SerialInterruptMask: serialHandler,
	JoypadInterruptMask: joypadHandler,
}

func (cpu *CPU) handleInterrupts() {
	defer cpu.handleIME()

	// Check for pending interrupts
	triggered := cpu.mmu.Read(ieAddr) & cpu.mmu.Read(ifAddr) & 0x1F

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

func (cpu *CPU) RequestInterrupt(interruptMask uint8) {
	IF := cpu.mmu.Read(ifAddr)
	IF |= interruptMask
	cpu.mmu.Write(ifAddr, IF)
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

		IF := cpu.mmu.Read(ifAddr)
		IF &= ^interruptMask
		cpu.mmu.Write(ifAddr, IF)
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

func (cpu *CPU) Halted() bool {
	return cpu.halted
}

func (cpu *CPU) SpeedSwitchHalted() bool {
	return cpu.speedSwitchHaltedTicks > 0
}
