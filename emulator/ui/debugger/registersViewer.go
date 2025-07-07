package debugger

import (
	"fmt"
	"github.com/danielecanzoneri/gb-emulator/emulator/internal/gameboy"
	"github.com/ebitenui/ebitenui/widget"
)

type registersViewer struct {
	*widget.Container

	// Sound
	ch1, ch2, ch3, ch4 *panel
	control            *panel
	waveRam            *panel

	// CPU
	cpu        *panel
	interrupts *panel

	lcd   *panel
	timer *panel
}

func newRegisterViewer() *registersViewer {
	rv := &registersViewer{
		ch1:        newSoundCh1Panel(),
		ch2:        newSoundCh2Panel(),
		ch3:        newSoundCh3Panel(),
		ch4:        newSoundCh4Panel(),
		control:    newSoundControlPanel(),
		waveRam:    newWaveRamPanel(),
		cpu:        newCpuPanel(),
		interrupts: newInterruptsPanel(),
		lcd:        newLcdPanel(),
		timer:      newTimerPanel(),
	}

	soundRow := newContainer(widget.DirectionHorizontal,
		rv.ch1, rv.ch2, rv.ch3, rv.ch4,
	)
	cpuColumn := newContainer(widget.DirectionVertical,
		rv.cpu, rv.interrupts,
	)
	timerSoundControlColumn := newContainer(widget.DirectionVertical,
		rv.control, rv.timer,
	)
	centralPanels := newContainer(widget.DirectionVertical,
		newContainer(widget.DirectionHorizontal,
			cpuColumn,               // left
			timerSoundControlColumn, // right
		),
		rv.waveRam,
	)
	secondRow := newContainer(widget.DirectionHorizontal,
		rv.lcd, centralPanels,
	)

	rv.Container = newContainer(widget.DirectionVertical,
		soundRow, secondRow,
	)
	return rv
}

func (rv *registersViewer) Sync(gb *gameboy.GameBoy) {
	rv.ch1.Sync(gb)
	rv.ch2.Sync(gb)
	rv.ch3.Sync(gb)
	rv.ch4.Sync(gb)
	rv.control.Sync(gb)
	rv.waveRam.Sync(gb)
	rv.cpu.Sync(gb)
	rv.interrupts.Sync(gb)
	rv.lcd.Sync(gb)
	rv.timer.Sync(gb)
}

func newCpuPanel() *panel {
	entries := []panelEntry{
		{name: "AF", valueSync: func(gb *gameboy.GameBoy) string { return fmt.Sprintf("%04X", gb.CPU.ReadAF()) }},
		{name: "BC", valueSync: func(gb *gameboy.GameBoy) string { return fmt.Sprintf("%04X", gb.CPU.ReadBC()) }},
		{name: "DE", valueSync: func(gb *gameboy.GameBoy) string { return fmt.Sprintf("%04X", gb.CPU.ReadDE()) }},
		{name: "HL", valueSync: func(gb *gameboy.GameBoy) string { return fmt.Sprintf("%04X", gb.CPU.ReadHL()) }},
		{name: "PC", valueSync: func(gb *gameboy.GameBoy) string { return fmt.Sprintf("%04X", gb.CPU.ReadPC()) }},
		{name: "SP", valueSync: func(gb *gameboy.GameBoy) string { return fmt.Sprintf("%04X", gb.CPU.ReadSP()) }},
	}
	return newPanel("CPU", entries...)
}

func newInterruptsPanel() *panel {
	imeString := map[bool]string{true: "enabled", false: "disabled"}
	entries := []panelEntry{
		{name: "FF0F IF", valueSync: func(gb *gameboy.GameBoy) string { return fmt.Sprintf("%02X", gb.Memory.DebugRead(0xFF0F)) }},
		{name: "FFFF IE", valueSync: func(gb *gameboy.GameBoy) string { return fmt.Sprintf("%02X", gb.Memory.DebugRead(0xFFFF)) }},
		{name: "IME", valueSync: func(gb *gameboy.GameBoy) string { return imeString[gb.CPU.IME] }},
	}
	return newPanel("Interrupts", entries...)
}

