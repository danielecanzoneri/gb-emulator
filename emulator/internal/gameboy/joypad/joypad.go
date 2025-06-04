package joypad

type Joypad struct {
	selectButtons uint8
	selectDPad    uint8

	startDown uint8
	selectUp  uint8
	bLeft     uint8
	aRight    uint8
}

func (j *Joypad) Reset() {
	j.selectButtons = 0
	j.selectDPad = 0
	j.startDown = 0
	j.selectUp = 0
	j.bLeft = 0
	j.aRight = 0
}

func New() *Joypad {
	return &Joypad{startDown: 1, selectUp: 1, bLeft: 1, aRight: 1}
}

func (j *Joypad) Write(v uint8) {
	j.selectButtons = (v & 0x20) >> 5
	j.selectDPad = (v & 0x10) >> 4
}

func (j *Joypad) Read() uint8 {
	return 0xC0 | (j.selectButtons << 5) | (j.selectDPad << 4) |
		(j.startDown << 3) | (j.selectUp << 2) | (j.bLeft << 1) | j.aRight
}
