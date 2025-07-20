package debugger

import (
	"math"
)

const (
	keyRepeatDelay    = 24 // 400 millisecond delay
	keyRepeatInterval = 4  // Repeat 15 times a second
)

// InputHandlers stores a collection of input handler functions that are called
// during each frame update to process keyboard shortcuts for debugger controls.
// These handlers are registered by toolbar menu entries
// and are executed by the main UI's input handling loop to provide responsive
// keyboard shortcuts for debugger operations.
var InputHandlers = []func(){}

func computeRowsToScroll(variation float64) int {
	direction := 1
	if variation < 0 {
		direction = -1
	}

	var amount int
	if v := int(math.Abs(variation)); v > 0 {
		amount = 1 << (v - 1)
	}
	return -amount * direction
}
