package joypad

import "github.com/hajimehoshi/ebiten/v2"

const (
	KeyStart = iota
	KeySelect
	KeyB
	KeyA
	KeyDown
	KeyUp
	KeyLeft
	KeyRight
)

var buttonsKeyMapping = map[int]ebiten.Key{
	KeyStart:  ebiten.KeyX,
	KeySelect: ebiten.KeyZ,
	KeyB:      ebiten.KeyA,
	KeyA:      ebiten.KeyS,
}
var dPadKeyMapping = map[int]ebiten.Key{
	KeyDown:  ebiten.KeyDown,
	KeyUp:    ebiten.KeyUp,
	KeyLeft:  ebiten.KeyLeft,
	KeyRight: ebiten.KeyRight,
}

func (jp *Joypad) DetectKeysPressed() {
	jp.selectUp = 1
	jp.startDown = 1
	jp.bLeft = 1
	jp.aRight = 1

	if jp.selectButtons == 0 {
		if ebiten.IsKeyPressed(buttonsKeyMapping[KeyStart]) {
			jp.startDown = 0
		}
		if ebiten.IsKeyPressed(buttonsKeyMapping[KeySelect]) {
			jp.selectUp = 0
		}
		if ebiten.IsKeyPressed(buttonsKeyMapping[KeyB]) {
			jp.bLeft = 0
		}
		if ebiten.IsKeyPressed(buttonsKeyMapping[KeyA]) {
			jp.aRight = 0
		}
	}
	if jp.selectDPad == 0 {
		if ebiten.IsKeyPressed(dPadKeyMapping[KeyDown]) {
			jp.startDown = 0
		}
		if ebiten.IsKeyPressed(dPadKeyMapping[KeyUp]) {
			jp.selectUp = 0
		}
		if ebiten.IsKeyPressed(dPadKeyMapping[KeyLeft]) {
			jp.bLeft = 0
		}
		if ebiten.IsKeyPressed(dPadKeyMapping[KeyRight]) {
			jp.aRight = 0
		}
	}
}
