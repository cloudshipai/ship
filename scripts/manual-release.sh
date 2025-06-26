#!/bin/bash
set -e

VERSION="${1:-v0.3.0}"

echo "Creating manual release for $VERSION..."

# Build the binaries
echo "Building binaries..."
goreleaser release --snapshot --clean

# Check if dist directory exists
if [ ! -d "dist" ]; then
    echo "Error: dist directory not found. Build failed."
    exit 1
fi

echo ""
echo "Build complete! Artifacts in dist/"
echo ""
echo "To manually create a GitHub release:"
echo "1. Go to https://github.com/cloudshipai/ship/releases/new"
echo "2. Choose tag: $VERSION"
echo "3. Release title: Ship CLI $VERSION"
echo "4. Upload these files from dist/:"
ls -la dist/*.tar.gz dist/*.zip 2>/dev/null || echo "No archives found"

echo ""
echo "Or use GitHub CLI (if authenticated):"
echo "gh release create $VERSION dist/*.tar.gz dist/*.zip --title \"Ship CLI $VERSION\" --notes \"See CHANGELOG.md for details\""