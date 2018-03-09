package main

import (
    "flag"
    "fmt"
    "log"
    "os"
    "path/filepath"
    "runtime"
    "strings"
)

func usage() {
    os.Stderr.WriteString(testUsage + "\n\n" +
        strings.TrimSpace(testFlag1) + "\n\n\t" +
        strings.TrimSpace(testFlag2) + "\n")
    os.Exit(2)
}

func main() {
    flag.Usage = usage
    flag.Parse()
    log.SetFlags(0)

    args := flag.Args()
    if len(args) < 1 {
        usage()
    }

    // Diagnose common mistake: GOPATH==GOROOT.
    // This setting is equivalent to not setting GOPATH at all,
    // which is not what most people want when they do it.
    if gopath := os.Getenv("GOPATH"); gopath == runtime.GOROOT() {
        fmt.Fprintf(os.Stderr, "warning: GOPATH set to GOROOT (%s) has no effect\n", gopath)
    } else {
        for _, p := range filepath.SplitList(gopath) {
            // Note: using HasPrefix instead of Contains because a ~ can appear
            // in the middle of directory elements, such as /tmp/git-1.8.2~rc3
            // or C:\PROGRA~1. Only ~ as a path prefix has meaning to the shell.
            if strings.HasPrefix(p, "~") {
                fmt.Fprintf(os.Stderr, "go: GOPATH entry cannot start with shell metacharacter '~': %q\n", p)
                os.Exit(2)
            }
            if !filepath.IsAbs(p) {
                fmt.Fprintf(os.Stderr, "go: GOPATH entry is relative; must be absolute path: %q.\nRun 'go help gopath' for usage.\n", p)
                os.Exit(2)
            }
        }
    }

    if fi, err := os.Stat(goroot); err != nil || !fi.IsDir() {
        fmt.Fprintf(os.Stderr, "go: cannot find GOROOT directory: %v\n", goroot)
        os.Exit(2)
    }

    // Set environment (GOOS, GOARCH, etc) explicitly.
    // In theory all the commands we invoke should have
    // the same default computation of these as we do,
    // but in practice there might be skew
    // This makes sure we all agree.
    origEnv = os.Environ()
    for _, env := range mkEnv() {
        if os.Getenv(env.name) != env.value {
            os.Setenv(env.name, env.value)
        }
    }

    defer exit()

    var flg flag.FlagSet

    flg.Usage = func() { cmd.Usage() }

    /*
       cmd.Flag.Usage = func() { cmd.Usage() }
       if cmd.CustomFlags {
           args = args[1:]
       } else {
           cmd.Flag.Parse(args[1:])
           args = cmd.Flag.Args()
       }
       cmd.Run(cmd, args)
       exit()
       return
    */
    if (len(os.Args) > 1) && (os.Args[1] == "ok") {
        fmt.Println("Toto")
        return
    }

    fmt.Fprintf(os.Stderr, "go: unknown subcommand %q\nRun 'go help' for usage.\n", args[0])
    setExitStatus(2)
}
