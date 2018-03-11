# texit
[![Coverage Status](https://coveralls.io/repos/github/corebreaker/texit/badge.svg?branch=master)](https://coveralls.io/github/corebreaker/texit?branch=master)
[![GoDoc](https://godoc.org/github.com/corebreaker/texit?status.svg)](https://godoc.org/github.com/corebreaker/texit)
![Version](https://img.shields.io/badge/version-1.0.0-green.svg)
[![release](https://img.shields.io/badge/release%20-v10.2-0077b3.svg?style=flat-square)](https://github.com/corebreaker/texit/releases)

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
