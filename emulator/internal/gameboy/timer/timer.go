package timer

import (
	"github.com/danielecanzoneri/gb-emulator/emulator/internal/gameboy/audio"
	"github.com/danielecanzoneri/gb-emulator/pkg/util"
)

type Timer struct {
	// DIV  uint8
	TIMA uint8
	TMA  uint8
	TAC  uint8

	// Updated at every M-cycle (high part is DIV)
	systemCounter uint16

	// Falling edge detector to detect when to update TIMA
	prevState uint8

	// Falling edge detector to detect when to update APU counter
	APU       *audio.APU
	prevBit12 uint8

	// Flag to request interrupt at next step
	timaOverflow bool
	timaReloaded bool // Don't reload TIMA if is written after overflow

	// Callback to request interrupt
	RequestInterrupt func()
}

func (t *Timer) Reset() {
	t.TIMA = 0
	t.TMA = 0
	t.TAC = 0
	t.systemCounter = 0
	t.prevState = 0
	t.prevBit12 = 0
	t.timaOverflow = false
	t.timaReloaded = false
}

func (t *Timer) Cycle() {
	t.timaReloaded = false
	if t.timaOverflow {
		t.timaOverflow = false
		t.timaReloaded = true
		t.RequestInterrupt()
		t.TIMA = t.TMA
	}

	// Update DIV
	t.systemCounter += 4

	// Update TIMA
	t.detectFallingEdge()

	// Update DIV APU
	t.detectAPUFallingEdge()
}

func (t *Timer) detectFallingEdge() {
	var currBit uint16
	switch t.TAC & 0b11 {
	case 0: // 00 = bit 9
		currBit = (t.systemCounter >> 9) & 1
	case 1: // 01 = bit 3
		currBit = (t.systemCounter >> 3) & 1
	case 2: // 10 = bit 5
		currBit = (t.systemCounter >> 5) & 1
	case 3: // 11 = bit 7
		currBit = (t.systemCounter >> 7) & 1
	}

	// And with timer enable bit
	currState := uint8(currBit) & (t.TAC >> 2)
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

func (t *Timer) detectAPUFallingEdge() {
	// APU frame sequencer runs at 512Hz
	// DIV bit 12 toggles at 512Hz (4194304/2^13)
	currBit := util.ReadBit(t.systemCounter, 12)

	// Detect falling edge
	if t.prevBit12 == 1 && currBit == 0 {
		t.APU.StepCounter()
	}
	t.prevBit12 = uint8(currBit)
}
