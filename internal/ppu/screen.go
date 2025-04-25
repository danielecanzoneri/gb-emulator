package ppu

const (
	frameWidth  = 160
	frameHeight = 144
)

func (ppu *PPU) drawLine() {
	// Array that contains current line objects pixels
	var frameLine = [frameWidth]uint8{}

	if ppu.bgWindowEnabled {
		yWindow := ppu.LY - ppu.WY
		yBackground := ppu.SCY + ppu.LY
		for x := 0; x < frameWidth; x++ {
			// TODO - avoid if inside for
			if ppu.LY >= ppu.WY && uint8(x)+7 >= ppu.WX {
				// We're drawing the window
				xWindow := uint8(x) + 7 - ppu.WX
				tileAddr := ppu.windowTileMapAddr + getTileMapOffset(xWindow, yWindow)
				tileId := ppu.ReadVRAM(tileAddr)

				tile := ppu.ReadTileBGWindow(tileId)
				pixels := tile.getRowPixels(yWindow & 0b111)
				frameLine[x] = ppu.BGP.getColor(pixels[x&0b111])
			} else {
				// We're drawing the background
				xBackground := ppu.SCX + uint8(x) // Auto wrap around
				tileAddr := ppu.bgTileMapAddr + getTileMapOffset(xBackground, yBackground)
				tileId := ppu.ReadVRAM(tileAddr)

				// TODO - obviously optimize
				tile := ppu.ReadTileBGWindow(tileId)
				pixels := tile.getRowPixels(yBackground & 0b111)
				frameLine[x] = ppu.BGP.getColor(pixels[x&0b111])
			}
		}
	}

	if ppu.objEnabled {
		var pixelLine = [frameWidth]uint8{}
		// Draw objects with priority
		for _, obj := range ppu.objsLY {
			// Object row to draw is: obj.y - 16 + LY
			rowPixels := obj.getRow(ppu.LY - yOffset + obj.y)

			// Draw pixels on the line if no other pixel with higher priority was drawn
			for i, px := range rowPixels {
				x := uint8(i) + obj.x
				if xOffset <= x && x < frameWidth+xOffset {
					if pixelLine[x-xOffset] == 0 {
						pixelLine[x-xOffset] = obj.palette.getColor(px)
					}
				}
			}
		}

		// Write objects pixel on the line
		for i, b := range pixelLine {
			if b > 0 {
				frameLine[i] = b
			}
		}
	}

	// Set current line
	ppu.framebuffer[ppu.LY] = frameLine
}

func getTileMapOffset(x, y uint8) uint16 {
	return (uint16(x) >> 3) | ((uint16(y) & 0xF8) << 2)
}
