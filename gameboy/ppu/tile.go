package ppu

// Tile contains 8 pixel (2 bit per pixel) per row (8 rows)
type Tile [tileSize]uint8

func (t *Tile) read(offset uint16) uint8 {
	return t[offset]
}

func (t *Tile) write(offset uint16, value uint8) {
	t[offset] = value
}

func (t *Tile) getRowPixels(row uint8) [8]uint8 {
	// First byte (2*row) specifies the least significant bit of the color ID of each pixel,
	// second byte (2*row+1) specifies the most significant bit.
	return [8]uint8{
		// lsb                  |   msb
		((t[2*row] >> 7) & 0b1) | ((t[2*row+1] >> 6) & 0b10),
		((t[2*row] >> 6) & 0b1) | ((t[2*row+1] >> 5) & 0b10),
		((t[2*row] >> 5) & 0b1) | ((t[2*row+1] >> 4) & 0b10),
		((t[2*row] >> 4) & 0b1) | ((t[2*row+1] >> 3) & 0b10),
		((t[2*row] >> 3) & 0b1) | ((t[2*row+1] >> 2) & 0b10),
		((t[2*row] >> 2) & 0b1) | ((t[2*row+1] >> 1) & 0b10),
		((t[2*row] >> 1) & 0b1) | ((t[2*row+1] >> 0) & 0b10),
		((t[2*row] >> 0) & 0b1) | ((t[2*row+1] << 1) & 0b10),
	}
}

func (t *Tile) GetRow(attr tileAttr, row uint8) [8]uint8 {
	if attr.yFlip() {
		row = 7 - row
	}

	pixels := t.getRowPixels(row)

	if attr.xFlip() {
		pixels[0], pixels[7] = pixels[7], pixels[0]
		pixels[1], pixels[6] = pixels[6], pixels[1]
		pixels[2], pixels[5] = pixels[5], pixels[2]
		pixels[3], pixels[4] = pixels[4], pixels[3]
	}
	return pixels
}

func (ppu *PPU) ReadTileObj(tileId uint8, vRAMBank uint8) *Tile {
	return &ppu.vRAM.tileData[vRAMBank][tileId]
}

func (ppu *PPU) ReadTileBGWindow(tileId uint8, vRAMBank uint8) *Tile {
	tileNum := uint16(tileId)

	// In this case tileId is a signed int8 with starting address 0x9000
	if ppu.bgWindowTileDataArea == 0 {
		if tileNum < 128 {
			tileNum += 256
		}
	}

	return &ppu.vRAM.tileData[vRAMBank][tileNum]
}

func (ppu *PPU) getBGWindowPixelRow(tileAddr uint16, tileY uint8) ([8]uint8, tileAttr) {
	// Address in the tilemap of the tile
	tileMapAddr := tileAddr - 0x9800

	tileId := ppu.vRAM.tileMaps[0][tileMapAddr]
	var tile *Tile

	// In CGB Mode, an additional map of 32×32 bytes is stored in VRAM Bank 1
	// (each byte defines attributes for the corresponding tile-number map entry in VRAM Bank 0)
	if ppu.cgb {
		attr := tileAttr(ppu.vRAM.tileMaps[1][tileMapAddr])
		tile = ppu.ReadTileBGWindow(tileId, attr.bank())
		return tile.GetRow(attr, tileY&0b111), attr
	} else {
		tile = ppu.ReadTileBGWindow(tileId, 0)
		return tile.getRowPixels(tileY & 0b111), 0
	}
}
