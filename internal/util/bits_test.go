package util

import "testing"

func TestCombineBytes(t *testing.T) {
	var high, low uint8 = 0x12, 0x34
	word := CombineBytes(high, low)

	if word != 0x1234 {
		t.Errorf("CombineBytes(): high=%02X, low=%02X -> expected 0x1234, got %04X", high, low, word)
	}
}

func TestSplitWord(t *testing.T) {
	var word uint16 = 0x1234
	high, low := SplitWord(word)

	if high != 0x12 {
		t.Errorf("SplitWord(): word=%04X -> high: expected 0x12, got %02X", word, high)
	}
	if low != 0x34 {
		t.Errorf("SplitWord(): word=%04X -> low: expected 0x34, got %02X", word, low)
	}
}

func TestCheckBit(t *testing.T) {
	var b uint8 = 0b00001111
	bit := uint8(3)

	if ReadBit(b, bit) != 1 {
		t.Errorf("checkBit(): b=%08b, bit=%d -> expected true, got false", b, bit)
	}

	bit = 4
	if ReadBit(b, bit) != 0 {
		t.Errorf("checkBit(): b=%08b, bit=%d -> expected false, got true", b, bit)
	}
}

func TestSetBit(t *testing.T) {
	var b uint8 = 0b00001111
	bit := uint8(4)
	value := uint8(1)

	SetBit(&b, bit, value)
	if b != 0b00011111 {
		t.Errorf("SetBit(): b=%08b, bit=%d, value=%v -> expected 0b00011111, got %08b", b, bit, value, b)
	}

	value = uint8(0)
	SetBit(&b, bit, value)
	if b != 0b00001111 {
		t.Errorf("SetBit(): b=%08b, bit=%d, value=%v -> expected 0b00001111, got %08b", b, bit, value, b)
	}
}
func TestSumWordsWithCarry(t *testing.T) {
	n1, n2 := uint16(0x00FF), uint16(0x0001)
	sum, carry, halfCarry := SumWordsWithCarry(n1, n2)

	if sum != n1+n2 {
		t.Errorf("SumWordsWithCarry(): n1=%04X, n2=%04X -> expected sum=%04X, got %04X", n1, n2, n1+n2, sum)
	}
	if carry != 0 {
		t.Errorf("SumWordsWithCarry(): n1=%04X, n2=%04X -> expected carry=0, got %d", n1, n2, carry)
	}
	if halfCarry != 0 {
		t.Errorf("SumWordsWithCarry(): n1=%04X, n2=%04X -> expected halfCarry=0, got %d", n1, n2, halfCarry)
	}

	n1, n2 = uint16(0x0FFF), uint16(0x0001)
	sum, carry, halfCarry = SumWordsWithCarry(n1, n2)

	if sum != n1+n2 {
		t.Errorf("SumWordsWithCarry(): n1=%04X, n2=%04X -> expected sum=%04X, got %04X", n1, n2, n1+n2, sum)
	}
	if carry != 0 {
		t.Errorf("SumWordsWithCarry(): n1=%04X, n2=%04X -> expected carry=0, got %d", n1, n2, carry)
	}
	if halfCarry != 1 {
		t.Errorf("SumWordsWithCarry(): n1=%04X, n2=%04X -> expected halfCarry=1, got %d", n1, n2, halfCarry)
	}

	n1, n2 = uint16(0xFFFF), uint16(0x0001)
	sum, carry, halfCarry = SumWordsWithCarry(n1, n2)

	if sum != n1+n2 {
		t.Errorf("SumWordsWithCarry(): n1=%04X, n2=%04X -> expected sum=%04X, got %04X", n1, n2, n1+n2, sum)
	}
	if carry != 1 {
		t.Errorf("SumWordsWithCarry(): n1=%04X, n2=%04X -> expected carry=1, got %d", n1, n2, carry)
	}
	if halfCarry != 1 {
		t.Errorf("SumWordsWithCarry(): n1=%04X, n2=%04X -> expected halfCarry=1, got %d", n1, n2, halfCarry)
	}
}

