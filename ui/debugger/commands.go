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
	defer d.Sync()

	d.gameBoy.CPU.ExecuteInstruction()
}

func (d *Debugger) Next() {
	d.Continue()
	d.NextInstruction = true
	d.CallDepth = 0
}

func (d *Debugger) Continue() {
	d.Continued = true

	// Unselect current entry
	d.disassembler.currentInstruction = -1

	// TODO Disable control buttons
	d.disassembler.refresh()
}

func (d *Debugger) Stop() {
	defer d.Sync()

	d.Continued = false
	// TODO Enable control buttons
}

func (d *Debugger) Reset() {
	defer d.Sync()

	d.gameBoy.Reset()
}
