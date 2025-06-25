# Release Process

This document describes the release process for Ship CLI.

## Overview

Ship CLI uses:
- **GoReleaser** for building multi-platform binaries
- **GitHub Actions** for automated releases
- **Conventional Commits** for automatic changelog generation
- **Semantic Versioning** for version numbers

## Prerequisites

1. **GitHub CLI**: Install from https://cli.github.com/
2. **Go 1.21+**: Required for building
3. **Git**: With push access to the repository
4. **Docker**: For testing container builds (optional)

## Release Types

### Patch Release (x.y.Z)
For bug fixes and minor improvements:
```bash
./scripts/create-release.sh
# Select option 1 (Patch)
```

### Minor Release (x.Y.0)
For new features that are backwards compatible:
```bash
./scripts/create-release.sh
# Select option 2 (Minor)
```

### Major Release (X.0.0)
For breaking changes:
```bash
./scripts/create-release.sh
# Select option 3 (Major)
```

## Step-by-Step Release Process

### 1. Prepare for Release

```bash
# Ensure you're on main branch
git checkout main
git pull origin main

# Run tests
go test ./...

# Check GoReleaser config
goreleaser check
```

### 2. Create Release

```bash
# Run the release script
./scripts/create-release.sh

# The script will:
# - Show current version
# - Ask for version bump type
# - Run tests
# - Create and push a git tag
# - Trigger GitHub Actions
```

### 3. Monitor Release

After pushing the tag, GitHub Actions will:
1. Build binaries for all platforms
2. Create Docker images
3. Generate changelog from commits
4. Create GitHub release with artifacts
5. Update Homebrew tap (if configured)

Monitor progress at: https://github.com/cloudshipai/ship/actions

### 4. Verify Release

Once complete, verify:
- GitHub release page: https://github.com/cloudshipai/ship/releases
- Docker Hub: https://hub.docker.com/r/cloudshipai/ship
- Test installation:
  ```bash
  # Linux/macOS
  curl -sSL https://github.com/cloudshipai/ship/releases/latest/download/ship_linux_amd64.tar.gz | tar xz
  ./ship --version
  ```

## Commit Message Convention

Use conventional commits for automatic changelog generation:

- `feat:` New feature
- `fix:` Bug fix
- `docs:` Documentation changes
- `chore:` Maintenance tasks
- `test:` Test additions/changes
- `refactor:` Code refactoring
- `perf:` Performance improvements

Example:
```bash
git commit -m "feat: add support for terraform 1.5"
git commit -m "fix: resolve config loading issue"
git commit -m "docs: update installation instructions"
```

## Manual Release (Emergency)

If automated release fails:

```bash
# Create tag manually
git tag -a v1.2.3 -m "Release v1.2.3"
git push origin v1.2.3

# Build with GoReleaser locally
export GITHUB_TOKEN=your_token
goreleaser release --clean
```

## Post-Release

After a successful release:

1. **Update documentation** if needed
2. **Announce** in relevant channels
3. **Monitor** for any issues
4. **Update** the roadmap

## Troubleshooting

### Build Failures
- Check Go version: `go version` (should be 1.21+)
- Clear GoReleaser cache: `goreleaser release --clean --snapshot`

### GitHub Actions Issues
- Check secrets are configured:
  - `GITHUB_TOKEN` (automatic)
  - `DOCKER_USERNAME` and `DOCKER_PASSWORD`
  - `HOMEBREW_TAP_TOKEN` (if using Homebrew)

### Tag Issues
- Delete local tag: `git tag -d v1.2.3`
- Delete remote tag: `git push origin :refs/tags/v1.2.3`

## Release Artifacts

Each release includes:

- **Binaries**: Linux, macOS, Windows (amd64, arm64)
- **Archives**: `.tar.gz` for Unix, `.zip` for Windows
- **Checksums**: `checksums.txt` with SHA256 hashes
- **Docker Images**: Multi-architecture images
- **Packages**: DEB, RPM, APK formats
- **Changelog**: Auto-generated from commits

## Security

- All binaries are built in GitHub Actions
- Checksums are provided for verification
- Docker images are multi-arch and minimal
- No sensitive data in release artifacts