package ppu

type tileAttr uint8

func (attr tileAttr) bgPriority() bool {
	return (attr & 0x80) > 0
}

func (attr tileAttr) yFlip() bool {
	return (attr & 0x40) > 0
}

func (attr tileAttr) xFlip() bool {
	return (attr & 0x20) > 0
}

func (attr tileAttr) dmgPalette() uint8 {
	return uint8(attr&0x10) >> 4
}

func (attr tileAttr) bank() uint8 {
	return uint8(attr&0x08) >> 3
}

func (attr tileAttr) cgbPalette() uint8 {
	return uint8(attr & 0x07)
}
