package cpu

func (cpu *CPU) SetHooks(callHook func(), retHook func()) {
	cpu.callHook = callHook
	cpu.retHook = retHook
}
