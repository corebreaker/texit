package main

import (
    "log"
    "os"
    "sync"
)

var (
    atexitFuncs []func()
    exitStatus  = 0
    exitMu      sync.Mutex
)

func setExitStatus(n int) {
    exitMu.Lock()
    if exitStatus < n {
        exitStatus = n
    }
    exitMu.Unlock()
}

func atexit(f func()) {
    atexitFuncs = append(atexitFuncs, f)
}

func exit() {
    for _, f := range atexitFuncs {
        f()
    }

    os.Exit(exitStatus)
}

func fatalf(format string, args ...interface{}) {
    errorf(format, args...)
    exit()
}

func errorf(format string, args ...interface{}) {
    log.Printf(format, args...)
    setExitStatus(1)
}
