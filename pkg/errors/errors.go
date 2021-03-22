package errors

import "fmt"

// ExitCode is an error indicating a specific exit code is used upon exit of
// shuttle.
type ExitCode struct {
	Code    int
	Message string
}

func (e *ExitCode) Error() string {
	return fmt.Sprintf("exit code %d - %s", e.Code, e.Message)
}

func NewExitCode(code int, format string, args ...interface{}) error {
	return &ExitCode{
		Code:    code,
		Message: fmt.Sprintf(format, args...),
	}
}
