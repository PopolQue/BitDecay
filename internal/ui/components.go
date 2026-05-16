package ui

import (
	"fmt"
	"image/color"

	"github.com/popolque/firstbitengi/internal/game"
	"github.com/popolque/firstbitengi/internal/models"
	"github.com/google/uuid"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

func (r *Renderer) drawMenu(screen *ebiten.Image) {
	screen.Fill(color.RGBA{10, 50, 80, 255})
	ebitenutil.DebugPrintAt(screen, "CAMPERS - SEASON 1", 320, 150)
	vector.DrawFilledRect(screen, 300, 245, 200, 40, color.RGBA{100, 100, 100, 255}, true)
	ebitenutil.DebugPrintAt(screen, "[S] START SOLO GAME", 330, 258)
	vector.DrawFilledRect(screen, 300, 305, 200, 40, color.RGBA{100, 100, 100, 255}, true)
	ebitenutil.DebugPrintAt(screen, "[T] TUTORIAL", 355, 318)
	vector.DrawFilledRect(screen, 300, 365, 200, 40, color.RGBA{100, 100, 100, 255}, true)
	ebitenutil.DebugPrintAt(screen, "[M] MULTIPLAYER", 345, 378)
	ebitenutil.DebugPrintAt(screen, "Press corresponding key to choose", 310, 450)
}

func (r *Renderer) drawTopBar(screen *ebiten.Image) {
	res := r.Engine.State.Resources
	status := fmt.Sprintf("Round: %d/12 | AP: %d | Money: $%d | Power: %d | Water: %d",
		r.Engine.State.Round, r.Engine.State.ActionPoints, res.Money, res.Power, res.Water)
	vector.DrawFilledRect(screen, 0, 0, 800, 30, color.RGBA{0, 0, 0, 180}, true)
	ebitenutil.DebugPrintAt(screen, status, 10, 10)
}

func (r *Renderer) drawActiveEvent(screen *ebiten.Image) {
	if r.Engine.State.ActiveEvent != nil {
		ev := r.Engine.State.ActiveEvent
		vector.DrawFilledRect(screen, 10, 40, 250, 45, color.RGBA{200, 50, 50, 200}, true)
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("EVENT: %s", ev.Name), 15, 45)
		ebitenutil.DebugPrintAt(screen, ev.Description, 15, 65)
	}
}

func (r *Renderer) drawAssetIcon(screen *ebiten.Image, x, y float32, asset models.Asset) {
	icon := "🏠"
	switch asset.Type {
	case models.AssetTent: icon = "⛺"
	case models.AssetCaravan: icon = "🚐"
	case models.AssetGenerator: icon = "⚡"
	case models.AssetWaterTank: icon = "💧"
	}
	if len(asset.Occupants) > 0 {
		vector.DrawFilledCircle(screen, x, y, 15, color.RGBA{255, 255, 0, 150}, true)
	}
	ebitenutil.DebugPrintAt(screen, icon, int(x)-8, int(y)-8)
}

func (r *Renderer) drawNPCPool(screen *ebiten.Image) {
	ebitenutil.DebugPrintAt(screen, "GUEST POOL", 620, 50)
	for i, npc := range r.Engine.State.GuestPool {
		y := 80 + i*70
		if r.Engine.State.SelectedNPCID == npc.ID {
			vector.DrawFilledRect(screen, 618, float32(y)-2, 174, 64, color.RGBA{255, 255, 255, 255}, true)
		}
		vector.DrawFilledRect(screen, 620, float32(y), 170, 60, color.RGBA{50, 50, 50, 200}, true)
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%s (%s)", npc.Name, npc.Type), 625, y+5)
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("S: %d | Pay: %d", npc.GroupSize, npc.IncomePerNight), 625, y+25)
	}
}

func (r *Renderer) drawInstructions(screen *ebiten.Image) {
	ebitenutil.DebugPrintAt(screen, "Click to select Tile/NPC | 1-3: Buy Asset | T: Buy Land ($100) | A: Accept", 10, 530)
	ebitenutil.DebugPrintAt(screen, "U: Upgrade Asset | SPACE: End Turn | ESC: Menu", 10, 550)
}

func (r *Renderer) HandleInput() {
	q, r_coord, _ := game.ParseTile(r.Engine.State.SelectedTile)
	
	if ebiten.IsKeyPressed(ebiten.Key1) {
		r.Engine.BuyAsset(models.AssetTent, q, r_coord)
	}
	if ebiten.IsKeyPressed(ebiten.Key2) {
		r.Engine.BuyAsset(models.AssetCaravan, q, r_coord)
	}
	if ebiten.IsKeyPressed(ebiten.Key3) {
		r.Engine.BuyAsset(models.AssetGenerator, q, r_coord)
	}
	if ebiten.IsKeyPressed(ebiten.KeyT) {
		r.Engine.BuyTile(q, r_coord)
	}
	if ebiten.IsKeyPressed(ebiten.KeyU) && r.Engine.State.SelectedAssetID != uuid.Nil {
		r.Engine.UpgradeAsset(r.Engine.State.SelectedAssetID)
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		if r.Engine.State.SelectedNPCID != uuid.Nil && r.Engine.State.SelectedAssetID != uuid.Nil {
			r.Engine.AcceptGuest(r.Engine.State.SelectedNPCID, r.Engine.State.SelectedAssetID)
		}
	}
}
