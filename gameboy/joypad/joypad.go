package joypad

import "github.com/danielecanzoneri/lucky-boy/util"

type Joypad struct {
	selectButtons uint8
	selectDPad    uint8

	startDown uint8
	selectUp  uint8
	bLeft     uint8
	aRight    uint8

	RequestInterrupt func()
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
