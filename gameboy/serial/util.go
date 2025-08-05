package serial

import "github.com/danielecanzoneri/gb-emulator/util"

// isTransferring checks whether a serial data exchange is happening
func (port *Port) isTransferring() bool {
	return util.ReadBit(port.SC, 7) > 0
}

// isMaster returns whether the game boy has initiated the exchange
func (port *Port) isMaster() bool {
	return util.ReadBit(port.SC, 0) > 0
}

// isSlave returns whether the game boy has requested an exchange
func (port *Port) isSlave() bool {
	return util.ReadBit(port.SC, 0) == 0
}
