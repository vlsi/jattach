//go:build (linux || darwin) && (amd64 || arm64)
// +build linux darwin
// +build amd64 arm64

package jattach

import "github.com/vlsi/jattach/v2/src/posix"

// callJattach delegates to the POSIX-specific implementation
func callJattach(pid int, args []string, outFd int, errFd int) (int, error) {
	return posix.CallJattach(pid, args, outFd, errFd)
}
