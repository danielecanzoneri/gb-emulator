package audio

import (
	"github.com/danielecanzoneri/gb-emulator/pkg/util"
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

	if nr10Addr <= addr && addr <= nr14Addr {
		apu.channel1.WriteRegister(addr, v)
	} else if nr21Addr <= addr && addr <= nr24Addr {
		apu.channel2.WriteRegister(addr, v)
	} else if nr30Addr <= addr && addr <= nr34Addr {
		apu.channel3.WriteRegister(addr, v)
	} else if nr41Addr <= addr && addr <= nr44Addr {
		apu.channel4.WriteRegister(addr, v)
	} else if waveRAMAddr <= addr && addr < waveRAMAddr+waveRAMSize {
		apu.channel3.WriteWRAM(addr, v)
	} else {
		switch addr {
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
			panic("APU: unknown addr " + strconv.FormatUint(uint64(addr), 16))
		}
	}
}

func (apu *APU) IORead(addr uint16) uint8 {
	if nr10Addr <= addr && addr <= nr14Addr {
		return apu.channel1.ReadRegister(addr)
	} else if nr21Addr <= addr && addr <= nr24Addr {
		return apu.channel2.ReadRegister(addr)
	} else if nr30Addr <= addr && addr <= nr34Addr {
		return apu.channel3.ReadRegister(addr)
	} else if nr41Addr <= addr && addr <= nr44Addr {
		return apu.channel4.ReadRegister(addr)
	} else if waveRAMAddr <= addr && addr < waveRAMAddr+waveRAMSize {
		return apu.channel3.ReadWRAM(addr)
	} else {
		switch addr {
		case nr50Addr:
			return apu.nr50
		case nr51Addr:
			return apu.nr51
		case nr52Addr:
			return apu.readNR52()
		default:
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
			util.SetBit(&nr52, 0, 1)
		}
		if apu.channel2.IsActive() {
			util.SetBit(&nr52, 1, 1)
		}
		if apu.channel3.IsActive() {
			util.SetBit(&nr52, 2, 1)
		}
		if apu.channel4.IsActive() {
			util.SetBit(&nr52, 3, 1)
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
