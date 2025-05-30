package client

import (
	"fmt"
	"github.com/danielecanzoneri/gb-emulator/pkg/protocol"
)

func (c *Client) Step() error {
	return c.sendCommand(protocol.MessageTypeStep, nil)
}

func (c *Client) Continue() error {
	return c.sendCommand(protocol.MessageTypeContinue, nil)
}

func (c *Client) Pause() error {
	return c.sendCommand(protocol.MessageTypePause, nil)
}

func (c *Client) ToggleBreakpoint(address uint16, state bool) error {
	return c.sendCommand(
		protocol.MessageTypeBreakpoint,
		map[string]any{
			"address": address,
			"set":     state,
		},
	)
}

func (c *Client) sendCommand(cmdType protocol.MessageType, data map[string]any) error {
	if !c.connected {
		return fmt.Errorf("not connected to debug server")
	}

	cmd := protocol.Message{
		Type:    cmdType,
		Payload: data,
	}

	return c.conn.WriteJSON(cmd)
}
