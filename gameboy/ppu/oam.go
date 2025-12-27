package ppu

import (
	"cmp"
	"slices"
)

const (
	objsLimit = 10

	yObjOffset = 16
	xObjOffset = 8

	OAMSize = 0xA0
)

type Object struct {
	y         uint8 // Byte 0
	x         uint8 // Byte 1
	tileIndex uint8 // Byte 2
	flags     uint8 // Byte 3
}

func (obj *Object) Read(addr uint8) uint8 {
	switch addr & 0b11 {
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

func (obj *Object) write(addr uint8, value uint8) {
	switch addr & 0b11 {
	case 0:
		obj.y = value
	case 1:
		obj.x = value
	case 2:
		obj.tileIndex = value
	case 3:
		obj.flags = value
	default:
		panic("should never happen")
	}
}

type OAM struct {
	Data [OAMSize / 4]Object

	// Disabled during mode 2 (OAM scan) and 3 (drawing)
	readDisabled  bool
	writeDisabled bool

	// OAM bug flags
	buggedRead  bool
	buggedWrite bool
	buggedRow   uint8
}

func (oam *OAM) Read(addr uint16) uint8 {
	if addr >= 0xFEA0 || oam.readDisabled {
		return 0xFF
	}
	return oam.read(uint8(addr - OAMStartAddr))
}
func (oam *OAM) read(addr uint8) uint8 {
	objectId := addr >> 2
	return oam.Data[objectId].Read(addr)
}

func (oam *OAM) Write(addr uint16, value uint8) {
	if addr >= 0xFEA0 || oam.writeDisabled {
		return
	}
	oam.write(uint8(addr-OAMStartAddr), value)
}
func (oam *OAM) write(addr uint8, value uint8) {
	objectId := addr >> 2
	oam.Data[objectId].write(addr, value)
}

func (ppu *PPU) DMAWrite(index uint16, value uint8) {
	ppu.oam.Write(OAMStartAddr+index, value)
}

func (ppu *PPU) GetObjectRow(obj *Object, row uint8) [8]uint8 {
	// CGB only
	vRAMBank := uint8(0)
	if ppu.Cgb {
		vRAMBank = TileAttribute(obj.flags).Bank()
	}

	if !ppu.obj8x16Size { // 8-row object
		tile := ppu.ReadTileObj(obj.tileIndex, vRAMBank)
		return tile.GetRow(TileAttribute(obj.flags), row&0x7)
	}

	// 16 row object
	row &= 0xF
	if TileAttribute(obj.flags).YFlip() {
		// Flip tile (just switch 4-th bit)
		bit4 := row & 0x8
		row = (row &^ 0x8) | (^bit4 & 0x8)
	}

	// 2-tiles object, fetch the correct one
	if row < 8 {
		tile := ppu.ReadTileObj(obj.tileIndex&0xFE, vRAMBank)
		return tile.GetRow(TileAttribute(obj.flags), row)
	} else {
		tile := ppu.ReadTileObj(obj.tileIndex|0b01, vRAMBank)
		return tile.GetRow(TileAttribute(obj.flags), row-8)
	}
}

func (ppu *PPU) searchOAM() {
	// objHeight is used to check which objects are currently on the line
	var objHeight uint8 = 8
	if ppu.obj8x16Size {
		objHeight = 16
	}

	ppu.numObjs = 0

	// Scan OAM and select objects that lie in current line
	for _, obj := range ppu.oam.Data {
		// obj is on the line if obj.y <= LY+16 < obj.y + height
		if obj.y <= ppu.LY+yObjOffset && ppu.LY+yObjOffset < obj.y+objHeight {
			ppu.objsLY[ppu.numObjs] = &obj
			ppu.numObjs++

			if ppu.numObjs == objsLimit {
				break
			}
		}
	}

	// Sort objects by priority (lower x have priority)
	objs := ppu.objsLY[0:ppu.numObjs]

	// In CGB mode, only the object’s location in OAM determines its priority.
	// The earlier the object, the higher its priority.
	if !ppu.Cgb || ppu.DmgCompatibility {
		slices.SortStableFunc(objs, func(a, b *Object) int {
			return cmp.Compare(a.x, b.x)
		})
	}
}
