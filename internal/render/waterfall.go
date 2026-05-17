package render

import (
	"image"
	"image/color"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/popolque/firstbitengi/internal/ui"
)

type MatrixColumn struct {
	X       int
	Y       float64
	Height  int
	Speed   float64
	Symbols []string
}

type WaterfallRenderer struct {
	columns   []*MatrixColumn
	tick      int
	offscreen *ebiten.Image
}

func NewWaterfallRenderer() *WaterfallRenderer {
	return &WaterfallRenderer{
		offscreen: ebiten.NewImage(ui.ScreenWidth, ui.ScreenHeight),
	}
}

func (w *WaterfallRenderer) Update(corruption float64) {
	w.tick++

	// Spawn new columns
	if len(w.columns) < 60 && w.tick%5 == 0 { // Increased density and spawn rate
		w.spawnColumn(corruption)
	}

	// Update columns
	for i := len(w.columns) - 1; i >= 0; i-- {
		c := w.columns[i]
		c.Y += c.Speed

		// If column is completely off screen (below), remove it
		if int(c.Y) > ui.ScreenHeight {
			w.columns = append(w.columns[:i], w.columns[i+1:]...)
		}
	}
}

func (w *WaterfallRenderer) spawnColumn(corruption float64) {
	x := rand.Intn(ui.ScreenWidth)
	speed := rand.Float64()*3 + 1 // Slightly faster
	height := rand.Intn(15) + 5   // Slightly taller
	
	charset := "01  "
	if corruption > 75 {
		charset = "01#@!?% "
	}

	symbols := make([]string, height)
	for i := 0; i < height; i++ {
		symbols[i] = string(charset[rand.Intn(len(charset))])
	}

	w.columns = append(w.columns, &MatrixColumn{
		X:       x,
		Y:       float64(-height * 14),
		Height:  height,
		Speed:   speed,
		Symbols: symbols,
	})
}

func (w *WaterfallRenderer) Draw(screen *ebiten.Image) {
	w.offscreen.Clear()
	lineHeight := 14

	for _, c := range w.columns {
		for i, sym := range c.Symbols {
			y := int(c.Y) + i*lineHeight
			
			// Only draw if not inside the widget rect
			pt := image.Pt(c.X, y)
			if !pt.In(ui.WidgetRect) && y >= 0 && y < ui.ScreenHeight {
				ebitenutil.DebugPrintAt(w.offscreen, sym, c.X, y)
			}
		}
	}

	// Draw the white text with neon green tint
	op := &ebiten.DrawImageOptions{}
	op.ColorScale.ScaleWithColor(color.RGBA{57, 255, 20, 255}) // Neon Green
	screen.DrawImage(w.offscreen, op)
}
