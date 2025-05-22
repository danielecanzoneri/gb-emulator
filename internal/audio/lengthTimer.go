package audio

type channel interface {
	Disable()
}

type LengthTimer struct {
	channel channel

	// How many times timer must be ticked before disabling the channel
	length uint8 // Bits 5-0 of NRx1
	// Keep count of elapsed ticks
	timer uint

	Max     uint // 64 for ch1, ch2, ch4, 256 for ch3
	Enabled bool // Bit 6 of NRx4
}

func (lt *LengthTimer) Step() {
	if !lt.Enabled {
		return
	}

	lt.timer++
	if lt.timer == lt.Max {
		lt.channel.Disable()
	}
}

func (lt *LengthTimer) Set(timer uint8) {
	lt.length = timer
}

func (lt *LengthTimer) Trigger() {
	if lt.timer == 0 {
		lt.timer = uint(lt.length)
	}
}
