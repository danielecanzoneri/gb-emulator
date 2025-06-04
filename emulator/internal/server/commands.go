package server

import (
	"github.com/danielecanzoneri/gb-emulator/pkg/protocol"
	"log"
)

func (s *Server) handleCommand(cmd protocol.Message) {
	switch cmd.Type {
	case protocol.MessageTypePause:
		s.Pause()
		s.sendState()
	case protocol.MessageTypeStep:
		s.Step()
		s.sendState()
	case protocol.MessageTypeContinue:
		s.Continue()
	case protocol.MessageTypeBreakpoint:
		payload := cmd.Payload
		addr := uint16(payload["address"].(float64))
		set := payload["set"].(bool)
		s.Breakpoint(addr, set)
	case protocol.MessageTypeReset:
		s.Reset()
		s.sendState()
	}
}

func (s *Server) sendBreakpointHit() {
	message := protocol.Message{
		Type: protocol.MessageTypeBreakpointHit,
	}
	if err := s.client.WriteJSON(message); err != nil {
		log.Println("[WARN] failed to send breakpoint hit to client:", err)
	}
}

func (s *Server) sendState() {
	state := s.GetState()
	message := protocol.StateMessage{
		Type:    protocol.MessageTypeState,
		Payload: *state,
	}

	if err := s.client.WriteJSON(message); err != nil {
		log.Println("[WARN] failed to send state to client:", err)
	}
}
