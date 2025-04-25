package gameboy

import (
	"github.com/danielecanzoneri/gb-emulator/internal/cpu"
	"github.com/danielecanzoneri/gb-emulator/internal/memory"
	"github.com/danielecanzoneri/gb-emulator/internal/ppu"
	"github.com/danielecanzoneri/gb-emulator/internal/timer"
)

type GameBoy struct {
	CPU    *cpu.CPU
	Timer  *timer.Timer
	Memory *memory.MMU
	PPU    *ppu.PPU

	cycles int
}

func Init() *GameBoy {
	t := &timer.Timer{}
	m := &memory.MMU{Timer: t}
	c := &cpu.CPU{Timer: t, MMU: m}

	// Set interrupt request for timer
	t.RequestInterrupt = cpu.RequestTimerInterruptFunc(c)

	return &GameBoy{CPU: c, Timer: t, Memory: m}
}

func (gb *GameBoy) Reset() {
	gb.CPU.A = 0x01
	gb.CPU.F = 0xB0
	gb.CPU.B = 0x00
	gb.CPU.C = 0x13
	gb.CPU.D = 0x00
	gb.CPU.E = 0xD8
	gb.CPU.H = 0x01
	gb.CPU.L = 0x4D
	gb.CPU.SP = 0xFFFE
	gb.CPU.PC = 0x0100
}

func (gb *GameBoy) Load(romData []byte) {
	// Load ROM into memory
	for i, b := range romData {
		gb.Memory.Write(uint16(i), b)
	}
}

func (gb *GameBoy) Run() {
	// Simplified loop
	for {
		cycles := gb.CPU.ExecuteInstruction()
		gb.Timer.Step(cycles)
		gb.PPU.Step(cycles)
	}
}
