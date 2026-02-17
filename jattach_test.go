//go:build (linux || darwin || windows) && (amd64 || arm64 || 386)
// +build linux darwin windows
// +build amd64 arm64 386

package jattach

import (
	"bufio"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"
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
	_, err := callJattach(0, []string{"properties"}, 1, 2)
	if err == nil {
		t.Error("Expected error for invalid PID, got nil")
	}

	_, err = callJattach(-1, []string{"properties"}, 1, 2)
	if err == nil {
		t.Error("Expected error for negative PID, got nil")
	}
}

func TestAttach_NoCommand(t *testing.T) {
	// Test with no command
	_, err := callJattach(1, []string{}, 1, 2)
	if err == nil {
		t.Error("Expected error for empty command, got nil")
	}
}

func TestIntegration_ThreadDump(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	javaPath, err := exec.LookPath("java")
	if err != nil {
		t.Skip("java not found in PATH")
	}

	javacPath, err := exec.LookPath("javac")
	if err != nil {
		t.Skip("javac not found in PATH")
	}

	// Compile the test fixture
	compileCmd := exec.Command(javacPath, "-d", "testdata", "testdata/SleepLoop.java")
	compileCmd.Dir = "."
	if out, err := compileCmd.CombinedOutput(); err != nil {
		t.Fatalf("Failed to compile SleepLoop.java: %v\n%s", err, out)
	}
	t.Cleanup(func() {
		os.Remove("testdata/SleepLoop.class")
	})

	// Launch the Java process
	cmd := exec.Command(javaPath, "-cp", "testdata", "SleepLoop")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatalf("Failed to create stdout pipe: %v", err)
	}
	if err := cmd.Start(); err != nil {
		t.Fatalf("Failed to start java process: %v", err)
	}
	t.Cleanup(func() {
		cmd.Process.Kill()
		cmd.Wait()
	})

	// Wait for "READY" signal with timeout
	ready := make(chan struct{})
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			if strings.TrimSpace(scanner.Text()) == "READY" {
				close(ready)
				return
			}
		}
	}()

	select {
	case <-ready:
	case <-time.After(30 * time.Second):
		t.Fatal("Timed out waiting for Java process to start")
	}

	// Take thread dump
	output, err := GetThreadDump(cmd.Process.Pid)
	if err != nil {
		t.Fatalf("GetThreadDump failed: %v", err)
	}

	lower := strings.ToLower(output)
	if !strings.Contains(lower, "sleeploop") && !strings.Contains(lower, "sleep") {
		t.Errorf("Thread dump does not mention SleepLoop or sleep.\nOutput:\n%s", output)
	}
}
