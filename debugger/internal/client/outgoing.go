package client

import (
	"github.com/danielecanzoneri/gb-emulator/pkg/protocol"
	"log"
)

func (c *Client) Step() {
	c.sendCommand(protocol.MessageTypeStep, nil)
}

func (c *Client) Resume() {
	c.sendCommand(protocol.MessageTypeResume, nil)
}

func (c *Client) Continue() {
	c.sendCommand(protocol.MessageTypeContinue, nil)
}

func (c *Client) Pause() {
	c.sendCommand(protocol.MessageTypePause, nil)
}

func (c *Client) Reset() {
	c.sendCommand(protocol.MessageTypeReset, nil)
}

func (c *Client) Breakpoint(address uint16, state bool) {
	c.sendCommand(
		protocol.MessageTypeBreakpoint,
		map[string]any{
			"address": address,
			"set":     state,
		},
	)
}

func (c *Client) sendCommand(cmdType protocol.MessageType, data map[string]any) {
	if !c.connected {
		log.Println("[WARN] not connected to debug server")
		return
	}

	cmd := protocol.Message{
		Type:    cmdType,
		Payload: data,
	}

	if err := c.conn.WriteJSON(cmd); err != nil {
		log.Println("[WARN] failed to send command:", err)
	}
}
