package debugger

import (
	"fmt"
	"strings"
)

func (d *Debugger) updateCPU() string {
	// Update CPU registers

	af := d.cpu.ReadAF()
	bc := d.cpu.ReadBC()
	de := d.cpu.ReadDE()
	hl := d.cpu.ReadHL()
	pc := d.cpu.ReadPC()
	sp := d.cpu.ReadSP()

	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("AF: %04X\n", af))
	builder.WriteString(fmt.Sprintf("BC: %04X\n", bc))
	builder.WriteString(fmt.Sprintf("DE: %04X\n", de))
	builder.WriteString(fmt.Sprintf("HL: %04X\n", hl))
	builder.WriteString(fmt.Sprintf("PC: %04X\n", pc))
	builder.WriteString(fmt.Sprintf("SP: %04X", sp))

	return builder.String()
}

func (d *Debugger) updateInterrupts() string {
	// Update Interrupts panel
	IF := d.mem.DebugRead(0xFF0F)
	IE := d.mem.DebugRead(0xFFFF)
	ime := d.cpu.InterruptEnabled()

	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("FF0F IF: %02X\n", IF))
	builder.WriteString(fmt.Sprintf("FFFF IE: %02X\n", IE))
	builder.WriteString("IME ")
	if ime {
		builder.WriteString("enabled")
	} else {
		builder.WriteString("disabled")
	}

	return builder.String()
}

func (d *Debugger) updateLCD() string {
	// Update LCD panel
	lcdc := d.mem.DebugRead(0xFF40)
	stat := d.mem.DebugRead(0xFF41)
	scy := d.mem.DebugRead(0xFF42)
	scx := d.mem.DebugRead(0xFF43)
	ly := d.mem.DebugRead(0xFF44)
	lyc := d.mem.DebugRead(0xFF45)
	dma := d.mem.DebugRead(0xFF46)
	bgp := d.mem.DebugRead(0xFF47)
	obp0 := d.mem.DebugRead(0xFF48)
	obp1 := d.mem.DebugRead(0xFF49)
	wy := d.mem.DebugRead(0xFF4A)
	wx := d.mem.DebugRead(0xFF4B)

	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("FF40 LCDC: %02X\n", lcdc))
	builder.WriteString(fmt.Sprintf("FF41 STAT: %02X\n", stat))
	builder.WriteString(fmt.Sprintf("FF42 SCY:  %02X\n", scy))
	builder.WriteString(fmt.Sprintf("FF43 SCX:  %02X\n", scx))
	builder.WriteString(fmt.Sprintf("FF44 LY:   %02X\n", ly))
	builder.WriteString(fmt.Sprintf("FF45 LYC:  %02X\n", lyc))
	builder.WriteString(fmt.Sprintf("FF46 DMA:  %02X\n", dma))
	builder.WriteString(fmt.Sprintf("FF47 BGP:  %02X\n", bgp))
	builder.WriteString(fmt.Sprintf("FF48 OBP0: %02X\n", obp0))
	builder.WriteString(fmt.Sprintf("FF49 OBP1: %02X\n", obp1))
	builder.WriteString(fmt.Sprintf("FF4A WY:   %02X\n", wy))
	builder.WriteString(fmt.Sprintf("FF4B WX:   %02X", wx))

	return builder.String()
}

func (d *Debugger) updateCh1() string {
	// Update Sound Channel 1 panel
	nr10 := d.mem.DebugRead(0xFF10)
	nr11 := d.mem.DebugRead(0xFF11)
	nr12 := d.mem.DebugRead(0xFF12)
	nr13 := d.mem.DebugRead(0xFF13)
	nr14 := d.mem.DebugRead(0xFF14)

	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("FF10 NR10: %02X\n", nr10))
	builder.WriteString(fmt.Sprintf("FF11 NR11: %02X\n", nr11))
	builder.WriteString(fmt.Sprintf("FF12 NR12: %02X\n", nr12))
	builder.WriteString(fmt.Sprintf("FF13 NR13: %02X\n", nr13))
	builder.WriteString(fmt.Sprintf("FF14 NR14: %02X", nr14))

	return builder.String()
}

