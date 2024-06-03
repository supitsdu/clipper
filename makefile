# Define the output directory for the binaries
OUT_DIR := bin

# Retrieve the version from git tags
VERSION := $(shell git describe --tags --always)

# Define the names for the binaries with version
WINDOWS_BIN := clipper_windows_amd64_$(VERSION).exe
LINUX_BIN := clipper_linux_amd64_$(VERSION)
DARWIN_BIN := clipper_darwin_amd64_$(VERSION)

# Define the build targets for each platform
.PHONY: all windows linux darwin clean checksums test help

# Default target: build binaries for all platforms
all: windows linux darwin

# Build binary for Windows
windows: $(OUT_DIR)
	GOOS=windows GOARCH=amd64 go build -ldflags="-X main.version=$(VERSION)" -o $(OUT_DIR)/$(WINDOWS_BIN)

# Build binary for Linux
linux: $(OUT_DIR)
	GOOS=linux GOARCH=amd64 go build -ldflags="-X main.version=$(VERSION)" -o $(OUT_DIR)/$(LINUX_BIN)

# Build binary for macOS
darwin: $(OUT_DIR)
	GOOS=darwin GOARCH=amd64 go build -ldflags="-X main.version=$(VERSION)" -o $(OUT_DIR)/$(DARWIN_BIN)

# Generate SHA256 checksums for each binary
checksums: $(OUT_DIR)
	@echo "Generating SHA256 checksums..."
	@for binary in $(WINDOWS_BIN) $(LINUX_BIN) $(DARWIN_BIN); do \
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
	@echo "  all        - Build binaries for all platforms"
	@echo "  windows    - Build binary for Windows"
	@echo "  linux      - Build binary for Linux"
	@echo "  darwin     - Build binary for macOS"
	@echo "  checksums  - Generate SHA256 checksums for binaries"
	@echo "  test       - Run tests"
	@echo "  clean      - Clean the built binaries"
	@echo "  help       - Show this help message"
