package audio

func (apu *APU) StepCounter() {
	// Rate   256 Hz      64 Hz       128 Hz
	// ---------------------------------------
	// 7      -           Clock       -
	// 6      Clock       -           Clock
	// 5      -           -           -
	// 4      Clock       -           -
	// 3      -           -           -
	// 2      Clock       -           Clock
	// 1      -           -           -
	// 0      Clock       -           -
	// ---------------------------------------
	// Step   Length Ctr  Vol Env     Sweep
	apu.audioCounter++

	if apu.audioCounter%2 == 0 {
		apu.channel1.stepSoundLength()
		apu.channel2.stepSoundLength()
	}
	if (apu.audioCounter-2)%4 == 0 {
		apu.channel1.stepSweep()
	}
	if (apu.audioCounter-7)%8 == 0 {
		apu.channel1.stepVolume()
		apu.channel2.stepVolume()
	}
}
