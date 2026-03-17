//go:build linux && amd64

package barretenberg

// #cgo LDFLAGS: -L${SRCDIR}/../lib/linux_amd64 -lbarretenberg -lstdc++ -lm -lpthread
import "C"
