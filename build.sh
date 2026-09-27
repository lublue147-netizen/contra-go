#!/bin/bash
set -e

echo "=== Building Contra Go WebAssembly ==="

# Create dist directory
mkdir -p dist

# Check if Go is installed, if not try to find or install it
if ! command -v go &> /dev/null; then
    if [ -x "/tmp/gotmp/go/bin/go" ]; then
        export PATH="/tmp/gotmp/go/bin:$PATH"
    elif [ -x "/usr/local/go/bin/go" ]; then
        export PATH="/usr/local/go/bin:$PATH"
    fi
fi

if ! command -v go &> /dev/null; then
    echo "Go is not found in PATH. Attempting automatic download for remote build container..."
    ARCH="amd64"
    if [ "$(uname -m)" = "aarch64" ] || [ "$(uname -m)" = "arm64" ]; then
        ARCH="arm64"
    fi
    GO_VERSION="1.22.5"
    curl -fsSL "https://go.dev/dl/go${GO_VERSION}.linux-${ARCH}.tar.gz" | tar -xz -C /tmp
    export PATH="/tmp/go/bin:$PATH"
fi

echo "Using Go: $(go version)"

# Build WebAssembly binary
echo "Compiling cmd/wasm to dist/contra.wasm..."
GOOS=js GOARCH=wasm go build -ldflags="-s -w" -o dist/contra.wasm ./cmd/wasm

# Copy web assets
echo "Copying web assets to dist/..."
cp web/index.html dist/index.html
cp web/style.css dist/style.css

# Copy wasm_exec.js from Go installation if available, otherwise from web/
if [ -f "$(go env GOROOT)/misc/wasm/wasm_exec.js" ]; then
    cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" dist/wasm_exec.js
elif [ -f "$(go env GOROOT)/lib/wasm/wasm_exec.js" ]; then
    cp "$(go env GOROOT)/lib/wasm/wasm_exec.js" dist/wasm_exec.js
elif [ -f "web/wasm_exec.js" ]; then
    cp web/wasm_exec.js dist/wasm_exec.js
fi

# Copy vercel.json and _headers to dist for static hosting
if [ -f "vercel.json" ]; then
    cp vercel.json dist/vercel.json
fi
if [ -f "web/_headers" ]; then
    cp web/_headers dist/_headers
elif [ -f "_headers" ]; then
    cp _headers dist/_headers
fi

echo "=== Build Complete! ==="
ls -lh dist/
