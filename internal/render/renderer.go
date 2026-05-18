package render

import (
	_ "embed"
	"fmt"
	"image"
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/popolque/firstbitengi/internal/format"
	"github.com/popolque/firstbitengi/internal/model"
	"github.com/popolque/firstbitengi/internal/ui"
)

//go:embed crt.kage
var crtShaderSrc []byte

type Renderer struct {
	waterfall  *WaterfallRenderer
	glitch     *GlitchSystem
	crtShader  *ebiten.Shader
	offscreen  *ebiten.Image
	normalFace *text.GoTextFace
	smallFace  *text.GoTextFace
	largeFace  *text.GoTextFace
}

func NewRenderer() *Renderer {
	s, err := ebiten.NewShader(crtShaderSrc)
	if err != nil {
		log.Printf("failed to load shader: %v\n", err)
	}

	r := &Renderer{
		waterfall: NewWaterfallRenderer(),
		glitch:    NewGlitchSystem(),
		crtShader: s,
		offscreen: ebiten.NewImage(ui.ScreenWidth, ui.ScreenHeight),
	}

	if ui.MainFaceSource != nil {
		r.normalFace = &text.GoTextFace{Source: ui.MainFaceSource, Size: 18}
		r.smallFace = &text.GoTextFace{Source: ui.MainFaceSource, Size: 14}
		r.largeFace = &text.GoTextFace{Source: ui.MainFaceSource, Size: 32}
	}

	return r
}

func (r *Renderer) Draw(screen *ebiten.Image, state *model.GameState) {
	// Recreate offscreen if resolution changed
	if r.offscreen.Bounds().Dx() != ui.ScreenWidth || r.offscreen.Bounds().Dy() != ui.ScreenHeight {
		r.offscreen = ebiten.NewImage(ui.ScreenWidth, ui.ScreenHeight)
	}

	r.offscreen.Clear()
	r.drawToImage(r.offscreen, state)

	if r.crtShader == nil {
		screen.DrawImage(r.offscreen, nil)
		return
	}

	op := &ebiten.DrawRectShaderOptions{}
	op.Images[0] = r.offscreen
	screen.DrawRectShader(ui.ScreenWidth, ui.ScreenHeight, r.crtShader, op)
}

const asciiHeader = `
  ____  _____ _______        ____   _______ ____       __     __
 |  _ \|_   _|__   __|      |  __ \|  ____/ ____|   /\ \ \   / /
 | |_) | | |    | |         | |  | | |__ | |       /  \ \ \_/ / 
 |  _ <  | |    | |         | |  | |  __|| |      / /\ \ \   /  
 | |_) |_| |_   | |         | |__| | |___| |____ / ____ \ | |   
 |____/|_____|  |_|         |_____/|______\_____/_/    \_\|_|   
 
 `

func (r *Renderer) DrawText(screen *ebiten.Image, str string, x, y int, clr color.Color) {
	ebitenutil.DebugPrintAt(screen, str, x, y)
}

func (r *Renderer) drawToImage(screen *ebiten.Image, state *model.GameState) {
	r.waterfall.Update(state.Corruption)
	r.glitch.Step(state.Corruption)

	screen.Fill(color.Black)
	r.waterfall.Draw(screen)
	if !ui.IsPortrait() {
		r.drawHeader(screen)
	}
	r.drawWidget(screen, state)
	r.glitch.Draw(screen)
}

func (r *Renderer) drawHeader(screen *ebiten.Image) {
	lines := 8
	charWidth := 6
	headerWidth := 67 * charWidth
	wr := ui.GetWidgetRect()
	x := (ui.ScreenWidth - headerWidth) / 2
	y := wr.Min.Y - (lines * 14) - 20

	if y < 10 {
		return
	}

	ebitenutil.DrawRect(screen, float64(x-20), float64(y-10), float64(headerWidth+40), float64(lines*14+20), color.Black)

	tempImg := ebiten.NewImage(headerWidth+20, lines*14+20)
	ebitenutil.DebugPrint(tempImg, asciiHeader)
	
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(x), float64(y))
	op.ColorScale.ScaleWithColor(color.RGBA{57, 255, 20, 255})
	screen.DrawImage(tempImg, op)
}

