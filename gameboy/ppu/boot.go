package ppu

import (
	"github.com/danielecanzoneri/lucky-boy/gameboy/cartridge"
	"github.com/danielecanzoneri/lucky-boy/util"
)

func (ppu *PPU) SkipDMGBoot() {
	ppu.dots = 400
	ppu.internalState = new(vBlank)
	ppu.interruptMode = 1
	ppu.internalStateLength = 56

	// vRAM
	ppu.vRAM.tileData[0][1].raw = [16]uint8{0xF0, 0, 0xF0, 0, 0xFC, 0, 0xFC, 0, 0xFC, 0, 0xFC, 0, 0xF3, 0, 0xF3, 0}
	ppu.vRAM.tileData[0][2].raw = [16]uint8{0x3C, 0, 0x3C, 0, 0x3C, 0, 0x3C, 0, 0x3C, 0, 0x3C, 0, 0x3C, 0, 0x3C, 0}
	ppu.vRAM.tileData[0][3].raw = [16]uint8{0xF0, 0, 0xF0, 0, 0xF0, 0, 0xF0, 0, 0, 0, 0, 0, 0xF3, 0, 0xF3, 0}
	ppu.vRAM.tileData[0][4].raw = [16]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0xCF, 0, 0xCF, 0}
	ppu.vRAM.tileData[0][5].raw = [16]uint8{0, 0, 0, 0, 0x0F, 0, 0x0F, 0, 0x3F, 0, 0x3F, 0, 0x0F, 0, 0x0F, 0}
	ppu.vRAM.tileData[0][6].raw = [16]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0xC0, 0, 0xC0, 0, 0x0F, 0, 0x0F, 0}
	ppu.vRAM.tileData[0][7].raw = [16]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0xF0, 0, 0xF0, 0}
	ppu.vRAM.tileData[0][8].raw = [16]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0xF3, 0, 0xF3, 0}
	ppu.vRAM.tileData[0][9].raw = [16]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0xC0, 0, 0xC0, 0}
	ppu.vRAM.tileData[0][10].raw = [16]uint8{0x03, 0, 0x03, 0, 0x03, 0, 0x03, 0, 0x03, 0, 0x03, 0, 0xFF, 0, 0xFF, 0}
	ppu.vRAM.tileData[0][11].raw = [16]uint8{0xC0, 0, 0xC0, 0, 0xC0, 0, 0xC0, 0, 0xC0, 0, 0xC0, 0, 0xC3, 0, 0xC3, 0}
	ppu.vRAM.tileData[0][12].raw = [16]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0xFC, 0, 0xFC, 0}
	ppu.vRAM.tileData[0][13].raw = [16]uint8{0xF3, 0, 0xF3, 0, 0xF0, 0, 0xF0, 0, 0xF0, 0, 0xF0, 0, 0xF0, 0, 0xF0, 0}
	ppu.vRAM.tileData[0][14].raw = [16]uint8{0x3C, 0, 0x3C, 0, 0xFC, 0, 0xFC, 0, 0xFC, 0, 0xFC, 0, 0x3C, 0, 0x3C, 0}
	ppu.vRAM.tileData[0][15].raw = [16]uint8{0xF3, 0, 0xF3, 0, 0xF3, 0, 0xF3, 0, 0xF3, 0, 0xF3, 0, 0xF3, 0, 0xF3, 0}
	ppu.vRAM.tileData[0][16].raw = [16]uint8{0xF3, 0, 0xF3, 0, 0xC3, 0, 0xC3, 0, 0xC3, 0, 0xC3, 0, 0xC3, 0, 0xC3, 0}
	ppu.vRAM.tileData[0][17].raw = [16]uint8{0xCF, 0, 0xCF, 0, 0xCF, 0, 0xCF, 0, 0xCF, 0, 0xCF, 0, 0xCF, 0, 0xCF, 0}
	ppu.vRAM.tileData[0][18].raw = [16]uint8{0x3C, 0, 0x3C, 0, 0x3F, 0, 0x3F, 0, 0x3C, 0, 0x3C, 0, 0x0F, 0, 0x0F, 0}
	ppu.vRAM.tileData[0][19].raw = [16]uint8{0x3C, 0, 0x3C, 0, 0xFC, 0, 0xFC, 0, 0, 0, 0, 0, 0xFC, 0, 0xFC, 0}
	ppu.vRAM.tileData[0][20].raw = [16]uint8{0xFC, 0, 0xFC, 0, 0xF0, 0, 0xF0, 0, 0xF0, 0, 0xF0, 0, 0xF0, 0, 0xF0, 0}
	ppu.vRAM.tileData[0][21].raw = [16]uint8{0xF3, 0, 0xF3, 0, 0xF3, 0, 0xF3, 0, 0xF3, 0, 0xF3, 0, 0xF0, 0, 0xF0, 0}
	ppu.vRAM.tileData[0][22].raw = [16]uint8{0xC3, 0, 0xC3, 0, 0xC3, 0, 0xC3, 0, 0xC3, 0, 0xC3, 0, 0xFF, 0, 0xFF, 0}
	ppu.vRAM.tileData[0][23].raw = [16]uint8{0xCF, 0, 0xCF, 0, 0xCF, 0, 0xCF, 0, 0xCF, 0, 0xCF, 0, 0xC3, 0, 0xC3, 0}
	ppu.vRAM.tileData[0][24].raw = [16]uint8{0x0F, 0, 0x0F, 0, 0x0F, 0, 0x0F, 0, 0x0F, 0, 0x0F, 0, 0xFC, 0, 0xFC, 0}
	ppu.vRAM.tileData[0][25].raw = [16]uint8{0x3C, 0, 0x42, 0, 0xB9, 0, 0xA5, 0, 0xB9, 0, 0xA5, 0, 0x42, 0, 0x3C, 0}
	for _, t := range ppu.vRAM.tileData[0] {
		t.updatePixels()
	}

	for i := range 13 { // 260: 1 ... 271: 12
		ppu.vRAM.tileMaps[0][0x103+i] = uint8(i)
	}
	ppu.vRAM.tileMaps[0][0x110] = 25
	for i := range 12 { // 292: 13 ... 303: 24
		ppu.vRAM.tileMaps[0][0x124+i] = uint8(i + 13)
	}

	ppu.LCDC = 0x91
	ppu.STAT = 0x81
	ppu.active = true
	ppu.LY = 0x99
	ppu.BGP = 0xFC
	ppu.windowTileMapAddr = 0x9800
	ppu.bgWindowTileDataArea = 1
	ppu.bgTileMapAddr = 0x9800
	ppu.bgWindowEnabled = true
}

