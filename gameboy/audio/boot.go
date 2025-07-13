package audio

func (apu *APU) SkipBoot() {
	// Channel 1
	apu.channel1.dacEnabled = true
	apu.channel1.active = true
	apu.channel1.sweep.timer = 5
	apu.channel1.sweep.shadow = 0x7C1
	apu.channel1.waveDuty = 2
	apu.channel1.lengthTimer.length = 64
	apu.channel1.envelope.volumeInit = 0xF
	apu.channel1.envelope.pace = 3
	apu.channel1.envelope.timer = 1
	apu.channel1.period = 0x7C1
	apu.channel1.periodCounter = 0x7F2
	apu.channel1.wavePosition = 3

	apu.nr50 = 0x77
	apu.nr51 = 0xF3
	apu.active = true
	apu.frameSequencer.position = 24
}
