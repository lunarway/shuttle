package cmd

import (
	"testing"
)

func TestRoot(t *testing.T) {
	strings := func(s ...string) []string {
		return s
	}
	testCases := []testCase{
		{
			name:      "test moonbase build",
			input:     strings("-p", "../examples/no-plan-project", "run", "hello"),
			stdoutput: "Hello no plan project\n",
			erroutput: "",
			err:       nil,
		},
	}
	executeTestCases(t, testCases)
}
