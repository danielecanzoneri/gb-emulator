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
	// SCX % 8 pixels are discarded from the leftmost background tile
	penaltyDotsBG := int(ppu.SCX % 8)
	penaltyDotsObj := 0

	// Flag that is set to true when x+7 >= wx and used to increment window Y counter
	windowsRendered := false

	// Tiles considered in the OBJ penalty algorithm (x ranges from 0 to 167+7, so we have at most 22 tiles
	var tileObjectsPenalties [(167+7)>>3 + 1]bool
	// Objects considered in the OBJ penalty algorithm (6 penalty dots per object)
	var objectsPenalties [10]bool

	for x := uint8(0); x < uint8(FrameWidth); x++ {
		if ppu.bgWindowEnabled {
			// TODO - obviously optimize
			var tile *Tile
			var tileX, tileY uint8
			var tileBaseAddr uint16

			if ppu.windowEnabled && ppu.LY >= ppu.WY && x+7 >= ppu.WX {
				windowsRendered = true

				// We're drawing the window
				tileX = x + 7 - ppu.WX
				tileY = ppu.wyCounter
				tileBaseAddr = ppu.windowTileMapAddr
			} else {
				// We're drawing the background
				tileX = ppu.SCX + x // Auto wrap around
				tileY = ppu.SCY + ppu.LY
				tileBaseAddr = ppu.bgTileMapAddr
			}

			tileAddr := tileBaseAddr + getTileMapOffset(tileX, tileY)
			tileId := ppu.vRAM.read(tileAddr)

			tile = ppu.ReadTileBGWindow(tileId)
			objPixels := tile.getRowPixels(tileY & 0b111)
			ppu.backBuffer[ppu.LY][x] = ppu.BGP.getColor(objPixels[tileX&0b111])
		}

		// Render objects
		if !ppu.objEnabled {
			continue
		}

		// Draw objects with priority
		for i := range ppu.numObjs {
			// Traverse objects in reverse so objects that come first will be displayed with priority
			obj := ppu.objsLY[ppu.numObjs-i-1]
			if obj.x >= 168 { // Out of range tile
				continue
			}

			// OBJ penalty algorithm (TODO - move out of loop)
			objX := obj.x + (ppu.SCX & 0b111)

			// Only the OBJ’s leftmost pixel matters here.
			// 1. Determine the tile (background or window) that the pixel is within. (This is affected by horizontal scrolling and/or the window!)
			tileId := objX >> 3

			// 2. If that tile has not been considered by a previous OBJ yet:
			if !tileObjectsPenalties[tileId] {
				tileObjectsPenalties[tileId] = true

				//    - Count how many of that tile’s pixels are strictly to the right of The Pixel.
				//    - Subtract 2.
				//    - Incur this many dots of penalty, or zero if negative (from waiting for the BG fetch to finish).
				penaltyDotsObj += max(5-int(objX&7), 0)
			}

			// 3. Incur a flat, 6-dot penalty (from fetching the OBJ’s tile).
			if !objectsPenalties[i] {
				objectsPenalties[i] = true
				penaltyDotsObj += 6
			}

			if !(obj.x <= x+8 && x+8 < obj.x+8) {
				// The object is not in this pixel
				continue
			}

			// Object row to draw is: LY + 16 - y
			rowPixels := ppu.GetObjectRow(obj, ppu.LY+yObjOffset-obj.y)

			// Draw pixel if no other pixel with higher priority was drawn
			px := rowPixels[(x+8)-obj.x]

			// Draw if pixel is not transparent and if no BG pixel has higher priority
			// (pixel is transparent if color id is 0 (not if the color itself is 0)
			if px > 0 && (ppu.backBuffer[ppu.LY][x] == 0 || !tileAttr(obj.flags).bgPriority()) {
				ppu.backBuffer[ppu.LY][x] = ppu.OBP[tileAttr(obj.flags).dmgPalette()].getColor(px)
			}
		}
	}

	if windowsRendered {
		ppu.wyCounter++
		penaltyDotsBG += 6 // 6-dot penalty is incurred while the BG fetcher is being set up for the window.
	}

	// Round objects penalty dots to M-cycle (TODO - investigate why it doesn't work otherwise)
	return penaltyDotsBG + (penaltyDotsObj & ^3)
}

func getTileMapOffset(x, y uint8) uint16 {
	return (uint16(x) >> 3) | ((uint16(y) & 0xF8) << 2)
}
