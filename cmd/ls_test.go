package cmd

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLs(t *testing.T) {
	strings := func(s ...string) []string {
		return s
	}
	tt := []struct {
		name   string
		input  []string
		output string
	}{
		{
			name:   "list one action",
			input:  strings("-p", "../examples/no-plan-project", "ls"),
			output: "Available Scripts:\n  hello        \n",
		},
		{
			name:   "list actions",
			input:  strings("-p", "../examples/repo-project/", "ls"),
			output: "Pulling latest plan changes on master\nAvailable Scripts:\n  build        Build the docker image\n  deploy       Deploys the image to a kubernetes environment\n  push         Push the docker image\n  say          Say something\n  test         Run test for the project\n",
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
