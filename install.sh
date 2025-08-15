#!/bin/bash

# Ship CLI Installation Script
# Usage: curl -fsSL https://raw.githubusercontent.com/cloudshipai/ship/main/install.sh | bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Configuration
REPO="cloudshipai/ship"
BINARY_NAME="ship"
INSTALL_DIR=""  # Will be determined by get_install_dir()
SHIP_VERSION="v0.6.3"  # Updated by CI on release

# Banner
print_banner() {
    echo -e "${BLUE}"
    echo "🚢 Ship CLI Installation Script"
    echo "================================="
    echo -e "${NC}"
    echo "CloudShip AI Team Curated MCP Servers & SDK"
    echo "Run infrastructure tools securely in containers"
    echo ""
}

# Logging functions
log_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

log_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

log_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

log_error() {
    echo -e "${RED}❌ $1${NC}"
    exit 1
}

# Platform detection
detect_platform() {
    local os arch
    
    # Detect OS
    case "$(uname -s)" in
        Darwin)
            os="darwin"
            ;;
        Linux)
            os="linux"
            ;;
        MINGW*|MSYS*|CYGWIN*)
            os="windows"
            ;;
        *)
            log_error "Unsupported operating system: $(uname -s)"
            ;;
    esac
    
    # Detect architecture
    case "$(uname -m)" in
        x86_64|amd64)
            arch="amd64"
            ;;
        arm64|aarch64)
            arch="arm64"
            ;;
        i386|i686)
            arch="386"
            ;;
        armv7*)
            arch="arm"
            ;;
        *)
            log_error "Unsupported architecture: $(uname -m)"
            ;;
    esac
    
    echo "${os}_${arch}"
}

# Check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Download file with fallback methods
download_file() {
    local url="$1"
    local output="$2"
    
    if command_exists curl; then
        log_info "Downloading with curl..."
        curl -fsSL -o "$output" "$url"
    elif command_exists wget; then
        log_info "Downloading with wget..."
        wget -q -O "$output" "$url"
    else
        log_error "Neither curl nor wget is available. Please install one of them."
    fi
}

# Get latest release version from GitHub API
get_latest_version() {
    local api_url="https://api.github.com/repos/$REPO/releases/latest"
    local version
    
    if command_exists curl; then
        version=$(curl -fsSL "$api_url" | grep '"tag_name":' | sed -E 's/.*"tag_name": "([^"]+)".*/\1/')
    elif command_exists wget; then
        version=$(wget -qO- "$api_url" | grep '"tag_name":' | sed -E 's/.*"tag_name": "([^"]+)".*/\1/')
    else
        log_error "Cannot fetch latest version. Please install curl or wget."
    fi
    
    if [ -z "$version" ]; then
        log_warning "Could not determine latest version, using fallback"
        version="v0.1.0"  # Fallback version
    fi
    
    echo "$version"
}

# Determine the best install directory
get_install_dir() {
    # Priority order for installation paths
    local user_bin="$HOME/.local/bin"
    local homebrew_bin="/opt/homebrew/bin"  # Apple Silicon macOS
    local old_homebrew_bin="/usr/local/bin"  # Intel macOS / Linux
    local system_bin="/usr/bin"
    
    # Always prefer user-local installation first
    if [ -w "$HOME" ]; then
        mkdir -p "$user_bin" 2>/dev/null
        if [ -d "$user_bin" ]; then
            echo "$user_bin"
            return 0
        fi
    fi
    
    # macOS: Check Homebrew paths (these are usually user-writable)
    if [[ "$(uname -s)" == "Darwin" ]]; then
        if [ -d "$homebrew_bin" ] && [ -w "$homebrew_bin" ]; then
            echo "$homebrew_bin"
            return 0
        fi
        if [ -d "$old_homebrew_bin" ] && [ -w "$old_homebrew_bin" ]; then
            echo "$old_homebrew_bin"
            return 0
        fi
    fi
    
    # Check if we can write to /usr/local/bin
    if [ -w "$old_homebrew_bin" ] || ([ -d "$old_homebrew_bin" ] && command_exists sudo && sudo -n true 2>/dev/null); then
        echo "$old_homebrew_bin"
        return 0
    fi
    
    # Last resort: system bin (requires sudo)
    if command_exists sudo; then
        echo "$system_bin"
        return 0
    fi
    
    # Fallback to user bin even if we couldn't create it initially
    echo "$user_bin"
    return 0
}

