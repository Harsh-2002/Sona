# GitHub Workflow Setup

This workflow automatically builds Sona for all platforms and uploads binaries to your MinIO S3 bucket on every push.

## Required Secrets

Set these secrets in your GitHub repository settings:

1. Go to **Settings** → **Secrets and variables** → **Actions**
2. Add the following secrets:

### `MINIO_ENDPOINT`
Your MinIO endpoint URL
```
https://s3.srvr.site
```

### `MINIO_ACCESS_KEY`
Your MinIO access key

### `MINIO_SECRET_KEY`
Your MinIO secret key

## How It Works

1. **Trigger**: Push to main/master branch
2. **Build**: Creates binaries for 6 platforms:
   - Linux (AMD64, ARM64)
   - macOS (Intel, Apple Silicon)
   - Windows (AMD64, ARM64)
3. **Upload**: Saves new binaries to MinIO bucket `artifact/sona/` (overwrites old ones)

## Setup Required

**Important**: Create the `sona` folder manually in your MinIO bucket:
1. Go to your MinIO console
2. Navigate to the `artifact` bucket
3. Create a folder named `sona`
4. The workflow will automatically upload binaries there

## Binary Names

The workflow creates standardized names:
- `sona-linux-amd64` - Linux Intel/AMD
- `sona-linux-arm64` - Linux ARM64
- `sona-darwin-amd64` - macOS Intel
- `sona-darwin-arm64` - macOS Apple Silicon
- `sona-windows-amd64.exe` - Windows Intel/AMD
- `sona-windows-arm64.exe` - Windows ARM64

## What Happens on Push

Every time you push to main/master:
1. All 6 platform binaries are built
2. New binaries are uploaded to MinIO (overwrites old ones)
3. No GitHub releases - just MinIO storage

## No Manual Triggers

This workflow runs automatically on every push. No need to:
- Create git tags
- Run manual workflows
- Manage versions
- Handle releases

Just push your code and the binaries will be built and uploaded automatically!
