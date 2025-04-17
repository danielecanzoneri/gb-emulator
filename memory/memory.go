package memory

const MemorySize = 0x10000 // 64KB

type Memory struct {
	data [MemorySize]uint8
}

func (m *Memory) Read(addr uint16) uint8 {
	return m.data[addr]
}

func (m *Memory) Write(addr uint16, value uint8) {
	m.data[addr] = value
}

func (m *Memory) ReadWord(addr uint16) uint16 {
	return uint16(m.data[addr]) | (uint16(m.data[addr+1]) << 8)
}

func (m *Memory) WriteWord(addr uint16, value uint16) {
	m.data[addr] = uint8(value)
	m.data[addr+1] = uint8(value >> 8)
}
