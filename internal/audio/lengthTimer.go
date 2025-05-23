package audio

type channel interface {
	Disable()
}

type LengthTimer struct {
	channel channel

	// How many times timer must be ticked before disabling the channel
	length uint // Bits 5-0 of NRx1 (max 64), bits 7-0 of NR31 (max 256)
	// Keep count of elapsed ticks
	timer uint

	enabled bool // Bit 6 of NRx4
}

func (lt *LengthTimer) Step() {
	if !lt.enabled {
		return
	}

	lt.timer--
	if lt.timer == 0 {
		lt.channel.Disable()
		lt.enabled = false
	}
}

func (lt *LengthTimer) Set(timer uint) {
	lt.length = timer
}

func (lt *LengthTimer) Enable(enabled bool) {
	if enabled {
		lt.enabled = true
		lt.timer = lt.length
	} else {
		lt.enabled = false
	}
}

func (lt *LengthTimer) Trigger() {
	if lt.timer == 0 {
		lt.timer = lt.length
	}
}
