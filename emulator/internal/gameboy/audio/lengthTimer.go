package audio

type LengthTimer struct {
	// Pointer to the channel active flag
	channelEnabled *bool

	// How many times timer must be ticked before disabling the channel
	length uint // Bits 5-0 of NRx1 (max 64), bits 7-0 of NR31 (max 256)

	enabled bool // Bit 6 of NRx4
}

func (lt *LengthTimer) Step() {
	if !lt.enabled {
		return
	}

	if lt.length > 0 {
		lt.length--

		if lt.length == 0 {
			*lt.channelEnabled = false
		}
	}
}

func (lt *LengthTimer) Set(timer uint) {
	lt.length = timer
}

func (lt *LengthTimer) Enable(enabled bool) {
	lt.enabled = enabled
}

func (lt *LengthTimer) Trigger(max uint) {
	if lt.length == 0 {
		lt.length = max
	}
}
