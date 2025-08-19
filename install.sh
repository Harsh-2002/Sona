#!/bin/sh

# Sona Installer Script
# Supports install, upgrade, reinstall, and uninstall operations
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

# Function to get binary name for platform
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
    # Try to get latest version from GitHub releases
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
    
    # Create install directory if it doesn't exist
    mkdir -p "$INSTALL_DIR"
    mkdir -p "$(dirname "$VERSION_FILE")"
    
    # Download binary to temp file first
    if command -v curl >/dev/null 2>&1; then
        curl -L -o "$temp_file" "$download_url"
    elif command -v wget >/dev/null 2>&1; then
        wget -O "$temp_file" "$download_url"
    else
        print_status "$RED" "❌ Error: Neither curl nor wget found. Please install one of them."
        exit 1
    fi
    
    # Verify download
    if [ ! -f "$temp_file" ] || [ ! -s "$temp_file" ]; then
        print_status "$RED" "❌ Error: Download failed or file is empty"
        exit 1
    fi
    
    # Move to final location
    mv "$temp_file" "$INSTALL_DIR/$BINARY_NAME"
    
    # Make binary executable
    chmod +x "$INSTALL_DIR/$BINARY_NAME"
    
    # Save version info
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
    
    print_status "$BLUE" "📦 Installing $binary_name for $platform"
    
    # Download and install binary
    download_binary "$binary_name"
    
    print_status "$GREEN" "🎉 Installation completed!"
    print_status "$GREEN" "✅ Sona is now available system-wide as '$BINARY_NAME'"
    print_status "$BLUE" "📋 Test it with: $BINARY_NAME --help"
}

# Function to upgrade sona
upgrade_sona() {
    local current_version=$(get_installed_version)
    local latest_version=$(get_latest_version)
    
    print_status "$BLUE" "🔄 Checking for updates..."
    print_status "$BLUE" "Current version: $current_version"
    print_status "$BLUE" "Latest version: $latest_version"
    
    if [ "$current_version" = "$latest_version" ] && [ "$latest_version" != "latest" ]; then
        print_status "$GREEN" "✅ Sona is already up to date!"
        return 0
    fi
    
    print_status "$YELLOW" "🔄 Upgrading Sona..."
    
    # Remove old version and reinstall
    rm -f "$INSTALL_DIR/$BINARY_NAME"
    rm -f "$VERSION_FILE"
    
    # Reinstall
    install_sona
    
    print_status "$GREEN" "🎉 Upgrade completed!"
}

# Function to reinstall sona
reinstall_sona() {
    print_status "$YELLOW" "🔄 Reinstalling Sona..."
    
    # Remove old version
    rm -f "$INSTALL_DIR/$BINARY_NAME"
    rm -f "$VERSION_FILE"
    
    # Reinstall
    install_sona
    
    print_status "$GREEN" "🎉 Reinstallation completed!"
}

# Function to uninstall sona
uninstall_sona() {
    print_status "$YELLOW" "🗑️  Uninstalling Sona..."
    
    # Remove binary
    if [ -f "$INSTALL_DIR/$BINARY_NAME" ]; then
        rm -f "$INSTALL_DIR/$BINARY_NAME"
        print_status "$GREEN" "✅ Removed binary from $INSTALL_DIR/"
    else
        print_status "$YELLOW" "⚠️  Binary not found in $INSTALL_DIR/"
    fi
    
    # Remove version file
    if [ -f "$VERSION_FILE" ]; then
        rm -f "$VERSION_FILE"
        print_status "$GREEN" "✅ Removed version file"
    fi
    
    # Remove config directory if empty
    if [ -d "$HOME/.sona" ] && [ -z "$(ls -A "$HOME/.sona" 2>/dev/null)" ]; then
        rmdir "$HOME/.sona"
        print_status "$GREEN" "✅ Removed empty config directory"
    fi
    
    print_status "$GREEN" "🎉 Uninstallation completed!"
}

# Function to show current status
show_status() {
    print_status "$BLUE" "📊 Sona Installation Status"
    print_status "$BLUE" "=========================="
    
    if is_installed; then
        print_status "$GREEN" "✅ Sona is installed"
        print_status "$BLUE" "Location: $(which $BINARY_NAME)"
        print_status "$BLUE" "Version: $(get_installed_version)"
        
        # Test if binary works
        if "$BINARY_NAME" --version >/dev/null 2>&1; then
            print_status "$GREEN" "✅ Binary is working correctly"
        else
            print_status "$RED" "❌ Binary has issues"
        fi
    else
        print_status "$RED" "❌ Sona is not installed"
    fi
    
    # Check dependencies
    print_status "$BLUE" "\n🔧 System Dependencies:"
    if command -v curl >/dev/null 2>&1; then
        print_status "$GREEN" "✅ curl: Available"
    else
        print_status "$RED" "❌ curl: Not found"
    fi
    
    if command -v wget >/dev/null 2>&1; then
        print_status "$GREEN" "✅ wget: Available"
    else
        print_status "$YELLOW" "⚠️  wget: Available (backup)"
    fi
}

# Function to show help
show_help() {
    cat << EOF
🚀 Sona Installer Script

Usage: $0 [OPTIONS]

OPTIONS:
    -i, --install       Install Sona (default if no option specified)
    -u, --upgrade       Upgrade existing installation
    -r, --reinstall     Reinstall Sona (removes old version first)
    -d, --uninstall     Uninstall Sona
    -s, --status        Show installation status
    -h, --help          Show this help message

EXAMPLES:
    $0                    # Install Sona
    $0 --install         # Install Sona
    $0 --upgrade         # Upgrade existing installation
    $0 --reinstall       # Reinstall Sona
    $0 --uninstall       # Remove Sona
    $0 --status          # Check installation status

NOTES:
    - Requires root privileges for system-wide installation
    - Automatically detects your platform (Linux, macOS, Windows)
    - Supports AMD64 and ARM64 architectures
    - Downloads from official Sona releases

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
    
    # Parse command line arguments
    while [ $# -gt 0 ]; do
        case $1 in
            -i|--install)
                action="install"
                shift
                ;;
            -u|--upgrade)
                action="upgrade"
                shift
                ;;
            -r|--reinstall)
                action="reinstall"
                shift
                ;;
            -d|--uninstall)
                action="uninstall"
                shift
                ;;
            -s|--status)
                action="status"
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
        "upgrade")
            check_root
            if is_installed; then
                upgrade_sona
            else
                print_status "$YELLOW" "⚠️  Sona is not installed. Installing instead..."
                install_sona
            fi
            ;;
        "reinstall")
            check_root
            reinstall_sona
            ;;
        "uninstall")
            check_root
            uninstall_sona
            ;;
        "status")
            show_status
            ;;
        *)
            print_status "$RED" "❌ Invalid action: $action"
            exit 1
            ;;
    esac
}

# Run main function
main "$@"
