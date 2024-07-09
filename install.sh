#!/bin/sh
# The install script for Clipper made by @supitsdu
# Any copyright is dedicated to the Public Domain.
# https://creativecommons.org/publicdomain/zero/1.0/

# --- Constants & Configuration ---
BINARY_DIR="/usr/local/bin"
GITHUB_REPO="supitsdu/clipper"

# ANSI Escape Codes for Colors and Styling
COLOR_RESET="\033[0m"
COLOR_DARK="\033[2m"
COLOR_RED="\033[31m"
COLOR_GREEN="\033[32m"
COLOR_BLUE="\033[34m"
COLOR_YELLOW="\033[33m"
STYLE_BOLD="\033[1m"

# --- Helper Functions ---

# Print colored/styled messages
log() {
    color="$1"
    shift
    printf "%b%b${COLOR_RESET}\n" "$color" "$@"
}

# Get the latest release version from GitHub
get_latest_version() {
    case "$GCMD" in
    "curl")
        curl -s "https://api.github.com/repos/${GITHUB_REPO}/releases/latest" |
            grep '"tag_name":' | sed -E 's/.*"v([^"]+)".*/\1/'
        ;;
    "wget")
        wget -qO- "https://api.github.com/repos/${GITHUB_REPO}/releases/latest" |
            grep '"tag_name":' | sed -E 's/.*"v([^"]+)".*/\1/'
        ;;
    esac
}

# Check if Clipper is installed and get its version
get_installed_version() {
    clipper -v 2>/dev/null | sed -E 's/.*\sv([^"]+)\s.*/\1/' || true # Silently handle 'command not found'
}

# Compare semantic versions (v1.2.3)
is_newer_version() {
    latest="$1"
    installed="$2"
    [ "$(printf "%s\n%s\n" "$latest" "$installed" | sort -V | sed -n '2p')" != "$installed" ]
}

# Set the download URL based on OS and Architecture
set_binary_url() {
    if [ -z "$latest_version" ]; then
        log "$COLOR_RED" "Failed to fetch the latest version. Cannot set binary URL."
        exit 1
    fi

    os=$(uname -s)
    arch=$(uname -m)

    case "$os" in
    "Linux")
        case "$arch" in
        "x86_64")
            echo "https://github.com/${GITHUB_REPO}/releases/download/v$latest_version/clipper_linux_amd64_v$latest_version"
            ;;
        "arm")
            echo "https://github.com/${GITHUB_REPO}/releases/download/v$latest_version/clipper_linux_arm_v$latest_version"
            ;;
        "aarch64" | "arm64")
            echo "https://github.com/${GITHUB_REPO}/releases/download/v$latest_version/clipper_linux_arm64_v$latest_version"
            ;;
        *)
            log "$COLOR_RED" "Unsupported architecture: $arch"
            exit 1
            ;;
        esac
        ;;
    "Darwin")
        case "$arch" in
        "x86_64")
            echo "https://github.com/${GITHUB_REPO}/releases/download/v$latest_version/clipper_darwin_amd64_v$latest_version"
            ;;
        "arm64")
            echo "https://github.com/${GITHUB_REPO}/releases/download/v$latest_version/clipper_darwin_arm64_v$latest_version"
            ;;
        *)
            log "$COLOR_RED" "Unsupported architecture: $arch"
            exit 1
            ;;
        esac
        ;;
    *)
        log "$COLOR_RED" "Unsupported OS: $os"
        exit 1
        ;;
    esac
}

# Download the binary to a temporary location
download_binary() {
    url=$1
    case "$GCMD" in
    "curl")
        curl -#L --fail -o "$TMP_BINARY" "$url" || {
            cleanup
            log "$COLOR_RED" "Failed to download Clipper."
            exit 1
        }
        ;;
    "wget")
        wget --show-progress -qO "$TMP_BINARY" "$url" || {
            cleanup
            log "$COLOR_RED" "Failed to download Clipper."
            exit 1
        }
        ;;
    esac
}

# Install or upgrade the binary
install_binary() {
    if ! (sudo mv -f "$TMP_BINARY" "$BINARY_DIR/clipper" && sudo chmod +x "$BINARY_DIR/clipper"); then
        cleanup
        log "$COLOR_RED" "Failed to install Clipper."
        exit 1
    fi
}

# Prompt the user for confirmation (unless -y is used)
confirm_action() {
    if [ -z "$AUTO_CONFIRM" ]; then
        printf "Do you want to proceed? [y/N] "
        read -r REPLY
        case "$REPLY" in
        [Yy]*) ;;
        *) exit 0 ;;
        esac
        echo # Adds new line for UX purpose only (Temporary)
    fi
}

# Remove or clean up binaries
uninstall() {
    if ! command -v clipper >/dev/null 2>&1; then
        log "${STYLE_BOLD}" "Clipper is not installed."
        log "${COLOR_BLUE}" "Run this script with the '-y' option to install Clipper."
        exit 0
    fi

    log "${COLOR_YELLOW}" "Uninstalling $(clipper -v)..."

    if ! sudo rm -f "$(command -v clipper)"; then
        log "$COLOR_RED" "Failed to uninstall $(clipper -v)."
        exit 1
    fi

    cleanup # Clean up any residual temporary files

    log "$COLOR_BLUE" "Clipper has been successfully uninstalled."
}

