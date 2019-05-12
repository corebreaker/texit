package texit

import (
	"testing"
)

func TestFuncName(t *testing.T) {
	name, line, err := func() (string, int, error) {
		return testFuncName()
	}()

	if err != nil {
		t.Fatal(err)
	}

	const funcName = "TestFuncName"
	const lineNum = 10

	if name != funcName {
		t.Fatalf("Bad function name (%s != %s)", name, funcName)
	}

	if line != lineNum {
		t.Fatalf("Bad file line (%d != %d)", line, lineNum)
	}

	type tMsg struct {
		name string
		line int
		err  error
	}

	resp := make(chan tMsg)

	go func() {
		name, line, err := testFuncName()
		resp <- tMsg{name, line, err}
	}()

	res := <-resp
	if res.err == nil {
		t.Errorf("No error returned, name=%s:%d", res.name, res.line)
	}

	if res.err != ErrTestFuncNameNotFound {
		t.Errorf("Bad error: %v", res.err)
	}
}
