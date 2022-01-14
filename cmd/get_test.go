package cmd

import (
	"testing"
)

func TestGet(t *testing.T) {
	strings := func(s ...string) []string {
		return s
	}
	testCases := []testCase{
		{
			name:      "get variable",
			input:     strings("-p", "../examples/repo-project", "get", "docker.baseImage"),
			stdoutput: "golang",
			erroutput: "",
			err:       nil,
		},
	}
	executeTestCases(t, testCases)
}
