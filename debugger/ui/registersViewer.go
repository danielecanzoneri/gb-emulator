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
	ch1     *soundCh1Panel
	ch2     *soundCh2Panel
	ch3     *soundCh3Panel
	ch4     *soundCh4Panel
	control *soundControlPanel
	waveRam *waveRamPanel

	// CPU
	cpu        *cpuPanel
	interrupts *interruptsPanel

	lcd   *lcdPanel
	timer *timerPanel
}

func newRegisterViewer() *registersViewer {
	v := &registersViewer{
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
	v.ExtendBaseWidget(v)
	return v
}

func (v *registersViewer) Update(state *debug.GameBoyState) {
	v.ch1.Update(state)
	v.ch2.Update(state)
	v.ch3.Update(state)
	v.ch4.Update(state)
	v.control.Update(state)
	v.waveRam.Update(state)
	v.cpu.Update(state)
	v.interrupts.Update(state)
	v.lcd.Update(state)
	v.timer.Update(state)
}

func (v *registersViewer) CreateRenderer() fyne.WidgetRenderer {
	soundRow := container.NewHBox(v.ch1, v.ch2, v.ch3, v.ch4)
	cpuColumn := container.NewVBox(v.cpu, v.interrupts)
	timerSoundControlColumn := container.NewVBox(v.control, v.timer)
	centralPanels := container.NewVBox(
		container.NewHBox(
			cpuColumn,               // left
			timerSoundControlColumn, // right
		),
		v.waveRam,
	)

	c := container.NewBorder(
		soundRow, nil, v.lcd, nil,
		centralPanels,
	)
	return widget.NewSimpleRenderer(c)
}

func uint8BindToString(u binding.Int) binding.String {
	return binding.NewSprintf("%02X", u)
}

func uint16BindToString(u binding.Int) binding.String {
	return binding.NewSprintf("%04X", u)
}

func boolBindToString(b binding.Bool) binding.String {
	return binding.BoolToString(b)
}

type cpuPanel struct {
	widget.BaseWidget
	AF binding.Int
	BC binding.Int
	DE binding.Int
	HL binding.Int
	PC binding.Int
	SP binding.Int
}

func newCpuPanel() *cpuPanel {
	p := &cpuPanel{
		AF: binding.NewInt(),
		BC: binding.NewInt(),
		DE: binding.NewInt(),
		HL: binding.NewInt(),
		PC: binding.NewInt(),
		SP: binding.NewInt(),
	}

	p.ExtendBaseWidget(p)
	return p
}

func (p *cpuPanel) Update(state *debug.GameBoyState) {
	_ = p.AF.Set(int(state.AF))
	_ = p.BC.Set(int(state.BC))
	_ = p.DE.Set(int(state.DE))
	_ = p.HL.Set(int(state.HL))
	_ = p.PC.Set(int(state.PC))
	_ = p.SP.Set(int(state.SP))
}

func (p *cpuPanel) CreateRenderer() fyne.WidgetRenderer {
	titles := []string{
		"AF", "BC", "DE", "HL", "PC", "SP",
	}
	values := []binding.String{
		uint16BindToString(p.AF),
		uint16BindToString(p.BC),
		uint16BindToString(p.DE),
		uint16BindToString(p.HL),
		uint16BindToString(p.PC),
		uint16BindToString(p.SP),
	}

	return newPanelRenderer(
		p.Theme(),
		"CPU",
		titles, values,
	)
}

type interruptsPanel struct {
	widget.BaseWidget
	IF  binding.Int // 0xFF0F
	IE  binding.Int // 0xFFFF
	IME binding.Bool
}

func newInterruptsPanel() *interruptsPanel {
	p := &interruptsPanel{
		IF:  binding.NewInt(),
		IE:  binding.NewInt(),
		IME: binding.NewBool(),
	}

	p.ExtendBaseWidget(p)
	return p
}

func (p *interruptsPanel) Update(state *debug.GameBoyState) {
	_ = p.IF.Set(int(state.Memory[0xFF0F]))
	_ = p.IE.Set(int(state.Memory[0xFFFF]))
	_ = p.IME.Set(state.IME)
}

func (p *interruptsPanel) CreateRenderer() fyne.WidgetRenderer {
	titles := []string{
		"FF0F IF",
		"FFFF IE",
		"IME",
	}
	values := []binding.String{
		uint8BindToString(p.IF),
		uint8BindToString(p.IE),
		boolBindToString(p.IME),
	}

	return newPanelRenderer(
		p.Theme(),
		"Interrupts",
		titles, values,
	)
}

type lcdPanel struct {
	widget.BaseWidget

	LCDC binding.Int // 0xFF40
	STAT binding.Int // 0xFF41
	SCY  binding.Int // 0xFF42
	SCX  binding.Int // 0xFF43
	LY   binding.Int // 0xFF44
	LYC  binding.Int // 0xFF45
	DMA  binding.Int // 0xFF46
	BGP  binding.Int // 0xFF47
	OBP0 binding.Int // 0xFF48
	OBP1 binding.Int // 0xFF49
	WY   binding.Int // 0xFF4A
	WX   binding.Int // 0xFF4B
}

func newLcdPanel() *lcdPanel {
	p := &lcdPanel{
		LCDC: binding.NewInt(),
		STAT: binding.NewInt(),
		SCY:  binding.NewInt(),
		SCX:  binding.NewInt(),
		LY:   binding.NewInt(),
		LYC:  binding.NewInt(),
		DMA:  binding.NewInt(),
		BGP:  binding.NewInt(),
		OBP0: binding.NewInt(),
		OBP1: binding.NewInt(),
		WY:   binding.NewInt(),
		WX:   binding.NewInt(),
	}

	p.ExtendBaseWidget(p)
	return p
}

func (p *lcdPanel) Update(state *debug.GameBoyState) {
	_ = p.LCDC.Set(int(state.Memory[0xFF40]))
	_ = p.STAT.Set(int(state.Memory[0xFF41]))
	_ = p.SCY.Set(int(state.Memory[0xFF42]))
	_ = p.SCX.Set(int(state.Memory[0xFF43]))
	_ = p.LY.Set(int(state.Memory[0xFF44]))
	_ = p.LYC.Set(int(state.Memory[0xFF45]))
	_ = p.DMA.Set(int(state.Memory[0xFF46]))
	_ = p.BGP.Set(int(state.Memory[0xFF47]))
	_ = p.OBP0.Set(int(state.Memory[0xFF48]))
	_ = p.OBP1.Set(int(state.Memory[0xFF49]))
	_ = p.WY.Set(int(state.Memory[0xFF4A]))
	_ = p.WX.Set(int(state.Memory[0xFF4B]))
}

func (p *lcdPanel) CreateRenderer() fyne.WidgetRenderer {
	titles := []string{
		"FF40 LCDC",
		"FF41 STAT",
		"FF42 SCY",
		"FF43 SCX",
		"FF44 LY",
		"FF45 LYC",
		"FF46 DMA",
		"FF47 BGP",
		"FF48 OBP0",
		"FF49 OBP1",
		"FF4A WY",
		"FF4B WX",
	}
	values := []binding.String{
		uint8BindToString(p.LCDC),
		uint8BindToString(p.STAT),
		uint8BindToString(p.SCY),
		uint8BindToString(p.SCX),
		uint8BindToString(p.LY),
		uint8BindToString(p.LYC),
		uint8BindToString(p.DMA),
		uint8BindToString(p.BGP),
		uint8BindToString(p.OBP0),
		uint8BindToString(p.OBP1),
		uint8BindToString(p.WY),
		uint8BindToString(p.WX),
	}

	return newPanelRenderer(
		p.Theme(),
		"LCD",
		titles, values,
	)
}

type soundCh1Panel struct {
	widget.BaseWidget

	NR10 binding.Int // 0xFF10
	NR11 binding.Int // 0xFF11
	NR12 binding.Int // 0xFF12
	NR13 binding.Int // 0xFF13
	NR14 binding.Int // 0xFF14
}

func newSoundCh1Panel() *soundCh1Panel {
	p := &soundCh1Panel{
		NR10: binding.NewInt(),
		NR11: binding.NewInt(),
		NR12: binding.NewInt(),
		NR13: binding.NewInt(),
		NR14: binding.NewInt(),
	}

	p.ExtendBaseWidget(p)
	return p
}

func (p *soundCh1Panel) Update(state *debug.GameBoyState) {
	_ = p.NR10.Set(int(state.Memory[0xFF10]))
	_ = p.NR11.Set(int(state.Memory[0xFF11]))
	_ = p.NR12.Set(int(state.Memory[0xFF12]))
	_ = p.NR13.Set(int(state.Memory[0xFF13]))
	_ = p.NR14.Set(int(state.Memory[0xFF14]))
}

func (p *soundCh1Panel) CreateRenderer() fyne.WidgetRenderer {
	titles := []string{
		"FF10 NR10",
		"FF11 NR11",
		"FF12 NR12",
		"FF13 NR13",
		"FF14 NR14",
	}
	values := []binding.String{
		uint8BindToString(p.NR10),
		uint8BindToString(p.NR11),
		uint8BindToString(p.NR12),
		uint8BindToString(p.NR13),
		uint8BindToString(p.NR14),
	}

	return newPanelRenderer(
		p.Theme(),
		"Ch1 (Square)",
		titles, values,
	)
}

type soundCh2Panel struct {
	widget.BaseWidget

	NR21 binding.Int // 0xFF16
	NR22 binding.Int // 0xFF17
	NR23 binding.Int // 0xFF18
	NR24 binding.Int // 0xFF19
}

func newSoundCh2Panel() *soundCh2Panel {
	p := &soundCh2Panel{
		NR21: binding.NewInt(),
		NR22: binding.NewInt(),
		NR23: binding.NewInt(),
		NR24: binding.NewInt(),
	}
	p.ExtendBaseWidget(p)
	return p
}

func (p *soundCh2Panel) Update(state *debug.GameBoyState) {
	_ = p.NR21.Set(int(state.Memory[0xFF16]))
	_ = p.NR22.Set(int(state.Memory[0xFF17]))
	_ = p.NR23.Set(int(state.Memory[0xFF18]))
	_ = p.NR24.Set(int(state.Memory[0xFF19]))
}

func (p *soundCh2Panel) CreateRenderer() fyne.WidgetRenderer {
	titles := []string{
		"FF16 NR21",
		"FF17 NR22",
		"FF18 NR23",
		"FF19 NR24",
	}
	values := []binding.String{
		uint8BindToString(p.NR21),
		uint8BindToString(p.NR22),
		uint8BindToString(p.NR23),
		uint8BindToString(p.NR24),
	}

	return newPanelRenderer(
		p.Theme(),
		"Ch2 (Square)",
		titles, values,
	)
}

type soundCh3Panel struct {
	widget.BaseWidget

	NR30 binding.Int // 0xFF1A
	NR31 binding.Int // 0xFF1B
	NR32 binding.Int // 0xFF1C
	NR33 binding.Int // 0xFF1D
	NR34 binding.Int // 0xFF1E
}

func newSoundCh3Panel() *soundCh3Panel {
	p := &soundCh3Panel{
		NR30: binding.NewInt(),
		NR31: binding.NewInt(),
		NR32: binding.NewInt(),
		NR33: binding.NewInt(),
		NR34: binding.NewInt(),
	}

	p.ExtendBaseWidget(p)
	return p
}

func (p *soundCh3Panel) Update(state *debug.GameBoyState) {
	_ = p.NR30.Set(int(state.Memory[0xFF1A]))
	_ = p.NR31.Set(int(state.Memory[0xFF1B]))
	_ = p.NR32.Set(int(state.Memory[0xFF1C]))
	_ = p.NR33.Set(int(state.Memory[0xFF1D]))
	_ = p.NR34.Set(int(state.Memory[0xFF1E]))
}

func (p *soundCh3Panel) CreateRenderer() fyne.WidgetRenderer {
	titles := []string{
		"FF1A NR30",
		"FF1B NR31",
		"FF1C NR32",
		"FF1D NR33",
		"FF1E NR34",
	}
	values := []binding.String{
		uint8BindToString(p.NR30),
		uint8BindToString(p.NR31),
		uint8BindToString(p.NR32),
		uint8BindToString(p.NR33),
		uint8BindToString(p.NR34),
	}

	return newPanelRenderer(
		p.Theme(),
		"Ch3 (Wave)",
		titles, values,
	)
}

type soundCh4Panel struct {
	widget.BaseWidget

	NR41 binding.Int // 0xFF20
	NR42 binding.Int // 0xFF21
	NR43 binding.Int // 0xFF22
	NR44 binding.Int // 0xFF23
}

func newSoundCh4Panel() *soundCh4Panel {
	p := &soundCh4Panel{
		NR41: binding.NewInt(),
		NR42: binding.NewInt(),
		NR43: binding.NewInt(),
		NR44: binding.NewInt(),
	}

	p.ExtendBaseWidget(p)
	return p
}

func (p *soundCh4Panel) Update(state *debug.GameBoyState) {
	_ = p.NR41.Set(int(state.Memory[0xFF20]))
	_ = p.NR42.Set(int(state.Memory[0xFF21]))
	_ = p.NR43.Set(int(state.Memory[0xFF22]))
	_ = p.NR44.Set(int(state.Memory[0xFF23]))
}

func (p *soundCh4Panel) CreateRenderer() fyne.WidgetRenderer {
	titles := []string{
		"FF20 NR41",
		"FF21 NR42",
		"FF22 NR43",
		"FF23 NR44",
	}
	values := []binding.String{
		uint8BindToString(p.NR41),
		uint8BindToString(p.NR42),
		uint8BindToString(p.NR43),
		uint8BindToString(p.NR44),
	}

	return newPanelRenderer(
		p.Theme(),
		"Ch4 (Noise)",
		titles, values,
	)
}

type soundControlPanel struct {
	widget.BaseWidget

	NR50 binding.Int // 0xFF24
	NR51 binding.Int // 0xFF25
	NR52 binding.Int // 0xFF26
}

func newSoundControlPanel() *soundControlPanel {
	p := &soundControlPanel{
		NR50: binding.NewInt(),
		NR51: binding.NewInt(),
		NR52: binding.NewInt(),
	}

	p.ExtendBaseWidget(p)
	return p
}

func (p *soundControlPanel) Update(state *debug.GameBoyState) {
	_ = p.NR50.Set(int(state.Memory[0xFF24]))
	_ = p.NR51.Set(int(state.Memory[0xFF25]))
	_ = p.NR52.Set(int(state.Memory[0xFF26]))
}

func (p *soundControlPanel) CreateRenderer() fyne.WidgetRenderer {
	titles := []string{
		"FF24 NR50",
		"FF25 NR51",
		"FF26 NR52",
	}
	values := []binding.String{
		uint8BindToString(p.NR50),
		uint8BindToString(p.NR51),
		uint8BindToString(p.NR52),
	}

	return newPanelRenderer(
		p.Theme(),
		"Sound Control",
		titles, values,
	)
}

type waveRamPanel struct {
	widget.BaseWidget
	// FF30 - FF3F
	values [16]binding.Int
}

func newWaveRamPanel() *waveRamPanel {
	p := new(waveRamPanel)
	for i := range p.values {
		p.values[i] = binding.NewInt()
	}

	p.ExtendBaseWidget(p)
	return p
}

func (p *waveRamPanel) Update(state *debug.GameBoyState) {
	for i, bind := range p.values {
		_ = bind.Set(int(state.Memory[0xFF30+uint16(i)]))
	}
}

func (p *waveRamPanel) CreateRenderer() fyne.WidgetRenderer {
	th := p.Theme()
	v := fyne.CurrentApp().Settings().ThemeVariant()

	// Panel title
	title := canvas.NewText("Wave RAM (FF30-FF3F)", th.Color(theme.ColorYellow, v))
	title.TextStyle = fyne.TextStyle{Bold: true}
	title.TextSize = 12

	// Create two vertical containers, one for titles and one for values
	waveRamString := binding.NewSprintf("%02X %02X %02X %02X %02X %02X %02X %02X  %02X %02X %02X %02X %02X %02X %02X %02X",
		p.values[0], p.values[1], p.values[2], p.values[3],
		p.values[4], p.values[5], p.values[6], p.values[7],
		p.values[8], p.values[9], p.values[10], p.values[11],
		p.values[12], p.values[13], p.values[14], p.values[15],
	)
	label := widget.NewLabelWithData(waveRamString)

	c := container.NewVBox(title, label)
	return widget.NewSimpleRenderer(c)
}

type timerPanel struct {
	widget.BaseWidget

	DIV  binding.Int // 0xFF04
	TIMA binding.Int // 0xFF05
	TMA  binding.Int // 0xFF06
	TAC  binding.Int // 0xFF07
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
		"FF07 TAC",
	}
	values := []binding.String{
		uint8BindToString(p.DIV),
		uint8BindToString(p.TIMA),
		uint8BindToString(p.TMA),
		uint8BindToString(p.TAC),
	}

	return newPanelRenderer(
		p.Theme(),
		"Timer",
		titles, values,
	)
}

func newPanelRenderer(th fyne.Theme, title string, names []string, data []binding.String) fyne.WidgetRenderer {
	v := fyne.CurrentApp().Settings().ThemeVariant()

	// Panel title
	titleText := canvas.NewText(title, th.Color(theme.ColorYellow, v))
	titleText.TextStyle = fyne.TextStyle{Bold: true}
	titleText.TextSize = 12

	// Create two vertical containers, one for titles and one for values
	titles := make([]fyne.CanvasObject, len(names))
	for i, name := range names {
		titles[i] = widget.NewLabel(name + " ")
	}
	values := make([]fyne.CanvasObject, len(data))
	for i, d := range data {
		values[i] = widget.NewLabelWithData(d)
	}

	c := container.NewBorder(
		titleText,
		nil,
		container.NewVBox(titles...),
		container.NewVBox(values...),
	)
	return widget.NewSimpleRenderer(c)
}
