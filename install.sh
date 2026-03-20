#!/bin/bash
set -e

REPO="Bau-sua/LazyMentor"
BINARY="lazymint"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m' # No Color

# Detect OS and architecture
OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"

# Map architecture names
case "$ARCH" in
    x86_64) ARCH="amd64" ;;
    aarch64|arm64) ARCH="arm64" ;;
    *) echo -e "${RED}Unsupported architecture: $ARCH${NC}" >&2; exit 1 ;;
esac

# Map OS names
case "$OS" in
    linux) PLATFORM="linux" ;;
    darwin) PLATFORM="darwin" ;;
    mingw*|cygwin*|msys*) PLATFORM="windows" ;;
    *) echo -e "${RED}Unsupported OS: $OS${NC}" >&2; exit 1 ;;
esac

# Construct release asset name
ASSET="${BINARY}-${PLATFORM}-${ARCH}"
if [ "$PLATFORM" = "windows" ]; then
    ASSET="${ASSET}.exe"
fi

# Create temp directory for download
TMP_DIR=$(mktemp -d)
cd "$TMP_DIR"

echo -n "Downloading LazyMentor... "
URL="https://github.com/${REPO}/releases/latest/download/${ASSET}"

if ! curl -fsSL "$URL" -o "${BINARY}"; then
    echo -e "${RED}FAILED${NC}"
    echo -e "${RED}Error downloading from $URL${NC}" >&2
    echo -e "${RED}Make sure the release exists at: https://github.com/${REPO}/releases${NC}" >&2
    exit 1
fi
echo -e "${GREEN}OK${NC}"

chmod +x "${BINARY}"

# Determine install directory
if [ -d "$HOME/.local/bin" ] && [ -w "$HOME/.local/bin" ]; then
    INSTALL_DIR="$HOME/.local/bin"
elif [ -w "/usr/local/bin" ]; then
    INSTALL_DIR="/usr/local/bin"
else
    INSTALL_DIR="$HOME/.local/bin"
    mkdir -p "$INSTALL_DIR"
fi

# Move binary to install location
mv "${BINARY}" "${INSTALL_DIR}/${BINARY}"

# Clean up temp directory
rm -rf "$TMP_DIR"

echo -e "${GREEN}✓${NC} Installed to: ${INSTALL_DIR}/${BINARY}"
echo ""

# Check if we have a terminal
if [ -t 0 ] && [ -t 1 ]; then
    # We have a terminal, run the TUI
    echo "Starting LazyMentor..."
    exec "${INSTALL_DIR}/${BINARY}"
else
    # No terminal, run in CLI mode
    echo "Running in CLI mode..."
    "${INSTALL_DIR}/${BINARY}" -install
fi
