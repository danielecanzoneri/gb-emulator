package cpu

import (
	"github.com/danielecanzoneri/gb-emulator/gameboy/cartridge"
	"github.com/danielecanzoneri/gb-emulator/gameboy/joypad"
	"github.com/danielecanzoneri/gb-emulator/gameboy/memory"
	"github.com/danielecanzoneri/gb-emulator/gameboy/ppu"
	"github.com/danielecanzoneri/gb-emulator/gameboy/timer"
)

func mockCPU() *CPU {
	p := ppu.New()
	c := &cartridge.MBC1{ROM: make([]uint8, 0x8000), RAM: make([]uint8, 0x2000), RAMBanks: 1, ROMBanks: 1}
	mem := &memory.MMU{PPU: p, Cartridge: c, Joypad: &joypad.Joypad{}, Timer: &timer.Timer{}, BootRomDisabled: true}
	mem.Write(0, 0x0A) // Enable RAM
	return &CPU{SP: 0xFFFE, MMU: mem}
}

func writeTestProgram(cpu *CPU, data ...byte) {
	for i, b := range data {
		addr := uint16(i) + cpu.PC
		if addr < 0x8000 {
			cpu.MMU.Cartridge.(*cartridge.MBC1).ROM[addr] = b
		} else {
			cpu.MMU.Write(addr, b)
		}
	}
}
