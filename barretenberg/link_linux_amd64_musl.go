//go:build linux && amd64 && muslc

package barretenberg

// #cgo LDFLAGS: -L${SRCDIR}/../lib/linux_amd64_musl -lbarretenberg -lc++ -lm -lpthread
import "C"
