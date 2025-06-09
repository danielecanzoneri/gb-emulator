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
	TestExpectedTile = [8]uint16{
		0b0010111111111000,
		0b0011000000001100,
		0b0011000000001100,
		0b0011000000001100,
		0b0011011111111100,
		0b0001010111011100,
		0b0011011101111000,
		0b0010111111100000,
	}
)

func TestReadTile(t *testing.T) {
	v := vRAM{}

	tileId := uint16(0)
	// Write tileData in vRAM
	for i := range tileSize {
		v.Data[tileId*tileSize+i] = TestTileData[i]
	}

	tile := v.readTile(tileId)
	for i, expected := range TestExpectedTile {
		if tile.data[i] != expected {
			t.Errorf("tileData[%d]: expected %08b, got %08b", i, expected, tile.data[i])
		}
	}
}

func TestReadTileBGWindow(t *testing.T) {
	ppu := New()

	tileId := uint16(300)
	// Write tileData in vRAM
	for i := range tileSize {
		ppu.vRAM.Data[tileId*tileSize+i] = TestTileData[i]
	}

	ppu.bgWindowTileDataArea = 0
	tile := ppu.ReadTileBGWindow(uint8(tileId & 0xFF))
	for i, expected := range TestExpectedTile {
		if tile.data[i] != expected {
			t.Errorf("tileData[%d]: expected %08b, got %08b", i, expected, tile.data[i])
		}
	}
}
