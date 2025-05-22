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

	// Counter
	audioCounter uint8

	// Buffer to store samples
	sampleRate    float64
	sampleCounter float64 // Counter used to produce samples at correct rate
	sampleBuffer  chan float32
}

func NewAPU(sampleRate float64, sampleBuffer chan float32) *APU {
	return &APU{
		sampleRate:   sampleRate,
		channel1:     NewSquareChannel(nr10Addr, nr11Addr, nr12Addr, nr13Addr, nr14Addr),
		channel2:     NewSquareChannel(0, nr21Addr, nr22Addr, nr23Addr, nr24Addr),
		channel3:     NewWaveChannel(),
		channel4:     NewNoiseChannel(),
		sampleBuffer: sampleBuffer,
	}
}

func (apu *APU) Cycle() {
	apu.channel1.Cycle()
	apu.channel2.Cycle()
	apu.channel3.Cycle()
	apu.channel4.Cycle()

	apu.sampleCounter++
	var cyclesPerSample = 4194304. / (apu.sampleRate * 4)
	if apu.sampleCounter >= cyclesPerSample {
		apu.sampleCounter -= cyclesPerSample

		left, right := apu.sample()
		apu.sampleBuffer <- left
		apu.sampleBuffer <- right
	}
}

func (apu *APU) sample() (left, right float32) {
	leftCh1, leftCh2, leftCh3, leftCh4 := apu.getLeftPanning()
	rightCh1, rightCh2, rightCh3, rightCh4 := apu.getRightPanning()
	leftVolume, rightVolume := apu.getVolume()

	if leftCh1 {
		left += apu.channel1.Output()
	}
	if leftCh2 {
		left += apu.channel2.Output()
	}
	if leftCh3 {
		left += apu.channel3.Output()
	}
	if leftCh4 {
		left += apu.channel4.Output()
	}
	left = (left / 4) * leftVolume

	if rightCh1 {
		right += apu.channel1.Output()
	}
	if rightCh2 {
		right += apu.channel2.Output()
	}
	if rightCh3 {
		right += apu.channel3.Output()
	}
	if rightCh4 {
		right += apu.channel4.Output()
	}
	right = (right / 4) * rightVolume

	return
}

func (apu *APU) enable() {
	if apu.active {
		return
	}

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
	apu.channel1.Reset()
	apu.channel2.Reset()
	apu.channel3.Reset()
	apu.channel4.Reset()
}
