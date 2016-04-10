package runtime

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"syscall"

	log "github.com/sirupsen/logrus"
)

// ChrootedProcess represents a process to be executed into a chroot sandbox
type ChrootedProcess struct {
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

// SandboxExec executes command in chroot sandbox
func (p *ChrootedProcess) SandboxExec(command string, args ...string) error {
	log.Debugf("Executing '%s %v' in sandbox %s", command, args, p.root)

	// TODO validation
	// make sure there is no other cmd being executed

	// add root dir as first arg to chroot, command as second
	chargs := []string{p.root, command}
	for _, arg := range args {
		chargs = append(chargs, arg)
	}

	p.cmd = exec.Command("chroot", chargs...)

	// get stdout from chrooted proc
	reader, err := p.cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("Error connecting process pipe: %s", err.Error())
	}

	// scan stream and send to output
	procscan := bufio.NewScanner(reader)
	go func() {
		for procscan.Scan() {
			fmt.Fprintf(p.outStream, "%s\n", procscan.Text())
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

	// wait for the process to exit
	err := p.cmd.Wait()
	if err != nil {
		return fmt.Errorf("Error waiting process: %s", err.Error())
	}
	return nil
}

// SendSignal sends TODO
func (p *ChrootedProcess) SendSignal(signal os.Signal) error {
	if p.cmd == nil || p.cmd.Process == nil {
		return fmt.Errorf("Process does not exists")
	}

	err := p.cmd.Process.Signal(signal)
	if err != nil {
		return err
	}

	return nil
}

// GetPID returns PID from the process (even if it's not running)
func (p *ChrootedProcess) GetPID() (int, error) {
	if p.cmd == nil || p.cmd.Process == nil {
		return 0, fmt.Errorf("Process does not exists")
	}

	return p.cmd.Process.Pid, nil
}

// GetExited returns boolean indicating if process has finished
func (p *ChrootedProcess) GetExited() (bool, error) {
	if _, err := p.GetPID(); err != nil {
		return true, err
	}

	if p.cmd.ProcessState != nil && p.cmd.ProcessState.Exited() {
		return true, nil
	}
	return false, nil
}

// GetExitStatus returns process exit status
func (p *ChrootedProcess) GetExitStatus() (int, error) {
	if p.cmd == nil || p.cmd.ProcessState == nil || !p.cmd.ProcessState.Exited() {
		return 0, fmt.Errorf("Process does not exists, %t", p.cmd.ProcessState.Exited())
	}

	return p.cmd.ProcessState.Sys().(syscall.WaitStatus).ExitStatus(), nil
}

// SetOutput sets output stream for sandboxed process
func (p *ChrootedProcess) SetOutput(out *os.File) {
	p.outStream = out
}
