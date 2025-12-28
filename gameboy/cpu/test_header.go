package cpu

import (
	"github.com/danielecanzoneri/lucky-boy/gameboy/audio"
	"github.com/danielecanzoneri/lucky-boy/gameboy/cartridge"
	"github.com/danielecanzoneri/lucky-boy/gameboy/joypad"
	"github.com/danielecanzoneri/lucky-boy/gameboy/mmu"
	"github.com/danielecanzoneri/lucky-boy/gameboy/ppu"
	"github.com/danielecanzoneri/lucky-boy/gameboy/serial"
	"github.com/danielecanzoneri/lucky-boy/gameboy/timer"
)

func mockCPU() *CPU {
	p := ppu.New(false)
	a := audio.New(48000, make(chan float32, 10), false)
	c := &cartridge.MBC1{ROM: make([]uint8, 0x8000), RAM: make([]uint8, 0x2000), RAMBanks: 1, ROMBanks: 1}
	mem := mmu.New(p, a, timer.New(a), joypad.New(), serial.NewPort(), false)
	mem.Cartridge = c
	mem.BootRomDisabled = true
	mem.Write(0, 0x0A) // Enable RAM

	cpu := New(mem, p, false)
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
