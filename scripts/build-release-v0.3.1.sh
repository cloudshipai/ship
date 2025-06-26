#!/bin/bash
set -e

echo "Building Ship CLI v0.3.1 Release"
echo "================================"
echo ""

# Clean and build with the correct tag
export VERSION=v0.3.1
export COMMIT=$(git rev-parse HEAD)
export DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ")

echo "Version: $VERSION"
echo "Commit: $COMMIT"
echo "Date: $DATE"
echo ""

# Build with GoReleaser using the actual tag
GITHUB_TOKEN="" goreleaser release --clean --skip=publish,announce,validate \
  --snapshot --config .goreleaser.yml

echo ""
echo "Build complete! Files in dist/"
echo ""

# Rename files to remove the snapshot suffix
cd dist/
for file in ship_*-next_*.tar.gz ship_*-next_*.zip; do
  if [ -f "$file" ]; then
    newname=$(echo "$file" | sed 's/0\.3\.2-next/0.3.1/')
    mv "$file" "$newname"
    echo "Renamed: $file -> $newname"
  fi
done

# Update checksums
rm -f checksums.txt
sha256sum ship_*.tar.gz ship_*.zip > checksums.txt

echo ""
echo "Release artifacts ready for upload:"
echo "=================================="
ls -lh ship_*.tar.gz ship_*.zip | awk '{print $9 " (" $5 ")"}'

echo ""
echo "To create the GitHub release:"
echo "1. Go to: https://github.com/cloudshipai/ship/releases/new"
echo "2. Choose existing tag: v0.3.1"
echo "3. Release title: Ship CLI v0.3.1"
echo "4. Upload these files:"
pwd
ls ship_*.tar.gz ship_*.zip

echo ""
echo "5. Add release notes:"
cat << 'EOF'
## Ship CLI Release v0.3.1

CloudshipAI CLI - Infrastructure analysis tools at your fingertips 🚀

### What's New
- ✅ Fixed CI/CD pipeline and release process
- ✅ Comprehensive release tooling with Makefile
- ✅ Automated version management
- ✅ Multi-platform binary releases

### Features
- 🔍 Terraform Linting with TFLint
- 🛡️ Security Scanning with Checkov and Trivy  
- 💰 Cost Estimation with Infracost and OpenInfraQuote
- 📝 Documentation Generation with terraform-docs
- 🐳 All tools run in containers via Dagger
- 🔧 Multi-platform support (Linux, macOS, Windows)

### Installation

```bash
# Linux/macOS
wget -qO- https://github.com/cloudshipai/ship/releases/download/v0.3.1/ship_$(uname -s)_$(uname -m).tar.gz | tar xz && sudo mv ship /usr/local/bin/

# Verify
ship version
```

### Checksums
See checksums.txt in release assets
EOF