package cmd

import (
	"errors"
	"testing"
)

func TestRun(t *testing.T) {
	strings := func(s ...string) []string {
		return s
	}
	testCases := []testCase{
		{
			name:      "std out echo",
			input:     strings("-p", "testdata/project", "run", "hello_stdout"),
			stdoutput: "Hello stdout\n",
			erroutput: "",
			err:       nil,
		},
		{
			name:      "std err echo",
			input:     strings("-p", "testdata/project", "run", "hello_stderr"),
			stdoutput: "",
			erroutput: "\x1b[31;1mHello stderr\x1b[0m\n",
			err:       nil,
		},
		{
			name:      "exit 0",
			input:     strings("-p", "testdata/project", "run", "exit_0"),
			stdoutput: "",
			erroutput: "",
			err:       nil,
		},
		{
			name:      "exit 1",
			input:     strings("-p", "testdata/project", "run", "exit_1"),
			stdoutput: "",
			erroutput: "Error: exit code 4 - Failed executing script `exit_1`: shell script `exit 1`\nExit code: 1\n",
			err:       errors.New("exit code 4 - Failed executing script `exit_1`: shell script `exit 1`\nExit code: 1"),
		},
	}
	executeTestCases(t, testCases)
}
