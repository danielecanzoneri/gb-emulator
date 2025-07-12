package joypad

type Joypad struct {
	selectButtons uint8
	selectDPad    uint8

	startDown uint8
	selectUp  uint8
	bLeft     uint8
	aRight    uint8
}

func New() *Joypad {
	return &Joypad{startDown: 1, selectUp: 1, bLeft: 1, aRight: 1}
}

func (jp *Joypad) Write(v uint8) {
	jp.selectButtons = (v & 0x20) >> 5
	jp.selectDPad = (v & 0x10) >> 4
}

func (jp *Joypad) Read() uint8 {
	return 0xC0 | (jp.selectButtons << 5) | (jp.selectDPad << 4) |
		(jp.startDown << 3) | (jp.selectUp << 2) | (jp.bLeft << 1) | jp.aRight
}
