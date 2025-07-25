package debugger

import (
	"github.com/danielecanzoneri/gb-emulator/ui/theme"
	"image/color"

	"github.com/danielecanzoneri/gb-emulator/gameboy"
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
	p := new(panel)

	// Create container (background image should account for padding)
	backgroundImage := image.NewBorderedNineSliceColor(theme.Debugger.Main.Color, color.Transparent, theme.Debugger.Padding)
	p.Container = widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Padding(widget.NewInsetsSimple(theme.Debugger.Padding)),
		)),
		widget.ContainerOpts.BackgroundImage(backgroundImage),
	)

	// Panel title
	titleLabel := newLabel(title, theme.Debugger.TitleColor)
	p.AddChild(titleLabel)

	// Two vertical containers: one with labels and one with values
	labels := newContainer(widget.DirectionVertical)
	for _, entry := range entries {
		if entry.name == "" {
			continue
		}

		l := newLabel(entry.name+":", theme.Debugger.LabelColor)
		labels.AddChild(l)
	}

	// Create a handler that syncs all entries
	p.Sync = func(gb *gameboy.GameBoy) {}
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
