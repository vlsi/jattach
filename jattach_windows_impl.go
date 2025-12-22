//go:build windows && (amd64 || arm64 || 386)
// +build windows
// +build amd64 arm64 386

package jattach

import "github.com/vlsi/jattach/v2/src/windows"

// callJattach delegates to the Windows-specific implementation
func callJattach(pid int, args []string, printOutput bool) (int, error) {
	return windows.CallJattach(pid, args, printOutput)
}
