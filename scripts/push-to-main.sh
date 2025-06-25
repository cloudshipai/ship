#!/bin/bash

# push-to-main.sh - Script to safely push changes to main repository

set -e

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Repository configuration
MAIN_REPO="cloudshipai/ship"
REMOTE_NAME="origin"

echo -e "${BLUE}Ship CLI - Push to Main Repository${NC}"
echo "===================================="

# Check if we're in a git repository
if ! git rev-parse --git-dir > /dev/null 2>&1; then
    echo -e "${RED}Error: Not in a git repository${NC}"
    exit 1
fi

# Check current remote
CURRENT_REMOTE=$(git config --get remote.origin.url || echo "none")
echo -e "${BLUE}Current remote: ${CURRENT_REMOTE}${NC}"

# Set up the correct remote if needed
if [[ ! "$CURRENT_REMOTE" =~ github\.com[:/]cloudshipai/ship ]]; then
    echo -e "${YELLOW}Setting up cloudshipai/ship as origin...${NC}"
    if git remote get-url origin &>/dev/null; then
        git remote set-url origin git@github.com:cloudshipai/ship.git
    else
        git remote add origin git@github.com:cloudshipai/ship.git
    fi
    echo -e "${GREEN}✓ Remote configured${NC}"
fi

# Fetch latest from remote
echo -e "${BLUE}Fetching latest from remote...${NC}"
git fetch origin

# Check current branch
CURRENT_BRANCH=$(git branch --show-current)
echo -e "${BLUE}Current branch: ${CURRENT_BRANCH}${NC}"

# Check for uncommitted changes
if ! git diff-index --quiet HEAD -- 2>/dev/null; then
    echo -e "${YELLOW}You have uncommitted changes:${NC}"
    git status --short
    echo
    read -p "Commit these changes? (y/n): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        read -p "Enter commit message: " COMMIT_MSG
        git add -A
        git commit -m "$COMMIT_MSG"
        echo -e "${GREEN}✓ Changes committed${NC}"
    else
        echo -e "${RED}Please commit or stash your changes first${NC}"
        exit 1
    fi
fi

# Run tests before pushing
echo -e "${BLUE}Running tests...${NC}"
if ! go test ./... > /dev/null 2>&1; then
    echo -e "${RED}✗ Tests failed!${NC}"
    echo "Run 'go test ./...' to see details"
    exit 1
fi
echo -e "${GREEN}✓ Tests passed${NC}"

# Check if we need to merge/rebase
if [ "$CURRENT_BRANCH" = "main" ]; then
    # On main branch - check if we're behind
    LOCAL=$(git rev-parse @)
    REMOTE=$(git rev-parse @{u} 2>/dev/null || echo "none")
    BASE=$(git merge-base @ @{u} 2>/dev/null || echo "none")

    if [ "$REMOTE" != "none" ] && [ "$LOCAL" != "$REMOTE" ]; then
        if [ "$LOCAL" = "$BASE" ]; then
            echo -e "${YELLOW}Your branch is behind origin/main${NC}"
            read -p "Pull latest changes? (y/n): " -n 1 -r
            echo
            if [[ $REPLY =~ ^[Yy]$ ]]; then
                git pull --rebase origin main
                echo -e "${GREEN}✓ Pulled latest changes${NC}"
            fi
        elif [ "$REMOTE" = "$BASE" ]; then
            echo -e "${GREEN}Your branch is ahead of origin/main${NC}"
        else
            echo -e "${YELLOW}Your branch has diverged from origin/main${NC}"
            echo "You may need to rebase or merge"
        fi
    fi
else
    # On feature branch
    echo -e "${YELLOW}You're on branch '${CURRENT_BRANCH}'${NC}"
    read -p "Push this branch? (y/n): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        exit 0
    fi
fi

# Show what will be pushed
echo -e "${BLUE}Changes to be pushed:${NC}"
git log --oneline origin/${CURRENT_BRANCH}..HEAD 2>/dev/null || git log --oneline -5

echo
read -p "Push to origin/${CURRENT_BRANCH}? (y/n): " -n 1 -r
echo

if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo -e "${BLUE}Pushing to origin/${CURRENT_BRANCH}...${NC}"
    git push origin ${CURRENT_BRANCH}
    echo -e "${GREEN}✓ Successfully pushed to origin/${CURRENT_BRANCH}${NC}"
    
    # Provide PR link if not on main
    if [ "$CURRENT_BRANCH" != "main" ]; then
        echo
        echo -e "${BLUE}Create a pull request:${NC}"
        echo "https://github.com/cloudshipai/ship/compare/${CURRENT_BRANCH}?expand=1"
    fi
else
    echo -e "${YELLOW}Push cancelled${NC}"
fi