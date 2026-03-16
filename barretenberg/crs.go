package barretenberg

/*
#cgo CFLAGS: -I${SRCDIR}/../include
#include "barretenberg_wrapper.h"
#include <stdlib.h>
*/
import "C"

import "unsafe"

// InitCRS initialises the Common Reference String (CRS) needed for proof
// verification. It must be called once before any Verify calls.
//
// path specifies the directory for CRS data files. If empty, the library
// checks the BB_CRS_PATH environment variable, then falls back to ~/.bb-crs.
//
// If download is true, any missing CRS data will be fetched from the Aztec
// CDN (crs.aztec-cdn.foundation) on first use. The download size depends on
// the circuit being verified and is performed via HTTP range requests.
//
// If download is false, the CRS files must already exist at the given path;
// verification will fail if they are missing or too small.
func InitCRS(path string, download bool) error {
	var cPath *C.char
	if path != "" {
		cPath = C.CString(path)
		defer C.free(unsafe.Pointer(cPath))
	}

	allowDownload := C.int(0)
	if download {
		allowDownload = 1
	}

	result := C.bb_init_crs(cPath, allowDownload)
	if result != C.BB_SUCCESS {
		return getLastError(int(result))
	}
	return nil
}

// CRSIsInitialized reports whether InitCRS has been called successfully.
func CRSIsInitialized() bool {
	return C.bb_crs_is_initialized() != 0
}
