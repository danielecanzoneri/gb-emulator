package gameboy

import (
	"github.com/danielecanzoneri/lucky-boy/gameboy/audio"
	"github.com/danielecanzoneri/lucky-boy/gameboy/cartridge"
	"github.com/danielecanzoneri/lucky-boy/gameboy/cpu"
	"github.com/danielecanzoneri/lucky-boy/gameboy/joypad"
	"github.com/danielecanzoneri/lucky-boy/gameboy/mmu"
	"github.com/danielecanzoneri/lucky-boy/gameboy/ppu"
	"github.com/danielecanzoneri/lucky-boy/gameboy/serial"
	"github.com/danielecanzoneri/lucky-boy/gameboy/timer"
	"log"
)

type SystemModel int

const (
	Auto SystemModel = iota
	DMG
	CGB
)

type GameBoy struct {
	CPU        *cpu.CPU
	SerialPort *serial.Port
	Timer      *timer.Timer
	Memory     *mmu.MMU
	PPU        *ppu.PPU
	Joypad     *joypad.Joypad
	APU        *audio.APU

	// DMG or CGB (Auto to automatically detect it based on cartridge)
	Model          SystemModel
	EmulationModel SystemModel // Actual model used to emulate

	sampleRate float64
	sampleBuff chan float32
}

func New(audioSampleBuffer chan float32, sampleRate float64) *GameBoy {
	gb := &GameBoy{
		sampleRate: sampleRate,
		sampleBuff: audioSampleBuffer,
	}

	return gb
}

func (gb *GameBoy) initComponents(rom cartridge.Cartridge) {
	isCGB := gb.EmulationModel == CGB

	gb.PPU = ppu.New(isCGB)
	gb.Joypad = joypad.New()
	gb.APU = audio.New(gb.sampleRate, gb.sampleBuff, isCGB)
	gb.SerialPort = serial.NewPort()
	gb.Timer = timer.New(gb.APU)

	gb.Memory = mmu.New(gb.PPU, gb.APU, gb.Timer, gb.Joypad, gb.SerialPort, isCGB)
	gb.CPU = cpu.New(gb.Memory, gb.PPU, isCGB)
	gb.Memory.IsCPUHalted = gb.CPU.Halted
	gb.Timer.DIVGlitched = gb.CPU.SpeedSwitchHalted
	gb.CPU.AddTicker(gb.SerialPort, gb.Timer, gb.PPU, gb.Memory, gb.APU)

	// Load ROM into memory
	gb.Memory.Cartridge = rom

	// Set interrupts request handler
	gb.PPU.RequestVBlankInterrupt = func() { gb.CPU.RequestInterrupt(cpu.VBlankInterruptMask) }
	gb.PPU.RequestSTATInterrupt = func() { gb.CPU.RequestInterrupt(cpu.STATInterruptMask) }
	gb.Timer.RequestInterrupt = func() { gb.CPU.RequestInterrupt(cpu.TimerInterruptMask) }
	gb.SerialPort.RequestInterrupt = func() { gb.CPU.RequestInterrupt(cpu.SerialInterruptMask) }
	gb.Joypad.RequestInterrupt = func() { gb.CPU.RequestInterrupt(cpu.JoypadInterruptMask) }
}

func (gb *GameBoy) Reset() {
	rom := gb.Memory.Cartridge
	bootRom := gb.Memory.BootRom

	// Load ROM and boot ROM
	gb.Load(rom)
	gb.LoadBootROM(bootRom)
}

func (gb *GameBoy) Load(rom cartridge.Cartridge) {
	// Detect system model
	if gb.Model == Auto {
		if rom.Header().CgbMode == cartridge.DmgOnly {
			gb.EmulationModel = DMG
		} else {
			gb.EmulationModel = CGB
		}
	} else if gb.Model == DMG {
		gb.EmulationModel = DMG

		if rom.Header().CgbMode == cartridge.CgbOnly {
			log.Println("WARNING: DMG doesn't support CGB only games, running as CGB")
		}
	} else {
		gb.EmulationModel = CGB
	}

	gb.initComponents(rom)

	// MBC3 RTC clocking
	if c, ok := rom.(cpu.Ticker); ok {
		gb.CPU.AddTicker(c)
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
	gb.Memory.DisableBootROM()

	if gb.EmulationModel == DMG {
		gb.CPU.SkipDMGBoot()
		gb.Timer.SkipDMGBoot()
		gb.Memory.SkipBoot()
		gb.PPU.SkipDMGBoot()
		gb.APU.SkipBoot()
	} else {
		// Find compatibility mode palette
		gb.Timer.SkipCGBBoot()
		gb.Memory.SkipBoot()
		titleChecksum := gb.PPU.SkipCGBBoot(gb.Memory.Cartridge)
		gb.CPU.SkipCGBBoot(titleChecksum)
		gb.APU.SkipBoot()
	}
}
