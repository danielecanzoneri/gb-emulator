package cpu

const (
	DIV_ADDR  = 0xFF04
	TIMA_ADDR = 0xFF05
	TMA_ADDR  = 0xFF06
	TAC_ADDR  = 0xFF07

	DIV_FREQ = 64
)

type Timer struct {
	cpu *CPU

	divCounter  uint
	timaCounter uint
}

func InitTimer(cpu *CPU) *Timer {
	return &Timer{cpu: cpu}
}

func (t *Timer) Update(cycles uint) {
	// Update DIV
	t.divCounter += cycles
	t.updateDIV()

	// Update TIMA if enabled
	enabled, freq := t.getFreq()
	if enabled {
		t.timaCounter += cycles
		t.updateTIMA(freq)
	}
}

func (t *Timer) updateDIV() {
	if t.divCounter >= DIV_FREQ {
		t.cpu.Mem.Write(DIV_ADDR, t.cpu.Mem.Read(DIV_ADDR)+1)
		t.divCounter -= DIV_FREQ
	}
}

func (t *Timer) updateTIMA(timaFreq uint) {
	for t.timaCounter >= timaFreq {
		timaValue := t.cpu.Mem.Read(TIMA_ADDR) + 1
		t.timaCounter -= timaFreq

		// Check overflow
		if timaValue == 0 {
			timaValue = t.cpu.Mem.Read(TMA_ADDR)
			t.cpu.requestInterrupt(TIMER_INT_MASK)
		}

		t.cpu.Mem.Write(TIMA_ADDR, timaValue)
	}
}

var TAC_CLOCK_MAPPING = map[uint8]uint{
	0b00: 256,
	0b01: 4,
	0b10: 16,
	0b11: 64,
}

func (t *Timer) getFreq() (bool, uint) {
	TAC := t.cpu.Mem.Read(TAC_ADDR)
	if TAC&0b100 == 0 {
		// Timer disabled
		return false, 0
	}
	return true, TAC_CLOCK_MAPPING[TAC&0b11]
}
