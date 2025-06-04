package server

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
}

func (s *Server) SetDebugger(cpu debug.CPUDebugger, mem debug.MemoryDebugger) {
	s.debugger = &Debugger{
		breakpoints: make(map[uint16]struct{}),
		cpu:         cpu,
		mem:         mem,
	}
}

func (s *Server) IsActive() bool {
	return s.debugger.active
}

func (s *Server) Pause() {
	s.debugger.active = true

	log.Println("Pausing...")
	s.debugger.paused = true
}

func (s *Server) Paused() bool {
	return s.debugger.paused
}

func (s *Server) Resume() {
	s.debugger.active = false
	s.debugger.paused = false
}

func (s *Server) Step() {
	log.Println("Stepping...")
	s.OnStep()
}

func (s *Server) Continue() {
	log.Println("Continuing...")
	s.debugger.paused = false
}

func (s *Server) Breakpoint(addr uint16, set bool) {
	log.Printf("Breakpoint: %04X, set: %v\n", addr, set)
	if set {
		s.debugger.breakpoints[addr] = struct{}{}
	} else {
		delete(s.debugger.breakpoints, addr)
	}
}

func (s *Server) IsBreakpoint(addr uint16) bool {
	_, ok := s.debugger.breakpoints[addr]
	return ok
}

func (s *Server) BreakpointHit() {
	s.Pause()
	s.sendBreakpointHit()
	s.sendState()
}

func (s *Server) GetState() *debug.GameBoyState {
	state := new(debug.GameBoyState)
	for i := range 0x10000 {
		state.Memory[uint16(i)] = s.debugger.mem.DebugRead(uint16(i))
	}
	state.AF = s.debugger.cpu.ReadAF()
	state.BC = s.debugger.cpu.ReadBC()
	state.DE = s.debugger.cpu.ReadDE()
	state.HL = s.debugger.cpu.ReadHL()
	state.PC = s.debugger.cpu.ReadPC()
	state.SP = s.debugger.cpu.ReadSP()
	state.IME = s.debugger.cpu.InterruptsEnabled()
	return state
}
