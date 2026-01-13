package debugger

import (
	"github.com/danielecanzoneri/lucky-boy/ui/graphics"
	"image/color"

	"github.com/danielecanzoneri/lucky-boy/gameboy"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
)

type panelEntry struct {
	name      string
	valueSync func(gb *gameboy.GameBoy) string
}

type panel struct {
	*widget.Container

	Sync func(gb *gameboy.GameBoy)
}

func newPanel(title string, entries ...panelEntry) *panel {
	return newPanelWithHeader(title, nil, entries...)
}

func newPanelWithHeader(title string, headerUpdate func(gb *gameboy.GameBoy) string, entries ...panelEntry) *panel {
	p := new(panel)

	// Create container (background image should account for padding)
	backgroundImage := image.NewBorderedNineSliceColor(theme.Debugger.Main.Color, color.Transparent, theme.Debugger.Padding)
	p.Container = widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Padding(&widget.Insets{
				Top:    theme.Debugger.Padding,
				Left:   2 * theme.Debugger.Padding,
				Right:  2 * theme.Debugger.Padding,
				Bottom: theme.Debugger.Padding,
			}),
		)),
		widget.ContainerOpts.BackgroundImage(backgroundImage),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Stretch: true,
			}),
		),
	)

	// Create a handler that syncs all entries
	p.Sync = func(gb *gameboy.GameBoy) {}

	// Panel title
	titleLabel := widget.NewText(
		widget.TextOpts.Text(title, &font, theme.Debugger.TitleColor),
		widget.TextOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionCenter, // Center title
			}),
		),
	)
	p.AddChild(titleLabel)

	// Add header if present
	if headerUpdate != nil {
		l := newLabel("", theme.Debugger.HeaderColor)
		p.AddChild(l)

		oldSync := p.Sync
		p.Sync = func(gb *gameboy.GameBoy) {
			oldSync(gb)
			l.Label = headerUpdate(gb)
		}
	}

	// Two vertical containers: one with labels and one with values
	labels := newContainer(widget.DirectionVertical)
	for _, entry := range entries {
		if entry.name == "" {
			continue
		}

		l := newLabel(entry.name+":", theme.Debugger.LabelColor)
		labels.AddChild(l)
	}

	values := newContainer(widget.DirectionVertical)
	for _, entry := range entries {
		l := newLabel("", theme.Debugger.LabelColor)
		values.AddChild(l)

		oldSync := p.Sync
		p.Sync = func(gb *gameboy.GameBoy) {
			oldSync(gb)
			l.Label = entry.valueSync(gb)
		}
	}

	core := widget.NewContainer(widget.ContainerOpts.Layout(
		widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionHorizontal),
			widget.RowLayoutOpts.Spacing(10),
		),
	))
	if len(labels.Children()) == 0 { // Do not add title container if all titles are empty
		core.AddChild(values)
	} else {
		core.AddChild(labels, values)
	}
	p.AddChild(core)

	return p
}
