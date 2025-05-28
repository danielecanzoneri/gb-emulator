package audio

import (
	"github.com/danielecanzoneri/gb-emulator/pkg/util"
)

type SquareChannel struct {
	dacEnabled bool
	active     bool

	// Sweep (NRx0)
	addrNRx0 uint16
	sweep    Sweep

	// Wave duty and length timer (NRx1)
	addrNRx1    uint16
	waveDuty    uint8 // Bits 7-6
	lengthTimer LengthTimer

	// Volume and Envelope (NRx2)
	addrNRx2 uint16
	envelope Envelope

	// Frequency and control
	addrNRx3 uint16
	addrNRx4 uint16
	period   uint16 // Bits 2-0 of NR24 and 7-0 of NRx3 (11 bits)

	periodCounter uint16
	wavePosition  uint8 // Varies from 0 to 7
}

func NewSquareChannel(addrNRx0, addrNRx1, addrNRx2, addrNRx3, addrNRx4 uint16) *SquareChannel {
	ch := &SquareChannel{
		addrNRx0: addrNRx0,
		addrNRx1: addrNRx1,
		addrNRx2: addrNRx2,
		addrNRx3: addrNRx3,
		addrNRx4: addrNRx4,
	}

	ch.sweep.channel = ch
	ch.lengthTimer.channel = ch

	return ch
}

func (ch *SquareChannel) Disable() {
	ch.active = false
}

func (ch *SquareChannel) IsActive() bool {
	return ch.active
}

func (ch *SquareChannel) Output() (sample float32) {
	if !(ch.dacEnabled && ch.active) {
		return
	}

	// Each channel outputs a value between 0 and 1
	if waveforms[ch.waveDuty][ch.wavePosition] {
		sample = float32(ch.envelope.Volume()) / 15
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
		ch.sweep.WriteRegister(v)

	case ch.addrNRx1:
		ch.waveDuty = (v >> 6) & 0b11

		// TODO - check if should be ^v ^ 0x3F
		ch.lengthTimer.Set(uint(64 - v&0x3F))

	case ch.addrNRx2:
		ch.dacEnabled = v&0xF8 > 0
		if !ch.dacEnabled {
			ch.active = false
		}

		ch.envelope.WriteRegister(v)

	case ch.addrNRx3:
		// Low 8 bits of period
		ch.period &= 0x700
		ch.period |= uint16(v)

	case ch.addrNRx4:
		ch.period = ch.period & 0xFF
		ch.period = ch.period | (uint16(v&0x7) << 8)

		ch.lengthTimer.Enable(v&0x40 > 0)

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
		return ch.sweep.ReadRegister()

	case ch.addrNRx1:
		// Length timer is write-only
		return ch.waveDuty<<6 | 0x3F

	case ch.addrNRx2:
		return ch.envelope.ReadRegister()

	case ch.addrNRx3:
		// Period is write-only
		return 0xFF

	case ch.addrNRx4:
		// Only length timer can be read
		var out uint8 = 0b10111111
		if ch.lengthTimer.enabled {
			util.SetBit(&out, 6, 1)
		}
		return out

	default:
		panic("SquareChannel: invalid address")
	}
}

func (ch *SquareChannel) Reset() {
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

	ch.sweep.Trigger()
	ch.lengthTimer.Trigger()
	ch.envelope.Trigger()
}
