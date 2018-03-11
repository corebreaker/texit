package texit

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
)

const (
	// Argument to notice the subprocess, try to have a pretty complicated argument to avoid collisions.
	_TEXIT_ARG = "TEXIT-X1RFWElUX0JFX0NSQVNIRVIK:"
)

// DoTestWithExit
func DoTestWithExit(func_to_test func()) (stdout, stderr string, status_code int, err error) {
	name, prog_line, err := test_func_name()
	if err != nil {
		return "", "", -1, err
	}

	line_arg := -1
	last_arg := os.Args[len(os.Args)-1]

	// If prefix matches the constant `_TEXIT_ARG`
	if strings.HasPrefix(last_arg, _TEXIT_ARG) {
		n, err := strconv.Atoi(last_arg[len(_TEXIT_ARG):])
		if err != nil {
			return "", "", -1, err
		}

		line_arg = n
	}

	// Only run the failing part when a specific command line argument is given
	// and line in source code matches the one that is passed as the command line argument suffix
	if line_arg == prog_line {
		func_to_test()

		os.Exit(0)
	}

	args := []string{"-test.run=" + name}
	for _, arg := range os.Args {
		if strings.HasPrefix(arg, "-test.coverprofile=") || strings.HasPrefix(arg, "-test.v=") || (arg == "-test.v") {
			args = append(args, arg)
		}
	}

	// Start the actual test in a different subprocess
	cmd := exec.Command(os.Args[0], append(args, fmt.Sprintf("%s%d", _TEXIT_ARG, prog_line))...)

	stdout_pipe, err := cmd.StdoutPipe()
	if err != nil {
		return "", "", -1, err
	}

	defer stdout_pipe.Close()

	stderr_pipe, _ := cmd.StderrPipe()
	if err != nil {
		return "", "", -1, err
	}

	defer stderr_pipe.Close()

	if err := cmd.Start(); err != nil {
		return "", "", -1, err
	}

	// Read the output stream
	stdout_bytes, err := ioutil.ReadAll(stdout_pipe)
	if err != nil {
		return "", "", -1, err
	}

	// Read the error stream
	stderr_bytes, err := ioutil.ReadAll(stderr_pipe)
	if err != nil {
		return "", "", -1, err
	}

	// Check that the program exited
	err = cmd.Wait()
	if e, ok := err.(*exec.ExitError); (err != nil) && (!(ok && e.Success())) {
		if ok && (!e.Success()) {
			status, ok := cmd.ProcessState.Sys().(syscall.WaitStatus)
			if ok {
				return string(stdout_bytes), string(stderr_bytes), status.ExitStatus(), err
			}
		}

		return string(stdout_bytes), string(stderr_bytes), -1, err
	}

	status, ok := cmd.ProcessState.Sys().(syscall.WaitStatus)
	if ok {
		status_code = status.ExitStatus()
	}

	return string(stdout_bytes), string(stderr_bytes), status_code, nil
}
