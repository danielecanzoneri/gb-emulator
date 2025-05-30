package debug

type CPUDebugger interface {
	ReadAF() uint16
	ReadBC() uint16
	ReadDE() uint16
	ReadHL() uint16
	ReadSP() uint16
	ReadPC() uint16
	InterruptsEnabled() bool
}

type MemoryDebugger interface {
	DebugRead(uint16) uint8
}

type Debugger interface {
	Pause()
	Resume()
	Step()
	Continue()
	Breakpoint(uint16, bool)
}
