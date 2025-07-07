package debugger

import (
	"fmt"
	"github.com/danielecanzoneri/gb-emulator/emulator/internal/gameboy"
	"github.com/ebitenui/ebitenui/widget"
	"golang.org/x/image/colornames"
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
	cpu *cpuPanel
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

type cpuPanel struct {
	*widget.Container

	labelAF *widget.Text
	labelBC *widget.Text
	labelDE *widget.Text
	labelHL *widget.Text
	labelPC *widget.Text
	labelSP *widget.Text
}

func newCpuPanel() *cpuPanel {
	p := new(cpuPanel)

	// Create container
	p.Container = newContainer(widget.DirectionVertical)

	// Panel title
	titleLabel := newLabel("CPU", colornames.Yellow)
	p.AddChild(titleLabel)

	// Two vertical containers: one with labels and one with values
	labels := newContainer(widget.DirectionVertical,
		newLabel("AF", colornames.White),
		newLabel("BC", colornames.White),
		newLabel("DE", colornames.White),
		newLabel("HL", colornames.White),
		newLabel("PC", colornames.White),
		newLabel("SP", colornames.White),
	)

	p.labelAF = newLabel("", colornames.White)
	p.labelBC = newLabel("", colornames.White)
	p.labelDE = newLabel("", colornames.White)
	p.labelHL = newLabel("", colornames.White)
	p.labelPC = newLabel("", colornames.White)
	p.labelSP = newLabel("", colornames.White)
	values := newContainer(widget.DirectionVertical,
		p.labelAF, p.labelBC, p.labelDE, p.labelHL, p.labelPC, p.labelSP,
	)

	core := widget.NewContainer(widget.ContainerOpts.Layout(
		widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionHorizontal),
			widget.RowLayoutOpts.Spacing(10),
		),
	))
	core.AddChild(labels, values)
	p.AddChild(core)

	return p
}

func (p *cpuPanel) Sync(gb *gameboy.GameBoy) {
	p.labelAF.Label = fmt.Sprintf("%04X", gb.CPU.ReadAF())
	p.labelBC.Label = fmt.Sprintf("%04X", gb.CPU.ReadBC())
	p.labelDE.Label = fmt.Sprintf("%04X", gb.CPU.ReadDE())
	p.labelHL.Label = fmt.Sprintf("%04X", gb.CPU.ReadHL())
	p.labelPC.Label = fmt.Sprintf("%04X", gb.CPU.ReadPC())
	p.labelSP.Label = fmt.Sprintf("%04X", gb.CPU.ReadSP())
}
