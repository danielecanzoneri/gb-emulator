package timer

import "strconv"

const (
	divAddr  = 0xFF04
	timaAddr = 0xFF05
	tmaAddr  = 0xFF06
	tacAddr  = 0xFF07

	tacMask uint8 = 0b00000111
)

func (t *Timer) Write(addr uint16, v uint8) {
	switch addr {
	case divAddr:
		t.systemCounter = 0
	case timaAddr:
		if !t.timaReloaded {
			t.TIMA = v
			t.timaOverflow = false
		}
	case tmaAddr:
		t.TMA = v
		if t.timaReloaded {
			// New TMA is written to TIMA if write happens while reloading
			t.TIMA = v
		}
	case tacAddr:
		t.TAC = v
		t.detectFallingEdge()
	default:
		panic("timer: unknown addr " + strconv.Itoa(int(addr)))
	}
}

func (t *Timer) Read(addr uint16) uint8 {
	switch addr {
	case divAddr:
		return uint8(t.systemCounter >> 8)
	case timaAddr:
		return t.TIMA
	case tmaAddr:
		return t.TMA
	case tacAddr:
		return ^tacMask | (t.TAC & tacMask)
	default:
		panic("timer: unknown addr " + strconv.Itoa(int(addr)))
	}
}