func (d *Debugger) updateCh2() string {
	// Update Sound Channel 2 panel
	nr21 := d.mem.DebugRead(0xFF16)
	nr22 := d.mem.DebugRead(0xFF17)
	nr23 := d.mem.DebugRead(0xFF18)
	nr24 := d.mem.DebugRead(0xFF19)

	var builder strings.Builder
	builder.WriteString("FF15 NR20: --\n")
	builder.WriteString(fmt.Sprintf("FF16 NR21: %02X\n", nr21))
	builder.WriteString(fmt.Sprintf("FF17 NR22: %02X\n", nr22))
	builder.WriteString(fmt.Sprintf("FF18 NR23: %02X\n", nr23))
	builder.WriteString(fmt.Sprintf("FF19 NR24: %02X", nr24))

	return builder.String()
}

func (d *Debugger) updateCh3() string {
	// Update Sound Channel 3 panel
	nr30 := d.mem.DebugRead(0xFF1A)
	nr31 := d.mem.DebugRead(0xFF1B)
	nr32 := d.mem.DebugRead(0xFF1C)
	nr33 := d.mem.DebugRead(0xFF1D)
	nr34 := d.mem.DebugRead(0xFF1E)

	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("FF1A NR30: %02X\n", nr30))
	builder.WriteString(fmt.Sprintf("FF1B NR31: %02X\n", nr31))
	builder.WriteString(fmt.Sprintf("FF1C NR32: %02X\n", nr32))
	builder.WriteString(fmt.Sprintf("FF1D NR33: %02X\n", nr33))
	builder.WriteString(fmt.Sprintf("FF1E NR34: %02X", nr34))

	return builder.String()
}

func (d *Debugger) updateCh4() string {
	// Update Sound Channel 4 panel
	nr41 := d.mem.DebugRead(0xFF20)
	nr42 := d.mem.DebugRead(0xFF21)
	nr43 := d.mem.DebugRead(0xFF22)
	nr44 := d.mem.DebugRead(0xFF23)

	var builder strings.Builder
	builder.WriteString("FF1F NR40: --\n")
	builder.WriteString(fmt.Sprintf("FF20 NR41: %02X\n", nr41))
	builder.WriteString(fmt.Sprintf("FF21 NR42: %02X\n", nr42))
	builder.WriteString(fmt.Sprintf("FF22 NR43: %02X\n", nr43))
	builder.WriteString(fmt.Sprintf("FF23 NR44: %02X", nr44))

	return builder.String()
}

func (d *Debugger) updateSoundControl() string {
	// Update Sound Control panel
	nr50 := d.mem.DebugRead(0xFF24)
	nr51 := d.mem.DebugRead(0xFF25)
	nr52 := d.mem.DebugRead(0xFF26)

	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("FF24 NR50: %02X\n", nr50))
	builder.WriteString(fmt.Sprintf("FF25 NR51: %02X\n", nr51))
	builder.WriteString(fmt.Sprintf("FF26 NR52: %02X", nr52))

	return builder.String()
}

func (d *Debugger) updateWaveRam() string {
	// Update WaveRam panel
	var builder strings.Builder
	for i := uint16(0); i < 16; i++ {
		b := d.mem.DebugRead(0xFF30 + i)
		builder.WriteString(fmt.Sprintf("%02X ", b))
		if i%4 == 3 && i < 15 {
			builder.WriteString("\n")
		}
	}

	return builder.String()
}

func (d *Debugger) updateOpcode() string {
	// Update Opcode panel
	pc := d.cpu.ReadPC()
	opcode := d.mem.DebugRead(pc)

	var opcodeHex, opcodeName string
	if opcode == 0xCB { // Prefix
		suffix := d.mem.DebugRead(pc + 1)
		opcodeHex = fmt.Sprintf("0xCB %02X", suffix)
		opcodeName = prefixedOpcodes[suffix]
	} else {
		opcodeHex = fmt.Sprintf("0x%02X", opcode)
		opcodeName = opcodes[opcode]
	}

	return opcodeHex + "\n" + opcodeName
}

func (d *Debugger) updateTimer() string {
	// Update Timer panel
	div := d.mem.DebugRead(0xFF04)
	tima := d.mem.DebugRead(0xFF05)
	tma := d.mem.DebugRead(0xFF06)
	tac := d.mem.DebugRead(0xFF07)

	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("FF04 DIV:  %02X\n", div))
	builder.WriteString(fmt.Sprintf("FF05 TIMA: %02X\n", tima))
	builder.WriteString(fmt.Sprintf("FF06 TMA:  %02X\n", tma))
	builder.WriteString(fmt.Sprintf("FF07 TAC:  %02X", tac))

	return builder.String()
}
