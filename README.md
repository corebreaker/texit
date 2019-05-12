# texit
[![Report](https://goreportcard.com/badge/github.com/corebreaker/texit?style=plastic)](https://goreportcard.com/report/github.com/corebreaker/gozlog)
[![Build Status](https://img.shields.io/travis/com/corebreaker/texit/master.svg?style=plastic)](https://travis-ci.com/corebreaker/texit)
[![Coverage Status](https://img.shields.io/coveralls/github/corebreaker/texit/master.svg?style=plastic)](https://coveralls.io/github/corebreaker/texit)
[![GoDoc](https://img.shields.io/badge/godoc-reference-5272B4.svg?style=plastic)](https://godoc.org/github.com/corebreaker/texit)
[![Release](https://img.shields.io/github/release/corebreaker/texit.svg?style=plastic)](https://github.com/corebreaker/texit/releases)
[![Libraries.io dependency status for GitHub repo](https://img.shields.io/librariesio/github/corebreaker/texit.svg?style=plastic)](https://libraries.io/github/corebreaker/texit)
[![GitHub license](https://img.shields.io/github/license/corebreaker/texit.svg?style=plastic)](https://github.com/corebreaker/texit/blob/master/LICENSE)

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
