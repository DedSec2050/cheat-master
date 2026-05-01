#!/bin/bash

################################################################################
# Cheat Master - Azure VM Linux Build Script
# Target: Ubuntu 24.04 (x64/amd64)
# Purpose: Cross-compile Go binaries for Azure VM deployment
################################################################################

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
PROJECT_NAME="cheat-master"
BUILD_DIR="dist"
GOOS="linux"
GOARCH="amd64"
GO_VERSION="1.20"  # Minimum Go version
UBUNTU_VERSION="24.04"

# Timestamp
BUILD_DATE=$(date -u '+%Y-%m-%dT%H:%M:%SZ')
GIT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")
VERSION=${VERSION:-"1.0.0"}

################################################################################
# Helper Functions
################################################################################

print_header() {
    echo -e "\n${BLUE}═══════════════════════════════════════════════════════${NC}"
    echo -e "${BLUE}$1${NC}"
    echo -e "${BLUE}═══════════════════════════════════════════════════════${NC}\n"
}

print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

print_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

################################################################################
# Pre-flight Checks
################################################################################

check_environment() {
    print_header "Pre-flight Environment Check"
    
    # Check Go installation
    if ! command -v go &> /dev/null; then
        print_error "Go is not installed"
        echo "Install from: https://golang.org/dl/"
        exit 1
    fi
    
    GO_INSTALLED=$(go version | awk '{print $3}')
    print_success "Go installed: $GO_INSTALLED"
    
    # Check Go version
    GO_MIN_VERSION=$(echo "$GO_VERSION" | sed 's/go//')
    GO_CURR_VERSION=$(go version | grep -oP 'go\K[0-9]+\.[0-9]+' | head -1)
    
    if [[ "$GO_CURR_VERSION" < "$GO_MIN_VERSION" ]]; then
        print_warning "Go version should be $GO_VERSION or higher (current: $GO_CURR_VERSION)"
    else
        print_success "Go version check passed"
    fi
    
    # Check Git
    if command -v git &> /dev/null; then
        print_success "Git available"
    else
        print_warning "Git not found - build metadata will be limited"
    fi
    
    # Check working directory
    if [ ! -f "go.mod" ]; then
        print_error "go.mod not found. Run from project root directory"
        exit 1
    fi
    print_success "Project root verified"
}

################################################################################
# Build Selection
################################################################################

select_build_mode() {
    print_header "Select Build Mode"
    
    echo "1) CLI Mode (cheat-cli)"
    echo "   └─ Single-user command-line interface"
    echo "   └─ Best for: Quick testing, one-off courses"
    echo ""
    echo "2) Service Mode (cheat-service) ⭐ RECOMMENDED for Azure VM"
    echo "   └─ Multi-user production server"
    echo "   └─ Best for: Azure VM deployment, concurrent jobs"
    echo ""
    echo "3) Both (CLI + Service)"
    echo "   └─ Build both binaries"
    echo ""
    
    read -p "Choose build mode (1-3, default 2): " BUILD_MODE
    BUILD_MODE=${BUILD_MODE:-2}
    
    case $BUILD_MODE in
        1)
            BINARIES=("cli")
            print_info "Selected: CLI mode"
            ;;
        2)
            BINARIES=("service")
            print_info "Selected: Service mode (Recommended for Azure VM)"
            ;;
        3)
            BINARIES=("cli" "service")
            print_info "Selected: Both CLI and Service"
            ;;
        *)
            print_error "Invalid selection"
            exit 1
            ;;
    esac
}

################################################################################
# Build Function
################################################################################

build_binary() {
    local binary_name=$1
    local cmd_path="cmd/${binary_name}"
    local output_file="${BUILD_DIR}/cheat-${binary_name}"
    
    if [ ! -d "$cmd_path" ]; then
        print_error "Build path not found: $cmd_path"
        return 1
    fi
    
    print_info "Building: $binary_name"
    print_info "Target: $GOOS/$GOARCH (Ubuntu 24.04)"
    
    # Build with linker flags for version info
    CGO_ENABLED=0 \
    GOOS=$GOOS \
    GOARCH=$GOARCH \
    go build \
        -v \
        -o "$output_file" \
        -ldflags="-s -w \
            -X 'main.Version=$VERSION' \
            -X 'main.BuildDate=$BUILD_DATE' \
            -X 'main.GitCommit=$GIT_COMMIT' \
            -X 'main.GoVersion=$(go version | awk '{print $3}')'" \
        "./$cmd_path" || {
        print_error "Build failed for $binary_name"
        return 1
    }
    
    # Verify binary exists
    if [ ! -f "$output_file" ]; then
        print_error "Binary not created: $output_file"
        return 1
    fi
    
    # Get file info
    local file_size=$(du -h "$output_file" | cut -f1)
    local file_info=$(file "$output_file")
    
    print_success "Built: $output_file ($file_size)"
    print_info "Binary info: $file_info"
}

################################################################################
# Checksum Generation
################################################################################

generate_checksums() {
    print_info "Generating SHA256 checksums..."
    
    cd "$BUILD_DIR"
    shasum -a 256 cheat-* > SHA256SUMS || {
        print_warning "Failed to generate checksums"
        return 1
    }
    
    print_success "Checksums generated: SHA256SUMS"
    cat SHA256SUMS
    cd - > /dev/null
}

