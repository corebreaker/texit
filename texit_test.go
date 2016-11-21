package texit

import (
    "fmt"
    "os"
    "testing"
)

func TestWithExit(t *testing.T) {
    MakeTestWithExit(f)

    stdout, stderr, status_code, err := MakeTestWithExit(func() {
        fmt.Println("Hello!")
        fmt.Fprintln(os.Stderr, "Salut !")
        os.Exit(0)
    })

    if err != nil {
        t.Fatal(err)
    }

    fmt.Println("C:", status_code)
    fmt.Println("O:", stdout)
    fmt.Println("E:", stderr)
}
