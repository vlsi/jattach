//go:build (linux || darwin || windows) && (amd64 || arm64 || 386)
// +build linux darwin windows
// +build amd64 arm64 386

// Package jattach provides Go bindings for the jattach utility, which allows
// sending commands to running JVM processes via the Dynamic Attach mechanism.
//
// This package wraps the native C implementation of jattach using CGo, providing
// a type-safe, idiomatic Go API for interacting with Java Virtual Machines.
//
// Supported platforms: Linux, macOS, and Windows on amd64, arm64, and 386 architectures.
//
// Example usage:
//
//	// Get thread dump
//	output, err := jattach.ThreadDump(12345)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println(output)
//
//	// Load Java agent
//	err = jattach.LoadAgent(12345, "/path/to/agent.jar", "options", false)
//	if err != nil {
//		log.Fatal(err)
//	}
//
// For more information about jattach commands, see:
// https://github.com/jattach/jattach
package jattach

import (
	"bytes"
	"fmt"
	"io"
	"os"
)

// Command represents a JVM attach command type.
type Command string

// Supported JVM attach commands
const (
	// Load loads a native or Java agent into the target JVM
	Load Command = "load"

	// ThreadDump prints all stack traces (equivalent to jstack)
	ThreadDump Command = "threaddump"

	// DumpHeap dumps the heap (equivalent to jmap)
	DumpHeap Command = "dumpheap"

	// InspectHeap prints heap histogram (equivalent to jmap -histo)
	InspectHeap Command = "inspectheap"

	// DataDump shows heap and thread summary
	DataDump Command = "datadump"

	// Properties prints system properties
	Properties Command = "properties"

	// AgentProperties prints agent properties
	AgentProperties Command = "agentProperties"

	// SetFlag modifies a manageable VM flag
	SetFlag Command = "setflag"

	// PrintFlag prints a VM flag value
	PrintFlag Command = "printflag"

	// Jcmd executes an arbitrary jcmd command
	Jcmd Command = "jcmd"
)

// Attach sends a command to a JVM process.
// The command output is printed to stdout.
// Returns the exit code from the JVM command.
func Attach(pid int, cmd Command, args ...string) (int, error) {
	cmdArgs := append([]string{string(cmd)}, args...)
	return callJattach(pid, cmdArgs, true)
}

// AttachWithOutput sends a command to a JVM process and captures the output.
// Unlike Attach, this function captures stdout instead of printing it.
// Returns the captured output, exit code, and any error.
func AttachWithOutput(pid int, cmd Command, args ...string) (string, int, error) {
	// Create a pipe to capture stdout
	r, w, err := os.Pipe()
	if err != nil {
		return "", 1, fmt.Errorf("failed to create pipe: %w", err)
	}

	// Save original stdout and restore it when done
	oldStdout := os.Stdout
	defer func() {
		os.Stdout = oldStdout
	}()

	// Redirect stdout to our pipe
	os.Stdout = w

	// Capture output in a goroutine
	outputChan := make(chan string, 1)
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outputChan <- buf.String()
	}()

	// Execute the command
	cmdArgs := append([]string{string(cmd)}, args...)
	exitCode, err := callJattach(pid, cmdArgs, true)

	// Close the write end of the pipe
	w.Close()

	// Read the captured output
	output := <-outputChan
	r.Close()

	return output, exitCode, err
}

// GetThreadDump retrieves a thread dump from the target JVM.
// This is equivalent to running jstack.
func GetThreadDump(pid int) (string, error) {
	output, exitCode, err := AttachWithOutput(pid, ThreadDump)
	if err != nil {
		return "", err
	}
	if exitCode != 0 {
		return "", fmt.Errorf("threaddump command failed with exit code %d", exitCode)
	}
	return output, nil
}

// HeapDump creates a heap dump file at the specified path.
// The filepath should be accessible by the target JVM process.
func HeapDump(pid int, filepath string) error {
	_, exitCode, err := AttachWithOutput(pid, DumpHeap, filepath)
	if err != nil {
		return err
	}
	if exitCode != 0 {
		return fmt.Errorf("dumpheap command failed with exit code %d", exitCode)
	}
	return nil
}

// GetHeapHistogram retrieves a heap histogram from the target JVM.
// This is equivalent to running jmap -histo.
func GetHeapHistogram(pid int) (string, error) {
	output, exitCode, err := AttachWithOutput(pid, InspectHeap)
	if err != nil {
		return "", err
	}
	if exitCode != 0 {
		return "", fmt.Errorf("inspectheap command failed with exit code %d", exitCode)
	}
	return output, nil
}

// LoadAgent loads a native agent (.so) or Java agent (.jar) into the target JVM.
//
// Parameters:
//   - pid: Process ID of the target JVM
//   - agentPath: Path to the agent file (must be accessible by the target JVM)
//   - options: Agent options string (can be empty)
//   - absolute: If true, agentPath is treated as absolute; if false, relative to JVM working directory
func LoadAgent(pid int, agentPath string, options string, absolute bool) error {
	absoluteStr := "false"
	if absolute {
		absoluteStr = "true"
	}

	args := []string{agentPath, absoluteStr}
	if options != "" {
		args = append(args, options)
	}

	_, exitCode, err := AttachWithOutput(pid, Load, args...)
	if err != nil {
		return err
	}
	if exitCode != 0 {
		return fmt.Errorf("load command failed with exit code %d", exitCode)
	}
	return nil
}

// GetSystemProperties retrieves all system properties from the target JVM.
func GetSystemProperties(pid int) (string, error) {
	output, exitCode, err := AttachWithOutput(pid, Properties)
	if err != nil {
		return "", err
	}
	if exitCode != 0 {
		return "", fmt.Errorf("properties command failed with exit code %d", exitCode)
	}
	return output, nil
}

// GetAgentProperties retrieves agent properties from the target JVM.
func GetAgentProperties(pid int) (string, error) {
	output, exitCode, err := AttachWithOutput(pid, AgentProperties)
	if err != nil {
		return "", err
	}
	if exitCode != 0 {
		return "", fmt.Errorf("agentProperties command failed with exit code %d", exitCode)
	}
	return output, nil
}

// ExecuteJcmd executes an arbitrary jcmd command on the target JVM.
// The jcmdArgs are passed directly to jcmd.
//
// Example:
//
//	output, err := jattach.ExecuteJcmd(12345, "VM.flags", "-all")
func ExecuteJcmd(pid int, jcmdArgs ...string) (string, error) {
	output, exitCode, err := AttachWithOutput(pid, Jcmd, jcmdArgs...)
	if err != nil {
		return "", err
	}
	if exitCode != 0 {
		return "", fmt.Errorf("jcmd command failed with exit code %d", exitCode)
	}
	return output, nil
}

// SetVMFlag sets a manageable VM flag to a new value.
func SetVMFlag(pid int, flagName string, value string) error {
	_, exitCode, err := AttachWithOutput(pid, SetFlag, flagName, value)
	if err != nil {
		return err
	}
	if exitCode != 0 {
		return fmt.Errorf("setflag command failed with exit code %d", exitCode)
	}
	return nil
}

// PrintVMFlag prints the value of a VM flag.
func PrintVMFlag(pid int, flagName string) (string, error) {
	output, exitCode, err := AttachWithOutput(pid, PrintFlag, flagName)
	if err != nil {
		return "", err
	}
	if exitCode != 0 {
		return "", fmt.Errorf("printflag command failed with exit code %d", exitCode)
	}
	return output, nil
}
