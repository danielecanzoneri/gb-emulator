package gameboy

import (
	"github.com/danielecanzoneri/gb-emulator/gameboy/audio"
	"github.com/danielecanzoneri/gb-emulator/gameboy/cartridge"
	"github.com/danielecanzoneri/gb-emulator/gameboy/cpu"
	"github.com/danielecanzoneri/gb-emulator/gameboy/joypad"
	"github.com/danielecanzoneri/gb-emulator/gameboy/memory"
	"github.com/danielecanzoneri/gb-emulator/gameboy/ppu"
	"github.com/danielecanzoneri/gb-emulator/gameboy/serial"
	"github.com/danielecanzoneri/gb-emulator/gameboy/timer"
)

type GameBoy struct {
	CPU        *cpu.CPU
	SerialPort *serial.Port
	Timer      *timer.Timer
	Memory     *memory.MMU
	PPU        *ppu.PPU
	Joypad     *joypad.Joypad
	APU        *audio.APU

	sampleRate float64
	sampleBuff chan float32
}

func New(audioSampleBuffer chan float32, sampleRate float64) *GameBoy {
	gb := &GameBoy{
		sampleRate: sampleRate,
		sampleBuff: audioSampleBuffer,
	}

	gb.initComponents()

	return gb
}

func (gb *GameBoy) initComponents() {
	gb.PPU = ppu.New()
	gb.Joypad = joypad.New()
	gb.APU = audio.NewAPU(gb.sampleRate, gb.sampleBuff)
	gb.SerialPort = serial.NewPort()
	gb.Timer = timer.New()
	gb.Timer.APU = gb.APU

	gb.Memory = &memory.MMU{Serial: gb.SerialPort, Timer: gb.Timer, PPU: gb.PPU, Joypad: gb.Joypad, APU: gb.APU}
	gb.CPU = cpu.New()
	gb.CPU.Timer = gb.Timer
	gb.CPU.MMU = gb.Memory
	gb.CPU.PPU = gb.PPU
	gb.CPU.AddCycler(gb.SerialPort, gb.Timer, gb.PPU, gb.Memory, gb.APU)

	// Set interrupt request for timer
	gb.Timer.RequestInterrupt = cpu.RequestTimerInterruptFunc(gb.CPU)
	// Set interrupt request for PPU
	gb.PPU.RequestVBlankInterrupt = cpu.RequestVBlankInterruptFunc(gb.CPU)
	gb.PPU.RequestSTATInterrupt = cpu.RequestSTATInterruptFunc(gb.CPU)
	// Set interrupt request for serial
	gb.SerialPort.RequestInterrupt = cpu.RequestSerialInterruptFunc(gb.CPU)
}

func (gb *GameBoy) Reset() {
	rom := gb.Memory.Cartridge
	bootRom := gb.Memory.BootRom

	// Reset all components
	gb.initComponents()

	// Load ROM and boot ROM
	gb.Load(rom)
	gb.LoadBootROM(bootRom)
}

func (gb *GameBoy) Load(rom cartridge.Cartridge) {
	// Load ROM into memory
	gb.Memory.Cartridge = rom

	// MBC3 RTC clocking
	if c, ok := rom.(interface{ Tick(uint) }); ok {
		gb.CPU.AddCycler(c)
	}
}

func (gb *GameBoy) LoadBootROM(bootRom []uint8) {
	if bootRom == nil {
		gb.skipBootROM()
		return
	}

	gb.Memory.BootRomDisabled = false
	gb.Memory.BootRom = bootRom
}

func (gb *GameBoy) skipBootROM() {
	gb.Memory.BootRomDisabled = true

	gb.CPU.SkipBoot()
	gb.Timer.SkipBoot()
	gb.Memory.SkipBoot()
	gb.PPU.SkipBoot()
	gb.APU.SkipBoot()
}
