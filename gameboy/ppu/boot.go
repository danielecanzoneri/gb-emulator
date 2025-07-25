package ppu

func (ppu *PPU) SkipBoot() {
	ppu.Dots = 400
	ppu.InternalState = new(mode1)
	ppu.interruptMode = 1
	ppu.InternalStateLength = 56

	// vRAM
	ppu.vRAM.tileData[1] = [16]uint8{0xF0, 0, 0xF0, 0, 0xFC, 0, 0xFC, 0, 0xFC, 0, 0xFC, 0, 0xF3, 0, 0xF3, 0}
	ppu.vRAM.tileData[2] = [16]uint8{0x3C, 0, 0x3C, 0, 0x3C, 0, 0x3C, 0, 0x3C, 0, 0x3C, 0, 0x3C, 0, 0x3C, 0}
	ppu.vRAM.tileData[3] = [16]uint8{0xF0, 0, 0xF0, 0, 0xF0, 0, 0xF0, 0, 0, 0, 0, 0, 0xF3, 0, 0xF3, 0}
	ppu.vRAM.tileData[4] = [16]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0xCF, 0, 0xCF, 0}
	ppu.vRAM.tileData[5] = [16]uint8{0, 0, 0, 0, 0x0F, 0, 0x0F, 0, 0x3F, 0, 0x3F, 0, 0x0F, 0, 0x0F, 0}
	ppu.vRAM.tileData[6] = [16]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0xC0, 0, 0xC0, 0, 0x0F, 0, 0x0F, 0}
	ppu.vRAM.tileData[7] = [16]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0xF0, 0, 0xF0, 0}
	ppu.vRAM.tileData[8] = [16]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0xF3, 0, 0xF3, 0}
	ppu.vRAM.tileData[9] = [16]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0xC0, 0, 0xC0, 0}
	ppu.vRAM.tileData[10] = [16]uint8{0x03, 0, 0x03, 0, 0x03, 0, 0x03, 0, 0x03, 0, 0x03, 0, 0xFF, 0, 0xFF, 0}
	ppu.vRAM.tileData[11] = [16]uint8{0xC0, 0, 0xC0, 0, 0xC0, 0, 0xC0, 0, 0xC0, 0, 0xC0, 0, 0xC3, 0, 0xC3, 0}
	ppu.vRAM.tileData[12] = [16]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0xFC, 0, 0xFC, 0}
	ppu.vRAM.tileData[13] = [16]uint8{0xF3, 0, 0xF3, 0, 0xF0, 0, 0xF0, 0, 0xF0, 0, 0xF0, 0, 0xF0, 0, 0xF0, 0}
	ppu.vRAM.tileData[14] = [16]uint8{0x3C, 0, 0x3C, 0, 0xFC, 0, 0xFC, 0, 0xFC, 0, 0xFC, 0, 0x3C, 0, 0x3C, 0}
	ppu.vRAM.tileData[15] = [16]uint8{0xF3, 0, 0xF3, 0, 0xF3, 0, 0xF3, 0, 0xF3, 0, 0xF3, 0, 0xF3, 0, 0xF3, 0}
	ppu.vRAM.tileData[16] = [16]uint8{0xF3, 0, 0xF3, 0, 0xC3, 0, 0xC3, 0, 0xC3, 0, 0xC3, 0, 0xC3, 0, 0xC3, 0}
	ppu.vRAM.tileData[17] = [16]uint8{0xCF, 0, 0xCF, 0, 0xCF, 0, 0xCF, 0, 0xCF, 0, 0xCF, 0, 0xCF, 0, 0xCF, 0}
	ppu.vRAM.tileData[18] = [16]uint8{0x3C, 0, 0x3C, 0, 0x3F, 0, 0x3F, 0, 0x3C, 0, 0x3C, 0, 0x0F, 0, 0x0F, 0}
	ppu.vRAM.tileData[19] = [16]uint8{0x3C, 0, 0x3C, 0, 0xFC, 0, 0xFC, 0, 0, 0, 0, 0, 0xFC, 0, 0xFC, 0}
	ppu.vRAM.tileData[20] = [16]uint8{0xFC, 0, 0xFC, 0, 0xF0, 0, 0xF0, 0, 0xF0, 0, 0xF0, 0, 0xF0, 0, 0xF0, 0}
	ppu.vRAM.tileData[21] = [16]uint8{0xF3, 0, 0xF3, 0, 0xF3, 0, 0xF3, 0, 0xF3, 0, 0xF3, 0, 0xF0, 0, 0xF0, 0}
	ppu.vRAM.tileData[22] = [16]uint8{0xC3, 0, 0xC3, 0, 0xC3, 0, 0xC3, 0, 0xC3, 0, 0xC3, 0, 0xFF, 0, 0xFF, 0}
	ppu.vRAM.tileData[23] = [16]uint8{0xCF, 0, 0xCF, 0, 0xCF, 0, 0xCF, 0, 0xCF, 0, 0xCF, 0, 0xC3, 0, 0xC3, 0}
	ppu.vRAM.tileData[24] = [16]uint8{0x0F, 0, 0x0F, 0, 0x0F, 0, 0x0F, 0, 0x0F, 0, 0x0F, 0, 0xFC, 0, 0xFC, 0}
	ppu.vRAM.tileData[25] = [16]uint8{0x3C, 0, 0x42, 0, 0xB9, 0, 0xA5, 0, 0xB9, 0, 0xA5, 0, 0x42, 0, 0x3C, 0}
	for i := range 13 { // 260: 1 ... 271: 12
		ppu.vRAM.tileMaps[0x103+i] = uint8(i)
	}
	ppu.vRAM.tileMaps[0x110] = 25
	for i := range 12 { // 292: 13 ... 303: 24
		ppu.vRAM.tileMaps[0x124+i] = uint8(i + 13)
	}

	ppu.LCDC = 0x91
	ppu.STAT = 0x81
	ppu.active = true
	ppu.LY = 0
	ppu.BGP = 0xFC
	ppu.windowTileMapAddr = 0x9800
	ppu.bgWindowTileDataArea = 1
	ppu.bgTileMapAddr = 0x9800
	ppu.bgWindowEnabled = true
}
