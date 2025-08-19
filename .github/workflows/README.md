# GitHub Workflow Setup

This workflow automatically builds Sona for all platforms and uploads binaries to your MinIO S3 bucket.

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

1. **Trigger**: Push a tag (e.g., `v1.0.0`) or manually run the workflow
2. **Build**: Creates binaries for 6 platforms:
   - Linux (AMD64, ARM64)
   - macOS (Intel, Apple Silicon)
   - Windows (AMD64, ARM64)
3. **Upload**: Saves binaries to MinIO bucket `artifact/sona/`
4. **Release**: Creates GitHub release with all binaries

## Binary Names

The workflow creates standardized names:
- `sona-linux-amd64` - Linux Intel/AMD
- `sona-linux-arm64` - Linux ARM64
- `sona-darwin-amd64` - macOS Intel
- `sona-darwin-arm64` - macOS Apple Silicon
- `sona-windows-amd64.exe` - Windows Intel/AMD
- `sona-windows-arm64.exe` - Windows ARM64

## Manual Run

You can manually trigger the workflow:
1. Go to **Actions** tab
2. Select **Build and Release Sona**
3. Click **Run workflow**
4. Enter version (e.g., `v1.0.0`)
5. Click **Run workflow**