# Cleanup temporary files
cleanup() {
    rm -rf "$TMP_DIR"
}

# Check dependencies (curl or wget)
check_deps() {
    getter_cmd=""

    for tool in "$@"; do
        [ -n "$getter_cmd" ] && continue
        if command -v "$tool" >/dev/null 2>&1; then
            [ "$GCMD" != "$tool" ] && [ -n "$GCMD" ] && continue
            getter_cmd="$tool"
        fi
    done

    if [ -z "$getter_cmd" ]; then
        cleanup

        if [ -n "$GCMD" ]; then
            log "$COLOR_RED" "The selected '$GCMD' is not installed on your system."
            log "$COLOR_RED" "Please install '$GCMD' or choose another supported option."
        else
            log "$COLOR_RED" "Neither 'curl' nor 'wget' is installed. Please install one of them to proceed."
        fi

        exit 1
    fi

    if [ -z "$GCMD" ]; then
        GCMD="$getter_cmd"
    fi
}

# --- Main Script Logic ---

# Parse command-line options for -y (auto confirm), -u (uninstall), -g (getter command), and -h (help)
while getopts ":yhug:" opt; do
    case "$opt" in
    u)
        uninstall
        exit 0
        ;;
    y)
        AUTO_CONFIRM=true
        ;;
    g)
        shift
        case "$OPTARG" in
        "curl") GCMD="curl" ;;
        "wget") GCMD="wget" ;;
        *)
            cleanup
            log "$COLOR_RED" "Unsupported command: '$OPTARG'. Please use 'curl' or 'wget'."
            exit 1
            ;;
        esac
        ;;
    h)
        log "${COLOR_GREEN}" "\nThis shell script installs or updates Clipper.\n"
        log "${STYLE_BOLD}" "Usage:${COLOR_RESET}"
        log "${COLOR_DARK}" "\n $0 ${COLOR_RESET}${COLOR_GREEN}<option> [arguments]${COLOR_RESET}\n"
        log "${STYLE_BOLD}" "Options:${COLOR_RESET}"
        log "${COLOR_GREEN}" "\n   -y${COLOR_RESET}\tAutomatically confirm prompts"
        log "${COLOR_GREEN}" "\n   -u${COLOR_RESET}\tUninstall Clipper"
        log "${COLOR_GREEN}" "\n   -g${COLOR_RESET}\tSpecify getter command (e.g., 'curl' or 'wget')"
        log "${COLOR_GREEN}" "\n   -h${COLOR_RESET}\tShow this help message"
        log "${COLOR_GREEN}" "\nClipper is a lightweight command-line tool for copying contents to the clipboard.\n"
        exit 0
        ;;
    \?) log "$COLOR_RED" "Invalid option: -$OPTARG" && exit 1 ;;
    :) log "$COLOR_RED" "Option -$OPTARG requires an argument." && exit 1 ;;
    esac
    shift
done

# Check dependencies (curl or wget)
check_deps "curl" "wget"

# Get latest and installed versions
latest_version="$(get_latest_version)"
installed_version="$(get_installed_version)"

# Check for Clipper and handle installation/update accordingly
if [ -z "$installed_version" ]; then
    log "$COLOR_YELLOW" "Installing Clipper...\n"
    confirm_action
else
    if is_newer_version "$latest_version" "$installed_version" || [ "$AUTO_CONFIRM" = "true" ]; then
        log "$COLOR_YELLOW" "Upgrading Clipper...\n"
        confirm_action
    else
        log "$COLOR_BLUE" "Clipper is already up-to-date (v$installed_version)."
        log "$COLOR_BLUE" "\n\t✦ No action needed. You're running the latest version! ✦"
        log "$COLOR_BLUE" "\nIf you'd like to reinstall Clipper, run this script with the '-y' option."
        exit 0
    fi
fi

# Create a temporary directory for the download
TMP_DIR="$(mktemp -d --suffix g.clipper)"
TMP_BINARY="$TMP_DIR/clipper"

# Download and install the binary
download_binary "$(set_binary_url)"
install_binary

# Display success messages
if is_newer_version "$latest_version" "$installed_version" || [ -n "$installed_version" ]; then
    log "${STYLE_BOLD}${COLOR_BLUE}" "\nClipper v$latest_version has been successfully installed!"
    log "${STYLE_BOLD}${COLOR_BLUE}" "\n To get started with Clipper:"
    log "${STYLE_BOLD}${COLOR_BLUE}" "\n  1. Ensure '/usr/local/bin' is in your PATH environment variable."
    log "${STYLE_BOLD}${COLOR_BLUE}" "\n  2. Type 'clipper --help' to see available commands and options."
else
    log "${STYLE_BOLD}${COLOR_BLUE}" "\nClipper has been updated to v$latest_version."
fi

log "${STYLE_BOLD}${COLOR_BLUE}" "\nFor detailed documentation and examples, visit the project repository:"
log "${STYLE_BOLD}${COLOR_BLUE}" "https://github.com/$GITHUB_REPO"
log "${STYLE_BOLD}${COLOR_BLUE}" "\nClipper is open source and licensed under the MIT License."

# Cleanup temporary files
cleanup
