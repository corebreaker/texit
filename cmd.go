package texit

import (
	"io"
	"os/exec"
)

type iCmd interface {
	Start() error
	StdoutPipe() (io.ReadCloser, error)
	StderrPipe() (io.ReadCloser, error)
	Wait() error
	GetExitStatus() int
}

type iExec interface {
	Exec(name string, arg ...string) iCmd
}

type tStdCmd struct {
	cmd *exec.Cmd
}

func (sc *tStdCmd) Start() error {
	return sc.cmd.Start()
}

func (sc *tStdCmd) StdoutPipe() (io.ReadCloser, error) {
	return sc.cmd.StdoutPipe()
}

func (sc *tStdCmd) StderrPipe() (io.ReadCloser, error) {
	return sc.cmd.StderrPipe()
}

func (sc *tStdCmd) Wait() error {
	return sc.cmd.Wait()
}

func (sc *tStdCmd) GetExitStatus() int {
	return sc.cmd.ProcessState.ExitCode()
}

type tStdExec struct{}

func (tStdExec) Exec(name string, arg ...string) iCmd {
	return &tStdCmd{exec.Command(name, arg...)}
}
