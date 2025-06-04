package client

import (
	"fmt"
	"github.com/danielecanzoneri/gb-emulator/pkg/debug"
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	conn      *websocket.Conn
	serverURL string
	connected bool

	// Function that will consume incoming state update from emulator
	OnState         func(*debug.GameBoyState)
	OnBreakpointHit func()
}

func New(host string, port int) *Client {
	serverURL := fmt.Sprintf("ws://%s:%d/ws", host, port)
	return &Client{
		serverURL: serverURL,
	}
}

func (c *Client) Connect() error {
	if c.connected {
		log.Println("Already connected")
		return nil
	}

	conn, _, err := websocket.DefaultDialer.Dial(c.serverURL, nil)
	if err != nil {
		return fmt.Errorf("failed to connect to debug server: %v", err)
	}

	c.conn = conn
	c.connected = true

	// Pause execution when connecting
	c.Pause()

	go c.listen()
	return nil
}

func (c *Client) Disconnect() error {
	if !c.connected {
		return nil
	}

	err := c.conn.Close()
	if err != nil {
		return fmt.Errorf("error closing connection: %v", err)
	}

	c.connected = false
	c.conn = nil
	return nil
}

func (c *Client) listen() {
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v\n", err)
			}
			break
		}

		c.handleMessage(message)
	}
}
