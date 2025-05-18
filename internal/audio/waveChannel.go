package audio

type WaveChannel struct {
}

func (ch *WaveChannel) IsActive() bool {
	return false
}

func (ch *WaveChannel) Output() (sample float32) {
	return
}

func (ch *WaveChannel) Cycle() {
}

func (ch *WaveChannel) WriteRegister(addr uint16, v uint8) {
}

func (ch *WaveChannel) ReadRegister(addr uint16) uint8 {
	return 0
}

func (ch *WaveChannel) Disable() {
}
