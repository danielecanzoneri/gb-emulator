package audio

import "github.com/danielecanzoneri/gb-emulator/gameboy/util"

// Envelope manages the NRx2 register
type Envelope struct {
	// How loud the channel initially is
	volumeInit uint8 // Bits 7-4
	// Envelope direction
	isIncreasing bool // Bit 3
	// How many times envelope must be ticked before increasing/decreasing
	pace uint8 // Bits 2-0

	// Keep count of elapsed ticks
	timer uint8
	// Current volume
	volume uint8
}

func (e *Envelope) Step() {
	if e.pace == 0 {
		// Envelope disabled
		return
	}

	e.timer--
	if e.timer == 0 {
		e.timer = e.pace

		// The digital value produced by the generator ranges between $0 and $F
		if e.isIncreasing && e.volume < 0xF {
			e.volume++
		} else if !e.isIncreasing && e.volume > 0 {
			e.volume--
		}
	}
}

func (e *Envelope) WriteRegister(v uint8) {
	e.pace = v & 0x7
	e.isIncreasing = util.ReadBit(v, 3) > 0
	e.volumeInit = v >> 4
}

func (e *Envelope) ReadRegister() uint8 {
	out := e.volumeInit<<4 | e.pace
	if e.isIncreasing {
		util.SetBit(&out, 3, 1)
	}
	return out
}

func (e *Envelope) Volume() uint8 {
	return e.volume
}

func (e *Envelope) Trigger() {
	e.timer = e.pace
	e.volume = e.volumeInit
}
