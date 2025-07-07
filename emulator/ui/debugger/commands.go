package debugger

func (d *Debugger) Toggle() {
	d.Active = !d.Active
	if d.Active {
		d.Stop()
	}
}

func (d *Debugger) CheckBreakpoint(addr uint16) bool {
	return d.disassembler.IsBreakpoint(addr)
}

func (d *Debugger) Step() {
	d.gameBoy.CPU.ExecuteInstruction()
	d.Sync() // Show updated data
}

func (d *Debugger) Continue() {
	d.Continued = true
	// Disable control buttons
}

func (d *Debugger) Stop() {
	d.Continued = false
	// Enable control buttons

	d.Sync() // Show updated data
}
