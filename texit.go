package texit

import (
    "io/ioutil"
    "os"
    "os/exec"
    "strings"
    "syscall"
)

const (
    TEXIT_ENVVAR = "_TEXIT_BE_CRASHER"
    TEXIT_ENVVAL = "1"
)

// Example
func MakeTestWithExit(func_to_test func()) (stdout, stderr string, status_code int, err error) {
    // Only run the failing part when a specific env variable is set
    if os.Getenv(TEXIT_ENVVAR) == TEXIT_ENVVAL {
        func_to_test()

        os.Exit(0)
    }

    name, err := func_name()
    if err != nil {
        return "", "", -1, err
    }

    args := []string{"-test.run=" + name}
    for _, arg := range os.Args {
        if strings.HasPrefix(arg, "-test.coverprofile=") {
            args = append(args, arg)

            break
        }
    }

    // Start the actual test in a different subprocess
    cmd := exec.Command(os.Args[0], args...)
    cmd.Env = append(os.Environ(), TEXIT_ENVVAR+"="+TEXIT_ENVVAL)

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
    if e, ok := err.(*exec.ExitError); (err != nil) && (!ok || e.Success()) {
        return "", "", -1, err
    }

    status, ok := cmd.ProcessState.Sys().(syscall.WaitStatus)
    if !ok {
        status_code = status.ExitStatus()
    }

    return string(stdout_bytes), string(stderr_bytes), status_code, nil
}
