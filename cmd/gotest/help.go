package main

import "strings"

const testUsage = "test [build/test flags] [packages] [build/test flags & test binary flags]"

const testFlag1 = `
In addition to the build flags, the flags handled by 'go test' itself are:

	-args
	    Pass the remainder of the command line (everything after -args)
	    to the test binary, uninterpreted and unchanged.
	    Because this flag consumes the remainder of the command line,
	    the package list (if present) must appear before this flag.

	-c
	    Compile the test binary to pkg.test but do not run it
	    (where pkg is the last element of the package's import path).
	    The file name can be changed with the -o flag.

	-exec xprog
	    Run the test binary using xprog. The behavior is the same as
	    in 'go run'. See 'go help run' for details.

	-i
	    Install packages that are dependencies of the test.
	    Do not run the test.

	-o file
	    Compile the test binary to the named file.
	    The test still runs (unless -c or -i is specified).

The test binary also accepts flags that control execution of the test; these
flags are also accessible by 'go test'.
`

const testFlag2 = `
	-bench regexp
	    Run (sub)benchmarks matching a regular expression.
	    The given regular expression is split into smaller ones by
	    top-level '/', where each must match the corresponding part of a
	    benchmark's identifier.
	    By default, no benchmarks run. To run all benchmarks,
	    use '-bench .' or '-bench=.'.

	-benchmem
	    Print memory allocation statistics for benchmarks.

	-benchtime t
	    Run enough iterations of each benchmark to take t, specified
	    as a time.Duration (for example, -benchtime 1h30s).
	    The default is 1 second (1s).

	-blockprofile block.out
	    Write a goroutine blocking profile to the specified file
	    when all tests are complete.
	    Writes test binary as -c would.

	-blockprofilerate n
	    Control the detail provided in goroutine blocking profiles by
	    calling runtime.SetBlockProfileRate with n.
	    See 'go doc runtime.SetBlockProfileRate'.
	    The profiler aims to sample, on average, one blocking event every
	    n nanoseconds the program spends blocked.  By default,
	    if -test.blockprofile is set without this flag, all blocking events
	    are recorded, equivalent to -test.blockprofilerate=1.

	-count n
	    Run each test and benchmark n times (default 1).
	    If -cpu is set, run n times for each GOMAXPROCS value.
	    Examples are always run once.

	-cover
	    Enable coverage analysis.

	-covermode set,count,atomic
	    Set the mode for coverage analysis for the package[s]
	    being tested. The default is "set" unless -race is enabled,
	    in which case it is "atomic".
	    The values:
		set: bool: does this statement run?
		count: int: how many times does this statement run?
		atomic: int: count, but correct in multithreaded tests;
			significantly more expensive.
	    Sets -cover.

	-coverpkg pkg1,pkg2,pkg3
	    Apply coverage analysis in each test to the given list of packages.
	    The default is for each test to analyze only the package being tested.
	    Packages are specified as import paths.
	    Sets -cover.

	-coverprofile cover.out
	    Write a coverage profile to the file after all tests have passed.
	    Sets -cover.

	-cpu 1,2,4
	    Specify a list of GOMAXPROCS values for which the tests or
	    benchmarks should be executed.  The default is the current value
	    of GOMAXPROCS.

	-cpuprofile cpu.out
	    Write a CPU profile to the specified file before exiting.
	    Writes test binary as -c would.

	-memprofile mem.out
	    Write a memory profile to the file after all tests have passed.
	    Writes test binary as -c would.

	-memprofilerate n
	    Enable more precise (and expensive) memory profiles by setting
	    runtime.MemProfileRate.  See 'go doc runtime.MemProfileRate'.
	    To profile all memory allocations, use -test.memprofilerate=1
	    and pass --alloc_space flag to the pprof tool.

	-outputdir directory
	    Place output files from profiling in the specified directory,
	    by default the directory in which "go test" is running.

	-parallel n
	    Allow parallel execution of test functions that call t.Parallel.
	    The value of this flag is the maximum number of tests to run
	    simultaneously; by default, it is set to the value of GOMAXPROCS.
	    Note that -parallel only applies within a single test binary.
	    The 'go test' command may run tests for different packages
	    in parallel as well, according to the setting of the -p flag
	    (see 'go help build').

	-run regexp
	    Run only those tests and examples matching the regular expression.
	    For tests the regular expression is split into smaller ones by
	    top-level '/', where each must match the corresponding part of a
	    test's identifier.

	-short
	    Tell long-running tests to shorten their run time.
	    It is off by default but set during all.bash so that installing
	    the Go tree can run a sanity check but not spend time running
	    exhaustive tests.

	-timeout t
	    If a test runs longer than t, panic.
	    The default is 10 minutes (10m).

	-trace trace.out
	    Write an execution trace to the specified file before exiting.

	-v
	    Verbose output: log all tests as they are run. Also print all
	    text from Log and Logf calls even if the test succeeds.

Each of these flags is also recognized with an optional 'test.' prefix,
as in -test.v. When invoking the generated test binary (the result of
'go test -c') directly, however, the prefix is mandatory.

The 'go test' command rewrites or removes recognized flags,
as appropriate, both before and after the optional package list,
before invoking the test binary.

For instance, the command

	go test -v -myflag testdata -cpuprofile=prof.out -x

will compile the test binary and then run it as

	pkg.test -test.v -myflag testdata -test.cpuprofile=prof.out

(The -x flag is removed because it applies only to the go command's
execution, not to the test itself.)

The test flags that generate profiles (other than for coverage) also
leave the test binary in pkg.test for use when analyzing the profiles.

When 'go test' runs a test binary, it does so from within the
corresponding package's source code directory. Depending on the test,
it may be necessary to do the same when invoking a generated test
binary directly.

The command-line package list, if present, must appear before any
flag not known to the go test command. Continuing the example above,
the package list would have to appear before -myflag, but could appear
on either side of -v.

To keep an argument for a test binary from being interpreted as a
known flag or a package name, use -args (see 'go help test') which
passes the remainder of the command line through to the test binary
uninterpreted and unaltered.

For instance, the command

	go test -v -args -x -v

will compile the test binary and then run it as

	pkg.test -test.v -x -v

Similarly,

	go test -args math

will compile the test binary and then run it as

	pkg.test math

In the first example, the -x and the second -v are passed through to the
test binary unchanged and with no effect on the go command itself.
In the second example, the argument math is passed through to the test
binary, instead of being interpreted as the package list.
`

var cmdHelp = `
'Go test' automates testing the packages named by the import paths.
It prints a summary of the test results in the format:

	ok   archive/tar   0.011s
	FAIL archive/zip   0.022s
	ok   compress/gzip 0.033s
	...

followed by detailed output for each failed package.

'Go test' recompiles each package along with any files with names matching
the file pattern "*_test.go".
Files whose names begin with "_" (including "_test.go") or "." are ignored.
These additional files can contain test functions, benchmark functions, and
example functions.  See 'go help testfunc' for more.
Each listed package causes the execution of a separate test binary.

Test files that declare a package with the suffix "_test" will be compiled as a
separate package, and then linked and run with the main test binary.

The go tool will ignore a directory named "testdata", making it available
to hold ancillary data needed by the tests.

By default, go test needs no arguments.  It compiles and tests the package
with source in the current directory, including tests, and runs the tests.

The package is built in a temporary directory so it does not interfere with the
non-test installation.

` + strings.TrimSpace(testFlag1) + ` See 'go help testflag' for details.

For more about build flags, see 'go help build'.
For more about specifying packages, see 'go help packages'.

See also: go build, go vet.
`
