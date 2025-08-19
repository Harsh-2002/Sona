#!/bin/bash

# Simple release script for Sona
# This creates a git tag and pushes it to trigger the GitHub workflow

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check if version is provided
if [ -z "$1" ]; then
    echo -e "${RED}Error: Please provide a version number${NC}"
    echo -e "Usage: $0 <version>"
    echo -e "Example: $0 1.0.0"
    exit 1
fi

VERSION=$1
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
