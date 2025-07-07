package debugger

import (
	"github.com/danielecanzoneri/gb-emulator/gameboy"
	"github.com/ebitenui/ebitenui/widget"
	"golang.org/x/image/colornames"
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

	// Create container
	p.Container = newContainer(widget.DirectionVertical)

	// Panel title
	titleLabel := newLabel(title, colornames.Yellow)
	p.AddChild(titleLabel)

	// Two vertical containers: one with labels and one with values
	labels := newContainer(widget.DirectionVertical)
	for _, entry := range entries {
		l := newLabel(entry.name, colornames.White)
		labels.AddChild(l)
	}

	// Create a handler that syncs all entries
	p.Sync = func(gb *gameboy.GameBoy) {}
	values := newContainer(widget.DirectionVertical)
	for _, entry := range entries {
		l := newLabel("", colornames.White)
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
	core.AddChild(labels, values)
	p.AddChild(core)

	return p
}
