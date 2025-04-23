package timer

const (
	divFreq = 64
)

type Timer struct {
	DIV  uint8
	TIMA uint8
	TMA  uint8
	TAC  uint8

	divCounter  uint
	timaCounter uint

	// Frequency to update TIMA (lower 2 bits of TAC)
	timaFreq uint
	enabled  bool

	// Flag to request interrupt at next step
	timaOverflow bool

	// Callback to request interrupt
	RequestInterrupt func()
}

func (t *Timer) Step(cycles uint) {
	// Update DIV
	t.divCounter += cycles
	t.updateDIV()

	// Update TIMA if enabled
	if t.enabled {
		t.timaCounter += cycles
		t.updateTIMA()
	}
}

func (t *Timer) updateDIV() {
	if t.divCounter >= divFreq {
		t.DIV++
		t.divCounter -= divFreq
	}
}

func (t *Timer) handleTIMAOverflow() {
	if t.timaOverflow {
		t.timaOverflow = false
		t.RequestInterrupt()
		t.TIMA = t.TMA
	}
}

func (t *Timer) updateTIMA() {
	t.handleTIMAOverflow()

	for t.timaCounter >= t.timaFreq {
		t.TIMA++
		t.timaCounter -= t.timaFreq

		// Check overflow
		if t.TIMA == 0 {
			t.timaOverflow = true
		}

		if t.timaCounter > 0 {
			t.handleTIMAOverflow()
		}
	}
}
