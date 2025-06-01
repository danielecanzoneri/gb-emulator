package server

import (
	"github.com/danielecanzoneri/gb-emulator/emulator/internal/debugger"
	"github.com/danielecanzoneri/gb-emulator/pkg/protocol"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Server struct {
	debugger *debugger.Debugger

	client *websocket.Conn // The debugging client
}

func New(debugger *debugger.Debugger) *Server {
	return &Server{
		debugger: debugger,
	}
}

func (s *Server) Start(port string) {
	http.HandleFunc("/ws", s.handleConnection)
	log.Printf("Debug server starting on port %s\n", port)

	go func() {
		if err := http.ListenAndServe(":"+port, nil); err != nil {
			log.Fatal(err)
		}
	}()
}

func (s *Server) Close() {
	s.closeClient()
}

func (s *Server) handleConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v\n", err)
		return
	}

	if s.client != nil {
		log.Printf("Client already connected\n")
		return
	}
	s.client = conn

	go s.clientHandler(conn)
}

func (s *Server) clientHandler(conn *websocket.Conn) {
	defer func() {
		// Delete client when disconnecting
		s.client = nil

		if err := conn.Close(); err != nil {
			log.Printf("WebSocket connection close error: %v\n", err)
		}
	}()

	// Main loop: handle commands received from the client
	for {
		var cmd protocol.Message
		err := conn.ReadJSON(&cmd)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v\n", err)
			}
			s.closeClient()
			break
		}

		s.handleCommand(cmd)
	}
}

func (s *Server) closeClient() {
	if s.client == nil {
		return
	}

	s.debugger.Resume()

	if err := s.client.Close(); err != nil {
		log.Printf("Error closing client: %v\n", err)
	}
}
