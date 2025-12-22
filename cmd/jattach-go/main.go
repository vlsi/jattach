package main

import (
	"fmt"
	"os"
	"strconv"

	jattach "github.com/vlsi/jattach/v2"
)

const version = "2.2-go"

func printUsage() {
	fmt.Printf("jattach-go %s\n", version)
	fmt.Println()
	fmt.Println("Usage: jattach-go <pid> <cmd> [args ...]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("    load  threaddump   dumpheap  setflag    properties")
	fmt.Println("    jcmd  inspectheap  datadump  printflag  agentProperties")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("    jattach-go 12345 threaddump")
	fmt.Println("    jattach-go 12345 load /path/to/agent.jar true")
	fmt.Println("    jattach-go 12345 jcmd VM.version")
}

func main() {
	if len(os.Args) < 3 {
		printUsage()
		os.Exit(1)
	}

	// Parse PID
	pidStr := os.Args[1]
	pid, err := strconv.Atoi(pidStr)
	if err != nil || pid <= 0 {
		fmt.Fprintf(os.Stderr, "%s is not a valid process ID\n", pidStr)
		os.Exit(1)
	}

	// Parse command and arguments
	cmd := jattach.Command(os.Args[2])
	args := []string{}
	if len(os.Args) > 3 {
		args = os.Args[3:]
	}

	// Execute the command
	exitCode, err := jattach.Attach(pid, cmd, args...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	os.Exit(exitCode)
}
