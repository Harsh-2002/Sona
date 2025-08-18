#!/bin/bash

# Cross-platform build script for Sona CLI tool

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}🚀 Building Sona CLI tool for multiple platforms...${NC}"

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo -e "${RED}❌ Go is not installed. Please install Go 1.21 or later.${NC}"
    exit 1
fi

# Clean build directory
echo -e "${YELLOW}🧹 Cleaning build directory...${NC}"
rm -rf build/
mkdir -p build/

# Install dependencies
echo -e "${YELLOW}📥 Installing dependencies...${NC}"
go mod tidy

# Build for multiple platforms
echo -e "${YELLOW}🔨 Building for multiple platforms...${NC}"

# Current platform
echo -e "${YELLOW}📱 Building for current platform...${NC}"
go build -ldflags="-s -w" -o build/sona cmd/sona/main.go

# Linux
echo -e "${YELLOW}🐧 Building for Linux...${NC}"
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o build/sona-linux-amd64 cmd/sona/main.go

# macOS (Intel)
echo -e "${YELLOW}🍎 Building for macOS (Intel)...${NC}"
GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o build/sona-darwin-amd64 cmd/sona/main.go

# macOS (Apple Silicon)
echo -e "${YELLOW}🍎 Building for macOS (Apple Silicon)...${NC}"
GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o build/sona-darwin-arm64 cmd/sona/main.go

# Windows
echo -e "${YELLOW}🪟 Building for Windows...${NC}"
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o build/sona-windows-amd64.exe cmd/sona/main.go

# Show build results
echo -e "${GREEN}✅ Build completed for all platforms!${NC}"
echo -e "${YELLOW}📁 Build outputs:${NC}"
ls -la build/

# Create checksums
echo -e "${YELLOW}🔍 Creating checksums...${NC}"
cd build
for file in *; do
    if [ -f "$file" ]; then
        sha256sum "$file" > "$file.sha256"
        echo "Created checksum for $file"
    fi
done
cd ..

echo -e "${GREEN}🎉 All builds completed successfully!${NC}"
echo -e "${YELLOW}💡 Binaries are ready in the build/ directory${NC}"