# Check if we need sudo for a given directory
needs_sudo() {
    local target_dir="$1"
    [ ! -w "$target_dir" ] && [ "$target_dir" != "$HOME/.local/bin" ]
}

# Install binary
install_binary() {
    local platform="$1"
    local version="$2"
    local temp_dir
    
    temp_dir=$(mktemp -d)
    cd "$temp_dir"
    
    # Construct download URL - GoReleaser uses ProjectName_Os_Arch format (no version)
    # Convert platform to GoReleaser format
    local goreleaser_platform
    case "$platform" in
        linux_amd64)
            goreleaser_platform="Linux_x86_64"
            ;;
        linux_arm64)
            goreleaser_platform="Linux_arm64"
            ;;
        darwin_amd64)
            goreleaser_platform="Darwin_x86_64"
            ;;
        darwin_arm64)
            goreleaser_platform="Darwin_arm64"
            ;;
        windows_amd64)
            goreleaser_platform="Windows_x86_64"
            ;;
        windows_arm64)
            goreleaser_platform="Windows_arm64"
            ;;
        *)
            log_error "Unsupported platform for GoReleaser: $platform"
            ;;
    esac
    
    local filename="ship_${goreleaser_platform}"
    local archive_name="${filename}.tar.gz"
    
    local download_url="https://github.com/$REPO/releases/download/$version/$archive_name"
    
    log_info "Downloading Ship CLI $version for $platform..."
    log_info "URL: $download_url"
    
    # Download the archive
    if ! download_file "$download_url" "$archive_name"; then
        log_error "Failed to download Ship CLI. Please check your internet connection and try again."
    fi
    
    # Extract the archive
    log_info "Extracting archive..."
    tar -xzf "$archive_name"
    
    # Find the binary
    local binary_path
    if [[ "$platform" == "windows"* ]]; then
        binary_path="$BINARY_NAME.exe"
    else
        binary_path="$BINARY_NAME"
    fi
    
    if [ ! -f "$binary_path" ]; then
        log_error "Binary not found in archive. Expected: $binary_path"
    fi
    
    # Make binary executable
    chmod +x "$binary_path"
    
    # Determine install directory
    INSTALL_DIR=$(get_install_dir)
    
    # Create directory if it doesn't exist
    if [ "$INSTALL_DIR" = "$HOME/.local/bin" ]; then
        mkdir -p "$INSTALL_DIR"
    elif [ ! -d "$INSTALL_DIR" ]; then
        if needs_sudo "$INSTALL_DIR"; then
            sudo mkdir -p "$INSTALL_DIR"
        else
            mkdir -p "$INSTALL_DIR"
        fi
    fi
    
    # Install binary
    log_info "Installing Ship CLI to $INSTALL_DIR..."
    
    if needs_sudo "$INSTALL_DIR"; then
        sudo cp "$binary_path" "$INSTALL_DIR/$BINARY_NAME"
        sudo chmod +x "$INSTALL_DIR/$BINARY_NAME"
    else
        cp "$binary_path" "$INSTALL_DIR/$BINARY_NAME"
        chmod +x "$INSTALL_DIR/$BINARY_NAME"
    fi
    
    # Check if install directory is in PATH and provide guidance
    if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
        echo ""
        log_warning "$INSTALL_DIR is not in your PATH."
        
        # Provide platform-specific PATH instructions
        if [[ "$(uname -s)" == "Darwin" ]]; then
            if [[ "$SHELL" == *"zsh"* ]]; then
                echo "Add it by running:"
                echo "  echo 'export PATH=\"$INSTALL_DIR:\$PATH\"' >> ~/.zshrc"
                echo "  source ~/.zshrc"
            else
                echo "Add it by running:"
                echo "  echo 'export PATH=\"$INSTALL_DIR:\$PATH\"' >> ~/.bash_profile"
                echo "  source ~/.bash_profile"
            fi
        else
            echo "Add it by running:"
            echo "  echo 'export PATH=\"$INSTALL_DIR:\$PATH\"' >> ~/.bashrc"
            echo "  source ~/.bashrc"
        fi
    fi
    
    # Cleanup
    cd /
    rm -rf "$temp_dir"
    
    log_success "Ship CLI installed successfully!"
}

