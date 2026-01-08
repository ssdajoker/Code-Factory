#!/bin/sh
set -e

# Code-Factory Installation Script
# Usage: curl -sSL https://raw.githubusercontent.com/ssdajoker/Code-Factory/main/scripts/install.sh | sh

REPO="ssdajoker/Code-Factory"
INSTALL_DIR="/usr/local/bin"
FALLBACK_DIR="$HOME/.local/bin"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Detect OS and architecture
detect_platform() {
    OS=$(uname -s | tr '[:upper:]' '[:lower:]')
    ARCH=$(uname -m)
    
    case "$OS" in
        linux*)
            OS="linux"
            ;;
        darwin*)
            OS="darwin"
            ;;
        mingw* | msys* | cygwin*)
            OS="windows"
            ;;
        *)
            echo "${RED}Unsupported OS: $OS${NC}"
            exit 1
            ;;
    esac
    
    case "$ARCH" in
        x86_64 | amd64)
            ARCH="amd64"
            ;;
        aarch64 | arm64)
            ARCH="arm64"
            ;;
        *)
            echo "${RED}Unsupported architecture: $ARCH${NC}"
            exit 1
            ;;
    esac
    
    BINARY_NAME="factory-${OS}-${ARCH}"
    if [ "$OS" = "windows" ]; then
        BINARY_NAME="${BINARY_NAME}.exe"
    fi
}

# Print banner
print_banner() {
    echo ""
    echo "${BLUE}╔══════════════════════════════════════════════════════════╗${NC}"
    echo "${BLUE}║                                                          ║${NC}"
    echo "${BLUE}║        🏭  SPEC-DRIVEN SOFTWARE FACTORY  🏭              ║${NC}"
    echo "${BLUE}║                                                          ║${NC}"
    echo "${BLUE}║        Installing Code-Factory...                        ║${NC}"
    echo "${BLUE}║                                                          ║${NC}"
    echo "${BLUE}╚══════════════════════════════════════════════════════════╝${NC}"
    echo ""
}

# Get latest release
get_latest_release() {
    echo "${BLUE}→ Fetching latest release...${NC}"
    LATEST_RELEASE=$(curl -s "https://api.github.com/repos/$REPO/releases/latest" | awk -F'"' '/"tag_name":/ {print $4; exit}')
    
    if [ -z "$LATEST_RELEASE" ]; then
        echo "${RED}✗ Failed to fetch latest release${NC}"
        echo "${YELLOW}Please check your internet connection or install manually from:${NC}"
        echo "  https://github.com/$REPO/releases"
        exit 1
    fi
    
    echo "${GREEN}✓ Latest version: $LATEST_RELEASE${NC}"
}

# Download binary
download_binary() {
    DOWNLOAD_URL="https://github.com/$REPO/releases/download/$LATEST_RELEASE/$BINARY_NAME"
    CHECKSUM_URL="https://github.com/$REPO/releases/download/$LATEST_RELEASE/checksums.txt"
    
    echo "${BLUE}→ Downloading $BINARY_NAME...${NC}"
    
    TMP_DIR=$(mktemp -d)
    TMP_FILE="$TMP_DIR/factory"
    
    if ! curl -sL "$DOWNLOAD_URL" -o "$TMP_FILE"; then
        echo "${RED}✗ Failed to download binary${NC}"
        echo "${YELLOW}URL: $DOWNLOAD_URL${NC}"
        exit 1
    fi
    
    echo "${GREEN}✓ Downloaded successfully${NC}"
    
    # Verify checksum (optional, skip if checksums.txt doesn't exist)
    if curl -sL "$CHECKSUM_URL" -o "$TMP_DIR/checksums.txt" 2>/dev/null; then
        echo "${BLUE}→ Verifying checksum...${NC}"
        cd "$TMP_DIR"
        if command -v sha256sum >/dev/null 2>&1; then
            grep "$BINARY_NAME" checksums.txt | sha256sum -c - || {
                echo "${RED}✗ Checksum verification failed${NC}"
                exit 1
            }
        elif command -v shasum >/dev/null 2>&1; then
            grep "$BINARY_NAME" checksums.txt | shasum -a 256 -c - || {
                echo "${RED}✗ Checksum verification failed${NC}"
                exit 1
            }
        else
            echo "${YELLOW}⚠ sha256sum not found, skipping checksum verification${NC}"
        fi
        echo "${GREEN}✓ Checksum verified${NC}"
    fi
}

# Install binary
install_binary() {
    echo "${BLUE}→ Installing factory...${NC}"
    
    # Try to install to /usr/local/bin first
    if [ -w "$INSTALL_DIR" ] || [ "$(id -u)" -eq 0 ]; then
        TARGET="$INSTALL_DIR/factory"
    else
        # Fallback to ~/.local/bin
        echo "${YELLOW}⚠ No write permission to $INSTALL_DIR, using $FALLBACK_DIR${NC}"
        mkdir -p "$FALLBACK_DIR"
        TARGET="$FALLBACK_DIR/factory"
    fi
    
    mv "$TMP_FILE" "$TARGET"
    chmod +x "$TARGET"
    
    echo "${GREEN}✓ Installed to $TARGET${NC}"
    
    # Check if install directory is in PATH
    case ":$PATH:" in
        *":$INSTALL_DIR:"* | *":$FALLBACK_DIR:"*)
            ;;
        *)
            echo "${YELLOW}⚠ $TARGET is not in your PATH${NC}"
            echo "${YELLOW}  Add this to your shell profile:${NC}"
            echo "    export PATH=\"\$PATH:$(dirname "$TARGET")\""
            ;;
    esac
}

# Verify installation
verify_installation() {
    echo "${BLUE}→ Verifying installation...${NC}"
    
    if command -v factory >/dev/null 2>&1; then
        VERSION=$(factory --version 2>&1 || echo "unknown")
        echo "${GREEN}✓ Factory installed successfully!${NC}"
        echo "${GREEN}  Version: $VERSION${NC}"
    else
        echo "${YELLOW}⚠ Factory installed but not found in PATH${NC}"
        echo "${YELLOW}  You may need to restart your shell or add it to PATH${NC}"
    fi
}

# Print next steps
print_next_steps() {
    echo ""
    echo "${GREEN}╔══════════════════════════════════════════════════════════╗${NC}"
    echo "${GREEN}║                                                          ║${NC}"
    echo "${GREEN}║              ✅  INSTALLATION COMPLETE!  ✅              ║${NC}"
    echo "${GREEN}║                                                          ║${NC}"
    echo "${GREEN}╚══════════════════════════════════════════════════════════╝${NC}"
    echo ""
    echo "${BLUE}Next Steps:${NC}"
    echo ""
    echo "  1. Initialize Factory in your project:"
    echo "     ${YELLOW}cd /path/to/your/project${NC}"
    echo "     ${YELLOW}factory init${NC}"
    echo ""
    echo "  2. Start the TUI:"
    echo "     ${YELLOW}factory${NC}"
    echo ""
    echo "  3. Learn more:"
    echo "     ${YELLOW}factory help${NC}"
    echo ""
    echo "${BLUE}Documentation:${NC} https://github.com/$REPO"
    echo ""
    echo "Happy building! 🏭"
    echo ""
}

# Main installation flow
main() {
    print_banner
    detect_platform
    get_latest_release
    download_binary
    install_binary
    verify_installation
    print_next_steps
}

main
