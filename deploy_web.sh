#!/bin/bash

# Configuration
OUT_DIR="docs"
WASM_BIN="bitdecay.wasm"

echo "🚀 Starting BIT-DECAY Web Deployment..."

# 1. Create output directory (GitHub Pages uses 'docs' or root)
mkdir -p $OUT_DIR

# 2. Build the WebAssembly binary
echo "🏗️  Building WASM binary..."
GOOS=js GOARCH=wasm go build -o $OUT_DIR/$WASM_BIN ./cmd/bitdecay

# 3. Copy Web Assets
echo "📄 Copying web assets..."
cp web/index.html $OUT_DIR/
cp web/wasm_exec.js $OUT_DIR/

echo "✅ Deployment package ready in /$OUT_DIR"
echo "👉 To host on GitHub Pages, push the '$OUT_DIR' folder and set 'Settings > Pages > Source' to '/docs'."
