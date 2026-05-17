package input

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/popolque/firstbitengi/internal/ui"
)

type InputSystem struct {
	MousePos      image.Point
	Clicked       bool // left mouse button just pressed
	ScrollDelta   int  // mouse wheel delta this frame
	RebootPressed bool // keyboard R or REBOOT button clicked
}

func NewInputSystem() *InputSystem {
	return &InputSystem{}
}

func (in *InputSystem) Poll() {
	mx, my := ebiten.CursorPosition()
	in.MousePos = image.Pt(mx, my)
	in.Clicked = inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft)
	_, wy := ebiten.Wheel()
	in.ScrollDelta = int(wy)
	in.RebootPressed = inpututil.IsKeyJustPressed(ebiten.KeyR)
}

// Hit-testing helpers used by engine
func (in *InputSystem) ClickerPressed() bool {
	// Support both mouse click in region and Space key
	mouseInClicker := in.Clicked && in.MousePos.In(ui.ClickerRegion)
	spacePressed := inpututil.IsKeyJustPressed(ebiten.KeySpace)
	return mouseInClicker || spacePressed
}

func (in *InputSystem) HardwareBuyPressed(rowRect image.Rectangle) bool {
	return in.Clicked && in.MousePos.In(rowRect)
}

func (in *InputSystem) RebootTriggered() bool {
	return in.RebootPressed || (in.Clicked && in.MousePos.In(ui.RebootBtnRect))
}
