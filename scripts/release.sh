#!/bin/bash

# Simple release script for Sona
# This creates a git tag and pushes it to trigger the GitHub workflow

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to get current version
get_current_version() {
    local latest_tag=$(git describe --tags --abbrev=0 2>/dev/null || echo "v0.0.0")
    echo "${latest_tag#v}"
}

# Function to bump version
bump_version() {
    local current_version=$1
    local bump_type=$2
    
    IFS='.' read -ra VERSION_PARTS <<< "$current_version"
    local major=${VERSION_PARTS[0]:-0}
    local minor=${VERSION_PARTS[1]:-0}
    local patch=${VERSION_PARTS[2]:-0}
    
    case $bump_type in
        major)
            major=$((major + 1))
            minor=0
            patch=0
            ;;
        minor)
            minor=$((minor + 1))
            patch=0
            ;;
        patch)
            patch=$((patch + 1))
            ;;
        *)
            echo -e "${RED}Error: Invalid bump type. Use major, minor, or patch${NC}"
            exit 1
            ;;
    esac
    
    echo "$major.$minor.$patch"
}

# Check if version is provided
if [ -z "$1" ]; then
    echo -e "${RED}Error: Please provide a version or bump type${NC}"
    echo -e "Usage: $0 <version|bump_type>"
    echo -e ""
    echo -e "Examples:"
    echo -e "  $0 1.0.0          # Specific version"
    echo -e "  $0 major           # Bump major version (1.0.0 -> 2.0.0)"
    echo -e "  $0 minor           # Bump minor version (1.0.0 -> 1.1.0)"
    echo -e "  $0 patch           # Bump patch version (1.0.0 -> 1.0.1)"
    echo -e ""
    echo -e "Current version: $(get_current_version)"
    exit 1
fi

# Determine the new version
if [[ "$1" =~ ^(major|minor|patch)$ ]]; then
    CURRENT_VERSION=$(get_current_version)
    VERSION=$(bump_version "$CURRENT_VERSION" "$1")
    echo -e "${BLUE}📊 Current version: $CURRENT_VERSION${NC}"
    echo -e "${BLUE}📈 Bumping $1 version to: $VERSION${NC}"
else
    VERSION=$1
fi

TAG="v$VERSION"

echo -e "${GREEN}🚀 Creating release for Sona $VERSION${NC}"

# Check if working directory is clean
if [ -n "$(git status --porcelain)" ]; then
    echo -e "${RED}Error: Working directory is not clean. Please commit or stash changes first.${NC}"
    git status --short
    exit 1
fi

# Check if we're on main branch
CURRENT_BRANCH=$(git branch --show-current)
if [ "$CURRENT_BRANCH" != "main" ]; then
    echo -e "${YELLOW}Warning: You're not on the main branch (current: $CURRENT_BRANCH)${NC}"
    read -p "Continue anyway? (y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        echo -e "${RED}Release cancelled.${NC}"
        exit 1
    fi
fi

# Check if tag already exists
if git tag -l | grep -q "^$TAG$"; then
    echo -e "${RED}Error: Tag $TAG already exists${NC}"
    exit 1
fi

echo -e "${YELLOW}📝 Creating tag $TAG...${NC}"
git tag -a "$TAG" -m "Release $TAG"

echo -e "${YELLOW}📤 Pushing tag to remote...${NC}"
git push origin "$TAG"

echo -e "${GREEN}✅ Release $TAG created and pushed!${NC}"
echo -e "${YELLOW}📋 The GitHub workflow will now:${NC}"
echo -e "   1. Build Sona for all platforms"
echo -e "   2. Upload binaries to MinIO S3"
echo -e ""
echo -e "${YELLOW}🔍 Check the progress at:${NC}"
echo -e "   https://github.com/$(git config --get remote.origin.url | sed 's/.*github.com[:/]\([^/]*\/[^/]*\).*/\1/')/actions"
