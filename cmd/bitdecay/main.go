package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/popolque/firstbitengi/internal/audio"
	"github.com/popolque/firstbitengi/internal/engine"
)

func main() {
	audio.Init()
	audio.LoadSounds()

	game := engine.NewGame()

	ebiten.SetWindowTitle("BIT-DECAY // MNEMONIC_OVERRIDE_3.0")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