func (r *Renderer) drawWidget(screen *ebiten.Image, state *model.GameState) {
	neonGreen := color.RGBA{57, 255, 20, 255}
	bgGreen := color.RGBA{0, 20, 0, 240}
	wr := ui.GetWidgetRect()

	ebitenutil.DrawRect(screen, float64(wr.Min.X), float64(wr.Min.Y),
		float64(wr.Dx()), float64(wr.Dy()), bgGreen)
	r.drawRectBorder(screen, wr, neonGreen)

	r.drawHUD(screen, state)

	// Rhythm Display
	cr := ui.GetClickerRect()
	ry := cr.Max.Y + 10
	rx := cr.Min.X
	
	loopLen := 32.0
	loopPos := math.Mod(state.AudioTime, loopLen) / loopLen
	r.DrawText(screen, "SEQ_PROG", rx, ry, color.RGBA{100, 255, 100, 255})
	r.drawProgressBar(screen, rx + 100, ry+2, 200, 6, loopPos, neonGreen)

	if state.Combo > 0 {
		r.DrawText(screen, fmt.Sprintf("COMBO: %d (%.1fX)", state.Combo, state.ComboMultiplier), rx, ry+18, neonGreen)
	} else {
		r.DrawText(screen, "SYNC_LOST: KEEP_THE_BEAT", rx, ry+18, color.RGBA{200, 0, 0, 255})
	}

	// 32nd Note Blinking Dot (SYNCED WITH ENGINE)
	interval := 0.0625
	beatPos := math.Mod(state.AudioTime, interval)
	tolerance := 0.025 // Match engine
	isHit := beatPos < tolerance || beatPos > (interval-tolerance)

	dotX := rx + 105
	dotY := ry + 36
	dotSize := 8.0
	
	ebitenutil.DrawRect(screen, float64(dotX)-2, float64(dotY)-2, dotSize+4, dotSize+4, color.RGBA{0, 40, 0, 255})
	if isHit {
		ebitenutil.DrawRect(screen, float64(dotX), float64(dotY), dotSize, dotSize, neonGreen)
		r.DrawText(screen, "HIT_WINDOW", dotX+20, dotY-6, neonGreen)
	} else {
		ebitenutil.DrawRect(screen, float64(dotX), float64(dotY), dotSize, dotSize, color.RGBA{0, 80, 0, 255})
	}

	r.drawHardwareList(screen, state)
	r.drawUpgradeList(screen, state)
	r.drawSystemLog(screen, state)
	r.drawClicker(screen, state, neonGreen)
	r.drawReboot(screen, state, neonGreen)

	if state.PacketActive {
		r.drawPacketIntercept(screen, state, neonGreen)
	}

	if state.RebootPending {
		r.drawRebootDialog(screen, state, neonGreen)
	}
}

