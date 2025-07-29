package ui

import (
	"encoding/binary"
	"github.com/ebitengine/oto/v3"
	"io"
	"math"
	"path/filepath"
)

const (
	sampleRate = 44100
	channels   = 2
	format     = oto.FormatFloat32LE

	bufferSize = 8192
)

var (
	audioFileBuffer         []byte
	audioFileBufferPosition int
)

func newAudioPlayer(r io.Reader) (*oto.Player, error) {
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

			// Write also to file buffer
			if ui.audioFile != nil {
				binary.LittleEndian.PutUint32(audioFileBuffer[audioFileBufferPosition:], math.Float32bits(sample))
				audioFileBufferPosition += 4
			}

		default:
			// If debugger is active and paused, return silence
			if ui.Paused || (ui.debugger.Active && !ui.debugger.Running) { // If paused, return silence
				binary.LittleEndian.PutUint32(buf[bufferPosition:], math.Float32bits(0))
				bufferPosition += 4
				continue
			}

			ui.gameBoy.Joypad.DetectKeysPressed()
			ui.gameBoy.CPU.ExecuteInstruction()

			if ui.debugger.Active {
				pc := ui.gameBoy.CPU.ReadPC()
				switch {
				// Stop if next instruction
				case ui.debugger.NextInstruction && ui.debugger.CallDepth <= 0:
					ui.debugger.NextInstruction = false
					ui.debugger.Stop()

				// Stop if breakpoint
				case ui.debugger.CheckBreakpoint(pc):
					ui.debugger.Stop()
				}
			}
		}
	}
	if ui.audioFile != nil {
		_, _ = ui.audioFile.Write(audioFileBuffer[:audioFileBufferPosition])
		audioFileBufferPosition = 0
	}

	return bufferPosition, nil
}

func (ui *UI) RecordAudio() (string, error) {
	// Extract the relative path
	fileName := filepath.Base(ui.fileName)
	f, err := createAudioFile(fileName)
	if err != nil {
		ui.audioPlayer = nil
		return "", err
	}

	// Create buffer
	audioFileBuffer = make([]byte, bufferSize)

	ui.audioFile = f
	return f.Name(), nil
}
