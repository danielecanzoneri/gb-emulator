package debugger

// InputHandlers stores a collection of input handler functions that are called
// during each frame update to process keyboard shortcuts for debugger controls.
// These handlers are registered by toolbar menu entries
// and are executed by the main UI's input handling loop to provide responsive
// keyboard shortcuts for debugger operations.
var InputHandlers = []func(){}
