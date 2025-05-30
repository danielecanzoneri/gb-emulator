package debug

type GameBoyState struct {
	Memory [0x10000]uint8 `json:"memory"`
	AF     uint16         `json:"AF"`
	BC     uint16         `json:"BC"`
	DE     uint16         `json:"DE"`
	HL     uint16         `json:"HL"`
	PC     uint16         `json:"PC"`
	SP     uint16         `json:"SP"`
	IME    bool           `json:"IME"`
}

func (state *GameBoyState) GetMap() map[string]any {
	return map[string]any{
		"memory": state.Memory,
		"AF":     state.AF,
		"BC":     state.BC,
		"DE":     state.DE,
		"HL":     state.HL,
		"PC":     state.PC,
		"SP":     state.SP,
		"IME":    state.IME,
	}
}
