package texit

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func isVerbose() bool {
	for _, arg := range os.Args {
		if strings.HasPrefix(arg, "-test.v=") || (arg == "-test.v") {
			return true
		}
	}

	return false
}

// Test func call
func TestWithProgLine(t *testing.T) {
	progLine := func() string {
		_, line, _ := testFuncName()

		return fmt.Sprint(line + 5)
	}()

	os.Args = append(os.Args, _TEXIT_ARG+progLine)
	texitExit = func(_ int) {}

	_, _, _, _ = DoTestWithExit(func() {})
}

func TestWithExit1(t *testing.T) {
	stdout, stderr, statusCode, err := DoTestWithExit(func() {
		fmt.Println("Hello OUT!")
		_, _ = fmt.Fprintln(os.Stderr, "Hello ERR!")

		os.Exit(0)
	})

	if err != nil {
		t.Fatal(err, stderr, stdout)
	}

	if statusCode != 0 {
		t.Fatal("Bad status code: ", statusCode)
	}

	stdoutRef, stderrRef := "Hello OUT!", "Hello ERR!"
	if isVerbose() {
		funcName, _, err := testFuncName()
		if err != nil {
			t.Fatal(err)
		}

		stdoutRef = fmt.Sprintf("=== RUN   %s\n%s", funcName, stdoutRef)
	}

	if strings.TrimSpace(stdout) != stdoutRef {
		t.Fatal("Bad standard output: ", stdout, fmt.Sprintf("  `%s`", strings.TrimSpace(stdout)))
	}

	if strings.TrimSpace(stderr) != stderrRef {
		t.Fatal("Bad error output: ", stderr)
	}
}

func TestWithExit2(t *testing.T) {
	stdout, stderr, statusCode, err := DoTestWithExit(func() {
		fmt.Println("Hello OUT!")
		_, _ = fmt.Fprintln(os.Stderr, "Hello ERR!")

		os.Exit(123)
	})

	if err == nil {
		t.Fatal("Error expected", stderr, stdout)
	}

	if statusCode != 123 {
		t.Fatalf("Status code expected: %d != 1", statusCode)
	}

	stdoutRef, stderrRef := "Hello OUT!", "Hello ERR!"
	if isVerbose() {
		funcName, _, err := testFuncName()
		if err != nil {
			t.Fatal(err)
		}

		stdoutRef = fmt.Sprintf("=== RUN   %s\n%s", funcName, stdoutRef)
	}

	if strings.TrimSpace(stdout) != stdoutRef {
		t.Fatal("Bad standard output: ", stdout, fmt.Sprintf("  `%s`", strings.TrimSpace(stdout)))
	}

	if strings.TrimSpace(stderr) != stderrRef {
		t.Fatal("Bad error output: ", stderr)
	}
}

// Example and for code coverage
func TestWithExitDirectCall(t *testing.T) {
	signal := make(chan bool)
	defer close(signal)

	go func() {
		_, _, _, _ = DoTestWithExit(func() {
			_, _, _, _ = DoTestWithExit(func() {
			})
		})

		signal <- true
	}()

	<-signal
}

// Simulate call with arg
func TestCallToFront(t *testing.T) {
	last := len(os.Args)

	os.Args = append(os.Args, _TEXIT_ARG+"A")
	_, _, _, _ = DoTestWithExit(nil)
	os.Args[last] = _TEXIT_ARG + "1"
	_, _, _, _ = DoTestWithExit(nil)
}

type tReader struct {
	err error
}

func (rd *tReader) Read(p []byte) (int, error) { return 0, rd.err }
func (rd *tReader) Close() error               { return nil }

type tTestCmd struct {
	soReader   io.ReadCloser
	seReader   io.ReadCloser
	sePipeErr  error
	soPipeErr  error
	startErr   error
	waitErr    error
	exitStatus int
}

func (tc *tTestCmd) Start() error                         { return tc.startErr }
func (tc *tTestCmd) StdoutPipe() (io.ReadCloser, error)   { return tc.soReader, tc.soPipeErr }
func (tc *tTestCmd) StderrPipe() (io.ReadCloser, error)   { return tc.seReader, tc.sePipeErr }
func (tc *tTestCmd) GetExitStatus() int                   { return tc.exitStatus }
func (tc *tTestCmd) Wait() error                          { return tc.waitErr }
func (tc *tTestCmd) Exec(name string, arg ...string) iCmd { return tc }

//
func TestExecCommandWithError(t *testing.T) {
	err := errors.New("Error")

	soRd, seRd := &tReader{}, &tReader{}

	cmd := &tTestCmd{soPipeErr: err, soReader: soRd, seReader: seRd}
	cmdMaker = cmd

	_, _, _, _ = DoTestWithExit(nil)

	cmd.soPipeErr, cmd.sePipeErr = nil, err
	_, _, _, _ = DoTestWithExit(nil)

	cmd.sePipeErr, cmd.startErr = nil, err
	_, _, _, _ = DoTestWithExit(nil)

	cmd.startErr, soRd.err = nil, err
	_, _, _, _ = DoTestWithExit(nil)

	soRd.err, seRd.err = io.EOF, err
	_, _, _, _ = DoTestWithExit(nil)

	seRd.err, cmd.waitErr = io.EOF, err
	_, _, _, _ = DoTestWithExit(nil)

	cmd.waitErr = &exec.ExitError{ProcessState: &os.ProcessState{}}
	_, _, _, _ = DoTestWithExit(nil)
}
