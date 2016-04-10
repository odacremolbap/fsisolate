package runtime

import "os"

// ChrootedProcess represents a process to be executed into a chroot sandbox
type ChrootedProcess struct {
	outStream *os.File
	errStream *os.File
	root      string
}

// NewChrootProcess returns a chroot process structure
func NewChrootProcess(root string) *ChrootedProcess {

	return &ChrootedProcess{
		outStream: os.Stdout,
		errStream: os.Stderr,
		root:      root,
	}

}

// SandboxExec executes command in chroot sandbox
func (p *ChrootedProcess) SandboxExec(command string, args ...string) {

}

// SendSignal sends TODO
func (p *ChrootedProcess) SendSignal(signal os.Signal) error {
	return nil
}

// GetStatus TODO
func (p *ChrootedProcess) GetStatus() (string, error) {
	return "", nil
}
