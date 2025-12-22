//go:build (linux || darwin) && (amd64 || arm64)
// +build linux darwin
// +build amd64 arm64

package posix

/*
#cgo CFLAGS: -I${SRCDIR} -O3
#cgo linux CFLAGS: -D_GNU_SOURCE

#include <stdlib.h>
#include "psutil.h"

// Forward declaration of the jattach function
extern int jattach(int pid, int argc, char** argv, int print_output);
*/
import "C"
import (
	"fmt"
	"unsafe"
)

// CallJattach is the low-level CGo wrapper for the jattach C function.
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
