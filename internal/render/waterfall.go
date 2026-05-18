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
	X         int
	Y         int
	MaxHeight int
	Symbols   []string
	Life      int // Frames until it relocates
}

func NewMatrixColumn(x, y, maxHeight, life int) *MatrixColumn {
	return &MatrixColumn{
		X:         x,
		Y:         y,
		MaxHeight: maxHeight,
		Life:      life,
		Symbols:   make([]string, maxHeight),
	}
}

type WaterfallRenderer struct {
	columns   []*MatrixColumn
	tick      int
	offscreen *ebiten.Image
}

func NewWaterfallRenderer() *WaterfallRenderer {
	wr := &WaterfallRenderer{
		offscreen: ebiten.NewImage(ui.ScreenWidth, ui.ScreenHeight),
	}
	return wr
}

func (w *WaterfallRenderer) Update(corruption float64) {
	w.tick++

	// Maintain a fixed number of columns
	targetCount := 60
	if len(w.columns) < targetCount {
		w.spawnColumn(corruption)
	}

	charset := "01  "
	if corruption > 75 {
		charset = "01#@!?% "
	}

	// Update columns
	for i := len(w.columns) - 1; i >= 0; i-- {
		c := w.columns[i]
		c.Life--

		// Randomly change symbols
		if w.tick%3 == 0 {
			for j := 0; j < len(c.Symbols); j++ {
				if rand.Float64() > 0.2 {
					c.Symbols[j] = string(charset[rand.Intn(len(charset))])
				}
			}
		}

		// If column life is over, remove it
		if c.Life <= 0 {
			w.columns = append(w.columns[:i], w.columns[i+1:]...)
		}
	}
}

func (w *WaterfallRenderer) spawnColumn(corruption float64) {
	x := rand.Intn(ui.ScreenWidth)
	y := rand.Intn(ui.ScreenHeight)
	maxHeight := rand.Intn(15) + 5
	life := rand.Intn(100) + 50

	w.columns = append(w.columns, NewMatrixColumn(x, y, maxHeight, life))
}

func (w *WaterfallRenderer) Draw(screen *ebiten.Image) {
	w.offscreen.Clear()
	lineHeight := 14

	for _, c := range w.columns {
		for i, sym := range c.Symbols {
			currY := c.Y + i*lineHeight

			// Only draw if not inside the widget rect
			pt := image.Pt(c.X, currY)
			if !pt.In(ui.WidgetRect) && currY >= 0 && currY < ui.ScreenHeight {
				ebitenutil.DebugPrintAt(w.offscreen, sym, c.X, currY)
			}
		}
	}

	// Draw the white text with neon green tint
	op := &ebiten.DrawImageOptions{}
	op.ColorScale.ScaleWithColor(color.RGBA{57, 255, 20, 255}) // Neon Green
	screen.DrawImage(w.offscreen, op)
}