func (r *Renderer) drawHUD(screen *ebiten.Image, state *model.GameState) {
	mr := ui.GetMetricsRect()
	r.drawRectBorder(screen, mr, color.RGBA{0, 200, 0, 255})
	
	r.DrawText(screen, fmt.Sprintf("BITS: %s", format.FormatBits(state.Bits)), mr.Min.X+10, mr.Min.Y+15, color.White)
	r.DrawText(screen, fmt.Sprintf("RANK: %s", state.GetRank()), mr.Min.X+10, mr.Min.Y+40, color.White)
	r.DrawText(screen, fmt.Sprintf("GHz: %.3fX", state.GHzMultiplier), mr.Min.X+10, mr.Min.Y+65, color.White)

	startX := mr.Min.X + 250
	if ui.IsPortrait() {
		startX = mr.Min.X + 10
		startY := mr.Min.Y + 90
		r.DrawText(screen, fmt.Sprintf("PWR: %.0fW", state.PowerUsage), startX, startY, color.White)
		r.drawProgressBar(screen, startX+80, startY+2, 120, 8, state.PowerUsage/state.PowerCapacity, color.RGBA{200, 200, 0, 255})
		
		r.DrawText(screen, fmt.Sprintf("THM: %.1fC", state.HeatLevel), startX, startY+20, color.White)
		r.drawProgressBar(screen, startX+80, startY+22, 120, 8, state.HeatLevel/100.0, color.RGBA{0, 255, 0, 255})
		
		r.DrawText(screen, fmt.Sprintf("ENT: %.1f%%", state.Entropy), startX, startY+40, color.White)
		r.drawProgressBar(screen, startX+80, startY+42, 120, 8, state.Entropy/100.0, color.RGBA{255, 165, 0, 255})
	} else {
		r.DrawText(screen, fmt.Sprintf("TOTAL: %s", format.FormatBits(state.TotalBitsEarned)), mr.Min.X+10, mr.Min.Y+90, color.White)

		r.DrawText(screen, fmt.Sprintf("POWER: %.0fW / %.0fW", state.PowerUsage, state.PowerCapacity), startX, mr.Min.Y+15, color.White)
		r.drawProgressBar(screen, startX, mr.Min.Y+30, 250, 8, state.PowerUsage/state.PowerCapacity, color.RGBA{200, 200, 0, 255})

		r.DrawText(screen, fmt.Sprintf("THERMAL: %.1fC", state.HeatLevel), startX, mr.Min.Y+50, color.White)
		r.drawProgressBar(screen, startX, mr.Min.Y+65, 250, 8, state.HeatLevel/100.0, color.RGBA{0, 255, 0, 255})

		startX2 := mr.Min.X + 650
		r.DrawText(screen, fmt.Sprintf("ENTROPY: %.1f%%", state.Entropy), startX2, mr.Min.Y+15, color.White)
		r.drawProgressBar(screen, startX2, mr.Min.Y+30, 250, 8, state.Entropy/100.0, color.RGBA{255, 165, 0, 255})

		r.DrawText(screen, fmt.Sprintf("CORRUPTION: %.1f%%", state.Corruption), startX2, mr.Min.Y+50, color.White)
		r.drawProgressBar(screen, startX2, mr.Min.Y+65, 250, 8, state.Corruption/100.0, color.RGBA{255, 0, 255, 255})
	}
}

func (r *Renderer) drawProgressBar(screen *ebiten.Image, x, y, w, h int, ratio float64, clr color.Color) {
	if ratio > 1.0 { ratio = 1.0 }
	if ratio < 0 { ratio = 0 }
	ebitenutil.DrawRect(screen, float64(x), float64(y), float64(w), float64(h), color.RGBA{0, 50, 0, 255})
	ebitenutil.DrawRect(screen, float64(x), float64(y), float64(w)*ratio, float64(h), clr)
	r.drawRectBorder(screen, image.Rect(x, y, x+w, y+h), color.RGBA{0, 100, 0, 255})
}

func (r *Renderer) drawClicker(screen *ebiten.Image, state *model.GameState, neonGreen color.Color) {
	cr := ui.GetClickerRect()
	clr := neonGreen
	if state.ClickerFlash { clr = color.White }
	r.drawRectBorder(screen, cr, clr)
	r.DrawText(screen, "MANUAL OVERRIDE", cr.Min.X+10, cr.Min.Y+10, color.White)
}

func (r *Renderer) drawHardwareList(screen *ebiten.Image, state *model.GameState) {
	hr := ui.GetHardwareRect()
	r.drawRectBorder(screen, hr, color.RGBA{0, 200, 0, 255})
	r.DrawText(screen, "[ HARDWARE ]", hr.Min.X, hr.Min.Y-15, color.White)
	
	x := hr.Min.X + 5
	startY := hr.Min.Y + 5
	rowHeight := 40
	if ui.IsPortrait() { rowHeight = 30 }

	for i, def := range model.AllHardware {
		y := startY + i*rowHeight - state.ScrollOffset
		if y < hr.Min.Y || y > hr.Max.Y-20 {
			continue
		}

		owned := state.Hardware[def.ID]
		cost := model.CurrentCost(def, owned)
		label := fmt.Sprintf("%s (x%d) - %s", def.Name, owned, format.FormatBits(cost))
		var clr color.Color = color.White
		if state.Bits < cost { clr = color.RGBA{0, 100, 0, 255} }
		r.DrawText(screen, label, x, y, clr)
	}
}

