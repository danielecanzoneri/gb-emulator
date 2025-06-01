package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/danielecanzoneri/gb-emulator/pkg/debug"
)

type registersViewer struct {
	widget.BaseWidget

	// Sound
	//ch1     soundCh1Panel
	//ch2     soundCh2Panel
	//ch3     soundCh3Panel
	//ch4     soundCh4Panel
	//control soundControlPanel
	//waveRam waveRamPanel
	//
	//// CPU
	//cpu        cpuPanel
	//interrupts interruptsPanel
	//
	//lcd   lcdPanel
	timer *timerPanel
}

func newRegisterViewer() *registersViewer {
	v := &registersViewer{
		timer: newTimerPanel(),
	}
	v.ExtendBaseWidget(v)
	return v
}

func (v *registersViewer) Update(state *debug.GameBoyState) {
	v.timer.Update(state)
}

func (v *registersViewer) CreateRenderer() fyne.WidgetRenderer {
	//soundRow := container.NewHBox(v.ch1, v.ch2, v.ch3, v.ch4)
	//cpuColumn := container.NewVBox(v.cpu, v.interrupts)
	//timerSoundControlColumn := container.NewVBox(v.control, v.timer)
	//centralPanels := container.NewHBox(
	//	cpuColumn,               // left
	//	timerSoundControlColumn, // right
	//)

	c := container.NewBorder(
		nil, nil, nil, nil,
		v.timer,
	)
	return widget.NewSimpleRenderer(c)
}

func comparator[T uint8 | uint16](a, b T) bool {
	return a < b
}

func uint8BindingSprintf(u binding.Int) binding.String {
	return binding.NewSprintf("%02X", u)
}

//func uint16BindingSprintf(u binding.Item[uint16]) binding.String {
//	return binding.NewSprintf("%04X", u)
//}
//
//func newCpuPanel() *cpuPanel {
//	p := new(cpuPanel)
//	p.ExtendBaseWidget(p)
//	return p
//}
//
//type cpuPanel struct {
//	widget.BaseWidget
//	AF uint16
//	BC uint16
//	DE uint16
//	HL uint16
//	PC uint16
//	SP uint16
//}
//
//func newInterruptsPanel() *interruptsPanel {
//	p := new(interruptsPanel)
//	p.ExtendBaseWidget(p)
//	return p
//}
//
//type interruptsPanel struct {
//	widget.BaseWidget
//	IF  uint8 // 0xFF0F
//	IE  uint8 // 0xFFFF
//	IME bool
//}
//
//func newLcdPanel() *lcdPanel {
//	p := new(lcdPanel)
//	p.ExtendBaseWidget(p)
//	return p
//}
//
//type lcdPanel struct {
//	widget.BaseWidget
//	LCDC uint8 // 0xFF40
//	STAT uint8 // 0xFF41
//	SCY  uint8 // 0xFF42
//	SCX  uint8 // 0xFF43
//	LY   uint8 // 0xFF44
//	LYC  uint8 // 0xFF45
//	DMA  uint8 // 0xFF46
//	BGP  uint8 // 0xFF47
//	OBP0 uint8 // 0xFF48
//	OBP1 uint8 // 0xFF49
//	WY   uint8 // 0xFF4A
//	WX   uint8 // 0xFF4B
//}
//
//func newSoundCh1Panel() *soundCh1Panel {
//	p := new(soundCh1Panel)
//	p.ExtendBaseWidget(p)
//	return p
//}
//
//type soundCh1Panel struct {
//	widget.BaseWidget
//	NR10 uint8 // 0xFF10
//	NR11 uint8 // 0xFF11
//	NR12 uint8 // 0xFF12
//	NR13 uint8 // 0xFF13
//	NR14 uint8 // 0xFF14
//}
//
//func newSoundCh2Panel() *soundCh2Panel {
//	p := new(soundCh2Panel)
//	p.ExtendBaseWidget(p)
//	return p
//}
//
//type soundCh2Panel struct {
//	widget.BaseWidget
//	NR21 uint8 // 0xFF16
//	NR22 uint8 // 0xFF17
//	NR23 uint8 // 0xFF18
//	NR24 uint8 // 0xFF19
//}
//
//func newSoundCh3Panel() *soundCh3Panel {
//	p := new(soundCh3Panel)
//	p.ExtendBaseWidget(p)
//	return p
//}
//
//type soundCh3Panel struct {
//	widget.BaseWidget
//	NR30 uint8 // 0xFF1A
//	NR31 uint8 // 0xFF1B
//	NR32 uint8 // 0xFF1C
//	NR33 uint8 // 0xFF1D
//	NR34 uint8 // 0xFF1E
//}
//
//func newSoundCh4Panel() *soundCh4Panel {
//	p := new(soundCh4Panel)
//	p.ExtendBaseWidget(p)
//	return p
//}
//
//type soundCh4Panel struct {
//	widget.BaseWidget
//	NR41 uint8 // 0xFF20
//	NR42 uint8 // 0xFF21
//	NR43 uint8 // 0xFF22
//	NR44 uint8 // 0xFF23
//}
//
//func newSoundControlPanel() *soundControlPanel {
//	p := new(soundControlPanel)
//	p.ExtendBaseWidget(p)
//	return p
//}
//
//type soundControlPanel struct {
//	widget.BaseWidget
//	NR50 uint8 // 0xFF24
//	NR51 uint8 // 0xFF25
//	NR52 uint8 // 0xFF26
//}
//
//func newWaveRamPanel() *waveRamPanel {
//	p := new(waveRamPanel)
//	p.ExtendBaseWidget(p)
//	return p
//}
//
//type waveRamPanel struct {
//	widget.BaseWidget
//	// FF30 - FF3F
//	values [16]uint8
//}

