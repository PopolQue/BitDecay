package audio

import (
	"bytes"
	"io"
	"log"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
	"github.com/popolque/firstbitengi/assets"
)

var (
	audioCtx    *audio.Context
	clickPlayer *audio.Player
	loopPlayer  *audio.Player
	alarmPlayer *audio.Player
)

func Init() {
	if audioCtx == nil {
		audioCtx = audio.NewContext(44100)
	}
}

func LoadSounds() {
	var err error
	clickPlayer, err = loadWAV("sfx/BitDecayClick.wav", false)
	if err != nil {
		log.Printf("failed to load click sound: %v\n", err)
	}

	loopPlayer, err = loadWAV("sfx/BitDecayLoop.wav", true)
	if err != nil {
		log.Printf("failed to load loop sound: %v\n", err)
	}
	
	// Use click sound for alarm if no specific alarm sound exists
	alarmPlayer, _ = loadWAV("sfx/BitDecayClick.wav", false)
}

func StartLoop() {
	if loopPlayer != nil && !loopPlayer.IsPlaying() {
		loopPlayer.Play()
	}
}

func Update() {
	// If we use NewInfiniteLoop, it will handle looping automatically once Play() is called.
	// We don't need manual check here if we use NewInfiniteLoop correctly.
}

func PlayClick() {
	if clickPlayer != nil {
		clickPlayer.Rewind()
		clickPlayer.Play()
	}
}

func PlayAlarm() {
	if alarmPlayer != nil {
		alarmPlayer.Rewind()
		alarmPlayer.Play()
	}
}

func loadWAV(path string, loop bool) (*audio.Player, error) {
	b, err := assets.FS.ReadFile(path)
	if err != nil {
		return nil, err
	}

	d, err := wav.DecodeWithSampleRate(44100, bytes.NewReader(b))
	if err != nil {
		return nil, err
	}

	if loop {
		infiniteLoop := audio.NewInfiniteLoop(d, d.Length())
		return audioCtx.NewPlayer(infiniteLoop)
	}

	decoded, err := io.ReadAll(d)
	if err != nil {
		return nil, err
	}

	return audio.NewPlayerFromBytes(audioCtx, decoded), nil
}
