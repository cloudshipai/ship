#!/bin/bash
set -e

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m'

echo -e "${BLUE}Ship CLI - Build and Release${NC}"
echo "============================="
echo ""

# Check if GITHUB_TOKEN is set
if [ -z "$GITHUB_TOKEN" ]; then
    echo -e "${YELLOW}Warning: GITHUB_TOKEN not set${NC}"
    echo ""
    echo "Options:"
    echo "1. Set GITHUB_TOKEN environment variable"
    echo "2. Use GitHub CLI: gh auth login"
    echo "3. Create release manually after building"
    echo ""
    read -p "Continue with local build only? (y/N) " -n 1 -r
    echo ""
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        exit 0
    fi
    LOCAL_ONLY=true
else
    echo -e "${GREEN}✓ GITHUB_TOKEN found${NC}"
    LOCAL_ONLY=false
fi

# Get the version
CURRENT_TAG=$(git describe --tags --abbrev=0 2>/dev/null || echo "v0.0.0")
echo -e "${BLUE}Current version:${NC} $CURRENT_TAG"

if [ "$LOCAL_ONLY" = true ]; then
    echo ""
    echo -e "${BLUE}Building release artifacts locally...${NC}"
    make release-snapshot
    
    echo ""
    echo -e "${GREEN}✓ Build complete!${NC}"
    echo ""
    echo "Artifacts are in the dist/ directory:"
    ls -lh dist/*.tar.gz dist/*.zip 2>/dev/null | awk '{print "  " $9 " (" $5 ")"}'
    
    echo ""
    echo "To create a GitHub release manually:"
    echo "1. Go to: https://github.com/cloudshipai/ship/releases/new"
    echo "2. Choose tag: $CURRENT_TAG (or create new)"
    echo "3. Upload the files from dist/"
else
    echo ""
    echo -e "${BLUE}Creating GitHub release...${NC}"
    
    # Run GoReleaser
    if make release; then
        echo ""
        echo -e "${GREEN}✓ Release created successfully!${NC}"
        echo ""
        echo "View release at:"
        echo "https://github.com/cloudshipai/ship/releases/tag/$CURRENT_TAG"
    else
        echo ""
        echo -e "${RED}✗ Release failed${NC}"
        echo ""
        echo "You can try:"
        echo "1. Check the error messages above"
        echo "2. Run 'make release-snapshot' for local build"
        echo "3. Create release manually on GitHub"
        exit 1
    fi
fi

echo ""
echo -e "${BLUE}Installation command:${NC}"
echo "wget -qO- https://github.com/cloudshipai/ship/releases/download/$CURRENT_TAG/ship_\$(uname -s)_\$(uname -m).tar.gz | tar xz && sudo mv ship /usr/local/bin/"