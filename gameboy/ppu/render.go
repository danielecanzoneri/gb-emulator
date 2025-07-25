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

// renderLine returns the number of penalty dots incurred to draw this line
func (ppu *PPU) renderLine() int {
	d := 0
	d += ppu.renderBackground(ppu.backBuffer[ppu.LY][:])
	d += ppu.renderObjects(ppu.backBuffer[ppu.LY][:])
	return d
}

func (ppu *PPU) renderBackground(pixels []uint8) int {
	if !ppu.bgWindowEnabled {
		return 0
	}

	// SCX % 8 pixels are discarded from the leftmost tile
	penaltyDots := int(ppu.SCX % 8)

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

				tile := ppu.ReadTileBGWindow(tileId)
				objPixels := tile.getRowPixels(yWindow & 0b111)
				pixels[x] = ppu.BGP.getColor(objPixels[xWindow&0b111])
			} else {
				// We're drawing the background
				xBackground := ppu.SCX + uint8(x) // Auto wrap around
				tileAddr := ppu.bgTileMapAddr + getTileMapOffset(xBackground, yBackground)
				tileId := ppu.vRAM.read(tileAddr)

				// TODO - obviously optimize
				tile := ppu.ReadTileBGWindow(tileId)
				objPixels := tile.getRowPixels(yBackground & 0b111)
				pixels[x] = ppu.BGP.getColor(objPixels[xBackground&0b111])
			}
		}
	}

	if windowsRendered {
		ppu.wyCounter++
		penaltyDots += 6 // 6-dot penalty is incurred while the BG fetcher is being set up for the window.
	}
	return penaltyDots
}

func (ppu *PPU) renderObjects(pixels []uint8) int {
	if !ppu.objEnabled {
		return 0
	}

	var penaltyDots int

	var pixelLine = [FrameWidth]uint8{}
	for x := 0; x < FrameWidth; x++ {
		pixelLine[x] = 0xFF // To not confuse with value 0 (0xFF means transparent)
	}
	var pixelBGPriority = [FrameWidth]bool{} // Pixel priority for BG/Window over obj

	// Tiles considered in the OBJ penalty algorithm (x ranges from 0 to 167+7, so we have at most 22 tiles
	var tileObjectsPenalties [(167 + 7) >> 3]bool

	// Draw objects with priority
	for i := range ppu.numObjs {
		obj := ppu.objsLY[i]
		if obj.x >= 168 {
			continue
		}

		// OBJ penalty algorithm
		x := obj.x + (ppu.SCX & 0b111)

		// Only the OBJ’s leftmost pixel matters here.
		// 1. Determine the tile (background or window) that the pixel is within. (This is affected by horizontal scrolling and/or the window!)
		tileId := x >> 3

		// 2. If that tile has not been considered by a previous OBJ yet:
		if !tileObjectsPenalties[tileId] {
			tileObjectsPenalties[tileId] = true

			//    - Count how many of that tile’s pixels are strictly to the right of The Pixel.
			//    - Subtract 2.
			//    - Incur this many dots of penalty, or zero if negative (from waiting for the BG fetch to finish).
			penaltyDots += max(5-int(x&7), 0)
		}

		// 3. Incur a flat, 6-dot penalty (from fetching the OBJ’s tile).
		penaltyDots += 6

		// Object row to draw is: LY + 16 - y
		rowPixels := ppu.GetObjectRow(obj, ppu.LY+yObjOffset-obj.y)

		// Draw pixels on the line if no other pixel with higher priority was drawn
		for i, px := range rowPixels {
			// Pixel is transparent if color id is 0 (not if the color itself is 0)
			if px > 0 {
				x := uint8(i) + obj.x
				if xObjOffset <= x && x < FrameWidth+xObjOffset {
					if pixelLine[x-xObjOffset] == 0xFF {
						pixelLine[x-xObjOffset] = ppu.OBP[obj.paletteId].getColor(px)
						pixelBGPriority[x-xObjOffset] = obj.bgPriority
					}
				}
			}
		}
	}

	// Write objects pixel on the line
	for i, b := range pixelLine {
		// Draw if pixel is not transparent and if no BG pixel has higher priority
		if b != 0xFF && (pixels[i] == 0 || !pixelBGPriority[i]) {
			pixels[i] = b
		}
	}

	// Round to M-cycle (TODO - investigate why it doesn't work otherwise)
	return penaltyDots & ^3
}

func getTileMapOffset(x, y uint8) uint16 {
	return (uint16(x) >> 3) | ((uint16(y) & 0xF8) << 2)
}
