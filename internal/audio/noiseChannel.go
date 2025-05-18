package audio

type NoiseChannel struct {
}

func (ch *NoiseChannel) IsActive() bool {
	return false
}

func (ch *NoiseChannel) Output() (sample float32) {
	return
}

func (ch *NoiseChannel) Cycle() {
}

func (ch *NoiseChannel) WriteRegister(addr uint16, v uint8) {
}

func (ch *NoiseChannel) ReadRegister(addr uint16) uint8 {
	return 0
}

func (ch *NoiseChannel) Disable() {
}
