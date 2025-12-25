package main

import (
	"log/slog"
	"sync"
	"time"

	"github.com/hajimehoshi/ebiten/v2/audio"
)

const FrameSize = 4 // 1 frame = 4 bytes (stereo 16bit)

type AudioManager struct {
	Started bool
	Stream  *AudioStream
	Player  *audio.Player
}

type AudioStream struct {
	// Ring buffer
	Buffer            [16384]uint8
	ReadPos, WritePos int
	CurrentSize       int

	m sync.Mutex
}

func NewAudioManager() *AudioManager {
	a := &AudioManager{
		Stream: &AudioStream{},
	}

	context := audio.NewContext(32768)

	m, err := context.NewPlayer(a)
	if err != nil {
		panic(err)
	}
	slog.Info("audio player created successfully", slog.Int("sampleRate", context.SampleRate()))
	a.Player = m
	a.Player.SetVolume(App.Config.Audio.Volume)
	a.Player.SetBufferSize(time.Second / 20) // wasmのときは time.Second / 16 の方が良いかも

	return a
}

func (a *AudioManager) Update() {
	if !a.Started {
		a.Player.Play()
		a.Started = true
	}
}

func (a *AudioManager) Close() error {
	err := a.Player.Close()
	if err != nil {
		return err
	}
	slog.Info("audio player closed successfully")
	return nil
}

func (a *AudioManager) Write(p []uint8) (int, error) {
	if a.Started {
		a.Stream.m.Lock()
		defer a.Stream.m.Unlock()
		n := len(p)

		// check available space
		available := len(a.Stream.Buffer) - a.Stream.CurrentSize
		if n > available {
			n = available
		}

		// write data to ring buffer
		for i := 0; i < n; i++ {
			a.Stream.Buffer[a.Stream.WritePos] = p[i]
			a.Stream.WritePos = (a.Stream.WritePos + 1) % len(a.Stream.Buffer)
		}
		a.Stream.CurrentSize += n

		return n, nil
	}

	return len(p), nil
}

func (a *AudioManager) Read(p []uint8) (int, error) {
	a.Stream.m.Lock()
	defer a.Stream.m.Unlock()

	n := len(p)
	if n > a.Stream.CurrentSize {
		n = a.Stream.CurrentSize
	}

	for i := 0; i < n; i++ { // read data from ring buffer
		p[i] = a.Stream.Buffer[a.Stream.ReadPos]
		a.Stream.ReadPos = (a.Stream.ReadPos + 1) % len(a.Stream.Buffer)
	}
	a.Stream.CurrentSize -= n

	for i := n; i < len(p); i++ {
		p[i] = 0 // zero-fill if not enough data
	}

	return len(p), nil
}
