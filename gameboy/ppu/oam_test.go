package ppu

import (
	"testing"
)

func TestObjectSelection8x8(t *testing.T) {
	ppu := New()
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
		writeObjectsOAM(&ppu.oam, objects)
		ppu.searchOAM()

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
		writeObjectsOAM(&ppu.oam, objects)
		ppu.searchOAM()

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
	ppu := New()
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
		writeObjectsOAM(&ppu.oam, objects)
		ppu.searchOAM()

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
		writeObjectsOAM(&ppu.oam, objects)
		ppu.searchOAM()

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
		o.objectsData[i] = obj
	}
}
