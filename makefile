# Makefile for building Clipper

# Define the output directory for the binaries
OUT_DIR := bin

# Define version info (set dynamically during build)
VERSION ?= $(shell git describe --tags --always)  # Get latest tag or commit hash as version

REPO_URL := github.com/supitsdu/clipper

# Define the build targets for each platform
.PHONY: all windows linux linux_arm linux_arm64 darwin darwin_arm64 clean checksums test help

# Default target: build binaries for all platforms
all: windows linux linux_arm linux_arm64 darwin darwin_arm64

# Generic build function with simplified binary name and embedded metadata
define build
GOOS=$(1) GOARCH=$(2) go build \
	-ldflags="-X '$(REPO_URL)/cli/options.Version=$(VERSION)' -X '$(REPO_URL)/cli/options.BuildMetadata=$(1)/$(2)'" \
	-o $(OUT_DIR)/clipper_$(1)_$(2)_$(VERSION) ./cli
endef

# Build binaries for each platform, calling the generic build function with appropriate arguments
windows: $(OUT_DIR)
	$(call build,windows,amd64,windows)

# Build binary for Linux (amd64)
linux: $(OUT_DIR)
	$(call build,linux,amd64,linux)

# Build binary for Linux (arm)
linux_arm: $(OUT_DIR)
	$(call build,linux,arm,linux)

# Build binary for Linux (arm64)
linux_arm64: $(OUT_DIR)
	$(call build,linux,arm64,linux)

# Build binary for macOS (amd64)
darwin: $(OUT_DIR)
	$(call build,darwin,amd64,darwin)

# Build binary for macOS (arm64)
darwin_arm64: $(OUT_DIR)
	$(call build,darwin,arm64,darwin)

# Generate SHA256 checksums for each binary
checksums: $(OUT_DIR)
	@echo "Generating SHA256 checksums..."
	@for binary in $(WINDOWS_BIN) $(LINUX_BIN) $(LINUX_ARM_BIN) $(LINUX_ARM64_BIN) $(DARWIN_BIN) $(DARWIN_ARM64_BIN); do \
		sha256sum $(OUT_DIR)/$$binary > $(OUT_DIR)/$$binary.sha256; \
	done
	@echo "Checksum files generated successfully."

# Run tests
test:
	go test ./...

# Clean the built binaries
clean:
	rm -rf $(OUT_DIR)

# Create the output directory if it doesn't exist
$(OUT_DIR):
	mkdir -p $(OUT_DIR)

# Show help message
help:
	@echo "Makefile targets:"
	@echo "  all          - Build binaries for all platforms"
	@echo "  windows      - Build binary for Windows (amd64)"
	@echo "  linux        - Build binary for Linux (amd64)"
	@echo "  linux_arm    - Build binary for Linux (arm)"
	@echo "  linux_arm64  - Build binary for Linux (arm64)"
	@echo "  darwin       - Build binary for macOS (amd64)"
	@echo "  darwin_arm64 - Build binary for macOS (arm64)"
	@echo "  checksums    - Generate SHA256 checksums for binaries"
	@echo "  test         - Run tests"
	@echo "  clean        - Clean the built binaries"
	@echo "  version      - Show the latest git tag version"
	@echo "  help         - Show this help message"

# Show the latest git tag version
version:
	@echo $(VERSION)
