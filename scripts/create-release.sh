#!/bin/bash

# create-release.sh - Script to create a new release with proper tagging and changelog

set -e

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Check if gh CLI is installed
if ! command -v gh &> /dev/null; then
    echo -e "${RED}Error: GitHub CLI (gh) is not installed${NC}"
    echo "Install it from: https://cli.github.com/"
    exit 1
fi

# Check if we're in a git repository
if ! git rev-parse --git-dir > /dev/null 2>&1; then
    echo -e "${RED}Error: Not in a git repository${NC}"
    exit 1
fi

# Get the current version from git tags
CURRENT_VERSION=$(git describe --tags --abbrev=0 2>/dev/null || echo "v0.0.0")
echo -e "${BLUE}Current version: ${CURRENT_VERSION}${NC}"

# Parse version components
VERSION_WITHOUT_V=${CURRENT_VERSION#v}
IFS='.' read -ra VERSION_PARTS <<< "$VERSION_WITHOUT_V"
MAJOR=${VERSION_PARTS[0]:-0}
MINOR=${VERSION_PARTS[1]:-0}
PATCH=${VERSION_PARTS[2]:-0}

# Prompt for version bump type
echo -e "${YELLOW}Select version bump type:${NC}"
echo "1) Patch (v$MAJOR.$MINOR.$((PATCH + 1)))"
echo "2) Minor (v$MAJOR.$((MINOR + 1)).0)"
echo "3) Major (v$((MAJOR + 1)).0.0)"
echo "4) Custom version"
read -p "Enter choice (1-4): " choice

case $choice in
    1)
        NEW_VERSION="v$MAJOR.$MINOR.$((PATCH + 1))"
        ;;
    2)
        NEW_VERSION="v$MAJOR.$((MINOR + 1)).0"
        ;;
    3)
        NEW_VERSION="v$((MAJOR + 1)).0.0"
        ;;
    4)
        read -p "Enter custom version (e.g., v1.2.3): " NEW_VERSION
        if [[ ! $NEW_VERSION =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
            echo -e "${RED}Error: Invalid version format. Use vX.Y.Z${NC}"
            exit 1
        fi
        ;;
    *)
        echo -e "${RED}Invalid choice${NC}"
        exit 1
        ;;
esac

echo -e "${GREEN}New version will be: ${NEW_VERSION}${NC}"

# Ensure we're on main branch
CURRENT_BRANCH=$(git branch --show-current)
if [ "$CURRENT_BRANCH" != "main" ]; then
    echo -e "${YELLOW}Warning: Not on main branch (currently on $CURRENT_BRANCH)${NC}"
    read -p "Switch to main branch? (y/n): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        git checkout main
        git pull origin main
    else
        echo -e "${RED}Aborted: Releases should be created from main branch${NC}"
        exit 1
    fi
fi

# Ensure working directory is clean
if ! git diff-index --quiet HEAD --; then
    echo -e "${RED}Error: Working directory is not clean${NC}"
    echo "Please commit or stash your changes first"
    exit 1
fi

# Pull latest changes
echo -e "${BLUE}Pulling latest changes...${NC}"
git pull origin main

# Run tests
echo -e "${BLUE}Running tests...${NC}"
if ! go test ./...; then
    echo -e "${RED}Error: Tests failed${NC}"
    exit 1
fi

# Generate release notes
echo -e "${BLUE}Generating release notes...${NC}"
RELEASE_NOTES=$(mktemp)

# Get commit messages since last tag
if [ "$CURRENT_VERSION" != "v0.0.0" ]; then
    git log $CURRENT_VERSION..HEAD --pretty=format:"- %s" | grep -v "Merge pull request" > $RELEASE_NOTES
else
    git log --pretty=format:"- %s" | grep -v "Merge pull request" > $RELEASE_NOTES
fi

# Add header to release notes
cat > $RELEASE_NOTES.final << EOF
# Ship CLI ${NEW_VERSION}

## What's Changed

$(cat $RELEASE_NOTES)

## Installation

\`\`\`bash
# macOS/Linux
curl -sSL https://github.com/cloudshipai/ship/releases/download/${NEW_VERSION}/ship_${NEW_VERSION}_\$(uname -s)_\$(uname -m).tar.gz | tar xz
sudo mv ship /usr/local/bin/

# Or using Go
go install github.com/cloudshipai/ship/cmd/ship@${NEW_VERSION}
\`\`\`

**Full Changelog**: https://github.com/cloudshipai/ship/compare/${CURRENT_VERSION}...${NEW_VERSION}
EOF

# Show release notes for review
echo -e "${BLUE}Release notes:${NC}"
cat $RELEASE_NOTES.final
echo

read -p "Continue with release? (y/n): " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    rm $RELEASE_NOTES $RELEASE_NOTES.final
    echo -e "${YELLOW}Release cancelled${NC}"
    exit 0
fi

# Create and push tag
echo -e "${BLUE}Creating tag ${NEW_VERSION}...${NC}"
git tag -a $NEW_VERSION -m "Release $NEW_VERSION"
git push origin $NEW_VERSION

echo -e "${GREEN}✓ Tag created and pushed${NC}"
echo -e "${BLUE}GitHub Actions will now build and create the release${NC}"
echo -e "${BLUE}Monitor progress at: https://github.com/cloudshipai/ship/actions${NC}"

# Clean up
rm $RELEASE_NOTES $RELEASE_NOTES.final

echo -e "${GREEN}✓ Release process initiated successfully!${NC}"