func (r *Renderer) drawUpgradeList(screen *ebiten.Image, state *model.GameState) {
	ur := ui.GetUpgradeRect()
	r.drawRectBorder(screen, ur, color.RGBA{0, 200, 0, 255})
	r.DrawText(screen, "[ UPGRADES ]", ur.Min.X, ur.Min.Y-15, color.White)
	
	x := ur.Min.X + 5
	startY := ur.Min.Y + 5
	rowHeight := 40
	if ui.IsPortrait() { rowHeight = 30 }

	for i, def := range model.AllUpgrades {
		y := startY + i*rowHeight - state.ScrollOffset
		if y < ur.Min.Y || y > ur.Max.Y-20 {
			continue
		}

		owned := state.Upgrades[def.ID]
		var clr color.Color = color.White
		if owned {
			r.DrawText(screen, fmt.Sprintf("[X] %s", def.Name), x, y, color.RGBA{0, 150, 0, 255})
			continue
		}

		if state.Bits < def.Cost { clr = color.RGBA{0, 100, 0, 255} }
		r.DrawText(screen, fmt.Sprintf("%s - %s", def.Name, format.FormatBits(def.Cost)), x, y, clr)
	}
}

func (r *Renderer) drawSystemLog(screen *ebiten.Image, state *model.GameState) {
	lr := ui.GetLogRect()
	r.drawRectBorder(screen, lr, color.RGBA{0, 200, 0, 255})
	r.DrawText(screen, "[ LOG ]", lr.Min.X, lr.Min.Y-15, color.White)

	logX := lr.Min.X + 5
	logY := lr.Min.Y + 5
	lineHeight := 14
	for i, msg := range state.MessageLog {
		r.DrawText(screen, msg, logX, logY+i*lineHeight, color.White)
	}
}

func (r *Renderer) drawReboot(screen *ebiten.Image, state *model.GameState, neonGreen color.Color) {
	rr := ui.GetRebootRect()
	threshold := state.GetRebootThreshold()
	clr := neonGreen
	if state.TotalBitsEarned < threshold {
		clr = color.RGBA{0, 80, 0, 255}
	}
	r.drawRectBorder(screen, rr, clr)
	label := "REBOOT"
	if !ui.IsPortrait() { label = "SYSTEM_REBOOT" }
	r.DrawText(screen, label, rr.Min.X+5, rr.Min.Y+20, color.White)
}

func (r *Renderer) drawPacketIntercept(screen *ebiten.Image, state *model.GameState, neonGreen color.Color) {
	rect := ui.GetPacketRect()
	clr := neonGreen
	if (int(state.PacketTimer*10) % 4) < 2 { clr = color.White }
	ebitenutil.DrawRect(screen, float64(rect.Min.X), float64(rect.Min.Y), float64(rect.Dx()), float64(rect.Dy()), color.RGBA{0, 100, 0, 100})
	r.drawRectBorder(screen, rect, clr)
	r.DrawText(screen, ">> PACKET <<", rect.Min.X+5, rect.Min.Y+20, color.White)
}

func (r *Renderer) drawRebootDialog(screen *ebiten.Image, state *model.GameState, neonGreen color.Color) {
	wr := ui.GetWidgetRect()
	rect := image.Rect(wr.Min.X+20, wr.Min.Y+20, wr.Max.X-20, wr.Max.Y-20)
	ebitenutil.DrawRect(screen, float64(rect.Min.X), float64(rect.Min.Y), float64(rect.Dx()), float64(rect.Dy()), color.Black)
	r.drawRectBorder(screen, rect, neonGreen)

	r.DrawText(screen, ">> REBOOT? <<", rect.Min.X+20, rect.Min.Y+50, color.White)
	r.DrawText(screen, "[Y] CONFIRM / [N] ABORT", rect.Min.X+20, rect.Min.Y+150, color.White)
}

func (r *Renderer) drawRectBorder(screen *ebiten.Image, rect image.Rectangle, clr color.Color) {
	ebitenutil.DrawLine(screen, float64(rect.Min.X), float64(rect.Min.Y), float64(rect.Max.X), float64(rect.Min.Y), clr)
	ebitenutil.DrawLine(screen, float64(rect.Max.X), float64(rect.Min.Y), float64(rect.Max.X), float64(rect.Max.Y), clr)
	ebitenutil.DrawLine(screen, float64(rect.Max.X), float64(rect.Max.Y), float64(rect.Min.X), float64(rect.Max.Y), clr)
	ebitenutil.DrawLine(screen, float64(rect.Min.X), float64(rect.Max.Y), float64(rect.Min.X), float64(rect.Min.Y), clr)
}
