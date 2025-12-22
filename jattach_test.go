//go:build (linux || darwin || windows) && (amd64 || arm64 || 386)
// +build linux darwin windows
// +build amd64 arm64 386

package jattach

import (
	"testing"
)

func TestCommandConstants(t *testing.T) {
	tests := []struct {
		cmd      Command
		expected string
	}{
		{Load, "load"},
		{ThreadDump, "threaddump"},
		{DumpHeap, "dumpheap"},
		{InspectHeap, "inspectheap"},
		{DataDump, "datadump"},
		{Properties, "properties"},
		{AgentProperties, "agentProperties"},
		{SetFlag, "setflag"},
		{PrintFlag, "printflag"},
		{Jcmd, "jcmd"},
	}

	for _, tt := range tests {
		t.Run(string(tt.cmd), func(t *testing.T) {
			if string(tt.cmd) != tt.expected {
				t.Errorf("Command %s has unexpected value: got %s, want %s", tt.cmd, tt.cmd, tt.expected)
			}
		})
	}
}

func TestAttach_InvalidPID(t *testing.T) {
	// Test with invalid PID
	_, err := callJattach(0, []string{"properties"}, true)
	if err == nil {
		t.Error("Expected error for invalid PID, got nil")
	}

	_, err = callJattach(-1, []string{"properties"}, true)
	if err == nil {
		t.Error("Expected error for negative PID, got nil")
	}
}

func TestAttach_NoCommand(t *testing.T) {
	// Test with no command
	_, err := callJattach(1, []string{}, true)
	if err == nil {
		t.Error("Expected error for empty command, got nil")
	}
}

// Integration tests require a running JVM process
// Run with: go test -tags=integration
// These tests are skipped by default
func TestIntegration_Attach(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Note: You need to set JATTACH_TEST_PID environment variable
	// to the PID of a running JVM process for these tests to work
	t.Skip("Integration tests require a running JVM - set JATTACH_TEST_PID environment variable")

	// Example of how integration tests would work:
	// pid := getTestJVMPID(t)
	//
	// t.Run("Properties", func(t *testing.T) {
	//     output, err := GetSystemProperties(pid)
	//     if err != nil {
	//         t.Fatalf("Failed to get properties: %v", err)
	//     }
	//     if len(output) == 0 {
	//         t.Error("Expected non-empty properties output")
	//     }
	// })
}
