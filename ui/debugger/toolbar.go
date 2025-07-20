package debugger

import (
	"fmt"
	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/event"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	goimage "image"
)

var (
	toolbarMenuImage = &widget.ButtonImage{
		Idle:    image.NewNineSliceColor(toolbarMenuColor),
		Hover:   image.NewNineSliceColor(toolbarMenuHoverColor),
		Pressed: image.NewNineSliceColor(toolbarMenuPressedColor),
	}
	toolbarEntryImage = &widget.ButtonImage{
		Idle:    image.NewNineSliceColor(toolbarEntryColor),
		Hover:   image.NewNineSliceColor(toolbarEntryHoverColor),
		Pressed: image.NewNineSliceColor(toolbarEntryPressedColor),
	}
	toolbarButtonTextColor = &widget.ButtonTextColor{
		Idle: toolbarTextColor,
	}
)

type toolbar struct {
	*widget.Container

	// Menu containing run options
	runMenu        *toolbarMenu
	stepButton     *toolbarMenuEntry
	continueButton *toolbarMenuEntry
	stopButton     *toolbarMenuEntry
	resetButton    *toolbarMenuEntry

	// Pointer to the UI for showing the window
	ui *ebitenui.UI
}

func (d *Debugger) newToolbar() *toolbar {
	t := &toolbar{ui: d.UI}
	t.Container = widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(toolbarBackgroundColor)),
		widget.ContainerOpts.Layout(widget.NewRowLayout()),
		widget.ContainerOpts.WidgetOpts(
			// Make the toolbar fill the whole horizontal space of the screen.
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{Stretch: true}),
		),
	)

	// Run menu
	t.runMenu = t.newMenu("Run")
	t.stepButton = t.runMenu.newEntryWithShortcut("Step", d.Step,
		ebiten.KeyF3)
	t.stepButton = t.runMenu.newEntryWithShortcut("Next", d.Next,
		ebiten.KeyF8)
	t.continueButton = t.runMenu.newEntryWithShortcut("Continue", d.Continue,
		ebiten.KeyF9)
	t.stopButton = t.runMenu.newEntryWithShortcut("Stop", d.Stop,
		ebiten.KeyShift, ebiten.KeyF9)
	t.resetButton = t.runMenu.newEntryWithShortcut("Reset", d.Reset,
		ebiten.KeyControl, ebiten.KeyR)
	t.AddChild(t.runMenu)
	return t
}

type toolbarMenu struct {
	*widget.Button

	title   string
	entries []*toolbarMenuEntry
}

type toolbarMenuEntry struct {
	*widget.Button

	title    string
	shortcut string
	onClick  func()
}

func (t *toolbar) newMenu(title string) *toolbarMenu {
	tMenu := &toolbarMenu{entries: make([]*toolbarMenuEntry, 0)}
	// Create a button for the toolbar.
	tMenu.Button = widget.NewButton(
		widget.ButtonOpts.Image(toolbarMenuImage),
		widget.ButtonOpts.Text(title, font, toolbarButtonTextColor),
		widget.ButtonOpts.TextPosition(widget.TextPositionStart, widget.TextPositionCenter), // Align text on the left
		widget.ButtonOpts.TextPadding(widget.NewInsetsSimple(padding)),
		widget.ButtonOpts.WidgetOpts(
			widget.WidgetOpts.MinSize(64, 0),
		),
	)

	tMenu.Button.ClickedEvent.AddHandler(event.WrapHandler(func(args *widget.ButtonClickedEventArgs) {
		t.openMenu(tMenu)
	}))
	return tMenu
}

func (t *toolbarMenu) createEntry(label string, shortcut string, onClick func()) *toolbarMenuEntry {
	entry := &toolbarMenuEntry{title: label, shortcut: shortcut, onClick: onClick}
	t.entries = append(t.entries, entry)

	// Create a button for a menu entry.
	entry.Button = widget.NewButton(
		widget.ButtonOpts.Image(toolbarEntryImage),
		widget.ButtonOpts.Text("", font, toolbarButtonTextColor),
		widget.ButtonOpts.TextPosition(widget.TextPositionStart, widget.TextPositionCenter), // Align text on the left
		widget.ButtonOpts.TextPadding(widget.Insets{Left: 16, Right: 64}),
		widget.ButtonOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{Stretch: true}),
		),
		// Handler
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			onClick()
		}),
	)
	return entry
}

func (t *toolbarMenu) newEntry(label string, onClick func()) *toolbarMenuEntry {
	defer t.relayout()
	return t.createEntry(label, "", onClick)
}

func (t *toolbarMenu) newEntryWithShortcut(label string, onClick func(), shortcutKeys ...ebiten.Key) *toolbarMenuEntry {
	defer t.relayout()

	// Shortcut mnemonic
	shortcut := shortcutKeys[0].String()
	for _, key := range shortcutKeys[1:] {
		shortcut = shortcut + "+" + key.String()
	}

	// Shortcut handler (all keys except last are modifier)
	shortcutHandler := func() {
		for _, key := range shortcutKeys[:len(shortcutKeys)-1] {
			if !ebiten.IsKeyPressed(key) {
				return
			}
		}
		d := inpututil.KeyPressDuration(shortcutKeys[len(shortcutKeys)-1])
		if d == 1 { // Key just pressed
			onClick()
		}
		if d > keyRepeatDelay && (d-keyRepeatDelay)%keyRepeatInterval == 0 { // Repeat if held down
			onClick()
		}
	}
	InputHandlers = append(InputHandlers, shortcutHandler)

	return t.createEntry(label, shortcut, onClick)
}

// relayout toolbar entries text
func (t *toolbarMenu) relayout() {
	maxLabelLength := 0
	for _, entry := range t.entries {
		maxLabelLength = max(maxLabelLength, len(entry.title))
	}

	// Reformat title length to align shortcuts
	for _, entry := range t.entries {
		titleFormatted := fmt.Sprintf("%-*s  %s", maxLabelLength, entry.title, entry.shortcut)
		entry.Button.Text().Label = titleFormatted
	}
}

func (t *toolbar) openMenu(menu *toolbarMenu) {
	c := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(toolbarMenuColor)),
		widget.ContainerOpts.Layout(
			widget.NewRowLayout(
				widget.RowLayoutOpts.Direction(widget.DirectionVertical),
				widget.RowLayoutOpts.Spacing(padding),
				widget.RowLayoutOpts.Padding(widget.Insets{Top: padding, Bottom: padding}),
			),
		),
		// Set minimum width for the menu.
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.MinSize(128, 0)),
	)

	for _, entry := range menu.entries {
		c.AddChild(entry)
	}

	w, h := c.PreferredSize()

	widgetRect := menu.GetWidget().Rect
	window := widget.NewWindow(
		widget.WindowOpts.Contents(c),
		widget.WindowOpts.CloseMode(widget.CLICK), // Close the menu if the user clicks outside of it.
		// Position the menu below the menu button that it belongs to.
		widget.WindowOpts.Location(
			goimage.Rect(
				widgetRect.Min.X,
				widgetRect.Min.Y+widgetRect.Max.Y,
				widgetRect.Min.X+w,
				widgetRect.Min.Y+widgetRect.Max.Y+widgetRect.Min.Y+h,
			),
		),
	)

	// Add the menu to the UI.
	t.ui.AddWindow(window)
}