func (ppu *PPU) SkipCGBBoot(rom cartridge.Cartridge) (titleChecksum uint8) {
	ppu.dots = 400
	ppu.internalState = new(vBlank)
	ppu.interruptMode = 1
	ppu.internalStateLength = 56

	ppu.LCDC = 0x91
	ppu.STAT = 0x81
	ppu.active = true
	ppu.LY = 0x99
	ppu.windowTileMapAddr = 0x9800
	ppu.bgWindowTileDataArea = 1
	ppu.bgTileMapAddr = 0x9800
	ppu.bgWindowEnabled = true

	return ppu.pickDMGCompatibilityPalettes(rom)
}

var titleChecksums = []uint8{
	/*  $0- $F */ 0x00, 0x88, 0x16, 0x36, 0xD1, 0xDB, 0xF2, 0x3C, 0x8C, 0x92, 0x3D, 0x5C, 0x58, 0xC9, 0x3E, 0x70,
	/* $10-$1F */ 0x1D, 0x59, 0x69, 0x19, 0x35, 0xA8, 0x14, 0xAA, 0x75, 0x95, 0x99, 0x34, 0x6F, 0x15, 0xFF, 0x97,
	/* $20-$2F */ 0x4B, 0x90, 0x17, 0x10, 0x39, 0xF7, 0xF6, 0xA2, 0x49, 0x4E, 0x43, 0x68, 0xE0, 0x8B, 0xF0, 0xCE,
	/* $30-$3F */ 0x0C, 0x29, 0xE8, 0xB7, 0x86, 0x9A, 0x52, 0x01, 0x9D, 0x71, 0x9C, 0xBD, 0x5D, 0x6D, 0x67, 0x3F,
	/* $40     */ 0x6B,
	// Ambiguous (must use 4th letter)
	/* $41-$4E */ 0xB3, 0x46, 0x28, 0xA5, 0xC6, 0xD3, 0x27, 0x61, 0x18, 0x66, 0x6A, 0xBF, 0x0D, 0xF4,
}

