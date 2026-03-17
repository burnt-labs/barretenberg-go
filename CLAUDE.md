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
  link_*.go            Platform-specific CGo linker flags
  doc.go               Package documentation
  *_test.go            Tests
  testdata/statics/    Binary test vectors (vk, proof, public_inputs)
wrapper/               C++ wrapper shim (barretenberg_wrapper.cpp)
include/               C header (barretenberg_wrapper.h)
lib/<platform>/        Static archives (gitignored; built locally or downloaded from GitHub Releases)
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

- **Release assets** — Pre-built static archives are uploaded as GitHub Release assets, not committed to the repo. Consumers download the right platform archive at build time. Run `make build` for local development.
- **CGo paths** — `${SRCDIR}` in link files resolves to `barretenberg/`, so paths to `lib/` and `include/` use `../` prefix.
- **Platform-specific C++ stdlib** — All platforms link `libc++`. This is set in both the link_*.go files and build-wrapper.sh.
- **Debug symbol stripping** — Build script strips debug symbols to reduce archive size (~544MB → ~48MB on darwin).
- **Checksums** — All downloads are SHA256-verified against checksums.json. Update this file when bumping the Aztec version.

## CI

`.github/workflows/release.yml` builds all 4 platforms on push to main (when source files change) or on workflow_dispatch. The `release` job uses `go-semantic-release` to determine the next semver from conventional commit messages, creates a GitHub Release, and uploads the 4 `libbarretenberg_<platform>.a` archives as release assets.

Conventional commit prefixes: `feat:` = minor bump, `fix:` = patch bump, `feat!:` or `BREAKING CHANGE:` = major bump.

## Testing

Tests use binary test vectors in `barretenberg/testdata/statics/`. These were generated with Noir + bb CLI at Aztec v4.0.4. See `barretenberg/testdata/README.md` for regeneration instructions.

Run tests: `go test -v -count=1 ./barretenberg/`

## Upstream dependency

This wraps Aztec's barretenberg library from [aztec-packages](https://github.com/AztecProtocol/aztec-packages). The pinned version is in `checksums.json` under `aztec_tag`. To upgrade:

1. Update `aztec_tag` in checksums.json
2. Update SHA256 checksums for all 4 platform tarballs
3. Rebuild: `make build-all`
4. Run tests
5. Push to main (CI will create a release with the new archives)

## CRS (Common Reference String)

Proof verification requires BN254 SRS (Structured Reference String) data. The C++ wrapper calls `bb::srs::init_net_crs_factory()` inside `bb_verify_proof`, which sets up a factory that fetches SRS points from the Aztec CDN (`crs.aztec-cdn.foundation`) via HTTP range requests on demand.

- **Storage path** — Controlled by `BB_CRS_PATH` env var, falls back to `~/.bb-crs`. Only a 64-byte seed file (`bn254_g1.dat`) is written to disk; actual SRS points are fetched into memory.
- **Download size** — Depends on circuit size. Larger circuits need more SRS points. The download happens on the first verification after process start (~0.4s for the test circuit).
- **No pre-init needed** — The factory is initialized per-call inside `bb_verify_proof`. There is no separate `InitCRS` step required.
- **Validator considerations** — For validators, the network fetch adds latency only on the first verification after process start. Subsequent verifications reuse in-memory SRS data. Set `BB_CRS_PATH` to a persistent directory to cache the seed file across restarts.

## Consumer

Primary consumer is [xion](https://github.com/burnt-labs/xion) (`x/zk/` module) which imports this package for on-chain UltraHonk proof verification.
