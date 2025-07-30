package serial

import (
	"net"
)

const (
	SBAddr = 0xFF01
	SCAddr = 0xFF02

	SCMask = 0x7E
)

type Port struct {
	// TCP socket
	Conn net.Conn

	SB uint8
	// Serial control (bit 7: transfer enable, bit 0: clock select)
	SC uint8

	// Counter used to determine serial clock
	clockCounter int
	// Exchange one bit at a time, when all bit are exchanged, set SC bit 7 to 0 and request interrupt
	exchangedBit int

	RequestInterrupt func()
}

func (port *Port) Tick(ticks uint) {
	// Only the master determines when an exchange takes place
	if !port.isMaster() {
		return
	}

	// Serial clock runs at 8 kHz, since game boy runs at 4 MHz
	// each serial clock happens once every 4 MHz / 8 kHz = 512 game boy ticks
	port.clockCounter += int(ticks)
	if port.clockCounter >= 8*512 {
		port.clockCounter = 0
		port.send()
	}
}
