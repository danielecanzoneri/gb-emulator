package client

import (
	"encoding/json"
	"github.com/danielecanzoneri/gb-emulator/pkg/protocol"
	"log"
)

func (c *Client) handleMessage(message []byte) {
	var msg protocol.Message
	if err := json.Unmarshal(message, &msg); err != nil {
		log.Printf("Error parsing message: %v\n", err)
		return
	}

	switch msg.Type {
	case protocol.MessageTypeState:
		var stateMsg protocol.StateMessage
		if err := json.Unmarshal(message, &stateMsg); err != nil {
			log.Printf("Error parsing state message: %v\n", err)
		}
		state := stateMsg.Payload

		// Consume emulator state
		c.OnState(&state)
	case protocol.MessageTypeBreakpointHit:
		c.OnBreakpointHit()
	default:
		log.Println("[WARN] unknown message type:", msg.Type)
	}
}
