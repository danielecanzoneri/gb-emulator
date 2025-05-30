package client

import (
	"fmt"
	"github.com/danielecanzoneri/gb-emulator/pkg/protocol"
	"log"
)

func (c *Client) handleMessage(msg protocol.Message) {
	switch msg.Type {
	case protocol.MessageTypeState:
		state := msg.Payload
		fmt.Println(state)
	default:
		log.Println("[WARN] unknown message type:", msg.Type)
	}
}
