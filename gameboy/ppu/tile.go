package ppu

// Tile contains 8 pixel (2 bit per pixel) per row (8 rows)
type Tile struct {
	raw [tileSize]uint8

	// Cached tile pixels
	pixels [8][8]uint8
}

func (t *Tile) read(offset uint16) uint8 {
	return t.raw[offset]
}

func (t *Tile) write(offset uint16, value uint8) {
	t.raw[offset] = value
	t.updatePixels()
}

func (t *Tile) updatePixels() {
	// For each row, first byte specifies the least significant bit of the color ID of each pixel,
	// second byte specifies the most significant bit.
	for row := range t.pixels {
		// lsb                  |   msb
		t.pixels[row][0] = ((t.raw[2*row] >> 7) & 0b1) | ((t.raw[2*row+1] >> 6) & 0b10)
		t.pixels[row][1] = ((t.raw[2*row] >> 6) & 0b1) | ((t.raw[2*row+1] >> 5) & 0b10)
		t.pixels[row][2] = ((t.raw[2*row] >> 5) & 0b1) | ((t.raw[2*row+1] >> 4) & 0b10)
		t.pixels[row][3] = ((t.raw[2*row] >> 4) & 0b1) | ((t.raw[2*row+1] >> 3) & 0b10)
		t.pixels[row][4] = ((t.raw[2*row] >> 3) & 0b1) | ((t.raw[2*row+1] >> 2) & 0b10)
		t.pixels[row][5] = ((t.raw[2*row] >> 2) & 0b1) | ((t.raw[2*row+1] >> 1) & 0b10)
		t.pixels[row][6] = ((t.raw[2*row] >> 1) & 0b1) | ((t.raw[2*row+1] >> 0) & 0b10)
		t.pixels[row][7] = ((t.raw[2*row] >> 0) & 0b1) | ((t.raw[2*row+1] << 1) & 0b10)
	}
}

func (t *Tile) GetRow(attr TileAttribute, row uint8) [8]uint8 {
	if attr.YFlip() {
		row = 7 - row
	}

	pixels := t.pixels[row]

	if attr.XFlip() {
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

func (ppu *PPU) GetTileId(tileAddress uint16) uint8 {
	return ppu.vRAM.tileMaps[0][tileAddress]
}

func (ppu *PPU) GetBGWindowPixelRow(tileAddr uint16, tileY uint8) ([8]uint8, TileAttribute) {
	// Address in the tilemap of the tile
	tileMapAddr := tileAddr - 0x9800

	tileId := ppu.GetTileId(tileMapAddr)
	var tile *Tile

	// In CGB Mode, an additional map of 32×32 bytes is stored in VRAM Bank 1
	// (each byte defines attributes for the corresponding tile-number map entry in VRAM Bank 0)
	var attr = TileAttribute(0)
	if ppu.Cgb && !ppu.DmgCompatibility {
		attr = TileAttribute(ppu.vRAM.tileMaps[1][tileMapAddr])
	}

	tile = ppu.ReadTileBGWindow(tileId, attr.Bank())
	return tile.GetRow(attr, tileY&0b111), attr
}
