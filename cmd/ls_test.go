package cmd

import (
	"testing"
)

func TestLs(t *testing.T) {
	strings := func(s ...string) []string {
		return s
	}
	testCases := []testCase{
		{
			name:      "list one action",
			input:     strings("-p", "../examples/no-plan-project", "ls"),
			stdoutput: "Available Scripts:\n  hello        \n",
			erroutput: "",
			err:       nil,
		},
		{
			name:      "list actions",
			input:     strings("-p", "../examples/repo-project/", "ls"),
			stdoutput: "Pulling latest plan changes on master\nAvailable Scripts:\n  build        Build the docker image\n  deploy       Deploys the image to a kubernetes environment\n  push         Push the docker image\n  say          Say something\n  test         Run test for the project\n",
			erroutput: "",
			err:       nil,
		},
	}
	executeTestCases(t, testCases)
}
