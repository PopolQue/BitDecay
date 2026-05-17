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
    ____  ____________________  _______________  __  __
   / __ )/  _/_  __/ ____/ __ \/ ____/ ____/   \ \ \/ /
  / __  |/ /  / / / __/ / / / / __/ / /   / /| |  \  / 
 / /_/ // /  / / / /___/ /_/ / /___/ /___/ ___ |  / /  
/____/___/ /_/ /_____/_____/_____/\____/_/  |_| /_/   

`

func (r *Renderer) DrawText(screen *ebiten.Image, str string, x, y float64, face *text.GoTextFace, clr color.Color) {
	if face == nil {
		return
	}
	op := &text.DrawOptions{}
	op.GeoM.Translate(x, y)
	op.ColorScale.ScaleWithColor(clr)
	text.Draw(screen, str, face, op)
}

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
	if r.normalFace == nil {
		ebitenutil.DebugPrintAt(screen, "BIT-DECAY (FONT MISSING)", ui.ScreenWidth/2-60, 20)
		return
	}

	// Manual height calculation because text.Measure seems unreliable for multi-line here
	lineHeight := 22.0
	lines := 5.0
	totalH := lines * lineHeight
	
	w, _ := text.Measure(asciiHeader, r.normalFace, 0)
	x := (float64(ui.ScreenWidth) - w) / 2
	y := float64(ui.WidgetY) - totalH - 40 // More padding

	if y < 10 {
		y = 10
	}

	// Draw black background to prevent matrix overlap
	ebitenutil.DrawRect(screen, x-20, y-10, w+40, totalH+20, color.Black)

	r.DrawText(screen, asciiHeader, x, y, r.normalFace, color.RGBA{57, 255, 20, 255})
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
		textClr := color.RGBA{0, 150, 0, 255}
		if state.ActiveTab == t.name {
			clr = color.RGBA{0, 255, 0, 255}
			textClr = color.RGBA{0, 255, 0, 255}
			// Fill active tab slightly
			ebitenutil.DrawRect(screen, float64(t.rect.Min.X), float64(t.rect.Min.Y),
				float64(t.rect.Dx()), float64(t.rect.Dy()), color.RGBA{0, 40, 0, 255})
		}
		r.drawRectBorder(screen, t.rect, clr)
		r.DrawText(screen, t.name, float64(t.rect.Min.X+20), float64(t.rect.Min.Y+10), r.smallFace, textClr)
	}
}

func (r *Renderer) drawHUD(screen *ebiten.Image, state *model.GameState) {
	r.drawRectBorder(screen, ui.MetricsHUDRect, color.RGBA{0, 80, 0, 255})
	green := color.RGBA{0, 200, 0, 255}
	r.DrawText(screen, fmt.Sprintf("BITS: %s", format.FormatBits(state.Bits)), float64(ui.MetricsHUDRect.Min.X+10), float64(ui.MetricsHUDRect.Min.Y+10), r.normalFace, green)
	r.DrawText(screen, fmt.Sprintf("TOTAL: %s", format.FormatBits(state.TotalBitsEarned)), float64(ui.MetricsHUDRect.Min.X+10), float64(ui.MetricsHUDRect.Min.Y+40), r.smallFace, green)
	r.DrawText(screen, fmt.Sprintf("ENTROPY: %.1f%%", state.Entropy), float64(ui.MetricsHUDRect.Min.X+10), float64(ui.MetricsHUDRect.Min.Y+65), r.smallFace, green)
	r.DrawText(screen, fmt.Sprintf("CORRUPTION: %.1f%%", state.Corruption), float64(ui.MetricsHUDRect.Min.X+10), float64(ui.MetricsHUDRect.Min.Y+85), r.smallFace, green)
	r.DrawText(screen, fmt.Sprintf("GHz MULT: %.3fX", state.GHzMultiplier), float64(ui.MetricsHUDRect.Min.X+300), float64(ui.MetricsHUDRect.Min.Y+10), r.normalFace, green)
}

func (r *Renderer) drawHardwareList(screen *ebiten.Image, state *model.GameState) {
	x := float64(ui.ListRect.Min.X + 10)
	startY := float64(ui.ListRect.Min.Y + 10)
	rowHeight := 60.0
	green := color.RGBA{0, 180, 0, 255}
	brightGreen := color.RGBA{0, 255, 0, 255}

	for i, def := range model.AllHardware {
		y := startY + float64(i)*rowHeight - float64(state.ScrollOffset)
		if y < float64(ui.ListRect.Min.Y) || y > float64(ui.ListRect.Max.Y-int(rowHeight)) {
			continue
		}

		owned := state.Hardware[def.ID]
		cost := model.CurrentCost(def, owned)

		clr := green
		status := "[ ]"
		if state.Bits >= cost {
			status = "[!]"
			clr = brightGreen
		}

		label := fmt.Sprintf("%s %s (Owned: %d) - Cost: %s", status, def.Name, owned, format.FormatBits(cost))
		r.DrawText(screen, label, x, y, r.smallFace, clr)
		r.DrawText(screen, fmt.Sprintf("  BPS: %.1f | Entropy: %.2f", def.BaseBPS, def.EntropyWeight), x, y+20, r.smallFace, green)
	}
}

func (r *Renderer) drawUpgradeList(screen *ebiten.Image, state *model.GameState) {
	x := float64(ui.ListRect.Min.X + 10)
	startY := float64(ui.ListRect.Min.Y + 10)
	rowHeight := 60.0
	green := color.RGBA{0, 180, 0, 255}
	brightGreen := color.RGBA{0, 255, 0, 255}
	purchasedClr := color.RGBA{0, 100, 0, 255}

	for i, def := range model.AllUpgrades {
		y := startY + float64(i)*rowHeight - float64(state.ScrollOffset)
		if y < float64(ui.ListRect.Min.Y) || y > float64(ui.ListRect.Max.Y-int(rowHeight)) {
			continue
		}

		owned := state.Upgrades[def.ID]
		if owned {
			r.DrawText(screen, fmt.Sprintf("[X] %s (PURCHASED)", def.Name), x, y, r.smallFace, purchasedClr)
			r.DrawText(screen, "  "+def.Description, x, y+20, r.smallFace, purchasedClr)
			continue
		}

		clr := green
		status := "[ ]"
		if state.Bits >= def.Cost {
			status = "[*]"
			clr = brightGreen
		}

		label := fmt.Sprintf("%s %s - Cost: %s", status, def.Name, format.FormatBits(def.Cost))
		r.DrawText(screen, label, x, y, r.smallFace, clr)
		r.DrawText(screen, "  "+def.Description, x, y+20, r.smallFace, green)
	}
}

func (r *Renderer) drawSystemTab(screen *ebiten.Image, state *model.GameState) {
	// Clicker Button
	clickerColor := color.RGBA{0, 255, 0, 255}
	if state.ClickerFlash {
		clickerColor = color.RGBA{100, 255, 100, 255}
	}
	r.drawRectBorder(screen, ui.ClickerRegion, clickerColor)
	r.DrawText(screen, "MANUAL OVERRIDE", float64(ui.ClickerRegion.Min.X+30), float64(ui.ClickerRegion.Min.Y+30), r.normalFace, clickerColor)

	// Reboot Button
	rebootColor := color.RGBA{0, 100, 0, 255}
	if state.TotalBitsEarned >= 1_000_000 {
		rebootColor = color.RGBA{0, 255, 0, 255}
	}
	r.drawRectBorder(screen, ui.RebootBtnRect, rebootColor)
	r.DrawText(screen, "SYSTEM_REBOOT", float64(ui.RebootBtnRect.Min.X+ui.RebootBtnRect.Dx()/2-80), float64(ui.RebootBtnRect.Min.Y+12), r.normalFace, rebootColor)
}

func (r *Renderer) drawRebootDialog(screen *ebiten.Image, state *model.GameState) {
	rect := image.Rect(ui.WidgetX+50, ui.WidgetY+150, ui.WidgetX+ui.WidgetWidth-50, ui.WidgetY+ui.WidgetHeight-150)
	ebitenutil.DrawRect(screen, float64(rect.Min.X), float64(rect.Min.Y), float64(rect.Dx()), float64(rect.Dy()), color.RGBA{0, 0, 0, 230})
	r.drawRectBorder(screen, rect, color.RGBA{0, 255, 0, 255})

	gain := math.Log10(state.TotalBitsEarned/1_000_000) * 0.1
	if gain < 0 {
		gain = 0
	}

	green := color.RGBA{0, 255, 0, 255}
	r.DrawText(screen, ">> SYSTEM_REBOOT INITIATED <<", float64(rect.Min.X+200), float64(rect.Min.Y+50), r.normalFace, green)
	r.DrawText(screen, fmt.Sprintf("GHz GAIN: +%.3fX", gain), float64(rect.Min.X+200), float64(rect.Min.Y+100), r.normalFace, green)
	r.DrawText(screen, "ALL HARDWARE AND BITS WILL BE WIPED.", float64(rect.Min.X+200), float64(rect.Min.Y+150), r.smallFace, green)
	r.DrawText(screen, "[CONFIRM (Left)]    [ABORT (Right)]", float64(rect.Min.X+200), float64(rect.Min.Y+220), r.normalFace, green)
}

func (r *Renderer) drawRectBorder(screen *ebiten.Image, rect image.Rectangle, clr color.Color) {
	ebitenutil.DrawLine(screen, float64(rect.Min.X), float64(rect.Min.Y), float64(rect.Max.X), float64(rect.Min.Y), clr)
	ebitenutil.DrawLine(screen, float64(rect.Max.X), float64(rect.Min.Y), float64(rect.Max.X), float64(rect.Max.Y), clr)
	ebitenutil.DrawLine(screen, float64(rect.Max.X), float64(rect.Max.Y), float64(rect.Min.X), float64(rect.Max.Y), clr)
	ebitenutil.DrawLine(screen, float64(rect.Min.X), float64(rect.Max.Y), float64(rect.Min.X), float64(rect.Min.Y), clr)
}
