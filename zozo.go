package texit

import (
    "fmt"
    "os"
    "path/filepath"

    "github.com/siddontang/go/ioutil2"
    "github.com/termie/go-shutil"
)

func f() {
    p := "C:\\Users\\dev\\go\\src\\github.com\\corebreaker\\texit\\test.sav"
    if !ioutil2.FileExists(p) {
        src := filepath.Dir(os.Args[0])
        err := shutil.CopyTree(src, p, nil)
        if err != nil {
            panic(err)
        }
    }

    fmt.Println("Args:", os.Args)
}
