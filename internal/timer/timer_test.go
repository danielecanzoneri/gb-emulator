package timer

import (
	"testing"
)

func TestTimer_UpdateDIV(t *testing.T) {
	timer := new(Timer)
	timer.divCounter = divFreq - 1

	timer.Step(1)

	if timer.DIV != 1 {
		t.Errorf("DIV not incremented")
	}

	if timer.divCounter != 0 {
		t.Errorf("divCounter not reset")
	}
}

func TestTimer_UpdateTIMA_NoOverflow(t *testing.T) {
	timer := new(Timer)

	interruptSet := false
	interruptRequestFunc := func() {
		interruptSet = true
	}
	timer.SetInterruptRequestFunc(interruptRequestFunc)

	var TIMA, TMA uint8 = 0xFE, 0x10
	timer.TIMA = TIMA
	timer.TMA = TMA
	timer.Write(tacAddr, 0b101) // timaFreq = 4

	timer.Step(3)
	if timer.TIMA != TIMA {
		t.Errorf("TIMA: got %02X, expected %02X", timer.TIMA, TIMA)
	}

	timer.Step(1)
	if timer.TIMA != TIMA+1 {
		t.Errorf("failed to increment TIMA")
	}
	if timer.timaCounter != 0 {
		t.Errorf("timaCounter not reset to 0")
	}
	if interruptSet {
		t.Errorf("interrupt should have not been requested")
	}
}

func TestTimer_UpdateTIMA_Overflow(t *testing.T) {
	timer := new(Timer)

	interruptSet := false
	interruptRequestFunc := func() {
		interruptSet = true
	}
	timer.SetInterruptRequestFunc(interruptRequestFunc)

	var TIMA, TMA uint8 = 0xFF, 0x10
	timer.TIMA = TIMA
	timer.TMA = TMA
	timer.Write(tacAddr, 0b101) // timaFreq = 4

	timer.Step(4)

	// Overflow, interrupt requested one cycle later
	if timer.TIMA != 0 {
		t.Errorf("TIMA should stay 0 after overflow: got %02X", timer.TIMA)
	}
	if interruptSet {
		t.Errorf("interrupt requested too early")
	}

	timer.Step(1)
	// Request interrupt
	if timer.TIMA != TMA {
		t.Errorf("TIMA not reset to TMA: got %02X, expected %02X", timer.TIMA, TMA)
	}
	if !interruptSet {
		t.Errorf("interrupt not requested")
	}
}

func TestTimer_WriteTIMAWhenOverflow(t *testing.T) {
	timer := new(Timer)

	interruptSet := false
	interruptRequestFunc := func() {
		interruptSet = true
	}
	timer.SetInterruptRequestFunc(interruptRequestFunc)

	var TIMA, TMA uint8 = 0xFF, 0x10
	timer.TIMA = TIMA
	timer.TMA = TMA
	timer.Write(tacAddr, 0b101) // timaFreq = 4

	var newTIMA uint8 = 0x12

	timer.Step(3)
	timer.Write(timaAddr, newTIMA)
	timer.Step(2)

	// TIMA should not be reset and interrupt not requested
	if timer.TIMA != newTIMA {
		t.Errorf("TIMA: got %02X, expected %02X", timer.TIMA, newTIMA)
	}
	if interruptSet {
		t.Errorf("interrupt should have not been requested")
	}
}

func TestTimer_WriteTIMAAfterOverflow(t *testing.T) {
	timer := new(Timer)

	interruptSet := false
	interruptRequestFunc := func() {
		interruptSet = true
	}
	timer.SetInterruptRequestFunc(interruptRequestFunc)

	var TIMA, TMA uint8 = 0xFF, 0x10
	timer.TIMA = TIMA
	timer.TMA = TMA
	timer.Write(tacAddr, 0b101) // timaFreq = 4

	var newTIMA uint8 = 0x12

	timer.Step(4)
	timer.Write(timaAddr, newTIMA)
	timer.Step(1)

	// TIMA should be reset to TMA
	if timer.TIMA != TMA {
		t.Errorf("TIMA: got %02X, expected %02X", timer.TIMA, TMA)
	}
	if !interruptSet {
		t.Errorf("interrupt not requested")
	}
}

func TestTimer_WriteTMAAfterOverflow(t *testing.T) {
	timer := new(Timer)

	interruptSet := false
	interruptRequestFunc := func() {
		interruptSet = true
	}
	timer.SetInterruptRequestFunc(interruptRequestFunc)

	var TIMA, TMA uint8 = 0xFF, 0x10
	timer.TIMA = TIMA
	timer.TMA = TMA
	timer.Write(tacAddr, 0b101) // timaFreq = 4

	var newTMA uint8 = 0x20

	timer.Step(4)
	timer.Write(tmaAddr, newTMA)
	timer.Step(1)

	// TIMA should be reset to new TMA
	if timer.TIMA != newTMA {
		t.Errorf("TIMA: got %02X, expected %02X", timer.TIMA, newTMA)
	}
	if !interruptSet {
		t.Errorf("interrupt not requested")
	}
}
