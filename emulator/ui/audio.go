package ui

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"os"

	"github.com/ebitengine/oto/v3"
)

const (
	sampleRate = 44100
	channels   = 2
	format     = oto.FormatFloat32LE

	bufferSize = 8192
)

var AudioFile *os.File

func newAudioPlayer(r io.Reader) (*oto.Player, error) {
	file, err := os.OpenFile("audio.dat", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println("Error creating audio file:", err)
	} else {
		AudioFile = file
	}

	op := &oto.NewContextOptions{}
	op.SampleRate = sampleRate
	op.ChannelCount = channels
	op.Format = format

	ctx, ready, err := oto.NewContext(op)
	if err != nil {
		return nil, err
	}

	<-ready // Context ready

	p := ctx.NewPlayer(r)
	p.SetBufferSize(bufferSize)

	return p, nil
}

// Implements io.Reader interface for audio playback
func (ui *UI) Read(buf []byte) (n int, err error) {
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
			// If debugger is active and paused, return silence
			if ui.debugger.Active && !ui.debugger.Continue { // If paused, return silence
				binary.LittleEndian.PutUint32(buf[bufferPosition:], math.Float32bits(0))
				bufferPosition += 4
				continue
			}

			ui.gameBoy.Joypad.DetectKeysPressed()
			ui.gameBoy.CPU.ExecuteInstruction()

			if ui.debugger.Active { // && ui.debugger.Continue
				// Check breakpoint
				pc := ui.gameBoy.CPU.ReadPC()
				if ui.debugger.CheckBreakpoint(pc) {
					// Stop
					ui.debugger.Continue = false
				}
			}
		}
	}
	_, _ = AudioFile.Write(buf[:bufferPosition])

	return bufferPosition, nil
}
