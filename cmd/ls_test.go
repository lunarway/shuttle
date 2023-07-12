package cmd

import (
	"errors"
	"testing"
)

func TestLs(t *testing.T) {
	testCases := []testCase{
		{
			name:      "invalid shuttle.yaml file",
			input:     args("-p", "testdata/invalid-yaml", "ls"),
			stdoutput: "",
			erroutput: "Error: exit code 2 - Failed to parse shuttle configuration: yaml: unmarshal errors:\n  line 1: cannot unmarshal !!str `plan` into config.ShuttleConfig\n\nMake sure your 'shuttle.yaml' is valid.\n",
			err: errors.New(
				"exit code 2 - Failed to parse shuttle configuration: yaml: unmarshal errors:\n  line 1: cannot unmarshal !!str `plan` into config.ShuttleConfig\n\nMake sure your 'shuttle.yaml' is valid.",
			),
		},
		{
			name:      "list one action",
			input:     args("-p", "testdata/project", "ls"),
			stdoutput: "Available Scripts:\n  exit_0         \n  exit_1         \n  hello_stderr   \n  hello_stdout   \n  required_arg   \n",
			erroutput: "",
			err:       nil,
		},
	}
	executeTestCases(t, testCases)
}
