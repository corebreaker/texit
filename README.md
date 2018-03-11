# texit
[![Build Status](https://travis-ci.org/corebreaker/texit.svg?branch=master)](https://travis-ci.org/corebreaker/texit)
[![Coverage Status](https://coveralls.io/repos/github/corebreaker/texit/badge.svg)](https://coveralls.io/github/corebreaker/texit)
[![GoDoc](https://godoc.org/github.com/corebreaker/texit?status.svg)](https://godoc.org/github.com/corebreaker/texit)
[![Release](https://img.shields.io/github/release/corebreaker/texit.svg?style=plastic)](https://github.com/corebreaker/texit/releases)


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
