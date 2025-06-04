package protocol

import "github.com/danielecanzoneri/gb-emulator/pkg/debug"

// Message represents a generic message between emulator and debugger
type Message struct {
	Type    MessageType    `json:"type"`
	Payload map[string]any `json:"payload,omitempty"`
}

type StateMessage struct {
	Type    MessageType        `json:"type"`
	Payload debug.GameBoyState `json:"payload,omitempty"`
}

type MessageType string

// Debugger to emulator
const (
	MessageTypePause      MessageType = "pause"
	MessageTypeResume     MessageType = "resume"
	MessageTypeStep       MessageType = "step"
	MessageTypeContinue   MessageType = "continue"
	MessageTypeBreakpoint MessageType = "breakpoint"
	MessageTypeReset      MessageType = "reset"
)

// Emulator to debugger
const (
	MessageTypeBreakpointHit MessageType = "breakpoint_hit"
	MessageTypeState         MessageType = "state"
)
