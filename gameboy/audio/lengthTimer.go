package audio

import "github.com/danielecanzoneri/gb-emulator/util"

type LengthTimer struct {
	// Pointer to the channel active flag
	channelEnabled *bool
	// Max value that the counter can be (64 for ch1/2/4, 256 for ch3)
	max uint

	// How many times timer must be ticked before disabling the channel
	length  uint // Bits 5-0 of NRx1 (max 64), bits 7-0 of NR31 (max 256)
	enabled bool // Bit 6 of NRx4

	// Attributes to correctly handle quirks
	frameSequencer *frameSequencer
}

func NewLengthTimer(channelEnabled *bool, frameSequencer *frameSequencer, max uint) *LengthTimer {
	return &LengthTimer{
		channelEnabled: channelEnabled,
		frameSequencer: frameSequencer,
		max:            max,
	}
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

func (lt *LengthTimer) Trigger(nrx4 uint8) {
	trigger := util.ReadBit(nrx4, 7) > 0
	enabled := util.ReadBit(nrx4, 6) > 0

	if lt.frameSequencer.ShouldClockLength() {
		// Extra length clocking occurs when writing to NRx4 when the frame sequencer's next step is one that doesn't clock the length counter.
		// In this case, if the length counter was PREVIOUSLY disabled and now enabled and the length counter is not zero, it is decremented.
		// If this decrement makes it zero and trigger is clear, the channel is disabled.
		if enabled && !lt.enabled && lt.length > 0 {
			lt.length--
			if lt.length == 0 && !trigger {
				// If triggered length counter is reset below
				*lt.channelEnabled = false
			}
		}
	}

	// If a channel is triggered when the frame sequencer's next step is one that doesn't clock the length counter
	// and the length counter is now enabled and length is being set to 64 (256 for wave channel)
	// because it was previously zero, it is set to 63 instead (255 for wave channel).
	if trigger && lt.length == 0 {
		lt.length = lt.max

		// If a channel is triggered when the frame sequencer's next step is one that doesn't clock the length counter
		// and the length counter is now enabled and length is being set to 64 (256 for wave channel)
		// because it was previously zero, it is set to 63 instead (255 for wave channel).
		if lt.frameSequencer.ShouldClockLength() && enabled {
			lt.length--
		}
	}

	lt.enabled = enabled
}
