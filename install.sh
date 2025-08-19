#!/bin/bash

# Sona Installer Script
# Supports install (default) and uninstall operations
# Compatible with all shells (sh, bash, zsh, dash, etc.)

set -e

# Colors for output (POSIX compatible)
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
ENDPOINT="https://s3.srvr.site"
BUCKET="artifact"
FOLDER="sona"
INSTALL_DIR="/usr/local/bin"
BINARY_NAME="sona"
VERSION_FILE="/usr/local/share/sona/version"

# Function to print colored output
print_status() {
    local color=$1
    local message=$2
    printf "${color}%s${NC}\n" "$message"
}

# Function to detect platform
detect_platform() {
    local os=$(uname -s | tr '[:upper:]' '[:lower:]')
    local arch=$(uname -m | tr '[:upper:]' '[:lower:]')
    
    case $arch in
        x86_64) arch="x86_64" ;;
        amd64) arch="x86_64" ;;
        aarch64) arch="aarch64" ;;
        arm64) arch="aarch64" ;;
        armv7l) arch="armv7l" ;;
        *) arch="unknown" ;;
    esac
    
    case $os in
        linux) os="linux" ;;
        darwin) os="macos" ;;
        msys*|cygwin*|mingw*) os="windows" ;;
        *) os="unknown" ;;
    esac
    
    echo "${os}-${arch}"
}

# Function to get binary name for platform
get_binary_name() {
    local platform=$1
    local os=$(echo $platform | cut -d'-' -f1)
    local arch=$(echo $platform | cut -d'-' -f2)
    
    case $platform in
        linux-x86_64) echo "sona-linux-amd64" ;;
        linux-aarch64) echo "sona-linux-arm64" ;;
        linux-armv7l) echo "sona-linux-armv7" ;;
        macos-x86_64) echo "sona-darwin-amd64" ;;
        macos-aarch64) echo "sona-darwin-arm64" ;;
        windows-x86_64) echo "sona-windows-amd64.exe" ;;
        windows-aarch64) echo "sona-windows-arm64.exe" ;;
        *) echo "unknown" ;;
    esac
}

# Function to get current installed version
get_installed_version() {
    if [ -f "$VERSION_FILE" ]; then
        cat "$VERSION_FILE"
    else
        echo "0.0.0"
    fi
}

# Function to get latest version from GitHub
get_latest_version() {
    if command -v curl >/dev/null 2>&1; then
        curl -s "https://api.github.com/repos/root/sona-ai/releases/latest" | grep '"tag_name"' | cut -d'"' -f4 2>/dev/null || echo "latest"
    elif command -v wget >/dev/null 2>&1; then
        wget -qO- "https://api.github.com/repos/root/sona-ai/releases/latest" | grep '"tag_name"' | cut -d'"' -f4 2>/dev/null || echo "latest"
    else
        echo "latest"
    fi
}

# Function to check if sona is already installed
is_installed() {
    command -v "$BINARY_NAME" >/dev/null 2>&1
}

# Function to download binary
download_binary() {
    local binary_name=$1
    local download_url="https://s3.srvr.site/artifact/sona/$binary_name"
    local temp_file="/tmp/sona-temp"

    print_status "$BLUE" "📥 Downloading $binary_name..."
    print_status "$BLUE" "URL: $download_url"

    mkdir -p "$INSTALL_DIR"
    mkdir -p "$(dirname "$VERSION_FILE")"

    if command -v curl >/dev/null 2>&1; then
        curl -L -o "$temp_file" "$download_url"
    elif command -v wget >/dev/null 2>&1; then
        wget -O "$temp_file" "$download_url"
    else
        print_status "$RED" "❌ Error: Neither curl nor wget found. Please install one of them."
        exit 1
    fi

    if [ ! -f "$temp_file" ] || [ ! -s "$temp_file" ]; then
        print_status "$RED" "❌ Error: Download failed or file is empty"
        exit 1
    fi

    mv "$temp_file" "$INSTALL_DIR/$BINARY_NAME"
    chmod +x "$INSTALL_DIR/$BINARY_NAME"
    echo "$(get_latest_version)" > "$VERSION_FILE"
    print_status "$GREEN" "✅ Downloaded and installed sona to $INSTALL_DIR/"
}

