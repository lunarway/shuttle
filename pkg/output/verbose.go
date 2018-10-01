package output

import (
	"fmt"
)

// Verbose prints verbose output
func Verbose(verbose bool, msg string, args ...interface{}) {
	if verbose {
		fmt.Println(fmt.Sprintf(msg, args...))
	}
}
