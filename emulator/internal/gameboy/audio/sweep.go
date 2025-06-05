package audio

// Sweep manages the NRx0 register
type Sweep struct {
	channel *SquareChannel

	// How many times sweep must be ticked before frequency increasing/decreasing
	pace uint8 // Bits 6-4
	// If frequency will increase or decrease
	isDecreasing uint8 // Bit 3
	// Quantity that modify the frequency at each step
	step uint8 // Bits 2-0

	// Keep count of elapsed ticks
	timer uint8
	// Controls if sweep is active
	enabled bool

	// Current output period (inverse of frequency)
	shadow uint16
}

func (sw *Sweep) checkOverflow() uint16 {
	step := sw.shadow >> sw.step
	newFrequency := sw.shadow
	if sw.isDecreasing > 0 {
		newFrequency -= step
	} else {
		newFrequency += step
	}

	if newFrequency > 0x7FF {
		sw.channel.active = false
	}

	return newFrequency
}

func (sw *Sweep) resetTimer() {
	// Sweep timers treat a period of 0 as 8
	if sw.pace == 0 {
		sw.timer = 8
	} else {
		sw.timer = sw.pace
	}
}

func (sw *Sweep) Step() {
	sw.timer--
	if sw.timer == 0 {
		sw.resetTimer()

		if sw.enabled && sw.pace != 0 {
			newFrequency := sw.checkOverflow()

			if newFrequency <= 0x7FF && sw.step != 0 {
				sw.shadow = newFrequency
				sw.channel.period = newFrequency

				// Perform another overflow check without writing it back
				sw.checkOverflow()
			}
		}
	}
}

func (sw *Sweep) WriteRegister(v uint8) {
	sw.pace = (v >> 4) & 0b111
	sw.isDecreasing = (v >> 3) & 0b1
	sw.step = v & 0b111
}

func (sw *Sweep) ReadRegister() uint8 {
	return 0x80 | sw.pace<<4 | sw.isDecreasing<<3 | sw.step
}

func (sw *Sweep) Trigger() {
	sw.shadow = sw.channel.period
	sw.resetTimer()
	sw.enabled = sw.pace != 0 || sw.step != 0

	if sw.step != 0 {
		sw.checkOverflow()
	}
}