################################################################################
# Build Report
################################################################################

print_build_report() {
    print_header "Build Report"
    
    echo "Project:           $PROJECT_NAME"
    echo "Version:           $VERSION"
    echo "Build Date:        $BUILD_DATE"
    echo "Git Commit:        $GIT_COMMIT"
    echo ""
    echo "Target OS:         $GOOS"
    echo "Target Arch:       $GOARCH"
    echo "Ubuntu Version:    $UBUNTU_VERSION"
    echo ""
    echo "Build Directory:   $(pwd)/$BUILD_DIR"
    echo ""
    
    echo "Generated Binaries:"
    for binary in "${BINARIES[@]}"; do
        local file="${BUILD_DIR}/cheat-${binary}"
        if [ -f "$file" ]; then
            local size=$(du -h "$file" | cut -f1)
            echo "  ✓ $file ($size)"
        fi
    done
    echo ""
    
    if [ -f "${BUILD_DIR}/SHA256SUMS" ]; then
        echo "SHA256 Checksums available: ${BUILD_DIR}/SHA256SUMS"
        echo ""
        cat "${BUILD_DIR}/SHA256SUMS"
    fi
}

################################################################################
# Deployment Guidance
################################################################################

print_deployment_guide() {
    print_header "Deployment Guide for Azure VM (Ubuntu 24.04)"
    
    if [[ " ${BINARIES[@]} " =~ " service " ]]; then
        echo "📦 Service Mode (Recommended for Azure VM)"
        echo ""
        echo "Deploy to Azure VM:"
        echo "   $ scp dist/cheat-service azureuser@<vm-ip>:~/"
        echo "   $ ssh azureuser@<vm-ip>"
        echo "   $ chmod +x cheat-service"
        echo "   $ ./cheat-service"
        echo ""
        echo "For systemd setup, see: PRODUCTION_DEPLOYMENT.md"
        echo ""
    fi
    
    if [[ " ${BINARIES[@]} " =~ " cli " ]]; then
        echo "📦 CLI Mode (Single-user, quick testing)"
        echo ""
        echo "Deploy to Azure VM:"
        echo "   $ scp dist/cheat-cli azureuser@<vm-ip>:~/"
        echo "   $ ssh azureuser@<vm-ip>"
        echo "   $ chmod +x cheat-cli"
        echo "   $ ./cheat-cli '<course-slug>' '<email>' '<password>'"
        echo ""
    fi
}

################################################################################
# Verification
################################################################################

verify_binary() {
    local binary=$1
    local file="${BUILD_DIR}/cheat-${binary}"
    
    print_info "Verifying: cheat-${binary}"
    
    # Check if executable
    if [ -x "$file" ]; then
        print_success "Binary is executable"
    else
        print_warning "Binary is not executable, fixing..."
        chmod +x "$file"
    fi
    
    # Check file type
    local file_type=$(file "$file" | grep -o "x86-64" || echo "unknown")
    if [ "$file_type" = "x86-64" ]; then
        print_success "Correct architecture: x86-64 (amd64)"
    fi
    
    # Check if statically linked
    local file_output=$(file "$file")
    if echo "$file_output" | grep -q "statically linked"; then
        print_success "Statically linked (no dependencies needed on VM)"
    elif command -v ldd &> /dev/null; then
        # Try ldd on Linux systems
        local is_static=$(ldd "$file" 2>&1 | grep -c "not a dynamic executable" || echo "0")
        if [ "$is_static" -eq 1 ]; then
            print_success "Statically linked (no dependencies needed on VM)"
        else
            print_warning "Binary has some dynamic dependencies"
        fi
    fi
}



################################################################################
# Main Script
################################################################################

main() {
    clear
    
    print_header "Cheat Master - Azure VM Linux Build System"
    echo "Target: Ubuntu 24.04 (x64/amd64)"
    echo "Purpose: Cross-compile Go binaries for Azure VM deployment"
    
    # Run checks
    check_environment
    
    # Create build directory
    mkdir -p "$BUILD_DIR"
    print_success "Build directory ready: $BUILD_DIR"
    
    # Select build mode
    select_build_mode
    
    # Build binaries
    print_header "Building Binaries"
    for binary in "${BINARIES[@]}"; do
        if build_binary "$binary"; then
            verify_binary "$binary"
        fi
    done
    
    # Generate checksums
    generate_checksums
    
    # Print reports
    print_build_report
    print_deployment_guide
    
    # Final success message
    print_header "Build Complete ✅"
    echo "Raw binaries ready in: $BUILD_DIR"
    echo ""
    echo "Binary files:"
    for binary in "${BINARIES[@]}"; do
        echo "  → dist/cheat-${binary}"
    done
    echo ""
    echo "Deployment:"
    echo "  $ scp dist/cheat-* azureuser@<vm-ip>:~/"
    echo "  $ ssh azureuser@<vm-ip>"
    echo "  $ chmod +x cheat-*"
    echo "  $ ./cheat-service (or ./cheat-cli)"
    echo ""
}

# Run main function
main "$@"
