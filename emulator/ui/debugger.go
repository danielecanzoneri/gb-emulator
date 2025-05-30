package ui

import (
	"fmt"
	"log"
	"os/exec"
)

// startDebugger starts the debugger client
func (ui *UI) startDebugger() error {
	if ui.debuggerCmd != nil {
		return fmt.Errorf("debugger already started")
	}
	ui.debuggerCmd = exec.Command("go", "run", "./cmd/main.go")
	ui.debuggerCmd.Dir = "../../debugger/"

	go func() {
		err := ui.debuggerCmd.Run()
		if err != nil {
			log.Println("debugger error:", err)
		}

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
