package audio

import (
	"github.com/danielecanzoneri/gb-emulator/util"
)

const (
	// https://forums.nesdev.org/viewtopic.php?p=188035&sid=c7368fbb8c83fb2f5245902eb4ff5791#p188035
	triggerWaveCycleDelay = 3 // 3 APU cycles -> 6 ticks
)

type WaveChannel struct {
	// Bit 7 of NR30
	dacEnabled bool
	active     bool

	// CGB flag
	isCGB bool

	// Length timer (NR31)
	LengthTimer *LengthTimer

	// Volume (NR32)
	volume uint8 // Bits 6-5

	// Frequency and control (NR33 & NR34)
	period        uint16 // Bits 2-0 of NR34 and 7-0 of NR33 (11 bits)
	periodCounter uint16
	// wavePosition ranges between 0 and 31 where 0 is the upper nibble of FF30 and 31 is the lower nibble of FF3F
	wavePosition      uint8
	justRead          bool // When reading Wave RAM while on, it will return the byte just read in the buffer if access happen in same cycle
	triggerCycleDelay int  // Delay before advancing wavePosition after triggering

	// Wave RAM
	WaveRam      [16]uint8
	bufferSample uint8 // CH3 does not emit samples directly, but stores every sample read into a buffer, and emits that continuously;

	ticks int
}

func NewWaveChannel(fs *frameSequencer, isCGB bool) *WaveChannel {
	ch := &WaveChannel{
		isCGB: isCGB,
	}
	ch.LengthTimer = NewLengthTimer(&ch.active, fs, 256)

	return ch
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

	// volume has the following meaning
	// 00	Mute (No sound)
	// 01	100% volume (use samples read from Wave RAM as-is)
	// 10	50% volume (shift samples read from Wave RAM right once)
	// 11	25% volume (shift samples read from Wave RAM right twice)
	currSample := ch.bufferSample >> (ch.volume - 1)
	sample = float32(currSample) / 15
	return
}

func (ch *WaveChannel) Tick(ticks int) {
	if !ch.active {
		return
	}

	ch.ticks += ticks

	// Channel 3 clocks at 2097152 HZ
	for ch.ticks >= 2 {
		ch.ticks -= 2

		if ch.triggerCycleDelay > 0 {
			ch.triggerCycleDelay--
			continue
		}
		ch.justRead = false
		ch.periodCounter++

		// Frequency is 11 bits
		if ch.periodCounter&0x7FF == 0 {
			ch.periodCounter = ch.period
			ch.wavePosition = (ch.wavePosition + 1) % 32

			// 0 is the upper nibble of FF30 and 31 is the lower nibble of FF3F
			ch.bufferSample = ch.WaveRam[ch.wavePosition>>1]
			ch.bufferSample >>= 4 * (1 - ch.wavePosition&1) // If 0 take upper nibble
			ch.bufferSample &= 0xF
			ch.justRead = true
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
		ch.LengthTimer.Set(256 - uint(v))

	case nr32Addr:
		ch.volume = (v >> 5) & 0b11

	case nr33Addr:
		// Low 8 bits of period
		ch.period &= 0x700
		ch.period |= uint16(v)

	case nr34Addr:
		ch.period = ch.period & 0xFF
		ch.period = ch.period | (uint16(v&0x7) << 8)

		ch.Trigger(v)

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
		if ch.LengthTimer.enabled {
			util.SetBit(&out, 6, 1)
		}
		return out

	default:
		panic("WaveChannel: invalid address")
	}
}

func (ch *WaveChannel) WriteWRAM(addr uint16, v uint8) {
	if ch.active {
		if ch.isCGB || /* isDMG && */ ch.justRead {
			ch.WaveRam[ch.wavePosition>>1] = v
		}
		return
	}
	ch.WaveRam[addr-waveRAMAddr] = v
}

func (ch *WaveChannel) ReadWRAM(addr uint16) uint8 {
	// If the wave channel is enabled, accessing any byte from $FF30-$FF3F is equivalent to accessing the current byte selected by the waveform position.
	// Furthermore, on DMG accesses will only work in this manner if made within a couple of clocks of the wave channel accessing wave RAM;
	// if made at any other time, reads return $FF and writes have no effect. ON CGB, instead, it can be accessed at any time
	if ch.active {
		if ch.isCGB || /* isDMG && */ ch.justRead {
			return ch.WaveRam[ch.wavePosition>>1]
		}
		return 0xFF
	}
	return ch.WaveRam[addr-waveRAMAddr]
}

func (ch *WaveChannel) disable() {
	ch.WriteRegister(nr30Addr, 0)
	// On the DMG, length counters are unaffected by power
	// ch.WriteRegister(nr31Addr, 0)
	ch.WriteRegister(nr32Addr, 0)
	ch.WriteRegister(nr33Addr, 0)
	ch.WriteRegister(nr34Addr, 0)
	ch.bufferSample = 0
	ch.ticks = 0
}

func (ch *WaveChannel) Trigger(value uint8) {
	ch.LengthTimer.Trigger(value)

	trigger := util.ReadBit(value, 7) > 0
	if trigger && ch.dacEnabled {
		// Active channel only if DAC is enabled
		ch.active = true

		// Triggering the wave channel on the DMG while it reads a sample byte will alter the first four bytes of wave RAM.
		if !ch.isCGB && ch.periodCounter == ch.period+1 { // Don't actually know why this works, but it works so GG
			indexOfNextByte := ((ch.wavePosition + 1) % 32) >> 1
			// If the channel was reading one of the first four bytes, the only first byte will be rewritten with the byte being read.
			if indexOfNextByte < 4 {
				ch.WaveRam[0] = ch.WaveRam[indexOfNextByte]
			} else {
				// If the channel was reading one of the later 12 bytes, the first FOUR bytes of wave RAM will be rewritten
				// with the four aligned bytes that the read was from (bytes 4-7, 8-11, or 12-15);
				// for example if it were reading byte 9 when< it was retriggered, the first four bytes would be rewritten with the contents of bytes 8-11.
				copy(ch.WaveRam[:], ch.WaveRam[indexOfNextByte&0xFC:indexOfNextByte&0xFC+4])
			}
		}

		ch.periodCounter = ch.period
		ch.wavePosition = 0
		ch.triggerCycleDelay = triggerWaveCycleDelay
	}
}
