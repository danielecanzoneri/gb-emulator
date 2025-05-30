package server

import (
	"github.com/danielecanzoneri/gb-emulator/pkg/protocol"
)

func (s *Server) handleCommand(cmd protocol.Message) {
	switch cmd.Type {
	case protocol.MessageTypePause:
		s.debugger.Pause()
	case protocol.MessageTypeStep:
		s.debugger.Step()
	case protocol.MessageTypeContinue:
		s.debugger.Continue()
	case protocol.MessageTypeBreakpoint:
		addr := cmd.Payload["address"].(uint16)
		set := cmd.Payload["set"].(bool)
		s.debugger.Breakpoint(addr, set)
	}
}
