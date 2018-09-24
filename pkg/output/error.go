package output

import (
	"fmt"
	"os"
)

func ExitWithError(msg string) {
	fmt.Printf("\x1b[31;1m%s\x1b[0m\n", msg)
	os.Exit(1)
}

func CheckIfError(err error) {
	if err == nil {
		return
	}
	ExitWithError(fmt.Sprintf("error: %s", err))
}
