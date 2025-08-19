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
    local install_dir="$HOME/.local/bin"
    
    echo "📥 Downloading $binary_name..."
    echo "URL: $download_url"
    
    # Create install directory if it doesn't exist
    mkdir -p "$install_dir"
    
    # Download binary
    if command -v curl >/dev/null 2>&1; then
        curl -L -o "$install_dir/$binary_name" "$download_url"
    elif command -v wget >/dev/null 2>&1; then
        wget -O "$install_dir/$binary_name" "$download_url"
    else
        echo "Error: Neither curl nor wget found. Please install one of them."
        exit 1
    fi
    
    # Make binary executable
    chmod +x "$install_dir/$binary_name"
    
    echo "✅ Downloaded $binary_name to $install_dir/"
}

# Function to setup PATH
setup_path() {
    local install_dir="$HOME/.local/bin"
    local shell_rc=""
    
    # Detect shell and config file
    if [ -n "$ZSH_VERSION" ]; then
        shell_rc="$HOME/.zshrc"
    elif [ -n "$BASH_VERSION" ]; then
        shell_rc="$HOME/.bashrc"
    elif [ -n "$KSH_VERSION" ]; then
        shell_rc="$HOME/.kshrc"
    elif [ -n "$FCEDIT" ]; then
        shell_rc="$HOME/.kshrc"
    else
        # Default to .profile for other shells
        shell_rc="$HOME/.profile"
    fi
    
    # Check if PATH already includes install directory
    case ":$PATH:" in
        *":$install_dir:"*) 
            echo "✅ PATH already includes $install_dir"
            return 0
            ;;
    esac
    
    echo "📝 Adding $install_dir to PATH in $shell_rc"
    echo "" >> "$shell_rc"
    echo "# Sona binary path" >> "$shell_rc"
    echo "export PATH=\"\$PATH:$install_dir\"" >> "$shell_rc"
    
    echo "✅ PATH updated! Please restart your terminal or run:"
    echo "   . $shell_rc"
}

# Function to create symlink
create_symlink() {
    local binary_name=$1
    local install_dir="$HOME/.local/bin"
    local symlink_path="/usr/local/bin/sona"
    
    if [ "$(id -u)" -eq 0 ]; then
        echo "🔗 Creating system-wide symlink..."
        ln -sf "$install_dir/$binary_name" "$symlink_path"
        echo "✅ System-wide symlink created at $symlink_path"
        echo "   You can now run 'sona' from anywhere"
    else
        echo "💡 To create a system-wide symlink, run with sudo:"
        echo "   sudo $0"
    fi
}

# Main installation process
main() {
    echo "🚀 Sona Installer"
    echo "================"
    
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
    
    # Download binary
    download_binary "$binary_name"
    
    # Setup PATH
    setup_path
    
    # Create symlink if running as root
    create_symlink "$binary_name"
    
    echo "🎉 Installation completed!"
    echo "📋 Next steps:"
    echo "  1. Restart your terminal or run: . ~/.bashrc (or ~/.zshrc)"
    echo "  2. Test installation: sona --help"
    echo "  3. Start using Sona!"
}

# Run main function
main "$@"
