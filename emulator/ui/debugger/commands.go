package debugger

func (d *Debugger) Toggle() {
	d.Active = !d.Active
	d.Continue = false // In case toggling while continuing
}

func (d *Debugger) CheckBreakpoint(addr uint16) bool {
	return d.disassembler.IsBreakpoint(addr)
}
