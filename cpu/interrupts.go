package cpu

func (cpu *CPU) handleInterrupts() {
	if cpu._EIDelayed {
		cpu._EIDelayed = false
		cpu.IME = true
	}
}
