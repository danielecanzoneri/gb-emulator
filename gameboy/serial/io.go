package serial

import (
	"strconv"
)

func (port *Port) Read(addr uint16) uint8 {
	switch addr {
	case SBAddr:
		return port.SB
	case SCAddr:
		return port.SC | SCMask
	default:
		panic("Serial: unknown addr " + strconv.FormatUint(uint64(addr), 16))
	}
}

func (port *Port) Write(addr uint16, v uint8) {
	switch addr {
	case SBAddr:
		port.SB = v
	case SCAddr:
		port.SC = v
	default:
		panic("Serial: unknown addr " + strconv.FormatUint(uint64(addr), 16))
	}
}
