#!/bin/sh

# Clipper version to install
VERSION="1.3.0"

# Determine OS and Architecture
OS=$(uname -s)
ARCH=$(uname -m)

# Create a temporary directory for downloading the binary
TMP_DIR=$(mktemp -d)
TMP_BINARY="$TMP_DIR/clipper"

# Function to set the download URL based on OS and Architecture
set_binary_url() {
    case "$OS" in
    "Linux")
        case "$ARCH" in
        "x86_64")
            BINARY_URL="https://github.com/supitsdu/clipper/releases/download/v$VERSION/clipper_linux_amd64_v$VERSION"
            ;;
        "aarch64" | "arm64")
            BINARY_URL="https://github.com/supitsdu/clipper/releases/download/v$VERSION/clipper_linux_arm64_v$VERSION"
            ;;
        *)
            echo "Error: Unsupported architecture: $ARCH"
            exit 1
            ;;
        esac
        ;;
    "Darwin")
        case "$ARCH" in
        "x86_64")
            BINARY_URL="https://github.com/supitsdu/clipper/releases/download/v$VERSION/clipper_darwin_amd64_v$VERSION"
            ;;
        "arm64")
            BINARY_URL="https://github.com/supitsdu/clipper/releases/download/v$VERSION/clipper_darwin_arm64_v$VERSION"
            ;;
        *)
            echo "Error: Unsupported architecture: $ARCH"
            exit 1
            ;;
        esac
        ;;
    *)
        echo "Error: Unsupported OS: $OS"
        exit 1
        ;;
    esac
}

# Function to download the binary
download_binary() {
    echo "Downloading Clipper from $BINARY_URL..."
    curl -Lo "$TMP_BINARY" "$BINARY_URL" >/dev/null 2>&1 || {
        echo "Error: Failed to download Clipper."
        exit 1
    }
    echo "Download completed successfully."
}

# Function to make the binary executable
make_executable() {
    echo "Making Clipper executable..."
    chmod +x "$TMP_BINARY"
    echo "Clipper is now executable."
}

# Function to move the binary to /usr/local/bin
install_binary() {
    if [ -d "/usr/local/bin" ] && echo "$PATH" | grep -q "/usr/local/bin"; then
        echo "Installing Clipper to /usr/local/bin..."
        sudo mv -f "$TMP_BINARY" /usr/local/bin/clipper >/dev/null 2>&1 || {
            echo "Error: Failed to move Clipper to /usr/local/bin."
            exit 1
        }
        echo "Clipper installed to /usr/local/bin."
    else
        echo "Error: Directory /usr/local/bin does not exist or is not in your PATH."
        exit 1
    fi
}

# Function to verify the installation
verify_installation() {
    if command -v clipper >/dev/null 2>&1; then
        echo "Nice! Your $(clipper -v) was installed successfully!"
    else
        echo "Error: Installation failed."
        exit 1
    fi
}

# Main script execution
set_binary_url
download_binary
make_executable
install_binary
verify_installation

# Clean up temporary directory
rm -rf "$TMP_DIR"
