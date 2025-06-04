package server

import (
	"github.com/danielecanzoneri/gb-emulator/pkg/protocol"
	"log"
)

func (s *Server) handleCommand(cmd protocol.Message) {
	switch cmd.Type {
	case protocol.MessageTypePause:
		s.debugger.Pause()
		s.sendState()
	case protocol.MessageTypeStep:
		s.debugger.Step()
		s.sendState()
	case protocol.MessageTypeContinue:
		s.debugger.Continue()
	case protocol.MessageTypeBreakpoint:
		payload := cmd.Payload
		addr := uint16(payload["address"].(float64))
		set := payload["set"].(bool)
		s.debugger.Breakpoint(addr, set)
	}
}

func (s *Server) sendState() {
	state := s.debugger.GetState()
	message := protocol.StateMessage{
		Type:    protocol.MessageTypeState,
		Payload: *state,
	}

	if err := s.client.WriteJSON(message); err != nil {
		log.Println("[WARN] failed to send state to client:", err)
	}
}
