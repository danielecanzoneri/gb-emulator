package cpu

func combineBytes(high, low uint8) uint16 {
	return (uint16(high) << 8) | uint16(low)
}

func splitWord(word uint16) (high, low uint8) {
	high = uint8(word >> 8)
	low = uint8(word)
	return
}

func readBit(b uint8, bit uint8) uint8 {
	return (b >> bit) & 1
}

func setBit(b *uint8, bit uint8, value uint8) {
	*b = *b & ^(1 << bit)
	*b = *b | ((value & 1) << bit)
}

func sumWordsWithCarry(n1, n2 uint16) (sum uint16, carry, half_carry uint8) {
	s := uint32(n1) + uint32(n2)
	carry = uint8(s >> 16)
	half_carry = uint8((n1&0xFFF + n2&0xFFF) >> 12)
	sum = uint16(s)
	return
}

func sumBytesWithCarry(n1, n2 uint8) (sum uint8, carry, half_carry uint8) {
	s, carry, half_carry := sumWordsWithCarry(uint16(n1)<<8, uint16(n2)<<8)
	sum = uint8(s >> 8)
	return
}

func isByteZeroUint8(n uint8) uint8 {
	return 1 - (n|(-n))>>7
}

func isWordZeroUint8(n uint16) uint8 {
	return 1 - uint8((n|(-n))>>15)
}

func subBytesWithCarry(n1, n2 uint8) (sub uint8, carry, half_carry uint8) {
	sub = n1 - n2
	borrows := ((^n1 & n2) | ((^n1 | n2) & sub))
	carry = borrows >> 7
	half_carry = (borrows >> 3) & 0x1
	return
}

func byteToBCD(n uint8) uint8 {
	low_nibble := n % 10
	high_nibble := n / 10
	return high_nibble<<4 | low_nibble
}