func newLcdPanel() *panel {
	entries := []panelEntry{
		{name: "FF40 LCDC", valueSync: func(gb *gameboy.GameBoy) string { return fmt.Sprintf("%02X", gb.Memory.DebugRead(0xFF40)) }},
		{name: "FF41 STAT", valueSync: func(gb *gameboy.GameBoy) string { return fmt.Sprintf("%02X", gb.Memory.DebugRead(0xFF41)) }},
		{name: "FF42 SCY", valueSync: func(gb *gameboy.GameBoy) string { return fmt.Sprintf("%02X", gb.Memory.DebugRead(0xFF42)) }},
		{name: "FF43 SCX", valueSync: func(gb *gameboy.GameBoy) string { return fmt.Sprintf("%02X", gb.Memory.DebugRead(0xFF43)) }},
		{name: "FF44 LY", valueSync: func(gb *gameboy.GameBoy) string { return fmt.Sprintf("%02X", gb.Memory.DebugRead(0xFF44)) }},
		{name: "FF45 LYC", valueSync: func(gb *gameboy.GameBoy) string { return fmt.Sprintf("%02X", gb.Memory.DebugRead(0xFF45)) }},
		{name: "FF46 DMA", valueSync: func(gb *gameboy.GameBoy) string { return fmt.Sprintf("%02X", gb.Memory.DebugRead(0xFF46)) }},
		{name: "FF47 BGP", valueSync: func(gb *gameboy.GameBoy) string { return fmt.Sprintf("%02X", gb.Memory.DebugRead(0xFF47)) }},
		{name: "FF48 OBP0", valueSync: func(gb *gameboy.GameBoy) string { return fmt.Sprintf("%02X", gb.Memory.DebugRead(0xFF48)) }},
		{name: "FF49 OBP1", valueSync: func(gb *gameboy.GameBoy) string { return fmt.Sprintf("%02X", gb.Memory.DebugRead(0xFF49)) }},
		{name: "FF4A WY", valueSync: func(gb *gameboy.GameBoy) string { return fmt.Sprintf("%02X", gb.Memory.DebugRead(0xFF4A)) }},
	}
	return newPanel("LCD", entries...)
}

func newSoundCh1Panel() *panel {
	entries := []panelEntry{
		{name: "FF10 NR10", valueSync: func(gb *gameboy.GameBoy) string { return fmt.Sprintf("%02X", gb.Memory.DebugRead(0xFF10)) }},
		{name: "FF11 NR11", valueSync: func(gb *gameboy.GameBoy) string { return fmt.Sprintf("%02X", gb.Memory.DebugRead(0xFF11)) }},
		{name: "FF12 NR12", valueSync: func(gb *gameboy.GameBoy) string { return fmt.Sprintf("%02X", gb.Memory.DebugRead(0xFF12)) }},
		{name: "FF13 NR13", valueSync: func(gb *gameboy.GameBoy) string { return fmt.Sprintf("%02X", gb.Memory.DebugRead(0xFF13)) }},
		{name: "FF14 NR14", valueSync: func(gb *gameboy.GameBoy) string { return fmt.Sprintf("%02X", gb.Memory.DebugRead(0xFF14)) }},
	}
	return newPanel("Ch1 (Square)", entries...)
}

func newSoundCh2Panel() *panel {
	entries := []panelEntry{
		{name: "FF16 NR21", valueSync: func(gb *gameboy.GameBoy) string { return fmt.Sprintf("%02X", gb.Memory.DebugRead(0xFF16)) }},
		{name: "FF17 NR22", valueSync: func(gb *gameboy.GameBoy) string { return fmt.Sprintf("%02X", gb.Memory.DebugRead(0xFF17)) }},
		{name: "FF18 NR23", valueSync: func(gb *gameboy.GameBoy) string { return fmt.Sprintf("%02X", gb.Memory.DebugRead(0xFF18)) }},
		{name: "FF19 NR24", valueSync: func(gb *gameboy.GameBoy) string { return fmt.Sprintf("%02X", gb.Memory.DebugRead(0xFF19)) }},
	}
	return newPanel("Ch2 (Square)", entries...)
}

