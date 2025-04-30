package ppu

import (
	"cmp"
	"github.com/danielecanzoneri/gb-emulator/internal/util"
	"slices"
	"strconv"
)

const (
	objSize = 4

	objsLimit = 10

	yOffset = 16
	xOffset = 8
)

type Object struct {
	y uint8
	x uint8

	bgPriority bool
	yFlip      bool
	xFlip      bool
	palette    Palette

	tile1 *Tile
	tile2 *Tile // Only if 8x16 object
}

func (obj Object) getRow(row uint8) [8]uint8 {
	if (obj.tile2 == nil && row >= 8) || (obj.tile2 != nil && row >= 16) {
		panic("Invalid row: " + strconv.Itoa(int(row)))
	}

	if obj.yFlip {
		if obj.tile2 == nil { // height = 8
			row = 7 - row
		} else { // height = 16
			row = 15 - row
		}
	}
	var pixels [8]uint8
	if row < 8 {
		pixels = obj.tile1.getRowPixels(row)
	} else {
		pixels = obj.tile2.getRowPixels(row - 8)
	}

	if obj.xFlip {
		pixels[0], pixels[7] = pixels[7], pixels[0]
		pixels[1], pixels[6] = pixels[6], pixels[1]
		pixels[2], pixels[5] = pixels[5], pixels[2]
		pixels[3], pixels[4] = pixels[4], pixels[3]
	}
	return pixels
}

type Palette uint8

func (p Palette) getColor(id uint8) uint8 {
	var mask uint8 = 0b11
	id &= mask

	// Bit 7,6: id 3; Bit 5,4: id 2; Bit 3,2: id 1; Bit 1,0: id 0
	shift := id * 2
	return (uint8(p) >> shift) & mask
}

func (ppu *PPU) parseObject(objAddr uint8) *Object {
	data := ppu.OAM.Data[objAddr : objAddr+objSize]

	var tile1, tile2 *Tile
	if !ppu.obj8x16Size {
		tile1 = ppu.ReadTileObj(data[2])
	} else {
		tile1 = ppu.ReadTileObj(data[2] & 0xFE)
		tile2 = ppu.ReadTileObj(data[2] | 0b01)
	}

	// byte4 - Bit 7: priority; Bit 6: y-flip; Bit 5: x-flip; Bit 4: palette
	var palette Palette
	if util.ReadBit(data[3], 4) == 0 {
		palette = ppu.OBP0
	} else {
		palette = ppu.OBP1
	}

	return &Object{
		y:          data[0],
		x:          data[1],
		tile1:      tile1,
		tile2:      tile2,
		bgPriority: util.ReadBit(data[3], 7) > 0,
		yFlip:      util.ReadBit(data[3], 6) > 0,
		xFlip:      util.ReadBit(data[3], 5) > 0,
		palette:    palette,
	}
}

func (ppu *PPU) selectObjects() {
	// objHeight is used to check which objects are currently on the line
	var objHeight uint8 = 8
	if ppu.obj8x16Size {
		objHeight = 16
	}

	ppu.numObjs = 0

	// Scan OAM and select objects that lie in current line
	for of := 0; of < OAMSize; of += 4 {
		// obj is on the line if obj.y <= LY+16 < obj.y + height
		if ppu.OAM.Data[of] <= ppu.LY+yOffset && ppu.LY+yOffset < ppu.OAM.Data[of]+objHeight {
			ppu.objsLY[ppu.numObjs] = ppu.parseObject(uint8(of))
			ppu.numObjs++

			if ppu.numObjs == objsLimit {
				break
			}
		}
	}

	// Sort objects by priority (lower x have priority)
	objs := ppu.objsLY[0:ppu.numObjs]
	slices.SortStableFunc(objs, func(a, b *Object) int {
		return cmp.Compare(a.x, b.x)
	})
}
