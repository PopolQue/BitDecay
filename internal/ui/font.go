package ui

import (
	"bytes"
	"fmt"

	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/popolque/firstbitengi/assets"
)

var (
	MainFaceSource *text.GoTextFaceSource
)

func LoadFont(path string) error {
	// Path is relative to the assets/ folder because we embedded sfx/ and fonts/
	// If LoadFont was called with "assets/fonts/BPdotsLight.otf", we should change it to "fonts/BPdotsLight.otf"
	f, err := assets.FS.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read font file: %w", err)
	}

	s, err := text.NewGoTextFaceSource(bytes.NewReader(f))
	if err != nil {
		return fmt.Errorf("failed to create font source: %w", err)
	}

	MainFaceSource = s
	return nil
}
