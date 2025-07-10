package debugger

import (
	"math"
)

// InputHandlers stores a collection of input handler functions that are called
// during each frame update to process keyboard shortcuts for debugger controls.
// These handlers are registered by toolbar menu entries
// and are executed by the main UI's input handling loop to provide responsive
// keyboard shortcuts for debugger operations.
var InputHandlers = []func(){}

func computeRowsToScroll(variation float64, maxScroll int) int {
	direction := 1
	if variation < 0 {
		direction = -1
	}

	var amount int
	switch int(math.Abs(variation)) {
	case 0:
		amount = 0
	case 1:
		amount = 1
	case 2:
		amount = 4
	case 3:
		amount = 32
	case 4:
		amount = 128
	case 5:
		amount = 512
	default: // >= 6
		amount = maxScroll
	}
	return -amount * direction
}
