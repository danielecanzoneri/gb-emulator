package cpu

import (
	"testing"
)

func TestTimer_UpdateDIV(t *testing.T) {
	mockCPU := mockCPU()
	timer := mockCPU.timer
	timer.divCounter = DIV_FREQ - 1

	writeTestProgram(mockCPU, NOP_OPCODE)
	mockCPU.ExecuteInstruction()

	if mockCPU.Mem.Read(DIV_ADDR) != 1 {
		t.Errorf("DIV not incremented")
	}

	if timer.divCounter != 0 {
		t.Errorf("divCounter not reset")
	}
}

func TestTimer_UpdateTIMA_NoOverflow(t *testing.T) {
	mockCPU := mockCPU()
	timer := mockCPU.timer

	var TIMA, TMA uint8 = 0xFE, 0x10
	mockCPU.Mem.Write(IF_ADDR, 0)
	mockCPU.Mem.Write(TIMA_ADDR, TIMA)
	mockCPU.Mem.Write(TMA_ADDR, TMA)
	mockCPU.Mem.Write(TAC_ADDR, 0b101)

	writeTestProgram(mockCPU, NOP_OPCODE, NOP_OPCODE, NOP_OPCODE, NOP_OPCODE)
	mockCPU.ExecuteInstruction()
	mockCPU.ExecuteInstruction()
	mockCPU.ExecuteInstruction()

	if mockCPU.Mem.Read(TIMA_ADDR) != TIMA {
		t.Errorf("TIMA: got %02X, expected %02X", mockCPU.Mem.Read(TIMA_ADDR), TIMA)
	}

	mockCPU.ExecuteInstruction()
	if mockCPU.Mem.Read(TIMA_ADDR) != TIMA+1 {
		t.Errorf("failed to increment TIMA")
	}

	if timer.timaCounter != 0 {
		t.Errorf("timaCounter not reset to 0")
	}

	if mockCPU.Mem.Read(IF_ADDR)&TIMER_INT_MASK != 0 {
		t.Errorf("interrupt should have not been requested")
	}
}

func TestTimer_UpdateTIMA_Overflow(t *testing.T) {
	mockCPU := mockCPU()
	timer := mockCPU.timer

	var TIMA, TMA uint8 = 0xFF, 0x10
	mockCPU.Mem.Write(IF_ADDR, 0)
	mockCPU.Mem.Write(TIMA_ADDR, TIMA)
	mockCPU.Mem.Write(TMA_ADDR, TMA)
	mockCPU.Mem.Write(TAC_ADDR, 0b101)

	writeTestProgram(mockCPU, NOP_OPCODE)
	timer.timaCounter = 3

	mockCPU.ExecuteInstruction()

	if mockCPU.Mem.Read(TIMA_ADDR) != TMA {
		t.Errorf("TIMA not reset to TMA: got %02X, expected %02X", mockCPU.Mem.Read(TIMA_ADDR), TMA)
	}

	if mockCPU.Mem.Read(IF_ADDR)&TIMER_INT_MASK == 0 {
		t.Errorf("interrupt not requested")
	}
}
