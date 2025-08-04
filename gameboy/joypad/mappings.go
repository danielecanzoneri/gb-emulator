package joypad

import "github.com/hajimehoshi/ebiten/v2"

type gbKey uint8

const (
	KeyStart gbKey = iota
	KeySelect
	KeyB
	KeyA
	KeyDown
	KeyUp
	KeyLeft
	KeyRight
)

var buttonsKeyMapping = map[gbKey]ebiten.Key{
	KeyStart:  ebiten.KeyX,
	KeySelect: ebiten.KeyZ,
	KeyB:      ebiten.KeyA,
	KeyA:      ebiten.KeyS,
}
var dPadKeyMapping = map[gbKey]ebiten.Key{
	KeyDown:  ebiten.KeyDown,
	KeyUp:    ebiten.KeyUp,
	KeyLeft:  ebiten.KeyLeft,
	KeyRight: ebiten.KeyRight,
}

func (jp *Joypad) DetectKeysPressed() {
	if jp.selectButtons == 0 {
		if ebiten.IsKeyPressed(buttonsKeyMapping[KeyStart]) {
			if jp.startDown == 1 { // Detect high -> low transition
				jp.RequestInterrupt()
			}
			jp.startDown = 0
		}
		if ebiten.IsKeyPressed(buttonsKeyMapping[KeySelect]) {
			if jp.selectUp == 1 { // Detect high -> low transition
				jp.RequestInterrupt()
			}
			jp.selectUp = 0
		}
		if ebiten.IsKeyPressed(buttonsKeyMapping[KeyB]) {
			if jp.bLeft == 1 { // Detect high -> low transition
				jp.RequestInterrupt()
			}
			jp.bLeft = 0
		}
		if ebiten.IsKeyPressed(buttonsKeyMapping[KeyA]) {
			if jp.aRight == 1 { // Detect high -> low transition
				jp.RequestInterrupt()
			}
			jp.aRight = 0
		}
	}
	if jp.selectDPad == 0 {
		if ebiten.IsKeyPressed(dPadKeyMapping[KeyDown]) {
			if jp.startDown == 1 { // Detect high -> low transition
				jp.RequestInterrupt()
			}
			jp.startDown = 0
		}
		if ebiten.IsKeyPressed(dPadKeyMapping[KeyUp]) {
			if jp.selectUp == 1 { // Detect high -> low transition
				jp.RequestInterrupt()
			}
			jp.selectUp = 0
		}
		if ebiten.IsKeyPressed(dPadKeyMapping[KeyLeft]) {
			if jp.bLeft == 1 { // Detect high -> low transition
				jp.RequestInterrupt()
			}
			jp.bLeft = 0
		}
		if ebiten.IsKeyPressed(dPadKeyMapping[KeyRight]) {
			if jp.aRight == 1 { // Detect high -> low transition
				jp.RequestInterrupt()
			}
			jp.aRight = 0
		}
	}
}
