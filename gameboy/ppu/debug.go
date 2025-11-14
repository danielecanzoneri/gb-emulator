package ppu

// Debug methods for debugger access to internal PPU state
// These methods provide read-only access to internal implementation details
// that are needed for debugging but should not be part of the public API.

// DebugGetBGTileMapAddr returns the background tile map address
func (ppu *PPU) DebugGetBGTileMapAddr() uint16 {
	return ppu.bgTileMapAddr
}

// DebugGetTileData returns a copy of tile data for the specified bank and tile index
func (ppu *PPU) DebugGetTileData(bank uint8, tileIndex uint16) Tile {
	return ppu.vRAM.tileData[bank][tileIndex]
}

// DebugGetBGPalette returns a copy of the background palette
func (ppu *PPU) DebugGetBGPalette() [64]uint8 {
	return ppu.BGPalette
}

// DebugGetOBJPalette returns a copy of the object palette
func (ppu *PPU) DebugGetOBJPalette() [64]uint8 {
	return ppu.OBJPalette
}

// DebugGetOAMObject returns a pointer to the OAM object at the specified index
func (ppu *PPU) DebugGetOAMObject(index int) *Object {
	if index < 0 || index >= len(ppu.oam.Data) {
		return nil
	}
	return &ppu.oam.Data[index]
}

// DebugGetDots returns the number of dots elapsed rendering the current line
func (ppu *PPU) DebugGetDots() int {
	return ppu.dots
}

// DebugGetInternalState returns the current internal state (for debugging)
// Returns nil if no state is set
func (ppu *PPU) DebugGetInternalState() ppuInternalState {
	return ppu.internalState
}

// DebugGetInternalStateLength returns the remaining length of the current internal state
func (ppu *PPU) DebugGetInternalStateLength() int {
	return ppu.internalStateLength
}

// DebugGetMode returns the current PPU mode (0-3) extracted from STAT register
func (ppu *PPU) DebugGetMode() uint8 {
	return ppu.STAT & 3
}
