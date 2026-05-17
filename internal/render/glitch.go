package render

import (
	"image/color"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type GlitchSystem struct {
	TearLines  []TearLine
	NoiseAlpha uint8
	Tick       int
}

type TearLine struct {
	Y       int
	Width   int
	OffsetX int
	Life    int // frames remaining
}

func NewGlitchSystem() *GlitchSystem {
	return &GlitchSystem{}
}

func (g *GlitchSystem) Step(corruption float64) {
	g.Tick++
	if corruption < 50 {
		g.NoiseAlpha = 0
		g.TearLines = nil
		return
	}

	g.NoiseAlpha = uint8(((corruption - 50) / 50.0) * 150)

	// Spawn new tears above 75% corruption
	if corruption > 75 && g.Tick%8 == 0 {
		g.TearLines = append(g.TearLines, TearLine{
			Y:       rand.Intn(768),
			Width:   rand.Intn(400) + 100,
			OffsetX: rand.Intn(20) - 10,
			Life:    rand.Intn(6) + 2,
		})
	}

	// Age out old tears
	alive := g.TearLines[:0]
	for _, t := range g.TearLines {
		t.Life--
		if t.Life > 0 {
			alive = append(alive, t)
		}
	}
	g.TearLines = alive
}

func (g *GlitchSystem) Draw(screen *ebiten.Image) {
	if g.NoiseAlpha == 0 && len(g.TearLines) == 0 {
		return
	}

	// Horizontal noise bands
	for _, tear := range g.TearLines {
		ebitenutil.DrawRect(screen, float64(tear.OffsetX), float64(tear.Y),
			float64(tear.Width), 2,
			color.RGBA{0, 255, 70, g.NoiseAlpha})
	}

	// Full-screen noise pixel scatter
	if g.NoiseAlpha > 0 {
		for i := 0; i < int(g.NoiseAlpha/2); i++ {
			x, y := rand.Intn(1280), rand.Intn(768)
			screen.Set(x, y, color.RGBA{0, 255, 0, g.NoiseAlpha / 2})
		}
	}
}
