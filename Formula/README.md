# Sona Homebrew Formula

This directory contains the Homebrew formula for installing Sona via Homebrew.

## Installation

Users can install Sona using this formula by running:

```bash
# Add this tap to Homebrew
brew tap harsh-2002/sona

# Install Sona
brew install sona
```

## Updating

To update Sona to the latest version:

```bash
brew upgrade sona
```

## Formula Details

- **Class**: `Sona`
- **Description**: Audio transcription tool using AssemblyAI
- **Homepage**: https://github.com/Harsh-2002/Sona
- **License**: MIT
- **Dependencies**: Go (for building from source)

## Building from Source

The formula builds Sona from source using Go, which ensures compatibility and allows for customization.

## Notes

- This formula requires Go to be installed
- The binary is built with optimizations (`-ldflags="-s -w"`)
- Test command verifies the installation by checking the version
