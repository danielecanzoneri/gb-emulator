package ppu

import (
	"testing"
)

func TestDrawLineObj(t *testing.T) {
	ppu := New()
	ppu.obj8x16Size = false
	ppu.bgWindowEnabled = false
	ppu.objEnabled = true
	ppu.LY = 10
	ppu.OBP[0] = DMGPalette(0b11100100)

	tile := Tile([16]uint8{
		0b00111010, 0b11001010, // Row 0
		0, 0,
		0b00111010, 0b11001010, // Row 2
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
	})
	ppu.vRAM.tileData[0][0] = tile

	// Mock objects
	ppu.objsLY[1] = &Object{ // This has higher priority
		x:         8,
		y:         26,
		flags:     0,
		tileIndex: 0, // Expected row 0
	}
	ppu.objsLY[0] = &Object{
		x:         14,
		y:         21,
		flags:     0x30, // Expect flipped row 2
		tileIndex: 0,
	}
	ppu.numObjs = 2

	// Pixel 0-7 should be of object at index 1 and pixel 8-13 should be of object at index 0
	// Pixel 7 is of object at index 0 because it is transparent for object 1
	expectedFrameLine := [160]uint16{}
	copy(expectedFrameLine[0:16], []uint16{0b10, 0b10, 0b01, 0b01, 0b11, 0b00, 0b11, 0b11, 0b00, 0b11, 0b01, 0b01, 0b10, 0b10})

	// Call renderLine
	ppu.renderLine()

	frameLine := ppu.backBuffer[ppu.LY]
	for x := 0; x < FrameWidth; x++ {
		if frameLine[x] != expectedFrameLine[x] {
			t.Errorf("pixel[%d]: got %02b, expected %02b", x, frameLine[x], expectedFrameLine[x])
		}
	}
}

func TestDrawLineBG(t *testing.T) {
	ppu := New()
	ppu.BGP = DMGPalette(0b11100100)
	ppu.windowEnabled = false
	ppu.bgWindowTileDataArea = 0
	ppu.bgTileMapAddr = 0x9C00
	ppu.objEnabled = false
	ppu.bgWindowEnabled = true
	ppu.LY = 12
	ppu.SCY = 254
	ppu.SCX = 205

	var tileId uint8 = 1
	var tileAddr uint16 = 0x9010 - vRAMStartAddr
	tileData := [16]uint8{
		0xFF, 0xEE, 0xDD, 0xCC,
		0b00111010, 0b11001010, // would be 0b1010010111001100
		0xBB, 0xAA, 0x99, 0x88, 0x77, 0x66, 0x55, 0x44, 0x33, 0x22,
	}
	copy(ppu.vRAM.tileData[0][tileAddr/16][:], tileData[:])

	// row 12 starting at 254 -> row 266 = 10 -> row 2 of 2nd row of tiles
	// ranging from x = 205 (x=5 of 26th tile) to x=205+159=364=108 (x=4 of 14th tile)
	ppu.vRAM.tileMaps[0][0x400+32+25] = tileId
	ppu.vRAM.tileMaps[0][0x400+32+13] = tileId

	tileRow := [8]uint16{0b10, 0b10, 0b01, 0b01, 0b11, 0b00, 0b11, 0b00}
	expectedFrameLine := [160]uint16{}
	copy(expectedFrameLine[0:3], tileRow[5:])
	copy(expectedFrameLine[155:160], tileRow[:5])

	// Call renderLine
	ppu.renderLine()

	frameLine := ppu.backBuffer[ppu.LY]
	for x := 0; x < FrameWidth; x++ {
		if frameLine[x] != expectedFrameLine[x] {
			t.Errorf("pixel[%d]: got %02b, expected %02b", x, frameLine[x], expectedFrameLine[x])
		}
	}
}

func TestDrawLineWindow(t *testing.T) {
	ppu := New()
	ppu.BGP = DMGPalette(0b11100100)
	ppu.windowTileMapAddr = 0x9800
	ppu.windowEnabled = true
	ppu.bgWindowTileDataArea = 0
	ppu.bgTileMapAddr = 0x9C00
	ppu.objEnabled = false
	ppu.bgWindowEnabled = true
	ppu.LY = 2
	ppu.SCY = 0
	ppu.SCX = 0
	ppu.WY = 2
	ppu.WX = 17

	var tileBGId uint8 = 1
	var tileBGAddr uint16 = 0x9010 - vRAMStartAddr
	tileBGData := [16]uint8{
		0xFF, 0xEE, 0xDD, 0xCC,
		0b00111010, 0b11001010, // would be 0b1010010111001100
		0xBB, 0xAA, 0x99, 0x88, 0x77, 0x66, 0x55, 0x44, 0x33, 0x22,
	}
	copy(ppu.vRAM.tileData[0][tileBGAddr/16][:], tileBGData[:])

	var tileWindowId uint8 = 2
	var tileWindowAddr uint16 = 0x9020 - vRAMStartAddr
	tileWindowData := [16]uint8{
		0b10100010, 0b11110110, // would be 11, 10, 11, 10, 00, 10, 11, 00
		0xFF, 0xEE, 0xDD, 0xCC,
		0xBB, 0xAA, 0x99, 0x88, 0x77, 0x66, 0x55, 0x44, 0x33, 0x22,
	}
	copy(ppu.vRAM.tileData[0][tileWindowAddr/16][:], tileWindowData[:])

	// Line 2 should take the first 10 pixels (WX = 17) of the BG and the others from the window
	ppu.vRAM.tileMaps[0][0x400] = tileBGId
	ppu.vRAM.tileMaps[0][0x400+1] = tileBGId
	ppu.vRAM.tileMaps[0][0] = tileWindowId

	bgTileRow := [8]uint16{0b10, 0b10, 0b01, 0b01, 0b11, 0b00, 0b11, 0b00}
	windowTileRow := [8]uint16{0b11, 0b10, 0b11, 0b10, 0b00, 0b10, 0b11, 0b00}
	expectedFrameLine := [160]uint16{}
	copy(expectedFrameLine[0:8], bgTileRow[:])
	copy(expectedFrameLine[8:10], bgTileRow[:2])
	copy(expectedFrameLine[10:18], windowTileRow[:])

	// Call renderLine
	ppu.renderLine()

	frameLine := ppu.backBuffer[ppu.LY]
	for x := 0; x < FrameWidth; x++ {
		if frameLine[x] != expectedFrameLine[x] {
			t.Errorf("pixel[%d]: got %02b, expected %02b", x, frameLine[x], expectedFrameLine[x])
		}
	}
}
