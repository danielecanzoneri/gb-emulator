package gameboy

import (
	"fmt"
	"github.com/danielecanzoneri/gb-emulator/emulator/internal/gameboy/audio"
	"github.com/danielecanzoneri/gb-emulator/emulator/internal/gameboy/cartridge"
	"github.com/danielecanzoneri/gb-emulator/emulator/internal/gameboy/cpu"
	"github.com/danielecanzoneri/gb-emulator/emulator/internal/gameboy/joypad"
	"github.com/danielecanzoneri/gb-emulator/emulator/internal/gameboy/memory"
	"github.com/danielecanzoneri/gb-emulator/emulator/internal/gameboy/ppu"
	"github.com/danielecanzoneri/gb-emulator/emulator/internal/gameboy/timer"
)

type GameBoy struct {
	CPU    *cpu.CPU
	Timer  *timer.Timer
	Memory *memory.MMU
	PPU    *ppu.PPU
	Joypad *joypad.Joypad
	APU    *audio.APU

	sampleRate float64
	sampleBuff chan float32

	cycles uint
}

func (gb *GameBoy) Cycle() {
	gb.cycles++
}

func New(audioSampleBuffer chan float32, sampleRate float64) *GameBoy {
	gb := &GameBoy{
		sampleRate: sampleRate,
		sampleBuff: audioSampleBuffer,
	}
	gb.Reset()

	return gb
}

func (gb *GameBoy) Reset() {
	gb.PPU = ppu.New()
	gb.Joypad = joypad.New()
	gb.APU = audio.NewAPU(gb.sampleRate, gb.sampleBuff)
	gb.Timer = &timer.Timer{APU: gb.APU}
	gb.Memory = &memory.MMU{Timer: gb.Timer, PPU: gb.PPU, Joypad: gb.Joypad, APU: gb.APU}
	gb.CPU = &cpu.CPU{Timer: gb.Timer, MMU: gb.Memory}
	gb.CPU.AddCycler(gb.Timer, gb.PPU, gb.Memory, gb.APU)

	// Set interrupt request for timer
	gb.Timer.RequestInterrupt = cpu.RequestTimerInterruptFunc(gb.CPU)
	// Set interrupt request for PPU
	gb.PPU.RequestVBlankInterrupt = cpu.RequestVBlankInterruptFunc(gb.CPU)
	gb.PPU.RequestSTATInterrupt = cpu.RequestSTATInterruptFunc(gb.CPU)

	gb.CPU.AddCycler(gb)

	gb.CPU.Reset()
	gb.Memory.Reset()
}

func (gb *GameBoy) Load(romPath string) (string, error) {
	rom, err := cartridge.LoadROM(romPath)
	if err != nil {
		return "", fmt.Errorf("error loading the cartridge: %v", err)
	}

	// Load ROM into memory
	gb.Memory.CartridgeData = rom.Data
	gb.Memory.SetMBC(rom.Header)

	return rom.Header.Title, nil
}
