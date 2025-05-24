package main

import (
	"encoding/binary"
	"github.com/ebitengine/oto/v3"
	"io"
	"math"
)

const (
	sampleRate = 44100
	channels   = 2
	format     = oto.FormatFloat32LE

	bufferSize = 8192
)

type AudioPlayer struct {
	gb *GameBoy

	sampleBuffer chan float32
}

func newAudioPlayer(gb *GameBoy, sampleBuffer chan float32) (*oto.Player, error) {
	op := &oto.NewContextOptions{}
	op.SampleRate = sampleRate
	op.ChannelCount = channels
	op.Format = format

	ctx, ready, err := oto.NewContext(op)
	if err != nil {
		return nil, err
	}

	<-ready // Context ready

	a := &AudioPlayer{
		gb:           gb,
		sampleBuffer: sampleBuffer,
	}

	p := ctx.NewPlayer(a)
	p.SetBufferSize(bufferSize)

	return p, nil
}

func (a *AudioPlayer) Read(buf []byte) (n int, err error) {
	// If Game Boy is paused return silence and don't execute cpu instructions
	if a.gb.paused {
		for i := range len(buf) {
			buf[i] = 0
		}
		return len(buf), nil
	}

	bufferPosition := 0

	for bufferPosition < len(buf) {
		// If not enough samples have been produced, keep executing CPU instructions
		select {
		case sample, ok := <-a.sampleBuffer:
			if !ok {
				return 0, io.EOF
			}

			binary.LittleEndian.PutUint32(buf[bufferPosition:], math.Float32bits(sample))
			bufferPosition += 4

		default:
			a.gb.Joypad.DetectKeysPressed()
			a.gb.CPU.ExecuteInstruction()
		}
	}
	return bufferPosition, nil
}
