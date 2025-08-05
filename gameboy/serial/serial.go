package serial

import (
	"net"
)

type LinkState int

const (
	Disconnected LinkState = iota
	Connecting
	Connected
)

type Port struct {
	// TCP socket
	Conn net.Conn
	// Connection state
	State LinkState

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
		// It seems that at startup actual Game Boy timer has elapsed for eight ticks (check Timer)
		clockTimer: 512 - 8,
	}
}

func (port *Port) Tick(ticks int) {
	// Serial clock runs at 8 kHz, since game boy runs at 4 MHz
	// each serial clock happens once every 4 MHz / 8 kHz = 512 game boy ticks
	// Note that serial clock is always running even when not transmitting data
	port.clockTimer -= ticks
	if port.clockTimer <= 0 {
		// Restart master clock for the next bit
		port.clockTimer += 512

		if port.isTransferring() && port.isMaster() {
			if port.State == Connected {
				port.sendBit() // Send bit to slave

				// Block until a bit is received back from slave
				bitReceived := <-port.dataChannel
				port.handleIncomingBit(bitReceived)
			} else {
				// Emulate disconnected cable
				port.handleIncomingBit(1)
			}
		}
	}

	if port.isSlave() {
		select {
		// If slave received a bit, immediately send back lower bit of SB
		case bit := <-port.dataChannel:
			port.sendBit()
			port.handleIncomingBit(bit)

		default: // Non blocking
		}
	}
}
