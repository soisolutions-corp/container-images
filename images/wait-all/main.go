package main

import (
	"golang.org/x/sys/unix"
	"log"
	"os"
)

func main() {
	// Set up this process to adopt abandoned children
	must(unix.Prctl(unix.PR_SET_CHILD_SUBREAPER, uintptr(1), 0, 0, 0))

	pid := os.Getpid()
	ppid := os.Getppid()
	log.Printf("pid: %d, ppid: %d, args: %s", pid, ppid, os.Args)

	// Run the process
	proc := RunCommand(os.Args[1:])

	log.Printf("pid: %d, forked process", proc.cmd.Pid)

	exitCode := proc.Wait()

	log.Printf("exiting with code %d", exitCode.ExitCode())
	os.Exit(exitCode.ExitCode())
}
