#!/bin/bash

# Version management script for Sona CLI

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Get current version from git tags
get_current_version() {
    local latest_tag=$(git describe --tags --abbrev=0 2>/dev/null || echo "v0.0.0")
    echo "${latest_tag#v}"
}

# Increment version
increment_version() {
    local version=$1
    local increment_type=$2
    
    IFS='.' read -ra VERSION_PARTS <<< "$version"
    local major=${VERSION_PARTS[0]}
    local minor=${VERSION_PARTS[1]}
    local patch=${VERSION_PARTS[2]}
    
    case $increment_type in
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
            echo -e "${RED}Error: Invalid increment type. Use: major, minor, or patch${NC}"
            exit 1
            ;;
    esac
    
    echo "${major}.${minor}.${patch}"
}

# Show current version
show_version() {
    local current_version=$(get_current_version)
    echo -e "${BLUE}Current version: ${GREEN}v${current_version}${NC}"
}

# Show next version options
show_next_versions() {
    local current_version=$(get_current_version)
    echo -e "${BLUE}Current version: ${GREEN}v${current_version}${NC}"
    echo -e "${YELLOW}Next version options:${NC}"
    echo -e "  ${GREEN}patch${NC}: v$(increment_version $current_version patch)   (bug fixes, small changes)"
    echo -e "  ${GREEN}minor${NC}: v$(increment_version $current_version minor)   (new features, backward compatible)"
    echo -e "  ${GREEN}major${NC}: v$(increment_version $current_version major)   (breaking changes, major updates)"
}

# Create new version
create_version() {
    local increment_type=$1
    
    if [ -z "$increment_type" ]; then
        echo -e "${RED}Error: Please specify increment type (major, minor, or patch)${NC}"
        echo -e "Usage: $0 create [major|minor|patch]"
        exit 1
    fi
    
    local current_version=$(get_current_version)
    local new_version=$(increment_version $current_version $increment_type)
    local new_tag="v${new_version}"
    
    echo -e "${BLUE}Current version: ${GREEN}v${current_version}${NC}"
    echo -e "${BLUE}New version: ${GREEN}${new_tag}${NC}"
    
    # Check if working directory is clean
    if [ -n "$(git status --porcelain)" ]; then
        echo -e "${RED}Error: Working directory is not clean. Please commit or stash changes first.${NC}"
        git status --short
        exit 1
    fi
    
    # Check if we're on main branch
    local current_branch=$(git branch --show-current)
    if [ "$current_branch" != "main" ]; then
        echo -e "${YELLOW}Warning: You're not on the main branch (current: ${current_branch})${NC}"
        read -p "Continue anyway? (y/N): " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            echo -e "${RED}Version creation cancelled.${NC}"
            exit 1
        fi
    fi
    
    # Create and push tag
    echo -e "${YELLOW}Creating tag ${new_tag}...${NC}"
    git tag -a "$new_tag" -m "Release $new_tag"
    
    echo -e "${YELLOW}Pushing tag to remote...${NC}"
    git push origin "$new_tag"
    
    echo -e "${GREEN}✅ Version ${new_tag} created and pushed successfully!${NC}"
    echo -e "${BLUE}GitHub Actions will automatically build and release this version.${NC}"
}

# Show help
show_help() {
    echo -e "${BLUE}Sona CLI Version Management Script${NC}"
    echo
    echo -e "Usage: $0 [command] [options]"
    echo
    echo -e "Commands:"
    echo -e "  ${GREEN}show${NC}                    Show current version"
    echo -e "  ${GREEN}next${NC}                    Show next version options"
    echo -e "  ${GREEN}create [major|minor|patch]${NC}  Create new version and trigger release"
    echo -e "  ${GREEN}help${NC}                    Show this help message"
    echo
    echo -e "Examples:"
    echo -e "  $0 show                           # Show current version"
    echo -e "  $0 next                           # Show next version options"
    echo -e "  $0 create patch                   # Create patch version (1.0.0 → 1.0.1)"
    echo -e "  $0 create minor                   # Create minor version (1.0.0 → 1.1.0)"
    echo -e "  $0 create major                   # Create major version (1.0.0 → 2.0.0)"
    echo
    echo -e "Note: Creating a version will automatically trigger the GitHub Actions workflow"
    echo -e "to build and release the new version across all platforms."
}

# Main script logic
case "${1:-help}" in
    show)
        show_version
        ;;
    next)
        show_next_versions
        ;;
    create)
        create_version "$2"
        ;;
    help|--help|-h)
        show_help
        ;;
    *)
        echo -e "${RED}Unknown command: $1${NC}"
        echo
        show_help
        exit 1
        ;;
esac
