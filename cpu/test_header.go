package cpu

// Mock implementation of Memory interface
type MockMemory struct {
	memory map[uint16]uint8
}

func (m *MockMemory) Read(addr uint16) uint8 {
	b, ok := m.memory[addr]
	if ok {
		return b
	}
	return 0
}
func (m *MockMemory) Write(addr uint16, value uint8) {
	m.memory[addr] = value
}
func (m *MockMemory) ReadWord(addr uint16) uint16 {
	return uint16(m.Read(addr)) | (uint16(m.Read(addr+1)) << 8)
}
func (m *MockMemory) WriteWord(addr uint16, value uint16) {
	m.memory[addr] = uint8(value)
	m.memory[addr+1] = uint8(value >> 8)
}

func mockCPU() *CPU {
	mem := &MockMemory{memory: make(map[uint16]uint8)}
	return New(mem)
}

func writeTestProgram(cpu *CPU, data ...byte) {
	for i, b := range data {
		cpu.Mem.Write(uint16(i)+cpu.PC, b)
	}
}
