package texit

import (
	"testing"
)

func TestFuncName(t *testing.T) {
	name, line, err := func() (string, int, error) {
		return test_func_name()
	}()

	if err != nil {
		t.Fatal(err)
	}

	const func_name = "TestFuncName"
	const line_num = 10

	if name != func_name {
		t.Fatalf("Bad function name (%s != %s)", name, func_name)
	}

	if line != line_num {
		t.Fatalf("Bad file line (%d != %d)", line, line_num)
	}

	type tMsg struct {
		name string
		line int
		err  error
	}

	resp := make(chan tMsg)

	go func() {
		name, line, err := test_func_name()
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
