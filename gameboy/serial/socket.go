package serial

import (
	"github.com/danielecanzoneri/gb-emulator/util"
	"io"
	"log"
)

// Listen to incoming bytes
func (port *Port) Listen() {
	go func() {
		buf := make([]uint8, 1)

		for {
			_, err := port.Conn.Read(buf)

			switch {
			case err == nil: // Do nothing
			case err == io.EOF:
				return
			default:
				log.Println("Connection error:", err)
				continue
			}

			port.dataChannel <- buf[0]
		}
	}()
}

func (port *Port) handleIncomingBit(bitIn uint8) {
	// Set bit 0 of SB
	port.SB = (port.SB << 1) | (bitIn & 1)
	port.bitsTransferred++

	if port.bitsTransferred == 8 {
		port.bitsTransferred = 0

		// Disable transferring and request interrupt
		util.SetBit(&port.SC, 7, 0)
		port.RequestInterrupt()
	}
}

// Send bit 7 of SB to the other port
func (port *Port) sendBit() {
	bitToSend := util.ReadBit(port.SB, 7)
	_, err := port.Conn.Write([]uint8{bitToSend})

	if err != nil {
		log.Println("Connection error:", err)
	}
}
