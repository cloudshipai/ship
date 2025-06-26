#!/bin/bash
set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_info() {
    echo -e "${BLUE}$1${NC}"
}

print_success() {
    echo -e "${GREEN}$1${NC}"
}

print_warning() {
    echo -e "${YELLOW}$1${NC}"
}

print_error() {
    echo -e "${RED}$1${NC}"
}

# Check if git is clean
check_git_status() {
    if [[ -n $(git status -s) ]]; then
        print_error "❌ Git working directory is not clean. Commit or stash changes first."
        exit 1
    fi
}

# Get current version
get_current_version() {
    git describe --tags --abbrev=0 2>/dev/null || echo "v0.0.0"
}

# Increment version
increment_version() {
    local version=$1
    local increment_type=$2
    
    # Remove 'v' prefix
    version=${version#v}
    
    # Split version
    IFS='.' read -r major minor patch <<< "$version"
    
    case $increment_type in
        major)
            major=$((major + 1))
            minor=0
            patch=0
            ;;
        minor)
            minor=$((minor + 1))
            patch=0
            ;;
        patch)
            patch=$((patch + 1))
            ;;
        *)
            print_error "Invalid increment type: $increment_type"
            exit 1
            ;;
    esac
    
    echo "v${major}.${minor}.${patch}"
}

# Main script
main() {
    print_info "🚀 Ship CLI Release Script"
    print_info "========================="
    echo ""
    
    # Check prerequisites
    print_info "Checking prerequisites..."
    
    if ! command -v goreleaser &> /dev/null; then
        print_error "❌ goreleaser not found. Install it first:"
        echo "  brew install goreleaser"
        echo "  or"
        echo "  go install github.com/goreleaser/goreleaser@latest"
        exit 1
    fi
    
    if ! command -v gh &> /dev/null; then
        print_warning "⚠️  GitHub CLI not found. Install for better integration:"
        echo "  brew install gh"
    fi
    
    # Check git status
    check_git_status
    
    # Get version increment type
    INCREMENT_TYPE=${1:-patch}
    
    # Get current version
    CURRENT_VERSION=$(get_current_version)
    print_info "Current version: $CURRENT_VERSION"
    
    # Calculate new version
    if [[ $INCREMENT_TYPE == "custom" ]]; then
        read -p "Enter custom version (e.g., v1.2.3): " NEW_VERSION
        if [[ ! $NEW_VERSION =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
            print_error "❌ Invalid version format. Use v1.2.3 format."
            exit 1
        fi
    else
        NEW_VERSION=$(increment_version "$CURRENT_VERSION" "$INCREMENT_TYPE")
    fi
    
    print_info "New version: $NEW_VERSION"
    echo ""
    
    # Confirm release
    read -p "Create release $NEW_VERSION? (y/N) " -n 1 -r
    echo ""
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        print_warning "Release cancelled."
        exit 0
    fi
    
    # Pull latest changes
    print_info "Pulling latest changes..."
    git pull origin main
    
    # Run tests
    print_info "Running tests..."
    if ! make test; then
        print_error "❌ Tests failed. Fix them before releasing."
        exit 1
    fi
    
    # Check GoReleaser config
    print_info "Checking GoReleaser configuration..."
    if ! goreleaser check; then
        print_error "❌ GoReleaser configuration is invalid."
        exit 1
    fi
    
    # Create and push tag
    print_info "Creating tag $NEW_VERSION..."
    git tag -a "$NEW_VERSION" -m "Release $NEW_VERSION"
    
    # Update changelog
    print_info "Updating CHANGELOG.md..."
    cat > CHANGELOG.md.tmp << EOF
# Changelog

## [$NEW_VERSION] - $(date +%Y-%m-%d)

### Added
- See commit history for changes

EOF
    
    if [ -f CHANGELOG.md ]; then
        tail -n +2 CHANGELOG.md >> CHANGELOG.md.tmp
    fi
    
    mv CHANGELOG.md.tmp CHANGELOG.md
    git add CHANGELOG.md
    git commit -m "chore: update changelog for $NEW_VERSION" || true
    
    # Push changes
    print_info "Pushing changes and tag..."
    git push origin main
    git push origin "$NEW_VERSION"
    
    # Check for GitHub token
    if [ -z "$GITHUB_TOKEN" ]; then
        print_warning "⚠️  GITHUB_TOKEN not set. Trying to use gh auth token..."
        
        if command -v gh &> /dev/null && gh auth status &> /dev/null; then
            export GITHUB_TOKEN=$(gh auth token)
            print_success "✅ Using GitHub CLI token"
        else
            print_error "❌ No GitHub authentication found."
            print_info ""
            print_info "To create the release, either:"
            print_info "1. Set GITHUB_TOKEN environment variable:"
            print_info "   export GITHUB_TOKEN=your-github-personal-access-token"
            print_info "   make release"
            print_info ""
            print_info "2. Or authenticate with GitHub CLI:"
            print_info "   gh auth login"
            print_info "   make release"
            print_info ""
            print_info "3. Or create the release manually at:"
            print_info "   https://github.com/cloudshipai/ship/releases/new"
            print_info "   Tag: $NEW_VERSION"
            exit 1
        fi
    fi
    
    # Create release with GoReleaser
    print_info "Creating release with GoReleaser..."
    if goreleaser release --clean; then
        print_success "✅ Release $NEW_VERSION created successfully!"
        print_info ""
        print_info "📦 Release URL: https://github.com/cloudshipai/ship/releases/tag/$NEW_VERSION"
        print_info ""
        print_info "🎉 Installation command:"
        echo "wget -qO- https://github.com/cloudshipai/ship/releases/download/$NEW_VERSION/ship_\$(uname -s)_\$(uname -m).tar.gz | tar xz && sudo mv ship /usr/local/bin/"
    else
        print_error "❌ GoReleaser failed. Check the logs above."
        print_info ""
        print_info "You can try creating the release manually:"
        print_info "1. Build locally: make release-snapshot"
        print_info "2. Upload at: https://github.com/cloudshipai/ship/releases/new"
        exit 1
    fi
}

# Run main function
main "$@"