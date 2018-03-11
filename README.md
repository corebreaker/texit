# texit
Yes, you can use os.Exit() in Go tests.

Inspired by https://talks.golang.org/2014/testing.slide#23

## Example in a testing function

```golang
import (
    "os"

    "github.com/corebreaker/texit"
)

func TestWithExitDirectCall(t *testing.T) {
	stdout, stderr, status, err := DoTestWithExit(func() {
        // Something to do

        os.Exit(0)
	})

    // …
}
```
