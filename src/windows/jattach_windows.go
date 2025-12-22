//go:build windows && (amd64 || arm64 || 386)
// +build windows
// +build amd64 arm64 386

package windows

/*
#cgo CFLAGS: -I${SRCDIR} -O2 -D_CRT_SECURE_NO_WARNINGS
#cgo LDFLAGS: -ladvapi32

#include <stdlib.h>

// Forward declaration of the jattach function
extern int jattach(int pid, int argc, char** argv, int print_output);
*/
import "C"
import (
	"fmt"
	"unsafe"
)

// CallJattach is the low-level CGo wrapper for the jattach C function (Windows implementation).
// It handles C string conversion and memory management.
// Returns the exit code from the jattach function.
func CallJattach(pid int, args []string, printOutput bool) (int, error) {
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

	// Determine print_output flag
	printOutputInt := C.int(0)
	if printOutput {
		printOutputInt = C.int(1)
	}

	// Call the C function
	ret := C.jattach(C.int(pid), argc, &argv[0], printOutputInt)

	return int(ret), nil
}
