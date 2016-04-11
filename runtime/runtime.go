package runtime

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"sync"
	"syscall"
)

// ProcessState is the state in which a process can be
type ProcessState string

// Possible process states
const (
	NotStarted ProcessState = "not started"
	Running    ProcessState = "running"
	Finished   ProcessState = "finished"
)

// ChrootedProcess represents a process to be executed into a chroot sandbox
// outputStream is accessed using a method
// root shouldn't change, it can only be set on creation
// cmd is set when the process is started.
type ChrootedProcess struct {
	sync.Mutex
	outStream *os.File
	root      string
	cmd       *exec.Cmd
}

// NewChrootProcess returns a chroot process structure
func NewChrootProcess(root string) *ChrootedProcess {
	return &ChrootedProcess{
		outStream: os.Stdout,
		root:      root,
	}
}

// SetOutput sets output stream for sandboxed process
func (p *ChrootedProcess) SetOutput(out *os.File) {
	p.outStream = out
}

// Exec executes command in chroot sandbox
func (p *ChrootedProcess) Exec(command string, args ...string) error {
	p.Lock()
	defer p.Unlock()

	if p.getState() == Running {
		return fmt.Errorf("Error starting process: there is another process executing in this chroot")
	}

	// Execution command is chroot, so we need to use command as first argument
	// and actual command arguments behind
	chargs := []string{p.root, command}
	for _, arg := range args {
		chargs = append(chargs, arg)
	}

	p.cmd = exec.Command("chroot", chargs...)

	// get stdout from chrooted process
	reader, err := p.cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("Error getting chrooted process stdout: %s", err.Error())
	}

	// send chrooted process output to our output stream
	procscan := bufio.NewScanner(reader)
	go func() {
		for procscan.Scan() {
			// TODO parametrize this prefix
			fmt.Fprintf(p.outStream, "[CHROOT]%s\n", procscan.Text())
		}
	}()

	// start process
	err = p.cmd.Start()
	if err != nil {
		return fmt.Errorf("Error starting process: %s", err.Error())
	}
	return nil
}

// Wait waits for the execution to end
func (p *ChrootedProcess) Wait() error {
	if p.getState() != Running {
		return fmt.Errorf("Error waiting process: process has not started")
	}

	err := p.cmd.Wait()
	if err != nil {
		return fmt.Errorf("Error waiting process: %s", err.Error())
	}
	return nil
}

// SendSignal sends signal to the chrooted process
func (p *ChrootedProcess) SendSignal(signal os.Signal) error {

	// the important matter is that p.cmd.Process is not nil to avoid panic
	if p.getState() != Running {
		return fmt.Errorf("Error sending signal process: process has not started")
	}

	err := p.cmd.Process.Signal(signal)
	if err != nil {
		return err
	}

	return nil
}

// GetPID returns PID from the process (started or finished)
func (p *ChrootedProcess) GetPID() (int, error) {

	if p.getState() == NotStarted {
		return 0, fmt.Errorf("Error getting PID: process has not started")
	}

	return p.cmd.Process.Pid, nil
}

// GetExitStatus returns exit status once the process has finished
func (p *ChrootedProcess) GetExitStatus() (int, error) {

	if p.getState() != Finished {
		return 0, fmt.Errorf("Error getting exit status: process is not finished")
	}

	return p.cmd.ProcessState.Sys().(syscall.WaitStatus).ExitStatus(), nil
}

// GetState returns the process state
func (p *ChrootedProcess) GetState() ProcessState {
	return p.getState()
}

// getState internal method that doesn't lock
func (p *ChrootedProcess) getState() ProcessState {
	// process not started
	if p.cmd == nil || p.cmd.Process == nil {
		return NotStarted
	}
	// process running
	if p.cmd.ProcessState == nil || p.cmd.ProcessState.Exited() == false {
		return Running
	}
	// process finished
	return Finished
}
