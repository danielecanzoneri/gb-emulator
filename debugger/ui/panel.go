package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

const (
	panelRowTextPadding = 2
)

type panelRowTitle interface {
	Title() string
	Display() bool
}

type panelRowValue interface {
	Value() string
}

type namedAddress struct {
	addr uint16
	name string
}

func (r *namedAddress) Title() string {
	return fmt.Sprintf("%04X %s", r.addr, r.name)
}

func (r *namedAddress) Display() bool {
	return true
}

type register struct {
	name string
}

func (r *register) Title() string {
	return r.name
}

func (r *register) Display() bool {
	return true
}

type unnamedRegister struct {
}

func (r *unnamedRegister) Title() string {
	return ""
}

func (r *unnamedRegister) Display() bool {
	return false
}

type uint8Value struct {
	value uint8
}

func (v *uint8Value) Value() string {
	return fmt.Sprintf("%02X", v.value)
}

type uint16Value struct {
	value uint16
}

func (v *uint16Value) Value() string {
	return fmt.Sprintf("%04X", v.value)
}

type flagValue struct {
	value bool
}

func (v *flagValue) Value() string {
	if v.value {
		return "enabled"
	}
	return "disabled"
}

type arrayValue struct {
	values [16]uint8
}

func (v *arrayValue) Value() string {
	return fmt.Sprintf("%02X %02X %02X %02X %02X %02X %02X %02X %02X %02X %02X %02X %02X %02X %02X %02X",
		v.values[0], v.values[1], v.values[2], v.values[3], v.values[4], v.values[5], v.values[6], v.values[7],
		v.values[8], v.values[9], v.values[10], v.values[11], v.values[12], v.values[13], v.values[14], v.values[15])
}

type panelRow struct {
	title panelRowTitle
	value panelRowValue
}

func newPanelRow(title panelRowTitle, value panelRowValue) *panelRow {
	return &panelRow{
		title: title,
		value: value,
	}
}

func (r *panelRow) Title() string {
	return r.title.Title()
}

func (r *panelRow) Value() string {
	return r.value.Value()
}

type panel struct {
	widget.BaseWidget

	title string
	rows  []*panelRow
}

func newPanel(title string) *panel {
	p := &panel{
		title: title,
		rows:  make([]*panelRow, 0),
	}
	p.ExtendBaseWidget(p)
	return p
}

func (p *panel) AddRow(title panelRowTitle, value panelRowValue) {
	// Add the new row
	newRow := newPanelRow(title, value)
	p.rows = append(p.rows, newRow)
}

func (p *panel) CreateRenderer() fyne.WidgetRenderer {
	th := p.Theme()
	v := fyne.CurrentApp().Settings().ThemeVariant()

	// Panel title
	title := canvas.NewText(p.title, th.Color(theme.ColorYellow, v))
	title.TextStyle = fyne.TextStyle{Bold: true}
	title.TextSize = 12

	// Create two vertical containers, one for titles and one for values
	titles := make([]fyne.CanvasObject, len(p.rows))
	values := make([]fyne.CanvasObject, len(p.rows))
	for i, row := range p.rows {
		titles[i] = widget.NewLabel(row.Title())
		values[i] = widget.NewLabel(row.Value())
	}

	c := container.NewBorder(
		title,
		nil,
		container.NewVBox(titles...),
		container.NewVBox(values...),
	)
	return widget.NewSimpleRenderer(c)
}

type panelOnlyValue struct {
	widget.BaseWidget

	title string
	rows  []panelRowValue
}

func newPanelOnlyValue(title string) *panelOnlyValue {
	p := &panelOnlyValue{
		title: title,
		rows:  make([]panelRowValue, 0),
	}
	p.ExtendBaseWidget(p)
	return p
}

func (p *panelOnlyValue) AddRow(value panelRowValue) {
	// Add the new row
	p.rows = append(p.rows, value)
}

func (p *panelOnlyValue) CreateRenderer() fyne.WidgetRenderer {
	th := p.Theme()
	v := fyne.CurrentApp().Settings().ThemeVariant()

	// Panel title
	title := canvas.NewText(p.title, th.Color(theme.ColorYellow, v))
	title.TextStyle = fyne.TextStyle{Bold: true}
	title.TextSize = 12

	// Create a vertical container for the rows
	items := make([]fyne.CanvasObject, len(p.rows)+1)
	items[0] = title
	for i, row := range p.rows {
		items[i+1] = widget.NewLabel(row.Value())
	}

	vBox := container.NewVBox(items...)
	return widget.NewSimpleRenderer(vBox)
}
