package timer

import (
	"testing"
)

func TestTimer_UpdateDIV(t *testing.T) {
	timer := new(Timer)
	timer.systemCounter = 0x00FC

	timer.Tick(4)

	if timer.Read(divAddr) != 1 {
		t.Errorf("DIV not incremented")
	}
}

func TestTimer_UpdateTIMA_NoOverflow(t *testing.T) {
	timer := new(Timer)

	interruptSet := false
	interruptRequestFunc := func() {
		interruptSet = true
	}
	timer.RequestInterrupt = interruptRequestFunc

	var TIMA, TMA uint8 = 0xFE, 0x10
	timer.TIMA = TIMA
	timer.TMA = TMA
	timer.Write(tacAddr, 0b101) // enabled and increment every 4 cycles

	timer.Tick(12)
	if timer.TIMA != TIMA {
		t.Errorf("TIMA: got %02X, expected %02X", timer.TIMA, TIMA)
	}

	timer.Tick(4)
	if timer.TIMA != TIMA+1 {
		t.Errorf("failed to increment TIMA")
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
	timer.RequestInterrupt = interruptRequestFunc

	var TIMA, TMA uint8 = 0xFF, 0x10
	timer.TIMA = TIMA
	timer.TMA = TMA
	timer.Write(tacAddr, 0b101) // timaFreq = 4

	timer.Tick(16)

	// Overflow, interrupt requested one cycle later
	if timer.TIMA != 0 {
		t.Errorf("TIMA should stay 0 after overflow: got %02X", timer.TIMA)
	}
	if interruptSet {
		t.Errorf("interrupt requested too early")
	}

	timer.Tick(4)
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
	timer.RequestInterrupt = interruptRequestFunc

	var TIMA, TMA uint8 = 0xFF, 0x10
	timer.TIMA = TIMA
	timer.TMA = TMA
	timer.Write(tacAddr, 0b101) // timaFreq = 4

	var newTIMA uint8 = 0x12

	timer.Tick(16) // TIMA should be 0
	timer.Write(timaAddr, newTIMA)
	timer.Tick(4) // Tima should not have been reset and interrupt not requested

	// TIMA should not be reset and interrupt not requested
	if timer.TIMA != newTIMA {
		t.Errorf("TIMA should not have been reset to TMA")
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
	timer.RequestInterrupt = interruptRequestFunc

	var TIMA, TMA uint8 = 0xFF, 0x10
	timer.TIMA = TIMA
	timer.TMA = TMA
	timer.Write(tacAddr, 0b101) // timaFreq = 4

	var newTIMA uint8 = 0x80

	timer.Tick(20)
	timer.Write(timaAddr, newTIMA)
	timer.Tick(4)

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
	timer.RequestInterrupt = interruptRequestFunc

	var TIMA, TMA uint8 = 0xFF, 0x10
	timer.TIMA = TIMA
	timer.TMA = TMA
	timer.Write(tacAddr, 0b101) // timaFreq = 4

	var newTMA uint8 = 0x20

	timer.Tick(20)
	timer.Write(tmaAddr, newTMA)
	timer.Tick(4)

	// TIMA should be reset to new TMA
	if timer.TIMA != newTMA {
		t.Errorf("TIMA: got %02X, expected %02X", timer.TIMA, newTMA)
	}
	if !interruptSet {
		t.Errorf("interrupt not requested")
	}
}