# Function to install sona
install_sona() {
    local platform=$(detect_platform)
    print_status "$BLUE" "🔍 Detected platform: $platform"

    if [ "$platform" = "unknown-unknown" ]; then
        print_status "$RED" "❌ Unsupported platform: $(uname -s) $(uname -m)"
        print_status "$YELLOW" "Supported platforms:"
        print_status "$YELLOW" "  - Linux (AMD64, ARM64)"
        print_status "$YELLOW" "  - macOS (Intel, Apple Silicon)"
        print_status "$YELLOW" "  - Windows (AMD64, ARM64)"
        exit 1
    fi

    local binary_name=$(get_binary_name "$platform")
    if [ "$binary_name" = "unknown" ]; then
        print_status "$RED" "❌ Unsupported platform combination: $platform"
        exit 1
    fi

    if is_installed; then
        local current_version=$(get_installed_version)
        local latest_version=$(get_latest_version)
        
        if [ "$current_version" != "$latest_version" ] && [ "$latest_version" != "latest" ]; then
            print_status "$YELLOW" "🔄 Updating existing installation..."
            print_status "$BLUE" "Current version: $current_version"
            print_status "$BLUE" "Latest version: $latest_version"
            rm -f "$INSTALL_DIR/$BINARY_NAME"
            rm -f "$VERSION_FILE"
        else
            print_status "$GREEN" "✅ Sona is already up to date!"
            return 0
        fi
    fi

    print_status "$BLUE" "📦 Installing $binary_name for $platform"
    download_binary "$binary_name"
    print_status "$GREEN" "🎉 Installation completed!"
    print_status "$GREEN" "✅ Sona is now available system-wide as '$BINARY_NAME'"
    print_status "$BLUE" "📋 Test it with: $BINARY_NAME --help"
}

# Function to uninstall sona
uninstall_sona() {
    print_status "$YELLOW" "🗑️  Uninstalling Sona..."
    
    if [ -f "$INSTALL_DIR/$BINARY_NAME" ]; then
        rm -f "$INSTALL_DIR/$BINARY_NAME"
        print_status "$GREEN" "✅ Removed binary from $INSTALL_DIR/"
    else
        print_status "$YELLOW" "⚠️  Binary not found in $INSTALL_DIR/"
    fi
    
    if [ -f "$VERSION_FILE" ]; then
        rm -f "$VERSION_FILE"
        print_status "$GREEN" "✅ Removed version file"
    fi
    
    if [ -d "$HOME/.sona" ] && [ -z "$(ls -A "$HOME/.sona" 2>/dev/null)" ]; then
        rmdir "$HOME/.sona"
        print_status "$GREEN" "✅ Removed empty config directory"
    fi
    
    # Ask user about removing dependencies
    print_status "$BLUE" "\n🔧 Dependencies (yt-dlp and FFmpeg):"
    print_status "$YELLOW" "These were auto-installed by Sona. Do you want to remove them as well?"
    print_status "$YELLOW" "Note: This will remove the binaries from your system."
    
    read -p "Remove yt-dlp and FFmpeg? (y/N): " -r remove_deps
    
    if [[ $remove_deps =~ ^[Yy]$ ]]; then
        # Remove yt-dlp
        if command -v yt-dlp >/dev/null 2>&1; then
            local ytdlp_path=$(which yt-dlp)
            rm -f "$ytdlp_path"
            print_status "$GREEN" "✅ Removed yt-dlp from $ytdlp_path"
        fi
        
        # Remove FFmpeg
        if command -v ffmpeg >/dev/null 2>&1; then
            local ffmpeg_path=$(which ffmpeg)
            rm -f "$ffmpeg_path"
            print_status "$GREEN" "✅ Removed FFmpeg from $ffmpeg_path"
        fi
        
        print_status "$GREEN" "✅ All dependencies removed"
    else
        print_status "$BLUE" "ℹ️  Dependencies kept (yt-dlp and FFmpeg remain installed)"
    fi
    
    print_status "$GREEN" "🎉 Uninstallation completed!"
}

# Function to show help
show_help() {
    cat << EOF
🚀 Sona Installer Script

Usage: $0 [OPTIONS]

OPTIONS:
    -u, --uninstall    Uninstall Sona
    -h, --help         Show this help message

DEFAULT ACTION:
    If no option is specified, Sona will be installed (or updated if already installed)

EXAMPLES:
    $0                    # Install/Update Sona
    $0 --uninstall       # Uninstall Sona
    $0 --help            # Show this help message

NOTES:
    - Requires root privileges for system-wide installation
    - Automatically detects your platform (Linux, macOS, Windows)
    - Supports AMD64 and ARM64 architectures
    - Downloads from official Sona releases
    - Dependencies (yt-dlp, FFmpeg) are auto-installed when needed

For more information, visit: https://github.com/root/sona-ai
EOF
}

# Function to check if running as root
check_root() {
    if [ "$(id -u)" -ne 0 ]; then
        print_status "$RED" "❌ This installer requires root privileges for system-wide installation."
        print_status "$YELLOW" "💡 Please run with sudo:"
        print_status "$YELLOW" "   sudo $0 [OPTIONS]"
        exit 1
    fi
}

# Main function
main() {
    local action="install"

    while [ $# -gt 0 ]; do
        case $1 in
            -u|--uninstall)
                action="uninstall"
                shift
                ;;
            -h|--help)
                show_help
                exit 0
                ;;
            *)
                print_status "$RED" "❌ Unknown option: $1"
                print_status "$YELLOW" "Use --help for usage information"
                exit 1
                ;;
        esac
    done

    print_status "$BLUE" "🚀 Sona Installer"
    print_status "$BLUE" "================"

    case $action in
        "install")
            check_root
            install_sona
            ;;
        "uninstall")
            check_root
            uninstall_sona
            ;;
        *)
            print_status "$RED" "❌ Invalid action: $action"
            exit 1
            ;;
    esac
}

main "$@"
