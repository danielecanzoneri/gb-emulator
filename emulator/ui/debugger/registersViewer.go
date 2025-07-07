package debugger

import (
	"fmt"
	"github.com/danielecanzoneri/gb-emulator/emulator/internal/gameboy"
	"github.com/ebitenui/ebitenui/widget"
)

type registersViewer struct {
	*widget.Container

	// Sound
	//ch1     *soundCh1Panel
	//ch2     *soundCh2Panel
	//ch3     *soundCh3Panel
	//ch4     *soundCh4Panel
	//control *soundControlPanel
	//waveRam *waveRamPanel

	// CPU
	cpu *panel
	//interrupts *interruptsPanel

	//lcd   *lcdPanel
	//timer *timerPanel
}

func newRegisterViewer() *registersViewer {
	rv := &registersViewer{
		//ch1:        newSoundCh1Panel(),
		//ch2:        newSoundCh2Panel(),
		//ch3:        newSoundCh3Panel(),
		//ch4:        newSoundCh4Panel(),
		//control:    newSoundControlPanel(),
		//waveRam:    newWaveRamPanel(),
		cpu: newCpuPanel(),
		//interrupts: newInterruptsPanel(),
		//lcd:        newLcdPanel(),
		//timer:      newTimerPanel(),
	}
	rv.Container = newContainer(widget.DirectionVertical, rv.cpu)
	return rv
}

func (v *registersViewer) Sync(gb *gameboy.GameBoy) {
	//v.ch1.Sync(gb)
	//v.ch2.Sync(gb)
	//v.ch3.Sync(gb)
	//v.ch4.Sync(gb)
	//v.control.Sync(gb)
	//v.waveRam.Sync(gb)
	v.cpu.Sync(gb)
	//v.interrupts.Sync(gb)
	//v.lcd.Sync(gb)
	//v.timer.Sync(gb)
}

func newCpuPanel() *panel {
	entries := []panelEntry{
		{name: "AF", valueSync: func(gb *gameboy.GameBoy) string {
			return fmt.Sprintf("%04X", gb.CPU.ReadAF())
		}},
		{name: "BC", valueSync: func(gb *gameboy.GameBoy) string {
			return fmt.Sprintf("%04X", gb.CPU.ReadBC())
		}},
		{name: "DE", valueSync: func(gb *gameboy.GameBoy) string {
			return fmt.Sprintf("%04X", gb.CPU.ReadDE())
		}},
		{name: "HL", valueSync: func(gb *gameboy.GameBoy) string {
			return fmt.Sprintf("%04X", gb.CPU.ReadHL())
		}},
		{name: "PC", valueSync: func(gb *gameboy.GameBoy) string {
			return fmt.Sprintf("%04X", gb.CPU.ReadPC())
		}},
		{name: "SP", valueSync: func(gb *gameboy.GameBoy) string {
			return fmt.Sprintf("%04X", gb.CPU.ReadSP())
		}},
	}

	return newPanel("CPU", entries...)
}
