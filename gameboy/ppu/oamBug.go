package ppu

import "github.com/danielecanzoneri/gb-emulator/util"

func readOAMWord(ppu *PPU, row uint8, offset uint8) uint16 {
	low := ppu.OAM.read((row)<<3 + 2*offset)
	high := ppu.OAM.read((row)<<3 + 2*offset + 1)

	return util.CombineBytes(high, low)
}
func writeOAMWord(ppu *PPU, row uint8, offset uint8, word uint16) {
	high, low := util.SplitWord(word)
	ppu.OAM.write((row)<<3+2*offset, low)
	ppu.OAM.write((row)<<3+2*offset+1, high)
}

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
	a := readOAMWord(ppu, oamRow, 0)
	b := readOAMWord(ppu, oamRow-1, 0)
	c := readOAMWord(ppu, oamRow-1, 2)

	newA := bitwiseGlitch(a, b, c)
	writeOAMWord(ppu, oamRow, 0, newA)

	// Copy values of the last three words
	for i := uint8(1); i < 4; i++ {
		writeOAMWord(ppu, oamRow, i, readOAMWord(ppu, oamRow-1, i))
	}
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

	// A “read corruption” works similarly to a write corruption,
	// except the bitwise expression is b | (a & c).
	glitch := func(a, b, c uint16) uint16 {
		return b | (a & c)
	}
	ppu.triggerOAMBug(oamRow, glitch)
}

func (ppu *PPU) TriggerOAMBugWriteBeforeRead(address uint16) {
	// If a register is increased or decreased in the same M-cycle of a read,
	// this will effectively trigger both a read and a write in a single M-cycle,
	// resulting in a more complex corruption pattern:
	// - This corruption will not happen if the accessed row is one of
	//   the first four, as well as if it’s the last row:
	//   - The first word in the row preceding the currently accessed row is replaced
	//     with the following bitwise expression:
	//         (b & (a | c | d)) | (a & c & d)
	//     where a is the first word two rows before the currently accessed row,
	//     b is the first word in the preceding row (the word being corrupted),
	//     c is the first word in the currently accessed row,
	//     and d is the third word in the preceding row.
	//   - The contents of the preceding row is copied (after the corruption of
	//     the first word in it) both to the currently accessed row
	//     and to two rows before the currently accessed row.

	oamRow := ppu.oamBugTriggered(address)
	if oamRow < 4 || oamRow >= 19 {
		return
	}

	a := readOAMWord(ppu, oamRow-2, 0)
	b := readOAMWord(ppu, oamRow-1, 0)
	c := readOAMWord(ppu, oamRow, 0)
	d := readOAMWord(ppu, oamRow-1, 2)

	newB := (b & (a | c | d)) | (a & c & d)
	writeOAMWord(ppu, oamRow-1, 0, newB)

	// Copy preceding row
	for i := uint8(0); i < 4; i++ {
		w := readOAMWord(ppu, oamRow-1, i)
		writeOAMWord(ppu, oamRow-2, i, w)
		writeOAMWord(ppu, oamRow, i, w)
	}
}