var titleFourthLetters = []rune{
	/* Column    1    2    3    4    5    6    7    8    9    A    B    C    D    E */
	/* Row 0 */ 'B', 'E', 'F', 'A', 'A', 'R', 'B', 'E', 'K', 'E', 'K', ' ', 'R', '-',
	/* Row 1 */ 'U', 'R', 'A', 'R', ' ', 'I', 'N', 'A', 'I', 'L', 'I', 'C', 'E', ' ',
	/* Row 2 */ 'R',
}

// For each of these, the lower 5 bits indicate which row of 'paletteOffsets' to use.
// The upper 3 bits indicate how the palette shall be "shuffled":
//   - The third offset in the triplet (BGP) will always be unchanged.
//   - By default (if no bit is set), the third offset will be copied to all 3, so both OBJ palettes
//     will match BGP.
//   - If bit 0 is set, the first offset (OBP0's) will instead be unchanged.
//   - If bit 1 is set, the second offset (OBP1's) will instead use the first "source" offset.
//     (Thus, if bit 0 was also set, both OBJ palettes will be identical; if it wasn't, OBP0 will
//     match BGP, but OBP1 won't.)
//   - If bit 2 is set, the second offset (OBP1's) will instead be unchanged; this overrides bit 1 if
//     it was set.
var paletteIdsAndFlags = []uint8{
	0x7C, 0x08, 0x12, 0xA3, 0xA2, 0x07, 0x87, 0x4B, 0x20, 0x12, 0x65, 0xA8, 0x16, 0xA9, 0x86, 0xB1,
	0x68, 0xA0, 0x87, 0x66, 0x12, 0xA1, 0x30, 0x3C, 0x12, 0x85, 0x12, 0x64, 0x1B, 0x07, 0x06, 0x6F,
	0x6E, 0x6E, 0xAE, 0xAF, 0x6F, 0xB2, 0xAF, 0xB2, 0xA8, 0xAB, 0x6F, 0xAF, 0x86, 0xAE, 0xA2, 0xA2,
	0x12, 0xAF, 0x13, 0x12, 0xA1, 0x6E, 0xAF, 0xAF, 0xAD, 0x06, 0x4C, 0x6E, 0xAF, 0xAF, 0x12, 0x7C,
	0xAC,
	0xA8, 0x6A, 0x6E, 0x13, 0xA0, 0x2D, 0xA8, 0x2B, 0xAC, 0x64, 0xAC, 0x6D, 0x87, 0xBC,
	0x60, 0xB4, 0x13, 0x72, 0x7C, 0xB5, 0xAE, 0xAE, 0x7C, 0x7C, 0x65, 0xA2, 0x6C, 0x64,
	0x85,
}

var paletteOffsets = [][3]uint8{
	{16 * 8, 22 * 8, 8 * 8},
	{17 * 8, 4 * 8, 13 * 8},
	{27*8 + 6, 0 * 8, 14 * 8},
	{27*8 + 6, 4 * 8, 15 * 8},
	{4 * 8, 4 * 8, 7 * 8},
	{4 * 8, 22 * 8, 18 * 8},
	{4 * 8, 22 * 8, 20 * 8},
	{28 * 8, 22 * 8, 24 * 8},
	{19 * 8, 22*8 + 6, 9 * 8},
	{16 * 8, 28 * 8, 10 * 8},
	{3*8 + 6, 3*8 + 6, 11 * 8},
	{4 * 8, 23 * 8, 28 * 8},
	{17 * 8, 22 * 8, 2 * 8},
	{4 * 8, 0 * 8, 2 * 8},
	{4 * 8, 28 * 8, 3 * 8},
	{28 * 8, 3 * 8, 0 * 8},
	{3 * 8, 28 * 8, 4 * 8},
	{21 * 8, 28 * 8, 4 * 8},
	{3 * 8, 28 * 8, 0 * 8},
	{4 * 8, 3 * 8, 27 * 8},
	{25 * 8, 3 * 8, 28 * 8},
	{0 * 8, 28 * 8, 8 * 8},
	{5 * 8, 5 * 8, 5 * 8},
	{3 * 8, 28 * 8, 12 * 8},
	{4 * 8, 3 * 8, 28 * 8},
	{0 * 8, 0 * 8, 1 * 8},
	{28 * 8, 3 * 8, 6 * 8},
	{26 * 8, 26 * 8, 26 * 8},
	{4 * 8, 28 * 8, 29 * 8},
}

