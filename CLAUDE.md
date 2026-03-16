# CLAUDE.md — barretenberg-go

Go bindings for Aztec Barretenberg's UltraHonk proof verification system.

## Project structure

```
barretenberg/          Go package (import "github.com/burnt-labs/barretenberg-go/barretenberg")
  bindings.go          CGo bridge to C wrapper (low-level)
  verifier.go          High-level Verifier API
  vkey.go              VerificationKey type
  proof.go             Proof and PublicInputs types
  errors.go            Sentinel errors
  crs.go               Common Reference String (CRS/SRS) initialization
  link_*.go            Platform-specific CGo linker flags
  doc.go               Package documentation
  *_test.go            Tests
  testdata/statics/    Binary test vectors (vk, proof, public_inputs)
wrapper/               C++ wrapper shim (barretenberg_wrapper.cpp)
include/               C header (barretenberg_wrapper.h)
lib/<platform>/        Pre-built static archives (committed, not LFS)
scripts/               Build script (build-wrapper.sh)
checksums.json         Aztec release version, SHA256 checksums
```

## Build and test

```bash
make build          # Build libbarretenberg.a for current platform
make test           # Run Go tests (requires lib/<platform>/libbarretenberg.a)
make bench          # Run benchmarks
make build-all      # Cross-compile all 4 platforms
```

The build script (`scripts/build-wrapper.sh`) downloads Aztec's pre-built `libbb-external.a`, compiles the C++ wrapper shim, and merges them into `libbarretenberg.a`. Version and checksums come from `checksums.json`.

## Key conventions

- **No Git LFS** — Go module proxy can't resolve LFS pointers. Archives are committed directly (stripped to stay under GitHub's 100MB limit).
- **CGo paths** — `${SRCDIR}` in link files resolves to `barretenberg/`, so paths to `lib/` and `include/` use `../` prefix.
- **Platform-specific C++ stdlib** — linux_amd64 links `libstdc++` (matching Aztec's build), all others link `libc++`. This is set in both the link_*.go files and build-wrapper.sh.
- **Debug symbol stripping** — Build script strips debug symbols to reduce archive size (~544MB → ~48MB on darwin).
- **Checksums** — All downloads are SHA256-verified against checksums.json. Update this file when bumping the Aztec version.

## CI

`.github/workflows/release.yml` builds all 4 platforms on push of a semver tag (`v*.*.*`). The `commit-archives` job commits built archives to main. Note: `burnt-labs` org requires signed commits, so the CI bot commit step may fail — download artifacts and commit locally if needed.

## Testing

Tests use binary test vectors in `barretenberg/testdata/statics/`. These were generated with Noir + bb CLI at Aztec v4.0.4. See `barretenberg/testdata/README.md` for regeneration instructions.

Run tests: `go test -v -count=1 ./barretenberg/`

## Upstream dependency

This wraps Aztec's barretenberg library from [aztec-packages](https://github.com/AztecProtocol/aztec-packages). The pinned version is in `checksums.json` under `aztec_tag`. To upgrade:

1. Update `aztec_tag` in checksums.json
2. Update SHA256 checksums for all 4 platform tarballs
3. Rebuild: `make build-all`
4. Run tests
5. Tag a new semver release

## Consumer

Primary consumer is [xion](https://github.com/burnt-labs/xion) (`x/zk/` module) which imports this package for on-chain UltraHonk proof verification.
