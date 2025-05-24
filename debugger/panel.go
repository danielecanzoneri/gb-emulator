package debugger

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

const margin = 8

type Panel struct {
	updateFunc func() (title, content string)
	title      string
	content    string

	face *text.GoTextFace

	width, height int
}

func NewPanel(updateFunc func() (string, string), face *text.GoTextFace) *Panel {
	title, content := updateFunc()
	w, h := text.Measure(title+"\n"+content, face, face.Size)
	return &Panel{
		updateFunc: updateFunc,
		title:      title,
		content:    content,
		face:       face,
		width:      int(math.Ceil(w)) + 2*margin,
		height:     int(math.Ceil(h)) + 2*margin,
	}
}

func (p *Panel) Update() {
	p.title, p.content = p.updateFunc()
}

func (p *Panel) Draw(screen *ebiten.Image, x, y int) {
	options := &text.DrawOptions{}
	options.LineSpacing = p.face.Size

	// Draw the title
	options.ColorScale.ScaleWithColor(color.RGBA{R: 255, G: 255, B: 200, A: 255})
	options.GeoM.Translate(float64(x+margin), float64(y+margin))
	text.Draw(screen, p.title, p.face, options)

	// Draw the content
	options.ColorScale.Reset()
	text.Draw(screen, "\n"+p.content, p.face, options)
}

func (p *Panel) Layout() (width, height int) {
	return p.width, p.height
}
