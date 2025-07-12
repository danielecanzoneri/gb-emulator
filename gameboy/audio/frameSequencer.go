package audio

type frameSequencer struct {
	position uint8
}

func (fs *frameSequencer) ShouldClockLength() bool {
	return fs.position%2 == 0
}

func (fs *frameSequencer) ShouldClockSweep() bool {
	return (fs.position-2)%4 == 0
}

func (fs *frameSequencer) ShouldClockEnvelope() bool {
	return (fs.position-7)%8 == 0
}

func (fs *frameSequencer) Step() (clockLength, clockSweep, clockEnvelope bool) {
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
	fs.position++

	clockLength = fs.ShouldClockLength()
	clockSweep = fs.ShouldClockSweep()
	clockEnvelope = fs.ShouldClockEnvelope()
	return
}

func (apu *APU) StepFrameSequencer() {
	clockLength, clockSweep, clockEnvelope := apu.frameSequencer.Step()

	if clockLength {
		apu.channel1.lengthTimer.Step()
		apu.channel2.lengthTimer.Step()
		apu.channel3.lengthTimer.Step()
		apu.channel4.lengthTimer.Step()
	}
	if clockSweep {
		apu.channel1.sweep.Step()
	}
	if clockEnvelope {
		apu.channel1.envelope.Step()
		apu.channel2.envelope.Step()
		apu.channel4.envelope.Step()
	}
}
