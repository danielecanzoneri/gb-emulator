package debugger

import (
	"github.com/danielecanzoneri/gb-emulator/gameboy"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"image/color"
)

const panelsPadding = 5

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
	backgroundImage := image.NewBorderedNineSliceColor(mainColor, color.Transparent, panelsPadding)
	p.Container = widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Padding(widget.NewInsetsSimple(panelsPadding)),
		)),
		widget.ContainerOpts.BackgroundImage(backgroundImage),
	)

	// Panel title
	titleLabel := newLabel(title, titleColor)
	p.AddChild(titleLabel)

	// Two vertical containers: one with labels and one with values
	labels := newContainer(widget.DirectionVertical)
	for _, entry := range entries {
		if entry.name == "" {
			continue
		}

		l := newLabel(entry.name+":", labelColor)
		labels.AddChild(l)
	}

	// Create a handler that syncs all entries
	p.Sync = func(gb *gameboy.GameBoy) {}
	values := newContainer(widget.DirectionVertical)
	for _, entry := range entries {
		l := newLabel("", labelColor)
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
