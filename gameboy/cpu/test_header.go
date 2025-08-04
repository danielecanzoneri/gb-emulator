package cpu

import (
	"github.com/danielecanzoneri/gb-emulator/gameboy/audio"
	"github.com/danielecanzoneri/gb-emulator/gameboy/cartridge"
	"github.com/danielecanzoneri/gb-emulator/gameboy/joypad"
	"github.com/danielecanzoneri/gb-emulator/gameboy/mmu"
	"github.com/danielecanzoneri/gb-emulator/gameboy/ppu"
	"github.com/danielecanzoneri/gb-emulator/gameboy/serial"
	"github.com/danielecanzoneri/gb-emulator/gameboy/timer"
)

func mockCPU() *CPU {
	p := ppu.New()
	c := &cartridge.MBC1{ROM: make([]uint8, 0x8000), RAM: make([]uint8, 0x2000), RAMBanks: 1, ROMBanks: 1}
	mem := mmu.New(p, audio.NewAPU(48000, make(chan float32, 10)), timer.New(), joypad.New(), serial.NewPort())
	mem.Cartridge = c
	mem.BootRomDisabled = true
	mem.Write(0, 0x0A) // Enable RAM

	cpu := New(mem, p)
	cpu.SP = 0xFFFE
	return cpu
}

func writeTestProgram(cpu *CPU, data ...byte) {
	for i, b := range data {
		addr := uint16(i) + cpu.PC
		if addr < 0x8000 {
			cpu.mmu.Cartridge.(*cartridge.MBC1).ROM[addr] = b
		} else {
			cpu.mmu.Write(addr, b)
		}
	}
}
