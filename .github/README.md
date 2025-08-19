# GitHub Workflows

This directory contains automated workflows for Sona.

## 🚀 Release Workflow (`release.yml`)

**What it does:** Automatically creates releases when you create a version tag.

**When it runs:** When you push a tag like `v1.0.0`, `v1.1.0`, etc.

**What it builds:**
- `sona-linux-amd64` - Linux (64-bit)
- `sona-linux-arm64` - Linux (ARM64)
- `sona-darwin-amd64` - macOS (Intel)
- `sona-darwin-arm64` - macOS (Apple Silicon)
- `sona-windows-amd64.exe` - Windows (64-bit)

**What it creates:**
- GitHub release with all binaries
- SHA256 checksums for security
- Custom release notes with installation instructions

## 🧪 Test Workflow (`test.yml`)

**What it does:** Tests that Sona builds correctly.

**When it runs:** On every pull request and push to main branch.

**What it tests:**
- Code compiles without errors
- Binary runs and shows help
- Dependencies are correct

## 📝 How to Use

### Create a Release

#### Option 1: Use the Release Script (Recommended)
```bash
# Create a patch release (1.0.0 -> 1.0.1)
./scripts/release.sh patch

# Create a minor release (1.0.0 -> 1.1.0)
./scripts/release.sh minor

# Create a major release (1.0.0 -> 2.0.0)
./scripts/release.sh major

# With custom message
./scripts/release.sh patch "Bug fixes and improvements"
```

#### Option 2: Manual Git Commands
```bash
# Create a version tag
git tag v1.0.0
git push origin v1.0.0
```

#### What Happens Next:
GitHub Actions automatically:
- Builds Sona for all platforms
- Creates a new release
- Uploads all binaries

### Check Build Status

- Go to **Actions** tab in your GitHub repository
- See which workflows are running
- Check build logs if something fails

## 🔧 Customization

- **Add platforms:** Edit the build commands in `release.yml`
- **Change Go version:** Update `go-version` in both workflows
- **Add tests:** Add more test steps in `test.yml`

## 📋 Requirements

- **GitHub CLI:** The workflow uses `gh` command for uploading checksums
- **Go 1.23+:** For building the binaries
- **Ubuntu Latest:** The workflow runs on Ubuntu for consistent builds
