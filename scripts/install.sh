#!/bin/bash

# Mitosis Chain Installation Script
# Usage:
#   curl -sSL https://raw.githubusercontent.com/mitosis-org/chain/main/scripts/install.sh | bash
#   COMPONENT=mito curl -sSL https://raw.githubusercontent.com/mitosis-org/chain/main/scripts/install.sh | bash

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
REPO="mitosis-org/chain"
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"
VERSION="${VERSION:-latest}"
COMPONENT="${COMPONENT:-mitosisd}"  # Default to mitosisd, can be 'mito'

# Utility functions
log() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

error() {
    echo -e "${RED}[ERROR]${NC} $1"
    exit 1
}

# Detect platform
detect_platform() {
    local os arch

    case "$(uname -s)" in
        Linux*)  os="linux" ;;
        Darwin*) os="darwin" ;;
        *) error "Unsupported OS: $(uname -s)" ;;
    esac

    case "$(uname -m)" in
        x86_64) arch="amd64" ;;
        arm64|aarch64) arch="arm64" ;;
        *) error "Unsupported architecture: $(uname -m)" ;;
    esac

    echo "${os}-${arch}"
}

# Get latest release version for component
get_latest_version() {
    local component="$1"
    local version

    if [[ "$component" == "mito" ]]; then
        # Get latest mito release (mito/vX.Y.Z)
        version=$(curl -s "https://api.github.com/repos/${REPO}/releases" | \
                 grep '"tag_name":' | \
                 grep '"mito/' | \
                 head -n1 | \
                 cut -d'"' -f4 | \
                 sed 's/mito\///')
    else
        # Get latest mitosisd release (vX.Y.Z)
        version=$(curl -s "https://api.github.com/repos/${REPO}/releases" | \
                 grep '"tag_name":' | \
                 grep -v '"mito/' | \
                 head -n1 | \
                 cut -d'"' -f4)
    fi

    if [[ -z "$version" ]]; then
        error "Failed to get latest version for ${component}"
    fi
    echo "$version"
}

# Get tag name for component and version
get_tag_name() {
    local component="$1"
    local version="$2"

    if [[ "$component" == "mito" ]]; then
        echo "mito/${version}"
    else
        echo "${version}"
    fi
}

# Download and install binary
install_binary() {
    local binary="$1"
    local platform="$2"
    local tag_name="$3"
    local url="https://github.com/${REPO}/releases/download/${tag_name}/${binary}-${platform}"

    log "Downloading ${binary} from tag ${tag_name} for ${platform}..."

    # Create temporary directory
    local temp_dir
    temp_dir=$(mktemp -d)
    cd "$temp_dir"

    # Download binary
    if command -v curl >/dev/null 2>&1; then
        curl -L -o "$binary" "$url"
    elif command -v wget >/dev/null 2>&1; then
        wget -O "$binary" "$url"
    else
        error "curl or wget is required"
    fi

    # Download checksum
    if command -v curl >/dev/null 2>&1; then
        curl -L -o "${binary}.sha256" "${url}.sha256"
    elif command -v wget >/dev/null 2>&1; then
        wget -O "${binary}.sha256" "${url}.sha256"
    fi

    # Verify checksum
    if command -v sha256sum >/dev/null 2>&1; then
        sha256sum -c "${binary}.sha256" || error "Checksum verification failed"
    elif command -v shasum >/dev/null 2>&1; then
        shasum -a 256 -c "${binary}.sha256" || error "Checksum verification failed"
    else
        warn "Cannot verify checksum: sha256sum or shasum not found"
    fi

    # Make executable
    chmod +x "$binary"

    # Install to target directory
    if [[ -w "$INSTALL_DIR" ]]; then
        mv "$binary" "$INSTALL_DIR/"
    else
        log "Installing to $INSTALL_DIR (requires sudo)"
        sudo mv "$binary" "$INSTALL_DIR/"
    fi

    # Cleanup
    cd /
    rm -rf "$temp_dir"

    log "Successfully installed ${binary} to ${INSTALL_DIR}/${binary}"
}

# Main installation function
main() {
    log "🚀 Installing Mitosis Chain component: ${COMPONENT}"

    # Validate component
    if [[ "$COMPONENT" != "mitosisd" && "$COMPONENT" != "mito" ]]; then
        error "Invalid component: ${COMPONENT}. Must be 'mitosisd' or 'mito'"
    fi

    # Check prerequisites
    if ! command -v curl >/dev/null 2>&1 && ! command -v wget >/dev/null 2>&1; then
        error "curl or wget is required"
    fi

    # Detect platform
    local platform
    platform=$(detect_platform)
    log "Detected platform: $platform"

    # Get version
    if [[ "$VERSION" == "latest" ]]; then
        VERSION=$(get_latest_version "$COMPONENT")
    fi
    log "Installing version: $VERSION"

    # Get tag name
    local tag_name
    tag_name=$(get_tag_name "$COMPONENT" "$VERSION")
    log "Using release tag: $tag_name"

    # Install binary
    install_binary "$COMPONENT" "$platform" "$tag_name"

    log "✅ Installation complete!"
    echo

    # Show component-specific next steps
    if [[ "$COMPONENT" == "mitosisd" ]]; then
        log "🎯 Next steps for Mitosis Chain Node:"
        echo "  1. Initialize your node: mitosisd init [moniker] --chain-id mitosis-mainnet-1"
        echo "  2. Start the node: mitosisd start"
        echo "  3. Check version: mitosisd version"
        echo "  4. View logs: mitosisd start --log_level info"
    else
        log "🎯 Next steps for Mito CLI:"
        echo "  1. Configure RPC: mito config set-rpc https://rpc.mitosis.org"
        echo "  2. Set contracts: mito config set-contract --validator-manager 0x..."
        echo "  3. Check version: mito version"
        echo "  4. View help: mito --help"
    fi

    echo
    log "📖 Documentation: https://docs.mitosis.org/developers/"
    log "💬 Community: https://discord.gg/mitosis"

    # Show installation info for other component
    if [[ "$COMPONENT" == "mitosisd" ]]; then
        echo
        log "💡 To install the Mito CLI separately:"
        echo "  COMPONENT=mito curl -sSL https://raw.githubusercontent.com/mitosis-org/chain/main/scripts/install.sh | bash"
    else
        echo
        log "💡 To install the Mitosis Chain Node separately:"
        echo "  curl -sSL https://raw.githubusercontent.com/mitosis-org/chain/main/scripts/install.sh | bash"
    fi
}

# Run main function
main "$@"