package joypad

import "github.com/danielecanzoneri/lucky-boy/util"

// InputProvider is an interface for detecting key presses.
// This allows the joypad package to be independent from any specific input library.
type InputProvider interface {
	IsKeyPressed(key Key) bool
}

// Key represents a Game Boy key
type Key uint8

const (
	KeyStart Key = iota
	KeySelect
	KeyB
	KeyA
	KeyDown
	KeyUp
	KeyLeft
	KeyRight
)

type Joypad struct {
	selectButtons uint8
	selectDPad    uint8

	startDown uint8
	selectUp  uint8
	bLeft     uint8
	aRight    uint8

	RequestInterrupt func()
	inputProvider    InputProvider
}

func New() *Joypad {
	return &Joypad{startDown: 1, selectUp: 1, bLeft: 1, aRight: 1}
}

func (jp *Joypad) Write(v uint8) {
	jp.selectButtons = util.ReadBit(v, 5)
	jp.selectDPad = util.ReadBit(v, 4)

	// If neither buttons nor d-pad is selected ($30 was written), then the low nibble reads $F (all buttons released)
	if v&0x30 > 0 {
		jp.startDown = 1
		jp.selectUp = 1
		jp.bLeft = 1
		jp.aRight = 1
	}
}

func (jp *Joypad) Read() uint8 {
	return 0xC0 | (jp.selectButtons << 5) | (jp.selectDPad << 4) |
		(jp.startDown << 3) | (jp.selectUp << 2) | (jp.bLeft << 1) | jp.aRight
}

// SetInputProvider sets the input provider for detecting key presses
func (jp *Joypad) SetInputProvider(provider InputProvider) {
	jp.inputProvider = provider
}

// DetectKeysPressed detects key presses using the configured InputProvider
func (jp *Joypad) DetectKeysPressed() {
	if jp.inputProvider == nil {
		return
	}

	if jp.selectButtons == 0 {
		if jp.inputProvider.IsKeyPressed(KeyStart) {
			if jp.startDown == 1 { // Detect high -> low transition
				jp.RequestInterrupt()
			}
			jp.startDown = 0
		}
		if jp.inputProvider.IsKeyPressed(KeySelect) {
			if jp.selectUp == 1 { // Detect high -> low transition
				jp.RequestInterrupt()
			}
			jp.selectUp = 0
		}
		if jp.inputProvider.IsKeyPressed(KeyB) {
			if jp.bLeft == 1 { // Detect high -> low transition
				jp.RequestInterrupt()
			}
			jp.bLeft = 0
		}
		if jp.inputProvider.IsKeyPressed(KeyA) {
			if jp.aRight == 1 { // Detect high -> low transition
				jp.RequestInterrupt()
			}
			jp.aRight = 0
		}
	}
	if jp.selectDPad == 0 {
		if jp.inputProvider.IsKeyPressed(KeyDown) {
			if jp.startDown == 1 { // Detect high -> low transition
				jp.RequestInterrupt()
			}
			jp.startDown = 0
		}
		if jp.inputProvider.IsKeyPressed(KeyUp) {
			if jp.selectUp == 1 { // Detect high -> low transition
				jp.RequestInterrupt()
			}
			jp.selectUp = 0
		}
		if jp.inputProvider.IsKeyPressed(KeyLeft) {
			if jp.bLeft == 1 { // Detect high -> low transition
				jp.RequestInterrupt()
			}
			jp.bLeft = 0
		}
		if jp.inputProvider.IsKeyPressed(KeyRight) {
			if jp.aRight == 1 { // Detect high -> low transition
				jp.RequestInterrupt()
			}
			jp.aRight = 0
		}
	}
}
