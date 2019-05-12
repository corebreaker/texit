package texit

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const (
	// Argument to notice the subprocess, try to have a pretty complicated argument to avoid collisions.
	_TEXIT_ARG = "TEXIT-X1RFWElUX0JFX0NSQVNIRVIK:"
)

var texitExit = os.Exit
var cmdMaker iExec = tStdExec{}

func runTest(funcToTest func(), progLine int) (done bool, err error) {
	lineArg := -1
	lastArg := os.Args[len(os.Args)-1]

	// If prefix matches the constant `_TEXIT_ARG`
	if strings.HasPrefix(lastArg, _TEXIT_ARG) {
		n, err := strconv.Atoi(lastArg[len(_TEXIT_ARG):])
		if err != nil {
			return true, err
		}

		lineArg = n
	}

	// Only run the failing part when a specific command line argument is given
	// and line in source code matches the one that is passed as the command line argument suffix
	if lineArg == progLine {
		funcToTest()

		texitExit(0)
		done = true
	}

	return
}

func makeCommand(name string, progLine int) iCmd {
	args := []string{"-test.run=" + name}
	for _, arg := range os.Args {
		const coverprofileArg = "-test.coverprofile="
		const vArg = "-test.v"

		if strings.HasPrefix(arg, coverprofileArg) || strings.HasPrefix(arg, vArg+"=") || (arg == vArg) {
			args = append(args, arg)
		}
	}

	// Start the actual test in a different subprocess
	return cmdMaker.Exec(os.Args[0], append(args, fmt.Sprintf("%s%d", _TEXIT_ARG, progLine))...)
}

func readStreams(stdout, stderr io.ReadCloser) (outBytes, errBytes []byte, err error) {
	// Read the output stream
	outBytes, err = ioutil.ReadAll(stdout)
	if err != nil {
		return
	}

	// Read the error stream
	errBytes, err = ioutil.ReadAll(stderr)
	if err != nil {
		return
	}

	return
}

func execCmd(name string, progLine int) (cmd iCmd, outBytes, errBytes []byte, err error) {
	cmd = makeCommand(name, progLine)

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return nil, nil, nil, err
	}

	defer stdoutPipe.Close()

	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return nil, nil, nil, err
	}

	defer stderrPipe.Close()

	if err := cmd.Start(); err != nil {
		return nil, nil, nil, err
	}

	stdoutBytes, stderrBytes, err := readStreams(stdoutPipe, stderrPipe)
	if err != nil {
		return nil, nil, nil, err
	}

	// Check that the program exited
	return cmd, stdoutBytes, stderrBytes, cmd.Wait()
}

// DoTestWithExit the function for doing test on function with os.Exit
func DoTestWithExit(funcToTest func()) (stdout, stderr string, statusCode int, err error) {
	name, progLine, err := testFuncName()
	if err != nil {
		return "", "", -1, err
	}

	{
		done, err := runTest(funcToTest, progLine)
		if done {
			return "", "", -1, err
		}
	}

	cmd, outBytes, errBytes, err := execCmd(name, progLine)
	exitStatus := -1
	if cmd != nil {
		exitStatus = cmd.GetExitStatus()
	}

	if e, ok := err.(*exec.ExitError); (err != nil) && ok && e.Success() {
		e.Success()
		err = nil
	}

	return string(outBytes), string(errBytes), exitStatus, err
}
