#!/bin/sh

# Sona Installer Script
# Automatically detects platform and downloads the correct binary
# Compatible with all shells (sh, bash, zsh, dash, etc.)

set -e

# Colors for output (POSIX compatible)
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# MinIO configuration
ENDPOINT="https://s3.srvr.site"
BUCKET="artifact"
FOLDER="sona"

# Function to detect platform
detect_platform() {
    local os
    local arch
    
    # Detect OS
    case "$(uname -s)" in
        Linux*)     os="linux" ;;
        Darwin*)    os="darwin" ;;
        CYGWIN*|MINGW*|MSYS*) os="windows" ;;
        *)          os="unknown" ;;
    esac
    
    # Detect architecture
    case "$(uname -m)" in
        x86_64|amd64) arch="amd64" ;;
        aarch64|arm64) arch="arm64" ;;
        armv7l|armv8l) arch="arm64" ;;
        *)             arch="unknown" ;;
    esac
    
    echo "$os-$arch"
}

# Function to get binary name
get_binary_name() {
    local platform=$1
    
    case $platform in
        linux-amd64)   echo "sona-linux-amd64" ;;
        linux-arm64)   echo "sona-linux-arm64" ;;
        darwin-amd64)  echo "sona-darwin-amd64" ;;
        darwin-arm64)  echo "sona-darwin-arm64" ;;
        windows-amd64) echo "sona-windows-amd64.exe" ;;
        windows-arm64) echo "sona-windows-arm64.exe" ;;
        *)             echo "unknown" ;;
    esac
}

# Function to download binary
download_binary() {
    local binary_name=$1
    local download_url="https://s3.srvr.site/artifact/sona/$binary_name"
    local install_dir="/usr/local/bin"
    
    echo "📥 Downloading $binary_name..."
    echo "URL: $download_url"
    
    # Create install directory if it doesn't exist
    mkdir -p "$install_dir"
    
    # Download binary directly to system bin directory
    if command -v curl >/dev/null 2>&1; then
        curl -L -o "$install_dir/sona" "$download_url"
    elif command -v wget >/dev/null 2>&1; then
        wget -O "$install_dir/sona" "$download_url"
    else
        echo "Error: Neither curl nor wget found. Please install one of them."
        exit 1
    fi
    
    # Make binary executable
    chmod +x "$install_dir/sona"
    
    echo "✅ Downloaded and installed sona to $install_dir/"
}

# Function to check if running as root
check_root() {
    if [ "$(id -u)" -ne 0 ]; then
        echo "❌ This installer requires root privileges for system-wide installation."
        echo "💡 Please run with sudo:"
        echo "   sudo $0"
        exit 1
    fi
}

# Main installation process
main() {
    echo "🚀 Sona Installer"
    echo "================"
    
    # Check if running as root
    check_root
    
    # Detect platform
    local platform=$(detect_platform)
    echo "🔍 Detected platform: $platform"
    
    if [ "$platform" = "unknown-unknown" ]; then
        echo "❌ Unsupported platform: $(uname -s) $(uname -m)"
        echo "Supported platforms:"
        echo "  - Linux (AMD64, ARM64)"
        echo "  - macOS (Intel, Apple Silicon)"
        echo "  - Windows (AMD64, ARM64)"
        exit 1
    fi
    
    # Get binary name
    local binary_name=$(get_binary_name "$platform")
    if [ "$binary_name" = "unknown" ]; then
        echo "❌ Unsupported platform combination: $platform"
        exit 1
    fi
    
    echo "📦 Installing $binary_name for $platform"
    
    # Download and install binary
    download_binary "$binary_name"
    
    echo "🎉 Installation completed!"
    echo "✅ Sona is now available system-wide as 'sona'"
    echo "📋 Test it with: sona --help"
}

# Run main function
main "$@"
