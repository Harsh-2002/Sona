#!/bin/bash

# Build script for Sona CLI tool

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}🚀 Building Sona CLI tool...${NC}"

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo -e "${RED}❌ Go is not installed. Please install Go 1.21 or later.${NC}"
    exit 1
fi

# Get Go version
GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
echo -e "${YELLOW}📦 Go version: $GO_VERSION${NC}"

# Clean build directory
echo -e "${YELLOW}🧹 Cleaning build directory...${NC}"
rm -rf build/
mkdir -p build/

# Install dependencies
echo -e "${YELLOW}📥 Installing dependencies...${NC}"
go mod tidy

# Build for current platform
echo -e "${YELLOW}🔨 Building for current platform...${NC}"
go build -ldflags="-s -w" -o build/sona cmd/sona/main.go

# Check if build was successful
if [ -f "build/sona" ]; then
    echo -e "${GREEN}✅ Build successful!${NC}"
    echo -e "${GREEN}📁 Binary location: build/sona${NC}"
    
    # Show file size
    SIZE=$(du -h build/sona | cut -f1)
    echo -e "${GREEN}📏 Binary size: $SIZE${NC}"
    
    # Test the binary
    echo -e "${YELLOW}🧪 Testing binary...${NC}"
    if ./build/sona --help &> /dev/null; then
        echo -e "${GREEN}✅ Binary test passed!${NC}"
    else
        echo -e "${RED}❌ Binary test failed!${NC}"
        exit 1
    fi
else
    echo -e "${RED}❌ Build failed!${NC}"
    exit 1
fi

echo -e "${GREEN}🎉 Build completed successfully!${NC}"
echo -e "${YELLOW}💡 Run './build/sona --help' to see available commands${NC}"
