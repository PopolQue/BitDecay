package render

import (
	_ "embed"
	"fmt"
	"image"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/popolque/firstbitengi/internal/format"
	"github.com/popolque/firstbitengi/internal/model"
	"github.com/popolque/firstbitengi/internal/ui"
)

//go:embed crt.kage
var crtShaderSrc []byte

type Renderer struct {
	waterfall *WaterfallRenderer
	glitch    *GlitchSystem
	crtShader *ebiten.Shader
	offscreen *ebiten.Image
}

func NewRenderer() *Renderer {
	s, err := ebiten.NewShader(crtShaderSrc)
	if err != nil {
		fmt.Printf("failed to load shader: %v\n", err)
	}

	return &Renderer{
		waterfall: NewWaterfallRenderer(),
		glitch:    NewGlitchSystem(),
		crtShader: s,
		offscreen: ebiten.NewImage(ui.ScreenWidth, ui.ScreenHeight),
	}
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

func (r *Renderer) drawToImage(screen *ebiten.Image, state *model.GameState) {
	// Update visual systems
	r.waterfall.Update(state.Corruption)
	r.glitch.Step(state.Corruption)

	// Draw background
	screen.Fill(color.RGBA{0, 5, 0, 255})

	// Draw matrix around the widget
	r.waterfall.Draw(screen)

	// ASCII Header
	r.drawHeader(screen)

	// Central Widget
	r.drawWidget(screen, state)

	// Glitch Overlay
	r.glitch.Draw(screen)
}

func (r *Renderer) drawHeader(screen *ebiten.Image) {
	lines := 8
	charWidth := 6
	headerWidth := 67 * charWidth
	x := (ui.ScreenWidth - headerWidth) / 2
	y := ui.WidgetY - (lines * 14) - 5

	// Draw to a temporary image to apply neon green color scale
	tempImg := ebiten.NewImage(headerWidth+20, lines*14+10)
	ebitenutil.DebugPrintAt(tempImg, asciiHeader, 0, 0)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(x), float64(y))
	op.ColorScale.ScaleWithColor(color.RGBA{57, 255, 20, 255}) // Neon Green
	screen.DrawImage(tempImg, op)
}

func (r *Renderer) drawWidget(screen *ebiten.Image, state *model.GameState) {
	// Widget Background (darker than screen)
	ebitenutil.DrawRect(screen, float64(ui.WidgetRect.Min.X), float64(ui.WidgetRect.Min.Y),
		float64(ui.WidgetRect.Dx()), float64(ui.WidgetRect.Dy()), color.RGBA{0, 15, 0, 240})
	r.drawRectBorder(screen, ui.WidgetRect, color.RGBA{0, 150, 0, 255})

	// Tabs
	r.drawTabs(screen, state)

	// HUD (Always visible in content area?)
	r.drawHUD(screen, state)

	// Content based on ActiveTab
	switch state.ActiveTab {
	case "HARDWARE":
		r.drawHardwareList(screen, state)
	case "UPGRADES":
		r.drawUpgradeList(screen, state)
	case "SYSTEM":
		r.drawSystemTab(screen, state)
	}

	// Reboot Dialog (Overlays everything in widget)
	if state.RebootPending {
		r.drawRebootDialog(screen, state)
	}
}

func (r *Renderer) drawTabs(screen *ebiten.Image, state *model.GameState) {
	tabs := []struct {
		rect image.Rectangle
		name string
	}{
		{ui.Tab1Rect, "HARDWARE"},
		{ui.Tab2Rect, "UPGRADES"},
		{ui.Tab3Rect, "SYSTEM"},
	}

	for _, t := range tabs {
		clr := color.RGBA{0, 100, 0, 255}
		if state.ActiveTab == t.name {
			clr = color.RGBA{0, 255, 0, 255}
			// Fill active tab slightly
			ebitenutil.DrawRect(screen, float64(t.rect.Min.X), float64(t.rect.Min.Y),
				float64(t.rect.Dx()), float64(t.rect.Dy()), color.RGBA{0, 40, 0, 255})
		}
		r.drawRectBorder(screen, t.rect, clr)
		ebitenutil.DebugPrintAt(screen, t.name, t.rect.Min.X+20, t.rect.Min.Y+12)
	}
}

func (r *Renderer) drawHUD(screen *ebiten.Image, state *model.GameState) {
	r.drawRectBorder(screen, ui.MetricsHUDRect, color.RGBA{0, 80, 0, 255})
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("BITS: %s", format.FormatBits(state.Bits)), ui.MetricsHUDRect.Min.X+10, ui.MetricsHUDRect.Min.Y+10)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("TOTAL: %s", format.FormatBits(state.TotalBitsEarned)), ui.MetricsHUDRect.Min.X+10, ui.MetricsHUDRect.Min.Y+30)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("ENTROPY: %.1f%%", state.Entropy), ui.MetricsHUDRect.Min.X+10, ui.MetricsHUDRect.Min.Y+50)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("CORRUPTION: %.1f%%", state.Corruption), ui.MetricsHUDRect.Min.X+10, ui.MetricsHUDRect.Min.Y+70)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("GHz MULT: %.3fX", state.GHzMultiplier), ui.MetricsHUDRect.Min.X+250, ui.MetricsHUDRect.Min.Y+10)
}

