package gameboy

import (
	"github.com/danielecanzoneri/gb-emulator/gameboy/audio"
	"github.com/danielecanzoneri/gb-emulator/gameboy/cartridge"
	"github.com/danielecanzoneri/gb-emulator/gameboy/cpu"
	"github.com/danielecanzoneri/gb-emulator/gameboy/joypad"
	"github.com/danielecanzoneri/gb-emulator/gameboy/mmu"
	"github.com/danielecanzoneri/gb-emulator/gameboy/ppu"
	"github.com/danielecanzoneri/gb-emulator/gameboy/serial"
	"github.com/danielecanzoneri/gb-emulator/gameboy/timer"
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
	gb.APU = audio.NewAPU(gb.sampleRate, gb.sampleBuff)
	gb.SerialPort = serial.NewPort()
	gb.Timer = timer.New(gb.APU)

	gb.Memory = mmu.New(gb.PPU, gb.APU, gb.Timer, gb.Joypad, gb.SerialPort, isCGB)
	gb.CPU = cpu.New(gb.Memory, gb.PPU)
	gb.CPU.AddCycler(gb.SerialPort, gb.Timer, gb.PPU, gb.Memory, gb.APU)

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
	gb.Memory.DisableBootROM()

	gb.CPU.SkipBoot()
	gb.Timer.SkipBoot()
	gb.Memory.SkipBoot()
	gb.PPU.SkipBoot()
	gb.APU.SkipBoot()
}
