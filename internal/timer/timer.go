package timer

import "fmt"

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

func (t *Timer) String() string {
	return fmt.Sprintf(
		"DIV:%02X, TIMA:%02X, TMA:%02X, TAC:%02X, DIV counter:%d, TIMA counter:%d",
		t.DIV, t.TIMA, t.TMA, t.TAC, t.divCounter, t.timaCounter,
	)
}

func (t *Timer) Cycle() {
	// Update DIV
	t.divCounter++
	t.updateDIV()

	// Update TIMA if enabled
	if t.enabled {
		t.timaCounter++
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
