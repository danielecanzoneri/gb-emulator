package ppu

import "testing"

var (
	TestTileData = [16]uint8{
		0b00111100, 0b01111110,
		0b01000010, 0b01000010,
		0b01000010, 0b01000010,
		0b01000010, 0b01000010,
		0b01111110, 0b01011110,
		0b01111110, 0b00001010,
		0b01111100, 0b01010110,
		0b00111000, 0b01111100,
	}
	TestExpectedTile = [8][8]uint8{
		{0b00, 0b10, 0b11, 0b11, 0b11, 0b11, 0b10, 0b00},
		{0b00, 0b11, 0b00, 0b00, 0b00, 0b00, 0b11, 0b00},
		{0b00, 0b11, 0b00, 0b00, 0b00, 0b00, 0b11, 0b00},
		{0b00, 0b11, 0b00, 0b00, 0b00, 0b00, 0b11, 0b00},
		{0b00, 0b11, 0b01, 0b11, 0b11, 0b11, 0b11, 0b00},
		{0b00, 0b01, 0b01, 0b01, 0b11, 0b01, 0b11, 0b00},
		{0b00, 0b11, 0b01, 0b11, 0b01, 0b11, 0b10, 0b00},
		{0b00, 0b10, 0b11, 0b11, 0b11, 0b10, 0b00, 0b00},
	}
)

func TestReadTile(t *testing.T) {
	tile := Tile(TestTileData)
	for i, expected := range TestExpectedTile {
		row := tile.getRowPixels(uint8(i))
		if row != expected {
			t.Errorf("tileRow[%d]: expected %v, got %v", i, expected, row)
		}
	}
}

func TestReadTileBGWindow(t *testing.T) {
	ppu := New()

	tileId := uint16(300)
	// Write tileData in vRAM
	copy(ppu.vRAM.tileData[tileId][:], TestTileData[:])

	ppu.bgWindowTileDataArea = 0
	tile := ppu.ReadTileBGWindow(uint8(tileId & 0xFF))
	for i, expected := range TestExpectedTile {
		row := tile.getRowPixels(uint8(i))
		if row != expected {
			t.Errorf("tileRow[%d]: expected %v, got %v", i, expected, row)
		}
	}
}
