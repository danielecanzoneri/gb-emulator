package debugger

import (
	"fmt"
	goimage "image"

	"github.com/danielecanzoneri/gb-emulator/ui/theme"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/event"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

var (
	toolbarMenuImage = &widget.ButtonImage{
		Idle:    image.NewNineSliceColor(theme.Debugger.Toolbar.MenuColor),
		Hover:   image.NewNineSliceColor(theme.Debugger.Toolbar.MenuHoverColor),
		Pressed: image.NewNineSliceColor(theme.Debugger.Toolbar.MenuPressedColor),
	}
	toolbarEntryImage = &widget.ButtonImage{
		Idle:    image.NewNineSliceColor(theme.Debugger.Toolbar.EntryColor),
		Hover:   image.NewNineSliceColor(theme.Debugger.Toolbar.EntryHoverColor),
		Pressed: image.NewNineSliceColor(theme.Debugger.Toolbar.EntryPressedColor),
	}
	toolbarButtonTextColor = &widget.ButtonTextColor{
		Idle: theme.Debugger.Toolbar.TextColor,
	}
)

type toolbar struct {
	*widget.Container

	// Pointer to the UI for showing the window
	ui *ebitenui.UI
}

func (d *Debugger) newToolbar() *toolbar {
	t := &toolbar{ui: d.UI}
	t.Container = widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(theme.Debugger.Toolbar.BackgroundColor)),
		widget.ContainerOpts.Layout(widget.NewRowLayout()),
		widget.ContainerOpts.WidgetOpts(
			// Make the toolbar fill the whole horizontal space of the screen.
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{Stretch: true}),
		),
	)

	// Run menu
	runMenu := t.newMenu("Run")
	runMenu.addEntryWithShortcut("Step", d.Step,
		ebiten.KeyF3)
	runMenu.addEntryWithShortcut("Next", d.Next,
		ebiten.KeyF8)
	runMenu.addEntryWithShortcut("Continue", d.Continue,
		ebiten.KeyF9)
	runMenu.addEntryWithShortcut("Stop", d.Stop,
		ebiten.KeyShift, ebiten.KeyF9)
	runMenu.addEntryWithShortcut("Reset", d.Reset,
		ebiten.KeyControl, ebiten.KeyR)

	// PPU menu
	ppuMenu := t.newMenu("PPU")
	ppuMenu.addEntryWithShortcut("OAM Viewer", func() { d.showWindow(d.oamViewer) },
		ebiten.KeyShift, ebiten.KeyO)
	ppuMenu.addEntryWithShortcut("BG Viewer", func() { d.showWindow(d.bgViewer) },
		ebiten.KeyShift, ebiten.KeyB)
	ppuMenu.addEntryWithShortcut("TilesViewer", func() { d.showWindow(d.tilesViewer) },
		ebiten.KeyShift, ebiten.KeyT)
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
		widget.ButtonOpts.TextPadding(widget.NewInsetsSimple(theme.Debugger.Padding)),
		widget.ButtonOpts.WidgetOpts(
			widget.WidgetOpts.MinSize(64, 0),
		),
	)

	tMenu.Button.ClickedEvent.AddHandler(event.WrapHandler(func(args *widget.ButtonClickedEventArgs) {
		t.openMenu(tMenu)
	}))
	t.AddChild(tMenu)

	return tMenu
}

func (t *toolbarMenu) createEntry(label string, shortcut string, onClick func()) {
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
}

func (t *toolbarMenu) addEntry(label string, onClick func()) {
	defer t.relayout()
	t.createEntry(label, "", onClick)
}

func (t *toolbarMenu) addEntryWithShortcut(label string, onClick func(), shortcutKeys ...ebiten.Key) {
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

	t.createEntry(label, shortcut, onClick)
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
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(theme.Debugger.Toolbar.MenuColor)),
		widget.ContainerOpts.Layout(
			widget.NewRowLayout(
				widget.RowLayoutOpts.Direction(widget.DirectionVertical),
				widget.RowLayoutOpts.Spacing(theme.Debugger.Padding),
				widget.RowLayoutOpts.Padding(theme.Debugger.Insets),
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
