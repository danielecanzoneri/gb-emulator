package ppu

import (
	"cmp"
	"github.com/danielecanzoneri/gb-emulator/util"
	"slices"
	"strconv"
)

const (
	OAMSize     = 0xA0
	objectsSize = 4
)

type OAM struct {
	objectsData [OAMSize / objectsSize]Object

	// Disabled during mode 2 (OAM scan) and 3 (drawing)
	readDisabled  bool
	writeDisabled bool
}

func (oam *OAM) Read(addr uint16) uint8 {
	if oam.readDisabled {
		return 0xFF
	}

	addr -= OAMStartAddr
	objectId := addr / objectsSize
	return oam.objectsData[objectId].read(addr % objectsSize)
}

func (oam *OAM) Write(addr uint16, value uint8) {
	if oam.writeDisabled {
		return
	}

	addr -= OAMStartAddr
	objectId := addr / objectsSize
	oam.objectsData[objectId].write(addr%objectsSize, value)
}

const (
	objsLimit = 10

	yOffset = 16
	xOffset = 8
)

type Object struct {
	y uint8 // Byte 0
	x uint8 // Byte 1

	flags      uint8 // Byte 3
	bgPriority bool
	yFlip      bool
	xFlip      bool
	paletteId  uint8 // 0: OBP0, 1: OBP1

	tileIndex uint8 // Byte 2
}

func (obj *Object) read(addr uint16) uint8 {
	switch addr {
	case 0:
		return obj.y
	case 1:
		return obj.x
	case 2:
		return obj.tileIndex
	case 3:
		return obj.flags
	default:
		panic("should never happen")
	}
}

func (obj *Object) write(addr uint16, value uint8) {
	switch addr {
	case 0:
		obj.y = value
	case 1:
		obj.x = value
	case 2:
		obj.tileIndex = value
	case 3:
		obj.flags = value
		// Bit 7: priority; Bit 6: y-flip; Bit 5: x-flip; Bit 4: palette
		obj.bgPriority = util.ReadBit(value, 7) > 0
		obj.yFlip = util.ReadBit(value, 6) > 0
		obj.xFlip = util.ReadBit(value, 5) > 0
		obj.paletteId = util.ReadBit(value, 4)

	default:
		panic("should never happen")
	}
}

func (ppu *PPU) getObjectRow(obj *Object, row uint8) [8]uint8 {
	if (!ppu.obj8x16Size && row >= 8) || (ppu.obj8x16Size && row >= 16) {
		panic("Invalid row: " + strconv.Itoa(int(row)))
	}

	if obj.yFlip {
		if !ppu.obj8x16Size { // height = 8
			row = 7 - row
		} else { // height = 16
			row = 15 - row
		}
	}

	var pixels [8]uint8
	if !ppu.obj8x16Size {
		tile := ppu.ReadTileObj(obj.tileIndex)
		pixels = tile.getRowPixels(row)
	} else {
		if row < 8 {
			tile := ppu.ReadTileObj(obj.tileIndex & 0xFE)
			pixels = tile.getRowPixels(row)
		} else {
			tile := ppu.ReadTileObj(obj.tileIndex | 0b01)
			pixels = tile.getRowPixels(row - 8)
		}
	}

	if obj.xFlip {
		pixels[0], pixels[7] = pixels[7], pixels[0]
		pixels[1], pixels[6] = pixels[6], pixels[1]
		pixels[2], pixels[5] = pixels[5], pixels[2]
		pixels[3], pixels[4] = pixels[4], pixels[3]
	}
	return pixels
}

func (ppu *PPU) selectObjects() {
	// objHeight is used to check which objects are currently on the line
	var objHeight uint8 = 8
	if ppu.obj8x16Size {
		objHeight = 16
	}

	ppu.numObjs = 0

	// Scan OAM and select objects that lie in current line
	for _, obj := range ppu.oam.objectsData {
		// obj is on the line if obj.y <= LY+16 < obj.y + height
		if obj.y <= ppu.LY+yOffset && ppu.LY+yOffset < obj.y+objHeight {
			ppu.objsLY[ppu.numObjs] = &obj
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
