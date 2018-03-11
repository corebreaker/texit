package texit

import (
	"fmt"
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
	DoTestWithExit(func() {
		DoTestWithExit(func() {
		})
	})
}
