# BIT-DECAY // ANNOTATED CODEBASE

This document explains the technical rationale behind the implementation of BIT-DECAY. It is intended for developers who want to understand the *why* behind specific patterns and architectural choices.

---

## 1. Entry Point: `cmd/bitdecay/main.go`

The entry point orchestrates the bootstrap process. The order of operations here is critical for hardware and OS compatibility.

* **`audio.Init()`**: Called first because Ebitengine's audio context initialization can be sensitive to the thread it's called on. It must happen before assets are loaded.
* **`audio.LoadSounds()`**: Loads the embedded `.wav` files into memory. We do this at startup to avoid "stuttering" during gameplay that would occur if we loaded sounds on-demand.
* **`audio.StartLoop()`**: Initiates the background ambience immediately.
* **`engine.LoadAssets()`**: Specifically loads fonts. Since our UI is text-heavy, failing here is a hard error.
* **`ebiten.SetWindowResizingMode(...)`**: We use a fixed logical resolution (1280x768), but we allow window resizing. Ebitengine automatically handles the scaling and letterboxing, keeping our rendering logic simple.

---

## 2. The Engine Loop: `internal/engine/gameengine.go`

This is the most complex logic in the game, managing the "Two-Speed Loop".

### 2.1 Sub-Tick Accumulator

Ebitengine runs at 60 Frames Per Second (FPS). However, incremental games often feel "jittery" if logic runs that fast, and it makes balancing harder.

```go
const GameTickMs = 100.0 // 10 ticks per second
```

We decouple game logic (bits per second, entropy) from the frame rate.

* **`accumMs += UpdateMs`**: We track how much time has passed since the last frame.
* **`for ge.accumMs >= GameTickMs`**: This "catch-up" loop ensures that if the user's computer lags and skips a frame, the game logic "simulates" the missing time in the next frame. This prevents the economy from slowing down due to low FPS.

### 2.2 Input Handling

* **`handleInputs`**: Inputs are processed every frame (60 FPS), even if a logic tick hasn't happened. This ensures the "Manual Clicker" feels responsive and snappy.
* **`RebootTriggered`**: We check for both a keyboard 'R' and a mouse click on the reboot button. This "dual-input" approach is standard for accessible game design.

---

## 3. Data Model & State: `internal/model/state.go`

### 3.1 Floating Point Precision

* **`float64`**: We use `float64` for all currency (Bits). In incremental games, values quickly exceed the limits of `int32` (2 billion). `float64` provides enough precision for astronomical numbers (up to 1.8e308) while remaining performant.

### 3.2 State Sanitization

* **`Sanitize()`**: JSON unmarshaling in Go creates `nil` maps if they are empty in the file. `Sanitize()` ensures that `Hardware` and `Upgrades` maps are always initialized, preventing "null pointer" panics when the game tries to read or write to them.

---

## 4. Rendering Pipeline: `internal/render/renderer.go`

BIT-DECAY uses a **Layered Immediate-Mode** approach.

### 4.1 Off-Screen Buffer

* **`offscreen := ebiten.NewImage(...)`**: We draw the entire game to a separate "canvas" first. This allows us to apply the **CRT Kage Shader** to the *entire result* in one single GPU pass, which is much more efficient than applying it per-widget.

### 4.2 Waterfall Optimization

* **`MatrixColumn`**: Each column has its own pre-allocated symbol slice.
* **`Update` Logic**: We avoid using `append()` or `make()` inside the `Update` loop. Instead, we overwrite existing indices in the symbol slice. This prevents the "Garbage Collector" from causing micro-stutters in the animation.

---

## 5. Audio System: `internal/audio/audio.go`

### 5.1 Embedded Filesystem

* **`go:embed`**: We embed assets directly into the `.exe`. This means the user only needs one file to play the game, and we never have to worry about "File Not Found" errors due to the user running the game from the wrong folder.

### 5.2 Infinite Looping

* **`audio.NewInfiniteLoop`**: This is a specialized Ebitengine feature. Instead of checking `if player.IsPlaying()` and restarting, `NewInfiniteLoop` tells the audio hardware to loop the data automatically at the buffer level, resulting in a perfectly gapless "Steady Loop."

---

## 6. Persistence: `internal/persist/autosave.go`

### 6.1 Thread-Safe Snapshots

Saving a large file to disk takes time. If we did it on the main thread, the game would "freeze" for a few milliseconds every 30 seconds.

* **`go func()`**: We run the actual `os.WriteFile` in a background thread (Goroutine).
* **Deep Cloning**: Before sending the state to the background thread, we copy every map manually. This is vital because if the main thread modified the `Hardware` map *while* the background thread was saving it, the program would crash with a "Concurrent Map Write" error.
