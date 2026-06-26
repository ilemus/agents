# Build Cross-Platform Binaries for Desktop Applications

Go makes cross-compilation incredibly easy using environment variables (GOOS and GOARCH). You don't need a Mac to compile the Mac binary, and you don't need Linux to compile the Linux binary.

Run these commands from your project root:
```bash
# 1. Build for Fedora Linux (64-bit Intel/AMD architecture)
GOOS=linux GOARCH=amd64 go build -o build/linux-amd64/gopm main.go

# 2. Build for macOS Apple Silicon (M1/M2/M3/M4 ARM architecture)
GOOS=darwin GOARCH=arm64 go build -o build/macos-arm64/gopm main.go
```

## Makefile
This is an attempt to format the building of the cross-platform executables.

```makefile
BINARY_NAME=gopm
VERSION=1.0.0
BUILD_DIR=build

clean:
	rm -rf $(BUILD_DIR)

build-all: clean
	@echo "Building binaries..."
	# Linux AMD64
	mkdir -p $(BUILD_DIR)/linux-amd64
	GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/linux-amd64/$(BINARY_NAME) main.go
	cp -r reference_data.json $(BUILD_DIR)/linux-amd64/
	
	# macOS ARM64 (Apple Silicon)
	mkdir -p $(BUILD_DIR)/macos-arm64
	GOOS=darwin GOARCH=arm64 go build -o $(BUILD_DIR)/macos-arm64/$(BINARY_NAME) main.go
	cp -r reference_data.json $(BUILD_DIR)/macos-arm64/

package: build-all
	@echo "Packaging into archives..."
	# Package Linux as a .tar.gz (standard for Linux utilities)
	tar -czf $(BUILD_DIR)/$(BINARY_NAME)-v$(VERSION)-linux-amd64.tar.gz -C $(BUILD_DIR)/linux-amd64 .
	
	# Package macOS as a .zip
	cd $(BUILD_DIR)/macos-arm64 && zip -r ../$(BINARY_NAME)-v$(VERSION)-macos-arm64.zip .
	@echo "Done! Check the $(BUILD_DIR) folder."
```