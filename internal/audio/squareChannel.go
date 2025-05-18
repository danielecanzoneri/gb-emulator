package audio

import "github.com/danielecanzoneri/gb-emulator/internal/util"

type SquareChannel struct {
	dacEnabled bool
	active     bool

	// Sweep (NRx0)
	addrNRx0        uint16
	sweepPace       uint8 // Bits 6-4
	sweepDecreasing uint8 // Bit 3
	sweepStep       uint8 // Bits 2-0

	sweepCounter     uint8
	sweepEnabled     bool
	shadowRegister   uint16
	currentSweepPace uint8

	// Wave duty and length timer (NRx1)
	addrNRx1        uint16
	waveDuty        uint8 // Bits 7-6
	lengthTimerInit uint8 // Bits 5-0

	lengthTimer uint8

	// Volume and Envelope (NRx2)
	addrNRx2           uint16
	volume             uint8 // Bits 7-4
	envelopeIncreasing bool  // Bit 3
	envelopePace       uint8 // Bits 2-0

	envelopeTimer uint8
	currentVolume uint8

	// Frequency and control
	addrNRx3           uint16
	addrNRx4           uint16
	period             uint16 // Bits 2-0 of NR24 and 7-0 of NRx3 (11 bits)
	lengthTimerEnabled bool   // Bit 6 of NRx4

	periodCounter uint16
	wavePosition  uint8 // Varies from 0 to 7
}

func NewSquareChannel(addrNRx0, addrNRx1, addrNRx2, addrNRx3, addrNRx4 uint16) *SquareChannel {
	return &SquareChannel{
		addrNRx0: addrNRx0,
		addrNRx1: addrNRx1,
		addrNRx2: addrNRx2,
		addrNRx3: addrNRx3,
		addrNRx4: addrNRx4,
	}
}

func (ch *SquareChannel) IsActive() bool {
	return ch.active
}

func (ch *SquareChannel) Output() (sample float32) {
	// If a DAC is disabled, it fades to an analog value of 0, which corresponds to “digital 7.5”
	if !(ch.dacEnabled && ch.active) {
		return
	}

	// If a DAC is enabled, the digital range $0 to $F is linearly translated to the analog range -1 to 1.
	// The slope is negative: “digital 0” maps to “analog 1”, not “analog -1”.
	if waveforms[ch.waveDuty][ch.wavePosition] {
		sample = 1 - float32(ch.currentVolume)/7.5
	}
	return
}

func (ch *SquareChannel) Cycle() {
	ch.periodCounter++

	// Frequency is 11 bits
	if ch.periodCounter&0x7FF == 0 {
		ch.periodCounter = ch.period
		ch.wavePosition = (ch.wavePosition + 1) % 8
	}
}

func (ch *SquareChannel) WriteRegister(addr uint16, v uint8) {
	switch addr {
	case ch.addrNRx0:
		ch.sweepPace = (v >> 4) & 0b111
		ch.sweepDecreasing = (v >> 3) & 0b1
		ch.sweepStep = v & 0b111

	case ch.addrNRx1:
		ch.waveDuty = (v >> 6) & 0b11
		ch.lengthTimerInit = v & 0x3F

	case ch.addrNRx2:
		ch.dacEnabled = v&0xF8 > 0
		if !ch.dacEnabled {
			ch.active = false
		}

		ch.envelopePace = v & 0x7
		ch.envelopeIncreasing = v&0x8 > 0
		ch.volume = v >> 4

	case ch.addrNRx3:
		// Low 8 bits of period
		ch.period &= 0x700
		ch.period |= uint16(v)

	case ch.addrNRx4:
		ch.period = ch.period & 0xFF
		ch.period = ch.period | (uint16(v&0x7) << 8)

		ch.lengthTimerEnabled = ch.period&0x40 > 0

		// Bit 7 is trigger
		if v&0x80 > 0 {
			ch.trigger()
		}

	default:
		panic("SquareChannel: invalid address")
	}
}

func (ch *SquareChannel) ReadRegister(addr uint16) uint8 {
	switch addr {
	case ch.addrNRx0:
		return 0x80 | ch.sweepPace<<4 | ch.sweepDecreasing<<3 | ch.sweepStep

	case ch.addrNRx1:
		// Length timer is write-only
		return ch.waveDuty<<6 | 0x3F

	case ch.addrNRx2:
		out := ch.volume<<4 | ch.envelopePace
		if ch.envelopeIncreasing {
			util.SetBit(&out, 3, 1)
		}
		return out

	case ch.addrNRx3:
		// Period is write-only
		return 0xFF

	case ch.addrNRx4:
		// Only length timer can be read
		var out uint8 = 0b10111111
		if ch.lengthTimerEnabled {
			util.SetBit(&out, 6, 1)
		}
		return out

	default:
		panic("SquareChannel: invalid address")
	}
}

func (ch *SquareChannel) Disable() {
	ch.WriteRegister(ch.addrNRx0, 0)
	ch.WriteRegister(ch.addrNRx1, 0)
	ch.WriteRegister(ch.addrNRx2, 0)
	ch.WriteRegister(ch.addrNRx3, 0)
	ch.WriteRegister(ch.addrNRx4, 0)
}

func (ch *SquareChannel) trigger() {
	if !ch.dacEnabled {
		return
	}
	ch.active = true

	ch.shadowRegister = ch.period
	ch.sweepCounter = 0
	ch.sweepEnabled = ch.sweepPace != 0 || ch.sweepStep != 0

	if ch.sweepStep != 0 {
		ch.sweepOverflowCheck()
	}

	if ch.lengthTimer == 0 {
		ch.lengthTimer = ch.lengthTimerInit
	}

	ch.envelopeTimer = ch.envelopePace
	ch.currentVolume = ch.volume
}

func (ch *SquareChannel) stepSoundLength() {
	if !ch.lengthTimerEnabled || ch.lengthTimer >= 64 {
		return
	}

	ch.lengthTimer++
	if ch.lengthTimer == 64 {
		ch.active = false
	}
}

func (ch *SquareChannel) stepVolume() {
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

func (ch *SquareChannel) sweepOverflowCheck() uint16 {
	step := ch.shadowRegister >> ch.sweepStep
	newFrequency := ch.shadowRegister
	if ch.sweepDecreasing > 0 {
		newFrequency -= step
	} else {
		newFrequency += step
	}

	if newFrequency > 0x7FF {
		ch.active = false
	}

	return newFrequency
}

func (ch *SquareChannel) stepSweep() {
	if !ch.sweepEnabled || ch.sweepPace == 0 {
		return
	}

	ch.sweepCounter++
	if ch.sweepCounter == ch.currentSweepPace {
		newFrequency := ch.sweepOverflowCheck()

		ch.currentSweepPace = ch.sweepPace
		ch.sweepCounter = 0

		ch.shadowRegister = newFrequency
		ch.period = newFrequency

		// Perform another overflow check without writing it back
		ch.sweepOverflowCheck()
	}
}