var palettes = []uint8{
	0xff, 0x7f, 0xbf, 0x32, 0xd0, 0x0, 0x0, 0x0,
	0x9f, 0x63, 0x79, 0x42, 0xb0, 0x15, 0xcb, 0x4,
	0xff, 0x7f, 0x31, 0x6e, 0x4a, 0x45, 0x0, 0x0,
	0xff, 0x7f, 0xef, 0x1b, 0x0, 0x2, 0x0, 0x0,
	0xff, 0x7f, 0x1f, 0x42, 0xf2, 0x1c, 0x0, 0x0,
	0xff, 0x7f, 0x94, 0x52, 0x4a, 0x29, 0x0, 0x0,
	0xff, 0x7f, 0xff, 0x3, 0x2f, 0x1, 0x0, 0x0,
	0xff, 0x7f, 0xef, 0x3, 0xd6, 0x1, 0x0, 0x0,
	0xff, 0x7f, 0xb5, 0x42, 0xc8, 0x3d, 0x0, 0x0,
	0x74, 0x7e, 0xff, 0x3, 0x80, 0x1, 0x0, 0x0,
	0xff, 0x67, 0xac, 0x77, 0x13, 0x1a, 0x6b, 0x2d,
	0xd6, 0x7e, 0xff, 0x4b, 0x75, 0x21, 0x0, 0x0,
	0xff, 0x53, 0x5f, 0x4a, 0x52, 0x7e, 0x0, 0x0,
	0xff, 0x4f, 0xd2, 0x7e, 0x4c, 0x3a, 0xe0, 0x1c,
	0xed, 0x3, 0xff, 0x7f, 0x5f, 0x25, 0x0, 0x0,
	0x6a, 0x3, 0x1f, 0x2, 0xff, 0x3, 0xff, 0x7f,
	0xff, 0x7f, 0xdf, 0x1, 0x12, 0x1, 0x0, 0x0,
	0x1f, 0x23, 0x5f, 0x3, 0xf2, 0x0, 0x9, 0x0,
	0xff, 0x7f, 0xea, 0x3, 0x1f, 0x1, 0x0, 0x0,
	0x9f, 0x29, 0x1a, 0x0, 0xc, 0x0, 0x0, 0x0,
	0xff, 0x7f, 0x7f, 0x2, 0x1f, 0x0, 0x0, 0x0,
	0xff, 0x7f, 0xe0, 0x3, 0x6, 0x2, 0x20, 0x1,
	0xff, 0x7f, 0xeb, 0x7e, 0x1f, 0x0, 0x0, 0x7c,
	0xff, 0x7f, 0xff, 0x3f, 0x0, 0x7e, 0x1f, 0x0,
	0xff, 0x7f, 0xff, 0x3, 0x1f, 0x0, 0x0, 0x0,
	0xff, 0x3, 0x1f, 0x0, 0xc, 0x0, 0x0, 0x0,
	0xff, 0x7f, 0x3f, 0x3, 0x93, 0x1, 0x0, 0x0,
	0x0, 0x0, 0x0, 0x42, 0x7f, 0x3, 0xff, 0x7f,
	0xff, 0x7f, 0x8c, 0x7e, 0x0, 0x7c, 0x0, 0x0,
	0xff, 0x7f, 0xef, 0x1b, 0x80, 0x61, 0x0, 0x0,
}

