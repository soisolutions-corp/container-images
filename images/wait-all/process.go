package main

import (
	"errors"
	"github.com/shirou/gopsutil/v4/process"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// Command holds the processes that we run.
type Command struct {
	sync.RWMutex
	wg       sync.WaitGroup
	args     []string
	cmd      *os.Process
	ppid     int
	pids     []int
	exitcode int
}

func RunCommand(args []string) *Command {
	log.Printf("running command: %s", args)

	cmd := &Command{
		args: args,
		ppid: os.Getppid(),
	}

	cmd.Start()

	return cmd
}

func (c *Command) Pid() int32 {
	return int32(c.cmd.Pid)
}

func (c *Command) IsRunning() bool {
	running, err := process.PidExists(c.Pid())
	mustNot(err)

	return running
}

func (c *Command) Start() {
	go c.HandleSignals()
	go c.Supervise()

	// Give go routine time to start
	time.Sleep(1 * time.Millisecond)

	var err error
	c.cmd, err = os.StartProcess(os.Args[1], os.Args[1:], &os.ProcAttr{
		Env: os.Environ(),
		Files: []*os.File{
			os.Stdin,
			os.Stdout,
			os.Stderr,
		},
	})
	mustNot(err)

	c.wg.Add(1)

	return
}

// Wait waits for the supervised command to exit
func (c *Command) Wait() *os.ProcessState {
	if c.cmd == nil {
		panic("Command not started")
	}

	status, err := c.cmd.Wait()
	mustNot(err)

	log.Printf("process exited with status: %v", status)

	// WaitGroup lock must be released before we can exit
	// This allows for all supervised processes to cleanly exit
	c.wg.Wait()

	log.Printf("wait group is released")

	return status
}

// Children returns the process ids of all child processes
func (c *Command) Children() []int {
	// Get our process info
	proc, err := process.NewProcess(int32(os.Getpid()))
	mustNot(err)

	children, err := proc.Children()
	if err != nil {
		return []int{}
	}

	// output the pids from the children array
	var pids []int
	for _, child := range children {
		pids = append(pids, int(child.Pid))
	}

	return pids
}

// Signal forwards the signal to the requested pid, or the supervised pid if pid is nil
func (c *Command) Signal(sig os.Signal, pid *int) {
	if pid == nil {
		pid = &c.cmd.Pid
	}

	log.Printf("sending signal %s to pid %d", sig.String(), *pid)

	must(syscall.Kill(*pid, sig.(syscall.Signal)))
}

// SignalAll sends a signal to the entire process group
func (c *Command) SignalAll(sig os.Signal) {
	// Signal the process group
	pid := 0
	c.Signal(sig, &pid)
}

// SignalAllChildren sends a signal to all child processes of the supervisor
func (c *Command) SignalAllChildren(sig os.Signal) {
	for _, pid := range c.Children() {
		c.Signal(sig, &pid)
	}
}

// Cleanup sends signals to all known child processes of the supervisor
func (c *Command) Cleanup(sig os.Signal) {
	log.Printf("received %s signal, sending to children", sig.String())
	c.SignalAllChildren(sig)

	tick := time.NewTicker(100 * time.Millisecond)
	defer tick.Stop()

	timer := time.NewTimer(15 * time.Second)
	defer timer.Stop()

	for {
		select {
		case <-tick.C:
			if len(c.Children()) == 0 {
				log.Printf("all child processes have exited")
				c.Kill()
				return
			}
		case <-timer.C:
			log.Printf("timed out waiting for child processes to exit")
			c.Kill()
			return
		}
	}
}

func (c *Command) Kill() {
	log.Printf("sending SIGKILL to process group")
	c.SignalAll(syscall.SIGKILL)
	c.wg.Done()
}

func (c *Command) Supervise() {
	tick := time.NewTicker(1 * time.Second)
	defer tick.Stop()

	for {
		<-tick.C

		if c.cmd == nil {
			log.Printf("supervised process has not started yet")
			continue
		}

		if !c.IsRunning() {
			log.Printf("supervised process is not alive, calling for cleanup")
			c.Cleanup(syscall.SIGTERM)
			return
		}
	}
}

func (c *Command) HandleSignals() {
	ints := make(chan os.Signal, 3)
	signal.Notify(ints, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	other := make(chan os.Signal, 3)
	signal.Notify(other, syscall.SIGHUP)

	child := make(chan os.Signal, 3)
	signal.Notify(child, syscall.SIGCHLD)

	for {
		select {
		case sig := <-other:
			c.Signal(sig, nil)
		case sig := <-ints:
			c.Cleanup(sig)
		case sig := <-child:
			go c.ReapChild(sig)
		default:
			// signal receiver wasn't ready, go around for another attempt
		}
	}
}

func (c *Command) ReapChild(sig os.Signal) {
	var status syscall.WaitStatus

	log.Printf("reaping child process")

	// Gets the pid of the abandoned child process
	pid, err := syscall.Wait4(-1, &status, syscall.WNOHANG, nil)

	for errors.Is(err, syscall.EINTR) {
		// Waits for the specific child process to exit
		pid, err = syscall.Wait4(pid, &status, syscall.WNOHANG, nil)
		log.Printf("pid: %d, exit status: %d", pid, status.ExitStatus())
	}

	if errors.Is(err, syscall.ECHILD) {
		// no children are left to clean up, supervisor can exit
		log.Printf("no children left to clean up")
		c.wg.Done()
		return
	}

	if err != nil {
		log.Printf("error reaping child process: %v", err)
	}

	log.Printf("done reaping child process")
	return
}
