package cmd

import (
	"fmt"
	"testing"
)

func TestHas(t *testing.T) {
	strings := func(s ...string) []string {
		return s
	}

	testCases := []testCase{
		{
			name:      "has variable",
			input:     strings("-p", "../examples/repo-project", "has", "docker.baseImage"),
			stdoutput: "",
			erroutput: "",
			err:       nil,
		},
		{
			name:      "has wrong argument switch",
			input:     strings("-j", "../examples/repo-project", "has", "docker.baseImage"),
			stdoutput: "",
			erroutput: "",
			err:       fmt.Errorf("unknown shorthand flag: 'j' in -j"),
		},
	}
	executeTestCases(t, testCases)
}
