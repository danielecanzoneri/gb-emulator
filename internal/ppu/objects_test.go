package ppu

import (
	"testing"
)

func TestObjectSelection8x8(t *testing.T) {
	ppu := &PPU{}
	ppu.obj8x16Size = false
	ppu.LY = 2 // 3rd line

	t.Run("<=10", func(t *testing.T) {
		t.Log(ppu)
		objects := []Object{
			{y: 11}, // Shown
			{y: 14}, // Shown
			{y: 18}, // Shown
			{y: 10}, // Hidden
			{y: 19}, // Hidden
		}
		writeObjectsOAM(&ppu.OAM, objects)
		ppu.selectObjects()

		if ppu.numObjs != 3 {
			t.Errorf("Expected 3 objects, got %d", ppu.numObjs)
		}
		// Check correct obj position
		for i := 0; i < ppu.numObjs; i++ {
			if ppu.objsLY[i].y != objects[i].y {
				t.Errorf("object[%d].y, got %d, expected %d", i, ppu.objsLY[i].y, objects[i].y)
			}
		}
	})

	t.Run(">10", func(t *testing.T) {
		objects := []Object{
			{y: 11},         // Shown
			{y: 12},         // Shown
			{y: 13},         // Shown
			{y: 14},         // Shown
			{y: 15},         // Shown
			{y: 16},         // Shown
			{y: 17},         // Shown
			{y: 18},         // Shown
			{y: 14, x: 100}, // Shown
			{y: 15, x: 100}, // Shown
			{y: 11, x: 100}, // Hidden because of limit
		}
		writeObjectsOAM(&ppu.OAM, objects)
		ppu.selectObjects()

		if ppu.numObjs != objsLimit {
			t.Errorf("Expected 3 objects, got %d", ppu.numObjs)
		}
		// Check correct obj position
		for i := 0; i < ppu.numObjs; i++ {
			if ppu.objsLY[i].y != objects[i].y {
				t.Errorf("object[%d].y, got %d, expected %d", i, ppu.objsLY[i].y, objects[i].y)
			}
		}
	})
}

func TestObjectSelection8x16(t *testing.T) {
	ppu := &PPU{}
	ppu.obj8x16Size = true
	ppu.LY = 2 // 3rd line

	t.Run("<=10", func(t *testing.T) {
		objects := []Object{
			{y: 3},  // Shown
			{y: 11}, // Shown
			{y: 18}, // Shown
			{y: 2},  // Hidden
			{y: 19}, // Hidden
		}
		writeObjectsOAM(&ppu.OAM, objects)
		ppu.selectObjects()

		if ppu.numObjs != 3 {
			t.Errorf("Expected 3 objects, got %d", ppu.numObjs)
		}
		// Check correct obj position
		for i := 0; i < ppu.numObjs; i++ {
			if ppu.objsLY[i].y != objects[i].y {
				t.Errorf("object[%d].y, got %d, expected %d", i, ppu.objsLY[i].y, objects[i].y)
			}
		}
	})

	t.Run(">10", func(t *testing.T) {
		objects := []Object{
			{y: 3},  // Shown
			{y: 4},  // Shown
			{y: 5},  // Shown
			{y: 6},  // Shown
			{y: 7},  // Shown
			{y: 14}, // Shown
			{y: 15}, // Shown
			{y: 16}, // Shown
			{y: 17}, // Shown
			{y: 18}, // Shown
			{y: 10}, // Hidden because of limit
		}
		writeObjectsOAM(&ppu.OAM, objects)
		ppu.selectObjects()

		if ppu.numObjs != objsLimit {
			t.Errorf("Expected 3 objects, got %d", ppu.numObjs)
		}
		// Check correct obj position
		for i := 0; i < ppu.numObjs; i++ {
			if ppu.objsLY[i].y != objects[i].y {
				t.Errorf("object[%d].y, got %d, expected %d", i, ppu.objsLY[i].y, objects[i].y)
			}
		}
	})
}

func writeObjectsOAM(o *OAM, objs []Object) {
	for i, obj := range objs {
		o.data[objSize*i] = obj.y
		o.data[objSize*i+1] = obj.x
	}
}

func TestObjectParsing8x8(t *testing.T) {
	ppu := &PPU{}
	ppu.OBP0 = 0
	ppu.OBP1 = 0xFF
	ppu.obj8x16Size = false

	var objAddr uint8 = 0x41
	var objData = [4]uint8{
		16,   // y
		7,    // x
		0xE1, // tileAddr
		0xFF, // Flags
	}

	// Write obj data in OAM
	copy(ppu.OAM.data[objAddr:objAddr+4], objData[:])

	// Write tile data at obj addr
	copy(ppu.vRAM.data[16*int(objData[2]):16*(int(objData[2])+1)], TestTileData[:])

	obj := ppu.parseObject(objAddr)
	if obj.y != objData[0] {
		t.Errorf("y: expected %d, got %d", objData[0], obj.y)
	}
	if obj.x != objData[1] {
		t.Errorf("x: expected %d, got %d", objData[1], obj.x)
	}
	if !obj.bgPriority {
		t.Errorf("bgPriority: expected %t, got %t", true, obj.bgPriority)
	}
	if !obj.yFlip {
		t.Errorf("yFlip: expected %t, got %t", true, obj.yFlip)
	}
	if !obj.xFlip {
		t.Errorf("xFlip: expected %t, got %t", true, obj.xFlip)
	}
	if obj.palette != ppu.OBP1 {
		t.Error("palette: expected OBP1, got OBP0")
	}
	if obj.tile1.data != TestExpectedTile {
		t.Errorf("tile1: expected %v, got %v", TestExpectedTile, obj.tile1.data)
	}
	if obj.tile2 != nil {
		t.Errorf("tile2: expected nil")
	}
}

func TestObjectParsing8x16(t *testing.T) {
	ppu := &PPU{}
	ppu.OBP0 = 0
	ppu.OBP1 = 0xFF
	ppu.obj8x16Size = true

	var objAddr uint8 = 0x41
	var objData = [4]uint8{
		16,   // y
		7,    // x
		0xE1, // tileAddr
		0xFF, // Flags
	}

	// Write obj data in OAM
	copy(ppu.OAM.data[objAddr:objAddr+4], objData[:])

	// Write 2 tile data at 0xE0 and 0xE1
	copy(ppu.vRAM.data[16*0xE0:16*0xE1], TestTileData[:])
	copy(ppu.vRAM.data[16*0xE1:16*0xE2], TestTileData[:])

	obj := ppu.parseObject(objAddr)
	if obj.y != objData[0] {
		t.Errorf("y: expected %d, got %d", objData[0], obj.y)
	}
	if obj.x != objData[1] {
		t.Errorf("x: expected %d, got %d", objData[1], obj.x)
	}
	if !obj.bgPriority {
		t.Errorf("bgPriority: expected %t, got %t", true, obj.bgPriority)
	}
	if !obj.yFlip {
		t.Errorf("yFlip: expected %t, got %t", true, obj.yFlip)
	}
	if !obj.xFlip {
		t.Errorf("xFlip: expected %t, got %t", true, obj.xFlip)
	}
	if obj.palette != ppu.OBP1 {
		t.Error("palette: expected OBP1, got OBP0")
	}
	if obj.tile1.data != TestExpectedTile {
		t.Errorf("tile1: expected %v, got %v", TestExpectedTile, obj.tile1.data)
	}
	if obj.tile2.data != TestExpectedTile {
		t.Errorf("tile2: expected %v, got %v", TestExpectedTile, obj.tile2.data)
	}
}
