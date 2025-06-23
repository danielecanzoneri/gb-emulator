package ppu

const (
	vRAMSize uint16 = 0x2000

	// There are 384 tiles of size 16 byte in the range 8000-97FF
	tileNums        = 384
	tileSize uint16 = 16

	// Range 9800-9FFF contains 2 32x32 tile maps
	tileMapsSize = vRAMSize - tileNums*tileSize
)

type vRAM struct {
	tileData [tileNums]Tile
	tileMaps [tileMapsSize]uint8

	// Disabled during mode 3 (drawing)
	disabled bool
}

func (v *vRAM) Read(addr uint16) uint8 {
	if v.disabled {
		return 0xFF
	}
	return v.read(addr)
}

func (v *vRAM) read(addr uint16) uint8 {
	addr -= vRAMStartAddr

	if addr < tileNums*tileSize { // Tile data
		tileId := addr / tileSize
		tileOffset := addr % tileSize
		return v.tileData[tileId].read(tileOffset)
	}

	// Tile maps
	return v.tileMaps[addr-tileNums*tileSize]
}

func (v *vRAM) Write(addr uint16, value uint8) {
	if v.disabled {
		return
	}
	v.write(addr, value)
}

func (v *vRAM) write(addr uint16, value uint8) {
	addr -= vRAMStartAddr

	if addr < tileNums*tileSize { // Tile data
		tileId := addr / tileSize
		tileOffset := addr % tileSize
		v.tileData[tileId].write(tileOffset, value)
		return
	}

	// Tile maps
	v.tileMaps[addr-tileNums*tileSize] = value
}

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

func (ppu *PPU) ReadTileObj(tileId uint8) *Tile {
	return &ppu.vRAM.tileData[tileId]
}

func (ppu *PPU) ReadTileBGWindow(tileId uint8) *Tile {
	tileNum := uint16(tileId)

	// In this case tileId is a signed int8 with starting address 0x9000
	if ppu.bgWindowTileDataArea == 0 {
		if tileNum < 128 {
			tileNum += 256
		}
	}

	return &ppu.vRAM.tileData[tileNum]
}
