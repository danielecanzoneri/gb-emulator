package debugger

func (d *Debugger) Toggle() {
	d.Active = !d.Active
	if d.Active {
		defer d.Sync()
		d.Stop()
	}
}

func (d *Debugger) CheckBreakpoint(addr uint16) bool {
	return d.disassembler.IsBreakpoint(addr)
}

// Run commands

func (d *Debugger) Step() {
	if d.Running {
		return
	}

	defer d.Sync()

	d.gameBoy.CPU.ExecuteInstruction()
}

func (d *Debugger) Next() {
	if d.Running {
		return
	}

	d.Continue()
	d.NextInstruction = true
	d.CallDepth = 0
}

func (d *Debugger) Continue() {
	if d.Running {
		return
	}

	d.Running = true

	// Unselect current entry
	d.disassembler.currentInstruction = -1

	// TODO Disable control buttons
	d.disassembler.refresh()
}

func (d *Debugger) Stop() {
	if !d.Running {
		return
	}

	defer d.Sync()

	d.Running = false
	// TODO Enable control buttons
}

func (d *Debugger) Reset() {
	defer d.Sync()

	// Stop if active
	d.Stop()

	d.gameBoy.Reset()
	d.initHooks()
}

func (d *Debugger) initHooks() {
	callHook := func() {
		d.CallDepth++
	}
	retHook := func() {
		d.CallDepth--
	}
	d.gameBoy.CPU.SetHooks(callHook, retHook)
}
