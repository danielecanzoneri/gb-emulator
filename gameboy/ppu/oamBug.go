package ppu

import "github.com/danielecanzoneri/gb-emulator/util"

// Set flags for OAM bug if this is a glitched access
func (ppu *PPU) GlitchedOAMAccess(address uint16, isRead bool) {
	// Triggered when 16-bit bus is used when content is in range 0xFE00 - 0xFEFF and PPU is in mode 2
	if !(0xFE00 <= address && address < 0xFF00) {
		return
	}

	// This means we are in mode 2
	if st, ok := ppu.internalState.(*oamScan); ok {
		ppu.oam.buggedRow = st.rowAccessed
		if isRead {
			ppu.oam.buggedRead = true
		} else {
			ppu.oam.buggedWrite = true
		}
	}
}

func readOAMWord(ppu *PPU, row uint8, offset uint8) uint16 {
	low := ppu.oam.read((row)<<3 + 2*offset)
	high := ppu.oam.read((row)<<3 + 2*offset + 1)

	return util.CombineBytes(high, low)
}
func writeOAMWord(ppu *PPU, row uint8, offset uint8, word uint16) {
	high, low := util.SplitWord(word)
	ppu.oam.write((row)<<3+2*offset, low)
	ppu.oam.write((row)<<3+2*offset+1, high)
}

func (ppu *PPU) triggerOAMBugWrite() {
	if ppu.oam.buggedRow == 0 || ppu.oam.buggedRow > 19 {
		return
	}

	// A "write corruption" corrupts the currently access row in the following manner,
	// as long as it's not the first row (containing the first two objects):
	// - The first word in the row is replaced with this bitwise expression:
	//       ((a ^ c) & (b ^ c)) ^ c
	//   where a is the original value of that word, b is the first word in the preceding row,
	//   and c is the third word in the preceding row.
	// - The last three words are copied from the last three words in the preceding row.

	a := readOAMWord(ppu, ppu.oam.buggedRow, 0)
	b := readOAMWord(ppu, ppu.oam.buggedRow-1, 0)
	c := readOAMWord(ppu, ppu.oam.buggedRow-1, 2)

	newA := ((a ^ c) & (b ^ c)) ^ c
	writeOAMWord(ppu, ppu.oam.buggedRow, 0, newA)

	// Copy values of the last three words
	for i := uint8(1); i < 4; i++ {
		writeOAMWord(ppu, ppu.oam.buggedRow, i, readOAMWord(ppu, ppu.oam.buggedRow-1, i))
	}
}
func (ppu *PPU) triggerOAMBugRead() {
	if ppu.oam.buggedRow == 0 || ppu.oam.buggedRow > 19 {
		return
	}

	// A "read corruption" works similarly to a write corruption,
	// except the bitwise expression is b | (a & c).
	a := readOAMWord(ppu, ppu.oam.buggedRow, 0)
	b := readOAMWord(ppu, ppu.oam.buggedRow-1, 0)
	c := readOAMWord(ppu, ppu.oam.buggedRow-1, 2)

	newA := b | (a & c)
	writeOAMWord(ppu, ppu.oam.buggedRow, 0, newA)

	// Copy values of the last three words
	for i := uint8(1); i < 4; i++ {
		writeOAMWord(ppu, ppu.oam.buggedRow, i, readOAMWord(ppu, ppu.oam.buggedRow-1, i))
	}
}

func (ppu *PPU) triggerOAMBugWriteAndRead() {
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
	// - Regardless of whether the previous corruption occurred or not, a normal read corruption is then applied.

	defer ppu.triggerOAMBugRead()

	if ppu.oam.buggedRow < 4 || ppu.oam.buggedRow >= 19 {
		return
	}

	a := readOAMWord(ppu, ppu.oam.buggedRow-2, 0)
	b := readOAMWord(ppu, ppu.oam.buggedRow-1, 0)
	c := readOAMWord(ppu, ppu.oam.buggedRow, 0)
	d := readOAMWord(ppu, ppu.oam.buggedRow-1, 2)

	newB := (b & (a | c | d)) | (a & c & d)
	writeOAMWord(ppu, ppu.oam.buggedRow-1, 0, newB)

	// Copy preceding row
	for i := uint8(0); i < 4; i++ {
		w := readOAMWord(ppu, ppu.oam.buggedRow-1, i)
		writeOAMWord(ppu, ppu.oam.buggedRow-2, i, w)
		writeOAMWord(ppu, ppu.oam.buggedRow, i, w)
	}
}
