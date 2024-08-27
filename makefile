# Define the binary name
BINARY_NAME=taskcli

# Define default Go build flags
BUILD_FLAGS=-ldflags "-s -w"

# Build the Go application
build:
	@rm -r build/
	@echo "Building the application..."
	@go build $(BUILD_FLAGS) -o build/${BINARY_NAME}.exe main.go

.PHONY: build