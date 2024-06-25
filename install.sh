#!/bin/sh

# Determine OS and Architecture
OS=$(uname -s)
ARCH=$(uname -m)

# Check if necessary tools are available
check_commands() {
    for cmd in "$@"; do
        if ! command -v "$cmd" >/dev/null 2>&1; then
            echo "Error: $cmd is required but not installed."
            exit 1
        fi
    done
}

check_commands curl sudo

# Create a temporary directory for downloading the binary
TMP_DIR=$(mktemp -d)
TMP_BINARY="$TMP_DIR/clipper"

# Function to get the latest release version from GitHub
get_latest_version() {
    LATEST_VERSION=$(curl -s https://api.github.com/repos/supitsdu/clipper/releases/latest | grep '"tag_name":' | sed -E 's/.*"v([^"]+)".*/\1/')
    if [ -z "$LATEST_VERSION" ]; then
        echo "Error: Failed to fetch the latest version."
        exit 1
    fi
}

# Function to check if Clipper is installed and get its version
get_installed_version() {
    if ! command -v clipper >/dev/null 2>&1; then
        INSTALLED_VERSION="" # Not installed
    else
        INSTALLED_VERSION=$(clipper -v)
    fi
}

# Function to compare versions
is_newer_version() {
    [ "$(printf "%s\n%s\n" "$1" "$2" | sort -V | sed -n '2p')" != "$2" ]
}

# Function to set the download URL based on OS and Architecture
set_binary_url() {
    if [ -z "$LATEST_VERSION" ]; then
        echo "Error: Latest version not fetched yet. Cannot set binary URL."
        exit 1
    fi

    case "$OS" in
    "Linux")
        case "$ARCH" in
        "x86_64")
            BINARY_URL="https://github.com/supitsdu/clipper/releases/download/v$LATEST_VERSION/clipper_linux_amd64_v$LATEST_VERSION"
            ;;
        "aarch64" | "arm64")
            BINARY_URL="https://github.com/supitsdu/clipper/releases/download/v$LATEST_VERSION/clipper_linux_arm64_v$LATEST_VERSION"
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
            BINARY_URL="https://github.com/supitsdu/clipper/releases/download/v$LATEST_VERSION/clipper_darwin_amd64_v$LATEST_VERSION"
            ;;
        "arm64")
            BINARY_URL="https://github.com/supitsdu/clipper/releases/download/v$LATEST_VERSION/clipper_darwin_arm64_v$LATEST_VERSION"
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
    echo "Downloading latest Clipper release from GitHub..."
    if ! curl -Lo "$TMP_BINARY" "$BINARY_URL"; then
        echo "Error: Failed to download Clipper."
        cleanup
        exit 1
    fi
    echo "Download completed."
}

# Function to make the binary executable
make_executable() {
    if ! chmod +x "$TMP_BINARY"; then
        echo "Error: Failed to make Clipper executable."
        cleanup
        exit 1
    fi
}

# Function to move the binary to /usr/local/bin
install_binary() {
    if [ -d "/usr/local/bin" ] && echo "$PATH" | grep -q "/usr/local/bin"; then
        echo "Installing Clipper to /usr/local/bin..."
        if ! sudo mv -f "$TMP_BINARY" /usr/local/bin/clipper; then
            echo "Error: Failed to move Clipper to /usr/local/bin."
            cleanup
            exit 1
        fi
    else
        echo "Error: Directory /usr/local/bin does not exist or is not in your PATH."
        cleanup
        exit 1
    fi
}

# Function to verify the installation
verify_installation() {
    if command -v clipper >/dev/null 2>&1; then
        echo "Nice! Your $(clipper -v) was installed successfully!"
    else
        echo "Error: Installation failed."
        cleanup
        exit 1
    fi
}

# Function to prompt user for confirmation
user_input() {
    echo "$1"
    while true; do
        printf "Do you want to %s? (Y/n): " "$2"
        read -r confirm
        case "$confirm" in
        [Yy]* | '' | ' ') break ;;
        [Nn]*) exit 0 ;;
        *) echo "Invalid input. Please enter 'y' or 'n'." ;;
        esac
    done
}

# Function to uninstall Clipper
uninstall() {
    if command -v clipper >/dev/null 2>&1; then
        echo "Uninstalling $(clipper -v)..."
        if ! sudo rm -vf "$(command -v clipper)"; then
            echo "Failed to uninstall."
            exit 1
        fi
        cleanup
        echo "Done."
        exit 0
    else
        echo "Clipper is not installed. Exiting..."
        exit 0
    fi
}

# Function to log usage message
log_usage_message() {
    echo "Usage:"
    echo " $1 [options]"
    echo "Options:"
    echo "    -y  attempts to continue the installation without prompts."
    echo "    -h  display this usage message."
    echo "    -r  uninstall Clipper."
    echo "    -v  display current and latest version of Clipper."
}

# Function to clean up temporary files
cleanup() {
    rm -rf "$TMP_DIR" || {
        echo "Failed to cleanup temporary files: $?"
        exit 1
    }
}

log_version_message() {
    get_latest_version
    get_installed_version

    if [ -n "$INSTALLED_VERSION" ] && [ -n "$LATEST_VERSION" ]; then
        echo "$INSTALLED_VERSION is currently installed. (Latest: $LATEST_VERSION)"
        exit 0
    fi

    if [ -n "$INSTALLED_VERSION" ]; then
        echo "$INSTALLED_VERSION is currently installed."
        exit 0
    fi

    if [ -n "$LATEST_VERSION" ]; then
        echo "Clipper $LATEST_VERSION is the latest version. To install, run this script: sh $0 -y"
        exit 0
    fi

    echo "Failed to get the version of Clipper. See https://github.com/supitsdu/clipper/releases"
    exit 1
}

# Main script execution
while getopts "yhrv" opt; do
    case $opt in
    y) AUTO_CONFIRM=true ;;
    r) uninstall ;;
    h) log_usage_message "$0" && exit 0 ;;
    v) log_version_message ;;
    *) log_usage_message "$0" && exit 1 ;;
    esac
done

[ -n "$AUTO_CONFIRM" ] || user_input "This script will install Clipper on your system." "proceed"

get_latest_version
get_installed_version

if [ -n "$INSTALLED_VERSION" ]; then
    echo "$INSTALLED_VERSION is already installed."
    if [ -z "$AUTO_CONFIRM" ]; then
        if is_newer_version "$LATEST_VERSION" "$INSTALLED_VERSION"; then
            user_input "A newer version ($LATEST_VERSION) is available." "upgrade"
        else
            user_input "You have the latest version installed." "reinstall"
        fi
    fi
fi

set_binary_url
download_binary
make_executable
install_binary
verify_installation

# Clean up temporary directory
cleanup
