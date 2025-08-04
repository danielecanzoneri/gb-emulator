package audio

import (
	"github.com/danielecanzoneri/gb-emulator/util"
)

type NoiseChannel struct {
	dacEnabled bool
	active     bool

	lfsr             uint16 // State of LFSR
	frequencyCounter uint16 // It steps LFSR when 0

	// Length timer (NR41)
	lengthTimer *LengthTimer

	// Volume and Envelope (NR42)
	envelope Envelope

	// Frequency & randomness (NR43)
	clockShift   uint8 // Bits 7-4
	lfsrWidth    uint8 // Bit 3 (0 = 15-bit, 1 = 7-bit)
	clockDivider uint8 // Bits 2-0

	ticks int
}

func NewNoiseChannel(fs *frameSequencer) *NoiseChannel {
	ch := new(NoiseChannel)
	ch.lengthTimer = NewLengthTimer(&ch.active, fs, 64)

	ch.resetFrequency()
	return ch
}

func (ch *NoiseChannel) resetFrequency() {
	// The frequency at which the LFSR is clocked is
	//     262144 / (divider * 2^shift) Hz
	// that means every 4194304 / 262144 * (divider * 2^shift) ticks.
	// Divider = 0 is treated as divider = 0.5

	ch.frequencyCounter = 4
	if ch.clockDivider == 0 {
		ch.frequencyCounter >>= 1 // Divide by 2 = multiply 0.5
	} else {
		ch.frequencyCounter *= uint16(ch.clockDivider)
	}
	ch.frequencyCounter <<= ch.clockShift
}

func (ch *NoiseChannel) IsActive() bool {
	return ch.active
}

func (ch *NoiseChannel) Output() (sample float32) {
	if !(ch.dacEnabled && ch.active) {
		return
	}

	sample = float32(ch.lfsr&0b1) * float32(ch.envelope.Volume()) / 15
	return
}

func (ch *NoiseChannel) Tick(ticks int) {
	if !ch.active {
		return
	}

	ch.ticks += ticks

	// Channel 4 clocks at 1048576 HZ
	for ch.ticks >= 4 {
		ch.ticks -= 4

		ch.frequencyCounter--
		if ch.frequencyCounter == 0 {
			// Advance LFSR
			bit15 := ^(ch.lfsr ^ (ch.lfsr >> 1)) & 0b1
			util.SetBit(&ch.lfsr, 15, bit15)

			if ch.lfsrWidth == 1 { // 7-bit mode
				util.SetBit(&ch.lfsr, 7, bit15)
			}

			ch.lfsr >>= 1
			ch.resetFrequency()
		}
	}
}

func (ch *NoiseChannel) WriteRegister(addr uint16, v uint8) {
	switch addr {
	case nr41Addr:
		ch.lengthTimer.Set(uint(64 - v&0x3F))

	case nr42Addr:
		ch.dacEnabled = v&0xF8 > 0
		if !ch.dacEnabled {
			ch.active = false
		}

		ch.envelope.WriteRegister(v)

	case nr43Addr:
		ch.clockShift = v >> 4
		ch.lfsrWidth = (v >> 3) & 0b1
		ch.clockDivider = v & 0b111

	case nr44Addr:
		ch.Trigger(v)

	default:
		panic("NoiseChannel: invalid address")
	}
}

func (ch *NoiseChannel) ReadRegister(addr uint16) uint8 {
	switch addr {
	case nr41Addr:
		// Length timer is write-only
		return 0xFF

	case nr42Addr:
		return ch.envelope.ReadRegister()

	case nr43Addr:
		return ch.clockShift<<4 | ch.lfsrWidth<<3 | ch.clockDivider

	case nr44Addr:
		// Only length timer can be read
		var out uint8 = 0b10111111
		if ch.lengthTimer.enabled {
			util.SetBit(&out, 6, 1)
		}
		return out

	default:
		panic("NoiseChannel: invalid address")
	}
}

func (ch *NoiseChannel) disable() {
	// On the DMG, length counters are unaffected by power
	// ch.WriteRegister(nr41Addr, 0)
	ch.WriteRegister(nr42Addr, 0)
	ch.WriteRegister(nr43Addr, 0)
	ch.WriteRegister(nr44Addr, 0)

	ch.ticks = 0
}

func (ch *NoiseChannel) Trigger(value uint8) {
	ch.lengthTimer.Trigger(value)

	trigger := util.ReadBit(value, 7) > 0
	if ch.dacEnabled && trigger {
		// Active channel only if DAC is enabled
		ch.active = true

		ch.lfsr = 0
		ch.envelope.Trigger()
	}
}
