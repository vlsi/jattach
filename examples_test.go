//go:build (linux || darwin || windows) && (amd64 || arm64 || 386)
// +build linux darwin windows
// +build amd64 arm64 386

package jattach_test

import (
	"fmt"
	"log"

	jattach "github.com/vlsi/jattach/v2"
)

// ExampleGetThreadDump demonstrates how to get a thread dump from a JVM process.
func ExampleGetThreadDump() {
	pid := 12345 // Replace with actual JVM PID

	output, err := jattach.GetThreadDump(pid)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(output)
}

// ExampleGetSystemProperties demonstrates how to retrieve system properties.
func ExampleGetSystemProperties() {
	pid := 12345 // Replace with actual JVM PID

	props, err := jattach.GetSystemProperties(pid)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(props)
}

// ExampleLoadAgent demonstrates how to load a Java agent.
func ExampleLoadAgent() {
	pid := 12345 // Replace with actual JVM PID

	// Load agent with absolute path
	err := jattach.LoadAgent(pid, "/path/to/agent.jar", "option1,option2", true)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Agent loaded successfully")
}

// ExampleExecuteJcmd demonstrates how to execute a jcmd command.
func ExampleExecuteJcmd() {
	pid := 12345 // Replace with actual JVM PID

	// Get VM version
	output, err := jattach.ExecuteJcmd(pid, "VM.version")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(output)
}

// ExampleHeapDump demonstrates how to create a heap dump.
func ExampleHeapDump() {
	pid := 12345 // Replace with actual JVM PID

	err := jattach.HeapDump(pid, "/tmp/heapdump.hprof")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Heap dump created at /tmp/heapdump.hprof")
}

// ExampleGetHeapHistogram demonstrates how to get a heap histogram.
func ExampleGetHeapHistogram() {
	pid := 12345 // Replace with actual JVM PID

	histogram, err := jattach.GetHeapHistogram(pid)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(histogram)
}

// ExampleAttach demonstrates low-level attach with custom command.
func ExampleAttach() {
	pid := 12345 // Replace with actual JVM PID

	// Execute a command and print output to stdout
	exitCode, err := jattach.Attach(pid, jattach.Properties)
	if err != nil {
		log.Fatal(err)
	}

	if exitCode != 0 {
		fmt.Printf("Command failed with exit code %d\n", exitCode)
	}
}

// ExampleAttachWithOutput demonstrates capturing command output.
func ExampleAttachWithOutput() {
	pid := 12345 // Replace with actual JVM PID

	// Execute command and capture output
	output, exitCode, err := jattach.AttachWithOutput(pid, jattach.Jcmd, "VM.flags")
	if err != nil {
		log.Fatal(err)
	}

	if exitCode != 0 {
		fmt.Printf("Command failed with exit code %d\n", exitCode)
		return
	}

	fmt.Println(output)
}
