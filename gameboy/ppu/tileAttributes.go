package ppu

type TileAttribute uint8

func (attr TileAttribute) BGPriority() bool {
	return (attr & 0x80) > 0
}

func (attr TileAttribute) YFlip() bool {
	return (attr & 0x40) > 0
}

func (attr TileAttribute) XFlip() bool {
	return (attr & 0x20) > 0
}

func (attr TileAttribute) DMGPalette() uint8 {
	return uint8(attr&0x10) >> 4
}

func (attr TileAttribute) Bank() uint8 {
	return uint8(attr&0x08) >> 3
}

func (attr TileAttribute) CGBPalette() uint8 {
	return uint8(attr & 0x07)
}
