package audio

import (
	"github.com/danielecanzoneri/gb-emulator/pkg/util"
)

type WaveChannel struct {
	// Bit 7 of NR30
	dacEnabled bool
	active     bool

	// Length timer (NR31)
	lengthTimer LengthTimer

	// Volume (NR32)
	volume uint8 // Bits 6-5

	// Frequency and control (NR33 & NR34)
	period        uint16 // Bits 2-0 of NR34 and 7-0 of NR33 (11 bits)
	periodCounter uint16
	// wavePosition ranges between 0 and 31
	// where 0 is the upper nibble of FF30 and 31 is the lower nibble of FF3F
	wavePosition uint8

	// Wave RAM
	WaveRam [16]uint8
}

func NewWaveChannel() *WaveChannel {
	ch := new(WaveChannel)
	ch.lengthTimer.channel = ch

	return ch
}

func (ch *WaveChannel) Disable() {
	ch.active = false
}

func (ch *WaveChannel) IsActive() bool {
	return ch.active
}

func (ch *WaveChannel) Output() (sample float32) {
	if !(ch.dacEnabled && ch.active) {
		return
	}

	if ch.volume == 0 {
		return
	}

	// 0 is the upper nibble of FF30 and 31 is the lower nibble of FF3F
	nibble := ch.WaveRam[ch.wavePosition>>1]
	nibble = nibble >> (4 * (1 - ch.wavePosition&1)) // If 0 take upper nibble
	nibble &= 0xF

	// volume has the following meaning
	// 00	Mute (No sound)
	// 01	100% volume (use samples read from Wave RAM as-is)
	// 10	50% volume (shift samples read from Wave RAM right once)
	// 11	25% volume (shift samples read from Wave RAM right twice)
	nibble >>= ch.volume - 1
	sample = float32(nibble) / 15
	return
}

func (ch *WaveChannel) Cycle() {
	// The wave channel’s period divider is clocked once per two dots (twice per cycle)
	for range 2 {
		ch.periodCounter++

		// Frequency is 11 bits
		if ch.periodCounter&0x7FF == 0 {
			ch.periodCounter = ch.period
			ch.wavePosition = (ch.wavePosition + 1) % 32
		}
	}
}

func (ch *WaveChannel) WriteRegister(addr uint16, v uint8) {
	switch addr {
	case nr30Addr:
		ch.dacEnabled = util.ReadBit(v, 7) > 0
		if !ch.dacEnabled {
			ch.active = false
		}

	case nr31Addr:
		// Timer should count up to 256, but we make it count down to 0
		// TODO - check if should be ^v
		ch.lengthTimer.Set(256 - uint(v))

	case nr32Addr:
		ch.volume = (v >> 5) & 0b11

	case nr33Addr:
		// Low 8 bits of period
		ch.period &= 0x700
		ch.period |= uint16(v)

	case nr34Addr:
		ch.period = ch.period & 0xFF
		ch.period = ch.period | (uint16(v&0x7) << 8)

		ch.lengthTimer.Enable(v&0x40 > 0)

		// Bit 7 is trigger
		if v&0x80 > 0 {
			ch.trigger()
		}

	default:
		panic("WaveChannel: invalid address")
	}
}

func (ch *WaveChannel) ReadRegister(addr uint16) uint8 {
	switch addr {
	case nr30Addr:
		var out uint8 = 0xFF
		if !ch.dacEnabled {
			util.SetBit(&out, 7, 0)
		}
		return out

	case nr31Addr:
		// Length timer is write-only
		return 0xFF

	case nr32Addr:
		return 0x9F | (ch.volume << 5)

	case nr33Addr:
		// Period is write-only
		return 0xFF

	case nr34Addr:
		// Only length timer can be read
		var out uint8 = 0b10111111
		if ch.lengthTimer.enabled {
			util.SetBit(&out, 6, 1)
		}
		return out

	default:
		panic("WaveChannel: invalid address")
	}
}

func (ch *WaveChannel) WriteWRAM(addr uint16, v uint8) {
	if ch.active {
		return
	}
	ch.WaveRam[addr-waveRAMAddr] = v
}
func (ch *WaveChannel) ReadWRAM(addr uint16) uint8 {
	if ch.active {
		return 0xFF
	}
	return ch.WaveRam[addr-waveRAMAddr]
}

func (ch *WaveChannel) Reset() {
	ch.WriteRegister(nr30Addr, 0)
	ch.WriteRegister(nr31Addr, 0)
	ch.WriteRegister(nr32Addr, 0)
	ch.WriteRegister(nr33Addr, 0)
	ch.WriteRegister(nr34Addr, 0)
}

func (ch *WaveChannel) trigger() {
	if !ch.dacEnabled {
		return
	}
	ch.active = true

	ch.periodCounter = ch.period
	ch.lengthTimer.Trigger()
}
