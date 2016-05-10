package fsisolate

import (
	"os"
	"runtime"
	"syscall"
	"testing"
	"time"
)

func TestExecute(t *testing.T) {

	var testData = []struct {
		root       string   // process root
		exec       string   // executable binary
		args       []string // arguments to the executable
		execOK     bool     // whether start should return OK or error
		waitOK     bool     // whether wait should return OK or error
		exitStatus int      // expected exit status
	}{
		{"testdata/simple", "/loop-" + runtime.GOOS, nil, true, true, 0},
		{"testdata/simple", "/loop-" + runtime.GOOS, []string{"-i=1"}, true, true, 0},
		{"testdata/simple", "/loop-" + runtime.GOOS, []string{"-e=1", "-i=2"}, true, false, 1},
	}

	for _, td := range testData {

		// get defaults
		p := NewChrootProcess(td.root)
		p.SetOutput(nil)

		err := p.Exec(td.exec, td.args...)
		if err != nil {
			if td.execOK {
				t.Errorf("Exectuion for [%s]%s with args %+v returned an error: %s", td.root, td.exec, td.args, err)
				continue
			}
			// error expected, continue with next test
			continue
		}

		err = p.Wait()
		if err != nil && td.waitOK {
			t.Errorf("Waiting for [%s]%s with args %+v returned an unexpected error: %s", td.root, td.exec, td.args, err)
			continue
		}

		st, err := p.GetExitStatus()
		if err != nil {
			t.Errorf("Getting exit status for [%s]%s with args %+v returned an unexpected error: %s", td.root, td.exec, td.args, err)
		}

		if st != td.exitStatus {
			t.Errorf("Exit status for [%s]%s with args %+v returned %d but expected %d", td.root, td.exec, td.args, st, td.exitStatus)
		}
	}
}

func TestNonExecutingPID(t *testing.T) {

	// get defaults
	p := NewChrootProcess("")

	pid, err := p.GetPID()
	if pid != 0 {
		t.Errorf("PID before execution should be 0, but returned %d", pid)
	}
	if err == nil {
		t.Errorf("PID before execution should return error, but returned nil")
	}
}

func TestSendSignal(t *testing.T) {

	var testData = []struct {
		signalPre  os.Signal // signal sent before executing loop
		signalExec os.Signal // signal sent during exection
	}{
		{syscall.SIGUSR1, syscall.SIGUSR1},
		{syscall.SIGUSR1, syscall.SIGKILL},
	}

	for _, td := range testData {

		// get defaults
		p := NewChrootProcess("testdata/simple")
		p.SetOutput(nil)

		err := p.SendSignal(td.signalPre)
		if err == nil {
			t.Errorf("Sending signal '%s' to non executing process should failed", td.signalPre)
		}

		err = p.Exec("/loop-"+runtime.GOOS, "-i=5")
		if err != nil {
			t.Errorf("Execution for senddignal test '%s' returned an error: %s", td.signalExec, err)
			continue
		}

		// wait and send signal
		time.Sleep(2 * time.Second)
		err = p.SendSignal(td.signalExec)
		if err != nil {
			t.Errorf("Sending signal '%s' to executing process failed: %s", td.signalExec, err)
		}

		// dismiss returned err. That is being tested somewhere else
		p.Wait()

	}
}
