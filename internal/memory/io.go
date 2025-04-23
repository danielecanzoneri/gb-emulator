package memory

import "fmt"

const (
	ioRegisters = 0xFF00
	hRAM        = 0xFF80
)

const (
	divAddr  = 0xFF04
	timaAddr = 0xFF05
	tmaAddr  = 0xFF06
	tacAddr  = 0xFF07
)

func (mmu *MMU) WriteIO(addr uint16, v uint8) {
	switch addr {
	// Timer I/O
	case divAddr:
		fallthrough
	case timaAddr:
		fallthrough
	case tmaAddr:
		fallthrough
	case tacAddr:
		mmu.Timer.Write(addr, v)

	default:
		mmu.data[addr] = v
	}

	// Debug on serial port
	if addr == 0xFF01 {
		if v == 0 {
			fmt.Println()
		} else {
			fmt.Printf("%c", v)
		}
	}
}

func (mmu *MMU) ReadIO(addr uint16) uint8 {
	switch addr {
	// Timer I/O
	case divAddr:
		fallthrough
	case timaAddr:
		fallthrough
	case tmaAddr:
		fallthrough
	case tacAddr:
		return mmu.Timer.Read(addr)

	default:
		return mmu.data[addr]
	}
}
