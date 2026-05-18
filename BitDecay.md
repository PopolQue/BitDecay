# Technical Design Document

## BIT-DECAY // MNEMONIC_OVERRIDE_3.0

## Implementation Stack: Go + Ebitengine

> **Document Version:** 1.0  
> **Stack:** Go 1.24+ · Ebitengine v2.9+  
> **Target Platforms:** Windows · macOS · Linux · (WebAssembly)  
> **Document Status:** DRAFT

---

## Table of Contents

1. [Architecture Overview](#1-architecture-overview)
2. [Project Structure](#2-project-structure)
3. [Ebitengine Game Loop Integration](#3-ebitengine-game-loop-integration)
4. [Core Engine & State](#4-core-engine--state)
5. [Rendering Pipeline](#5-rendering-pipeline)
6. [Input System](#6-input-system)
7. [Glitch & Corruption Visual System](#7-glitch--corruption-visual-system)
8. [Hardware & Economy Systems](#8-hardware--economy-systems)
9. [Persistence & Autosave](#9-persistence--autosave)
10. [Prestige System (SYSTEM_REBOOT)](#10-prestige-system-system_reboot)
11. [Audio](#11-audio)
12. [Configuration & Balancing](#12-configuration--balancing)
13. [Build & Deployment](#13-build--deployment)
14. [Risk Register](#14-risk-register)

---

## 1. Architecture Overview

Ebitengine provides a single fixed-timestep loop via `Update()` + `Draw()`. Unlike Fyne's retained widget model, **all rendering is immediate-mode**: every frame redraws the screen from scratch by blitting images and drawing primitives onto `*ebiten.Image` targets. This gives total pixel-level control — ideal for the CRT glitch aesthetic.

``` Runtime
┌─────────────────────────────────────────────────────────────┐
│                     EBITENGINE RUNTIME                      │
│                                                             │
│  ebiten.RunGame(game)                                       │
│       │                                                     │
│       ├── Update() — called @ 60 TPS (fixed)                │
│       │       ├── InputSystem.Poll()                        │
│       │       ├── GameEngine.Tick() [every 6th → 100ms]     │
│       │       │       ├── HardwareRegistry.Compute()        │
│       │       │       ├── EntropyEngine.Step()              │
│       │       │       └── DecayEngine.Step()                │
│       │       ├── UIState.Animate()                         │
│       │       └── GlitchSystem.Step()                       │
│       │                                                     │
│       └── Draw(*ebiten.Image) — called every frame          │
│               ├── RenderBackground()                        │
│               ├── RenderWaterfall()                         │
│               ├── RenderHUD()                               │
│               ├── RenderHardwarePanel()                     │
│               ├── RenderClickButton()                       │
│               └── RenderGlitchOverlay()                     │
└─────────────────────────────────────────────────────────────┘
```

### Key Design Decisions — Why Ebitengine

| Concern | Ebitengine Approach |
| --- | --- |
| Full pixel control | Every frame drawn to `*ebiten.Image`; GPU-accelerated via OpenGL/Metal/DirectX |
| CRT scanline & glitch | Custom shader (Kage) for phosphor bloom, scanlines, chromatic aberration |
| Waterfall animation | Pixel-buffer written per-frame; no retained widget overhead |
| Cross-platform + WASM | `GOOS=js GOARCH=wasm` for browser deployment |
| No OS UI dependencies | Pure game window; no widget toolkit dependency |

---

## 2. Project Structure

``` Filestructure
bit-decay/
├── cmd/
│   └── bitdecay/
│       └── main.go               # ebiten.RunGame entrypoint
├── internal/
│   ├── engine/
│   │   ├── game.go               # Implements ebiten.Game interface
│   │   ├── gameengine.go         # Sub-tick accumulator, orchestrates subsystems
│   │   ├── entropy.go            # Entropy accumulation & decay math
│   │   ├── hardware.go           # Production computation
│   │   └── prestige.go           # SYSTEM_REBOOT logic
│   ├── model/
│   │   ├── state.go              # GameState (canonical truth)
│   │   ├── hardware_def.go       # Static hardware table
│   │   └── upgrade_def.go        # Upgrade tree
│   ├── render/
│   │   ├── renderer.go           # Top-level Draw() dispatcher
│   │   ├── crt.kage              # Kage shader: scanlines + phosphor glow
│   ├── input/
│   │   └── input.go              # Mouse + keyboard polling, click detection
│   ├── ui/
│   │   ├── layout.go             # Screen region constants
│   │   ├── font.go               # go-text/renderv2 font loader
│   │   └── scroll.go             # Generic scrollable region helper
│   ├── persist/
│   │   ├── save.go               # JSON save / load
│   │   └── autosave.go           # Background goroutine (30s interval)
│   └── format/
│       └── bits.go               # Bit → Brontobyte formatter
├── assets/
│   ├── assets.go                 # //go:embed for all assets
│   ├── fonts/
│   │   └── BPdotsLight.otf
│   └── sfx/
│       ├── BitDecayClick.wav
│       └── BitDecayLoop.wav
├── config/
│   └── balance.toml
├── save.json
├── go.mod
└── go.sum
```

---

## 3. Ebitengine Game Loop Integration

### 3.1 The `Game` Struct

```go
// internal/engine/game.go
package engine

import "github.com/hajimehoshi/ebiten/v2"

type Game struct {
    engine   *GameEngine
    renderer *render.Renderer
    input    *input.InputSystem
    uiState  *UIAnimationState
    glitch   *GlitchSystem
}

func (g *Game) Update() error {
    g.input.Poll()
    g.engine.Update(g.input)   // accumulates sub-ticks; fires game logic at 100ms
    g.uiState.Animate()
    g.glitch.Step(g.engine.State().Corruption)
    return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
    g.renderer.Draw(screen, g.engine.State(), g.uiState, g.glitch)
}

func (g *Game) Layout(outsideW, outsideH int) (int, int) {
    return 1280, 768 // fixed logical resolution; Ebitengine scales to window
}
```

### 3.2 Sub-Tick Accumulator

Ebitengine calls `Update()` at 60 TPS. Game logic must tick at 100ms (10 TPS). An accumulator converts between the two rates without using a separate goroutine.

```go
// internal/engine/gameengine.go
const GameTickMs = 100.0 // ms per game logic tick
const UpdateMs   = 1000.0 / 60.0 // ~16.67ms per Ebitengine Update

type GameEngine struct {
    state       *model.GameState
    accumMs     float64
    hardware    *HardwareRegistry
    entropyEng  *EntropyEngine
}

func (ge *GameEngine) Update(in *input.InputSystem) {
    ge.accumMs += UpdateMs
    for ge.accumMs >= GameTickMs {
        ge.accumMs -= GameTickMs
        ge.gameTick(GameTickMs / 1000.0) // deltaSeconds
    }
    // Manual click is processed every Update (not just game tick)
    if in.ClickerPressed() {
        ge.state.Bits += ge.manualClickValue()
        ge.state.ClickerFlash = true
    }
}

func (ge *GameEngine) gameTick(dt float64) {
    ge.hardware.Compute(ge.state, dt)
    ge.entropyEng.Step(ge.state, dt)
    applyDecay(ge.state, dt)
    clampState(ge.state)
}
```

---

## 4. Core Engine & State

### 4.1 GameState

```go
// internal/model/state.go
package model

type GameState struct {
    // Economy
    Bits            float64
    TotalBitsEarned float64

    // System health
    Entropy         float64 // 0.0 – 100.0
    Corruption      float64 // 0.0 – 100.0

    // Prestige
    GHzMultiplier   float64
    RebootCount     int

    // Hardware owned
    Hardware        map[string]int
    Upgrades        map[string]bool

    // UI animation signals (not persisted)
    ClickerFlash    bool   // consumed by renderer
    ScrollOffset    int    // hardware panel scroll position
    RebootPending   bool   // confirmation dialog visible
}
```

> **Thread Safety Note:** Unlike the Fyne implementation, there is **no mutex** on `GameState`. Ebitengine guarantees that `Update()` and `Draw()` are never called concurrently. The only background goroutine (autosave) works from a periodic **snapshot copy** of the state, not a pointer, avoiding races entirely.

### 4.2 Production & Entropy Formulas

``` formula
bitsPerSec = Σ(hardware[id] × bps[id] × upgradeMult[id]) × GHzMultiplier × (1 − corruptPenalty)
corruptPenalty = min(Corruption / 200.0, 0.5)

entropyDelta   = Σ(hardware[id] × entropyWeight[id]) × dt
if Entropy > 50:
    corruptDelta = (Entropy − 50) × 0.002 × dt
if Corruption > 75:
    decayRate    = Bits × (Corruption − 75) × 0.0001
    Bits        -= decayRate × dt
```

---

## 5. Rendering Pipeline

### 5.1 Layer Order

Each frame, `renderer.Draw()` paints layers in order onto the `screen *ebiten.Image`:

``` UI Layers
Layer 0 — Background grid (static, dark #000A00 fill + dim grid lines)
Layer 1 — Waterfall (scrolling pixel buffer blitted as texture)
Layer 2 — HUD panel (metrics text, progress bars drawn with vector rects)
Layer 3 — Hardware panel (clipped scrollable region with item rows)
Layer 4 — Click button (animated 0/1 glyph)
Layer 5 — CRT shader pass (full-screen Kage shader: scanlines + glow)
Layer 6 — Glitch overlay (horizontal tears, palette shift, at Corruption >75%)
Layer 7 — Dialog / confirmation (SYSTEM_REBOOT prompt, if pending)
```

### 5.2 CRT Kage Shader

The CRT post-process runs as a full-screen Kage shader applied to an off-screen `*ebiten.Image` that all lower layers render into. This is the cleanest way to apply global effects without per-widget overhead.

```kage
// shader/crt.kage
//go:build ignore

package main

var Time float

func Fragment(dst vec4, src vec2, color vec4) vec4 {
    // Scanline darkening
    scanline := sin(src.y * 3.14159 * 768.0)
    scanlineMask := clamp(scanline * 0.15 + 0.85, 0.0, 1.0)

    // Phosphor glow: sample with slight UV spread
    col := imageSrc0At(src)
    glow := imageSrc0At(src + vec2(0.001, 0.0)) * 0.15
    col = col + glow

    // Green channel boost for phosphor effect
    col.g = clamp(col.g * 1.1, 0.0, 1.0)

    return col * vec4(scanlineMask, scanlineMask, scanlineMask, 1.0)
}
```

```go
// Loading and applying the shader
var crtShader *ebiten.Shader

func init() {
    src, _ := os.ReadFile("shader/crt.kage")
    crtShader, _ = ebiten.NewShader(src)
}

func applyCRTShader(dst, src *ebiten.Image) {
    op := &ebiten.DrawRectShaderOptions{}
    op.Images[0] = src
    dst.DrawRectShader(src.Bounds().Dx(), src.Bounds().Dy(), crtShader, op)
}
```

### 5.3 Waterfall Renderer

The waterfall is a pre-allocated `*ebiten.Image` of fixed size. Each `Update()` call shifts the pixel buffer downward by one row and writes a new top row of random characters. This is done by drawing to an `image.RGBA` and uploading via `ebiten.Image.WritePixels()`.

```go
// internal/render/waterfall.go
type WaterfallRenderer struct {
    img     *ebiten.Image
    buf     *image.RGBA      // CPU-side pixel buffer
    cols    int
    charPx  int              // character cell width/height in pixels
    font    *opentype.Font
    corrupt *float64         // read from GameState each frame
}

var normalSet = []rune("01 ")
var glitchSet = []rune("01#@!?% ")

func (w *WaterfallRenderer) Update() {
    charset := normalSet
    if *w.corrupt > 75 {
        charset = glitchSet
    }
    // shift buffer down one row, draw new top row, upload pixels
    shiftDown(w.buf)
    drawTopRow(w.buf, charset, w.font)
    w.img.WritePixels(w.buf.Pix)
}

func (w *WaterfallRenderer) Draw(screen *ebiten.Image) {
    screen.DrawImage(w.img, nil)
}
```

### 5.4 Hardware Panel — Scrollable Region

Ebitengine has no built-in scroll widgets. The hardware panel uses a **sub-image clip** pattern:

1. All hardware rows are drawn onto a tall off-screen `*ebiten.Image` (`panelCanvas`).
2. A sub-image slice of `panelCanvas` is blitted to the screen using `SubImage(scrollRect)`.
3. Mouse wheel events in the panel region adjust `state.ScrollOffset`.

```go
func (r *HardwarePanelRenderer) Draw(screen *ebiten.Image, state *model.GameState) {
    r.panelCanvas.Clear()
    y := 0
    for _, def := range model.AllHardware {
        drawHardwareRow(r.panelCanvas, def, state, y)
        y += rowHeight
    }
    visibleRect := image.Rect(0, state.ScrollOffset, panelW, state.ScrollOffset+panelH)
    sub := r.panelCanvas.SubImage(visibleRect).(*ebiten.Image)
    op := &ebiten.DrawImageOptions{}
    op.GeoM.Translate(float64(panelX), float64(panelY))
    screen.DrawImage(sub, op)
}
```

### 5.5 Text Rendering

Text is rendered using `github.com/hajimehoshi/ebiten/v2/text/v2` with a pre-loaded `text.GoTextFace` backed by `BPdotsLight.otf`.

```go
var termFace *text.GoTextFace

func init() {
    // Loaded from embedded assets.FS
    f, _ := assets.FS.ReadFile("fonts/BPdotsLight.otf")
    tt, _ := text.NewGoTextFaceSource(bytes.NewReader(f))
    termFace = &text.GoTextFace{Source: tt, Size: 14}
}

func DrawText(dst *ebiten.Image, str string, x, y int, clr color.RGBA) {
    op := &text.DrawOptions{}
    op.GeoM.Translate(float64(x), float64(y))
    op.ColorScale.ScaleWithColor(clr)
    text.Draw(dst, str, termFace, op)
}
```

---

## 6. Input System

```go
// internal/input/input.go
package input

import (
    "image"
    "github.com/hajimehoshi/ebiten/v2"
    "github.com/hajimehoshi/ebiten/v2/inpututil"
)

type InputSystem struct {
    MousePos      image.Point
    Clicked       bool   // left mouse button just pressed
    ScrollDelta   int    // mouse wheel delta this frame
    RebootPressed bool   // keyboard R or REBOOT button clicked
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
    return in.Clicked && in.MousePos.In(ClickerRegion)
}

func (in *InputSystem) HardwareBuyPressed(rowRect image.Rectangle) bool {
    return in.Clicked && in.MousePos.In(rowRect)
}
```

### Screen Region Constants

```go
// internal/ui/layout.go
var (
    LeftPanelRect    = image.Rect(0,   0,   512, 768)
    RightPanelRect   = image.Rect(512, 0,   1280, 768)
    ClickerRegion    = image.Rect(156, 350, 356, 430)
    HardwarePanelRect= image.Rect(532, 220, 1260, 720)
    RebootBtnRect    = image.Rect(532, 730, 800, 760)
    MetricsHUDRect   = image.Rect(532, 20,  1260, 210)
)
```

---

## 7. Glitch & Corruption Visual System

### 7.1 Thresholds

| Corruption % | Visual Effect |
| --- | --- |
| 0 – 49 | Clean render; standard green palette |
| 50 – 74 | Subtle CRT flicker (shader `Time` uniform pulsed) |
| 75 – 89 | Waterfall charset → glitch set; random pixel noise bands in HUD |
| 90 – 99 | Horizontal tear scanlines; HUD text character substitution |
| 100 | CRITICAL: forced SYSTEM_REBOOT dialog; full-screen red flash |

### 7.2 GlitchSystem

```go
// internal/engine/glitch.go
type GlitchSystem struct {
    TearLines  []TearLine   // active horizontal tears
    NoiseAlpha uint8        // 0–200, controls noise band opacity
    Tick       int
}

type TearLine struct {
    Y      int
    Width  int
    OffsetX int
    Life   int // frames remaining
}

func (g *GlitchSystem) Step(corruption float64) {
    g.Tick++
    g.NoiseAlpha = uint8((max(0, corruption-50) / 50.0) * 200)

    // Spawn new tears above 75% corruption
    if corruption > 75 && g.Tick%8 == 0 {
        g.TearLines = append(g.TearLines, TearLine{
            Y: rand.Intn(768), Width: rand.Intn(400) + 100,
            OffsetX: rand.Intn(20) - 10, Life: rand.Intn(6) + 2,
        })
    }
    // Age out old tears
    alive := g.TearLines[:0]
    for _, t := range g.TearLines {
        t.Life--
        if t.Life > 0 { alive = append(alive, t) }
    }
    g.TearLines = alive
}
```

### 7.3 Glitch Overlay Renderer

```go
func (r *GlitchOverlayRenderer) Draw(screen *ebiten.Image, g *GlitchSystem) {
    // Horizontal noise bands
    for _, tear := range g.TearLines {
        px := image.Rect(tear.OffsetX, tear.Y, tear.OffsetX+tear.Width, tear.Y+2)
        ebitenutil.DrawRect(screen, float64(px.Min.X), float64(px.Min.Y),
            float64(px.Dx()), float64(px.Dy()),
            color.RGBA{0, 255, 70, g.NoiseAlpha})
    }
    // Full-screen noise pixel scatter
    if g.NoiseAlpha > 0 {
        for i := 0; i < int(g.NoiseAlpha/4); i++ {
            x, y := rand.Intn(1280), rand.Intn(768)
            screen.Set(x, y, color.RGBA{0, 255, 0, g.NoiseAlpha / 2})
        }
    }
}
```

---

## 8. Hardware & Economy Systems

### 8.1 Hardware Definition

```go
// internal/model/hardware_def.go
type HardwareDef struct {
    ID            string
    Name          string
    Tier          int
    BaseCost      float64
    CostScaling   float64  // per-purchase multiplier (default 1.15)
    BaseBPS       float64  // bits/second contribution
    EntropyWeight float64  // entropy/second (negative = reduction)
    Description   string
}

var AllHardware = []HardwareDef{
    {ID:"logic_gate",    Name:"Logic Gate",     Tier:1, BaseCost:10,       CostScaling:1.15, BaseBPS:0.1,   EntropyWeight: 0.01},
    {ID:"alu",           Name:"ALU",             Tier:1, BaseCost:100,      CostScaling:1.15, BaseBPS:0.8,   EntropyWeight: 0.05},
    {ID:"ecc_memory",    Name:"ECC Memory",      Tier:1, BaseCost:80,       CostScaling:1.12, BaseBPS:0.0,   EntropyWeight:-0.08},
    {ID:"heatsink",      Name:"Heatsink & Fan",  Tier:1, BaseCost:120,      CostScaling:1.12, BaseBPS:0.0,   EntropyWeight:-0.12},
    {ID:"quantum_core",  Name:"Quantum Core",    Tier:2, BaseCost:5000,     CostScaling:1.18, BaseBPS:15.0,  EntropyWeight: 0.80},
    {ID:"neural_link",   Name:"Neural Link",     Tier:2, BaseCost:12000,    CostScaling:1.18, BaseBPS:35.0,  EntropyWeight: 1.50},
    {ID:"ai_kernel",     Name:"AI Kernel",       Tier:3, BaseCost:500000,   CostScaling:1.20, BaseBPS:250.0, EntropyWeight: 3.00},
    {ID:"temp_buffer",   Name:"Temporal Buffer", Tier:3, BaseCost:1200000,  CostScaling:1.20, BaseBPS:120.0, EntropyWeight:-2.50},
}
```

### 8.2 Cost Calculation

```go
func CurrentCost(def HardwareDef, owned int) float64 {
    return def.BaseCost * math.Pow(def.CostScaling, float64(owned))
}
```

### 8.3 Manual Click Value

```go
func (ge *GameEngine) manualClickValue() float64 {
    baseClick := 1.0
    // Upgrades can boost click value
    for id, owned := range ge.state.Upgrades {
        if owned { baseClick *= upgradeClickMult(id) }
    }
    return baseClick * ge.state.GHzMultiplier
}
```

### 8.4 Bit Formatter

```go
// internal/format/bits.go
var units = []string{"Bits","Kilobits","Megabits","Gigabits","Terabits",
                     "Petabits","Exabits","Zettabits","Yottabits","Brontobytes"}

func FormatBits(b float64) string {
    idx := 0
    for b >= 1000 && idx < len(units)-1 {
        b /= 1000; idx++
    }
    return fmt.Sprintf("%.2f %s", b, units[idx])
}
```

---

## 9. Persistence & Autosave

### 9.1 Snapshot-Based Autosave

Because `GameState` has no mutex (relying on Ebitengine's single-goroutine update), autosave uses a **deep-copied snapshot** passed to the save goroutine via channel.

```go
// internal/persist/autosave.go
type Autosaver struct {
    snapshotCh chan model.GameState
    path       string
}

func (a *Autosaver) Start() {
    go func() {
        for snap := range a.snapshotCh {
            _ = Save(snap, a.path)
        }
    }()
}

// Called from Update() on the 30-second boundary
func (a *Autosaver) MaybeSnapshot(state *model.GameState, now time.Time) {
    select {
    case a.snapshotCh <- *state: // struct copy; maps need deep clone
    default: // previous save still in progress — skip
    }
}
```

> **Map deep clone:** `Hardware` and `Upgrades` maps are cloned via explicit `for k,v := range` copy before passing the snapshot.

### 9.2 Save Schema

```json
{
  "version": "1.0.0",
  "saved_at": "2025-08-14T22:31:00Z",
  "bits": 142857.0,
  "total_bits_earned": 9999999.0,
  "entropy": 38.2,
  "corruption": 12.5,
  "ghz_multiplier": 2.0,
  "reboot_count": 1,
  "hardware": { "logic_gate": 45, "alu": 12 },
  "upgrades": { "gate_efficiency_1": true }
}
```

---

## 10. Prestige System (SYSTEM_REBOOT)

### 10.1 Eligibility

`TotalBitsEarned >= 1,000,000` (configurable). Triggered by pressing `R` or clicking the REBOOT button. A full-screen terminal-style confirmation dialog rendered directly in `Draw()` intercepts all input until confirmed or cancelled.

### 10.2 Dialog Rendering

```go
// Rendered in Layer 7 when state.RebootPending == true
func drawRebootDialog(screen *ebiten.Image, state *model.GameState) {
    // Semi-transparent black fill
    ebitenutil.DrawRect(screen, 240, 234, 800, 300, color.RGBA{0, 0, 0, 210})
    // Border
    // ...DrawLine calls for green border...
    gain := computeGHzGain(state)
    DrawText(screen, ">> SYSTEM_REBOOT INITIATED <<", 340, 260, green)
    DrawText(screen, fmt.Sprintf("GHz GAIN: +%.3f×", gain), 340, 300, brightGreen)
    DrawText(screen, "[Y] CONFIRM    [N] ABORT", 340, 380, green)
}
```

### 10.3 GHz Reward Formula

``` formula
gain          = log10(TotalBitsEarned / 1,000,000) × 0.1
GHzMultiplier += gain
```

### 10.4 State Reset

```go
func (p *PrestigeManager) Reboot(state *model.GameState) {
    gain := math.Log10(state.TotalBitsEarned/1_000_000) * 0.1
    state.GHzMultiplier  += gain
    state.RebootCount++
    state.Bits            = 0
    state.TotalBitsEarned = 0
    state.Entropy         = 0
    state.Corruption      = 0
    state.Hardware        = make(map[string]int)
    state.Upgrades        = make(map[string]bool)
    state.RebootPending   = false
}
```

---

## 11. Audio

Ebitengine uses `github.com/hajimehoshi/ebiten/v2/audio` with WAV decoding. All audio assets are embedded via the `assets` package.

```go
// internal/audio/audio.go
var (
    audioCtx    *audio.Context
    clickPlayer *audio.Player
    loopPlayer  *audio.Player
    alarmPlayer *audio.Player
)

func Init() {
    audioCtx = audio.NewContext(44100)
    // Background loop uses audio.NewInfiniteLoop
    loopPlayer   = loadWAV("sfx/BitDecayLoop.wav", true)
    clickPlayer  = loadWAV("sfx/BitDecayClick.wav", false)
    alarmPlayer  = loadWAV("sfx/BitDecayClick.wav", false) // Placeholder
}

func PlayClick() {
    clickPlayer.Rewind()
    clickPlayer.Play()
}

// Alarm triggered in Update() when Corruption > 90 (throttled to once per 3s)
func PlayAlarm() { alarmPlayer.Rewind(); alarmPlayer.Play() }
```

---

## 12. Configuration & Balancing

```toml
# config/balance.toml

[economy]
autosave_interval_s   = 30
bits_format_threshold = 1000.0

[loop]
game_tick_ms          = 100   # logic tick interval
ebitengine_tps        = 60    # Update() calls per second

[entropy]
corruption_threshold  = 50.0
decay_threshold       = 75.0
max_corruption_penalty= 0.5

[prestige]
reboot_unlock_bits    = 1_000_000
ghz_gain_coefficient  = 0.1

[render]
waterfall_cols        = 24
waterfall_char_px     = 14
glitch_alpha_max      = 200
scanline_strength     = 0.15   # Kage shader uniform
```

---

## 13. Build & Deployment

### 13.1 Dependencies

```go
// go.mod
module github.com/popolque/firstbitengi

go 1.24.13

require (
    github.com/hajimehoshi/ebiten/v2 v2.9.9
)
```

### 13.2 Build Commands

```bash
# Development run
go run ./cmd/bitdecay

# Desktop release builds
GOOS=windows GOARCH=amd64 go build -o dist/bitdecay.exe ./cmd/bitdecay
GOOS=darwin  GOARCH=arm64 go build -o dist/bitdecay-mac  ./cmd/bitdecay
GOOS=linux   GOARCH=amd64 go build -o dist/bitdecay-linux ./cmd/bitdecay

# WebAssembly build (browser)
GOOS=js GOARCH=wasm go build -o dist/bitdecay.wasm ./cmd/bitdecay
cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" dist/
# Serve dist/ with any static file server
```

### 13.3 Shader Embedding

Kage shaders are embedded at compile time via Go's `//go:embed` directive within the rendering package:

```go
// internal/render/renderer.go

//go:embed crt.kage
var crtShaderSrc []byte

func loadCRTShader() (*ebiten.Shader, error) {
    return ebiten.NewShader(crtShaderSrc)
}
```

### 13.4 Minimum System Requirements

| | Minimum |
| --- | --- |
| OS | Windows 10 / macOS 11 / Ubuntu 20.04 |
| RAM | 64 MB |
| GPU | OpenGL 3.1 / Metal / DirectX 11 capable |
| Disk | < 25 MB |
| Browser (WASM) | Chrome 90+, Firefox 88+, Safari 15+ |

---

## 14. Risk Register

| Risk | Likelihood | Impact | Mitigation |
| --- | --- | --- | --- |
| `WritePixels()` waterfall causes GC pressure | Medium | Frame drops | Pre-allocate `[]byte` pixel slice; avoid per-frame allocs |
| Kage shader compile failure on older drivers | Low | No CRT effect | Fallback renderer skips shader pass; logs warning |
| Map snapshot race during autosave | Medium | Corrupt save | Deep-clone maps before passing snapshot to goroutine |
| WASM build: audio context blocked by browser autoplay policy | High | No sound | Defer audio init to first user click event |
| 60 TPS + heavy glitch rendering causes CPU spike | Medium | Battery drain | Cap noise scatter pixels; profile with `pprof` |
| Ebitengine sub-tick accumulator drift over long sessions | Low | Economy imbalance | Clamp `accumMs` to max 3× GameTickMs per Update |

---

***End of Technical Design Document — BIT-DECAY // Go + Ebitengine***
