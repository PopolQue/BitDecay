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
	r.drawHeader(screen)
	r.drawWidget(screen, state)
	r.glitch.Draw(screen)
}

func (r *Renderer) drawHeader(screen *ebiten.Image) {
	lines := 8
	charWidth := 6
	headerWidth := 67 * charWidth
	x := (ui.ScreenWidth - headerWidth) / 2
	y := ui.WidgetY - (lines * 14) - 20

	if y < 10 {
		y = 10
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

	ebitenutil.DrawRect(screen, float64(ui.WidgetRect.Min.X), float64(ui.WidgetRect.Min.Y),
		float64(ui.WidgetRect.Dx()), float64(ui.WidgetRect.Dy()), bgGreen)
	r.drawRectBorder(screen, ui.WidgetRect, neonGreen)

	r.drawHUD(screen, state)
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
	r.drawRectBorder(screen, ui.MetricsHUDRect, color.RGBA{0, 200, 0, 255})
	
	r.DrawText(screen, fmt.Sprintf("BITS: %s", format.FormatBits(state.Bits)), ui.MetricsHUDRect.Min.X+20, ui.MetricsHUDRect.Min.Y+15, color.White)
	r.DrawText(screen, fmt.Sprintf("TOTAL: %s", format.FormatBits(state.TotalBitsEarned)), ui.MetricsHUDRect.Min.X+20, ui.MetricsHUDRect.Min.Y+40, color.White)
	r.DrawText(screen, fmt.Sprintf("RANK: %s", state.GetRank()), ui.MetricsHUDRect.Min.X+20, ui.MetricsHUDRect.Min.Y+65, color.White)
	r.DrawText(screen, fmt.Sprintf("GHz MULT: %.3fX", state.GHzMultiplier), ui.MetricsHUDRect.Min.X+20, ui.MetricsHUDRect.Min.Y+90, color.White)

	startX := ui.MetricsHUDRect.Min.X + 350
	r.DrawText(screen, fmt.Sprintf("POWER: %.0fW / %.0fW", state.PowerUsage, state.PowerCapacity), startX, ui.MetricsHUDRect.Min.Y+15, color.White)
	r.drawProgressBar(screen, startX, ui.MetricsHUDRect.Min.Y+30, 250, 8, state.PowerUsage/state.PowerCapacity, color.RGBA{200, 200, 0, 255})

	r.DrawText(screen, fmt.Sprintf("THERMAL: %.1fC", state.HeatLevel), startX, ui.MetricsHUDRect.Min.Y+50, color.White)
	thermalClr := color.RGBA{0, 255, 0, 255}
	if state.HeatLevel > 80 { thermalClr = color.RGBA{255, 0, 0, 255} }
	r.drawProgressBar(screen, startX, ui.MetricsHUDRect.Min.Y+65, 250, 8, state.HeatLevel/100.0, thermalClr)

	r.DrawText(screen, fmt.Sprintf("RACK SPACE: %.0fU / %.0fU", state.SpaceUsage, state.SpaceCapacity), startX, ui.MetricsHUDRect.Min.Y+85, color.White)
	r.drawProgressBar(screen, startX, ui.MetricsHUDRect.Min.Y+100, 250, 8, state.SpaceUsage/state.SpaceCapacity, color.RGBA{0, 200, 255, 255})

	startX2 := ui.MetricsHUDRect.Min.X + 650
	r.DrawText(screen, fmt.Sprintf("ENTROPY: %.1f%%", state.Entropy), startX2, ui.MetricsHUDRect.Min.Y+15, color.White)
	r.drawProgressBar(screen, startX2, ui.MetricsHUDRect.Min.Y+30, 250, 8, state.Entropy/100.0, color.RGBA{255, 165, 0, 255})

	r.DrawText(screen, fmt.Sprintf("CORRUPTION: %.1f%%", state.Corruption), startX2, ui.MetricsHUDRect.Min.Y+50, color.White)
	r.drawProgressBar(screen, startX2, ui.MetricsHUDRect.Min.Y+65, 250, 8, state.Corruption/100.0, color.RGBA{255, 0, 255, 255})
}

func (r *Renderer) drawProgressBar(screen *ebiten.Image, x, y, w, h int, ratio float64, clr color.Color) {
	if ratio > 1.0 { ratio = 1.0 }
	if ratio < 0 { ratio = 0 }
	ebitenutil.DrawRect(screen, float64(x), float64(y), float64(w), float64(h), color.RGBA{0, 50, 0, 255})
	ebitenutil.DrawRect(screen, float64(x), float64(y), float64(w)*ratio, float64(h), clr)
	r.drawRectBorder(screen, image.Rect(x, y, x+w, y+h), color.RGBA{0, 100, 0, 255})
}

func (r *Renderer) drawClicker(screen *ebiten.Image, state *model.GameState, neonGreen color.Color) {
	clr := neonGreen
	if state.ClickerFlash { clr = color.White }
	r.drawRectBorder(screen, ui.ClickerRegion, clr)
	r.DrawText(screen, "MANUAL OVERRIDE", ui.ClickerRegion.Min.X+80, ui.ClickerRegion.Min.Y+60, color.White)
}

func (r *Renderer) drawHardwareList(screen *ebiten.Image, state *model.GameState) {
	r.drawRectBorder(screen, ui.HardwareListRect, color.RGBA{0, 200, 0, 255})
	r.DrawText(screen, "[ HARDWARE ]", ui.HardwareListRect.Min.X+10, ui.HardwareListRect.Min.Y-20, color.White)
	
	x := ui.HardwareListRect.Min.X + 10
	startY := ui.HardwareListRect.Min.Y + 10
	rowHeight := 60

	for i, def := range model.AllHardware {
		y := startY + i*rowHeight - state.ScrollOffset
		if y < ui.HardwareListRect.Min.Y+5 || y > ui.HardwareListRect.Max.Y-rowHeight {
			continue
		}

		owned := state.Hardware[def.ID]
		cost := model.CurrentCost(def, owned)
		label := fmt.Sprintf("%s (x%d) - %s", def.Name, owned, format.FormatBits(cost))
		var clr color.Color = color.White
		if state.Bits < cost { clr = color.RGBA{0, 100, 0, 255} }
		r.DrawText(screen, label, x, y, clr)
		
		infra := ""
		if def.WattsImpact != 0 { infra += fmt.Sprintf("%.0fW ", def.WattsImpact) }
		if def.ThermalImpact != 0 { infra += fmt.Sprintf("%.0fC ", def.ThermalImpact) }
		if def.SpaceImpact != 0 { infra += fmt.Sprintf("%.0fU ", def.SpaceImpact) }
		r.DrawText(screen, fmt.Sprintf("BPS: %.1f | %s", def.BaseBPS, infra), x, y+20, clr)
	}
}

func (r *Renderer) drawUpgradeList(screen *ebiten.Image, state *model.GameState) {
	r.drawRectBorder(screen, ui.UpgradeListRect, color.RGBA{0, 200, 0, 255})
	r.DrawText(screen, "[ UPGRADES ]", ui.UpgradeListRect.Min.X+10, ui.UpgradeListRect.Min.Y-20, color.White)
	
	x := ui.UpgradeListRect.Min.X + 10
	startY := ui.UpgradeListRect.Min.Y + 10
	rowHeight := 60

	for i, def := range model.AllUpgrades {
		y := startY + i*rowHeight - state.ScrollOffset
		if y < ui.UpgradeListRect.Min.Y+5 || y > ui.UpgradeListRect.Max.Y-rowHeight {
			continue
		}

		owned := state.Upgrades[def.ID]
		var clr color.Color = color.White
		if owned {
			r.DrawText(screen, fmt.Sprintf("[X] %s", def.Name), x, y, color.RGBA{0, 150, 0, 255})
			r.DrawText(screen, "PURCHASED", x, y+20, color.RGBA{0, 150, 0, 255})
			continue
		}

		if state.Bits < def.Cost { clr = color.RGBA{0, 100, 0, 255} }
		r.DrawText(screen, fmt.Sprintf("%s - %s", def.Name, format.FormatBits(def.Cost)), x, y, clr)
		r.DrawText(screen, def.Description, x, y+20, clr)
	}
}

func (r *Renderer) drawSystemLog(screen *ebiten.Image, state *model.GameState) {
	r.drawRectBorder(screen, ui.LogRect, color.RGBA{0, 200, 0, 255})
	r.DrawText(screen, "[ LOG ]", ui.LogRect.Min.X+10, ui.LogRect.Min.Y-20, color.White)

	logX := ui.LogRect.Min.X + 10
	logY := ui.LogRect.Min.Y + 10
	lineHeight := 18
	for i, msg := range state.MessageLog {
		r.DrawText(screen, msg, logX, logY+i*lineHeight, color.White)
	}
}
func (r *Renderer) drawReboot(screen *ebiten.Image, state *model.GameState, neonGreen color.Color) {
	threshold := state.GetRebootThreshold()
	clr := neonGreen
	if state.TotalBitsEarned < threshold {
		clr = color.RGBA{0, 80, 0, 255} // Very dim
	}
	r.drawRectBorder(screen, ui.RebootBtnRect, clr)
	label := fmt.Sprintf("SYSTEM_REBOOT (REQ: %s)", format.FormatBits(threshold))
	r.DrawText(screen, label, ui.RebootBtnRect.Min.X+ui.RebootBtnRect.Dx()/2-120, ui.RebootBtnRect.Min.Y+20, color.White)
}

func (r *Renderer) drawPacketIntercept(screen *ebiten.Image, state *model.GameState, neonGreen color.Color) {
	rect := ui.PacketRect
	clr := neonGreen
	if (int(state.PacketTimer*10) % 4) < 2 { clr = color.White }
	ebitenutil.DrawRect(screen, float64(rect.Min.X), float64(rect.Min.Y), float64(rect.Dx()), float64(rect.Dy()), color.RGBA{0, 100, 0, 100})
	r.drawRectBorder(screen, rect, clr)
	r.DrawText(screen, ">> PACKET <<", rect.Min.X+10, rect.Min.Y+20, color.White)
}

func (r *Renderer) drawRebootDialog(screen *ebiten.Image, state *model.GameState, neonGreen color.Color) {
	rect := image.Rect(ui.WidgetX+100, ui.WidgetY+100, ui.WidgetX+ui.WidgetWidth-100, ui.WidgetY+ui.WidgetHeight-100)
	ebitenutil.DrawRect(screen, float64(rect.Min.X), float64(rect.Min.Y), float64(rect.Dx()), float64(rect.Dy()), color.Black)
	r.drawRectBorder(screen, rect, neonGreen)

	threshold := state.GetRebootThreshold()
	gain := math.Log10(state.TotalBitsEarned/threshold+1.0) * 0.1
	if state.RebootCount == 0 {
		gain += 0.1
	}

	r.DrawText(screen, ">> SYSTEM_REBOOT INITIATED <<", rect.Min.X+150, rect.Min.Y+50, color.White)
	r.DrawText(screen, fmt.Sprintf("GHz GAIN: +%.3fX", gain), rect.Min.X+150, rect.Min.Y+100, color.White)
	r.DrawText(screen, "[Y] CONFIRM / [N] ABORT", rect.Min.X+150, rect.Min.Y+200, color.White)
}

func (r *Renderer) drawRectBorder(screen *ebiten.Image, rect image.Rectangle, clr color.Color) {
	ebitenutil.DrawLine(screen, float64(rect.Min.X), float64(rect.Min.Y), float64(rect.Max.X), float64(rect.Min.Y), clr)
	ebitenutil.DrawLine(screen, float64(rect.Max.X), float64(rect.Min.Y), float64(rect.Max.X), float64(rect.Max.Y), clr)
	ebitenutil.DrawLine(screen, float64(rect.Max.X), float64(rect.Max.Y), float64(rect.Min.X), float64(rect.Max.Y), clr)
	ebitenutil.DrawLine(screen, float64(rect.Min.X), float64(rect.Max.Y), float64(rect.Min.X), float64(rect.Min.Y), clr)
}
