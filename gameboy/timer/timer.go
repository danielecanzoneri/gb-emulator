package timer

import (
	"github.com/danielecanzoneri/lucky-boy/gameboy/audio"
	"github.com/danielecanzoneri/lucky-boy/util"
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
	apu       *audio.APU
	prevBit12 uint8

	// Flag to request interrupt at next step
	timaOverflow bool
	timaReloaded bool // Don't reload TIMA if is written after overflow

	// In double speed, Timer will tick at same rate of CPU, but APU will tick at normal speed
	speedFactor int // (0: normal, 1: double)

	// Callback to request interrupt
	RequestInterrupt func()

	DIVGlitched func() bool
}

func New(apu *audio.APU) *Timer {
	// It seems that at startup actual Game Boy timer has elapsed for eight ticks,
	// (maybe it's for a wrong boot rom emulation)
	t := &Timer{
		apu:         apu,
		speedFactor: 0,
		DIVGlitched: func() bool { return false },
	}

	t.Tick(8)
	return t
}

func (t *Timer) Tick(ticks int) {
	if t.DIVGlitched() {
		return
	}

	// Calculate how many ticks until next multiple of 4
	ticksToNextMultiple := int(4 - (t.systemCounter % 4))

	// Process ticks in batches
	for {
		if ticks < ticksToNextMultiple {
			t.systemCounter += uint16(ticks)
			return
		}

		// Update system counter for this batch
		t.systemCounter += uint16(ticksToNextMultiple)
		ticks -= ticksToNextMultiple

		// Happens every 4 ticks
		t.timaReloaded = false
		if t.timaOverflow {
			t.timaOverflow = false
			t.timaReloaded = true
			t.RequestInterrupt()
			t.TIMA = t.TMA
		}

		// Update TIMA
		t.detectFallingEdge()

		// Update ticks to next multiple for next iteration
		ticksToNextMultiple = 4
	}
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

	// Update DIV APU
	t.detectAPUFallingEdge()
}

func (t *Timer) detectAPUFallingEdge() {
	// APU frame sequencer runs at 512Hz (regardless of CPU speed)
	// DIV bit 12 toggles at 512Hz in normal speed (4194304/2^13)
	// DIV bit 13 toggles at 512Hz in double speed (8388608/2^14)
	currBit := util.ReadBit(t.systemCounter, 12+uint8(t.speedFactor))

	// Detect falling edge
	if t.prevBit12 == 1 && currBit == 0 {
		t.apu.StepFrameSequencer()
	}
	t.prevBit12 = uint8(currBit)
}

func (t *Timer) SwitchSpeed(doubleSpeed bool) {
	if doubleSpeed {
		t.speedFactor = 1
	} else {
		t.speedFactor = 0
	}
	t.systemCounter = 0
	t.detectFallingEdge()
}
