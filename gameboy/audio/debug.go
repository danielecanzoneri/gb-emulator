package audio

import "strconv"

func (apu *APU) DebugRead(addr uint16) uint8 {
	if nr10Addr <= addr && addr <= nr14Addr {
		return apu.channel1.ReadRegister(addr)
	} else if nr21Addr <= addr && addr <= nr24Addr {
		return apu.channel2.ReadRegister(addr)
	} else if nr30Addr <= addr && addr <= nr34Addr {
		return apu.channel3.ReadRegister(addr)
	} else if nr41Addr <= addr && addr <= nr44Addr {
		return apu.channel4.ReadRegister(addr)
	} else if waveRAMAddr <= addr && addr < waveRAMAddr+waveRAMSize {
		return apu.channel3.WaveRam[addr-waveRAMAddr]
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
