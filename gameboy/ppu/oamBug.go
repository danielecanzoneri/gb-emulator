package ppu

func (ppu *PPU) getOamAccessedRow() uint8 {
	return 0xFF
}

func (ppu *PPU) oamBugTriggered(address uint16) bool {
	// Triggered when 16-bit bus is used when content is in range 0xFE00 - 0xFEFF and PPU is in mode 2
	if !(0xFE00 <= address && address < 0xFF00) {
		return false
	}

	return ppu.STAT&3 == 2
}

func (ppu *PPU) TriggerOAMBug(address uint16) {
	if !ppu.oamBugTriggered(address) {
		return
	}

	// TODO - trigger bug
}
