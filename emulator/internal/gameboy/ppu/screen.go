package ppu

const (
	FrameWidth  = 160
	FrameHeight = 144
)

type Palette uint8

func (p Palette) getColor(id uint8) uint8 {
	var mask uint8 = 0b11
	id &= mask

	// Bit 7,6: id 3; Bit 5,4: id 2; Bit 3,2: id 1; Bit 1,0: id 0
	shift := id * 2
	return (uint8(p) >> shift) & mask
}

func (ppu *PPU) GetFrame() *[FrameHeight][FrameWidth]uint8 {
	return ppu.frontBuffer
}

func (ppu *PPU) emptyFrame() {
	for x := range FrameWidth {
		for y := range FrameHeight {
			ppu.backBuffer[y][x] = 0
		}
	}

	// Swap buffers
	ppu.frontBuffer = ppu.backBuffer
	ppu.backBuffer = new([FrameHeight][FrameWidth]uint8)
}

// drawLine returns the number of penalty dots incurred to draw this line
func (ppu *PPU) drawLine() int {
	// SCX % 8 pixels are discarded from the leftmost tile
	penaltyDots := int(ppu.SCX % 8)
	// Keep track of the background/window tiles under each pixel
	tilesUnderPixels := [FrameWidth]uint16{}

	// Array that contains current line objects pixels
	var frameLine = [FrameWidth]uint8{}

	// Flag that is set to true when x+7 >= wx and used to increment window Y counter
	windowsRendered := false

	if ppu.bgWindowEnabled {
		yWindow := ppu.wyCounter

		yBackground := ppu.SCY + ppu.LY
		for x := 0; x < FrameWidth; x++ {
			if ppu.windowEnabled && ppu.LY >= ppu.WY && uint8(x)+7 >= ppu.WX {
				windowsRendered = true

				// We're drawing the window
				xWindow := uint8(x) + 7 - ppu.WX
				tileAddr := ppu.windowTileMapAddr + getTileMapOffset(xWindow, yWindow)
				tileId := ppu.vRAM.read(tileAddr)

				tilesUnderPixels[x] = tileAddr

				tile := ppu.ReadTileBGWindow(tileId)
				pixels := tile.getRowPixels(yWindow & 0b111)
				frameLine[x] = ppu.BGP.getColor(pixels[xWindow&0b111])
			} else {
				// We're drawing the background
				xBackground := ppu.SCX + uint8(x) // Auto wrap around
				tileAddr := ppu.bgTileMapAddr + getTileMapOffset(xBackground, yBackground)
				tileId := ppu.vRAM.read(tileAddr)

				tilesUnderPixels[x] = tileAddr

				// TODO - obviously optimize
				tile := ppu.ReadTileBGWindow(tileId)
				pixels := tile.getRowPixels(yBackground & 0b111)
				frameLine[x] = ppu.BGP.getColor(pixels[xBackground&0b111])
			}
		}
	}

	if windowsRendered {
		// 6-dot penalty is incurred while the BG fetcher is being set up for the window.
		penaltyDots += 6
		ppu.wyCounter++
	}

	if ppu.objEnabled {
		var pixelLine = [FrameWidth]uint8{}
		for x := 0; x < FrameWidth; x++ {
			pixelLine[x] = 0xFF // To not confuse with value 0 (0xFF means transparent)
		}
		var pixelBGPriority = [FrameWidth]bool{} // Pixel priority for BG/Window over obj

		// Indexes of the previous tile considered in the OBJ penalty algorithm
		var previousTile uint16 = 0

		// Draw objects with priority
		for i := range ppu.numObjs {
			obj := ppu.objsLY[i]

			// OBJ penalty algorithm
			// Only the OBJ’s leftmost pixel matters here, transparent or not; it is designated as “The Pixel” in the following.
			// TODO - Understand if "The Pixel" of objects with x < 8 is the leftmost pixel of the object or the leftmost pixel on the screen.
			// 1. Determine the tile (background or window) that The Pixel is within. (This is affected by horizontal scrolling and/or the window!)
			thePixel := 0
			if obj.x > 8 {
				thePixel = int(obj.x) - 8
			}

			if thePixel >= FrameWidth {
				continue
			}

			tile := tilesUnderPixels[thePixel]
			// 2. If that tile has not been considered by a previous OBJ yet:
			if tile != previousTile {
				pixelsOnTheRight := 0
				//    - Count how many of that tile’s pixels are strictly to the right of The Pixel.
				for px := thePixel + 1; px < min(thePixel+8, FrameWidth); px++ {
					if tilesUnderPixels[px] != tile {
						break
					}
					pixelsOnTheRight++
				}
				//    - Subtract 2.
				pixelsOnTheRight -= 2
				//    - Incur this many dots of penalty, or zero if negative (from waiting for the BG fetch to finish).
				if pixelsOnTheRight > 0 {
					penaltyDots += pixelsOnTheRight
				}
			}
			previousTile = tile
			// 3. Incur a flat, 6-dot penalty (from fetching the OBJ’s tile).
			penaltyDots += 6

			// Object row to draw is: LY + 16 - y
			rowPixels := ppu.getObjectRow(obj, ppu.LY+yOffset-obj.y)

			// Draw pixels on the line if no other pixel with higher priority was drawn
			for i, px := range rowPixels {
				// Pixel is transparent if color id is 0 (not if the color itself is 0)
				if px > 0 {
					x := uint8(i) + obj.x
					if xOffset <= x && x < FrameWidth+xOffset {
						if pixelLine[x-xOffset] == 0xFF {
							pixelLine[x-xOffset] = ppu.OBP[obj.paletteId].getColor(px)
							pixelBGPriority[x-xOffset] = obj.bgPriority
						}
					}
				}
			}
		}

		// Write objects pixel on the line
		for i, b := range pixelLine {
			// Draw if pixel is not transparent and if no BG pixel has higher priority
			if b != 0xFF && (frameLine[i] == 0 || !pixelBGPriority[i]) {
				frameLine[i] = b
			}
		}
	}

	// Set current line of the back buffer
	ppu.backBuffer[ppu.LY] = frameLine

	return penaltyDots
}

func getTileMapOffset(x, y uint8) uint16 {
	return (uint16(x) >> 3) | ((uint16(y) & 0xF8) << 2)
}
