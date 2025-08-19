#!/bin/bash

# Sona Release Script
# Usage: ./scripts/release.sh [patch|minor|major] [message]

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if we're in a git repository
if ! git rev-parse --git-dir > /dev/null 2>&1; then
    print_error "Not in a git repository!"
    exit 1
fi

# Check if we have uncommitted changes
if ! git diff-index --quiet HEAD --; then
    print_warning "You have uncommitted changes. Please commit or stash them first."
    git status --short
    exit 1
fi

# Get current version from git tags
CURRENT_VERSION=$(git describe --tags --abbrev=0 2>/dev/null || echo "v0.0.0")
print_status "Current version: $CURRENT_VERSION"

# Parse version components
VERSION_TYPE=${1:-patch}
RELEASE_MESSAGE=${2:-"Release $VERSION_TYPE"}

# Validate version type
if [[ ! "$VERSION_TYPE" =~ ^(patch|minor|major)$ ]]; then
    print_error "Invalid version type. Use: patch, minor, or major"
    exit 1
fi

# Extract version numbers
CURRENT_MAJOR=$(echo $CURRENT_VERSION | sed 's/v\([0-9]*\)\.\([0-9]*\)\.\([0-9]*\)/\1/')
CURRENT_MINOR=$(echo $CURRENT_VERSION | sed 's/v\([0-9]*\)\.\([0-9]*\)\.\([0-9]*\)/\2/')
CURRENT_PATCH=$(echo $CURRENT_VERSION | sed 's/v\([0-9]*\)\.\([0-9]*\)\.\([0-9]*\)/\3/')

# Calculate new version
case $VERSION_TYPE in
    patch)
        NEW_PATCH=$((CURRENT_PATCH + 1))
        NEW_VERSION="v$CURRENT_MAJOR.$CURRENT_MINOR.$NEW_PATCH"
        ;;
    minor)
        NEW_MINOR=$((CURRENT_MINOR + 1))
        NEW_VERSION="v$CURRENT_MAJOR.$NEW_MINOR.0"
        ;;
    major)
        NEW_MAJOR=$((CURRENT_MAJOR + 1))
        NEW_VERSION="v$NEW_MAJOR.0.0"
        ;;
esac

print_status "New version will be: $NEW_VERSION"

# Confirm release
read -p "Do you want to create release $NEW_VERSION? (y/N): " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    print_warning "Release cancelled."
    exit 0
fi

# Create and push tag
print_status "Creating tag $NEW_VERSION..."
git tag -a "$NEW_VERSION" -m "$RELEASE_MESSAGE"

print_status "Pushing tag to remote..."
git push origin "$NEW_VERSION"

print_success "Release $NEW_VERSION created and pushed!"
print_status "GitHub Actions will now automatically:"
print_status "  - Build Sona for all platforms"
print_status "  - Create a new release"
print_status "  - Upload all binaries"
print_status ""
print_status "Check the Actions tab in your GitHub repository to monitor progress."
print_status "Release will be available at: https://github.com/$(git config --get remote.origin.url | sed 's/.*github.com[:/]\([^/]*\/[^/]*\).*/\1/')/releases/tag/$NEW_VERSION"
