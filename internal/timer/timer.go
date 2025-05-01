package timer

import (
	"fmt"
	"github.com/danielecanzoneri/gb-emulator/internal/util"
)

const (
	divFreq = 64
)

type Timer struct {
	// DIV  uint8
	TIMA uint8
	TMA  uint8
	TAC  uint8

	// Updated at every M-cycle (high part is DIV)
	systemCounter uint16

	// Frequency to update TIMA (lower 2 bits of TAC)
	bitFalling uint8
	prevState  uint8 // Falling edge detector to detect when to update
	enabled    bool

	// Flag to request interrupt at next step
	timaOverflow bool

	// Callback to request interrupt
	RequestInterrupt func()
}

func (t *Timer) String() string {
	return fmt.Sprintf(
		"DIV:%02X, TIMA:%02X, TMA:%02X, TAC:%02X, counter:%04X",
		t.Read(divAddr), t.TIMA, t.TMA, t.TAC, t.systemCounter,
	)
}

func (t *Timer) Cycle() {
	if t.timaOverflow {
		t.handleTIMAOverflow()
	}

	// Update DIV
	t.systemCounter += 4

	// Update TIMA if enabled
	if t.enabled {
		t.updateTIMA()
	}
}

func (t *Timer) handleTIMAOverflow() {
	t.timaOverflow = false
	t.RequestInterrupt()
	t.TIMA = t.TMA
}

func (t *Timer) updateTIMA() {
	currState := util.ReadBit16(t.systemCounter, t.bitFalling)
	fallingEdge := t.prevState == 1 && currState == 0
	t.prevState = currState

	if fallingEdge {
		t.TIMA++

		// Check overflow
		if t.TIMA == 0 {
			t.timaOverflow = true
		}
	}
}
