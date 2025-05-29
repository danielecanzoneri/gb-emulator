package ui

import (
	"fmt"
	"image/color"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

const (
	panelRowTextPadding = 2
)

type panelRow struct {
	widget.BaseWidget

	title       string
	content     string
	titleLength int // Length of the title (used for padding)
}

func newPanelRow(title string, content string) *panelRow {
	row := &panelRow{
		title:   title,
		content: content,
	}
	row.ExtendBaseWidget(row)
	return row
}

func (r *panelRow) CreateRenderer() fyne.WidgetRenderer {
	th := r.Theme()
	v := fyne.CurrentApp().Settings().ThemeVariant()

	// Text placeholder
	text := canvas.NewText(fmt.Sprintf("%s: %s", r.title, r.content), th.Color(theme.ColorNameForeground, v))
	text.TextStyle = fyne.TextStyle{Monospace: true}

	background := canvas.NewRectangle(color.Transparent)

	return &panelRowRenderer{
		entry:      r,
		text:       text,
		background: background,
	}
}

type panelRowRenderer struct {
	entry *panelRow

	text       *canvas.Text
	background *canvas.Rectangle
}

func (r *panelRowRenderer) Destroy() {}

func (r *panelRowRenderer) Layout(size fyne.Size) {
	r.background.Resize(size)
	r.text.Resize(fyne.NewSize(size.Width-2*panelRowTextPadding, size.Height))
	r.text.Move(fyne.NewPos(panelRowTextPadding, 0))
}

func (r *panelRowRenderer) MinSize() fyne.Size {
	return fyne.NewSize(r.text.MinSize().Width+2*panelRowTextPadding, r.text.MinSize().Height)
}

func (r *panelRowRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.background, r.text}
}

func (r *panelRowRenderer) Refresh() {
	if r.entry.title != "" {
		spaces := strings.Repeat(" ", r.entry.titleLength-len(r.entry.title))
		r.text.Text = fmt.Sprintf("%s: %s%s", r.entry.title, spaces, r.entry.content)
	} else {
		r.text.Text = r.entry.content
	}

	r.text.Refresh()
}

type Panel struct {
	widget.BaseWidget

	title string
	rows  []fyne.CanvasObject

	maxTitleLength int // Maximum length of the title (used for padding)
}

func NewPanel(title string) *Panel {
	p := &Panel{
		title: title,
		rows:  make([]fyne.CanvasObject, 0),
	}
	p.ExtendBaseWidget(p)
	return p
}

func (p *Panel) AddRow(title string, content string) {
	// Update the max title length
	if len(title) > p.maxTitleLength {
		p.maxTitleLength = len(title)

		// Update the rows with the new max title length
		for _, row := range p.rows {
			row.(*panelRow).titleLength = p.maxTitleLength
		}
	}

	// Add the new row
	newRow := newPanelRow(title, content)
	newRow.titleLength = p.maxTitleLength

	p.rows = append(p.rows, newRow)
}

func (p *Panel) CreateRenderer() fyne.WidgetRenderer {
	th := p.Theme()
	v := fyne.CurrentApp().Settings().ThemeVariant()

	// Panel title
	title := canvas.NewText(p.title, th.Color(theme.ColorYellow, v))
	title.TextStyle = fyne.TextStyle{Bold: true}
	title.TextSize = 12

	// Create a vertical box for the rows
	items := append([]fyne.CanvasObject{title}, p.rows...)
	vBox := container.NewVBox(items...)

	return widget.NewSimpleRenderer(vBox)
}
