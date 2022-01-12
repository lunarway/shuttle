package cmd

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRoot(t *testing.T) {
	strings := func(s ...string) []string {
		return s
	}
	tt := []struct {
		name   string
		input  []string
		output string
	}{
		{
			name:   "test moonbase build",
			input:  strings("-p", "../examples/no-plan-project", "run", "hello"),
			output: "Hello no plan project\n",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			rootCmd.SetOut(buf)
			rootCmd.SetErr(buf)
			rootCmd.SetArgs(tc.input)
			err := rootCmd.Execute()

			assert.Equal(t, tc.output, buf.String())
			assert.NoError(t, err)
		})
	}
}