func TestSumBytesWithCarry(t *testing.T) {
	n1, n2 := uint8(0x08), uint8(0x01)
	sum, carry, halfCarry := SumBytesWithCarry(n1, n2)

	if sum != n1+n2 {
		t.Errorf("SumBytesWithCarry(): n1=%02X, n2=%02X -> expected sum=%02X, got %02X", n1, n2, n1+n2, sum)
	}
	if carry != 0 {
		t.Errorf("SumBytesWithCarry(): n1=%02X, n2=%02X -> expected carry=0, got %d", n1, n2, carry)
	}
	if halfCarry != 0 {
		t.Errorf("SumBytesWithCarry(): n1=%02X, n2=%02X -> expected halfCarry=0, got %d", n1, n2, halfCarry)
	}

	n1, n2 = uint8(0x0F), uint8(0x01)
	sum, carry, halfCarry = SumBytesWithCarry(n1, n2)

	if sum != n1+n2 {
		t.Errorf("SumBytesWithCarry(): n1=%02X, n2=%02X -> expected sum=%02X, got %02X", n1, n2, n1+n2, sum)
	}
	if carry != 0 {
		t.Errorf("SumBytesWithCarry(): n1=%02X, n2=%02X -> expected carry=0, got %d", n1, n2, carry)
	}
	if halfCarry != 1 {
		t.Errorf("SumBytesWithCarry(): n1=%02X, n2=%02X -> expected halfCarry=1, got %d", n1, n2, halfCarry)
	}

	n1, n2 = uint8(0xFF), uint8(0x01)
	sum, carry, halfCarry = SumBytesWithCarry(n1, n2)

	if sum != n1+n2 {
		t.Errorf("SumBytesWithCarry(): n1=%02X, n2=%02X -> expected sum=%02X, got %02X", n1, n2, n1+n2, sum)
	}
	if carry != 1 {
		t.Errorf("SumBytesWithCarry(): n1=%02X, n2=%02X -> expected carry=1, got %d", n1, n2, carry)
	}
	if halfCarry != 1 {
		t.Errorf("SumBytesWithCarry(): n1=%02X, n2=%02X -> expected halfCarry=1, got %d", n1, n2, halfCarry)
	}
}

func TestIsByteZeroUint8(t *testing.T) {
	if IsByteZeroUint8(0) != 1 {
		t.Fatal("0: got 0, expected 1")
	}
	for i := 1; i <= 0xFF; i++ {
		if IsByteZeroUint8(uint8(i)) != 0 {
			t.Fatalf("%02X: got 1, expected 0", i)
		}
	}
}

func TestIsWordZeroUint8(t *testing.T) {
	if IsWordZeroUint8(0) != 1 {
		t.Fatal("0: got 0, expected 1")
	}
	for i := 1; i <= 0xFFFF; i++ {
		if IsWordZeroUint8(uint16(i)) != 0 {
			t.Fatalf("%04X: got 1, expected 0", i)
		}
	}
}

func TestSubBytesWithCarry(t *testing.T) {
	tests := map[string]struct {
		n1   uint8
		n2   uint8
		expH uint8
		expC uint8
	}{
		"C0-H0": {n1: 0xFF, n2: 0x0F, expC: 0, expH: 0},
		"C1-H1": {n1: 0x00, n2: 0x01, expC: 1, expH: 1},
		"C1-H0": {n1: 0x0F, n2: 0x10, expC: 1, expH: 0},
		"C0-H1": {n1: 0x80, n2: 0x01, expC: 0, expH: 1},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			sub, carry, halfCarry := SubBytesWithCarry(test.n1, test.n2)

			if sub != test.n1-test.n2 {
				t.Errorf("n1=%02X, n2=%02X -> expected sub=%02X, got %02X", test.n1, test.n2, test.n1-test.n2, sub)
			}
			if carry != test.expC {
				t.Errorf("n1=%02X, n2=%02X -> expected carry=%x, got %d", test.n1, test.n2, carry, test.expC)
			}
			if halfCarry != test.expH {
				t.Errorf("n1=%02X, n2=%02X -> expected halfCarry=%x, got %d", test.n1, test.n2, halfCarry, test.expH)
			}
		})
	}
}

func TestByteToBCD(t *testing.T) {
	tests := map[string]struct {
		n        uint8
		expected uint8
	}{
		"0":  {n: 0, expected: 0x00},
		"19": {n: 19, expected: 0x19},
		"99": {n: 99, expected: 0x99},
		"8":  {n: 8, expected: 0x8},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			if ByteToBCD(test.n) != test.expected {
				t.Fatalf("got %02X, expected %02X", ByteToBCD(test.n), test.expected)
			}
		})
	}
}
