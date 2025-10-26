package ppu

const (
	FrameWidth  = 160
	FrameHeight = 144
)

type Palette interface {
	GetColor(uint8) uint16
}

type DMGPalette uint8

func (p DMGPalette) GetColor(id uint8) uint16 {
	var mask uint8 = 0b11
	id &= mask

	// Bit 7,6: id 3; Bit 5,4: id 2; Bit 3,2: id 1; Bit 1,0: id 0
	shift := id * 2
	return uint16((uint8(p) >> shift) & mask)
}

type CGBPalette []uint8

func (p CGBPalette) GetColor(id uint8) uint16 {
	// Each color is stored as little-endian RGB555
	return uint16(p[2*id]) | (uint16(p[2*id+1]) << 8)
}

func (ppu *PPU) GetFrame() *[FrameHeight][FrameWidth]uint16 {
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
	ppu.backBuffer = new([FrameHeight][FrameWidth]uint16)
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
		// CGB flag that signals if this pixel BG has higher priority than objs
		var bgPriority = false

		// In CGB mode the LCDC.0 has a different meaning, it is the BG/Window master priority
		if ppu.bgWindowEnabled || ppu.cgb {
			// TODO - obviously optimize
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
			bgPixels, tileAttributes := ppu.getBGWindowPixelRow(tileAddr, tileY)

			// Pixel color (0-3)
			color := bgPixels[tileX&0b111]
			var palette Palette = ppu.BGP // DMG palette

			if ppu.cgb {
				bgPriority = tileAttributes.bgPriority()
				paletteId := tileAttributes.cgbPalette()
				palette = CGBPalette(ppu.BGPalette[8*paletteId : 8*paletteId+8])
			}

			ppu.backBuffer[ppu.LY][x] = palette.GetColor(color)
		}

		// Render objects
		if !ppu.objEnabled {
			continue
		}

		// Draw objects with priority
		for i := range ppu.numObjs {
			obj := ppu.objsLY[i]
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

			if ppu.cgb {
				// If the BG color index is 0, the OBJ will always have priority;
				// Otherwise, if LCDC bit 0 is clear, the OBJ will always have priority;
				// Otherwise, if both the BG Attributes and the OAM Attributes have bit 7 clear, the OBJ will have priority;
				// Otherwise, BG will have priority.
				if px > 0 {
					// In CGB mode the LCDC.0 has a different meaning, it is the BG/Window master priority
					if ppu.backBuffer[ppu.LY][x] == 0 || !ppu.bgWindowEnabled || (!bgPriority && !tileAttr(obj.flags).bgPriority()) {
						paletteId := tileAttr(obj.flags).cgbPalette()
						palette := CGBPalette(ppu.OBJPalette[8*paletteId : 8*paletteId+8])

						ppu.backBuffer[ppu.LY][x] = palette.GetColor(px)
					}
				}
			} else {
				// If object pixel is transparent (px == 0), draw background pixel
				// Otherwise:
				//  - If background pixel is 0, draw object pixel
				//  - If both object and background pixel are not 0, draw pixel based on
				//    object attributes BG/Window priority (bit 7)
				if px > 0 { // Object pixel is transparent
					if ppu.backBuffer[ppu.LY][x] == 0 || !tileAttr(obj.flags).bgPriority() {
						palette := ppu.OBP[tileAttr(obj.flags).dmgPalette()]
						ppu.backBuffer[ppu.LY][x] = palette.GetColor(px)
					}
				}
			}

			// The first object that impacts this pixel will be the one displayed
			break
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
