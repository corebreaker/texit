package texit

import (
    "errors"
    "runtime"
    "strings"
)

// Error message in the case of no testing function was found by `func_name` function
const _FUNCNAME_ERRMSG = "Test function has not found"

// Returns the name of testing function
// An error is returned if no testing function is found
// (this function call is not in the call graph of a testing function)
func func_name() (string, error) {
    // Program counter
    var pc uintptr = 1

    // Populate the stack trace starting to 1 cause 0 is this function (`func_name`)
    for i := 1; pc != 0; i++ {
        // Retreive runtime informations
        ptr, _, _, ok := runtime.Caller(i)
        pc = ptr

        // If there isn't significant information, go to next level
        if (pc == 0) || (!ok) {
            continue
        }

        // Retreive called function
        f := runtime.FuncForPC(pc)

        // Get the current functiont name
        name := f.Name()

        // Extract the short name by substracting the package name
        idx := strings.LastIndexByte(f.Name(), '.')
        if idx > 0 {
            name = name[(idx + 1):]
        }

        // If the prefix of the function is "Test"
        if strings.HasPrefix(name, "Test") {
            // So, it's a testing function
            return name, nil
        }
    }

    // Returns an error if no testing function was found
    return "", errors.New(_FUNCNAME_ERRMSG)
}
