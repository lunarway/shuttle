package cmd

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHasNoErr(t *testing.T) {
	strings := func(s ...string) []string {
		return s
	}
	tt := []struct {
		name   string
		input  []string
		output string
	}{
		{
			name:   "has variable",
			input:  strings("-p", "../examples/repo-project", "has", "docker.baseImage"),
			output: "",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			rootCmd.SetOut(buf)
			rootCmd.SetErr(buf)
			rootCmd.SetArgs(tc.input)
			err := rootCmd.Execute()

			assert.NoError(t, err)

			if tc.output != "" {
				assert.Equal(t, tc.output, buf.String())
			}

		})
	}
}