type timerPanel struct {
	widget.BaseWidget
	// 0xFF04
	DIV binding.Int
	// 0xFF05
	TIMA binding.Int
	// 0xFF06
	TMA binding.Int
	// 0xFF07
	TAC binding.Int
}

func newTimerPanel() *timerPanel {
	p := &timerPanel{
		DIV:  binding.NewInt(),
		TIMA: binding.NewInt(),
		TMA:  binding.NewInt(),
		TAC:  binding.NewInt(),
	}

	p.ExtendBaseWidget(p)
	return p
}

func (p *timerPanel) Update(state *debug.GameBoyState) {
	_ = p.DIV.Set(int(state.Memory[0xFF04]))
	_ = p.TIMA.Set(int(state.Memory[0xFF05]))
	_ = p.TMA.Set(int(state.Memory[0xFF06]))
	_ = p.TAC.Set(int(state.Memory[0xFF07]))
}

func (p *timerPanel) CreateRenderer() fyne.WidgetRenderer {
	titles := []string{
		"FF04 DIV",
		"FF05 TIMA",
		"FF06 TMA",
		"FF07 TAC"}
	values := []binding.String{
		uint8BindingSprintf(p.DIV),
		uint8BindingSprintf(p.TIMA),
		uint8BindingSprintf(p.TMA),
		uint8BindingSprintf(p.TAC),
	}

	return newPanelRenderer(
		p.Theme(),
		titles, values,
	)
}

func newPanelRenderer(th fyne.Theme, names []string, data []binding.String) fyne.WidgetRenderer {
	v := fyne.CurrentApp().Settings().ThemeVariant()

	// Panel title
	title := canvas.NewText("Timer", th.Color(theme.ColorYellow, v))
	title.TextStyle = fyne.TextStyle{Bold: true}
	title.TextSize = 12

	// Create two vertical containers, one for titles and one for values
	titles := make([]fyne.CanvasObject, len(names))
	for i, name := range names {
		titles[i] = widget.NewLabel(name)
	}
	values := make([]fyne.CanvasObject, len(data))
	for i, d := range data {
		values[i] = widget.NewLabelWithData(d)
	}

	c := container.NewBorder(
		title,
		nil,
		container.NewVBox(titles...),
		container.NewVBox(values...),
	)
	return widget.NewSimpleRenderer(c)
}
