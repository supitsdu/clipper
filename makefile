# Define the output directory for the binaries
OUT_DIR := bin

# Retrieve the version from git tags
VERSION := $(shell git describe --tags --always)

# Define the names for the binaries with version
WINDOWS_BIN := clipper_windows_amd64_$(VERSION).exe
LINUX_BIN := clipper_linux_amd64_$(VERSION)
LINUX_ARM_BIN := clipper_linux_arm_$(VERSION)
LINUX_ARM64_BIN := clipper_linux_arm64_$(VERSION)
DARWIN_BIN := clipper_darwin_amd64_$(VERSION)
DARWIN_ARM64_BIN := clipper_darwin_arm64_$(VERSION)

# Define the build targets for each platform
.PHONY: all windows linux linux_arm linux_arm64 darwin darwin_arm64 clean checksums test help

# Default target: build binaries for all platforms
all: windows linux linux_arm linux_arm64 darwin darwin_arm64

# Build binary for Windows
windows: $(OUT_DIR)
	GOOS=windows GOARCH=amd64 go build -ldflags="-X main.version=$(VERSION)" -o $(OUT_DIR)/$(WINDOWS_BIN)

# Build binary for Linux (amd64)
linux: $(OUT_DIR)
	GOOS=linux GOARCH=amd64 go build -ldflags="-X main.version=$(VERSION)" -o $(OUT_DIR)/$(LINUX_BIN)

# Build binary for Linux (arm)
linux_arm: $(OUT_DIR)
	GOOS=linux GOARCH=arm go build -ldflags="-X main.version=$(VERSION)" -o $(OUT_DIR)/$(LINUX_ARM_BIN)

# Build binary for Linux (arm64)
linux_arm64: $(OUT_DIR)
	GOOS=linux GOARCH=arm64 go build -ldflags="-X main.version=$(VERSION)" -o $(OUT_DIR)/$(LINUX_ARM64_BIN)

# Build binary for macOS (amd64)
darwin: $(OUT_DIR)
	GOOS=darwin GOARCH=amd64 go build -ldflags="-X main.version=$(VERSION)" -o $(OUT_DIR)/$(DARWIN_BIN)

# Build binary for macOS (arm64)
darwin_arm64: $(OUT_DIR)
	GOOS=darwin GOARCH=arm64 go build -ldflags="-X main.version=$(VERSION)" -o $(OUT_DIR)/$(DARWIN_ARM64_BIN)

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