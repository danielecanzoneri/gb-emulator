package cpu

import (
	"github.com/danielecanzoneri/gb-emulator/emulator/internal/gameboy/cartridge"
	"github.com/danielecanzoneri/gb-emulator/emulator/internal/gameboy/joypad"
	"github.com/danielecanzoneri/gb-emulator/emulator/internal/gameboy/memory"
	"github.com/danielecanzoneri/gb-emulator/emulator/internal/gameboy/ppu"
	"github.com/danielecanzoneri/gb-emulator/emulator/internal/gameboy/timer"
)

func mockCPU() *CPU {
	p := &ppu.PPU{}
	c := &cartridge.Cartridge{Header: &cartridge.Header{ROMBanks: 1}, Data: make([]byte, 0x8000)}
	c.MBC = cartridge.NewMBC(c.Header)
	mem := &memory.MMU{PPU: p, Cartridge: c, Joypad: &joypad.Joypad{}, Timer: &timer.Timer{}}
	return &CPU{SP: 0xFFFE, MMU: mem}
}

func writeTestProgram(cpu *CPU, data ...byte) {
	for i, b := range data {
		addr := uint16(i) + cpu.PC
		if addr < 0x8000 {
			cpu.MMU.Cartridge.Data[addr] = b
		} else {
			cpu.MMU.Write(addr, b)
		}
	}
}
