package debugger

import (
	"github.com/danielecanzoneri/gb-emulator/pkg/debug"
	"log"
)

// Debugger implements debug.Debugger interface
type Debugger struct {
	active bool

	paused      bool
	breakpoints map[uint16]struct{}

	// Game boy memory and register
	cpu debug.CPUDebugger
	mem debug.MemoryDebugger

	// Function to execute on step
	onStep func()
}

func NewDebugger(cpu debug.CPUDebugger, mem debug.MemoryDebugger, onStep func()) *Debugger {
	return &Debugger{
		breakpoints: make(map[uint16]struct{}),
		cpu:         cpu,
		mem:         mem,
		onStep:      onStep,
	}
}

func (debugger *Debugger) IsActive() bool {
	return debugger.active
}

func (debugger *Debugger) Pause() {
	debugger.active = true

	log.Println("Pausing...")
	debugger.paused = true
}

func (debugger *Debugger) Paused() bool {
	return debugger.paused
}

func (debugger *Debugger) Resume() {
	debugger.active = false
	debugger.paused = false
}

func (debugger *Debugger) Step() {
	log.Println("Stepping...")
	debugger.onStep()
}

func (debugger *Debugger) Continue() {
	log.Println("Continuing...")
	debugger.paused = false
}

func (debugger *Debugger) Breakpoint(addr uint16, set bool) {
	log.Println("Breakpoint:", addr, "set:", set)
	if set {
		debugger.breakpoints[addr] = struct{}{}
	} else {
		delete(debugger.breakpoints, addr)
	}
}

func (debugger *Debugger) IsBreakpoint(addr uint16) bool {
	_, ok := debugger.breakpoints[addr]
	return ok
}

func (debugger *Debugger) GetState() *debug.GameBoyState {
	state := new(debug.GameBoyState)
	for i := range 0x10000 {
		state.Memory[uint16(i)] = debugger.mem.DebugRead(uint16(i))
	}
	state.AF = debugger.cpu.ReadAF()
	state.BC = debugger.cpu.ReadBC()
	state.DE = debugger.cpu.ReadDE()
	state.HL = debugger.cpu.ReadHL()
	state.PC = debugger.cpu.ReadPC()
	state.SP = debugger.cpu.ReadSP()
	state.IME = debugger.cpu.InterruptsEnabled()
	return state
}
