package serial

import (
	"github.com/danielecanzoneri/gb-emulator/util"
	"io"
	"log"
)

// Listen to incoming bytes from the other port
func (port *Port) Listen() {
	go func() {
		buf := make([]byte, 1)

		for {
			_, err := port.Conn.Read(buf)
			log.Printf("Received: %08b\n", buf[0])

			switch {
			case err == nil: // Do nothing
			case err == io.EOF:
				return
			default:
				log.Println("Connection error:", err)
				continue
			}

			// Send back the byte to the master
			if port.isSlave() {
				port.send()
			}

			// Take only the first bit and shift everything to the left
			port.SB = buf[0]

			// Both master and slave
			util.SetBit(&port.SC, 7, 0)
			port.RequestInterrupt()
		}
	}()
}

// Send SB to the other port
func (port *Port) send() {
	_, err := port.Conn.Write([]byte{port.SB})

	if err != nil {
		log.Println("Connection error:", err)
	}

	log.Println("Sent ", port.SB)
}
