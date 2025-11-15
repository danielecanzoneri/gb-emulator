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
	// CGB register
	bankNumber uint8

	// One for each bank
	tileData [2][tileNums]Tile
	tileMaps [2][tileMapsSize]uint8

	// Disabled during mode 3 (drawing)
	readDisabled  bool
	writeDisabled bool
}

func (v *vRAM) Read(addr uint16) uint8 {
	if v.readDisabled {
		return 0xFF
	}
	return v.read(addr)
}

func (v *vRAM) read(addr uint16) uint8 {
	addr -= vRAMStartAddr

	if addr < tileNums*tileSize { // Tile data
		tileId := addr / tileSize
		tileOffset := addr % tileSize
		return v.tileData[v.bankNumber][tileId].read(tileOffset)
	}

	// Tile maps
	return v.tileMaps[v.bankNumber][addr-tileNums*tileSize]
}

func (v *vRAM) Write(addr uint16, value uint8) {
	if v.writeDisabled {
		return
	}
	v.write(addr, value)
}

func (v *vRAM) write(addr uint16, value uint8) {
	addr -= vRAMStartAddr

	if addr < tileNums*tileSize { // Tile data
		tileId := addr / tileSize
		tileOffset := addr % tileSize
		v.tileData[v.bankNumber][tileId].write(tileOffset, value)
		return
	}

	// Tile maps
	v.tileMaps[v.bankNumber][addr-tileNums*tileSize] = value
}

func (ppu *PPU) VDMAWrite(index uint16, value uint8) {
	ppu.vRAM.Write(vRAMStartAddr+index, value)
}
