package output

import (
	"fmt"
	"os"
)

func ExitWithError(msg string) {
	ExitWithErrorCode(1, msg)
}

func ExitWithErrorCode(code int, msg string) {
	fmt.Printf("\x1b[31;1m%s\x1b[0m\n", "shuttle failed\n"+msg)
	os.Exit(code)
}

func CheckIfError(err error) {
	if err == nil {
		return
	}
	ExitWithError(fmt.Sprintf("error: %s", err))
}
