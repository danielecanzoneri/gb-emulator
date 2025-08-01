package serial

import (
	"net"
	"sync"
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
	// Start clock only when master receives a bit from slave
	waitingForReply bool

	packetCount int
	lock        sync.Mutex

	RequestInterrupt func()
}

func (port *Port) Tick(ticks uint) {
	port.lock.Lock()
	defer port.lock.Unlock()

	// Only the master determines when an exchange takes place
	if !port.isTransferring() || !port.isMaster() {
		return
	}

	port.clockTimer -= int(ticks)
	if port.clockTimer <= 0 {
		// If master hasn't received bit from slave, wait
		if !port.waitingForReply {
			port.sendBit()
			port.waitingForReply = true

			// Restart master clock for the next bit
			port.clockTimer = 512
		}
	}
}
