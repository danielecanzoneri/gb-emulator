package serial

import "github.com/danielecanzoneri/gb-emulator/util"

// isActive checks whether a serial data exchange is happening
func (port *Port) isActive() bool {
	return port.Conn != nil && util.ReadBit(port.SC, 7) > 0
}

// isMaster returns whether the game boy has initiated the exchange
func (port *Port) isMaster() bool {
	return port.isActive() && util.ReadBit(port.SC, 0) > 0
}

// isSlave returns whether the game boy has requested an exchange
func (port *Port) isSlave() bool {
	return port.isActive() && util.ReadBit(port.SC, 0) == 0
}
