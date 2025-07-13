#!/bin/bash
# Ship CLI Cleanup Script
# This script performs comprehensive cleanup and maintenance

set -e

echo "🧹 Ship CLI Cleanup & Maintenance"
echo "=================================="

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${GREEN}✓${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}⚠${NC} $1"  
}

print_error() {
    echo -e "${RED}✗${NC} $1"
}

# Check if we're in the right directory
if [ ! -f "go.mod" ] || [ ! -d "internal" ]; then
    print_error "This script must be run from the Ship CLI root directory"
    exit 1
fi

echo "1. Cleaning build artifacts..."
rm -f ship ship-*
rm -rf dist/
rm -f *.exe *.exe~
rm -f coverage.out *.coverprofile profile.cov
print_status "Build artifacts cleaned"

echo "2. Formatting Go code..."
go fmt ./...
print_status "Go code formatted"

echo "3. Running go mod tidy..."
go mod tidy
print_status "Dependencies tidied"

echo "4. Running go vet..."
if go vet ./...; then
    print_status "Go vet passed"
else
    print_error "Go vet failed - please fix issues before proceeding"
    exit 1
fi

echo "5. Running basic tests..."
if go test -short ./...; then
    print_status "Basic tests passed"
else
    print_warning "Some tests failed - check test output above"
fi

echo "6. Checking for security vulnerabilities..."
if command -v govulncheck >/dev/null 2>&1; then
    govulncheck ./...
    print_status "Vulnerability check completed"
else
    print_warning "govulncheck not installed - run: go install golang.org/x/vuln/cmd/govulncheck@latest"
fi

echo "7. Checking for common issues..."

# Check for TODO/FIXME comments
TODO_COUNT=$(find . -name "*.go" -exec grep -l "TODO\|FIXME\|XXX\|HACK" {} \; | wc -l)
if [ $TODO_COUNT -gt 0 ]; then
    print_warning "Found $TODO_COUNT files with TODO/FIXME comments"
fi

# Check for committed secrets (basic check)
if grep -r "api.key\|secret\|password\|token" . --include="*.go" --include="*.yaml" --include="*.yml" | grep -v "test" | grep -v "example" | grep -v "TODO" >/dev/null; then
    print_warning "Potential secrets found in code - please review"
fi

# Check .gitignore coverage
if [ -f "ship" ]; then
    print_warning "Ship binary still exists - should be in .gitignore"
fi

echo "8. Repository size check..."
REPO_SIZE=$(du -sh . | cut -f1)
print_status "Repository size: $REPO_SIZE"

if [ -d ".git" ]; then
    GIT_SIZE=$(du -sh .git | cut -f1)
    echo "   Git history size: $GIT_SIZE"
fi

echo ""
echo "🎉 Cleanup completed!"
echo ""
echo "Next steps:"
echo "- Review any warnings above"
echo "- Run 'make test' for full test suite"
echo "- Run 'make lint' if golangci-lint is available"
echo "- Consider running 'make release-check' before releases"