package audio

import (
	"github.com/danielecanzoneri/gb-emulator/internal/util"
	"strconv"
)

const (
	nr10Addr = 0xFF10
	nr11Addr = 0xFF11
	nr12Addr = 0xFF12
	nr13Addr = 0xFF13
	nr14Addr = 0xFF14

	nr21Addr = 0xFF16
	nr22Addr = 0xFF17
	nr23Addr = 0xFF18
	nr24Addr = 0xFF19

	nr30Addr = 0xFF1A
	nr31Addr = 0xFF1B
	nr32Addr = 0xFF1C
	nr33Addr = 0xFF1D
	nr34Addr = 0xFF1E

	nr41Addr = 0xFF20
	nr42Addr = 0xFF21
	nr43Addr = 0xFF22
	nr44Addr = 0xFF23

	nr50Addr = 0xFF24
	nr51Addr = 0xFF25
	nr52Addr = 0xFF26

	waveRAMAddr = 0xFF30
	waveRAMSize = 16
)

func (apu *APU) IOWrite(addr uint16, v uint8) {
	// Only NR52 and wave RAM can be written when APU is disabled
	if !apu.active && !(addr == nr52Addr || addr >= waveRAMAddr) {
		return
	}

	switch addr {
	case nr10Addr:
		fallthrough
	case nr11Addr:
		fallthrough
	case nr12Addr:
		fallthrough
	case nr13Addr:
		fallthrough
	case nr14Addr:
		apu.channel1.WriteRegister(addr, v)
	case nr21Addr:
		fallthrough
	case nr22Addr:
		fallthrough
	case nr23Addr:
		fallthrough
	case nr24Addr:
		apu.channel2.WriteRegister(addr, v)
	case nr30Addr:
		fallthrough
	case nr31Addr:
		fallthrough
	case nr32Addr:
		fallthrough
	case nr33Addr:
		fallthrough
	case nr34Addr:
		apu.channel3.WriteRegister(addr, v)
	case nr41Addr:
		fallthrough
	case nr42Addr:
		fallthrough
	case nr43Addr:
		fallthrough
	case nr44Addr:
		apu.channel4.WriteRegister(addr, v)
	case nr50Addr:
		apu.nr50 = v
	case nr51Addr:
		apu.nr51 = v
	case nr52Addr:
		active := v&0x80 > 0
		if active {
			apu.enable()
		} else {
			apu.disable()
		}
	default:
		if waveRAMAddr <= addr && addr < waveRAMAddr+waveRAMSize {
			apu.WaveRam[addr-waveRAMAddr] = v
		} else {
			panic("APU: unknown addr " + strconv.FormatUint(uint64(addr), 16))
		}
	}
}

func (apu *APU) IORead(addr uint16) uint8 {
	switch addr {
	case nr10Addr:
		fallthrough
	case nr11Addr:
		fallthrough
	case nr12Addr:
		fallthrough
	case nr13Addr:
		fallthrough
	case nr14Addr:
		return apu.channel1.ReadRegister(addr)
	case nr21Addr:
		fallthrough
	case nr22Addr:
		fallthrough
	case nr23Addr:
		fallthrough
	case nr24Addr:
		return apu.channel2.ReadRegister(addr)
	case nr30Addr:
		fallthrough
	case nr31Addr:
		fallthrough
	case nr32Addr:
		fallthrough
	case nr33Addr:
		fallthrough
	case nr34Addr:
		return apu.channel3.ReadRegister(addr)
	case nr41Addr:
		fallthrough
	case nr42Addr:
		fallthrough
	case nr43Addr:
		fallthrough
	case nr44Addr:
		return apu.channel4.ReadRegister(addr)
	case nr50Addr:
		return apu.nr50
	case nr51Addr:
		return apu.nr51
	case nr52Addr:
		return apu.readNR52()
	default:
		if waveRAMAddr <= addr && addr < waveRAMAddr+waveRAMSize {
			return apu.WaveRam[addr-waveRAMAddr]
		} else {
			panic("APU: unknown addr " + strconv.FormatUint(uint64(addr), 16))
		}
	}
}

func (apu *APU) readNR52() uint8 {
	// Bits 6-4 are unused
	var nr52 uint8 = 0x70
	if apu.active {
		util.SetBit(&nr52, 7, 1)

		if apu.channel1.IsActive() {
			util.SetBit(&nr52, 0, 0)
		}
		if apu.channel2.IsActive() {
			util.SetBit(&nr52, 1, 0)
		}
		if apu.channel3.IsActive() {
			util.SetBit(&nr52, 2, 0)
		}
		if apu.channel4.IsActive() {
			util.SetBit(&nr52, 3, 0)
		}
	}
	return nr52
}

// Returns the volume for left and right between 0 and 1
func (apu *APU) getVolume() (left, right float32) {
	// A value of 0 is treated as a volume of 1 (very quiet), and a value of 7 is treated as a volume of 8 (no volume reduction).
	// Importantly, the amplifier never mutes a non-silent input.
	leftU8 := (apu.nr50 & 0x70) >> 4
	rightU8 := apu.nr50 & 0x7
	left = float32(leftU8+1) / 8
	right = float32(rightU8+1) / 8
	return
}

func (apu *APU) getLeftPanning() (ch1, ch2, ch3, ch4 bool) {
	ch1 = apu.nr51&0x10 > 0
	ch2 = apu.nr51&0x20 > 0
	ch3 = apu.nr51&0x40 > 0
	ch4 = apu.nr51&0x80 > 0
	return
}

func (apu *APU) getRightPanning() (ch1, ch2, ch3, ch4 bool) {
	ch1 = apu.nr51&0x1 > 0
	ch2 = apu.nr51&0x2 > 0
	ch3 = apu.nr51&0x4 > 0
	ch4 = apu.nr51&0x8 > 0
	return
}
