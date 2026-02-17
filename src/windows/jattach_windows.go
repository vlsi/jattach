//go:build windows && (amd64 || arm64 || 386)
// +build windows
// +build amd64 arm64 386

package windows

/*
#cgo CFLAGS: -I${SRCDIR} -O2 -D_CRT_SECURE_NO_WARNINGS
#cgo LDFLAGS: -ladvapi32

#include <stdlib.h>
#include <io.h>
#include <fcntl.h>

// Forward declaration of the jattach function
extern int jattach(int pid, int argc, char** argv, int out_fd, int err_fd);

// Convert a Windows HANDLE (from Go's os.File.Fd()) to a C runtime file descriptor.
// Go's os.File.Fd() returns a Windows HANDLE, but jattach's C code uses _write()
// which requires a C runtime file descriptor.
static int handle_to_cfd(intptr_t handle) {
    if (handle < 0) return -1;
    return _open_osfhandle(handle, _O_WRONLY);
}
*/
import "C"
import (
	"fmt"
	"unsafe"
)

// CallJattach is the low-level CGo wrapper for the jattach C function (Windows implementation).
// It handles C string conversion and memory management.
// Returns the exit code from the jattach function.
func CallJattach(pid int, args []string, outFd int, errFd int) (int, error) {
	if pid <= 0 {
		return 1, fmt.Errorf("invalid PID: %d", pid)
	}

	if len(args) == 0 {
		return 1, fmt.Errorf("no command specified")
	}

	// Convert Go strings to C strings
	argc := C.int(len(args))
	argv := make([]*C.char, len(args))

	for i, arg := range args {
		argv[i] = C.CString(arg)
		defer C.free(unsafe.Pointer(argv[i]))
	}

	// Convert Windows HANDLEs to C runtime file descriptors.
	// Go's os.File.Fd() returns Windows HANDLEs, but the C code uses _write()
	// which requires C runtime file descriptors.
	cOutFd := C.handle_to_cfd(C.intptr_t(outFd))
	cErrFd := C.handle_to_cfd(C.intptr_t(errFd))

	// Call the C function
	ret := C.jattach(C.int(pid), argc, &argv[0], C.int(cOutFd), C.int(cErrFd))

	return int(ret), nil
}
