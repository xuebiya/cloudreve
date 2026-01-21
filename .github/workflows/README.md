# GitHub Actions Workflows

This directory contains GitHub Actions workflows for automated building and testing.

## Workflows

### 1. Build (`build.yml`)

**Trigger**: 
- Push to `main` branch
- Tags starting with `v*`
- Manual trigger via workflow_dispatch

**Platforms**:
- Windows (amd64, arm64)
- Linux (amd64, arm64, arm)
- macOS (amd64, arm64)

**Outputs**:
- Compiled binaries for all platforms
- SHA256 checksums
- GitHub Release (for tags)

**Usage**:
```bash
# Create a release
git tag v4.0.0-cloud189
git push origin v4.0.0-cloud189

# Or manually trigger from GitHub Actions tab
```

### 2. Test Build (`test-build.yml`)

**Trigger**:
- Push to `main` or `dev` branch
- Pull requests to `main`

**Purpose**:
- Quick build test for Linux amd64
- Verify code compiles successfully
- Run unit tests (if available)

## Creating a Release

### Step 1: Tag the release

```bash
# Create and push a tag
git tag -a v4.0.0-cloud189 -m "Release v4.0.0 with Cloud189 driver"
git push origin v4.0.0-cloud189
```

### Step 2: Wait for build

The GitHub Actions workflow will automatically:
1. Build binaries for all platforms
2. Generate SHA256 checksums
3. Create a GitHub Release
4. Upload all artifacts

### Step 3: Download binaries

Go to the [Releases](https://github.com/xuebiya/cloudreve/releases) page to download the compiled binaries.

## Manual Build Trigger

You can manually trigger a build from the GitHub Actions tab:

1. Go to **Actions** tab
2. Select **Build Cloudreve with Cloud189 Driver**
3. Click **Run workflow**
4. Select branch and click **Run workflow**

## Build Artifacts

Each build produces:
- `cloudreve-{os}-{arch}[.exe]` - Compiled binary
- `cloudreve-{os}-{arch}[.exe].sha256` - SHA256 checksum

### Verify Checksum

**Linux/macOS**:
```bash
sha256sum -c cloudreve-linux-amd64.sha256
```

**Windows (PowerShell)**:
```powershell
$hash = (Get-FileHash cloudreve-windows-amd64.exe -Algorithm SHA256).Hash
$expected = (Get-Content cloudreve-windows-amd64.exe.sha256).Split()[0]
if ($hash -eq $expected) { "OK" } else { "FAILED" }
```

## Troubleshooting

### Build fails with "assets.zip not found"

The workflow automatically creates `assets.zip`. If it fails:
1. Check the "Create assets.zip" step in the workflow log
2. Ensure the workflow has write permissions

### Release not created

Releases are only created for tags starting with `v*`:
- ✅ `v4.0.0-cloud189`
- ✅ `v1.0.0`
- ❌ `release-1.0.0`
- ❌ `4.0.0`

### Artifact upload fails

Check that:
1. The build step completed successfully
2. The binary file exists
3. GitHub Actions has sufficient storage quota

## Local Testing

To test the build locally before pushing:

```bash
# Install dependencies
go mod download

# Create assets
mkdir -p assets/build application/statics
echo '{"name":"cloudreve-frontend","version":"4.0.0"}' > assets/build/version.json
cd assets && zip -r assets.zip build/ && cd ..
cp assets/assets.zip application/statics/assets.zip

# Build
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "-s -w" -o cloudreve-test ./main.go
```

## Workflow Status

[![Build Status](https://github.com/xuebiya/cloudreve/actions/workflows/build.yml/badge.svg)](https://github.com/xuebiya/cloudreve/actions/workflows/build.yml)
[![Test Build](https://github.com/xuebiya/cloudreve/actions/workflows/test-build.yml/badge.svg)](https://github.com/xuebiya/cloudreve/actions/workflows/test-build.yml)

## Notes

- Builds use Go 1.21
- CGO is disabled for static binaries
- Binaries are stripped (`-s -w`) to reduce size
- All builds are cross-compiled on Ubuntu
