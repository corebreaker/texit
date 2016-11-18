package texit

import (
    "io/ioutil"
    "os"
    "os/exec"
    "strings"
    "syscall"
    "testing"
)

const TEXIT_ENVVAR = "_TEXIT_BE_CRASHER"
const TEXIT_ENVVAL = "1"

// Example
func MakeTestWithExit(func_to_test func()) (stdout, stderr string, status_code int, err error) {
    // Only run the failing part when a specific env variable is set
    if os.Getenv(ENVVAR) == TEXIT_ENVVAL {
        func_to_test()

        return
    }

    // Start the actual test in a different subprocess
    cmd := exec.Command(os.Args[0], "-test.run=TestCrashes")
    cmd.Env = append(os.Environ(), TEXIT_ENVVAR+"="+TEXIT_ENVVAL)

    stdout_pipe, err := cmd.StdoutPipe()
    if err != nil {
        return "", "", err
    }

    defer stdout_pipe.Close()

    stderr_pipe, _ := cmd.StderrPipe()
    if err != nil {
        return "", "", err
    }

    defer stderr_pipe.Close()

    if err := cmd.Start(); err != nil {
        return "", "", err
    }

    // Read the output stream
    stdout_bytes, err := ioutil.ReadAll(stdout_pipe)
    if err != nil {
        return "", "", err
    }

    // Read the error stream
    stderr_bytes, err := ioutil.ReadAll(stderr_pipe)
    if err != nil {
        return "", "", err
    }

    // Check that the program exited
    err := cmd.Wait()
    if e, ok := err.(*exec.ExitError); !ok || e.Success() {
        return "", "", err
    }

    status := cmd.ProcessState.Sys().(syscall.WaitStatus)

    return string(stdout_bytes), string(stderr_bytes), int(status.ExitCode), nil
}
