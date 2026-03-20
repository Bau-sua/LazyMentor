#!/bin/bash
set -e

REPO="Bau-sua/LazyMentor"
BINARY="lazymint"

# Detect OS and architecture
OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"

# Map architecture names
case "$ARCH" in
    x86_64) ARCH="amd64" ;;
    aarch64|arm64) ARCH="arm64" ;;
    *) echo "Unsupported architecture: $ARCH"; exit 1 ;;
esac

# Map OS names
case "$OS" in
    linux) PLATFORM="linux" ;;
    darwin) PLATFORM="darwin" ;;
    mingw*|cygwin*|msys*) PLATFORM="windows" ;;
    *) echo "Unsupported OS: $OS"; exit 1 ;;
esac

# Construct release asset name
ASSET="${BINARY}-${PLATFORM}-${ARCH}"
if [ "$PLATFORM" = "windows" ]; then
    ASSET="${ASSET}.exe"
fi

# Create temp directory
TMP_DIR=$(mktemp -d)
cd "$TMP_DIR"

echo "Downloading LazyMentor..."
curl -fsSL "https://github.com/${REPO}/releases/latest/download/${ASSET}" -o "${BINARY}"

chmod +x "${BINARY}"

# Install to ~/.local/bin or /usr/local/bin
if [ -w "$HOME/.local/bin" ]; then
    INSTALL_DIR="$HOME/.local/bin"
elif [ -w "/usr/local/bin" ]; then
    INSTALL_DIR="/usr/local/bin"
else
    INSTALL_DIR="$HOME/.local/bin"
    mkdir -p "$INSTALL_DIR"
fi

mv "${BINARY}" "${INSTALL_DIR}/${BINARY}"

echo ""
echo "✓ LazyMentor installed to: ${INSTALL_DIR}/${BINARY}"
echo ""
echo "Run: ${BINARY}"
