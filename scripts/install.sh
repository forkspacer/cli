#!/usr/bin/env bash

set -e

# Forkspacer CLI Install Script
# This script automatically downloads and installs the latest Forkspacer CLI binary

# Configuration
GITHUB_REPO="forkspacer/cli"
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"
BINARY_NAME="forkspacer"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Helper functions
info() {
    echo -e "${BLUE}ℹ${NC} $1"
}

success() {
    echo -e "${GREEN}✓${NC} $1"
}

error() {
    echo -e "${RED}✗${NC} $1"
}

warning() {
    echo -e "${YELLOW}⚠${NC} $1"
}

# Detect OS
detect_os() {
    case "$(uname -s)" in
        Darwin*)
            OS="darwin"
            ;;
        Linux*)
            OS="linux"
            ;;
        MINGW*|MSYS*|CYGWIN*)
            error "Windows is not supported by this install script"
            error "Please download the Windows binary manually from:"
            error "https://github.com/${GITHUB_REPO}/releases/latest"
            exit 1
            ;;
        *)
            error "Unsupported operating system: $(uname -s)"
            exit 1
            ;;
    esac
}

# Detect architecture
detect_arch() {
    case "$(uname -m)" in
        x86_64|amd64)
            ARCH="amd64"
            ;;
        arm64|aarch64)
            ARCH="arm64"
            ;;
        *)
            error "Unsupported architecture: $(uname -m)"
            exit 1
            ;;
    esac
}

# Check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Check prerequisites
check_prerequisites() {
    if ! command_exists curl && ! command_exists wget; then
        error "Neither curl nor wget found. Please install one of them and try again."
        exit 1
    fi

    if ! command_exists tar; then
        error "tar not found. Please install tar and try again."
        exit 1
    fi
}

# Get latest release version
get_latest_version() {
    if command_exists curl; then
        VERSION=$(curl -sSL "https://api.github.com/repos/${GITHUB_REPO}/releases/latest" | grep '"tag_name"' | cut -d'"' -f4)
    elif command_exists wget; then
        VERSION=$(wget -qO- "https://api.github.com/repos/${GITHUB_REPO}/releases/latest" | grep '"tag_name"' | cut -d'"' -f4)
    fi

    if [ -z "$VERSION" ]; then
        error "Failed to fetch latest version"
        exit 1
    fi
}

# Download and extract binary
download_binary() {
    DOWNLOAD_URL="https://github.com/${GITHUB_REPO}/releases/download/${VERSION}/${BINARY_NAME}-${OS}-${ARCH}.tar.gz"
    TEMP_DIR=$(mktemp -d)
    TEMP_FILE="${TEMP_DIR}/${BINARY_NAME}.tar.gz"

    info "Downloading ${BINARY_NAME} ${VERSION} for ${OS}/${ARCH}..."

    if command_exists curl; then
        curl -sSL -o "${TEMP_FILE}" "${DOWNLOAD_URL}"
    elif command_exists wget; then
        wget -q -O "${TEMP_FILE}" "${DOWNLOAD_URL}"
    fi

    if [ $? -ne 0 ]; then
        error "Failed to download ${BINARY_NAME}"
        error "URL: ${DOWNLOAD_URL}"
        rm -rf "${TEMP_DIR}"
        exit 1
    fi

    success "Downloaded ${BINARY_NAME} ${VERSION}"

    # Extract binary
    info "Extracting binary..."
    tar -xzf "${TEMP_FILE}" -C "${TEMP_DIR}"

    if [ ! -f "${TEMP_DIR}/${BINARY_NAME}" ]; then
        error "Binary not found in archive"
        rm -rf "${TEMP_DIR}"
        exit 1
    fi

    success "Extracted binary"
}

# Install binary
install_binary() {
    info "Installing to ${INSTALL_DIR}..."

    # Check if we need sudo
    if [ -w "${INSTALL_DIR}" ]; then
        SUDO=""
    else
        if ! command_exists sudo; then
            error "${INSTALL_DIR} is not writable and sudo is not available"
            error "Please run this script as root or choose a different installation directory"
            error "Example: INSTALL_DIR=\$HOME/bin $0"
            rm -rf "${TEMP_DIR}"
            exit 1
        fi
        SUDO="sudo"
        warning "Installation requires sudo privileges"
    fi

    # Create install directory if it doesn't exist
    $SUDO mkdir -p "${INSTALL_DIR}"

    # Install binary
    $SUDO mv "${TEMP_DIR}/${BINARY_NAME}" "${INSTALL_DIR}/${BINARY_NAME}"
    $SUDO chmod +x "${INSTALL_DIR}/${BINARY_NAME}"

    # Clean up
    rm -rf "${TEMP_DIR}"

    success "Installed ${BINARY_NAME} to ${INSTALL_DIR}/${BINARY_NAME}"
}

# Verify installation
verify_installation() {
    if ! command_exists "${BINARY_NAME}"; then
        warning "${BINARY_NAME} is installed but not in PATH"
        warning "Please add ${INSTALL_DIR} to your PATH:"
        echo ""
        echo "    export PATH=\"${INSTALL_DIR}:\$PATH\""
        echo ""
        warning "Add this line to your shell configuration file (~/.bashrc, ~/.zshrc, etc.)"
        return
    fi

    INSTALLED_VERSION=$(${BINARY_NAME} version | grep "Version:" | awk '{print $2}')
    success "${BINARY_NAME} ${INSTALLED_VERSION} installed successfully!"

    # Show next steps
    echo ""
    info "Next steps:"
    echo "  1. Verify installation:    ${BINARY_NAME} version"
    echo "  2. Get help:               ${BINARY_NAME} --help"
    echo "  3. Enable completion:      ${BINARY_NAME} completion --help"
    echo "  4. List workspaces:        ${BINARY_NAME} workspace list"
    echo ""
    info "Documentation: https://github.com/${GITHUB_REPO}#readme"
}

# Main installation flow
main() {
    echo ""
    echo "Forkspacer CLI Installer"
    echo "========================"
    echo ""

    check_prerequisites
    detect_os
    detect_arch
    get_latest_version

    info "Detected system: ${OS}/${ARCH}"
    info "Latest version: ${VERSION}"
    echo ""

    download_binary
    install_binary
    verify_installation

    echo ""
    success "Installation complete!"
}

# Run main function
main
