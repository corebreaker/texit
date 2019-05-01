package texit

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
)

func is_verbose() bool {
	for _, arg := range os.Args {
		if strings.HasPrefix(arg, "-test.v=") || (arg == "-test.v") {
			return true
		}
	}

	return false
}

// Test func call
func TestWithProgLine(t *testing.T) {
	prog_line := func() string {
		_, line, _ := test_func_name()

		return fmt.Sprint(line + 5)
	}()

	os.Args = append(os.Args, _TEXIT_ARG+prog_line)
	texitExit = func(_ int) {}

	DoTestWithExit(func() {})
}

func TestWithExit1(t *testing.T) {
	stdout, stderr, status_code, err := DoTestWithExit(func() {
		fmt.Println("Hello OUT!")
		fmt.Fprintln(os.Stderr, "Hello ERR!")

		os.Exit(0)
	})

	if err != nil {
		t.Fatal(err, stderr, stdout)
	}

	if status_code != 0 {
		t.Fatal("Bad status code: ", status_code)
	}

	stdout_ref, stderr_ref := "Hello OUT!", "Hello ERR!"
	if is_verbose() {
		func_name, _, err := test_func_name()
		if err != nil {
			t.Fatal(err)
		}

		stdout_ref = fmt.Sprintf("=== RUN   %s\n%s", func_name, stdout_ref)
	}

	if strings.TrimSpace(stdout) != stdout_ref {
		t.Fatal("Bad standard output: ", stdout, fmt.Sprintf("  `%s`", strings.TrimSpace(stdout)))
	}

	if strings.TrimSpace(stderr) != stderr_ref {
		t.Fatal("Bad error output: ", stderr)
	}
}

func TestWithExit2(t *testing.T) {
	stdout, stderr, status_code, err := DoTestWithExit(func() {
		fmt.Println("Hello OUT!")
		fmt.Fprintln(os.Stderr, "Hello ERR!")

		os.Exit(123)
	})

	if err == nil {
		t.Fatal("Error expected", stderr, stdout)
	}

	if status_code != 123 {
		t.Fatalf("Status code expected: %d != 1", status_code)
	}

	stdout_ref, stderr_ref := "Hello OUT!", "Hello ERR!"
	if is_verbose() {
		func_name, _, err := test_func_name()
		if err != nil {
			t.Fatal(err)
		}

		stdout_ref = fmt.Sprintf("=== RUN   %s\n%s", func_name, stdout_ref)
	}

	if strings.TrimSpace(stdout) != stdout_ref {
		t.Fatal("Bad standard output: ", stdout, fmt.Sprintf("  `%s`", strings.TrimSpace(stdout)))
	}

	if strings.TrimSpace(stderr) != stderr_ref {
		t.Fatal("Bad error output: ", stderr)
	}
}

// Example and for code coverage
func TestWithExitDirectCall(t *testing.T) {
	signal := make(chan bool)
	defer close(signal)

	go func() {
		DoTestWithExit(func() {
			DoTestWithExit(func() {
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
	DoTestWithExit(nil)
	os.Args[last] = _TEXIT_ARG + "1"
	DoTestWithExit(nil)
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

	f := func() {}

	DoTestWithExit(f)

	cmd.soPipeErr, cmd.sePipeErr = nil, err
	DoTestWithExit(f)

	cmd.sePipeErr, cmd.startErr = nil, err
	DoTestWithExit(f)

	cmd.startErr, soRd.err = nil, err
	DoTestWithExit(f)

	soRd.err, seRd.err = io.EOF, err
	DoTestWithExit(f)

	seRd.err, cmd.waitErr = io.EOF, err
	DoTestWithExit(f)
}
