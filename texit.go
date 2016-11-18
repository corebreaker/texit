package texit

import (
    "io/ioutil"
    "os"
    "os/exec"
    "strings"
    "testing"
)

// Example
func TestCrashes(t *testing.T) {
    // Only run the failing part when a specific env variable is set
    if os.Getenv("BE_CRASHER") == "1" {
        Crashes(42)
        return
    }

    // Start the actual test in a different subprocess
    cmd := exec.Command(os.Args[0], "-test.run=TestCrashes")
    cmd.Env = append(os.Environ(), "BE_CRASHER=1")
    stdout, _ := cmd.StderrPipe()
    if err := cmd.Start(); err != nil {
        t.Fatal(err)
    }

    // Check that the log fatal message is what we expected
    gotBytes, _ := ioutil.ReadAll(stdout)
    got := string(gotBytes)
    expected := "It crashes because you gave the answer"
    if !strings.HasSuffix(got[:len(got)-1], expected) {
        t.Fatalf("Unexpected log message. Got %s but should contain %s", got[:len(got)-1], expected)
    }

    // Check that the program exited
    err := cmd.Wait()
    if e, ok := err.(*exec.ExitError); !ok || e.Success() {
        t.Fatalf("Process ran with err %v, want exit status 1", err)
    }
}
