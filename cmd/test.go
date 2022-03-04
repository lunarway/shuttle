package cmd

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testCase struct {
	name      string
	input     []string
	stdoutput string
	erroutput string
	err       error
}

func executeTestCases(t *testing.T, testCases []testCase) {
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			stdBuf := new(bytes.Buffer)
			errBuf := new(bytes.Buffer)

			rootCmd, _ := initializedRoot(stdBuf, errBuf)
			rootCmd.SetArgs(tc.input)

			err := rootCmd.Execute()
			if tc.err == nil {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, tc.err, err.Error())
			}
			assert.Equal(t, tc.stdoutput, stdBuf.String())
			assert.Equal(t, tc.erroutput, errBuf.String())
			t.Logf("STDOUT: %s", stdBuf.String())
			t.Logf("ERROUT: %s", errBuf.String())
		})
	}
}
