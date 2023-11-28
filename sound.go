package main

import (
	"log"

	"github.com/veandco/go-sdl2/mix"
)

type Sound struct {
	wav *mix.Chunk
}

const (
	SOUND_FILEPATH = "assets/beep.wav"
)

func (s *Sound) InitSound() {
	err := mix.OpenAudio(44100, mix.DEFAULT_FORMAT, 2, 4096)
	if err != nil {
		log.Fatalln("Failed to open audio")
	}

	s.wav, err = mix.LoadWAV(SOUND_FILEPATH)

	if err != nil {
		log.Fatalln("Failed to load wav")
	}
}

func (s *Sound) Play() {
	if mix.Playing(1) == 0 {
		s.wav.Play(1, 10)
	}
}

func (s *Sound) Stop() {
	if mix.Playing(1) == 1 {
		mix.HaltChannel(1)
	}
}