func (r *Renderer) drawHardwareList(screen *ebiten.Image, state *model.GameState) {
	x := ui.ListRect.Min.X + 10
	startY := ui.ListRect.Min.Y + 10
	rowHeight := 60

	for i, def := range model.AllHardware {
		y := startY + i*rowHeight - state.ScrollOffset
		if y < ui.ListRect.Min.Y || y > ui.ListRect.Max.Y-rowHeight {
			continue
		}

		owned := state.Hardware[def.ID]
		cost := model.CurrentCost(def, owned)

		status := "[ ]"
		if state.Bits >= cost {
			status = "[!]"
		}

		label := fmt.Sprintf("%s %s (Owned: %d) - Cost: %s", status, def.Name, owned, format.FormatBits(cost))
		ebitenutil.DebugPrintAt(screen, label, x, y)
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("  BPS: %.1f | Entropy: %.2f", def.BaseBPS, def.EntropyWeight), x, y+20)
	}
}

func (r *Renderer) drawUpgradeList(screen *ebiten.Image, state *model.GameState) {
	x := ui.ListRect.Min.X + 10
	startY := ui.ListRect.Min.Y + 10
	rowHeight := 60

	for i, def := range model.AllUpgrades {
		y := startY + i*rowHeight - state.ScrollOffset
		if y < ui.ListRect.Min.Y || y > ui.ListRect.Max.Y-rowHeight {
			continue
		}

		owned := state.Upgrades[def.ID]
		if owned {
			ebitenutil.DebugPrintAt(screen, fmt.Sprintf("[X] %s (PURCHASED)", def.Name), x, y)
			ebitenutil.DebugPrintAt(screen, "  "+def.Description, x, y+20)
			continue
		}

		status := "[ ]"
		if state.Bits >= def.Cost {
			status = "[*]"
		}

		label := fmt.Sprintf("%s %s - Cost: %s", status, def.Name, format.FormatBits(def.Cost))
		ebitenutil.DebugPrintAt(screen, label, x, y)
		ebitenutil.DebugPrintAt(screen, "  "+def.Description, x, y+20)
	}
}

func (r *Renderer) drawSystemTab(screen *ebiten.Image, state *model.GameState) {
	// Clicker Button
	clickerColor := color.RGBA{0, 255, 0, 255}
	if state.ClickerFlash {
		clickerColor = color.RGBA{100, 255, 100, 255}
	}
	r.drawRectBorder(screen, ui.ClickerRegion, clickerColor)
	ebitenutil.DebugPrintAt(screen, "MANUAL OVERRIDE (CLICK)", ui.ClickerRegion.Min.X+30, ui.ClickerRegion.Min.Y+30)

	// Reboot Button
	rebootColor := color.RGBA{0, 100, 0, 255}
	if state.TotalBitsEarned >= 1_000_000 {
		rebootColor = color.RGBA{0, 255, 0, 255}
	}
	r.drawRectBorder(screen, ui.RebootBtnRect, rebootColor)
	ebitenutil.DebugPrintAt(screen, "SYSTEM_REBOOT", ui.RebootBtnRect.Min.X+ui.RebootBtnRect.Dx()/2-50, ui.RebootBtnRect.Min.Y+12)
}

func (r *Renderer) drawRebootDialog(screen *ebiten.Image, state *model.GameState) {
	rect := image.Rect(ui.WidgetX+50, ui.WidgetY+150, ui.WidgetX+ui.WidgetWidth-50, ui.WidgetY+ui.WidgetHeight-150)
	ebitenutil.DrawRect(screen, float64(rect.Min.X), float64(rect.Min.Y), float64(rect.Dx()), float64(rect.Dy()), color.RGBA{0, 0, 0, 230})
	r.drawRectBorder(screen, rect, color.RGBA{0, 255, 0, 255})

	gain := math.Log10(state.TotalBitsEarned/1_000_000) * 0.1
	if gain < 0 {
		gain = 0
	}

	ebitenutil.DebugPrintAt(screen, ">> SYSTEM_REBOOT INITIATED <<", rect.Min.X+250, rect.Min.Y+50)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("GHz GAIN: +%.3fX", gain), rect.Min.X+250, rect.Min.Y+100)
	ebitenutil.DebugPrintAt(screen, "ALL HARDWARE AND BITS WILL BE WIPED.", rect.Min.X+250, rect.Min.Y+150)
	ebitenutil.DebugPrintAt(screen, "[CONFIRM (Left)]    [ABORT (Right)]", rect.Min.X+250, rect.Min.Y+220)
}

func (r *Renderer) drawRectBorder(screen *ebiten.Image, rect image.Rectangle, clr color.Color) {
	ebitenutil.DrawLine(screen, float64(rect.Min.X), float64(rect.Min.Y), float64(rect.Max.X), float64(rect.Min.Y), clr)
	ebitenutil.DrawLine(screen, float64(rect.Max.X), float64(rect.Min.Y), float64(rect.Max.X), float64(rect.Max.Y), clr)
	ebitenutil.DrawLine(screen, float64(rect.Max.X), float64(rect.Max.Y), float64(rect.Min.X), float64(rect.Max.Y), clr)
	ebitenutil.DrawLine(screen, float64(rect.Min.X), float64(rect.Max.Y), float64(rect.Min.X), float64(rect.Min.Y), clr)
}
