package audio

type Disabler interface {
	Disable()
}

type LengthTimer struct {
	channel Disabler

	// How many times timer must be ticked before disabling the channel
	length uint8 // Bits 5-0 of NRx1
	// Keep count of elapsed ticks
	timer uint8

	Enabled bool // Bit 6 of NRx4
}

func (lt *LengthTimer) Step() {
	if !lt.Enabled || lt.timer >= 64 {
		return
	}

	lt.timer++
	if lt.timer == 64 {
		lt.channel.Disable()
	}
}

func (lt *LengthTimer) Set(timer uint8) {
	lt.length = timer & 0x3F
}

func (lt *LengthTimer) Trigger() {
	if lt.timer == 0 {
		lt.timer = lt.length
	}
}
