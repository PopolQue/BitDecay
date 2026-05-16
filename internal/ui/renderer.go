package ui

import (
	"image/color"

	"github.com/popolque/firstbitengi/internal/game"
	"github.com/popolque/firstbitengi/internal/models"
	"github.com/hajimehoshi/ebiten/v2"
)

type Renderer struct {
	Engine *game.Engine
}

func NewRenderer(e *game.Engine) *Renderer {
	return &Renderer{Engine: e}
}

var emptyImage = ebiten.NewImage(3, 3)

func (r *Renderer) Draw(screen *ebiten.Image) {
	if r.Engine.State.Screen == models.ScreenMenu {
		r.drawMenu(screen)
		return
	}

	screen.Fill(color.RGBA{109, 76, 65, 255}) // Wood board background

	r.drawTopBar(screen)
	r.drawActiveEvent(screen)
	r.drawHexBoard(screen)
	r.drawNPCPool(screen)
	r.drawInstructions(screen)
}
