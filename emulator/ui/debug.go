package ui

import (
	"fmt"
	"log"
	"os/exec"
)

// DebuggerState implements debug.Debugger interface
type DebuggerState struct {
	active bool

	paused      bool
	step        bool
	breakpoints map[uint16]struct{}
}

func NewDebuggerState() *DebuggerState {
	return &DebuggerState{
		breakpoints: make(map[uint16]struct{}),
	}
}

func (state *DebuggerState) Pause() {
	state.active = true

	log.Println("Pausing...")
	state.paused = true
}

func (state *DebuggerState) Resume() {
	state.active = false
	state.paused = false
	state.step = false
}

func (state *DebuggerState) Step() {
	log.Println("Stepping...")
	state.step = true
}

func (state *DebuggerState) Continue() {
	log.Println("Continuing...")
	state.paused = false
}

func (state *DebuggerState) Breakpoint(addr uint16, set bool) {
	log.Println("Breakpoint:", addr, "set:", set)
	if set {
		state.breakpoints[addr] = struct{}{}
	} else {
		delete(state.breakpoints, addr)
	}
}

// startDebugger starts the debugger client
func (ui *UI) startDebugger() error {
	if ui.debuggerCmd != nil {
		return fmt.Errorf("debugger already started")
	}
	ui.debuggerCmd = exec.Command("go", "run", "./cmd/main.go")
	ui.debuggerCmd.Dir = "../../debugger/"

	go func() {
		out, err := ui.debuggerCmd.Output()
		if err != nil {
			log.Println("debugger error:", err)
		}
		fmt.Println(string(out))

		ui.debuggerCmd = nil
	}()
	return nil
}

func (ui *UI) stopDebugger() error {
	if ui.debuggerCmd != nil {
		return ui.debuggerCmd.Process.Kill()
	}
	return nil
}
