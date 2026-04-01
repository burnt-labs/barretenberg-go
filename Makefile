.PHONY: build test clean bench

PLATFORM := $(shell go env GOOS)_$(shell go env GOARCH)

# Build libbarretenberg.a for the current platform
build:
	./scripts/build-wrapper.sh --platform $(PLATFORM)

# Build for a specific platform (e.g., make build-linux_amd64)
build-%:
	./scripts/build-wrapper.sh --platform $*

# Build for all platforms (requires cross-compilation toolchains)
build-all: build-linux_amd64 build-linux_arm64 build-darwin_amd64 build-darwin_arm64

# Run tests (requires lib/<platform>/libbarretenberg.a to exist)
test:
	go test -v -count=1 ./...

# Run benchmarks
bench:
	go test -bench=. -benchmem -run=^$$ ./...

# Clean temporary build artifacts (does not remove committed lib/*.a files)
clean:
	rm -rf /tmp/bb-build-*
