package main

import (
    "bytes"
    "go/build"
    "io/ioutil"
    "os"
    "os/exec"
    "path/filepath"
    "runtime"
    "strings"
    "sync"
)

var (
    goroot    = filepath.Clean(runtime.GOROOT())
    gobin     = os.Getenv("GOBIN")
    gorootBin = filepath.Join(goroot, "bin")
    gorootPkg = filepath.Join(goroot, "pkg")
    gorootSrc = filepath.Join(goroot, "src")
    toolDir   = build.ToolDir

    goarch    string
    goos      string
    exeSuffix string
    gopath    []string

    origEnv      []string
    buildContext = build.Default
)

func init() {
    goarch = buildContext.GOARCH
    goos = buildContext.GOOS

    if goos == "windows" {
        exeSuffix = ".exe"
    }

    gopath = filepath.SplitList(buildContext.GOPATH)
}

// envForDir returns a copy of the environment
// suitable for running in the given directory.
// The environment is the current process's environment
// but with an updated $PWD, so that an os.Getwd in the
// child will be faster.
func envForDir(dir string, base []string) []string {
    // Internally we only use rooted paths, so dir is rooted.
    // Even if dir is not rooted, no harm done.
    return mergeEnvLists([]string{"PWD=" + dir}, base)
}

// envList returns the value of the given environment variable broken
// into fields, using the default value when the variable is empty.
func envList(key, def string) []string {
    v := os.Getenv(key)
    if v == "" {
        v = def
    }
    return strings.Fields(v)
}

// mergeEnvLists merges the two environment lists such that
// variables with the same name in "in" replace those in "out".
// This always returns a newly allocated slice.
func mergeEnvLists(in, out []string) []string {
    out = append([]string(nil), out...)
NextVar:
    for _, inkv := range in {
        k := strings.SplitAfterN(inkv, "=", 2)[0]
        for i, outkv := range out {
            if strings.HasPrefix(outkv, k) {
                out[i] = inkv
                continue NextVar
            }
        }
        out = append(out, inkv)
    }
    return out
}

type envVar struct {
    name, value string
}

type builder struct {
    work      string          // the temporary work directory (ends in filepath.Separator)
    flagCache map[string]bool // a cache of supported compiler flags
    exec      sync.Mutex
}

func (b *builder) init() {
    workdir, err := ioutil.TempDir("", "go-build")
    if err != nil {
        fatalf("%s", err)
    }

    b.work = workdir
    atexit(func() { os.RemoveAll(workdir) })
}

// gccArchArgs returns arguments to pass to gcc based on the architecture.
func (b *builder) gccArchArgs() []string {
    switch goarch {
    case "386":
        return []string{"-m32"}
    case "amd64", "amd64p32":
        return []string{"-m64"}
    case "arm":
        return []string{"-marm"} // not thumb
    case "s390x":
        return []string{"-m64", "-march=z196"}
    case "mips64", "mips64le":
        return []string{"-mabi=64"}
    }
    return nil
}

// gccSupportsFlag checks to see if the compiler supports a flag.
func (b *builder) gccSupportsFlag(flag string) bool {
    b.exec.Lock()
    defer b.exec.Unlock()

    if b, ok := b.flagCache[flag]; ok {
        return b
    }

    if b.flagCache == nil {
        src := filepath.Join(b.work, "trivial.c")
        if err := ioutil.WriteFile(src, []byte{}, 0666); err != nil {
            return false
        }
        b.flagCache = make(map[string]bool)
    }

    cmdArgs := append(envList("CC", defaultCC), flag, "-c", "trivial.c")

    cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
    cmd.Dir = b.work
    cmd.Env = mergeEnvLists([]string{"LC_ALL=C"}, envForDir(cmd.Dir, os.Environ()))
    out, err := cmd.CombinedOutput()
    supported := err == nil && !bytes.Contains(out, []byte("unrecognized"))
    b.flagCache[flag] = supported

    return supported
}

func (b *builder) ccompilerCmd(envvar, defcmd, objdir string) []string {
    // NOTE: env.go's mkEnv knows that the first three
    // strings returned are "gcc", "-I", objdir (and cuts them off).

    compiler := envList(envvar, defcmd)
    a := []string{compiler[0], "-I", objdir}
    a = append(a, compiler[1:]...)

    // Definitely want -fPIC but on Windows gcc complains
    // "-fPIC ignored for target (all code is position independent)"
    if goos != "windows" {
        a = append(a, "-fPIC")
    }

    a = append(a, b.gccArchArgs()...)
    // gcc-4.5 and beyond require explicit "-pthread" flag
    // for multithreading with pthread library.
    if buildContext.CgoEnabled {
        switch goos {
        case "windows":
            a = append(a, "-mthreads")
        default:
            a = append(a, "-pthread")
        }
    }

    if strings.Contains(a[0], "clang") {
        // disable ASCII art in clang errors, if possible
        a = append(a, "-fno-caret-diagnostics")
        // clang is too smart about command-line arguments
        a = append(a, "-Qunused-arguments")
    }

    // disable word wrapping in error messages
    a = append(a, "-fmessage-length=0")

    // Tell gcc not to include the work directory in object files.
    if b.gccSupportsFlag("-fdebug-prefix-map=a=b") {
        a = append(a, "-fdebug-prefix-map="+b.work+"=/tmp/go-build")
    }

    // Tell gcc not to include flags in object files, which defeats the
    // point of -fdebug-prefix-map above.
    if b.gccSupportsFlag("-gno-record-gcc-switches") {
        a = append(a, "-gno-record-gcc-switches")
    }

    // On OS X, some of the compilers behave as if -fno-common
    // is always set, and the Mach-O linker in 6l/8l assumes this.
    // See https://golang.org/issue/3253.
    if goos == "darwin" {
        a = append(a, "-fno-common")
    }

    return a
}

// gccCmd returns a gcc command line prefix
// defaultCC is defined in zdefaultcc.go, written by cmd/dist.
func (b *builder) gccCmd(objdir string) []string {
    return b.ccompilerCmd("CC", defaultCC, objdir)
}

// gxxCmd returns a g++ command line prefix
// defaultCXX is defined in zdefaultcc.go, written by cmd/dist.
func (b *builder) gxxCmd(objdir string) []string {
    return b.ccompilerCmd("CXX", defaultCXX, objdir)
}

func mkEnv() []envVar {
    var b builder

    b.init()

    env := []envVar{
        {"GOARCH", goarch},
        {"GOBIN", gobin},
        {"GOEXE", exeSuffix},
        {"GOHOSTARCH", runtime.GOARCH},
        {"GOHOSTOS", runtime.GOOS},
        {"GOOS", goos},
        {"GOPATH", os.Getenv("GOPATH")},
        {"GORACE", os.Getenv("GORACE")},
        {"GOROOT", goroot},
        {"GOTOOLDIR", toolDir},

        // disable escape codes in clang errors
        {"TERM", "dumb"},
    }

    if goos != "plan9" {
        cmd := b.gccCmd(".")
        env = append(env, envVar{"CC", cmd[0]})
        env = append(env, envVar{"GOGCCFLAGS", strings.Join(cmd[3:], " ")})
        cmd = b.gxxCmd(".")
        env = append(env, envVar{"CXX", cmd[0]})
    }

    if buildContext.CgoEnabled {
        env = append(env, envVar{"CGO_ENABLED", "1"})
    } else {
        env = append(env, envVar{"CGO_ENABLED", "0"})
    }

    return env
}
