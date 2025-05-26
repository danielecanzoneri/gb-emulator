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
		apu.channel1.lengthTimer.Step()
		apu.channel2.lengthTimer.Step()
		apu.channel3.lengthTimer.Step()
		apu.channel4.lengthTimer.Step()
	}
	if (apu.audioCounter-2)%4 == 0 {
		apu.channel1.sweep.Step()
	}
	if (apu.audioCounter-7)%8 == 0 {
		apu.channel1.envelope.Step()
		apu.channel2.envelope.Step()
		apu.channel4.envelope.Step()
	}
}
