package ppu

import "github.com/danielecanzoneri/gb-emulator/util"

// oamBugTriggered returns the OAM row affected by the OAM bug or 0xFF if not triggered
func (ppu *PPU) oamBugTriggered(address uint16) uint8 {
	// Triggered when 16-bit bus is used when content is in range 0xFE00 - 0xFEFF and PPU is in mode 2
	if !(0xFE00 <= address && address < 0xFF00) {
		return 0xFF
	}

	// This means we are in mode 2
	if st, ok := ppu.InternalState.(*oamScan); ok {
		return st.rowAccessed
	}
	return 0xFF
}

func (ppu *PPU) triggerOAMBug(oamRow uint8, bitwiseGlitch func(a, b, c uint16) uint16) {
	prevRow := [8]uint8{
		ppu.OAM.read((oamRow-1)<<3 + 0),
		ppu.OAM.read((oamRow-1)<<3 + 1),
		ppu.OAM.read((oamRow-1)<<3 + 2),
		ppu.OAM.read((oamRow-1)<<3 + 3),
		ppu.OAM.read((oamRow-1)<<3 + 4),
		ppu.OAM.read((oamRow-1)<<3 + 5),
		ppu.OAM.read((oamRow-1)<<3 + 6),
		ppu.OAM.read((oamRow-1)<<3 + 7),
	}
	a := util.CombineBytes(ppu.OAM.read(oamRow<<3+1), ppu.OAM.read(oamRow<<3+0))
	b := util.CombineBytes(prevRow[1], prevRow[0])
	c := util.CombineBytes(prevRow[5], prevRow[4])

	newByte1, newByte0 := util.SplitWord(bitwiseGlitch(a, b, c))

	// Copy values in row
	ppu.OAM.write(oamRow<<3+0, newByte0)
	ppu.OAM.write(oamRow<<3+1, newByte1)
	ppu.OAM.write(oamRow<<3+2, prevRow[2])
	ppu.OAM.write(oamRow<<3+3, prevRow[3])
	ppu.OAM.write(oamRow<<3+4, prevRow[4])
	ppu.OAM.write(oamRow<<3+5, prevRow[5])
	ppu.OAM.write(oamRow<<3+6, prevRow[6])
	ppu.OAM.write(oamRow<<3+7, prevRow[7])
}
func (ppu *PPU) TriggerOAMBugWrite(address uint16) {
	oamRow := ppu.oamBugTriggered(address)
	if oamRow == 0 || oamRow > 19 {
		return
	}

	// TODO - maybe optimize
	// A “write corruption” corrupts the currently access row in the following manner,
	// as long as it’s not the first row (containing the first two objects):
	// - The first word in the row is replaced with this bitwise expression:
	//       ((a ^ c) & (b ^ c)) ^ c
	//   where a is the original value of that word, b is the first word in the preceding row,
	//   and c is the third word in the preceding row.
	// - The last three words are copied from the last three words in the preceding row.
	glitch := func(a, b, c uint16) uint16 {
		return ((a ^ c) & (b ^ c)) ^ c
	}
	ppu.triggerOAMBug(oamRow, glitch)
}
func (ppu *PPU) TriggerOAMBugRead(address uint16) {
	oamRow := ppu.oamBugTriggered(address)
	if oamRow == 0 || oamRow > 19 {
		return
	}

	// TODO - trigger bug
	// A “read corruption” works similarly to a write corruption,
	// except the bitwise expression is b | (a & c).
	glitch := func(a, b, c uint16) uint16 {
		return b | (a & c)
	}
	ppu.triggerOAMBug(oamRow, glitch)
}
