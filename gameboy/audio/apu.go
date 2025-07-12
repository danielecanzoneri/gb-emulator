package audio

type APU struct {
	channel1 *SquareChannel
	channel2 *SquareChannel
	channel3 *WaveChannel
	channel4 *NoiseChannel

	// Master (NR50, NR51, NR52)
	nr50   uint8 // bits 6-4 left volume, bits 2-0 right volume
	nr51   uint8 // bits 7-4 left panning, bits 3-0 right panning
	active bool  // Bit 7 of NR52

	// FrameSequencer
	frameSequencer frameSequencer

	// Buffer to store samples
	sampleRate    float64
	sampleCounter float64 // Counter used to produce samples at correct rate
	sampleBuffer  chan float32

	// Control manually channels
	Ch1Enabled bool
	Ch2Enabled bool
	Ch3Enabled bool
	Ch4Enabled bool
}

func NewAPU(sampleRate float64, sampleBuffer chan float32) *APU {
	apu := &APU{
		sampleRate:   sampleRate,
		sampleBuffer: sampleBuffer,
		Ch1Enabled:   true, Ch2Enabled: true, Ch3Enabled: true, Ch4Enabled: true,
	}
	apu.channel1 = NewSquareChannel(nr10Addr, nr11Addr, nr12Addr, nr13Addr, nr14Addr, &apu.frameSequencer)
	apu.channel2 = NewSquareChannel(0, nr21Addr, nr22Addr, nr23Addr, nr24Addr, &apu.frameSequencer)
	apu.channel3 = NewWaveChannel(&apu.frameSequencer)
	apu.channel4 = NewNoiseChannel(&apu.frameSequencer)

	return apu
}

func (apu *APU) Tick(ticks uint) {
	if !apu.active {
		return
	}

	apu.channel1.Tick(ticks)
	apu.channel2.Tick(ticks)
	apu.channel3.Tick(ticks)
	apu.channel4.Tick(ticks)

	apu.sampleCounter += float64(ticks)
	var ticksPerSample = 4194304. / apu.sampleRate
	for apu.sampleCounter >= ticksPerSample {
		apu.sampleCounter -= ticksPerSample

		left, right := apu.sample()
		apu.sampleBuffer <- left
		apu.sampleBuffer <- right
	}
}

func (apu *APU) sample() (left, right float32) {
	if !apu.active {
		return
	}

	leftCh1, leftCh2, leftCh3, leftCh4 := apu.getLeftPanning()
	rightCh1, rightCh2, rightCh3, rightCh4 := apu.getRightPanning()
	leftVolume, rightVolume := apu.getVolume()

	if leftCh1 && apu.Ch1Enabled {
		left += apu.channel1.Output()
	}
	if leftCh2 && apu.Ch2Enabled {
		left += apu.channel2.Output()
	}
	if leftCh3 && apu.Ch3Enabled {
		left += apu.channel3.Output()
	}
	if leftCh4 && apu.Ch4Enabled {
		left += apu.channel4.Output()
	}
	left = (left / 4) * leftVolume

	if rightCh1 && apu.Ch1Enabled {
		right += apu.channel1.Output()
	}
	if rightCh2 && apu.Ch2Enabled {
		right += apu.channel2.Output()
	}
	if rightCh3 && apu.Ch3Enabled {
		right += apu.channel3.Output()
	}
	if rightCh4 && apu.Ch4Enabled {
		right += apu.channel4.Output()
	}
	right = (right / 4) * rightVolume

	return
}

func (apu *APU) enable() {
	if apu.active {
		return
	}

	// When powered on, the frame sequencer is reset so that the next step will be 0
	apu.frameSequencer.position = 7

	// Set bit 7 of NR52
	apu.active = true
}

func (apu *APU) disable() {
	if !apu.active {
		return
	}

	apu.nr50 = 0
	apu.nr51 = 0
	// Reset bit 7 of NR52
	apu.active = false

	// Reset all registers except NR52 and wave RAM
	apu.channel1.disable()
	apu.channel2.disable()
	apu.channel3.disable()
	apu.channel4.disable()
}
