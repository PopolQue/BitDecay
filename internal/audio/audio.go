package audio

import (
	"bytes"
	"io"
	"os"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
)

var (
	audioCtx    *audio.Context
	clickPlayer *audio.Player
	alarmPlayer *audio.Player
)

func Init() {
	audioCtx = audio.NewContext(44100)
}

func LoadSounds() {
	// Placeholder for loading sounds
	// clickPlayer = loadOGG("assets/sfx/click.ogg")
	// alarmPlayer = loadOGG("assets/sfx/alarm.ogg")
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

func loadOGG(path string) *audio.Player {
	f, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer f.Close()

	d, err := vorbis.DecodeWithSampleRate(44100, f)
	if err != nil {
		return nil
	}

	b, err := io.ReadAll(d)
	if err != nil {
		return nil
	}

	return audio.NewPlayerFromBytes(audioCtx, b)
}

func loadOGGFromBytes(b []byte) *audio.Player {
	d, err := vorbis.DecodeWithSampleRate(44100, bytes.NewReader(b))
	if err != nil {
		return nil
	}

	data, err := io.ReadAll(d)
	if err != nil {
		return nil
	}

	return audio.NewPlayerFromBytes(audioCtx, data)
}
