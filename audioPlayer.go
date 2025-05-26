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

func (ui *UI) initAudioPlayer() error {
	op := &oto.NewContextOptions{}
	op.SampleRate = sampleRate
	op.ChannelCount = channels
	op.Format = format

	ctx, ready, err := oto.NewContext(op)
	if err != nil {
		return err
	}

	<-ready // Context ready

	p := ctx.NewPlayer(ui)
	p.SetBufferSize(bufferSize)

	// Store reference to player
	ui.audioPlayer = p

	return nil
}

func (ui *UI) Read(buf []byte) (n int, err error) {
	// If Game Boy is paused return silence and don't execute cpu instructions
	if ui.paused || ui.debugging {
		if ui.debugging && ui.stepInstruction {
			ui.stepInstruction = false
			ui.gameBoy.CPU.ExecuteInstruction()
		}
		for i := range len(buf) {
			buf[i] = 0
		}
		return len(buf), nil
	}

	bufferPosition := 0

	for bufferPosition < len(buf) {
		// If not enough samples have been produced, keep executing CPU instructions
		select {
		case sample, ok := <-ui.audioBuffer:
			if !ok {
				return 0, io.EOF
			}

			binary.LittleEndian.PutUint32(buf[bufferPosition:], math.Float32bits(sample))
			bufferPosition += 4

		default:
			ui.gameBoy.Joypad.DetectKeysPressed()
			ui.gameBoy.CPU.ExecuteInstruction()
		}
	}

	return bufferPosition, nil
}
