package texit

import (
    "testing"
)

func TestFuncName(t *testing.T) {
    name, err := func_name()
    if err != nil {
        t.Fatal(err)
    }

    if name != "TestFuncName" {
        t.Fail()
    }

    type tMsg struct {
        name string
        err  error
    }

    resp := make(chan tMsg)

    go func() {
        name, err := func_name()
        resp <- tMsg{name, err}
    }()

    res := <-resp
    if res.err == nil {
        t.Errorf("No error returned, name=%s")
    }

    if res.err.Error() != _FUNCNAME_ERRMSG {
        t.Errorf("Bad error: %v", res.err)
    }
}