func newSoundCh3Panel() *panel {
	entries := []panelEntry{
		{name: "FF1A NR30", valueSync: func(gb *gameboy.GameBoy) string { return fmt.Sprintf("%02X", gb.Memory.DebugRead(0xFF1A)) }},
		{name: "FF1B NR31", valueSync: func(gb *gameboy.GameBoy) string { return fmt.Sprintf("%02X", gb.Memory.DebugRead(0xFF1B)) }},
		{name: "FF1C NR32", valueSync: func(gb *gameboy.GameBoy) string { return fmt.Sprintf("%02X", gb.Memory.DebugRead(0xFF1C)) }},
		{name: "FF1D NR33", valueSync: func(gb *gameboy.GameBoy) string { return fmt.Sprintf("%02X", gb.Memory.DebugRead(0xFF1D)) }},
		{name: "FF1E NR34", valueSync: func(gb *gameboy.GameBoy) string { return fmt.Sprintf("%02X", gb.Memory.DebugRead(0xFF1E)) }},
	}
	return newPanel("Ch3 (Wave)", entries...)
}

func newSoundCh4Panel() *panel {
	entries := []panelEntry{
		{name: "FF20 NR41", valueSync: func(gb *gameboy.GameBoy) string { return fmt.Sprintf("%02X", gb.Memory.DebugRead(0xFF20)) }},
		{name: "FF21 NR42", valueSync: func(gb *gameboy.GameBoy) string { return fmt.Sprintf("%02X", gb.Memory.DebugRead(0xFF21)) }},
		{name: "FF22 NR43", valueSync: func(gb *gameboy.GameBoy) string { return fmt.Sprintf("%02X", gb.Memory.DebugRead(0xFF22)) }},
		{name: "FF23 NR44", valueSync: func(gb *gameboy.GameBoy) string { return fmt.Sprintf("%02X", gb.Memory.DebugRead(0xFF23)) }},
	}
	return newPanel("Ch4 (Noise)", entries...)
}

func newSoundControlPanel() *panel {
	entries := []panelEntry{
		{name: "FF24 NR50", valueSync: func(gb *gameboy.GameBoy) string { return fmt.Sprintf("%02X", gb.Memory.DebugRead(0xFF24)) }},
		{name: "FF25 NR51", valueSync: func(gb *gameboy.GameBoy) string { return fmt.Sprintf("%02X", gb.Memory.DebugRead(0xFF25)) }},
		{name: "FF26 NR52", valueSync: func(gb *gameboy.GameBoy) string { return fmt.Sprintf("%02X", gb.Memory.DebugRead(0xFF26)) }},
	}
	return newPanel("Sound Control", entries...)
}

func newWaveRamPanel() *panel {
	entries := []panelEntry{
		{name: "", valueSync: func(gb *gameboy.GameBoy) string {
			return fmt.Sprintf("%02X %02X %02X %02X  %02X %02X %02X %02X",
				gb.Memory.DebugRead(0xFF30),
				gb.Memory.DebugRead(0xFF31),
				gb.Memory.DebugRead(0xFF32),
				gb.Memory.DebugRead(0xFF33),
				gb.Memory.DebugRead(0xFF34),
				gb.Memory.DebugRead(0xFF35),
				gb.Memory.DebugRead(0xFF36),
				gb.Memory.DebugRead(0xFF37),
			)
		}},
		{name: "", valueSync: func(gb *gameboy.GameBoy) string {
			return fmt.Sprintf("%02X %02X %02X %02X  %02X %02X %02X %02X",
				gb.Memory.DebugRead(0xFF38),
				gb.Memory.DebugRead(0xFF39),
				gb.Memory.DebugRead(0xFF3A),
				gb.Memory.DebugRead(0xFF3B),
				gb.Memory.DebugRead(0xFF3C),
				gb.Memory.DebugRead(0xFF3D),
				gb.Memory.DebugRead(0xFF3E),
				gb.Memory.DebugRead(0xFF3F),
			)
		}},
	}

	return newPanel("Wave RAM (FF30-FF3F)", entries...)
}

func newTimerPanel() *panel {
	entries := []panelEntry{
		{name: "FF04 DIV", valueSync: func(gb *gameboy.GameBoy) string { return fmt.Sprintf("%02X", gb.Memory.DebugRead(0xFF04)) }},
		{name: "FF05 TIMA", valueSync: func(gb *gameboy.GameBoy) string { return fmt.Sprintf("%02X", gb.Memory.DebugRead(0xFF05)) }},
		{name: "FF06 TMA", valueSync: func(gb *gameboy.GameBoy) string { return fmt.Sprintf("%02X", gb.Memory.DebugRead(0xFF06)) }},
		{name: "FF07 TAC", valueSync: func(gb *gameboy.GameBoy) string { return fmt.Sprintf("%02X", gb.Memory.DebugRead(0xFF07)) }},
	}
	return newPanel("Timer", entries...)
}
