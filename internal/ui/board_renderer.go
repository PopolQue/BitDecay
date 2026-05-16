package ui

import (
	"fmt"
	"image/color"
	"math"

	"github.com/popolque/firstbitengi/internal/constants"
	"github.com/popolque/firstbitengi/internal/game"
	"github.com/google/uuid"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

func (r *Renderer) drawHexBoard(screen *ebiten.Image) {
	centerX, centerY := 350.0, 300.0
	size := float64(constants.HexSize)

	for key, owner := range r.Engine.State.Board {
		q, r_coord, _ := game.ParseTile(key)
		
		x := centerX + size*(math.Sqrt(3)*float64(q)+math.Sqrt(3)/2.0*float64(r_coord))
		y := centerY + size*(3.0/2.0*float64(r_coord))*0.75 

		r.drawIsometricHex(screen, float32(x), float32(y), key, owner)
		
		for _, asset := range r.Engine.State.Assets {
			if asset.Q == q && asset.R == r_coord {
				r.drawAssetIcon(screen, float32(x), float32(y), asset)
			}
		}
	}
}

func (r *Renderer) drawIsometricHex(screen *ebiten.Image, x, y float32, key string, owner int) {
	size := float32(constants.HexSize)
	
	path := vector.Path{}
	for i := 0; i < 6; i++ {
		angle := 2.0 * math.Pi / 6 * (float64(i) + 0.5)
		px := x + size*float32(math.Cos(angle))
		py := y + size*float32(math.Sin(angle))*0.75 
		if i == 0 {
			path.MoveTo(px, py)
		} else {
			path.LineTo(px, py)
		}
	}
	path.Close()

	col := color.RGBA{174, 213, 129, 255} 
	if owner == 0 {
		col = color.RGBA{124, 179, 66, 255} 
	} else if r.Engine.State.SelectedTile == key {
		col = color.RGBA{255, 167, 38, 255} 
	}

	vertices, indices := path.AppendVerticesAndIndicesForFilling(nil, nil)
	for i := range vertices {
		vertices[i].ColorR = float32(col.R) / 255
		vertices[i].ColorG = float32(col.G) / 255
		vertices[i].ColorB = float32(col.B) / 255
		vertices[i].ColorA = 1
	}
	screen.DrawTriangles(vertices, indices, emptyImage, &ebiten.DrawTrianglesOptions{})
	
	for i := 0; i < 6; i++ {
		angle1 := 2.0 * math.Pi / 6 * (float64(i) + 0.5)
		angle2 := 2.0 * math.Pi / 6 * (float64(i+1) + 0.5)
		px1 := x + size*float32(math.Cos(angle1))
		py1 := y + size*float32(math.Sin(angle1))*0.75
		px2 := x + size*float32(math.Cos(angle2))
		py2 := y + size*float32(math.Sin(angle2))*0.75
		vector.StrokeLine(screen, px1, py1, px2, py2, 1, color.RGBA{0, 0, 0, 80}, true)
	}
}

func hexRound(q, r float64) (int, int) {
	s := -q - r
	rq := math.Round(q)
	rr := math.Round(r)
	rs := math.Round(s)

	dq := math.Abs(rq - q)
	dr := math.Abs(rr - r)
	ds := math.Abs(rs - s)

	if dq > dr && dq > ds {
		rq = -rr - rs
	} else if dr > ds {
		rr = -rq - rs
	} else {
		rs = -rq - rr
	}
	
	return int(rq), int(rr)
}

func (r *Renderer) HandleMouseInput() {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		mx, my := ebiten.CursorPosition()
		if mx >= 620 {
			for i, npc := range r.Engine.State.GuestPool {
				y := 80 + i*70
				if my >= y && my <= y+60 {
					r.Engine.State.SelectedNPCID = npc.ID
					return
				}
			}
		} else {
			centerX, centerY := 350.0, 300.0
			size := float64(constants.HexSize)
			
			dx := float64(mx) - centerX
			dy := (float64(my) - centerY) / 0.75 
			
			q := (math.Sqrt(3)/3.0*dx - 1.0/3.0*dy) / size
			r_axial := (2.0 / 3.0 * dy) / size
			
			q_int, r_int := hexRound(q, r_axial)
			
			key := fmt.Sprintf("%d,%d", q_int, r_int)
			if _, exists := r.Engine.State.Board[key]; exists {
				r.Engine.State.SelectedTile = key
				r.Engine.State.SelectedAssetID = uuid.Nil
				for _, asset := range r.Engine.State.Assets {
					if asset.Q == q_int && asset.R == r_int {
						r.Engine.State.SelectedAssetID = asset.ID
						break
					}
				}
			}
		}
	}
}
