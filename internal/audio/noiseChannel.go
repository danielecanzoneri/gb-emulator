package audio

import (
	"github.com/danielecanzoneri/gb-emulator/internal/util"
)

type NoiseChannel struct {
	dacEnabled bool
	active     bool

	lfsr             uint16 // State of LFSR
	frequencyCounter uint16 // It steps LFSR when 0

	// Wave duty and length timer (NR41)
	lengthTimerInit uint8 // Bits 5-0 (write only)

	lengthTimer uint8

	// Volume and Envelope (NR42)
	volume             uint8 // Bits 7-4
	envelopeIncreasing bool  // Bit 3
	envelopePace       uint8 // Bits 2-0

	envelopeTimer uint8
	currentVolume uint8

	// Frequency & randomness (NR43)
	clockShift   uint8 // Bits 7-4
	lfsrWidth    uint8 // Bit 3 (0 = 15-bit, 1 = 7-bit)
	clockDivider uint8 // Bits 2-0

	// Control (NR44)
	lengthTimerEnabled bool // Bit 6
}

func NewNoiseChannel() *NoiseChannel {
	ch := new(NoiseChannel)
	ch.resetFrequency()
	return ch
}

func (ch *NoiseChannel) resetFrequency() {
	// The frequency at which the LFSR is clocked is
	//     262144 / (divider * 2^shift) Hz
	// that means every 4194304 / 262144 * (divider * 2^shift) ticks.
	// Divider = 0 is treated as divider = 0.5

	ch.frequencyCounter = 4 * (2 << ch.clockShift)
	if ch.clockDivider == 0 {
		ch.frequencyCounter >>= 1 // Divide by 2 = multiply 0.5
	} else {
		ch.frequencyCounter *= uint16(ch.clockDivider)
	}
}

func (ch *NoiseChannel) IsActive() bool {
	return ch.active
}

func (ch *NoiseChannel) Output() (sample float32) {
	// If a DAC is disabled, it fades to an analog value of 0, which corresponds to “digital 7.5”
	if !(ch.dacEnabled && ch.active) {
		return
	}

	if ch.lfsr&0b1 > 0 {
		sample = 1 - float32(ch.currentVolume)/7.5
	}
	return
}

func (ch *NoiseChannel) Cycle() {
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

func (ch *NoiseChannel) WriteRegister(addr uint16, v uint8) {
	switch addr {
	case nr41Addr:
		ch.lengthTimerInit = v & 0x3F

	case nr42Addr:
		ch.dacEnabled = v&0xF8 > 0
		if !ch.dacEnabled {
			ch.active = false
		}

		ch.envelopePace = v & 0x7
		ch.envelopeIncreasing = v&0x8 > 0
		ch.volume = v >> 4

	case nr43Addr:
		ch.clockShift = v >> 4
		ch.lfsrWidth = (v >> 3) & 0b1
		ch.clockDivider = v & 0b111

	case nr44Addr:
		ch.lengthTimerEnabled = v&0x40 > 0

		// Bit 7 is trigger
		if v&0x80 > 0 {
			ch.trigger()
		}

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
		out := ch.volume<<4 | ch.envelopePace
		if ch.envelopeIncreasing {
			util.SetBit(&out, 3, 1)
		}
		return out

	case nr43Addr:
		return ch.clockShift<<4 | ch.lfsrWidth<<3 | ch.clockDivider

	case nr44Addr:
		// Only length timer can be read
		var out uint8 = 0b10111111
		if ch.lengthTimerEnabled {
			util.SetBit(&out, 6, 1)
		}
		return out

	default:
		panic("NoiseChannel: invalid address")
	}
}

func (ch *NoiseChannel) Disable() {
	ch.WriteRegister(nr41Addr, 0)
	ch.WriteRegister(nr42Addr, 0)
	ch.WriteRegister(nr43Addr, 0)
	ch.WriteRegister(nr44Addr, 0)
}

func (ch *NoiseChannel) trigger() {
	if !ch.dacEnabled {
		return
	}
	ch.active = true

	// Reset LFSR bits
	ch.lfsr = 0

	if ch.lengthTimer == 0 {
		ch.lengthTimer = ch.lengthTimerInit
	}

	ch.envelopeTimer = ch.envelopePace
	ch.currentVolume = ch.volume
}

func (ch *NoiseChannel) stepSoundLength() {
	if !ch.lengthTimerEnabled || ch.lengthTimer >= 64 {
		return
	}

	ch.lengthTimer++
	if ch.lengthTimer == 64 {
		ch.active = false
	}
}

func (ch *NoiseChannel) stepVolume() {
	if ch.envelopePace == 0 {
		// Envelope disabled
		return
	}

	ch.envelopeTimer--
	if ch.envelopeTimer == 0 {
		ch.envelopeTimer = ch.envelopePace

		// The digital value produced by the generator ranges between $0 and $F
		if ch.envelopeIncreasing && ch.currentVolume < 0xF {
			ch.currentVolume++
		} else if !ch.envelopeIncreasing && ch.currentVolume > 0 {
			ch.currentVolume--
		}
	}
}
