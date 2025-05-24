package debugger

import "github.com/hajimehoshi/ebiten/v2"

// Layout represents how panels are arranged in a container
type Layout int

type Component interface {
	Draw(image *ebiten.Image, x, y int)
	Update()
	Layout() (width, height int)
}

const (
	Horizontal Layout = iota
	Vertical
)

// Container represents a group of panels arranged horizontally or vertically
type Container struct {
	components    []Component
	width, height int
	layout        Layout
}

// NewContainer creates a new container at the specified position with given layout
func NewContainer(layout Layout) *Container {
	return &Container{
		components: make([]Component, 0),
		layout:     layout,
	}
}

// AddComponent adds a component to the container and updates container dimensions
func (c *Container) AddComponent(component Component) *Container {
	c.components = append(c.components, component)

	// Update container dimensions
	width, height := component.Layout()
	if c.layout == Horizontal {
		c.width += width
		c.height = max(c.height, height)
	} else { // Vertical
		c.width = max(c.width, width)
		c.height += height
	}
	return c
}

func (c *Container) Update() {
	for _, component := range c.components {
		component.Update()
	}
}

// Draw renders all components in the container
func (c *Container) Draw(image *ebiten.Image, x, y int) {
	for _, component := range c.components {
		component.Draw(image, x, y)
		xComponent, yComponent := component.Layout()

		// Update the position of the next component
		if c.layout == Horizontal {
			x += xComponent
		} else { // Vertical
			y += yComponent
		}
	}
}

func (c *Container) Layout() (width, height int) {
	return c.width, c.height
}
