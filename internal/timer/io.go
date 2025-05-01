package timer

import "strconv"

const (
	divAddr  = 0xFF04
	timaAddr = 0xFF05
	tmaAddr  = 0xFF06
	tacAddr  = 0xFF07
)

var tacClockToBitFallMapping = map[uint8]uint8{
	0b00: 9,
	0b01: 3,
	0b10: 5,
	0b11: 7,
}

func (t *Timer) Write(addr uint16, v uint8) {
	switch addr {
	case divAddr:
		t.systemCounter = 0
	case timaAddr:
		t.TIMA = v
		// t.timaCounter = 0
	case tmaAddr:
		t.TMA = v
	case tacAddr:
		t.TAC = v
		t.enabled = v&0b100 > 0
		t.bitFalling = tacClockToBitFallMapping[v&0b11]
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
		return t.TAC
	default:
		panic("timer: unknown addr " + strconv.Itoa(int(addr)))
	}
}
