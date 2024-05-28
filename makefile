# Define the output directory for the binaries
OUT_DIR := bin

# Define the names for the binaries
WINDOWS_BIN := clipper_windows_amd64.exe
LINUX_BIN := clipper_linux_amd64
DARWIN_BIN := clipper_darwin_amd64

# Define the build targets for each platform
.PHONY: all windows linux darwin clean

# Default target: build binaries for all platforms
all: windows linux darwin

# Build binary for Windows
windows:
	GOOS=windows GOARCH=amd64 go build -o $(OUT_DIR)/$(WINDOWS_BIN)

# Build binary for Linux
linux:
	GOOS=linux GOARCH=amd64 go build -o $(OUT_DIR)/$(LINUX_BIN)

# Build binary for macOS
darwin:
	GOOS=darwin GOARCH=amd64 go build -o $(OUT_DIR)/$(DARWIN_BIN)

# Generate SHA256 checksums for each binary
checksums:
	@echo "Generating SHA256 checksums..."
	@for binary in $(WINDOWS_BIN) $(LINUX_BIN) $(DARWIN_BIN); do \
		sha256sum $(OUT_DIR)/$$binary > $(OUT_DIR)/$$binary.sha256; \
	done
	@echo "Checksum files generated successfully."

# Clean the built binaries
clean:
	rm -rf $(OUT_DIR)
