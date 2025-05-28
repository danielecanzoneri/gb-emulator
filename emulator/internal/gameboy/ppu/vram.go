package ppu

import (
	"github.com/danielecanzoneri/gb-emulator/pkg/util"
)

const (
	vRAMSize uint16 = 0x2000

	tileSize uint16 = 16
)

type vRAM struct {
	Data [vRAMSize]uint8
}

func (v *vRAM) read(addr uint16) uint8 {
	return v.Data[addr]
}

func (v *vRAM) write(addr uint16, value uint8) {
	v.Data[addr] = value
}

type Tile struct {
	// data contains 8 pixel (2 bit per pixel) per row
	data [8]uint16
}

func (t *Tile) getRowPixels(row uint8) [8]uint8 {
	return [8]uint8{
		uint8((t.data[row] >> 14) & 0b11),
		uint8((t.data[row] >> 12) & 0b11),
		uint8((t.data[row] >> 10) & 0b11),
		uint8((t.data[row] >> 8) & 0b11),
		uint8((t.data[row] >> 6) & 0b11),
		uint8((t.data[row] >> 4) & 0b11),
		uint8((t.data[row] >> 2) & 0b11),
		uint8(t.data[row] & 0b11),
	}
}

// readTile one of the 384 tiles from vRAM
func (v *vRAM) readTile(tileNum uint16) *Tile {
	tile := &Tile{}

	rawTileData := v.Data[tileSize*tileNum : tileSize*(tileNum+1)]
	for i := 0; i < int(tileSize); i += 2 {
		// First byte specifies the least significant bit of the color ID of each pixel,
		// second byte specifies the most significant bit.
		least := util.SpreadBits(rawTileData[i])
		most := util.SpreadBits(rawTileData[i+1]) << 1

		tile.data[i/2] = least | most
	}

	return tile
}

func (ppu *PPU) ReadTileObj(tileId uint8) *Tile {
	return ppu.vRAM.readTile(uint16(tileId))
}

func (ppu *PPU) ReadTileBGWindow(tileId uint8) *Tile {
	tileNum := uint16(tileId)

	// In this case tileId is a signed int8 with starting address 0x9000
	if ppu.bgWindowTileDataArea == 0 {
		if tileNum < 128 {
			tileNum += 256
		}
	}

	return ppu.vRAM.readTile(tileNum)
}
