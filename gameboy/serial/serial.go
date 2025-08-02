package serial

import (
	"net"
)

type Port struct {
	// TCP socket
	Conn net.Conn

	SB uint8
	// Serial control (bit 7: transfer enable, bit 0: clock select)
	SC uint8

	// Counter used to determine serial clock
	clockTimer int
	// Exchange one bit at a time, when all bit are exchanged, set SC bit 7 to 0 and request interrupt
	bitsTransferred int
	// Channel where data is received from socket
	dataChannel chan uint8

	RequestInterrupt func()
}

func NewPort() *Port {
	return &Port{
		// Synchronous channel
		dataChannel: make(chan uint8),
	}
}

func (port *Port) Tick(ticks uint) {
	if port.isSlave() {
		select {
		// If slave received a bit, immediately send back lower bit of SB
		case bit := <-port.dataChannel:
			port.sendBit()
			port.handleIncomingBit(bit)

		default: // Non blocking
		}

	} else {
		// Only the master determines when an exchange takes place
		if !port.isTransferring() {
			return
		}

		port.clockTimer -= int(ticks)
		if port.clockTimer <= 0 {
			port.sendBit() // Send bit to slave

			// Block until receive bit back from slave
			bitReceived := <-port.dataChannel
			port.handleIncomingBit(bitReceived)

			// Restart master clock for the next bit
			port.clockTimer += 512
		}
	}
}
