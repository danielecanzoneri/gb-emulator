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
