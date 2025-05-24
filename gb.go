package main

import (
	"github.com/danielecanzoneri/gb-emulator/debugger"
	"github.com/danielecanzoneri/gb-emulator/internal/audio"
	"github.com/danielecanzoneri/gb-emulator/internal/cartridge"
	"github.com/danielecanzoneri/gb-emulator/internal/cpu"
	"github.com/danielecanzoneri/gb-emulator/internal/joypad"
	"github.com/danielecanzoneri/gb-emulator/internal/memory"
	"github.com/danielecanzoneri/gb-emulator/internal/ppu"
	"github.com/danielecanzoneri/gb-emulator/internal/timer"
	"github.com/ebitengine/oto/v3"
	"github.com/hajimehoshi/ebiten/v2"
)

type GameBoy struct {
	CPU    *cpu.CPU
	Timer  *timer.Timer
	Memory *memory.MMU
	PPU    *ppu.PPU
	Joypad *joypad.Joypad
	APU    *audio.APU

	debugger *debugger.Debugger

	cycles uint
	paused bool

	gameTitle string

	debugString      string
	debugStringTimer uint
}

func (gb *GameBoy) Cycle() {
	gb.cycles++
}

func Init() (*GameBoy, *oto.Player) {
	sampleBuffer := make(chan float32, bufferSize)

	p := ppu.New()
	j := joypad.New()
	a := audio.NewAPU(sampleRate, sampleBuffer)
	t := &timer.Timer{APU: a}
	m := &memory.MMU{Timer: t, PPU: p, Joypad: j, APU: a}
	c := &cpu.CPU{Timer: t, MMU: m}
	c.AddCycler(t, p, m, a)

	// Set interrupt request for timer
	t.RequestInterrupt = cpu.RequestTimerInterruptFunc(c)
	// Set interrupt request for PPU
	p.RequestVBlankInterrupt = cpu.RequestVBlankInterruptFunc(c)
	p.RequestSTATInterrupt = cpu.RequestSTATInterruptFunc(c)

	gb := &GameBoy{
		CPU: c, Timer: t, PPU: p, Memory: m, Joypad: j,
		APU: a,
	}
	c.AddCycler(gb)

	// Create Debugger
	gb.debugger = debugger.NewDebugger(gb.Memory, gb.CPU)

	player, err := newAudioPlayer(gb, sampleBuffer)
	if err != nil {
		panic(err)
	}

	return gb, player
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

	// gb.Timer.DIV = 0x1E
	gb.Timer.TAC = 0xF8

	gb.Memory.Write(0xFF0F, 0xE1) // IF

	gb.PPU.Write(0xFF40, 0x91) // LCDC
	gb.PPU.Write(0xFF41, 0x81)
	gb.PPU.Write(0xFF47, 0xFC) // BGP
	gb.PPU.Write(0xFF48, 0x00) // OBP0
	gb.PPU.Write(0xFF49, 0x00) // OBP1
}

func (gb *GameBoy) Load(rom *cartridge.Rom) {
	// Load ROM into memory
	gb.Memory.CartridgeData = rom.Data
	gb.Memory.SetMBC(rom.Header)

	gb.gameTitle = rom.Header.Title
	ebiten.SetWindowTitle(gb.gameTitle)
}

func (gb *GameBoy) Pause() {
	gb.paused = !gb.paused

	if gb.paused {
		ebiten.SetWindowTitle(gb.gameTitle + " (paused)")
	} else {
		ebiten.SetWindowTitle(gb.gameTitle)
	}
}

// ToggleDebugger enables/disables visualization of I/O registers
func (gb *GameBoy) ToggleDebugger() {
	gb.debugger.ToggleVisibility()
}
