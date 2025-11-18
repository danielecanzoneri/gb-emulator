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

// When in DMG compatibility mode, the CGB palettes are still being used.
// BGP, OBP0, and OBP1 actually index into the CGB palettes instead of the DMG’s shades of grey.
func (p DMGPalette) ConvertToCGB(cgbPaletteData []uint8) CGBPalette {
	id0 := (p >> 0) & 0b11
	id1 := (p >> 2) & 0b11
	id2 := (p >> 4) & 0b11
	id3 := (p >> 6) & 0b11
	palette := CGBPalette{
		cgbPaletteData[2*id0],
		cgbPaletteData[2*id0+1],
		cgbPaletteData[2*id1],
		cgbPaletteData[2*id1+1],
		cgbPaletteData[2*id2],
		cgbPaletteData[2*id2+1],
		cgbPaletteData[2*id3],
		cgbPaletteData[2*id3+1],
	}

	return palette
}

func (ppu *PPU) GetFrame() *[FrameHeight][FrameWidth]uint16 {
	return ppu.frontBuffer
}

func (ppu *PPU) emptyFrame() {
	var blankPixel uint16 = 0
	if ppu.Cgb {
		blankPixel = 0xFFFF
	}

	for x := range FrameWidth {
		for y := range FrameHeight {
			ppu.backBuffer[y][x] = blankPixel
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

	// Flag that is set to true when x+7 >= wx and used to increment window Y counter
	windowsRendered := false

	for x := uint8(0); x < uint8(FrameWidth); x++ {
		// CGB flag that signals if this pixel BG has higher priority than objs
		var bgPriority = false
		var bgPixel uint8 = 0

		// In CGB mode the LCDC.0 has a different meaning, it is the BG/Window master priority
		if ppu.bgWindowEnabled || ppu.Cgb {
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
			bgPixels, tileAttributes := ppu.GetBGWindowPixelRow(tileAddr, tileY)

			// Pixel color (0-3)
			bgPixel = bgPixels[tileX&0b111]
			var bgPalette Palette = ppu.BGP // DMG palette

			if ppu.Cgb {
				bgPriority = tileAttributes.BGPriority()

				// If we are in compatibility mode, use the boot computed palette
				if ppu.DmgCompatibility {
					bgPalette = ppu.BGP.ConvertToCGB(ppu.BGPalette[0:8])
				} else {
					paletteId := tileAttributes.CGBPalette()
					bgPalette = CGBPalette(ppu.BGPalette[8*paletteId : 8*paletteId+8])
				}
			}

			ppu.backBuffer[ppu.LY][x] = bgPalette.GetColor(bgPixel)
		}

		// Render objects
		if !ppu.objEnabled {
			continue
		}

		// Draw objects with priority
		for i := range ppu.numObjs {
			obj := ppu.objsLY[i]
			if !(obj.x <= x+8 && x+8 < obj.x+8) {
				// The current pixel does not overlap the object
				continue
			}

			// Object row to draw is: LY + 16 - y
			rowPixels := ppu.GetObjectRow(obj, ppu.LY+yObjOffset-obj.y)

			// Draw pixel if no other pixel with higher priority was drawn
			px := rowPixels[(x+8)-obj.x]

			// If object pixel is transparent (px == 0), draw background pixel
			if px > 0 {
				if ppu.Cgb {
					// - If the BG color index is 0, the OBJ will always have priority;
					// - If LCDC bit 0 is clear, the OBJ will always have priority;
					// - If both the BG Attributes and the OAM Attributes have bit 7 clear, the OBJ will have priority;
					// Otherwise, BG will have priority.
					// In CGB mode the LCDC.0 has a different meaning, it is the BG/Window master priority
					if bgPixel == 0 || !ppu.bgWindowEnabled || (!bgPriority && !TileAttribute(obj.flags).BGPriority()) {
						var palette Palette
						// If we are in compatibility mode, use the boot computed palette
						if ppu.DmgCompatibility {
							paletteId := TileAttribute(obj.flags).DMGPalette()
							palette = ppu.OBP[paletteId].ConvertToCGB(ppu.OBJPalette[8*paletteId : 8*paletteId+8])
						} else {
							paletteId := TileAttribute(obj.flags).CGBPalette()
							palette = CGBPalette(ppu.OBJPalette[8*paletteId : 8*paletteId+8])
						}
						ppu.backBuffer[ppu.LY][x] = palette.GetColor(px)
					}
				} else { // DMG
					// If background pixel is 0, draw object pixel
					// If both object and background pixel are not 0, draw pixel based on
					//    object attributes BG/Window priority (bit 7)
					if bgPixel == 0 || !TileAttribute(obj.flags).BGPriority() {
						palette := ppu.OBP[TileAttribute(obj.flags).DMGPalette()]
						ppu.backBuffer[ppu.LY][x] = palette.GetColor(px)
					}
				}

				// The first object that impacts this pixel will be the one displayed
				break
			}
		}
	}

	if windowsRendered {
		ppu.wyCounter++
		penaltyDotsBG += 6 // 6-dot penalty is incurred while the BG fetcher is being set up for the window.
	}

	penaltyDotsObj := 0

	// Tiles considered in the OBJ penalty algorithm (x ranges from 0 to 167+7, so we have at most 22 tiles
	var tileObjectsPenalties [(167+7)>>3 + 1]bool
	// Objects considered in the OBJ penalty algorithm (6 penalty dots per object)
	var objectsPenalties [10]bool

	// OBJ penalty algorithm
	for i := range ppu.numObjs {
		obj := ppu.objsLY[i]
		if obj.x >= 168 { // Out of range tile
			continue
		}

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
	}

	// Round objects penalty dots to M-cycle (TODO - investigate why it doesn't work otherwise)
	return penaltyDotsBG + (penaltyDotsObj & ^3)
}

func getTileMapOffset(x, y uint8) uint16 {
	return (uint16(x) >> 3) | ((uint16(y) & 0xF8) << 2)
}
