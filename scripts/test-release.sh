#!/bin/bash
set -e

echo "Testing GoReleaser locally..."
echo "================================"

# Check if we have a tag
if ! git describe --tags --exact-match 2>/dev/null; then
    echo "No tag found on current commit. Creating a test tag..."
    git tag -a v0.2.0 -m "Release v0.2.0" || echo "Tag v0.2.0 already exists"
fi

# Run GoReleaser in snapshot mode first
echo ""
echo "Running GoReleaser in snapshot mode..."
goreleaser release --snapshot --clean

echo ""
echo "Snapshot build complete! Check the dist/ directory for artifacts."
echo ""
echo "To create a real release:"
echo "1. Push the tag: git push origin v0.2.0"
echo "2. GitHub Actions will automatically create the release"
echo ""
echo "Or manually create a release with:"
echo "GITHUB_TOKEN=your-token-here goreleaser release --clean"