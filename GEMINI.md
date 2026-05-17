# Project: BIT-DECAY (firstEbitengi)

## Overview
`BIT-DECAY` is an incremental/clicker game featuring a unique "glitch & corruption" visual aesthetic. It is built using **Go 1.24+** and the **Ebitengine (v2.9+)** game library. The project aims for a CRT-inspired look with phosphor bloom, scanlines, and chromatic aberration, implemented via custom Kage shaders.

### Core Technologies
- **Language:** Go
- **Game Engine:** [Ebitengine](https://ebitengine.org/) (Immediate-mode rendering)
- **Shaders:** Kage (Go-like shader language for Ebitengine)
- **Configuration:** TOML (via `github.com/BurntSushi/toml`)
- **Persistence:** JSON

## Architecture
The project follows a modular structure (as defined in the Technical Design Document):
- `cmd/bitdecay/`: Entry point (`main.go`).
- `internal/engine/`: Core game loop integration and logic orchestration.
- `internal/model/`: Game state definitions and static data.
- `internal/render/`: Immediate-mode drawing logic for various UI components.
- `internal/shader/`: CRT and glitch visual effect shaders.
- `internal/persist/`: JSON save/load and autosave logic.

### Game Loop
- **Ebitengine Loop:** Runs at 60 TPS (`Update()` and `Draw()`).
- **Game Logic Tick:** Fixed 10 TPS (100ms) interval, managed via a sub-tick accumulator within the main loop to ensure consistent game progression regardless of frame rate.

## Building and Running

### Development
```bash
# Run the game directly
go run ./cmd/bitdecay
```

### Build Commands
```bash
# Desktop (Windows)
GOOS=windows GOARCH=amd64 go build -o dist/bitdecay.exe ./cmd/bitdecay

# Desktop (macOS)
GOOS=darwin GOARCH=arm64 go build -o dist/bitdecay-mac ./cmd/bitdecay

# WebAssembly
GOOS=js GOARCH=wasm go build -o dist/bitdecay.wasm ./cmd/bitdecay
```

## Development Conventions
- **Resolution:** Logical resolution is fixed at **1280x768**. Ebitengine handles scaling to the window size.
- **Rendering:** All UI is custom-drawn using Ebitengine primitives; no external widget toolkits are used.
- **State Management:** The `GameState` struct is the "canonical truth". Subsystems should read from it and only the `GameEngine` should modify it during `gameTick`.
- **Concurrency:** Use deep-cloning for the game state before passing it to the `persist` package to avoid race conditions during background autosaves.
- **Shaders:** Use `//go:embed` to include `.kage` files in the binary.

## Key Files
- `BitDecay.md`: Comprehensive Technical Design Document (TDD). **Refer to this for detailed system specs.**
- `go.mod`: Project dependencies.
- `config/balance.toml`: Game balancing parameters (planned).
