package runtime

import "os"

// ChrootedProcess TODO
type ChrootedProcess struct {
}

// ExecNewRoot executes command in chroot sandbox
func ExecNewRoot(root string, command ...string) {

}

// SendSignal sends TODO
func (p *ChrootedProcess) SendSignal(signal os.Signal) error {
	return nil
}

// GetStatus TODO
func (p *ChrootedProcess) GetStatus() (string, error) {
	return "", nil
}
