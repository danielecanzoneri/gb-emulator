package util

func CombineBytes(high, low uint8) uint16 {
	return (uint16(high) << 8) | uint16(low)
}

func SplitWord(word uint16) (high, low uint8) {
	high = uint8(word >> 8)
	low = uint8(word)
	return
}

func ReadBit(b uint8, bit uint8) uint8 {
	return (b >> bit) & 1
}

func SetBit(b *uint8, bit uint8, value uint8) {
	*b = *b & ^(1 << bit)
	*b = *b | ((value & 1) << bit)
}

func SumWordsWithCarry(n1, n2 uint16) (sum uint16, carry, halfCarry uint8) {
	s := uint32(n1) + uint32(n2)
	carry = uint8(s >> 16)
	halfCarry = uint8((n1&0xFFF + n2&0xFFF) >> 12)
	sum = uint16(s)
	return
}

func SumBytesWithCarry(n1, n2 uint8) (sum uint8, carry, halfCarry uint8) {
	s, carry, halfCarry := SumWordsWithCarry(uint16(n1)<<8, uint16(n2)<<8)
	sum = uint8(s >> 8)
	return
}

func IsByteZeroUint8(n uint8) uint8 {
	return 1 - (n|(-n))>>7
}

func IsWordZeroUint8(n uint16) uint8 {
	return 1 - uint8((n|(-n))>>15)
}

func SubBytesWithCarry(n1, n2 uint8) (sub uint8, carry, halfCarry uint8) {
	sub = n1 - n2
	borrows := (^n1 & n2) | ((^n1 | n2) & sub)
	carry = borrows >> 7
	halfCarry = (borrows >> 3) & 0x1
	return
}

func ByteToBCD(n uint8) uint8 {
	lowNibble := n % 10
	highNibble := n / 10
	return highNibble<<4 | lowNibble
}

func SpreadBits(b uint8) uint16 {
	x := uint16(b)
	x = (x | (x << 4)) & 0x0F0F
	x = (x | (x << 2)) & 0x3333
	x = (x | (x << 1)) & 0x5555
	return x
}
