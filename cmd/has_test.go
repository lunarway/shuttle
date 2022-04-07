package cmd

import (
	"errors"
	"fmt"
	"testing"
)

func TestHas(t *testing.T) {
	testCases := []testCase{
		{
			name:      "bool variable",
			input:     args("-p", "testdata/project", "has", "boolVar"),
			stdoutput: "",
			erroutput: "",
			err:       nil,
		},
		{
			name:      "has wrong argument switch",
			input:     args("-j", "testdata/project", "has", "docker.baseImage"),
			stdoutput: "",
			erroutput: "",
			err:       fmt.Errorf("unknown shorthand flag: 'j' in -j"),
		},
		{
			name:      "stdout",
			input:     args("-p", "testdata/project", "has", "--stdout", "boolVar"),
			stdoutput: "true",
			erroutput: "",
			err:       nil,
		},
		{
			name:      "not existing",
			input:     args("-p", "testdata/project", "has", "unknown"),
			stdoutput: "",
			erroutput: "",
			err:       errors.New("exit code 1 - "),
		},
		{
			name:      "not existing stdout",
			input:     args("-p", "testdata/project", "has", "--stdout", "unknown"),
			stdoutput: "false",
			erroutput: "",
			err:       nil,
		},
		{
			name:      "script",
			input:     args("-p", "testdata/project", "has", "--script", "hello_stdout"),
			stdoutput: "",
			erroutput: "",
			err:       nil,
		},
		{
			name:      "script",
			input:     args("-p", "testdata/project", "has", "--script", "unknown"),
			stdoutput: "",
			erroutput: "",
			err:       errors.New("exit code 1 - "),
		},
	}
	executeTestCases(t, testCases)
}