func findPaletteFlags(rom cartridge.Cartridge) (titleChecksum uint8, paletteIdFlags uint8) {
	// First we check that the licensee in the header is Nintendo (code 01). If not, use palette 0.
	// Check if old licensee code is 0x33 (i.e. new licensee code must be used)
	if rom.Read(0x014B) == 0x33 {
		// If new licensee code is not "01", use palettes 00
		if rom.Read(0x0144) != '0' || rom.Read(0x0145) != '1' {
			paletteIdFlags = paletteIdsAndFlags[0]
			return
		}
	} else {
		// If old licensee code is not 01, use palette 00
		if rom.Read(0x14B) != 0x01 {
			paletteIdFlags = paletteIdsAndFlags[0]
			return
		}
	}

	// In the CGB boot rom, when picking the compatibility mode palette,
	// if the old licensee code is $01, or the old licensee code is $33 and the new licensee code is "01" ($30 $31),
	// then B is the sum of all 16 title bytes.
	// Otherwise, B is $00. As indicated by the "+ 1" in the "AGB (DMG mode)" column, if on AGB, that value is increased by 12.

	// Continue algorithm

	// Compute the sum of all 16 game title bytes, storing this as the "title checksum"
	for i := uint16(0); i < 16; i++ {
		titleChecksum += rom.Read(0x134 + i)
	}

	// Find checksum index in the checksums table
	index := -1
	for i, checksum := range titleChecksums {
		if checksum == titleChecksum {
			index = i
			break
		}
	}

	// If index not found, use palette 0
	if index == -1 {
		paletteIdFlags = paletteIdsAndFlags[0]
		return
	}
	// If the index is 64 or below, the index is used as-is as the palettes ID, and the algorithm ends
	if index <= 64 {
		paletteIdFlags = paletteIdsAndFlags[index]
		return
	}

	// Otherwise the fourth letter is searched in another table: we take the index last digit ($43 -> 3)
	// and search the letter in the corresponding column: if we find it, we add the index of the letter
	// in this table (14 * row) to the previous index, otherwise we return 0
	columnIndex := (index % 0x40) - 1
	fourthLetter := rune(rom.Read(0x137))

	// First row
	if fourthLetter == titleFourthLetters[columnIndex] {
		paletteIdFlags = paletteIdsAndFlags[index]
		return
	}
	// Second row
	if fourthLetter == titleFourthLetters[14+columnIndex] {
		paletteIdFlags = paletteIdsAndFlags[index+14]
		return
	}
	// Third row
	if columnIndex == 0 && fourthLetter == titleFourthLetters[14*2+columnIndex] {
		paletteIdFlags = paletteIdsAndFlags[index+14*2]
		return
	}

	paletteIdFlags = paletteIdsAndFlags[0]
	return
}

func (ppu *PPU) pickDMGCompatibilityPalettes(rom cartridge.Cartridge) (titleChecksum uint8) {
	var idAndFlags uint8
	titleChecksum, idAndFlags = findPaletteFlags(rom)

	// The lower 5 bits indicate which row of 'paletteOffsets' to use. The upper 3 bits indicate how the palette shall be "shuffled".
	id := idAndFlags & 0x1F
	flags := idAndFlags >> 5

	// The third offset in the triplet (BGP) will always be unchanged.
	bgpIndex := paletteOffsets[id][2]

	// By default, the third offset will be copied to all 3, so both OBJ palettes will match BGP.
	obp0Index := bgpIndex
	obp1Index := bgpIndex

	// If bit 0 is set, the first offset (OBP0's) will instead be unchanged.
	if util.ReadBit(flags, 0) > 0 {
		obp0Index = paletteOffsets[id][0]
	}

	// If bit 1 is set, the second offset (OBP1's) will instead use the first "source" offset.
	// (Thus, if bit 0 was also set, both OBJ palettes will be identical; if it wasn't, OBP0 will
	//  match BGP, but OBP1 won't.)
	if util.ReadBit(flags, 1) > 0 {
		obp1Index = paletteOffsets[id][0]
	}

	// If bit 2 is set, the second offset (OBP1's) will instead be unchanged; this overrides bit 1 if it was set.
	if util.ReadBit(flags, 2) > 0 {
		obp1Index = paletteOffsets[id][1]
	}

	copy(ppu.OBJPalette[0:8], palettes[obp0Index:obp0Index+8])
	copy(ppu.OBJPalette[8:16], palettes[obp1Index:obp1Index+8])
	copy(ppu.BGPalette[0:8], palettes[bgpIndex:bgpIndex+8])

	return
}