# Verify installation
verify_installation() {
    log_info "Verifying installation..."
    
    if command_exists "$BINARY_NAME"; then
        local version_output
        version_output=$($BINARY_NAME version 2>&1 || echo "unknown")
        log_success "Ship CLI is installed and working!"
        echo "  Version: $version_output"
        echo "  Location: $(which $BINARY_NAME)"
    else
        log_error "Installation verification failed. Ship CLI command not found in PATH."
    fi
}

# Print next steps
print_next_steps() {
    echo ""
    echo -e "${PURPLE}🎉 Installation Complete!${NC}"
    echo ""
    echo -e "${CYAN}Next Steps:${NC}"
    echo ""
    echo -e "${YELLOW}CLI Usage:${NC}"
    echo "1. Run Terraform analysis:"
    echo "   cd your-terraform-project"
    echo "   ${BINARY_NAME} tf lint        # Linting with TFLint"
    echo "   ${BINARY_NAME} tf checkov     # Security scanning"
    echo "   ${BINARY_NAME} tf cost        # Cost estimation"
    echo "   ${BINARY_NAME} tf docs        # Generate documentation"
    echo ""
    echo "2. Generate infrastructure diagrams:"
    echo "   ${BINARY_NAME} tf diagram . --hcl -o infrastructure.png"
    echo ""
    echo "3. Start MCP server for AI assistants:"
    echo "   ${BINARY_NAME} mcp"
    echo ""
    echo -e "${YELLOW}Ship SDK Development:${NC}"
    echo "1. Create custom MCP server:"
    echo "   mkdir my-mcp-server && cd my-mcp-server"
    echo "   go mod init my-mcp-server"
    echo "   go get github.com/cloudshipai/ship/pkg/ship"
    echo ""
    echo "2. Use Ship tools in your server:"
    echo "   go get github.com/cloudshipai/ship/pkg/tools"
    echo ""
    echo -e "${CYAN}Documentation:${NC}"
    echo "• CLI Reference: https://github.com/cloudshipai/ship#-command-reference"
    echo "• Ship SDK Overview: https://github.com/cloudshipai/ship/blob/main/docs/ship-sdk-overview.md"
    echo "• Quick Start Guide: https://github.com/cloudshipai/ship/blob/main/docs/ship-sdk-quickstart.md"
    echo "• API Reference: https://github.com/cloudshipai/ship/blob/main/docs/ship-sdk-api-reference.md"
    echo "• GitHub: https://github.com/$REPO"
    echo ""
    echo -e "${GREEN}Happy shipping! 🚢${NC}"
}

# Main installation flow
main() {
    print_banner
    
    # Parse command line arguments
    while [[ $# -gt 0 ]]; do
        case $1 in
            --version)
                SHIP_VERSION="v0.6.3"
                shift 2
                ;;
            --install-dir)
                INSTALL_DIR="$2"
                shift 2
                ;;
            --help)
                echo "Usage: $0 [OPTIONS]"
                echo ""
                echo "Options:"
                echo "  --version VERSION    Install specific version (default: latest)"
                echo "  --install-dir DIR    Install directory (default: auto-detect)"
                echo "  --help               Show this help message"
                exit 0
                ;;
            *)
                log_error "Unknown option: $1"
                ;;
        esac
    done
    
    # Check prerequisites
    log_info "Checking prerequisites..."
    
    if ! command_exists tar; then
        log_error "tar is required but not installed."
    fi
    
    # Detect platform
    log_info "Detecting platform..."
    local platform
    platform=$(detect_platform)
    log_success "Detected platform: $platform"
    
    # Get version
    if [ "$SHIP_VERSION" = "latest" ]; then
        log_info "Fetching latest version..."
        SHIP_VERSION=$(get_latest_version)
    fi
    log_success "Target version: $SHIP_VERSION"
    
    # Install
    install_binary "$platform" "$SHIP_VERSION"
    
    # Verify
    verify_installation
    
    # Print next steps
    print_next_steps
}

# Handle Ctrl+C
trap 'echo -e "\n${RED}Installation interrupted by user.${NC}"; exit 1' INT

# Run main function
main "$@"