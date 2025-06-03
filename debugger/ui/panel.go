package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

const (
	panelTitleMargin    = 1
	panelBorderPadding  = 4
	panelColumnsPadding = 8
)

type panelRenderer struct {
	title       *canvas.Text
	titleColumn fyne.CanvasObject
	valueColumn fyne.CanvasObject
}

func (p panelRenderer) Destroy() {
}

func (p panelRenderer) Layout(_ fyne.Size) {
	titleW, titleH := p.title.MinSize().Components()
	c1W, c1H := p.titleColumn.MinSize().Components()
	c2W, c2H := p.valueColumn.MinSize().Components()

	// Add padding on the left and on top of title
	p.title.Resize(fyne.NewSize(titleW, titleH))
	p.title.Move(fyne.NewPos(panelBorderPadding, panelBorderPadding))

	// Position columns with padding from title
	p.titleColumn.Resize(fyne.NewSize(c1W, c1H))
	p.titleColumn.Move(fyne.NewPos(panelBorderPadding, panelBorderPadding+titleH+panelTitleMargin))
	p.valueColumn.Resize(fyne.NewSize(c2W, c2H))
	p.valueColumn.Move(fyne.NewPos(panelBorderPadding+c1W+panelColumnsPadding, panelBorderPadding+titleH+panelTitleMargin))
}

func (p panelRenderer) MinSize() fyne.Size {
	c1Size := p.titleColumn.MinSize()
	c2Size := p.valueColumn.MinSize()
	columnsHeight := max(c1Size.Height, c2Size.Height)
	columnsWidth := c1Size.Width + c2Size.Width + panelColumnsPadding

	titleSize := p.title.Size()
	return fyne.NewSize(
		max(titleSize.Width, columnsWidth)+2*panelBorderPadding,
		titleSize.Height+panelTitleMargin+columnsHeight+2*panelBorderPadding,
	)
}

func (p panelRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{p.title, p.valueColumn, p.titleColumn}
}

func (p panelRenderer) Refresh() {
	p.title.Refresh()
	p.titleColumn.Refresh()
	p.valueColumn.Refresh()
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

	return &panelRenderer{
		title:       titleText,
		titleColumn: container.NewVBox(titles...),
		valueColumn: container.NewVBox(values...),
	}
}
