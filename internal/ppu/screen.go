package ppu

const (
	FrameWidth  = 160
	FrameHeight = 144
)

func (ppu *PPU) drawLine() {
	// Array that contains current line objects pixels
	var frameLine = [FrameWidth]uint8{}

	if ppu.bgWindowEnabled {
		yWindow := ppu.LY - ppu.WY
		yBackground := ppu.SCY + ppu.LY
		for x := 0; x < FrameWidth; x++ {
			// TODO - avoid if inside for
			if ppu.windowEnabled && ppu.LY >= ppu.WY && uint8(x)+7 >= ppu.WX {
				// We're drawing the window
				xWindow := uint8(x) + 7 - ppu.WX
				tileAddr := ppu.windowTileMapAddr + getTileMapOffset(xWindow, yWindow)
				tileId := ppu.ReadVRAM(tileAddr)

				tile := ppu.ReadTileBGWindow(tileId)
				pixels := tile.getRowPixels(yWindow & 0b111)
				frameLine[x] = ppu.BGP.getColor(pixels[xWindow&0b111])
			} else {
				// We're drawing the background
				xBackground := ppu.SCX + uint8(x) // Auto wrap around
				tileAddr := ppu.bgTileMapAddr + getTileMapOffset(xBackground, yBackground)
				tileId := ppu.ReadVRAM(tileAddr)

				// TODO - obviously optimize
				tile := ppu.ReadTileBGWindow(tileId)
				pixels := tile.getRowPixels(yBackground & 0b111)
				frameLine[x] = ppu.BGP.getColor(pixels[xBackground&0b111])
			}
		}
	}

	if ppu.objEnabled {
		var pixelLine = [FrameWidth]uint8{}
		var pixelBGPriority = [FrameWidth]bool{} // Pixel priority for BG/Window over obj

		// Draw objects with priority
		for i := range ppu.numObjs {
			obj := ppu.objsLY[i]
			// Object row to draw is: LY + 16 - y
			rowPixels := obj.getRow(ppu.LY + yOffset - obj.y)

			// Draw pixels on the line if no other pixel with higher priority was drawn
			for i, px := range rowPixels {
				x := uint8(i) + obj.x
				if xOffset <= x && x < FrameWidth+xOffset {
					if pixelLine[x-xOffset] == 0 {
						pixelLine[x-xOffset] = obj.palette.getColor(px)
						pixelBGPriority[x-xOffset] = obj.bgPriority
					}
				}
			}
		}

		// Write objects pixel on the line
		for i, b := range pixelLine {
			// Draw if pixel is not transparent and if no BG pixel has higher priority
			if b > 0 && (frameLine[i] == 0 || !pixelBGPriority[i]) {
				frameLine[i] = b
			}
		}
	}

	// Set current line
	ppu.Framebuffer[ppu.LY] = frameLine
}

func getTileMapOffset(x, y uint8) uint16 {
	return (uint16(x) >> 3) | ((uint16(y) & 0xF8) << 2)
}
