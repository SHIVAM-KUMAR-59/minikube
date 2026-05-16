#!/bin/bash

set -e

REPO="SHIVAM-KUMAR-59/minikube"
VERSION="v1.0.0"
INSTALL_DIR="/usr/local/bin"
BASE_URL="https://github.com/$REPO/releases/download/$VERSION"

# Colors
RESET="\033[0m"
BOLD="\033[1m"
CYAN="\033[36m"
GREEN="\033[32m"
RED="\033[31m"
GRAY="\033[90m"
YELLOW="\033[33m"

step()    { echo -e "  ${BOLD}${CYAN}→${RESET} $1"; }
success() { echo -e "  ${GREEN}✓${RESET} $1"; }
fail()    { echo -e "  ${RED}✗${RESET} $1"; exit 1; }
info()    { echo -e "  ${GRAY}  $1${RESET}"; }

# ─── Detect platform ────────────────────────────────────────

OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

if [ "$ARCH" = "x86_64" ]; then
    ARCH="amd64"
elif [ "$ARCH" = "arm64" ] || [ "$ARCH" = "aarch64" ]; then
    ARCH="arm64"
else
    fail "Unsupported architecture: $ARCH"
fi

if [ "$OS" = "darwin" ]; then
    PLATFORM="darwin"
elif [ "$OS" = "linux" ]; then
    PLATFORM="linux"
else
    fail "Unsupported OS: $OS"
fi

# ─── Header ─────────────────────────────────────────────────

echo ""
echo -e "  ${BOLD}${CYAN}MiniKube${RESET} ${GRAY}$VERSION${RESET}"
echo -e "  ${GRAY}$PLATFORM/$ARCH${RESET}"
echo ""
echo -e "  ${GRAY}────────────────────────────────────────${RESET}"
echo ""

# ─── Step 1: Binaries ───────────────────────────────────────

echo -e "  ${BOLD}[1/3] Binaries${RESET}"

download() {
    BINARY=$1
    URL="$BASE_URL/$BINARY-$PLATFORM-$ARCH"

    step "Downloading $BINARY"
    info "$URL"

    HTTP_CODE=$(curl -sL --max-time 60 -w "%{http_code}" "$URL" -o "/tmp/$BINARY")

    if [ "$HTTP_CODE" != "200" ]; then
        fail "Failed to download $BINARY (HTTP $HTTP_CODE)"
    fi

    chmod +x "/tmp/$BINARY"
    sudo mv "/tmp/$BINARY" "$INSTALL_DIR/$BINARY"
    success "$BINARY installed"
}

sudo mkdir -p "$INSTALL_DIR"

download minik
download minik-server
download minik-worker

echo ""

# ─── Step 2: Dashboard ──────────────────────────────────────

echo -e "  ${BOLD}[2/3] Dashboard${RESET}"

step "Downloading dashboard.zip"
HTTP_CODE=$(curl -sL --max-time 120 -w "%{http_code}" "$BASE_URL/dashboard.zip" -o /tmp/dashboard.zip)
if [ "$HTTP_CODE" != "200" ]; then
    fail "Failed to download dashboard (HTTP $HTTP_CODE)"
fi
success "Downloaded"

step "Extracting..."
sudo unzip -q -o /tmp/dashboard.zip -d "$INSTALL_DIR"
sudo rm -f /tmp/dashboard.zip
success "Extracted to $INSTALL_DIR/dashboard"

step "Installing dependencies (this may take a moment)..."
cd "$INSTALL_DIR/dashboard" && npm install --silent 2>/dev/null
sudo chown -R "$USER" "$INSTALL_DIR/dashboard"
success "Dependencies installed"

echo ""

# ─── Step 3: Verify ─────────────────────────────────────────

echo -e "  ${BOLD}[3/3] Verifying${RESET}"

verify() {
    BINARY=$1
    if command -v "$BINARY" &>/dev/null; then
        success "$BINARY found at $(which $BINARY)"
    else
        echo -e "  ${YELLOW}⚠${RESET}  $BINARY not found in PATH — you may need to restart your shell"
    fi
}

verify minik
verify minik-server
verify minik-worker

echo ""
echo -e "  ${GRAY}────────────────────────────────────────${RESET}"
echo -e "  ${GREEN}${BOLD}✓ MiniKube installed successfully${RESET}"
echo ""
echo -e "  ${BOLD}Get started:${RESET}"
echo -e "    ${CYAN}minik cluster start --workers 2${RESET}"
echo -e "    ${CYAN}minik dashboard${RESET}"
echo ""