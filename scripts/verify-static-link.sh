#!/usr/bin/env bash
#
# verify-static-link.sh — link the barretenberg test binary the way a consumer
# does (Zig musl toolchain, muslc tag, external static link) and check that the
# result is genuinely static.
#
# Usage:
#   ./scripts/verify-static-link.sh <amd64|arm64> <output-path>
#
# Requires the pinned Zig toolchain wrappers on PATH (scripts/install-zig.sh)
# and lib/linux_<arch>_musl/libbarretenberg.a.
#
# Fails when:
#   - the linker warns that a symbol needs shared libraries at runtime
#     ("statically linked applications requires ...")
#   - the binary has an ELF interpreter or dynamic dependencies
#   - the requested default stack (PT_GNU_STACK) is below 8 MiB; musl sizes new
#     threads from it and its built-in default is only 128 KiB
#
# The binary is executed when the host architecture matches.

set -euo pipefail

if [[ $# -ne 2 ]]; then
    echo "Usage: $0 <amd64|arm64> <output-path>" >&2
    exit 1
fi

ARCH="$1"
OUT="$(cd "$(dirname "$2")" && pwd)/$(basename "$2")"
case "$ARCH" in
    amd64) ZIG_TRIPLE="x86_64-linux-musl"; HOST_MACHINE="x86_64" ;;
    arm64) ZIG_TRIPLE="aarch64-linux-musl"; HOST_MACHINE="aarch64" ;;
    *) echo "Unsupported architecture: $ARCH" >&2; exit 1 ;;
esac

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
LIB="$REPO_ROOT/lib/linux_${ARCH}_musl/libbarretenberg.a"
MIN_STACK=$((8 * 1024 * 1024))

if [[ ! -f "$LIB" ]] || head -c 64 "$LIB" | grep -q 'git-lfs'; then
    echo "$LIB is missing or is an LFS pointer" >&2
    exit 1
fi

LINK_LOG="$(mktemp)"
RUN_LOG=""
trap 'rm -f "$LINK_LOG" ${RUN_LOG:+"$RUN_LOG"}' EXIT

(
    cd "$REPO_ROOT"
    GOOS=linux GOARCH="$ARCH" CGO_ENABLED=1 \
        CC="${ZIG_TRIPLE}-zig-cc" \
        CXX="${ZIG_TRIPLE}-zig-c++" \
        go test -c -tags muslc \
        -ldflags "-linkmode=external -extldflags '-static -lm'" \
        -o "$OUT" ./barretenberg
) 2>&1 | tee "$LINK_LOG"

if grep -q 'statically linked applications requires' "$LINK_LOG"; then
    echo "linux/$ARCH: the static link still needs shared libraries at runtime" >&2
    exit 1
fi
if readelf -lW "$OUT" | grep -q INTERP; then
    echo "linux/$ARCH: test binary unexpectedly has an ELF interpreter" >&2
    exit 1
fi
if readelf -dW "$OUT" 2>&1 | grep -q NEEDED; then
    echo "linux/$ARCH: test binary unexpectedly has dynamic dependencies" >&2
    exit 1
fi

STACK_HEX="$(readelf -lW "$OUT" | awk '$1 == "GNU_STACK" { print $6; exit }')"
if [[ -z "$STACK_HEX" ]] || (( STACK_HEX < MIN_STACK )); then
    echo "linux/$ARCH: PT_GNU_STACK size ${STACK_HEX:-unset} is below 8 MiB; musl threads would get a small stack" >&2
    exit 1
fi
echo "linux/$ARCH: static, no interpreter, no dynamic deps, stack request $((STACK_HEX / 1024 / 1024)) MiB"

if [[ "$(uname -m)" == "$HOST_MACHINE" ]]; then
    # Run from the package directory so the testdata vectors resolve; a run
    # that skips them proves nothing.
    RUN_LOG="$(mktemp)"
    (cd "$REPO_ROOT/barretenberg" && "$OUT" -test.v -test.count=1) 2>&1 | tee "$RUN_LOG"
    if grep -q -- '--- SKIP' "$RUN_LOG"; then
        echo "linux/$ARCH: tests were skipped; expected the full suite to run" >&2
        exit 1
    fi
else
    echo "linux/$ARCH: host is $(uname -m); link verified, execution skipped"
fi
