package debug

type Debugger interface {
	Pause()
	Resume()
	Step()
	Continue()
	Breakpoint(uint16, bool)
}
