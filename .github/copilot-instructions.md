# Copilot Instructions — barretenberg-go

Go (CGo) bindings for Aztec Barretenberg's UltraHonk zero-knowledge proof verification. Single Go package, zero Go dependencies, ~15 source files. Primary consumer is [xion](https://github.com/burnt-labs/xion) (`x/zk/` module) for on-chain proof verification.

## Repository layout

```
barretenberg/              Go package — all Go source lives here
  bindings.go              CGo bridge (low-level C calls, vkeyHandle)
  verifier.go              High-level Verifier API (NewVerifier, Verify)
  vkey.go                  VerificationKey parsing and validation
  proof.go                 Proof and PublicInputs types
  errors.go                Sentinel errors and error codes (must match C enum)
  link_<os>_<arch>.go      Platform-specific CGo LDFLAGS (4 files)
  doc.go                   Package godoc
  *_test.go                Tests (3 files, ~900 lines total)
  testdata/statics/        Binary test vectors: vk, proof, public_inputs
wrapper/                   C++ wrapper shim (barretenberg_wrapper.cpp, 277 lines)
include/                   C header (barretenberg_wrapper.h) — defines bb_error_t enum and API
scripts/build-wrapper.sh   Build script: downloads Aztec .a, compiles wrapper, merges archives
lib/<os>_<arch>/           Static archives (gitignored; built by `make build` or downloaded from releases)
checksums.json             Aztec release version (aztec_tag) and SHA256 checksums for all 4 platform tarballs
.github/workflows/release.yml   CI: builds all 4 platforms, creates GitHub Release with assets
Makefile                   Convenience targets
go.mod                     Module: github.com/burnt-labs/barretenberg-go (go 1.25, zero dependencies, no go.sum)
```

## Build and test

**Prerequisite:** The static archive `lib/<platform>/libbarretenberg.a` must exist before any Go command (test, vet, build) will work. Without it, CGo linking fails immediately. Always run `make build` first on a fresh clone.

```bash
make build          # Build libbarretenberg.a for current OS/arch (~2-5 min, downloads ~50MB from Aztec)
make test           # Runs: go test -v -count=1 ./...
make bench          # Runs: go test -bench=. -benchmem -run=^$ ./...
make clean          # Removes /tmp/bb-build-* temp dirs (does NOT delete lib/)
```

**Build script prerequisites** (for `make build`): `clang++` with C++20, `python3`, `curl`, `ar`, `git`. The script respects `$CXX` and `$AR` env vars.

**Tests complete in <1 second.** All 29 tests pass. The `TestVerifyInvalidProof` test prints a `UltraVerifier: verification failed` line to stderr — this is expected, not an error.

**Cross-compilation:** `make build-<platform>` where platform is `linux_amd64`, `linux_arm64`, `darwin_amd64`, or `darwin_arm64`. Darwin amd64 cross-compiles from arm64 runners using `-target x86_64-apple-macos10.15 -isysroot $(xcrun --show-sdk-path)`.

**No linter config exists.** Standard `go vet ./...` works once the library is built.

## CI pipeline

`.github/workflows/release.yml` triggers on push to main (when source/build files change) or `workflow_dispatch`.

1. **build** job (4-platform matrix): builds archive, runs tests (skipped for cross-compiled darwin_amd64), uploads artifact
2. **release** job: downloads all 4 artifacts, runs `go-semantic-release` (dry) to determine version from conventional commits (`feat:` = minor, `fix:` = patch, `feat!:` / `BREAKING CHANGE:` = major), creates GitHub Release with `libbarretenberg_<os>_<arch>.a` assets

Runners: `ubuntu-latest` (linux_amd64), `ubuntu-24.04-arm` (linux_arm64), `macos-15` (darwin_amd64 cross-compile), `macos-latest` (darwin_arm64). The darwin_amd64 build requires macOS 15+ (Xcode 16) for full C++20 support in Aztec headers.

## Key facts for making changes

**Error codes must stay in sync across three files.** If adding or changing an error code, update all three:
1. `include/barretenberg_wrapper.h` — `bb_error_t` C enum
2. `wrapper/barretenberg_wrapper.cpp` — C++ implementation returning codes
3. `barretenberg/errors.go` — Go constants and `errorFromCode()` switch

**CGo link flags are per-platform.** `link_linux_amd64.go` uses `-lstdc++` (libstdc++); all other platforms use `-lc++` (libc++). This matches Aztec's pre-built archives (amd64 built with GCC, arm64 with Zig/clang). The divergence does not affect verification correctness — barretenberg's cryptographic operations (BN254 field arithmetic, polynomial commitments, Poseidon2 hashing) are custom implementations independent of C++ stdlib behavior.

**Thread safety pattern.** All public types (`Verifier`, `VerificationKey`, `vkeyHandle`) use `sync.RWMutex`. Reads take `RLock`, mutations and `Close()` take full `Lock`. Always check `isClosed` under the lock before operating on C pointers. Follow this pattern when adding new methods.

**Memory management.** Go-to-C calls use `runtime.Pinner` to pin byte slices during C calls (prevents GC from moving memory). C resources are freed via explicit `Close()` with `runtime.SetFinalizer` as a safety net.

**Field elements are 32 bytes, big-endian** (BN254 scalar field). Public inputs are concatenated 32-byte elements. The C wrapper validates `pub_len == num_inputs * 32`.

**CRS (Common Reference String).** G1 and G2 points are hardcoded in the C++ wrapper — no network download. Initialization happens once per process via `std::once_flag`.

**Test vectors** in `barretenberg/testdata/statics/` were generated with Aztec v4.0.4. Regeneration instructions are in `barretenberg/testdata/README.md`.

**Upstream version.** Pinned in `checksums.json` under `aztec_tag`. To upgrade: update `aztec_tag`, update all 4 SHA256 checksums, rebuild all platforms, run tests.

## Common pitfalls

- **`go test` fails with linker errors** → Run `make build` first. The static archive must exist.
- **Changing the C header or wrapper** → Rebuild the archive (`make build`) before testing. Go caches CGo artifacts aggressively; run `go clean -cache` if stale.
- **Adding a new public method to Verifier/VerificationKey** → Must acquire the appropriate mutex lock and check `isClosed` before accessing the underlying `vkeyHandle`.
- **Modifying build-wrapper.sh platform flags** → Keep `link_<os>_<arch>.go` LDFLAGS in sync. The build script's `EXTRA_LDFLAGS` variable documents what the link file should contain.

Trust these instructions. Only search the codebase if information here is incomplete or found to be incorrect.
