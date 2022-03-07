package cmd

import (
	"fmt"
	"testing"
)

func TestHas(t *testing.T) {
	testCases := []testCase{
		{
			name:      "has variable",
			input:     args("-p", "../examples/repo-project", "has", "docker.baseImage"),
			stdoutput: "",
			erroutput: "",
			err:       nil,
		},
		{
			name:      "has wrong argument switch",
			input:     args("-j", "../examples/repo-project", "has", "docker.baseImage"),
			stdoutput: "",
			erroutput: "",
			err:       fmt.Errorf("unknown shorthand flag: 'j' in -j"),
		},
	}
	executeTestCases(t, testCases)
}
