package cmd

import (
	"testing"
)

func TestVersion(t *testing.T) {
	testCases := []testCase{
		{
			name:      "no args",
			input:     args("version"),
			stdoutput: "<dev-version>\n",
			erroutput: "",
			err:       nil,
		},
		{
			name:      "with commit",
			input:     args("version", "--commit"),
			stdoutput: "<unspecified-commit>\n",
			erroutput: "",
			err:       nil,
		},
	}
	executeTestCases(t, testCases)
}
